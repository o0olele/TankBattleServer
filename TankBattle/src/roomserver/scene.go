package main

import (
	common "common"

	"github.com/golang/glog"
)

type Scene struct {
	players map[uint32]*ScenePlayer
	room    *Room

	//self      common.Stat // 自身当前坐标
	//next      common.Stat // 使用next坐标进行计算，便于丢弃
	//selfMutex sync.Mutex
	//others    []common.Stat // 其他玩家信息
	//outters   []uint32      // 不再视野内玩家id
	//hasMove   bool          // 标识是否移动，用于后续优化（游戏开始时发送所有玩家列表，游戏中发送移动的玩家信息）

	bullets     []*common.RetBullet
	scenePlayer map[uint32]*ScenePlayer
}

func (this *Scene) Init(room *Room) {
	this.room = room
	this.players = make(map[uint32]*ScenePlayer)
	for _, p := range this.room.players {
		this.players[p.playerInfo.id] = NewScenePlayer(p, this)
	}
}

func (this *Scene) AddPlayer(player *ScenePlayer) {
	this.players[player.id] = player
}

// func (this *Scene) SceneMsg() []byte {
// 	this.UpdatePos()
// 	this.getBullet()
// 	var users common.RetSceneMsg
// 	users.Users = []common.Stat{}
// 	users.Outter = []uint32{}

// 	users.Users = append(users.Users, this.self)
// 	users.Users = append(users.Users, this.others...)
// 	users.Outter = append(users.Outter, this.outters...)
// 	users.Bullets = append(users.Bullets, this.bullets...)

// 	bytes, err := json.Marshal(users)
// 	if nil != err {
// 		glog.Error("[Scene] Scene Msg Error ", err)
// 		return nil
// 	}

// 	return bytes
// }

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
		this.players[req.Userid].movereq = req
		//this.players[req.Userid].UpdateSelfPos(req.Direct)
		// angle = xxx  speed = power
	case common.PlayerTurn:
		req, ok := op.args.(common.ReqMoveMsg)
		if !ok {
			glog.Info("[Turn] turn arg error")
			return
		}
		this.players[req.Userid].self.playerInfo.pos.Ag = req.Direct
	}
}

func (this *Scene) UpdatePos() {
	for _, p := range this.players {
		p.DoMove()
		//p.UpdatePos()
	}
	//p.UpdatePos();
}
