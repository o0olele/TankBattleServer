package main

import (
	common "common"
	"encoding/json"
	"math"
	"sync"

	"github.com/golang/glog"
)

type Scene struct {
	self      common.Pos
	next      common.Pos
	selfMutex sync.Mutex
	others    []common.Pos
	outters   []uint32
	hasMove   bool

	speed float64

	room *Room
}

func (this *Scene) CaculateNext(direct uint32) {
	this.next.X = this.self.X + math.Sin(float64(direct)*math.Pi/180)*this.speed
	this.next.Y = this.self.Y + math.Cos(float64(direct)*math.Pi/180)*this.speed
}

func (this *Scene) UpdateSelfPos(direct uint32) {
	this.selfMutex.Lock()

	/*if 0 == this.speed {
		this.hasMove = false
		return
	}*/
	this.CaculateNext(direct)
	this.self = this.next
	//this.hasMove = true
	this.UpdateSpeed(0)

	this.selfMutex.Unlock()
}

func (this *Scene) UpdatePos() {
	this.others = []common.Pos{}
	this.outters = []uint32{}
	for _, user := range this.room.players {
		/*if !user.scene.hasMove {
			continue
		}*/
		if math.Abs(user.scene.self.X-this.self.X) < common.SceneHeight/2 &&
			math.Abs(user.scene.self.Y-this.self.Y) < common.SceneWidth/2 {
			this.others = append(this.others, common.Pos{Id: user.id, X: user.scene.self.X, Y: user.scene.self.Y})
		} else {
			this.outters = append(this.outters, user.id)
		}
	}
}

func (this *Scene) UpdateSpeed(s float64) {
	this.speed = s
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
