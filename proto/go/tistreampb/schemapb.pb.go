// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.8.0
// source: schemapb.proto

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

type Table_OP int32

const (
	Table_CREATE Table_OP = 0
	Table_DROP   Table_OP = 1
	Table_ALTER  Table_OP = 2
)

// Enum value maps for Table_OP.
var (
	Table_OP_name = map[int32]string{
		0: "CREATE",
		1: "DROP",
		2: "ALTER",
	}
	Table_OP_value = map[string]int32{
		"CREATE": 0,
		"DROP":   1,
		"ALTER":  2,
	}
)

func (x Table_OP) Enum() *Table_OP {
	p := new(Table_OP)
	*p = x
	return p
}

func (x Table_OP) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Table_OP) Descriptor() protoreflect.EnumDescriptor {
	return file_schemapb_proto_enumTypes[0].Descriptor()
}

func (Table_OP) Type() protoreflect.EnumType {
	return &file_schemapb_proto_enumTypes[0]
}

func (x Table_OP) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Table_OP.Descriptor instead.
func (Table_OP) EnumDescriptor() ([]byte, []int) {
	return file_schemapb_proto_rawDescGZIP(), []int{0, 0}
}

type Schema_OP int32

const (
	Schema_CREATE Schema_OP = 0
	Schema_DROP   Schema_OP = 1
)

// Enum value maps for Schema_OP.
var (
	Schema_OP_name = map[int32]string{
		0: "CREATE",
		1: "DROP",
	}
	Schema_OP_value = map[string]int32{
		"CREATE": 0,
		"DROP":   1,
	}
)

func (x Schema_OP) Enum() *Schema_OP {
	p := new(Schema_OP)
	*p = x
	return p
}

func (x Schema_OP) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Schema_OP) Descriptor() protoreflect.EnumDescriptor {
	return file_schemapb_proto_enumTypes[1].Descriptor()
}

func (Schema_OP) Type() protoreflect.EnumType {
	return &file_schemapb_proto_enumTypes[1]
}

func (x Schema_OP) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Schema_OP.Descriptor instead.
func (Schema_OP) EnumDescriptor() ([]byte, []int) {
	return file_schemapb_proto_rawDescGZIP(), []int{1, 0}
}

type Table struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name            string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id              uint64   `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Op              Table_OP `protobuf:"varint,3,opt,name=op,proto3,enum=tistreampb.Table_OP" json:"op,omitempty"`
	Ts              uint64   `protobuf:"varint,4,opt,name=ts,proto3" json:"ts,omitempty"`
	Statement       string   `protobuf:"bytes,5,opt,name=statement,proto3" json:"statement,omitempty"`
	CreateStatement string   `protobuf:"bytes,6,opt,name=create_statement,json=createStatement,proto3" json:"create_statement,omitempty"`
}

func (x *Table) Reset() {
	*x = Table{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schemapb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Table) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Table) ProtoMessage() {}

func (x *Table) ProtoReflect() protoreflect.Message {
	mi := &file_schemapb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Table.ProtoReflect.Descriptor instead.
func (*Table) Descriptor() ([]byte, []int) {
	return file_schemapb_proto_rawDescGZIP(), []int{0}
}

func (x *Table) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Table) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Table) GetOp() Table_OP {
	if x != nil {
		return x.Op
	}
	return Table_CREATE
}

func (x *Table) GetTs() uint64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *Table) GetStatement() string {
	if x != nil {
		return x.Statement
	}
	return ""
}

func (x *Table) GetCreateStatement() string {
	if x != nil {
		return x.CreateStatement
	}
	return ""
}

type Schema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id        uint64    `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Charset   string    `protobuf:"bytes,3,opt,name=charset,proto3" json:"charset,omitempty"`
	Collation string    `protobuf:"bytes,4,opt,name=collation,proto3" json:"collation,omitempty"`
	Op        Schema_OP `protobuf:"varint,5,opt,name=op,proto3,enum=tistreampb.Schema_OP" json:"op,omitempty"`
	Ts        uint64    `protobuf:"varint,6,opt,name=ts,proto3" json:"ts,omitempty"`
}

func (x *Schema) Reset() {
	*x = Schema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schemapb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Schema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Schema) ProtoMessage() {}

func (x *Schema) ProtoReflect() protoreflect.Message {
	mi := &file_schemapb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Schema.ProtoReflect.Descriptor instead.
func (*Schema) Descriptor() ([]byte, []int) {
	return file_schemapb_proto_rawDescGZIP(), []int{1}
}

func (x *Schema) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Schema) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Schema) GetCharset() string {
	if x != nil {
		return x.Charset
	}
	return ""
}

func (x *Schema) GetCollation() string {
	if x != nil {
		return x.Collation
	}
	return ""
}

func (x *Schema) GetOp() Schema_OP {
	if x != nil {
		return x.Op
	}
	return Schema_CREATE
}

func (x *Schema) GetTs() uint64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

type SchemaTables struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id        uint64   `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Charset   string   `protobuf:"bytes,3,opt,name=charset,proto3" json:"charset,omitempty"`
	Collation string   `protobuf:"bytes,4,opt,name=collation,proto3" json:"collation,omitempty"`
	Tables    []*Table `protobuf:"bytes,5,rep,name=tables,proto3" json:"tables,omitempty"`
}

func (x *SchemaTables) Reset() {
	*x = SchemaTables{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schemapb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemaTables) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemaTables) ProtoMessage() {}

func (x *SchemaTables) ProtoReflect() protoreflect.Message {
	mi := &file_schemapb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemaTables.ProtoReflect.Descriptor instead.
func (*SchemaTables) Descriptor() ([]byte, []int) {
	return file_schemapb_proto_rawDescGZIP(), []int{2}
}

func (x *SchemaTables) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SchemaTables) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SchemaTables) GetCharset() string {
	if x != nil {
		return x.Charset
	}
	return ""
}

func (x *SchemaTables) GetCollation() string {
	if x != nil {
		return x.Collation
	}
	return ""
}

func (x *SchemaTables) GetTables() []*Table {
	if x != nil {
		return x.Tables
	}
	return nil
}

type SchemasSnapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts      uint64          `protobuf:"varint,1,opt,name=ts,proto3" json:"ts,omitempty"`
	Schemas []*SchemaTables `protobuf:"bytes,2,rep,name=schemas,proto3" json:"schemas,omitempty"`
}

func (x *SchemasSnapshot) Reset() {
	*x = SchemasSnapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schemapb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SchemasSnapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SchemasSnapshot) ProtoMessage() {}

func (x *SchemasSnapshot) ProtoReflect() protoreflect.Message {
	mi := &file_schemapb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SchemasSnapshot.ProtoReflect.Descriptor instead.
func (*SchemasSnapshot) Descriptor() ([]byte, []int) {
	return file_schemapb_proto_rawDescGZIP(), []int{3}
}

func (x *SchemasSnapshot) GetTs() uint64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

func (x *SchemasSnapshot) GetSchemas() []*SchemaTables {
	if x != nil {
		return x.Schemas
	}
	return nil
}

var File_schemapb_proto protoreflect.FileDescriptor

var file_schemapb_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x22, 0xd1, 0x01, 0x0a,
	0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x02, 0x6f, 0x70,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x70, 0x62, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x4f, 0x50, 0x52, 0x02, 0x6f, 0x70,
	0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x74, 0x73,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x29,
	0x0a, 0x10, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x25, 0x0a, 0x02, 0x4f, 0x50, 0x12,
	0x0a, 0x0a, 0x06, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x44,
	0x52, 0x4f, 0x50, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x4c, 0x54, 0x45, 0x52, 0x10, 0x02,
	0x22, 0xb7, 0x01, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6c,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f,
	0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x02, 0x6f, 0x70, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62,
	0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x4f, 0x50, 0x52, 0x02, 0x6f, 0x70, 0x12, 0x0e,
	0x0a, 0x02, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x74, 0x73, 0x22, 0x1a,
	0x0a, 0x02, 0x4f, 0x50, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x00,
	0x12, 0x08, 0x0a, 0x04, 0x44, 0x52, 0x4f, 0x50, 0x10, 0x01, 0x22, 0x95, 0x01, 0x0a, 0x0c, 0x53,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6c,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f,
	0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x70, 0x62, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x06, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x73, 0x22, 0x55, 0x0a, 0x0f, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x53, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x74, 0x73, 0x12, 0x32, 0x0a, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x70, 0x62, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73,
	0x52, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_schemapb_proto_rawDescOnce sync.Once
	file_schemapb_proto_rawDescData = file_schemapb_proto_rawDesc
)

func file_schemapb_proto_rawDescGZIP() []byte {
	file_schemapb_proto_rawDescOnce.Do(func() {
		file_schemapb_proto_rawDescData = protoimpl.X.CompressGZIP(file_schemapb_proto_rawDescData)
	})
	return file_schemapb_proto_rawDescData
}

var file_schemapb_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_schemapb_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_schemapb_proto_goTypes = []interface{}{
	(Table_OP)(0),           // 0: tistreampb.Table.OP
	(Schema_OP)(0),          // 1: tistreampb.Schema.OP
	(*Table)(nil),           // 2: tistreampb.Table
	(*Schema)(nil),          // 3: tistreampb.Schema
	(*SchemaTables)(nil),    // 4: tistreampb.SchemaTables
	(*SchemasSnapshot)(nil), // 5: tistreampb.SchemasSnapshot
}
var file_schemapb_proto_depIdxs = []int32{
	0, // 0: tistreampb.Table.op:type_name -> tistreampb.Table.OP
	1, // 1: tistreampb.Schema.op:type_name -> tistreampb.Schema.OP
	2, // 2: tistreampb.SchemaTables.tables:type_name -> tistreampb.Table
	4, // 3: tistreampb.SchemasSnapshot.schemas:type_name -> tistreampb.SchemaTables
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_schemapb_proto_init() }
func file_schemapb_proto_init() {
	if File_schemapb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_schemapb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Table); i {
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
		file_schemapb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Schema); i {
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
		file_schemapb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchemaTables); i {
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
		file_schemapb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SchemasSnapshot); i {
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
			RawDescriptor: file_schemapb_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_schemapb_proto_goTypes,
		DependencyIndexes: file_schemapb_proto_depIdxs,
		EnumInfos:         file_schemapb_proto_enumTypes,
		MessageInfos:      file_schemapb_proto_msgTypes,
	}.Build()
	File_schemapb_proto = out.File
	file_schemapb_proto_rawDesc = nil
	file_schemapb_proto_goTypes = nil
	file_schemapb_proto_depIdxs = nil
}
