package main

import (
	common "common"
	"encoding/json"
	"fmt"
	"math"
	"sync"

	"github.com/golang/glog"
)

type Scene struct {
	self      common.Pos
	selfMutex sync.Mutex
	others    []common.Pos
	outters   []uint32

	width  float64
	height float64

	room *Room
}

func (this *Scene) UpdateSelfPos(direct uint32) {
	this.selfMutex.Lock()
	this.self.X += math.Sin(float64(direct) * math.Pi / 180)
	this.self.Y += math.Cos(float64(direct) * math.Pi / 180)
	fmt.Println(direct, math.Sin(float64(direct)*math.Pi/180), math.Cos(float64(direct)*math.Pi/180))
	this.selfMutex.Unlock()

	this.UpdatePos()
}

func (this *Scene) UpdatePos() {
	this.others = []common.Pos{}
	this.outters = []uint32{}
	for _, user := range this.room.players {
		if math.Abs(user.scene.self.X-this.self.X) < this.height/2 &&
			math.Abs(user.scene.self.Y-this.self.Y) < this.width/2 {
			this.others = append(this.others, common.Pos{Id: user.id, X: user.scene.self.X, Y: user.scene.self.Y})
		} else {
			this.outters = append(this.outters, user.id)
		}
	}
}

func (this *Scene) SceneMsg() []byte {
	this.UpdatePos()

	var users common.RetSceneMsg
	users.Users = []common.Pos{}
	users.Outter = []uint32{}

	users.Users = append(users.Users, this.self)
	users.Users = append(users.Users, this.others...)
	users.Outter = append(users.Outter, this.outters...)

	bytes, err := json.Marshal(users)
	if nil != err {
		glog.Error("[Scene] Scene Msg Error ", err)
		return nil
	}

	return bytes
}
