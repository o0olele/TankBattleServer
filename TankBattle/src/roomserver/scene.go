package main

import (
	common "common"
	"fmt"

	"github.com/golang/glog"
)

type Scene struct {
	players  map[uint32]*ScenePlayer
	room     *Room
	Obstacle *map[uint32]*common.Obstacle
}

func NewScene(room *Room) *Scene {
	scene := &Scene{
		room:    room,
		players: make(map[uint32]*ScenePlayer),
	}

	scene.init()

	return scene
}

func (this *Scene) init() {
	this.Obstacle = GenerateRandMap()
}

func (this *Scene) AddPlayer(p *PlayerTask) {
	fmt.Println("add player")
	this.players[p.id] = NewScenePlayer(p, this)
}

//定时发送
func (this *Scene) sendRoomMsg() {
	for _, p := range this.players {
		p.sendSceneMsg()
	}
}

func (this *Scene) SendOverMsg() {
	for _, player := range this.players {
		player.SendOverMsg()
	}
}

func (this *Scene) sendTime(t uint64) {
	for _, p := range this.players {
		p.sendTime(t)
	}
}
func (this *Scene) UpdateOP(op *opMsg) {
	fmt.Println(op.op)
	switch op.op {
	case common.PlayerMove:
		req, ok := op.args.(common.ReqMoveMsg)
		if !ok {
			glog.Info("[Move] move arg error")
			return
		}
		if this.players[req.Userid].self.HP == 0 {
			return
		}
		this.players[req.Userid].movereq = &req
	case common.PlayerTurn:
		req, ok := op.args.(common.ReqTurnMsg)
		if this.players[req.Userid].self.HP == 0 {
			return
		}
		if !ok {
			glog.Info("[Turn] turn arg error")
			return
		}
		this.players[req.Userid].turnreq = &req
	case common.AddBullet:
		req, ok := op.args.(common.ReqShootMsg)
		if this.players[req.Userid].self.HP == 0 {
			return
		}
		if !ok {
			glog.Info("[shoot] shoot arg error")
			return
		}
		this.players[req.Userid].shootreq = &req
	case common.Relive:
		id, ok := op.args.(uint32)
		if !ok {
			glog.Info("[relive] relive arg error")
			return
		}
		this.players[id].relive()
	}
}

func (this *Scene) UpdatePos() {
	for _, p := range this.players {
		p.DoMove()
		p.DoShoot()
	}
}
