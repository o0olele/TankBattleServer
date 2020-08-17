package main

import (
	common "common"
	"encoding/json"
	"math"

	"github.com/golang/glog"
)

type Scene struct {
	self   common.Pos
	others []common.Pos

	width  float64
	height float64

	room *Room
}

func (this *Scene) UpdateSelfPos(direct int) {
	this.self.X += math.Sin(float64(direct))
	this.self.Y += math.Cos(float64(direct))
}

func (this *Scene) UpdatePos() {
	this.others = []common.Pos{}
	for _, user := range this.room.players {
		if math.Abs(user.scene.self.X-this.self.X) < this.height/2 &&
			math.Abs(user.scene.self.Y-this.self.Y) < this.width/2 {
			this.others = append(this.others, common.Pos{Id: user.id, X: user.scene.self.X, Y: user.scene.self.Y})

		}
	}
}

func (this *Scene) SceneMsg() []byte {
	var users common.RetSceneMsg
	users.Users = []common.Pos{}

	users.Users = append(users.Users, this.self)
	users.Users = append(users.Users, this.others...)

	bytes, err := json.Marshal(users)
	if nil != err {
		glog.Error("[Scene] Scene Msg Error ", err)
		return nil
	}

	return bytes
}
