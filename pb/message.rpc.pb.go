// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.5
// source: pb/message.rpc.proto

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

type MqMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data      []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`                             //消息数据
	SessionId string `protobuf:"bytes,2,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`  //设备标识
	FromId    uint64 `protobuf:"varint,3,opt,name=from_id,json=fromId,proto3" json:"from_id,omitempty"`          //发送方ID,可能为0,根据业务处理
	ReceiveId uint64 `protobuf:"varint,4,opt,name=receive_id,json=receiveId,proto3" json:"receive_id,omitempty"` //消息接收ID,可能为0,根据业务事件类型
}

func (x *MqMsg) Reset() {
	*x = MqMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_message_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MqMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MqMsg) ProtoMessage() {}

func (x *MqMsg) ProtoReflect() protoreflect.Message {
	mi := &file_pb_message_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MqMsg.ProtoReflect.Descriptor instead.
func (*MqMsg) Descriptor() ([]byte, []int) {
	return file_pb_message_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *MqMsg) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *MqMsg) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *MqMsg) GetFromId() uint64 {
	if x != nil {
		return x.FromId
	}
	return 0
}

func (x *MqMsg) GetReceiveId() uint64 {
	if x != nil {
		return x.ReceiveId
	}
	return 0
}

var File_pb_message_rpc_proto protoreflect.FileDescriptor

var file_pb_message_rpc_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0f, 0x70, 0x62, 0x2f, 0x70,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x72, 0x0a, 0x05, 0x4d,
	0x71, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x72, 0x6f, 0x6d, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x66, 0x72, 0x6f, 0x6d, 0x49, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x49, 0x64, 0x42,
	0x0e, 0x5a, 0x0c, 0x2e, 0x2e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x50,
	0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_message_rpc_proto_rawDescOnce sync.Once
	file_pb_message_rpc_proto_rawDescData = file_pb_message_rpc_proto_rawDesc
)

func file_pb_message_rpc_proto_rawDescGZIP() []byte {
	file_pb_message_rpc_proto_rawDescOnce.Do(func() {
		file_pb_message_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_message_rpc_proto_rawDescData)
	})
	return file_pb_message_rpc_proto_rawDescData
}

var file_pb_message_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pb_message_rpc_proto_goTypes = []interface{}{
	(*MqMsg)(nil), // 0: pb.MqMsg
}
var file_pb_message_rpc_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_message_rpc_proto_init() }
func file_pb_message_rpc_proto_init() {
	if File_pb_message_rpc_proto != nil {
		return
	}
	file_pb_packet_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pb_message_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MqMsg); i {
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
			RawDescriptor: file_pb_message_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_message_rpc_proto_goTypes,
		DependencyIndexes: file_pb_message_rpc_proto_depIdxs,
		MessageInfos:      file_pb_message_rpc_proto_msgTypes,
	}.Build()
	File_pb_message_rpc_proto = out.File
	file_pb_message_rpc_proto_rawDesc = nil
	file_pb_message_rpc_proto_goTypes = nil
	file_pb_message_rpc_proto_depIdxs = nil
}