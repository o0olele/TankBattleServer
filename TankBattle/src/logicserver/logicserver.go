package main

import (
	"base/env"
	"base/gonet"
	"flag"

	"github.com/golang/glog"
)

type LoginService struct {
	gonet.Service
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

}
