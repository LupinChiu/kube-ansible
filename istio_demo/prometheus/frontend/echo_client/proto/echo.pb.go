// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: echo.proto

package proto

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

type EchoReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Datetime string `protobuf:"bytes,1,opt,name=datetime,proto3" json:"datetime,omitempty"`
}

func (x *EchoReq) Reset() {
	*x = EchoReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoReq) ProtoMessage() {}

func (x *EchoReq) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoReq.ProtoReflect.Descriptor instead.
func (*EchoReq) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{0}
}

func (x *EchoReq) GetDatetime() string {
	if x != nil {
		return x.Datetime
	}
	return ""
}

type EchoRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Datetime string `protobuf:"bytes,2,opt,name=datetime,proto3" json:"datetime,omitempty"`
}

func (x *EchoRes) Reset() {
	*x = EchoRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_echo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRes) ProtoMessage() {}

func (x *EchoRes) ProtoReflect() protoreflect.Message {
	mi := &file_echo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoRes.ProtoReflect.Descriptor instead.
func (*EchoRes) Descriptor() ([]byte, []int) {
	return file_echo_proto_rawDescGZIP(), []int{1}
}

func (x *EchoRes) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EchoRes) GetDatetime() string {
	if x != nil {
		return x.Datetime
	}
	return ""
}

var File_echo_proto protoreflect.FileDescriptor

var file_echo_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x65, 0x63, 0x68, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x07, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x12, 0x1a,
	0x0a, 0x08, 0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x39, 0x0a, 0x07, 0x45, 0x63,
	0x68, 0x6f, 0x52, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74,
	0x65, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74,
	0x65, 0x74, 0x69, 0x6d, 0x65, 0x32, 0x30, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x28, 0x0a,
	0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x63,
	0x68, 0x6f, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x63,
	0x68, 0x6f, 0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_echo_proto_rawDescOnce sync.Once
	file_echo_proto_rawDescData = file_echo_proto_rawDesc
)

func file_echo_proto_rawDescGZIP() []byte {
	file_echo_proto_rawDescOnce.Do(func() {
		file_echo_proto_rawDescData = protoimpl.X.CompressGZIP(file_echo_proto_rawDescData)
	})
	return file_echo_proto_rawDescData
}

var file_echo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_echo_proto_goTypes = []interface{}{
	(*EchoReq)(nil), // 0: proto.EchoReq
	(*EchoRes)(nil), // 1: proto.EchoRes
}
var file_echo_proto_depIdxs = []int32{
	0, // 0: proto.Echo.Send:input_type -> proto.EchoReq
	1, // 1: proto.Echo.Send:output_type -> proto.EchoRes
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_echo_proto_init() }
func file_echo_proto_init() {
	if File_echo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_echo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoReq); i {
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
		file_echo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoRes); i {
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
			RawDescriptor: file_echo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_echo_proto_goTypes,
		DependencyIndexes: file_echo_proto_depIdxs,
		MessageInfos:      file_echo_proto_msgTypes,
	}.Build()
	File_echo_proto = out.File
	file_echo_proto_rawDesc = nil
	file_echo_proto_goTypes = nil
	file_echo_proto_depIdxs = nil
}
