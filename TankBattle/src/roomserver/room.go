package main

import (
	"errors"
	"time"

	"github.com/golang/glog"
)

const (
	MaxPlayerNum uint32 = 4
)

//提供信息给roommgr管理房间
type Room struct {
	//mutex    sync.RWMutex
	id       uint32       //房间id
	roomtype uint32       //房间类型
	players  []PlayerTask //房间内的玩家
	curnum   uint32       //当前房间内玩家数
	isstart  bool
}

//返回给客户端的房间信息
type room struct {
	port uint32
}

//给一个玩家分配房间（已经加入房间）
func (this *Room) AddPlayer(player PlayerTask) error {
	//this.mutex.Lock()
	if this.curnum >= MaxPlayerNum {
		glog.Error("[Room] 房间已满")
		return errors.New("room is full")
	}
	this.curnum++
	this.players = append(this.players, player)
	//this.mutex.Unlock()
	return nil
}

func NewRoom(rtype, rid uint32) *Room {
	room := &Room{
		id:       rid,
		roomtype: rtype,
	}
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
	for i := 0; i < 30; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		glog.Info("game is running")
	}
	RoomMgr_GetMe().endchan <- this.id
}
