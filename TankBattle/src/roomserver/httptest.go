package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang/glog"
)

func startHttpTest() {
	http.HandleFunc("/getroom", GetRoom)
	http.HandleFunc("/getload", GetLoad)
	http.ListenAndServe("127.0.0.1:9999", nil)
	glog.Info("httptest端口绑定成功 9999")
}

func GetRoom(res http.ResponseWriter, req *http.Request) {
	id, _ := strconv.ParseUint(req.URL.Query()["id"][0], 10, 32)
	r, _ := RoomMgr_GetMe().GetRoom(PlayerTask{id: uint32(id)})
	glog.Info(r)
	jsonstr, _ := json.Marshal(r)
	glog.Info("jsonstr:", jsonstr)
	res.Write([]byte(jsonstr))
}

func GetLoad(res http.ResponseWriter, req *http.Request) {
	jsonstr, _ := json.Marshal(RoomMgr_GetMe().GetLoad())
	res.Write([]byte(jsonstr))
}
