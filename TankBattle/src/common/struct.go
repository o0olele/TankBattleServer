package common

type MsgType uint16

const (
	MsgType_Token  MsgType = 0
	MsgType_Move   MsgType = 1
	MsgType_Finsh  MsgType = 2
	MsgType_Shoot  MsgType = 3
	MsgType_Heart  MsgType = 4
	MsgType_Direct MsgType = 5
)
const (
	PlayerMove uint32 = 1
	PlayerTure uint32 = 2
	AddBullet  uint32 = 3
	BulletMove uint32 = 4
)
const (
	SceneSpeed  float64 = 0.2
	SceneWidth  float64 = 20
	SceneHeight float64 = 20
)

const (
	BulletSpeed float64 = 0.5
	BulletLife  int64   = 5
)

const (
	Task_TimeOut = 20
)

const (
	PlayerSize float64 = 1
	FullHP     uint32  = 100
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

// 客户请求
type ReqGetIDMsg struct {
	DeviceId string `json:"deviceId"`
	Ip       string `json:"ip"`
}

// 服务器返回
type RetGetIDMsg struct {
	Id uint32 `json:"id"`
}
type Bullet struct {
	Id     uint32
	Btype  uint32
	Pos    Stat
	Direct uint32
	Time   int64
}
type Move struct {
	Userid uint32 `json:"userid"`
	Pos    Pos
	HP     uint32
}
type ReMove struct {
	Userid uint32 `json:"userid"`
}
type Add struct {
	Userid uint32 `json:"userid"`
	Pos    Pos
	HP     uint32
}
type RetSceneMsg struct {
	Move    []Move      `json:"Move"`
	ReMove  []ReMove    `json:"ReMove"`
	Add     []Add       `json:"Add"`
	Bullets []RetBullet `json:"bullets"`
}

type RetTimeMsg struct {
	Time uint64 `json:"time"`
}
type RetBullet struct {
	Id  uint32
	Pos Stat
}

type RetOverMsg struct {
	End bool `json:"end"`
}

type ReqMoveMsg struct {
	Userid uint32
	Direct uint32
}
type ReqTurnMsg struct {
	Userid uint32
	Direct uint32
}
