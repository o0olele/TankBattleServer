package gonet

import (
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"

	"github.com/golang/glog"
)

type IService interface {
	Init() bool
	Reload()
	MainLoop()
	Final() bool
}

type Service struct {
	terminated bool
	Derived    IService
}

func (this *Service) Terminate() {
	this.terminated = true
}

func (this *Service) IsTerminated() bool {
	return this.terminated
}

func (this *Service) Main() bool {

	defer func() {
		// catch errors before panic
		if err := recover(); nil != err {
			glog.Error("[Unexpeted] ", err, "\n", string(debug.Stack()))
		}
	}()

	// catch system signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGHUP)
	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGHUP:
				this.Derived.Reload()
			default:
				this.Terminate()
			}
			glog.Info("[Service] Got Signal ", sig)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	if !this.Derived.Init() {
		return false
	}

	for !this.IsTerminated() {
		this.Derived.MainLoop()
	}

	this.Derived.Final()
	return true

}
