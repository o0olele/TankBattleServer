package main

import (
	common "common"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

type Scene struct {
	self      common.Stat // 自身当前坐标
	next      common.Stat // 使用next坐标进行计算，便于丢弃
	selfMutex sync.Mutex
	others    []common.Stat // 其他玩家信息
	outters   []uint32      // 不再视野内玩家id
	hasMove   bool          // 标识是否移动，用于后续优化（游戏开始时发送所有玩家列表，游戏中发送移动的玩家信息）

	speed   float64
	bullets []*common.RetBullet
	room    *Room
}

func (this *Scene) CaculateNext(direct uint32) {
	this.next.X = this.self.X + math.Sin(float64(direct)*math.Pi/180)*this.speed
	this.next.Y = this.self.Y + math.Cos(float64(direct)*math.Pi/180)*this.speed
	this.next.Id = this.self.Id
	this.next.HP = this.self.HP
}

func (this *Scene) UpdateSelfPos(direct uint32) {
	this.selfMutex.Lock()

	// 后续优化
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

//发射子弹
func (this *Scene) addBullet(direct uint32) {
	initpos := this.self
	initpos.Id = this.self.Id
	this.room.allbullet.Store(this.room.bulletcount, &common.Bullet{
		Id:     this.room.bulletcount,
		Btype:  this.self.Id,
		Pos:    initpos,
		Direct: direct,
		Time:   time.Now().Unix(),
	})
	atomic.StoreUint32(&this.room.bulletcount, this.room.bulletcount+1)
}

func updateBulletPos(bullet *common.Bullet, players map[uint32]*PlayerTask) {
	angle := bullet.Direct
	last := *bullet
	bullet.Pos.X += math.Sin(float64(angle)*math.Pi/180) * common.BulletSpeed
	bullet.Pos.Y += math.Cos(float64(angle)*math.Pi/180) * common.BulletSpeed
	for _, player := range players {
		if player.scene == nil {
			return
		}
		if player.scene.self.HP > 0 && beshoot(&last, bullet, player) {
			player.scene.self.HP--
			bullet.Time += common.BulletLife
		}
	}
}

func beshoot(last, next *common.Bullet, player *PlayerTask) bool {
	if player.scene == nil || last.Btype == player.scene.self.Id {
		return false
	}
	ndot := common.Dot{X: next.Pos.X, Y: next.Pos.Y}
	ldot := common.Dot{X: last.Pos.X, Y: last.Pos.Y}

	pdot := common.Dot{X: player.scene.self.X, Y: player.scene.self.Y}
	if common.GetDDDistance(ndot, pdot) < common.PlayerSize {
		return true
	}
	line := common.GetLine(ldot, ndot)
	if common.GetDLDistance(line, pdot) < common.PlayerSize {
		if common.TriCos(ldot, pdot, ndot) > 0 && common.TriCos(ndot, pdot, ldot) > 0 {
			return true
		}
	}
	return false
}

//获取视野内的子弹
func (this *Scene) getBullet() {
	this.bullets = []*common.RetBullet{}

	this.room.allbullet.Range(func(key, value interface{}) bool {
		bullet, ok := value.(*common.Bullet)
		if !ok {
			return false
		}
		if time.Now().Unix()-int64(bullet.Time) > common.BulletLife {
			this.room.allbullet.Delete(bullet.Id)
			return true
		}
		if math.Abs(bullet.Pos.X-this.self.X) < common.SceneHeight/2 &&
			math.Abs(bullet.Pos.Y-this.self.Y) < common.SceneWidth/2 {
			this.bullets = append(this.bullets, &common.RetBullet{Pos: bullet.Pos, Id: bullet.Id})
		}
		return true
	})

}
func (this *Scene) UpdatePos() {
	this.others = []common.Stat{}
	this.outters = []uint32{}
	for _, user := range this.room.players {
		/*后续优化
		  if !user.scene.hasMove {
			continue
		}*/
		if nil == user.scene {
			continue
		}

		if math.Abs(user.scene.self.X-this.self.X) < common.SceneHeight/2 &&
			math.Abs(user.scene.self.Y-this.self.Y) < common.SceneWidth/2 {
			this.others = append(this.others, common.Stat{Id: user.id, X: user.scene.self.X, Y: user.scene.self.Y, HP: user.scene.self.HP})
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
	this.getBullet()
	var users common.RetSceneMsg
	users.Users = []common.Stat{}
	users.Outter = []uint32{}

	users.Users = append(users.Users, this.self)
	users.Users = append(users.Users, this.others...)
	users.Outter = append(users.Outter, this.outters...)
	users.Bullets = append(users.Bullets, this.bullets...)

	bytes, err := json.Marshal(users)
	if nil != err {
		glog.Error("[Scene] Scene Msg Error ", err)
		return nil
	}
	if len(users.Bullets) != 0 {
		fmt.Println("--------------------------")
		fmt.Println(string(bytes))
	}
	return bytes
}
