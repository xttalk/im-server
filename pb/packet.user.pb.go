// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.5
// source: pb/packet.user.proto

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

type PacketGetUserInfoReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"` //用户名
}

func (x *PacketGetUserInfoReq) Reset() {
	*x = PacketGetUserInfoReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_packet_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PacketGetUserInfoReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketGetUserInfoReq) ProtoMessage() {}

func (x *PacketGetUserInfoReq) ProtoReflect() protoreflect.Message {
	mi := &file_pb_packet_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketGetUserInfoReq.ProtoReflect.Descriptor instead.
func (*PacketGetUserInfoReq) Descriptor() ([]byte, []int) {
	return file_pb_packet_user_proto_rawDescGZIP(), []int{0}
}

func (x *PacketGetUserInfoReq) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type PacketGetUserInfoResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RetCode RetCode `protobuf:"varint,1,opt,name=ret_code,json=retCode,proto3,enum=pb.RetCode" json:"ret_code,omitempty"` //响应状态码
	User    *User   `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`                                       //用户信息
}

func (x *PacketGetUserInfoResp) Reset() {
	*x = PacketGetUserInfoResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_packet_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PacketGetUserInfoResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PacketGetUserInfoResp) ProtoMessage() {}

func (x *PacketGetUserInfoResp) ProtoReflect() protoreflect.Message {
	mi := &file_pb_packet_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PacketGetUserInfoResp.ProtoReflect.Descriptor instead.
func (*PacketGetUserInfoResp) Descriptor() ([]byte, []int) {
	return file_pb_packet_user_proto_rawDescGZIP(), []int{1}
}

func (x *PacketGetUserInfoResp) GetRetCode() RetCode {
	if x != nil {
		return x.RetCode
	}
	return RetCode_Unknown
}

func (x *PacketGetUserInfoResp) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_pb_packet_user_proto protoreflect.FileDescriptor

var file_pb_packet_user_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x62, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0f, 0x70, 0x62, 0x2f, 0x70,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x70, 0x62, 0x2f,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x14, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x47, 0x65, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x5d, 0x0a, 0x15, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x26, 0x0a, 0x08, 0x72, 0x65, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x52,
	0x07, 0x72, 0x65, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x2e, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x50, 0x00, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pb_packet_user_proto_rawDescOnce sync.Once
	file_pb_packet_user_proto_rawDescData = file_pb_packet_user_proto_rawDesc
)

func file_pb_packet_user_proto_rawDescGZIP() []byte {
	file_pb_packet_user_proto_rawDescOnce.Do(func() {
		file_pb_packet_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_packet_user_proto_rawDescData)
	})
	return file_pb_packet_user_proto_rawDescData
}

var file_pb_packet_user_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pb_packet_user_proto_goTypes = []interface{}{
	(*PacketGetUserInfoReq)(nil),  // 0: pb.PacketGetUserInfoReq
	(*PacketGetUserInfoResp)(nil), // 1: pb.PacketGetUserInfoResp
	(RetCode)(0),                  // 2: pb.RetCode
	(*User)(nil),                  // 3: pb.User
}
var file_pb_packet_user_proto_depIdxs = []int32{
	2, // 0: pb.PacketGetUserInfoResp.ret_code:type_name -> pb.RetCode
	3, // 1: pb.PacketGetUserInfoResp.user:type_name -> pb.User
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pb_packet_user_proto_init() }
func file_pb_packet_user_proto_init() {
	if File_pb_packet_user_proto != nil {
		return
	}
	file_pb_packet_proto_init()
	file_pb_packet_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pb_packet_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PacketGetUserInfoReq); i {
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
		file_pb_packet_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PacketGetUserInfoResp); i {
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
			RawDescriptor: file_pb_packet_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_packet_user_proto_goTypes,
		DependencyIndexes: file_pb_packet_user_proto_depIdxs,
		MessageInfos:      file_pb_packet_user_proto_msgTypes,
	}.Build()
	File_pb_packet_user_proto = out.File
	file_pb_packet_user_proto_rawDesc = nil
	file_pb_packet_user_proto_goTypes = nil
	file_pb_packet_user_proto_depIdxs = nil
}