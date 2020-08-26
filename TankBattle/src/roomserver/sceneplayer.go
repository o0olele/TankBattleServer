package main

import (
	common "common"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type ScenePlayer struct {
	id         uint32      //玩家id
	self       *PlayerInfo //玩家个人信息
	playerTask *PlayerTask
	scene      *Scene
	curflag    map[uint32]*ScenePlayer
	lastflag   map[uint32]*ScenePlayer
	senddie    bool

	isMove bool
	next   common.Pos
	angle  uint32
	speed  float64
	drag   float64

	bullets    map[uint32]*common.Bullet
	bulletnum  uint32
	curbullet  map[uint32]*common.Bullet
	lastbullet map[uint32]*common.Bullet

	movereq  *common.ReqMoveMsg
	turnreq  *common.ReqTurnMsg
	shootreq *common.ReqShootMsg
}

func NewScenePlayer(player *PlayerTask, scene *Scene) *ScenePlayer {
	s := &ScenePlayer{
		id:    player.id,
		scene: scene,
		self: &PlayerInfo{
			id: player.id,
			HP: common.FullHP,
		},
		playerTask: player,
		//others:   make(map[uint32]*PlayerTask),
		curflag:    make(map[uint32]*ScenePlayer),
		lastflag:   make(map[uint32]*ScenePlayer),
		speed:      1,
		drag:       0.1,
		senddie:    false,
		bullets:    make(map[uint32]*common.Bullet),
		lastbullet: make(map[uint32]*common.Bullet),
		curbullet:  make(map[uint32]*common.Bullet),
	}
	return s
}

func (this *ScenePlayer) CaculateNext(direct uint32, power uint32) {
	this.speed = float64(power)
	this.next.X = this.self.pos.X + math.Sin(float64(direct)*math.Pi/180)*this.speed
	this.next.Y = this.self.pos.Y + math.Cos(float64(direct)*math.Pi/180)*this.speed
	this.UpdateSpeed()
}

//初始化子弹
func (this *ScenePlayer) addBullet(direct uint32) {
	initpos := this.self.pos
	this.bullets[this.bulletnum] = &common.Bullet{
		Id:     this.bulletnum,
		Btype:  this.id,
		Pos:    initpos,
		Direct: direct,
		Time:   time.Now().Unix(),
	}

	this.bulletnum = (this.bulletnum + 1) % 10000
}

//更新视野内子弹并判断是否击中自己
func (this *ScenePlayer) getBullet() {
	for _, p := range this.scene.players {
		for _, bullet := range p.bullets {
			if time.Now().Unix()-int64(bullet.Time) > common.BulletLife {
				delete(p.bullets, bullet.Id)
				continue
			}
			if math.Abs(bullet.Pos.X-this.self.pos.X) < common.SceneHeight/2 &&
				math.Abs(bullet.Pos.Y-this.self.pos.Y) < common.SceneWidth/2 {
				angle := bullet.Direct
				last := *bullet
				bullet.Pos.X += math.Sin(float64(angle)*math.Pi/180) * common.BulletSpeed
				bullet.Pos.Y += math.Cos(float64(angle)*math.Pi/180) * common.BulletSpeed
				if this.self.HP > 0 && this.beshoot(&last, bullet) {
					this.self.HP--
					if this.self.HP == 0 {
						this.senddie = true
					}
					bullet.Time += common.BulletLife
					delete(p.bullets, bullet.Id)
					continue
				}
				this.curbullet[bullet.Id] = bullet
			}

		}
	}

}

func (this *ScenePlayer) beshoot(last, next *common.Bullet) bool {
	if last.Btype == this.id {
		return false
	}
	ndot := common.Dot{X: next.Pos.X, Y: next.Pos.Y}
	ldot := common.Dot{X: last.Pos.X, Y: last.Pos.Y}

	pdot := common.Dot{X: this.self.pos.X, Y: this.self.pos.Y}
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
// func (this *ScenePlayer) getBullet() {
// 	this.bullets = []*common.RetBullet{}

// 	this.room.allbullet.Range(func(key, value interface{}) bool {
// 		bullet, ok := value.(*common.Bullet)
// 		if !ok {
// 			return false
// 		}

// 		if math.Abs(bullet.Pos.X-this.self.X) < common.SceneHeight/2 &&
// 			math.Abs(bullet.Pos.Y-this.self.Y) < common.SceneWidth/2 {
// 			this.bullets = append(this.bullets, &common.RetBullet{Pos: bullet.Pos, Id: bullet.Id})
// 		}
// 		return true
// 	})
// }

func (this *ScenePlayer) UpdateSpeed() {
	this.speed = math.Max(0, this.speed-this.drag)
}

//更新视野
func (this *ScenePlayer) UpdatePos() {

	for _, user := range this.scene.players {
		if user.self.id == this.self.id {
			continue
		}
		if math.Abs(user.self.pos.X-this.self.pos.X) < common.SceneHeight/2 &&
			math.Abs(user.self.pos.Y-this.self.pos.Y) < common.SceneWidth/2 {
			this.curflag[user.id] = user
		}
	}
}

func (this *ScenePlayer) setIsMove() {
	if math.Abs(this.speed) < 1e-5 {
		this.isMove = false
	}
}

func (this *ScenePlayer) DoShoot() {
	if this.shootreq != nil {
		this.addBullet(this.shootreq.Direct)
		this.shootreq = nil
	}
}
func (this *ScenePlayer) DoMove() {

	if this.movereq != nil {
		this.CaculateNext(this.movereq.Direct, this.movereq.Power)
		this.self.pos = this.next
		this.isMove = true
		this.movereq = nil
	}
	if this.turnreq != nil {
		this.self.pos.Ag = this.turnreq.Direct
		this.turnreq = nil
		this.isMove = true
	}
}

func aoi(last, cur map[uint32]bool) (add []uint32, remove []uint32, move []uint32) {
	for id := range last {
		if _, ok := cur[id]; !ok {
			remove = append(remove, id)
		} else {
			move = append(move, id)
		}
	}
	for id := range cur {
		if _, ok := last[id]; !ok {
			add = append(add, id)
		}
	}
	return
}

func (this *ScenePlayer) getMoveMsg(msg *common.RetSceneMsg) {
	var last, cur map[uint32]bool
	last = make(map[uint32]bool)
	cur = make(map[uint32]bool)
	for id := range this.lastflag {
		last[id] = true
	}
	for id := range this.curflag {
		cur[id] = true
	}
	add, remove, move := aoi(last, cur)
	for _, id := range add {
		msg.Add = append(msg.Add, common.Add{
			Userid: id,
			Pos:    this.curflag[id].self.pos,
			HP:     this.curflag[id].self.HP,
		})
	}
	for _, id := range move {
		if this.curflag[id].isMove {
			msg.Move = append(msg.Move, common.Move{
				Userid: id,
				Pos:    this.curflag[id].self.pos,
				HP:     this.curflag[id].self.HP,
			})
		}
	}
	for _, id := range remove {
		msg.ReMove = append(msg.ReMove, common.ReMove{
			Userid: id,
		})
	}
}
func (this *ScenePlayer) getBulletMsg(msg *common.RetSceneMsg) {
	var last, cur map[uint32]bool
	last = make(map[uint32]bool)
	cur = make(map[uint32]bool)
	for id := range this.lastbullet {
		last[id] = true
	}
	for id := range this.curbullet {
		cur[id] = true
	}
	add, remove, move := aoi(last, cur)
	for _, id := range add {
		msg.Bullets.Add = append(msg.Bullets.Add, common.RetBullet{
			Id:  id,
			Pos: this.curbullet[id].Pos,
		})
	}
	for _, id := range move {
		msg.Bullets.Move = append(msg.Bullets.Move, common.RetBullet{
			Id:  id,
			Pos: this.curbullet[id].Pos,
		})
	}
	for _, id := range remove {
		msg.Bullets.ReMove = append(msg.Bullets.ReMove, id)
	}
}

func (this *ScenePlayer) relive() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	this.self.HP = common.FullHP
	this.self.pos.X = float64(r.Intn(10))
	this.self.pos.Y = float64(r.Intn(10))
}

func (this *ScenePlayer) sendSceneMsg() {
	this.UpdatePos()
	this.getBullet()
	if this.senddie {
		this.playerTask.SendDieMsg()
		this.senddie = false
		return
	}
	if !this.isMove && len(this.curbullet) == 0 && len(this.curflag) == 0 {
		return
	}
	msg := &common.RetSceneMsg{
		Add:    []common.Add{},
		ReMove: []common.ReMove{},
		Move:   []common.Move{},
	}
	fmt.Println("--------", this.id, "---------")
	msg.Move = append(msg.Move, common.Move{
		Userid: this.id,
		Pos:    this.self.pos,
		HP:     this.self.HP,
	})
	this.getMoveMsg(msg)
	this.getBulletMsg(msg)
	// for k, v := range this.lastflag {
	// 	fmt.Println(k, v.self.pos)
	// }
	// for k, v := range this.curflag {
	// 	fmt.Println(k, v.self.pos)
	// }

	// for id := range this.lastflag {
	// 	//上一次存在，这次不存在，remove
	// 	if _, ok := this.curflag[id]; !ok {

	// 		fmt.Println("remove")
	// 		msg.ReMove = append(msg.ReMove, common.ReMove{
	// 			Userid: id,
	// 		})
	// 	} else { //上一次存在，这次存在，move(移动)
	// 		fmt.Println("move")
	// 		if this.curflag[id].isMove {
	// 			msg.Move = append(msg.Move, common.Move{
	// 				Userid: id,
	// 				Pos:    this.curflag[id].self.pos,
	// 				HP:     this.curflag[id].self.HP,
	// 			})
	// 		}

	// 	}
	// }

	// for id := range this.curflag {
	// 	//这次存在，上一次不存在,add
	// 	if _, ok := this.lastflag[id]; !ok {
	// 		fmt.Println("add")
	// 		msg.Add = append(msg.Add, common.Add{
	// 			Userid: id,
	// 			Pos:    this.curflag[id].self.pos,
	// 			HP:     this.curflag[id].self.HP,
	// 		})
	// 	}
	// }
	this.lastflag = this.curflag
	this.curflag = make(map[uint32]*ScenePlayer)

	this.lastbullet = this.curbullet
	this.curbullet = make(map[uint32]*common.Bullet)
	this.playerTask.SendSceneMsg(msg)
	this.setIsMove()
}
