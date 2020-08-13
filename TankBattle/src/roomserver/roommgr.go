package main

import (
	"common"
	"container/list"
	"errors"
	"sync"

	"github.com/golang/glog"
)

type RoomMgr struct {
	mutex      sync.Mutex
	unFullRoom Queue
	runRoom    map[uint32]*Room
	Load       uint32 //roomnum
	endchan    chan uint32
	rcount     uint32 //上一个房间id
}

var mroommgr *RoomMgr

func RoomMgr_GetMe() *RoomMgr {
	if mroommgr == nil {
		mroommgr = &RoomMgr{
			unFullRoom: Queue{queue: list.New()},
		}
	}
	return mroommgr
}
func (this *RoomMgr) start() {
	for {
		select {
		case rid := <-this.endchan:
			this.mutex.Lock()
			defer this.mutex.Unlock()
			this.runRoom[rid] = nil
			this.Load--
			delete(this.runRoom, rid)
		}
	}
}

//给玩家分配可用房间
func (this *RoomMgr) GetRoom(player PlayerTask) (*Room, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	room, ok := this.unFullRoom.Front().(*Room)

	if !ok && room == nil {
		glog.Info("当前没有空房，创建新房")
		rid := this.getNextRoomid()
		room = NewRoom(common.CommonRoom, rid)
		if !room.IsFull() {
			this.unFullRoom.Push(room)
		}
		this.Load++
	}

	err := room.AddPlayer(player)
	if err != nil {
		glog.Error("为玩家分配房间失败", room, player)
		return nil, errors.New("distribute room error")
	}
	glog.Info("为玩家", player, "分配房间", room)
	if room.IsFull() {
		room.Start()
		this.runRoom[room.id] = room
		this.unFullRoom.Pop()
	}
	return room, nil
}

func (this *RoomMgr) getNextRoomid() uint32 {
	this.rcount = (this.rcount + 1) % 10000
	return this.rcount
}

func (this *RoomMgr) GetLoad() uint32 {
	return this.Load
}
