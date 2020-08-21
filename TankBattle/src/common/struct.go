package common

type MsgType uint16

const (
	MsgType_Token MsgType = 0
	MsgType_Move  MsgType = 1
	MsgType_Finsh MsgType = 2
	MsgType_Shoot MsgType = 3
	MsgType_Heart MsgType = 4
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
)

type Stat struct {
	Id uint32
	X  float64
	Y  float64
	HP uint
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
type RetSceneMsg struct {
	Users   []Stat       `json:"users"`
	Outter  []uint32     `json:"outter"`
	Bullets []*RetBullet `json:"bullets"`
}

type RetTimeMsg struct {
	Time uint64 `json:"time"`
}
type RetBullet struct {
	Id  uint32
	Pos Stat
}
