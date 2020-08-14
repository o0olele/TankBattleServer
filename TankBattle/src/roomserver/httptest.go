package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
)

func startHttpTest() {
	x := http.NewServeMux()
	x.HandleFunc("/getroom", GetRoom)
	x.HandleFunc("/getload", GetLoad)
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		glog.Error("绑定9999失败")
	}
	srv := &http.Server{
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
		Handler:      x,
	}
	go srv.Serve(ln)
	glog.Info("httptest端口绑定成功 9999")
}

func GetRoom(res http.ResponseWriter, req *http.Request) {
	id, _ := strconv.ParseUint(req.URL.Query()["id"][0], 10, 32)
	r, _ := RoomMgr_GetMe().GetRoom(&PlayerTask{id: uint32(id)})
	glog.Info(r)
	jsonstr, _ := json.Marshal(r)
	glog.Info("jsonstr:", jsonstr)
	res.Write([]byte(jsonstr))
}

func GetLoad(res http.ResponseWriter, req *http.Request) {
	r, err := RoomMgr_GetMe().unFullRoom.Front().(*Room)
	fmt.Println("getload:", r, err, RoomMgr_GetMe().unFullRoom)
	//sglog.Info(RoomMgr_GetMe())
	//jsonstr, _ := json.Marshal(RoomMgr_GetMe().GetLoad())
	//res.Write([]byte(jsonstr))
}
