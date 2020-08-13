package main

import (
	"base/gonet"

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

func (this *PlayerTask) OnClose() {

}

func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {

	return true
}
