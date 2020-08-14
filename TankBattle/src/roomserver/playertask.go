package main

import (
	"base/gonet"
	"bytes"
	"encoding/binary"

	common "common"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type PlayerTask struct {
	wstask *gonet.WebSocketTask

	key  string
	id   uint32
	name string
	room *Room
}

func NewPlayerTask(conn *websocket.Conn) *PlayerTask {
	m := &PlayerTask{
		wstask: gonet.NewWebSocketTask(conn),
	}
	m.wstask.Derived = m

	return m
}

func (this *PlayerTask) Start() {
	this.wstask.Start()
	RoomMgr_GetMe().GetRoom(this)
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
	default:
	}
	return true
}
