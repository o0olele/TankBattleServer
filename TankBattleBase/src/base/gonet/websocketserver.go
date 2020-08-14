package gonet

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type IWebSocketServer interface {
	OnWSAccept(conn *websocket.Conn)
}

type WebSocketServer struct {
	Derived IWebSocketServer
}

var upgrader = websocket.Upgrader{}

func (this *WebSocketServer) WSBind(addr string) error {
	http.HandleFunc("/", this.WSListen)
	err := http.ListenAndServe(addr, nil)
	if nil != err {
		glog.Error("[WS] Init Fail, ", err)
		return err
	}
	glog.Info("[WS] Bind Success ", addr)
	return nil
}

func (this *WebSocketServer) WSListen(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		glog.Error("[WS] Failed to Upgrade", r.RemoteAddr, err)
		return
	}

	glog.Info("[WS] Connected ", r.RemoteAddr, w.Header().Get("Origin"), r.Header.Get("Origin"))

	if nil != this.Derived {
		this.Derived.OnWSAccept(c)
	}

}
