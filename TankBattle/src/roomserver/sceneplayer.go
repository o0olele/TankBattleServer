package main

import (
	common "common"
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

	curobstacle  map[uint32]*common.Obstacle
	lastobstacle map[uint32]*common.Obstacle

	movereq  *common.ReqMoveMsg
	turnreq  *common.ReqTurnMsg
	shootreq *common.ReqShootMsg
}

func NewScenePlayer(player *PlayerTask, scene *Scene) *ScenePlayer {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &ScenePlayer{
		id:    player.id,
		scene: scene,
		self: &PlayerInfo{
			id: player.id,
			HP: common.FullHP,
			pos: common.Pos{
				X: float64(r.Intn(int(common.MapWidth))),
				Y: float64(r.Intn(int(common.MapHeight))),
			},
		},
		playerTask: player,
		//others:   make(map[uint32]*PlayerTask),
		curflag:      make(map[uint32]*ScenePlayer),
		lastflag:     make(map[uint32]*ScenePlayer),
		speed:        0,
		drag:         0.001,
		senddie:      false,
		bullets:      make(map[uint32]*common.Bullet),
		lastbullet:   make(map[uint32]*common.Bullet),
		curbullet:    make(map[uint32]*common.Bullet),
		curobstacle:  make(map[uint32]*common.Obstacle),
		lastobstacle: make(map[uint32]*common.Obstacle),
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
func (this *ScenePlayer) bulletHitObstacle(last, next *common.Bullet) bool {
	for _, ob := range *this.scene.Obstacle {
		if this.beshoot(last, next, &ob.Pos) {
			return true
		}
	}
	return false
}

//更新视野内子弹并判断是否击中自己
func (this *ScenePlayer) getBullet() {
	for _, p := range this.scene.players {
		for _, bullet := range p.bullets {
			if time.Now().Unix()-int64(bullet.Time) > common.BulletLife || !this.isInMap(&bullet.Pos) {
				delete(p.bullets, bullet.Id)
				continue
			}
			if math.Abs(bullet.Pos.X-this.self.pos.X) < common.SceneHeight/2 &&
				math.Abs(bullet.Pos.Y-this.self.pos.Y) < common.SceneWidth/2 {
				angle := bullet.Direct
				last := *bullet
				bullet.Pos.X += math.Sin(float64(angle)*math.Pi/180) * common.BulletSpeed
				bullet.Pos.Y += math.Cos(float64(angle)*math.Pi/180) * common.BulletSpeed
				if this.bulletHitObstacle(&last, bullet) {
					bullet.Time += common.BulletLife
					delete(p.bullets, bullet.Id)
					continue
				}
				if this.self.HP > 0 && this.beshoot(&last, bullet, &this.self.pos) {
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

func (this *ScenePlayer) beshoot(last, next *common.Bullet, pos *common.Pos) bool {
	if last.Btype == this.id {
		return false
	}
	ndot := common.Dot{X: next.Pos.X, Y: next.Pos.Y}
	ldot := common.Dot{X: last.Pos.X, Y: last.Pos.Y}

	pdot := common.Dot{X: pos.X, Y: pos.Y}
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

func (this *ScenePlayer) UpdateSpeed() {
	this.speed = math.Max(0, this.speed-this.drag)
}

func (this *ScenePlayer) getObstacle() {
	for _, obstacle := range *this.scene.Obstacle {
		if math.Abs(obstacle.Pos.X-this.self.pos.X) < common.SceneHeight/2 &&
			math.Abs(obstacle.Pos.Y-this.self.pos.Y) < common.SceneWidth/2 {
			this.curobstacle[obstacle.Id] = obstacle
		}
	}
}

//更新视野
func (this *ScenePlayer) UpdatePos() {

	for _, user := range this.scene.players {
		if user.self.id == this.self.id || user.self.HP == 0 {
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

func (this *ScenePlayer) isInMap(pos *common.Pos) bool {
	if math.Abs(pos.X) > math.Abs(float64(common.MapWidth)) || math.Abs(pos.Y) > math.Abs(float64(common.MapHeight)) {
		return false
	}
	return true
}

func (this *ScenePlayer) isCollision(pos *common.Pos, size float64, ob *common.Obstacle) bool {
	if math.Abs(ob.Pos.X-pos.X) < math.Abs(float64(ob.Width)/2+size/2) && math.Abs(ob.Pos.Y-pos.Y) < math.Abs(float64(ob.Height)/2+size/2) {
		return true
	}
	return false
}

func (this *ScenePlayer) setInMap(pos *common.Pos) {
	mh := float64(common.MapHeight)
	mw := float64(common.MapWidth)
	if pos.X < (-mw) {
		pos.X = -mw
	}
	if pos.Y < (-mh) {
		pos.Y = -mh
	}
	if pos.X > mw {
		pos.X = mw
	}
	if pos.Y > mh {
		pos.Y = mh
	}

}

func (this *ScenePlayer) DoMove() {

	if this.movereq != nil {
		this.isMove = true

		this.CaculateNext(this.movereq.Direct, this.movereq.Power)
		this.setInMap(&this.next)
		for _, ob := range *this.scene.Obstacle {
			if this.isCollision(&this.self.pos, common.PlayerSize, ob) {
				this.next.X = ob.Pos.X + float64(ob.Width) + common.PlayerSize + 0.1
			}
			if this.isCollision(&this.next, common.PlayerSize, ob) {
				this.next = this.self.pos
				break
			}
		}
		this.self.pos = this.next
		this.movereq = nil
	}
	if this.turnreq != nil && this.turnreq.Direct != this.self.pos.Ag {
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
	this.lastflag = this.curflag
	this.curflag = make(map[uint32]*ScenePlayer)
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
	this.lastbullet = this.curbullet
	this.curbullet = make(map[uint32]*common.Bullet)
}

func (this *ScenePlayer) getObstacleMsg(msg *common.RetSceneMsg) {
	var last, cur map[uint32]bool
	last = make(map[uint32]bool)
	cur = make(map[uint32]bool)

	for id := range this.lastobstacle {
		last[id] = true
	}
	for id := range this.curobstacle {
		cur[id] = true
	}
	add, remove, _ := aoi(last, cur)
	for _, id := range add {
		msg.Obstacles.Add = append(msg.Obstacles.Add, *this.curobstacle[id])
	}
	for _, id := range remove {
		msg.Obstacles.ReMove = append(msg.Obstacles.ReMove, id)
	}

	this.lastobstacle = this.curobstacle
	this.curobstacle = make(map[uint32]*common.Obstacle)
}
func (this *ScenePlayer) relive() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	this.self.HP = common.FullHP
	this.self.pos.X = float64(r.Intn(10))
	this.self.pos.Y = float64(r.Intn(10))
}

func (this *ScenePlayer) sendSceneMsg() {
	this.getObstacle()
	this.UpdatePos()
	this.getBullet()
	if this.senddie {
		this.playerTask.SendDieMsg()
		this.senddie = false
		return
	}

	msg := &common.RetSceneMsg{
		Add:    []common.Add{},
		ReMove: []common.ReMove{},
		Move:   []common.Move{},
	}
	if this.isMove {
		msg.Move = append(msg.Move, common.Move{
			Userid: this.id,
			Pos:    this.self.pos,
			HP:     this.self.HP,
		})
	}
	this.getMoveMsg(msg)
	this.getBulletMsg(msg)
	this.getObstacleMsg(msg)
	if !this.isMove && len(msg.Add) == 0 && len(msg.Move) == 0 && len(msg.ReMove) == 0 && len(this.curbullet) == 0 {
		return
	}
	this.playerTask.SendSceneMsg(msg)
	this.setIsMove()
}
