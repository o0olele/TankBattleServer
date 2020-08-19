package main

import (
	"base/gonet"
	"bytes"
	"encoding/binary"
	"math/rand"
	"time"

	common "common"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type PlayerTask struct {
	wstask *gonet.WebSocketTask

	key    string
	id     uint32
	name   string
	room   *Room
	direct uint32

	scene *Scene
}

func NewPlayerTask(conn *websocket.Conn) *PlayerTask {
	m := &PlayerTask{
		wstask: gonet.NewWebSocketTask(conn),
		scene: &Scene{
			speed:   common.SceneSpeed,
			hasMove: true,
		},
	}
	m.wstask.Derived = m

	return m
}

func (this *PlayerTask) Start() {
	this.id = rand.New(rand.NewSource(time.Now().UnixNano())).Uint32() % 100
	this.wstask.Start()
	this.wstask.Verify()

	room, err := RoomMgr_GetMe().GetRoom(this)
	if nil != err {
		glog.Error("[roomserver] Allocate room fail ", err)
		return
	}

	this.scene.room = room
}

func (this *PlayerTask) OnClose() {
	this.wstask.Close()
}

func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {
	glog.Info("[WS] Parse Msg ", data)

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
		this.direct = angle
		this.scene.UpdateSpeed(common.SceneSpeed)
	default:
	}
	return true
}

func (this *PlayerTask) SendSceneMsg() bool {
	msg := this.scene.SceneMsg()
	if nil == msg {
		glog.Error("[Scene] Msg Nil")
		return false
	}

	return this.wstask.AsyncSend(msg, 0)
}

func (this *PlayerTask) Update() {
	this.scene.UpdateSelfPos(this.direct)
}

func (this *PlayerTask) UpdateOthers() {
	this.scene.UpdatePos()
}
