package main

import (
	"base/env"
	"base/gonet"
	"flag"
	"time"

	"github.com/golang/glog"
)

type LogicServer struct {
	gonet.Service
}

var mServer *LogicServer

func LogicServer_GetMe() *LogicServer {
	if nil == mServer {
		mServer = &LogicServer{}
	}
	mServer.Derived = mServer
	return mServer
}

func (this *LogicServer) Init() bool {
	if !StartHttpServer() {
		glog.Error("[Start] Http Server Fail")
		return false
	}

	if !LogicGrpcClient_GetMe().Init() {
		glog.Error("[gRPC] Client Init Fail")
		return false
	}

	return true
}

func (this *LogicServer) Reload() {

}

func (this *LogicServer) MainLoop() {
	time.Sleep(time.Second)
}

func (this *LogicServer) Final() bool {
	LogicGrpcClient_GetMe().Close()
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

	LogicServer_GetMe().Main()
}
