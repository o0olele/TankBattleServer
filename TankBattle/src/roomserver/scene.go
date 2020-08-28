package main

import (
	common "common"

	"github.com/golang/glog"
)

type Scene struct {
	players  map[uint32]*ScenePlayer
	room     *Room
	Obstacle *map[uint32]*common.Obstacle
}

func (this *Scene) Init(room *Room) {
	this.room = room
	this.players = make(map[uint32]*ScenePlayer)
	this.Obstacle = GenerateRandMap()
	for _, p := range this.room.players {
		this.players[p.id] = NewScenePlayer(p, this)
	}
}

func (this *Scene) AddPlayer(p *PlayerTask) {
	this.players[p.id] = NewScenePlayer(p, this)
}

//定时发送
func (this *Scene) sendRoomMsg() {

	for _, p := range this.players {
		p.sendSceneMsg()
	}
}

func (this *Scene) UpdateOP(op *opMsg) {
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
