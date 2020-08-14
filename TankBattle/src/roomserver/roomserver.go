package main

import (
	"base/env"
	"base/gonet"
	"flag"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type RoomServer struct {
	gonet.Service
	roomser *gonet.WebSocketServer
}

var mServer *RoomServer

func RoomServer_GetMe() *RoomServer {
	if nil == mServer {
		mServer = &RoomServer{
			roomser: &gonet.WebSocketServer{},
		}
		mServer.Derived = mServer
		mServer.roomser.Derived = mServer
	}

	return mServer
}

func (this *RoomServer) Init() bool {

	if !RoomGrpcClient_GetMe().Init() {
		glog.Error("[gRPC] Room Client Init Fail")
		return false
	}

	go func() {
		err := this.roomser.WSBind(env.Get("room", "listen"))
		if nil != err {
			glog.Error("[Start] Bind Port Fail")
		}
	}()

	return true
}

func (this *RoomServer) Reload() {

}

func (this *RoomServer) MainLoop() {
	time.Sleep(time.Second)
}

func (this *RoomServer) Final() bool {
	RoomGrpcClient_GetMe().Close()
	return true
}

func (this *RoomServer) OnWSAccept(conn *websocket.Conn) {
	glog.Info("[WS] Connected")
	NewPlayerTask(conn).Start()
}

var (
	config  = flag.String("config", "", "config file")
	logfile = flag.String("logfile", "", "log file")
)

func main() {
	flag.Parse()

	env.Load(*config)

	if "" != *logfile {
		glog.SetLogFile(*logfile)
	} else {
		glog.SetLogFile(env.Get("room", "log"))
	}
	defer glog.Flush()

	RoomServer_GetMe().Main()
}
