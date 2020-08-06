package main

import (
	"base/env"
	"flag"

	"github.com/golang/glog"
)

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

	StartGrpcServer()

}
