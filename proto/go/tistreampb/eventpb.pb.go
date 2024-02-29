// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.8.0
// source: eventpb.proto

package __

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

type OpType int32

const (
	OpType_UNKNOWN  OpType = 0
	OpType_PREWRITE OpType = 1
	OpType_COMMIT   OpType = 2
	OpType_ROLLBACK OpType = 3
)

// Enum value maps for OpType.
var (
	OpType_name = map[int32]string{
		0: "UNKNOWN",
		1: "PREWRITE",
		2: "COMMIT",
		3: "ROLLBACK",
	}
	OpType_value = map[string]int32{
		"UNKNOWN":  0,
		"PREWRITE": 1,
		"COMMIT":   2,
		"ROLLBACK": 3,
	}
)

func (x OpType) Enum() *OpType {
	p := new(OpType)
	*p = x
	return p
}

func (x OpType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpType) Descriptor() protoreflect.EnumDescriptor {
	return file_eventpb_proto_enumTypes[0].Descriptor()
}

func (OpType) Type() protoreflect.EnumType {
	return &file_eventpb_proto_enumTypes[0]
}

func (x OpType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpType.Descriptor instead.
func (OpType) EnumDescriptor() ([]byte, []int) {
	return file_eventpb_proto_rawDescGZIP(), []int{0}
}

type EventRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTs  uint64 `protobuf:"varint,1,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	CommitTs uint64 `protobuf:"varint,2,opt,name=commit_ts,json=commitTs,proto3" json:"commit_ts,omitempty"`
	OpType   OpType `protobuf:"varint,3,opt,name=op_type,json=opType,proto3,enum=tistreampb.OpType" json:"op_type,omitempty"`
	Key      []byte `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
	Value    []byte `protobuf:"bytes,5,opt,name=value,proto3" json:"value,omitempty"`
	OldValue []byte `protobuf:"bytes,6,opt,name=old_value,json=oldValue,proto3" json:"old_value,omitempty"`
}

func (x *EventRow) Reset() {
	*x = EventRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRow) ProtoMessage() {}

func (x *EventRow) ProtoReflect() protoreflect.Message {
	mi := &file_eventpb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRow.ProtoReflect.Descriptor instead.
func (*EventRow) Descriptor() ([]byte, []int) {
	return file_eventpb_proto_rawDescGZIP(), []int{0}
}

func (x *EventRow) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *EventRow) GetCommitTs() uint64 {
	if x != nil {
		return x.CommitTs
	}
	return 0
}

func (x *EventRow) GetOpType() OpType {
	if x != nil {
		return x.OpType
	}
	return OpType_UNKNOWN
}

func (x *EventRow) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *EventRow) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *EventRow) GetOldValue() []byte {
	if x != nil {
		return x.OldValue
	}
	return nil
}

type EventWatermark struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts           uint64 `protobuf:"varint,1,opt,name=ts,proto3" json:"ts,omitempty"`
	RangeId      uint64 `protobuf:"varint,2,opt,name=range_id,json=rangeId,proto3" json:"range_id,omitempty"`
	RangeVersion uint64 `protobuf:"varint,3,opt,name=range_version,json=rangeVersion,proto3" json:"range_version,omitempty"`
	RangeStart   []byte `protobuf:"bytes,4,opt,name=range_start,json=rangeStart,proto3" json:"range_start,omitempty"`
	RangeEnd     []byte `protobuf:"bytes,5,opt,name=range_end,json=rangeEnd,proto3" json:"range_end,omitempty"`
}

func (x *EventWatermark) Reset() {
	*x = EventWatermark{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventWatermark) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventWatermark) ProtoMessage() {}

func (x *EventWatermark) ProtoReflect() protoreflect.Message {
	mi := &file_eventpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventWatermark.ProtoReflect.Descriptor instead.
func (*EventWatermark) Descriptor() ([]byte, []int) {
	return file_eventpb_proto_rawDescGZIP(), []int{1}
}

func (x *EventWatermark) GetTs() uint64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *EventWatermark) GetRangeId() uint64 {
	if x != nil {
		return x.RangeId
	}
	return 0
}

func (x *EventWatermark) GetRangeVersion() uint64 {
	if x != nil {
		return x.RangeVersion
	}
	return 0
}

func (x *EventWatermark) GetRangeStart() []byte {
	if x != nil {
		return x.RangeStart
	}
	return nil
}

func (x *EventWatermark) GetRangeEnd() []byte {
	if x != nil {
		return x.RangeEnd
	}
	return nil
}

type EventBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rows       []*EventRow       `protobuf:"bytes,1,rep,name=rows,proto3" json:"rows,omitempty"`
	Watermarks []*EventWatermark `protobuf:"bytes,2,rep,name=watermarks,proto3" json:"watermarks,omitempty"`
}

func (x *EventBatch) Reset() {
	*x = EventBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eventpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventBatch) ProtoMessage() {}

func (x *EventBatch) ProtoReflect() protoreflect.Message {
	mi := &file_eventpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventBatch.ProtoReflect.Descriptor instead.
func (*EventBatch) Descriptor() ([]byte, []int) {
	return file_eventpb_proto_rawDescGZIP(), []int{2}
}

func (x *EventBatch) GetRows() []*EventRow {
	if x != nil {
		return x.Rows
	}
	return nil
}

func (x *EventBatch) GetWatermarks() []*EventWatermark {
	if x != nil {
		return x.Watermarks
	}
	return nil
}

var File_eventpb_proto protoreflect.FileDescriptor

var file_eventpb_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x22, 0xb4, 0x01, 0x0a, 0x08,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x77, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x5f, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x54, 0x73,
	0x12, 0x2b, 0x0a, 0x07, 0x6f, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e, 0x4f,
	0x70, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x6f, 0x70, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x6c, 0x64, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6f, 0x6c, 0x64, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x9e, 0x01, 0x0a, 0x0e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x57, 0x61, 0x74, 0x65,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x49, 0x64,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x72, 0x61, 0x6e, 0x67,
	0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x5f,
	0x65, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x72, 0x61, 0x6e, 0x67, 0x65,
	0x45, 0x6e, 0x64, 0x22, 0x72, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x28, 0x0a, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x6f, 0x77, 0x52, 0x04, 0x72, 0x6f, 0x77, 0x73, 0x12, 0x3a, 0x0a, 0x0a, 0x77,
	0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x57, 0x61, 0x74, 0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x0a, 0x77, 0x61, 0x74,
	0x65, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x2a, 0x3d, 0x0a, 0x06, 0x4f, 0x70, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0c,
	0x0a, 0x08, 0x50, 0x52, 0x45, 0x57, 0x52, 0x49, 0x54, 0x45, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06,
	0x43, 0x4f, 0x4d, 0x4d, 0x49, 0x54, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x4f, 0x4c, 0x4c,
	0x42, 0x41, 0x43, 0x4b, 0x10, 0x03, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_eventpb_proto_rawDescOnce sync.Once
	file_eventpb_proto_rawDescData = file_eventpb_proto_rawDesc
)

func file_eventpb_proto_rawDescGZIP() []byte {
	file_eventpb_proto_rawDescOnce.Do(func() {
		file_eventpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_eventpb_proto_rawDescData)
	})
	return file_eventpb_proto_rawDescData
}

var file_eventpb_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_eventpb_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_eventpb_proto_goTypes = []interface{}{
	(OpType)(0),            // 0: tistreampb.OpType
	(*EventRow)(nil),       // 1: tistreampb.EventRow
	(*EventWatermark)(nil), // 2: tistreampb.EventWatermark
	(*EventBatch)(nil),     // 3: tistreampb.EventBatch
}
var file_eventpb_proto_depIdxs = []int32{
	0, // 0: tistreampb.EventRow.op_type:type_name -> tistreampb.OpType
	1, // 1: tistreampb.EventBatch.rows:type_name -> tistreampb.EventRow
	2, // 2: tistreampb.EventBatch.watermarks:type_name -> tistreampb.EventWatermark
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_eventpb_proto_init() }
func file_eventpb_proto_init() {
	if File_eventpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_eventpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventRow); i {
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
		file_eventpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventWatermark); i {
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
		file_eventpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventBatch); i {
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
			RawDescriptor: file_eventpb_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_eventpb_proto_goTypes,
		DependencyIndexes: file_eventpb_proto_depIdxs,
		EnumInfos:         file_eventpb_proto_enumTypes,
		MessageInfos:      file_eventpb_proto_msgTypes,
	}.Build()
	File_eventpb_proto = out.File
	file_eventpb_proto_rawDesc = nil
	file_eventpb_proto_goTypes = nil
	file_eventpb_proto_depIdxs = nil
}
