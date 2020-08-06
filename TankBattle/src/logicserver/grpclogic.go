package main

import (
	"base/env"
	"context"
	"fmt"
	"io"

	proto "grpc"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type LogicGrpcClient struct {
	conn         *grpc.ClientConn
	mRouteClient proto.StreamLogicService_RouteClient

	roomInfoList map[string]*proto.ConnectRoomInfo
}

var mLogicGrpcClient *LogicGrpcClient

func LogicGrpcClient_GetMe() *LogicGrpcClient {
	if nil == mLogicGrpcClient {
		mLogicGrpcClient = &LogicGrpcClient{
			roomInfoList: make(map[string]*proto.ConnectRoomInfo),
		}
	}

	return mLogicGrpcClient
}

func (this *LogicGrpcClient) Init() bool {
	var err error
	this.conn, err = grpc.Dial(env.Get("logic", "grpc"), grpc.WithInsecure())
	if nil != err {
		glog.Error("[gRPC] Connect Fail ", err)
		return false
	}

	client := proto.NewStreamLogicServiceClient(this.conn)

	this.mRouteClient, err = client.Route(context.Background())
	if nil != err {
		glog.Error("[gRPC] Error ", err)
		return false
	}

	return true
}

func (this *LogicGrpcClient) Send(token string) bool {
	if nil == this.mRouteClient {
		glog.Error("[gRPC] Route Client is nil ")
		return false
	}

	this.mRouteClient.Send(&proto.LogicRequest{Token: token})

	for {
		stream, err := this.mRouteClient.Recv()
		if io.EOF == err {
			glog.Info("[gRPC] Client Got EOF")
			break
		}
		if nil != err {
			glog.Error("[gRPC] Client Error ", err)
			return false
		}

		fmt.Println("client recv:", string(stream.GetMInfo().GetIp()))
		this.roomInfoList[token] = stream.GetMInfo()
		break
	}

	return true
}

func (this *LogicGrpcClient) Close() {
	if nil != this.conn {
		this.conn.Close()
	}
}
