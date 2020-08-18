package main

import (
	"base/env"
	"base/gonet"
	"flag"
	"time"

	"github.com/golang/glog"
)

type RLoad struct {
	ip   uint32
	port uint32
	load uint32
}

type RcenterServer struct {
	gonet.Service
	rList []RLoad
}

var mServer *RcenterServer

func RcenterServer_GetMe() *RcenterServer {
	if nil == mServer {
		mServer = &RcenterServer{
			rList: make([]RLoad, 0),
		}
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

func (this *RcenterServer) RegisterRoomServer(ip, port, load uint32) {
	this.rList = append(this.rList, RLoad{ip: ip, port: port, load: load})
}

func (this *RcenterServer) GetRoomServer() (uint32, uint32) {
	// todo 负载均衡

	return this.rList[0].ip, this.rList[0].port
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
