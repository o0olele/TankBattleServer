package main

import (
	"base/env"
	"encoding/json"
	"fmt"
	proto "grpc"
	"io"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type LogicService struct {
}

func (this *LogicService) Route(conn proto.StreamLogicService_RouteServer) error {
	for {
		stream, err := conn.Recv()
		if io.EOF == err {
			glog.Info("[gRPC] Server Got EOF")
			return nil
		}

		if nil != err {
			glog.Error("[gRPC] Server Error ", err)
			return err
		}

		glog.Info("[gRPC] Server Recv: ", stream.Token)

		conn.Send(&proto.LogicResponse{
			MInfo: &proto.ConnectRoomInfo{
				Ip:   1,
				Port: 65,
			},
		})
	}

}

type RoomService struct {
}

func (this *RoomService) Route(conn proto.StreamRoomService_RouteServer) error {
	for {
		stream, err := conn.Recv()
		if io.EOF == err {
			glog.Info("[gRPC] Server Got EOF")
			return nil
		}

		if nil != err {
			glog.Error("[gRPC] Server Error ", err)
			return err
		}

		glog.Info("[gRPC] Server Recv: ", stream.Data)

		switch stream.Type {
		case proto.MsgType_Regist:
			var info proto.ConnectRoomInfo
			err := json.Unmarshal(stream.Data, &info)
			if nil != err {
				glog.Error("[Common] Json to Struct Error, ", err)
				return err
			}
			fmt.Println("Server Got Regist Msg ", info.Ip, ",", info.Port)
			break
		case proto.MsgType_Update:
			fmt.Println("Server Got Update Msg")
			break
		}

	}
}

func StartGrpcServer() bool {

	addr := env.Get("rcenter", "grpc")
	listen, err := net.Listen("tcp", addr)
	if nil != err {
		glog.Error("[Start] Bind Port Error, Port=", addr, ",", err)
		return false
	}

	s := grpc.NewServer()
	proto.RegisterStreamLogicServiceServer(s, &LogicService{})
	proto.RegisterStreamRoomServiceServer(s, &RoomService{})
	go func() {
		s.Serve(listen)
	}()

	glog.Info("[gRPC] Start Server Success, ", s.GetServiceInfo())

	return true
}
