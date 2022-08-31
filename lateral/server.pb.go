// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: server.proto

package lateral

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type QueueMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupIDs []*QueueMessage_Query `protobuf:"bytes,1,rep,name=GroupIDs,proto3" json:"GroupIDs,omitempty"`
}

func (x *QueueMessage) Reset() {
	*x = QueueMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueMessage) ProtoMessage() {}

func (x *QueueMessage) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueMessage.ProtoReflect.Descriptor instead.
func (*QueueMessage) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{0}
}

func (x *QueueMessage) GetGroupIDs() []*QueueMessage_Query {
	if x != nil {
		return x.GroupIDs
	}
	return nil
}

type QueueMessage_Query struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UC       string               `protobuf:"bytes,1,opt,name=UC,proto3" json:"UC,omitempty"`
	TS       *timestamp.Timestamp `protobuf:"bytes,2,opt,name=TS,proto3" json:"TS,omitempty"`
	Instance string               `protobuf:"bytes,3,opt,name=Instance,proto3" json:"Instance,omitempty"`
	Query    string               `protobuf:"bytes,4,opt,name=Query,proto3" json:"Query,omitempty"`
	Confirms []byte               `protobuf:"bytes,5,opt,name=Confirms,proto3" json:"Confirms,omitempty"`
}

func (x *QueueMessage_Query) Reset() {
	*x = QueueMessage_Query{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueMessage_Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueMessage_Query) ProtoMessage() {}

func (x *QueueMessage_Query) ProtoReflect() protoreflect.Message {
	mi := &file_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueMessage_Query.ProtoReflect.Descriptor instead.
func (*QueueMessage_Query) Descriptor() ([]byte, []int) {
	return file_server_proto_rawDescGZIP(), []int{0, 0}
}

func (x *QueueMessage_Query) GetUC() string {
	if x != nil {
		return x.UC
	}
	return ""
}

func (x *QueueMessage_Query) GetTS() *timestamp.Timestamp {
	if x != nil {
		return x.TS
	}
	return nil
}

func (x *QueueMessage_Query) GetInstance() string {
	if x != nil {
		return x.Instance
	}
	return ""
}

func (x *QueueMessage_Query) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *QueueMessage_Query) GetConfirms() []byte {
	if x != nil {
		return x.Confirms
	}
	return nil
}

var File_server_proto protoreflect.FileDescriptor

var file_server_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x6c, 0x61, 0x74, 0x65, 0x72, 0x61, 0x6c, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x01, 0x0a, 0x0c, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a, 0x08, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x49, 0x44, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x61,
	0x74, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x08, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49,
	0x44, 0x73, 0x1a, 0x91, 0x01, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02,
	0x55, 0x43, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x55, 0x43, 0x12, 0x2a, 0x0a, 0x02,
	0x54, 0x53, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x54, 0x53, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x73, 0x32, 0x44, 0x0a, 0x09, 0x53, 0x79, 0x6e, 0x63, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x12, 0x37, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x15, 0x2e, 0x6c,
	0x61, 0x74, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x1a, 0x15, 0x2e, 0x6c, 0x61, 0x74, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x42, 0x25, 0x5a, 0x23,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x6c, 0x69, 0x67, 0x6f,
	0x64, 0x65, 0x6e, 0x2f, 0x63, 0x68, 0x61, 0x73, 0x73, 0x69, 0x73, 0x2f, 0x6c, 0x61, 0x74, 0x65,
	0x72, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_server_proto_rawDescOnce sync.Once
	file_server_proto_rawDescData = file_server_proto_rawDesc
)

func file_server_proto_rawDescGZIP() []byte {
	file_server_proto_rawDescOnce.Do(func() {
		file_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_proto_rawDescData)
	})
	return file_server_proto_rawDescData
}

var file_server_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_server_proto_goTypes = []interface{}{
	(*QueueMessage)(nil),        // 0: lateral.QueueMessage
	(*QueueMessage_Query)(nil),  // 1: lateral.QueueMessage.Query
	(*timestamp.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_server_proto_depIdxs = []int32{
	1, // 0: lateral.QueueMessage.GroupIDs:type_name -> lateral.QueueMessage.Query
	2, // 1: lateral.QueueMessage.Query.TS:type_name -> google.protobuf.Timestamp
	0, // 2: lateral.SyncQueue.Queue:input_type -> lateral.QueueMessage
	0, // 3: lateral.SyncQueue.Queue:output_type -> lateral.QueueMessage
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_server_proto_init() }
func file_server_proto_init() {
	if File_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueMessage); i {
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
		file_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueMessage_Query); i {
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
			RawDescriptor: file_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_proto_goTypes,
		DependencyIndexes: file_server_proto_depIdxs,
		MessageInfos:      file_server_proto_msgTypes,
	}.Build()
	File_server_proto = out.File
	file_server_proto_rawDesc = nil
	file_server_proto_goTypes = nil
	file_server_proto_depIdxs = nil
}
