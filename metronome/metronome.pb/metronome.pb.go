//******************************************************************************
//metronome.pb
//
//This package provides the protocol buffer source code for the metronome message
//that is sent periodically to each gorotuine that is connected to the message_q.
//
//Rev Date     By  Reason
//--- -------- --- ---------------------------------------------------------------
//001			 tdm  original
//
//*****************************************************************************

// command line to generate this file: protoc -I=. --go_out=. metronome.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: metronome.proto

package metronome_pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Tick struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Second        int32                  `protobuf:"varint,1,opt,name=second,proto3" json:"second,omitempty"`
	Minute        int32                  `protobuf:"varint,2,opt,name=minute,proto3" json:"minute,omitempty"`
	Seconda_2     int32                  `protobuf:"varint,20,opt,name=Seconda_2,json=Seconda2,proto3" json:"Seconda_2,omitempty"`
	Seconda_3     int32                  `protobuf:"varint,21,opt,name=Seconda_3,json=Seconda3,proto3" json:"Seconda_3,omitempty"`
	Seconda_4     int32                  `protobuf:"varint,22,opt,name=Seconda_4,json=Seconda4,proto3" json:"Seconda_4,omitempty"`
	Seconda_5     int32                  `protobuf:"varint,23,opt,name=Seconda_5,json=Seconda5,proto3" json:"Seconda_5,omitempty"`
	Seconda_6     int32                  `protobuf:"varint,24,opt,name=Seconda_6,json=Seconda6,proto3" json:"Seconda_6,omitempty"`
	Seconda_10    int32                  `protobuf:"varint,25,opt,name=Seconda_10,json=Seconda10,proto3" json:"Seconda_10,omitempty"`
	Seconda_12    int32                  `protobuf:"varint,26,opt,name=Seconda_12,json=Seconda12,proto3" json:"Seconda_12,omitempty"`
	Seconda_15    int32                  `protobuf:"varint,27,opt,name=Seconda_15,json=Seconda15,proto3" json:"Seconda_15,omitempty"`
	Seconda_20    int32                  `protobuf:"varint,28,opt,name=Seconda_20,json=Seconda20,proto3" json:"Seconda_20,omitempty"`
	Seconda_30    int32                  `protobuf:"varint,29,opt,name=Seconda_30,json=Seconda30,proto3" json:"Seconda_30,omitempty"`
	Minutes_1     int32                  `protobuf:"varint,30,opt,name=Minutes_1,json=Minutes1,proto3" json:"Minutes_1,omitempty"`
	Minutes_2     int32                  `protobuf:"varint,31,opt,name=Minutes_2,json=Minutes2,proto3" json:"Minutes_2,omitempty"`
	Minutes_3     int32                  `protobuf:"varint,32,opt,name=Minutes_3,json=Minutes3,proto3" json:"Minutes_3,omitempty"`
	Minutes_4     int32                  `protobuf:"varint,33,opt,name=Minutes_4,json=Minutes4,proto3" json:"Minutes_4,omitempty"`
	Minutes_5     int32                  `protobuf:"varint,34,opt,name=Minutes_5,json=Minutes5,proto3" json:"Minutes_5,omitempty"`
	Minutes_6     int32                  `protobuf:"varint,35,opt,name=Minutes_6,json=Minutes6,proto3" json:"Minutes_6,omitempty"`
	Minutes_10    int32                  `protobuf:"varint,36,opt,name=Minutes_10,json=Minutes10,proto3" json:"Minutes_10,omitempty"`
	Minutes_12    int32                  `protobuf:"varint,37,opt,name=Minutes_12,json=Minutes12,proto3" json:"Minutes_12,omitempty"`
	Minutes_15    int32                  `protobuf:"varint,38,opt,name=Minutes_15,json=Minutes15,proto3" json:"Minutes_15,omitempty"`
	Minutes_20    int32                  `protobuf:"varint,39,opt,name=Minutes_20,json=Minutes20,proto3" json:"Minutes_20,omitempty"`
	Minutes_30    int32                  `protobuf:"varint,40,opt,name=Minutes_30,json=Minutes30,proto3" json:"Minutes_30,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tick) Reset() {
	*x = Tick{}
	mi := &file_metronome_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tick) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tick) ProtoMessage() {}

func (x *Tick) ProtoReflect() protoreflect.Message {
	mi := &file_metronome_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tick.ProtoReflect.Descriptor instead.
func (*Tick) Descriptor() ([]byte, []int) {
	return file_metronome_proto_rawDescGZIP(), []int{0}
}

func (x *Tick) GetSecond() int32 {
	if x != nil {
		return x.Second
	}
	return 0
}

func (x *Tick) GetMinute() int32 {
	if x != nil {
		return x.Minute
	}
	return 0
}

func (x *Tick) GetSeconda_2() int32 {
	if x != nil {
		return x.Seconda_2
	}
	return 0
}

func (x *Tick) GetSeconda_3() int32 {
	if x != nil {
		return x.Seconda_3
	}
	return 0
}

func (x *Tick) GetSeconda_4() int32 {
	if x != nil {
		return x.Seconda_4
	}
	return 0
}

func (x *Tick) GetSeconda_5() int32 {
	if x != nil {
		return x.Seconda_5
	}
	return 0
}

func (x *Tick) GetSeconda_6() int32 {
	if x != nil {
		return x.Seconda_6
	}
	return 0
}

func (x *Tick) GetSeconda_10() int32 {
	if x != nil {
		return x.Seconda_10
	}
	return 0
}

func (x *Tick) GetSeconda_12() int32 {
	if x != nil {
		return x.Seconda_12
	}
	return 0
}

func (x *Tick) GetSeconda_15() int32 {
	if x != nil {
		return x.Seconda_15
	}
	return 0
}

func (x *Tick) GetSeconda_20() int32 {
	if x != nil {
		return x.Seconda_20
	}
	return 0
}

func (x *Tick) GetSeconda_30() int32 {
	if x != nil {
		return x.Seconda_30
	}
	return 0
}

func (x *Tick) GetMinutes_1() int32 {
	if x != nil {
		return x.Minutes_1
	}
	return 0
}

func (x *Tick) GetMinutes_2() int32 {
	if x != nil {
		return x.Minutes_2
	}
	return 0
}

func (x *Tick) GetMinutes_3() int32 {
	if x != nil {
		return x.Minutes_3
	}
	return 0
}

func (x *Tick) GetMinutes_4() int32 {
	if x != nil {
		return x.Minutes_4
	}
	return 0
}

func (x *Tick) GetMinutes_5() int32 {
	if x != nil {
		return x.Minutes_5
	}
	return 0
}

func (x *Tick) GetMinutes_6() int32 {
	if x != nil {
		return x.Minutes_6
	}
	return 0
}

func (x *Tick) GetMinutes_10() int32 {
	if x != nil {
		return x.Minutes_10
	}
	return 0
}

func (x *Tick) GetMinutes_12() int32 {
	if x != nil {
		return x.Minutes_12
	}
	return 0
}

func (x *Tick) GetMinutes_15() int32 {
	if x != nil {
		return x.Minutes_15
	}
	return 0
}

func (x *Tick) GetMinutes_20() int32 {
	if x != nil {
		return x.Minutes_20
	}
	return 0
}

func (x *Tick) GetMinutes_30() int32 {
	if x != nil {
		return x.Minutes_30
	}
	return 0
}

var File_metronome_proto protoreflect.FileDescriptor

var file_metronome_proto_rawDesc = string([]byte{
	0x0a, 0x0f, 0x6d, 0x65, 0x74, 0x72, 0x6f, 0x6e, 0x6f, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x6d, 0x65, 0x74, 0x72, 0x6f, 0x6e, 0x6f, 0x6d, 0x65, 0x2e, 0x70, 0x62, 0x22,
	0xab, 0x05, 0x0a, 0x04, 0x74, 0x69, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x53, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x61, 0x5f, 0x32, 0x18, 0x14, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x61, 0x32, 0x12, 0x1b, 0x0a, 0x09, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61,
	0x5f, 0x33, 0x18, 0x15, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x61, 0x33, 0x12, 0x1b, 0x0a, 0x09, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x5f, 0x34, 0x18,
	0x16, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x34, 0x12,
	0x1b, 0x0a, 0x09, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x5f, 0x35, 0x18, 0x17, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x35, 0x12, 0x1b, 0x0a, 0x09,
	0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x5f, 0x36, 0x18, 0x18, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x36, 0x12, 0x1d, 0x0a, 0x0a, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x61, 0x5f, 0x31, 0x30, 0x18, 0x19, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x53,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x31, 0x30, 0x12, 0x1d, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x61, 0x5f, 0x31, 0x32, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x53, 0x65,
	0x63, 0x6f, 0x6e, 0x64, 0x61, 0x31, 0x32, 0x12, 0x1d, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x61, 0x5f, 0x31, 0x35, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x61, 0x31, 0x35, 0x12, 0x1d, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x61, 0x5f, 0x32, 0x30, 0x18, 0x1c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x53, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x61, 0x32, 0x30, 0x12, 0x1d, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61,
	0x5f, 0x33, 0x30, 0x18, 0x1d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x53, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x61, 0x33, 0x30, 0x12, 0x1b, 0x0a, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f,
	0x31, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73,
	0x31, 0x12, 0x1b, 0x0a, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x32, 0x18, 0x1f,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x32, 0x12, 0x1b,
	0x0a, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x33, 0x18, 0x20, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x08, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x33, 0x12, 0x1b, 0x0a, 0x09, 0x4d,
	0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x34, 0x18, 0x21, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x34, 0x12, 0x1b, 0x0a, 0x09, 0x4d, 0x69, 0x6e, 0x75,
	0x74, 0x65, 0x73, 0x5f, 0x35, 0x18, 0x22, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x4d, 0x69, 0x6e,
	0x75, 0x74, 0x65, 0x73, 0x35, 0x12, 0x1b, 0x0a, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73,
	0x5f, 0x36, 0x18, 0x23, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65,
	0x73, 0x36, 0x12, 0x1d, 0x0a, 0x0a, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x31, 0x30,
	0x18, 0x24, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x31,
	0x30, 0x12, 0x1d, 0x0a, 0x0a, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x31, 0x32, 0x18,
	0x25, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x31, 0x32,
	0x12, 0x1d, 0x0a, 0x0a, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x31, 0x35, 0x18, 0x26,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x31, 0x35, 0x12,
	0x1d, 0x0a, 0x0a, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x32, 0x30, 0x18, 0x27, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x32, 0x30, 0x12, 0x1d,
	0x0a, 0x0a, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x33, 0x30, 0x18, 0x28, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x73, 0x33, 0x30, 0x42, 0x10, 0x5a,
	0x0e, 0x2e, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x6f, 0x6e, 0x6f, 0x6d, 0x65, 0x2e, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_metronome_proto_rawDescOnce sync.Once
	file_metronome_proto_rawDescData []byte
)

func file_metronome_proto_rawDescGZIP() []byte {
	file_metronome_proto_rawDescOnce.Do(func() {
		file_metronome_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_metronome_proto_rawDesc), len(file_metronome_proto_rawDesc)))
	})
	return file_metronome_proto_rawDescData
}

var file_metronome_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_metronome_proto_goTypes = []any{
	(*Tick)(nil), // 0: metronome.pb.tick
}
var file_metronome_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_metronome_proto_init() }
func file_metronome_proto_init() {
	if File_metronome_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_metronome_proto_rawDesc), len(file_metronome_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_metronome_proto_goTypes,
		DependencyIndexes: file_metronome_proto_depIdxs,
		MessageInfos:      file_metronome_proto_msgTypes,
	}.Build()
	File_metronome_proto = out.File
	file_metronome_proto_goTypes = nil
	file_metronome_proto_depIdxs = nil
}
