package main

import (
	"base/env"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	common "common"

	"github.com/golang/glog"
	"github.com/rs/cors"
)

// /getname?Id=xxx
func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	randName := "YoRHa" + GetDateFormat()

	fmt.Fprintf(w, randName)
}

// /getid?DeviceId=xxx&Ip=xxx
func GetIDHandler(w http.ResponseWriter, r *http.Request) {
	msg := common.ReqGetIDMsg{
		DeviceId: r.FormValue("DeviceId"),
		Ip:       r.FormValue("Ip"),
	}

	glog.Info("[login] Get msg", msg)

	id, err := AccRedis_GetMe().GetIncID()
	if nil != err {
		glog.Error("[login] Get Inc Id Fail ", err)
		return
	}

	err = json.NewEncoder(w).Encode(&common.RetGetIDMsg{
		Id:   id,
		Name: "YoRHa",
	})
	if nil != err {
		glog.Error("[login] Return id Fail ", err)
		return
	}

	err = AccRedis_GetMe().SetDeviceIdAndIp(id, &msg)
	if nil != err {
		glog.Error("[login] Set Userinfo Fail ", err)
		return
	}

}

// /getroom?Id=xxx
func GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("Id"))
	if nil != err {
		glog.Error("[logic] get/parse user id fail ", err)
		return
	}

	info, err := GetVailabelRoomInfo("token")
	if nil != err {
		glog.Error("[logic] RPC get room info fail ", err)
		return
	}
	bytes, err := json.Marshal(common.Token{
		Id:   uint32(id),
		Time: time.Now().Unix(),
	})

	token, err := common.GetTokenSSL(bytes)
	if nil != err {
		glog.Error("[logic] get openssl token fail ", err)
		return
	}

	err = json.NewEncoder(w).Encode(&common.RetGetRoom{
		Ip:    "0.0.0.0",
		Port:  info.Port,
		Token: token,
	})
	if nil != err {
		glog.Error("[logic] return room info fail ", err)
		return
	}
	glog.Info("[logic] token ", token)
}

// 时间戳转年月日 时分秒
func GetDateFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func StartHttpServer() bool {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/getname", GetIDHandler)
	mux.HandleFunc("/getroom", GetRoomHandler)
	mux.HandleFunc("/getid", GetIDHandler)

	addr := env.Get("logic", "listen")
	handler := c.Handler(mux)
	http.ListenAndServe(addr, handler)

	glog.Info("[Start] Bind Port Success, Port=", addr)

	return true
}
