package common

type MsgType uint16

const (
	MsgType_Token MsgType = 0
	MsgType_Move  MsgType = 1
)

// 客户请求
type ReqGetIDMsg struct {
	DeviceId string `json:"deviceId"`
	Ip       string `json:"ip"`
}

// 服务器返回
type RetGetIDMsg struct {
	Id uint32 `json:"id"`
}
