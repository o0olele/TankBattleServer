package main

import (
	"base/env"
	"base/gonet"
	"flag"
	"time"

	"github.com/golang/glog"
)

type RcenterServer struct {
	gonet.Service
}

var mServer *RcenterServer

func RcenterServer_GetMe() *RcenterServer {
	if nil == mServer {
		mServer = &RcenterServer{}
		mServer.Derived = mServer
	}

	return mServer
}

func (this *RcenterServer) Init() bool {
	if !StartGrpcServer() {
		glog.Error("[gRPC] Start Server Fail")
		return false
	}

	return true
}

func (this *RcenterServer) Reload() {

}

func (this *RcenterServer) MainLoop() {
	time.Sleep(time.Second)
}

func (this *RcenterServer) Final() bool {
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
		glog.SetLogFile(env.Get("rcenter", "log"))
	}
	defer glog.Flush()

	RcenterServer_GetMe().Main()

}
