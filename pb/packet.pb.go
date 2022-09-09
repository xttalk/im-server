// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.5
// source: pb/packet.proto

package pb

import (
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

// 数据包命令
type Packet int32

const (
	Packet_PacketNone   Packet = 0   //不可用
	Packet_SystemError  Packet = 500 //系统错误
	Packet_UserNotLogin Packet = 403 //用户未登录
	// 未登录
	Packet_Login Packet = 1001 //登录鉴权
	// 个人信息
	Packet_GetProfile    Packet = 2001 //获取当前账号信息
	Packet_ModifyProfile Packet = 2002 //修改当前账号信息
	// 用户
	Packet_GetUser Packet = 3001 //获取用户信息
	// 好友
	Packet_GetFriendList Packet = 4001 //获取好友列表
	Packet_GetFriend     Packet = 4002 //获取好友信息
	Packet_RemoveFriend  Packet = 4003 //删除好友
	// 消息发送
	Packet_PrivateMsg     Packet = 8001 //私聊消息
	Packet_PrivateMsgAck  Packet = 8002 //私聊消息已送达
	Packet_PrivateMsgRead Packet = 8003 //私聊消息已读
)

// Enum value maps for Packet.
var (
	Packet_name = map[int32]string{
		0:    "PacketNone",
		500:  "SystemError",
		403:  "UserNotLogin",
		1001: "Login",
		2001: "GetProfile",
		2002: "ModifyProfile",
		3001: "GetUser",
		4001: "GetFriendList",
		4002: "GetFriend",
		4003: "RemoveFriend",
		8001: "PrivateMsg",
		8002: "PrivateMsgAck",
		8003: "PrivateMsgRead",
	}
	Packet_value = map[string]int32{
		"PacketNone":     0,
		"SystemError":    500,
		"UserNotLogin":   403,
		"Login":          1001,
		"GetProfile":     2001,
		"ModifyProfile":  2002,
		"GetUser":        3001,
		"GetFriendList":  4001,
		"GetFriend":      4002,
		"RemoveFriend":   4003,
		"PrivateMsg":     8001,
		"PrivateMsgAck":  8002,
		"PrivateMsgRead": 8003,
	}
)

func (x Packet) Enum() *Packet {
	p := new(Packet)
	*p = x
	return p
}

func (x Packet) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Packet) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_packet_proto_enumTypes[0].Descriptor()
}

func (Packet) Type() protoreflect.EnumType {
	return &file_pb_packet_proto_enumTypes[0]
}

func (x Packet) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Packet.Descriptor instead.
func (Packet) EnumDescriptor() ([]byte, []int) {
	return file_pb_packet_proto_rawDescGZIP(), []int{0}
}

// 状态码响应
type RetCode int32

const (
	RetCode_Unknown        RetCode = 0    //系统未知内部错误
	RetCode_Success        RetCode = 200  //请求成功,常规成功
	RetCode_Error          RetCode = 500  //系统错误
	RetCode_UserNotFound   RetCode = 3001 //没有找到用户
	RetCode_FriendNotFound RetCode = 4001 //不是好友关系
)

// Enum value maps for RetCode.
var (
	RetCode_name = map[int32]string{
		0:    "Unknown",
		200:  "Success",
		500:  "Error",
		3001: "UserNotFound",
		4001: "FriendNotFound",
	}
	RetCode_value = map[string]int32{
		"Unknown":        0,
		"Success":        200,
		"Error":          500,
		"UserNotFound":   3001,
		"FriendNotFound": 4001,
	}
)

func (x RetCode) Enum() *RetCode {
	p := new(RetCode)
	*p = x
	return p
}

func (x RetCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RetCode) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_packet_proto_enumTypes[1].Descriptor()
}

func (RetCode) Type() protoreflect.EnumType {
	return &file_pb_packet_proto_enumTypes[1]
}

func (x RetCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RetCode.Descriptor instead.
func (RetCode) EnumDescriptor() ([]byte, []int) {
	return file_pb_packet_proto_rawDescGZIP(), []int{1}
}

var File_pb_packet_proto protoreflect.FileDescriptor

var file_pb_packet_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x62, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x02, 0x70, 0x62, 0x2a, 0xed, 0x01, 0x0a, 0x06, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00,
	0x12, 0x10, 0x0a, 0x0b, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10,
	0xf4, 0x03, 0x12, 0x11, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x6f, 0x74, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x10, 0x93, 0x03, 0x12, 0x0a, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0xe9,
	0x07, 0x12, 0x0f, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x10,
	0xd1, 0x0f, 0x12, 0x12, 0x0a, 0x0d, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x10, 0xd2, 0x0f, 0x12, 0x0c, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x10, 0xb9, 0x17, 0x12, 0x12, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x46, 0x72, 0x69, 0x65, 0x6e,
	0x64, 0x4c, 0x69, 0x73, 0x74, 0x10, 0xa1, 0x1f, 0x12, 0x0e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x46,
	0x72, 0x69, 0x65, 0x6e, 0x64, 0x10, 0xa2, 0x1f, 0x12, 0x11, 0x0a, 0x0c, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x10, 0xa3, 0x1f, 0x12, 0x0f, 0x0a, 0x0a, 0x50,
	0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4d, 0x73, 0x67, 0x10, 0xc1, 0x3e, 0x12, 0x12, 0x0a, 0x0d,
	0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4d, 0x73, 0x67, 0x41, 0x63, 0x6b, 0x10, 0xc2, 0x3e,
	0x12, 0x13, 0x0a, 0x0e, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4d, 0x73, 0x67, 0x52, 0x65,
	0x61, 0x64, 0x10, 0xc3, 0x3e, 0x2a, 0x58, 0x0a, 0x07, 0x52, 0x65, 0x74, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x0c, 0x0a,
	0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x10, 0xc8, 0x01, 0x12, 0x0a, 0x0a, 0x05, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x10, 0xf4, 0x03, 0x12, 0x11, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x4e,
	0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xb9, 0x17, 0x12, 0x13, 0x0a, 0x0e, 0x46, 0x72,
	0x69, 0x65, 0x6e, 0x64, 0x4e, 0x6f, 0x74, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0xa1, 0x1f, 0x42,
	0x0e, 0x5a, 0x0c, 0x2e, 0x2e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_packet_proto_rawDescOnce sync.Once
	file_pb_packet_proto_rawDescData = file_pb_packet_proto_rawDesc
)

func file_pb_packet_proto_rawDescGZIP() []byte {
	file_pb_packet_proto_rawDescOnce.Do(func() {
		file_pb_packet_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_packet_proto_rawDescData)
	})
	return file_pb_packet_proto_rawDescData
}

var file_pb_packet_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_pb_packet_proto_goTypes = []interface{}{
	(Packet)(0),  // 0: pb.Packet
	(RetCode)(0), // 1: pb.RetCode
}
var file_pb_packet_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_packet_proto_init() }
func file_pb_packet_proto_init() {
	if File_pb_packet_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pb_packet_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_packet_proto_goTypes,
		DependencyIndexes: file_pb_packet_proto_depIdxs,
		EnumInfos:         file_pb_packet_proto_enumTypes,
	}.Build()
	File_pb_packet_proto = out.File
	file_pb_packet_proto_rawDesc = nil
	file_pb_packet_proto_goTypes = nil
	file_pb_packet_proto_depIdxs = nil
}
