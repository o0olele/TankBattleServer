package common

type MsgType uint16

const (
	MsgType_Token  MsgType = 0
	MsgType_Move   MsgType = 1
	MsgType_Finsh  MsgType = 2
	MsgType_Shoot  MsgType = 3
	MsgType_Heart  MsgType = 4
	MsgType_Direct MsgType = 5
	MsgType_Relife MsgType = 6
)
const (
	PlayerMove uint32 = 1
	PlayerTurn uint32 = 2
	AddBullet  uint32 = 3
	BulletMove uint32 = 4
	Relive     uint32 = 5
)
const (
	SceneSpeed  float64 = 0.2
	SceneWidth  float64 = 200
	SceneHeight float64 = 200
)

const (
	BulletSpeed float64 = 2
	BulletLife  int64   = 5
)

const (
	Task_TimeOut = 20
)

const (
	PlayerSize float64 = 3
	FullHP     uint32  = 100
)

const (
	MapHeight uint32 = 25
	MapWidth  uint32 = 25
)

type Pos struct {
	X  float64
	Y  float64
	Ag uint32
}
type Stat struct {
	Pos Pos
	HP  uint32
}

type Token struct {
	Id   uint32
	Time int64
}

type Bullet struct {
	Id     uint32
	Btype  uint32
	Pos    Pos
	Direct uint32
	Time   int64
}
type Move struct {
	Userid uint32
	Pos    Pos
	HP     uint32
}
type ReMove struct {
	Userid uint32
}
type Add struct {
	Userid uint32
	Pos    Pos
	HP     uint32
}

// 客户请求
type ReqGetIDMsg struct {
	DeviceId string
	Ip       string
}

type ReqMoveMsg struct {
	Userid uint32
	Direct uint32
	Power  uint32
}
type ReqTurnMsg struct {
	Userid uint32
	Direct uint32
}
type ReqShootMsg struct {
	Userid uint32
	Direct uint32
}

// 服务器返回
type RetGetIDMsg struct {
	Id   uint32
	Name string
}

type RetSceneMsg struct {
	Move      []Move
	ReMove    []ReMove
	Add       []Add
	Bullets   RetBulletMsg
	Obstacles RetObstacle
}

type RetTimeMsg struct {
	Time uint64
}
type RetBulletMsg struct {
	Move   []RetBullet
	ReMove []uint32
	Add    []RetBullet
}
type RetBullet struct {
	Id  uint32
	Pos Pos
}

type RetOverMsg struct {
	End bool
}

type RetGetRoom struct {
	Ip    string
	Port  uint32
	Token string
}

type RetYouDie struct {
	Die bool
}

type Obstacle struct {
	Id     uint32
	Pos    Pos
	Length uint32
	Width  uint32
	Height uint32
}

type RetObstacle struct {
	ReMove []uint32
	Add    []Obstacle
}
