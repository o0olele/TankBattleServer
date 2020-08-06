package main

import (
	"base/env"
	"base/gonet"
	"flag"
	"time"

	"github.com/golang/glog"
)

type RoomServer struct {
	gonet.Service
}

var mServer *RoomServer

func RoomServer_GetMe() *RoomServer {
	if nil == mServer {
		mServer = &RoomServer{}
		mServer.Derived = mServer
	}

	return mServer
}

func (this *RoomServer) Init() bool {

	if !RoomGrpcClient_GetMe().Init() {
		glog.Error("[gRPC] Room Client Init Fail")
		return false
	}
	return true
}

func (this *RoomServer) Reload() {

}

func (this *RoomServer) MainLoop() {
	time.Sleep(time.Second)
}

func (this *RoomServer) Final() bool {

	return true
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
		glog.SetLogFile(env.Get("logic", "log"))
	}
	defer glog.Flush()

	RoomServer_GetMe().Main()
}
