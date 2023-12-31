// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: blockevent.proto

package block

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BlockEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Event     string                 `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Index     int64                  `protobuf:"varint,4,opt,name=index,proto3" json:"index,omitempty"`
	// Types that are assignable to BlockEvent:
	//
	//	*BlockEvent_BlockIndex
	BlockEvent isBlockEvent_BlockEvent `protobuf_oneof:"block_event"`
}

func (x *BlockEvent) Reset() {
	*x = BlockEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blockevent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockEvent) ProtoMessage() {}

func (x *BlockEvent) ProtoReflect() protoreflect.Message {
	mi := &file_blockevent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockEvent.ProtoReflect.Descriptor instead.
func (*BlockEvent) Descriptor() ([]byte, []int) {
	return file_blockevent_proto_rawDescGZIP(), []int{0}
}

func (x *BlockEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *BlockEvent) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *BlockEvent) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *BlockEvent) GetIndex() int64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (m *BlockEvent) GetBlockEvent() isBlockEvent_BlockEvent {
	if m != nil {
		return m.BlockEvent
	}
	return nil
}

func (x *BlockEvent) GetBlockIndex() *BlockIndex {
	if x, ok := x.GetBlockEvent().(*BlockEvent_BlockIndex); ok {
		return x.BlockIndex
	}
	return nil
}

type isBlockEvent_BlockEvent interface {
	isBlockEvent_BlockEvent()
}

type BlockEvent_BlockIndex struct {
	BlockIndex *BlockIndex `protobuf:"bytes,8,opt,name=block_index,json=blockIndex,proto3,oneof"`
}

func (*BlockEvent_BlockIndex) isBlockEvent_BlockEvent() {}

var File_blockevent_proto protoreflect.FileDescriptor

var file_blockevent_proto_rawDesc = []byte{
	0x0a, 0x10, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x73, 0x74, 0x66, 0x78, 0x2e, 0x69, 0x6f, 0x2e, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd0, 0x01, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x3d, 0x0a, 0x0b, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x73, 0x74, 0x66, 0x78, 0x2e, 0x69, 0x6f, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x48, 0x00, 0x52, 0x0a, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x42, 0x0d, 0x0a, 0x0b, 0x62, 0x6c, 0x6f,
	0x63, 0x6b, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x15, 0x5a, 0x13, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x72, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blockevent_proto_rawDescOnce sync.Once
	file_blockevent_proto_rawDescData = file_blockevent_proto_rawDesc
)

func file_blockevent_proto_rawDescGZIP() []byte {
	file_blockevent_proto_rawDescOnce.Do(func() {
		file_blockevent_proto_rawDescData = protoimpl.X.CompressGZIP(file_blockevent_proto_rawDescData)
	})
	return file_blockevent_proto_rawDescData
}

var file_blockevent_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_blockevent_proto_goTypes = []interface{}{
	(*BlockEvent)(nil),            // 0: stfx.io.stream.BlockEvent
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*BlockIndex)(nil),            // 2: stfx.io.stream.BlockIndex
}
var file_blockevent_proto_depIdxs = []int32{
	1, // 0: stfx.io.stream.BlockEvent.timestamp:type_name -> google.protobuf.Timestamp
	2, // 1: stfx.io.stream.BlockEvent.block_index:type_name -> stfx.io.stream.BlockIndex
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_blockevent_proto_init() }
func file_blockevent_proto_init() {
	if File_blockevent_proto != nil {
		return
	}
	file_blockindex_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_blockevent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockEvent); i {
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
	file_blockevent_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*BlockEvent_BlockIndex)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_blockevent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_blockevent_proto_goTypes,
		DependencyIndexes: file_blockevent_proto_depIdxs,
		MessageInfos:      file_blockevent_proto_msgTypes,
	}.Build()
	File_blockevent_proto = out.File
	file_blockevent_proto_rawDesc = nil
	file_blockevent_proto_goTypes = nil
	file_blockevent_proto_depIdxs = nil
}
