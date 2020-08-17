package main

import (
	"base/env"
	"context"
	"encoding/json"
	"strconv"
	"time"

	proto "proto"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type RoomGrpcClient struct {
	conn         *grpc.ClientConn
	mRouteClient proto.StreamRoomService_RouteClient
}

var mRoomGrpcClient *RoomGrpcClient

func RoomGrpcClient_GetMe() *RoomGrpcClient {
	if nil == mRoomGrpcClient {
		mRoomGrpcClient = &RoomGrpcClient{}
	}

	return mRoomGrpcClient
}

func (this *RoomGrpcClient) Init() bool {
	var err error
	this.conn, err = grpc.Dial(env.Get("room", "grpc"), grpc.WithInsecure())
	if nil != err {
		glog.Error("[gRPC] Connect Fail ", err)
		return false
	}

	if !this.InitClient() {
		glog.Error("[gRPC] Init Client Fail")
		return false
	}

	return this.SendRegist()
}

func (this *RoomGrpcClient) InitClient() bool {
	var err error
	client := proto.NewStreamRoomServiceClient(this.conn)

	this.mRouteClient, err = client.Route(context.Background())
	if nil != err {
		glog.Error("[gRPC] Error ", err)
		return false
	}

	return true
}

func (this *RoomGrpcClient) SendRegist() bool {
	if nil == this.mRouteClient {
		glog.Error("[gRPC] Route Client is nil ")
		return false
	}

	port, err := strconv.Atoi(env.Get("room", "port"))
	if nil != err {
		glog.Error("[Common] String to Int Error, ", err)
		return false
	}

	bytes, err := json.Marshal(proto.ConnectRoomInfo{
		Ip:   0,
		Port: uint32(port),
	})
	if nil != err {
		glog.Error("[Common] Struct to Json Error, ", err)
		return false
	}

	this.mRouteClient.Send(&proto.RoomRequest{
		Type: proto.MsgType_Regist,
		Data: bytes,
	})

	return true
}

func (this *RoomGrpcClient) SendLoad() {
	for {
		time.Sleep(time.Second * 2)

		this.mRouteClient.Send(&proto.RoomRequest{
			Type: proto.MsgType_Update,
			Data: []byte(strconv.Itoa(int(RoomMgr_GetMe().GetLoad()))),
		})
	}

}

func (this *RoomGrpcClient) Close() {
	if nil != this.conn {
		this.conn.Close()
	}
}
