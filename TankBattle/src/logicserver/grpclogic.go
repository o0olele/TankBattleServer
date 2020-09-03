package main

import (
	"base/env"
	"context"

	proto "proto"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

func GetVailabelRoomInfo(token string) (*proto.ConnectRoomInfo, error) {
	conn, err := grpc.Dial(env.Get("logic", "grpc"), grpc.WithInsecure())
	if nil != err {
		glog.Error("[gRPC] Connect Fail ", err)
		return nil, err
	}

	client := proto.NewStreamLogicServiceClient(conn)

	result, err := client.Route(context.Background(), &proto.LogicRequest{Token: token})
	if nil != err {
		glog.Error("[gRPC] Client Call Error ", err)
		return nil, err
	}

	return result.MInfo, nil
}

// type RPCClient struct {
// 	client  *rpc.Client
// 	isclose int32
// 	addr    string
// }

// var mRPCClient *RPCClient

// func RPCClient_GetMe() *RPCClient {
// 	if nil == mRPCClient {
// 		mRPCClient = &RPCClient{
// 			isclose: 1,
// 		}
// 	}

// 	return mRPCClient
// }

// func (this *RPCClient) Init(addr string) bool {
// 	if "" == addr {
// 		return false
// 	}

// 	this.addr = addr
// 	return this.Connect()
// }

// func (this *RPCClient) Connect() bool {
// 	if 0 == atomic.LoadInt32(&this.isclose) {
// 		return false
// 	}

// 	client, err := rpc.Dial("tcp", this.addr)
// 	if nil != err {
// 		return false
// 	}

// 	if nil != client {
// 		this.client = client
// 	}

// 	atomic.StoreInt32(&this.isclose, 0)
// 	return true
// }

// func (this *RPCClient) ReConnect() bool {
// 	if atomic.CompareAndSwapInt32(&this.isclose, 0, 1) {
// 		for {
// 			if this.Connect() {
// 				break
// 			}
// 			time.Sleep(time.Second)
// 		}
// 		return true
// 	}
// 	return false
// }

// func (this *RPCClient) RemoteCall(name string, args interface{}, reply interface{}) error {
// 	if nil == this.client {
// 		return errors.New("error")
// 	}

// 	err := this.client.Call(name, args, reply)
// 	if nil != err {
// 		if this.ReConnect() {
// 			return this.client.Call(name, args, reply)
// 		} else {
// 			for i := 0; i < 30; i++ {
// 				if nil == this.client {
// 					continue
// 				}

// 				err = this.client.Call(name, args, reply)
// 				if err != rpc.ErrShutdown {
// 					break
// 				}
// 			}
// 		}

// 		return err
// 	}

// 	return nil
// }
