package main

import (
	"base/gonet"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"

	common "common"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type PlayerTask struct {
	wstask *gonet.WebSocketTask

	key        string
	playerInfo *PlayerInfo
	activetime time.Time
	room       *Room
}

func NewPlayerTask(conn *websocket.Conn) *PlayerTask {
	m := &PlayerTask{
		wstask:     gonet.NewWebSocketTask(conn),
		activetime: time.Now(),
		playerInfo: &PlayerInfo{HP: 100},
	}
	m.wstask.Derived = m
	return m
}

func (this *PlayerTask) Start() {
	this.playerInfo.id = rand.New(rand.NewSource(time.Now().UnixNano())).Uint32() % 100 // 待优化

	fmt.Println("new playertask", this.playerInfo)
	this.wstask.Start()
	this.wstask.Verify() // 待优化
	PlayerTaskMgr_GetMe().Add(this)
	RoomMgr_GetMe().GetRoom(this)
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

	case common.MsgType_Move:

		var angle uint32
		err := binary.Read(bytes.NewReader(data[4:]), binary.LittleEndian, &angle)
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
			Userid: this.playerInfo.id,
			Direct: angle,
		}
		this.room.opChan <- &opMsg{op: common.PlayerMove, args: req}

	case common.MsgType_Shoot:

		//this.scene.addBullet(this.direct)

	case common.MsgType_Finsh:
		this.room.Close()
	case common.MsgType_Heart:
		this.wstask.AsyncSend(data, flag)
	case common.MsgType_Direct:
		var angle, power uint32
		err := binary.Read(bytes.NewReader(data[4:8]), binary.LittleEndian, &angle)
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
			Userid: this.playerInfo.id,
			Direct: angle,
			Power:  power,
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
					glog.Info("[Player] Connection timeout, player id=", t.playerInfo.id)
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

	this.tasks[t.playerInfo.id] = t

	return true
}

func (this *PlayerTaskMgr) Del(t *PlayerTask) bool {
	if nil == t {
		glog.Error("[WS] Player Task Manager Del Fail, Nil")
		return false
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	_t, ok := this.tasks[t.playerInfo.id]
	if !ok {
		return false
	}
	if t != _t {
		glog.Error("[WS] Player Task Manager Del Fail, ", t.playerInfo.id, ",", &t, ",", &_t)
		return false
	}

	delete(this.tasks, t.playerInfo.id)

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
