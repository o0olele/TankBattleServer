package main

import (
	"base/env"
	common "common"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
)

const (
	MaxPlayerNum uint32 = 2
)

//提供信息给roommgr管理房间
type Room struct {
	//mutex    sync.RWMutex
	id          uint32                 //房间id
	roomtype    uint32                 //房间类型
	players     map[uint32]*PlayerTask //房间内的玩家
	curnum      uint32                 //当前房间内玩家数
	isstart     bool
	timeloop    uint64
	stopch      chan bool
	Isstop      bool
	totgametime uint64 //in second
	allbullet   sync.Map
	//allbullet   map[uint32]*common.Bullet
	endchan     chan bool
	bulletcount uint32
}

//返回给客户端的房间信息
type room struct {
	port uint32
}

//给一个玩家分配房间（已经加入房间）
func (this *Room) AddPlayer(player *PlayerTask) error {
	//this.mutex.Lock()
	if this.checkPlayer(player) {
		glog.Info("[Room] ", player.id, "玩家已经在[", this.id, "]房间里面了")
		return nil
	}
	if this.curnum >= MaxPlayerNum {
		glog.Error("[Room] 房间已满")
		return errors.New("room is full")
	}
	this.curnum++
	this.players[player.id] = player
	this.players[player.id].room = this
	//this.mutex.Unlock()
	return nil
}

func NewRoom(rtype, rid uint32) *Room {
	room := &Room{
		id:       rid,
		roomtype: rtype,
		players:  make(map[uint32]*PlayerTask),
		curnum:   0,
		isstart:  false,
		Isstop:   false,
		endchan:  make(chan bool),
		//allbullet:   make(map[uint32]*common.Bullet),
		bulletcount: 1,
	}
	room.totgametime, _ = strconv.ParseUint(env.Get("room", "time"), 10, 64)
	return room
}

func (this *Room) IsFull() bool {
	if this.curnum < MaxPlayerNum {
		return false
	}
	return true
}

func (this *Room) Start() {
	this.isstart = true
	this.GameLoop()
}

func (this *Room) GameLoop() {
	timeTicker := time.NewTicker(time.Millisecond * 10)
	stop := false
	for !stop {
		// SceneMsg, 用于同步场景，在球球里面是100ms一次，在该游戏中，SceneMsg可以包含
		// 以下信息：当前视野、当前视野内的玩家位置（包括自己）、当前视野内的子弹位置

		// 对于单次处理时间内玩家操作过多的问题，顺序处理即可，如果出现玩家卡顿的现象，再考虑优化
		// 例如给玩家的操作添加帧号，但是这个是可选项，不是必要的

		// 另外，服务器的逻辑帧率通常要高于客户端的渲染帧率，例如客户端可能是30，但是服务器可能就是
		// 60甚至更高
		select {
		case <-timeTicker.C:
			if this.timeloop%2 == 0 { //0.02s
				this.update()

			}
			if this.timeloop%2 == 0 { //0.1s
				this.sendRoomMsg()
			}
			if this.timeloop%100 == 0 { //1s
				this.sendTime(this.totgametime - this.timeloop/100)
			}
			if this.timeloop != 0 && this.timeloop%(this.totgametime*100) == 0 {
				stop = true
			}
			this.timeloop++
			if this.Isstop {
				stop = true
			}
		}
	}
	this.Close()
}
func (this *Room) Close() {
	if !this.Isstop {
		this.Isstop = true
		RoomMgr_GetMe().endchan <- this.id
	}
}

func (this *Room) checkPlayer(player *PlayerTask) bool {
	if _, ok := this.players[player.id]; !ok {
		return false
	}
	return true
}

func (this *Room) sendRoomMsg() {
	for _, p := range this.players {
		p.SendSceneMsg()
	}
}

func (this *Room) sendTime(t uint64) {
	for _, p := range this.players {
		t := common.RetTimeMsg{
			Time: t,
		}
		jstr, err := json.Marshal(t)
		if err != nil {
			glog.Error("[Time] marshal jsonMsg err")
			return
		}
		fmt.Println(string(jstr))
		p.wstask.AsyncSend(jstr, 0)
	}
}

func (this *Room) update() {
	for _, p := range this.players {
		p.Update()
	}

	for _, p := range this.players {
		p.UpdateOthers()
	}

	this.allbullet.Range(func(k, v interface{}) bool {
		b, ok := v.(*common.Bullet)
		if !ok {
			return false
		}
		updateBulletPos(b, this.players)
		return true
	})

}
