// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.10.1
// source: proto/server.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type MsgType int32

const (
	MsgType_Regist MsgType = 0 // 注册roomserver
	MsgType_Update MsgType = 1 // 更新roomserver信息
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0: "Regist",
		1: "Update",
	}
	MsgType_value = map[string]int32{
		"Regist": 0,
		"Update": 1,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_server_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_proto_server_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{0}
}

type ConnectRoomInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   uint32 `protobuf:"fixed32,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port uint32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *ConnectRoomInfo) Reset() {
	*x = ConnectRoomInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectRoomInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectRoomInfo) ProtoMessage() {}

func (x *ConnectRoomInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectRoomInfo.ProtoReflect.Descriptor instead.
func (*ConnectRoomInfo) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{0}
}

func (x *ConnectRoomInfo) GetIp() uint32 {
	if x != nil {
		return x.Ip
	}
	return 0
}

func (x *ConnectRoomInfo) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type LogicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MInfo *ConnectRoomInfo `protobuf:"bytes,1,opt,name=mInfo,proto3" json:"mInfo,omitempty"`
}

func (x *LogicResponse) Reset() {
	*x = LogicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogicResponse) ProtoMessage() {}

func (x *LogicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogicResponse.ProtoReflect.Descriptor instead.
func (*LogicResponse) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{1}
}

func (x *LogicResponse) GetMInfo() *ConnectRoomInfo {
	if x != nil {
		return x.MInfo
	}
	return nil
}

type LogicRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *LogicRequest) Reset() {
	*x = LogicRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogicRequest) ProtoMessage() {}

func (x *LogicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogicRequest.ProtoReflect.Descriptor instead.
func (*LogicRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{2}
}

func (x *LogicRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RoomResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RoomResponse) Reset() {
	*x = RoomResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomResponse) ProtoMessage() {}

func (x *RoomResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomResponse.ProtoReflect.Descriptor instead.
func (*RoomResponse) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{3}
}

type RoomRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.MsgType" json:"type,omitempty"`
	Data []byte  `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *RoomRequest) Reset() {
	*x = RoomRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomRequest) ProtoMessage() {}

func (x *RoomRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomRequest.ProtoReflect.Descriptor instead.
func (*RoomRequest) Descriptor() ([]byte, []int) {
	return file_proto_server_proto_rawDescGZIP(), []int{4}
}

func (x *RoomRequest) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_Regist
}

func (x *RoomRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_server_proto protoreflect.FileDescriptor

var file_proto_server_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x35, 0x0a, 0x0f, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x07, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x22, 0x3d, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x6d, 0x49, 0x6e, 0x66,
	0x6f, 0x22, 0x24, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x0e, 0x0a, 0x0c, 0x52, 0x6f, 0x6f, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x45, 0x0a, 0x0b, 0x52, 0x6f, 0x6f, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x73, 0x67,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x21,
	0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x10,
	0x01, 0x32, 0x4a, 0x0a, 0x12, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4c, 0x6f, 0x67, 0x69, 0x63,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x63, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f,
	0x67, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x4b, 0x0a,
	0x11, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x6f, 0x6f, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x36, 0x0a, 0x05, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_server_proto_rawDescOnce sync.Once
	file_proto_server_proto_rawDescData = file_proto_server_proto_rawDesc
)

func file_proto_server_proto_rawDescGZIP() []byte {
	file_proto_server_proto_rawDescOnce.Do(func() {
		file_proto_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_server_proto_rawDescData)
	})
	return file_proto_server_proto_rawDescData
}

var file_proto_server_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_server_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_server_proto_goTypes = []interface{}{
	(MsgType)(0),            // 0: proto.MsgType
	(*ConnectRoomInfo)(nil), // 1: proto.ConnectRoomInfo
	(*LogicResponse)(nil),   // 2: proto.LogicResponse
	(*LogicRequest)(nil),    // 3: proto.LogicRequest
	(*RoomResponse)(nil),    // 4: proto.RoomResponse
	(*RoomRequest)(nil),     // 5: proto.RoomRequest
}
var file_proto_server_proto_depIdxs = []int32{
	1, // 0: proto.LogicResponse.mInfo:type_name -> proto.ConnectRoomInfo
	0, // 1: proto.RoomRequest.type:type_name -> proto.MsgType
	3, // 2: proto.StreamLogicService.Route:input_type -> proto.LogicRequest
	5, // 3: proto.StreamRoomService.Route:input_type -> proto.RoomRequest
	2, // 4: proto.StreamLogicService.Route:output_type -> proto.LogicResponse
	4, // 5: proto.StreamRoomService.Route:output_type -> proto.RoomResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_server_proto_init() }
func file_proto_server_proto_init() {
	if File_proto_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectRoomInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogicResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogicRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_server_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_proto_server_proto_goTypes,
		DependencyIndexes: file_proto_server_proto_depIdxs,
		EnumInfos:         file_proto_server_proto_enumTypes,
		MessageInfos:      file_proto_server_proto_msgTypes,
	}.Build()
	File_proto_server_proto = out.File
	file_proto_server_proto_rawDesc = nil
	file_proto_server_proto_goTypes = nil
	file_proto_server_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// StreamLogicServiceClient is the client API for StreamLogicService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StreamLogicServiceClient interface {
	Route(ctx context.Context, in *LogicRequest, opts ...grpc.CallOption) (*LogicResponse, error)
}

type streamLogicServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamLogicServiceClient(cc grpc.ClientConnInterface) StreamLogicServiceClient {
	return &streamLogicServiceClient{cc}
}

func (c *streamLogicServiceClient) Route(ctx context.Context, in *LogicRequest, opts ...grpc.CallOption) (*LogicResponse, error) {
	out := new(LogicResponse)
	err := c.cc.Invoke(ctx, "/proto.StreamLogicService/Route", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StreamLogicServiceServer is the server API for StreamLogicService service.
type StreamLogicServiceServer interface {
	Route(context.Context, *LogicRequest) (*LogicResponse, error)
}

// UnimplementedStreamLogicServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStreamLogicServiceServer struct {
}

func (*UnimplementedStreamLogicServiceServer) Route(context.Context, *LogicRequest) (*LogicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Route not implemented")
}

func RegisterStreamLogicServiceServer(s *grpc.Server, srv StreamLogicServiceServer) {
	s.RegisterService(&_StreamLogicService_serviceDesc, srv)
}

func _StreamLogicService_Route_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StreamLogicServiceServer).Route(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.StreamLogicService/Route",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StreamLogicServiceServer).Route(ctx, req.(*LogicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _StreamLogicService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.StreamLogicService",
	HandlerType: (*StreamLogicServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Route",
			Handler:    _StreamLogicService_Route_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/server.proto",
}

// StreamRoomServiceClient is the client API for StreamRoomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StreamRoomServiceClient interface {
	Route(ctx context.Context, opts ...grpc.CallOption) (StreamRoomService_RouteClient, error)
}

type streamRoomServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamRoomServiceClient(cc grpc.ClientConnInterface) StreamRoomServiceClient {
	return &streamRoomServiceClient{cc}
}

func (c *streamRoomServiceClient) Route(ctx context.Context, opts ...grpc.CallOption) (StreamRoomService_RouteClient, error) {
	stream, err := c.cc.NewStream(ctx, &_StreamRoomService_serviceDesc.Streams[0], "/proto.StreamRoomService/Route", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamRoomServiceRouteClient{stream}
	return x, nil
}

type StreamRoomService_RouteClient interface {
	Send(*RoomRequest) error
	Recv() (*RoomResponse, error)
	grpc.ClientStream
}

type streamRoomServiceRouteClient struct {
	grpc.ClientStream
}

func (x *streamRoomServiceRouteClient) Send(m *RoomRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamRoomServiceRouteClient) Recv() (*RoomResponse, error) {
	m := new(RoomResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamRoomServiceServer is the server API for StreamRoomService service.
type StreamRoomServiceServer interface {
	Route(StreamRoomService_RouteServer) error
}

// UnimplementedStreamRoomServiceServer can be embedded to have forward compatible implementations.
type UnimplementedStreamRoomServiceServer struct {
}

func (*UnimplementedStreamRoomServiceServer) Route(StreamRoomService_RouteServer) error {
	return status.Errorf(codes.Unimplemented, "method Route not implemented")
}

func RegisterStreamRoomServiceServer(s *grpc.Server, srv StreamRoomServiceServer) {
	s.RegisterService(&_StreamRoomService_serviceDesc, srv)
}

func _StreamRoomService_Route_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamRoomServiceServer).Route(&streamRoomServiceRouteServer{stream})
}

type StreamRoomService_RouteServer interface {
	Send(*RoomResponse) error
	Recv() (*RoomRequest, error)
	grpc.ServerStream
}

type streamRoomServiceRouteServer struct {
	grpc.ServerStream
}

func (x *streamRoomServiceRouteServer) Send(m *RoomResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamRoomServiceRouteServer) Recv() (*RoomRequest, error) {
	m := new(RoomRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _StreamRoomService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.StreamRoomService",
	HandlerType: (*StreamRoomServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Route",
			Handler:       _StreamRoomService_Route_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/server.proto",
}