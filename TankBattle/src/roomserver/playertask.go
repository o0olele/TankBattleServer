package main

import (
	"base/gonet"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"runtime/debug"
	"sync"
	"time"

	common "common"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type PlayerTask struct {
	id         uint32
	wstask     *gonet.WebSocketTask
	self       *ScenePlayer
	key        string
	activetime time.Time
	room       *Room
}

func NewPlayerTask(conn *websocket.Conn) *PlayerTask {
	m := &PlayerTask{
		wstask:     gonet.NewWebSocketTask(conn),
		activetime: time.Now(),
	}
	m.wstask.Derived = m
	return m
}

func (this *PlayerTask) Start() {
	this.wstask.Start()
}

func (this *PlayerTask) Stop() bool {
	return this.wstask.Stop()
}

func (this *PlayerTask) OnClose() {
	this.wstask.Close()

	PlayerTaskMgr_GetMe().Del(this)

	this.room = nil
}

func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {
	glog.Info("[WS] Parse Msg ", data)
	this.activetime = time.Now()

	msgtype := common.MsgType(uint16(data[2]) | uint16(data[3])<<8)
	switch msgtype {
	case common.MsgType_Token:
		glog.Info("[room] Recv Token Msg ", data[4:])
		var id uint32
		err := binary.Read(bytes.NewReader(data[4:8]), binary.LittleEndian, &id)
		if nil != err {
			glog.Error("[WS] Endian Trans Fail")
			return false
		}
		this.id = id

		var token *common.Token
		token, err = common.DecryptTokenSSL(string(data[8:]))
		if nil != err {
			glog.Error("[room] decrypt openssl token fail ", err)
			return false
		}
		glog.Info("[room] decrypt openssl token ", token.Id, token.Time, time.Now().Unix())

		if time.Now().Unix()-token.Time > 30 {
			glog.Error("[room] player verify fail")
			return false
		}

		this.wstask.Verify()
		PlayerTaskMgr_GetMe().Add(this)
		RoomMgr_GetMe().GetRoom(this)
	case common.MsgType_Move:
		var angle, power uint32
		err := binary.Read(bytes.NewReader(data[4:]), binary.LittleEndian, &angle)
		err = binary.Read(bytes.NewReader(data[8:]), binary.LittleEndian, &power)
		if nil != err {
			glog.Error("[WS] Endian Trans Fail")
			return false
		}
		glog.Info("[WS] Parse Msg Move ", angle)

		if nil == this.room {
			return false
		}
		if this.room.Isstop {
			return false
		}
		req := common.ReqMoveMsg{
			Userid: this.id,
			Direct: angle,
			Power:  power,
		}
		this.room.opChan <- &opMsg{op: common.PlayerMove, args: req}

	case common.MsgType_Shoot:
		var angle uint32
		err := binary.Read(bytes.NewReader(data[4:]), binary.LittleEndian, &angle)
		if nil != err {
			glog.Error("[WS] Endian Trans Fail")
			return false
		}
		glog.Info("[WS] Parse Msg shoot ", angle)
		req := common.ReqShootMsg{
			Userid: this.id,
			Direct: angle,
		}
		this.room.opChan <- &opMsg{op: common.AddBullet, args: req}
		//this.scene.addBullet(this.direct)

	case common.MsgType_Finsh:
		this.room.Close()
	case common.MsgType_Heart:
		//this.wstask.AsyncSend(data, flag)
	case common.MsgType_Direct:
		var angle uint32
		err := binary.Read(bytes.NewReader(data[4:8]), binary.LittleEndian, &angle)
		if nil != err {
			glog.Error("[WS] Endian Trans Fail")
			return false
		}
		glog.Info("[WS] Parse Msg Turn ", angle)

		if nil == this.room {
			return false
		}
		if this.room.Isstop {
			return false
		}
		req := common.ReqTurnMsg{
			Userid: this.id,
			Direct: angle,
		}
		this.room.opChan <- &opMsg{op: common.PlayerTurn, args: req}
	default:
	}
	return true
}

func (this *PlayerTask) SendOverMsg() {
	bytes, _ := json.Marshal(common.RetOverMsg{End: true})
	this.wstask.AsyncSend(bytes, 0)
}

func (this *PlayerTask) SendSceneMsg(msg *common.RetSceneMsg) bool {
	// if nil == this.scene {
	// 	return false
	// }

	// msg := this.scene.SceneMsg()
	// if nil == msg {
	// 	glog.Error("[Scene] Msg Nil")
	// 	return false
	// }

	buf, _ := json.Marshal(*msg)
	//fmt.Println(string(buf))
	return this.wstask.AsyncSend(buf, 0)
}

type PlayerTaskMgr struct {
	mutex sync.RWMutex
	tasks map[uint32]*PlayerTask
}

var mPlayerTaskMgr *PlayerTaskMgr

func PlayerTaskMgr_GetMe() *PlayerTaskMgr {
	if nil == mPlayerTaskMgr {
		mPlayerTaskMgr = &PlayerTaskMgr{
			tasks: make(map[uint32]*PlayerTask),
		}
		go mPlayerTaskMgr.iTimeAction()
	}

	return mPlayerTaskMgr
}

func (this *PlayerTaskMgr) iTimeAction() {
	var (
		timeTicker = time.NewTicker(time.Second)
		loop       uint64
		ptasks     []*PlayerTask
	)
	defer func() {
		timeTicker.Stop()
		if err := recover(); nil != err {
			glog.Error("[Unexpeted] ", err, "\n", string(debug.Stack()))
		}
	}()

	for {
		select {
		case <-timeTicker.C:
			if 0 == loop%5 {
				now := time.Now()

				this.mutex.RLock()
				for _, t := range this.tasks {
					if now.Sub(t.activetime) > common.Task_TimeOut*time.Second {
						ptasks = append(ptasks, t)
					}
				}
				this.mutex.RUnlock()

				for _, t := range ptasks {
					if !t.Stop() {
						this.Del(t)
					}
					glog.Info("[Player] Connection timeout, player id=", t.id)
				}
				ptasks = ptasks[:0]
			}
			loop += 1
		}
	}
}

func (this *PlayerTaskMgr) Add(t *PlayerTask) bool {
	if nil == t {
		glog.Error("[WS] Player Task Manager Add Fail, Nil")
		return false
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.tasks[t.id] = t

	return true
}

func (this *PlayerTaskMgr) Del(t *PlayerTask) bool {
	if nil == t {
		glog.Error("[WS] Player Task Manager Del Fail, Nil")
		return false
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	_t, ok := this.tasks[t.id]
	if !ok {
		return false
	}
	if t != _t {
		glog.Error("[WS] Player Task Manager Del Fail, ", t.id, ",", &t, ",", &_t)
		return false
	}

	delete(this.tasks, t.id)

	return true
}

func (this *PlayerTaskMgr) Get(id uint32) *PlayerTask {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	t, ok := this.tasks[id]
	if !ok {
		return nil
	}

	return t
}
