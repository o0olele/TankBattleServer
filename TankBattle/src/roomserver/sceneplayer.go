package main

import (
	common "common"
	"fmt"
	"math"
)

type ScenePlayer struct {
	id       uint32 //玩家id
	self     *PlayerTask
	scene    *Scene
	others   map[uint32]*PlayerTask
	curflag  map[uint32]bool
	lastflag map[uint32]bool

	next  common.Pos
	speed float64
}

func NewScenePlayer(udata *PlayerTask, scene *Scene) *ScenePlayer {
	s := &ScenePlayer{
		id:       udata.playerInfo.id,
		scene:    scene,
		self:     udata,
		others:   make(map[uint32]*PlayerTask),
		curflag:  make(map[uint32]bool),
		lastflag: make(map[uint32]bool),
		speed:    1,
	}
	s.self.playerInfo = udata.playerInfo
	return s
}

func (this *ScenePlayer) CaculateNext(direct uint32) {
	this.next.X = this.self.playerInfo.pos.X + math.Sin(float64(direct)*math.Pi/180)*this.speed
	this.next.Y = this.self.playerInfo.pos.Y + math.Cos(float64(direct)*math.Pi/180)*this.speed
}

func (this *ScenePlayer) UpdateSelfPos(direct uint32) {

	this.CaculateNext(direct)
	this.self.playerInfo.pos = this.next
}

//发射子弹
func (this *ScenePlayer) addBullet(direct uint32) {
	// initpos := this.pos
	// initpos.Id = this.self.Id
	// this.room.allbullet.Store(this.room.bulletcount, &common.Bullet{
	// 	Id:     this.room.bulletcount,
	// 	Btype:  this.self.Id,
	// 	Pos:    initpos,
	// 	Direct: direct,
	// 	Time:   time.Now().Unix(),
	// })
	// atomic.StoreUint32(&this.room.bulletcount, this.room.bulletcount+1)
}

// func updateBulletPos(bullet *common.Bullet, players map[uint32]*PlayerTask) {
// 	angle := bullet.Direct
// 	last := *bullet
// 	bullet.Pos.X += math.Sin(float64(angle)*math.Pi/180) * common.BulletSpeed
// 	bullet.Pos.Y += math.Cos(float64(angle)*math.Pi/180) * common.BulletSpeed
// 	for _, player := range players {
// 		if player.scene == nil {
// 			return
// 		}
// 		if player.scene.self.HP > 0 && beshoot(&last, bullet, player) {
// 			player.scene.self.HP--
// 			bullet.Time += common.BulletLife
// 		}
// 	}
// }

// func beshoot(last, next *common.Bullet, player *PlayerTask) bool {
// 	if player.scene == nil || last.Btype == player.scene.self.Id {
// 		return false
// 	}
// 	ndot := common.Dot{X: next.Pos.X, Y: next.Pos.Y}
// 	ldot := common.Dot{X: last.Pos.X, Y: last.Pos.Y}

// 	pdot := common.Dot{X: player.scene.self.X, Y: player.scene.self.Y}
// 	if common.GetDDDistance(ndot, pdot) < common.PlayerSize {
// 		return true
// 	}
// 	line := common.GetLine(ldot, ndot)
// 	if common.GetDLDistance(line, pdot) < common.PlayerSize {
// 		if common.TriCos(ldot, pdot, ndot) > 0 && common.TriCos(ndot, pdot, ldot) > 0 {
// 			return true
// 		}
// 	}
// 	return false
// }

// //获取视野内的子弹
// func (this *ScenePlayer) getBullet() {
// 	this.bullets = []*common.RetBullet{}

// 	this.room.allbullet.Range(func(key, value interface{}) bool {
// 		bullet, ok := value.(*common.Bullet)
// 		if !ok {
// 			return false
// 		}
// 		if time.Now().Unix()-int64(bullet.Time) > common.BulletLife {
// 			this.room.allbullet.Delete(bullet.Id)
// 			return true
// 		}
// 		if math.Abs(bullet.Pos.X-this.self.X) < common.SceneHeight/2 &&
// 			math.Abs(bullet.Pos.Y-this.self.Y) < common.SceneWidth/2 {
// 			this.bullets = append(this.bullets, &common.RetBullet{Pos: bullet.Pos, Id: bullet.Id})
// 		}
// 		return true
// 	})
// }

func (this *ScenePlayer) UpdateSpeed(s float64) {
	this.speed = s
}

//更新视野
func (this *ScenePlayer) UpdatePos() {

	for _, user := range this.scene.players {
		if user.self.playerInfo.id == this.self.playerInfo.id {
			continue
		}
		if math.Abs(user.self.playerInfo.pos.X-this.self.playerInfo.pos.X) < common.SceneHeight/2 &&
			math.Abs(user.self.playerInfo.pos.Y-this.self.playerInfo.pos.Y) < common.SceneWidth/2 {
			this.curflag[user.id] = true
			this.others[user.id] = user.self
			//fmt.Println("add others", this.self.playerInfo.id, this.others[user.id].playerInfo)
		}
	}
}

func (this *ScenePlayer) sendSceneMsg() {
	msg := &common.RetSceneMsg{
		Add:     []common.Add{},
		ReMove:  []common.ReMove{},
		Move:    []common.Move{},
		Bullets: []common.RetBullet{},
	}
	// fmt.Println("------", this.self.playerInfo.id, "------")
	// fmt.Println("this.lastflag", this.lastflag)
	// fmt.Println("this.curflag", this.curflag)
	// for k, v := range this.others {
	// 	fmt.Println("others", k, v.playerInfo)
	// }
	for id := range this.lastflag {
		//上一次存在，这次不存在，remove

		if _, ok := this.curflag[id]; !ok {

			fmt.Println("remove")
			msg.ReMove = append(msg.ReMove, common.ReMove{
				Userid: id,
			})
		} else { //上一次不存在，这次存在，move

			fmt.Println("move")
			msg.Move = append(msg.Move, common.Move{
				Userid: id,
				Pos:    this.others[id].playerInfo.pos,
				HP:     this.others[id].playerInfo.HP,
			})
		}
	}

	for id := range this.curflag {
		//这次存在，上一次不存在,add
		if _, ok := this.lastflag[id]; !ok {
			fmt.Println("add", this.self.playerInfo.id, id)
			fmt.Println(this.others[id].playerInfo)
			msg.Add = append(msg.Add, common.Add{
				Userid: id,
				Pos:    this.others[id].playerInfo.pos,
				HP:     this.others[id].playerInfo.HP,
			})
		}
	}
	this.lastflag = this.curflag
	this.curflag = make(map[uint32]bool)
	this.others = make(map[uint32]*PlayerTask)
	//fmt.Println("------", this.self.playerInfo.id, "------")

	this.self.SendSceneMsg(msg)
}