package main

import (
	"base/env"
	proto "grpc"
	"io"
	"net"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type LogicService struct {
}

func (this *LogicService) Route(conn proto.StreamLogicService_RouteServer) error {
	for {
		stream, err := conn.Recv()
		if io.EOF == err {
			return nil
		}

		if nil != err {
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

func StartGrpcServer() bool {

	addr := env.Get("rcenter", "grpc")
	listen, err := net.Listen("tcp", addr)
	if nil != err {
		glog.Error("[Start] Bind Port Error, Port=", addr, ",", err)
		return false
	}

	s := grpc.NewServer()
	proto.RegisterStreamLogicServiceServer(s, &LogicService{})
	go func() {
		s.Serve(listen)
	}()

	glog.Info("[gRPC] Start Server Success, ", s.GetServiceInfo())

	for {
		time.Sleep(time.Second)
	}

	return true
}
