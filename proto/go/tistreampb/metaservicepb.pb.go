// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.8.0
// source: metaservicepb.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_metaservicepb_proto protoreflect.FileDescriptor

var file_metaservicepb_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6d, 0x65, 0x74, 0x61, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x70, 0x62, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70,
	0x62, 0x1a, 0x0c, 0x6d, 0x65, 0x74, 0x61, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32,
	0xe1, 0x03, 0x0a, 0x0b, 0x4d, 0x65, 0x74, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x4f, 0x0a, 0x12, 0x54, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x4e, 0x65, 0x77, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1b, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x70, 0x62, 0x2e, 0x48, 0x61, 0x73, 0x4e, 0x65, 0x77, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52,
	0x65, 0x71, 0x1a, 0x1c, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e,
	0x48, 0x61, 0x73, 0x4e, 0x65, 0x77, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x5e, 0x0a, 0x13, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x48, 0x65,
	0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x22, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x70, 0x62, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x48,
	0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x23, 0x2e, 0x74, 0x69,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x52, 0x0a, 0x0f, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62,
	0x65, 0x61, 0x74, 0x12, 0x1e, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62,
	0x2e, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74,
	0x52, 0x65, 0x71, 0x1a, 0x1f, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62,
	0x2e, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x6a, 0x0a, 0x17, 0x46, 0x65, 0x74, 0x63, 0x68, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x41, 0x64, 0x64, 0x72, 0x12,
	0x26, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x41, 0x64, 0x64, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x27, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x70, 0x62, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x41, 0x64, 0x64, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x61, 0x0a, 0x14, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x12, 0x23, 0x2e, 0x74, 0x69, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x53, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x24, 0x2e,
	0x74, 0x69, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x70, 0x62, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x52, 0x61, 0x6e, 0x67, 0x65, 0x53, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var file_metaservicepb_proto_goTypes = []interface{}{
	(*HasNewChangeReq)(nil),             // 0: tistreampb.HasNewChangeReq
	(*DispatcherHeartbeatReq)(nil),      // 1: tistreampb.DispatcherHeartbeatReq
	(*SorterHeartbeatReq)(nil),          // 2: tistreampb.SorterHeartbeatReq
	(*FetchSchemaRegistryAddrReq)(nil),  // 3: tistreampb.FetchSchemaRegistryAddrReq
	(*FetchRangeSorterAddrReq)(nil),     // 4: tistreampb.FetchRangeSorterAddrReq
	(*HasNewChangeResp)(nil),            // 5: tistreampb.HasNewChangeResp
	(*DispatcherHeartbeatResp)(nil),     // 6: tistreampb.DispatcherHeartbeatResp
	(*SorterHeartbeatResp)(nil),         // 7: tistreampb.SorterHeartbeatResp
	(*FetchSchemaRegistryAddrResp)(nil), // 8: tistreampb.FetchSchemaRegistryAddrResp
	(*FetchRangeSorterAddrResp)(nil),    // 9: tistreampb.FetchRangeSorterAddrResp
}
var file_metaservicepb_proto_depIdxs = []int32{
	0, // 0: tistreampb.MetaService.TenantHasNewChange:input_type -> tistreampb.HasNewChangeReq
	1, // 1: tistreampb.MetaService.DispatcherHeartbeat:input_type -> tistreampb.DispatcherHeartbeatReq
	2, // 2: tistreampb.MetaService.SorterHeartbeat:input_type -> tistreampb.SorterHeartbeatReq
	3, // 3: tistreampb.MetaService.FetchSchemaRegistryAddr:input_type -> tistreampb.FetchSchemaRegistryAddrReq
	4, // 4: tistreampb.MetaService.FetchRangeSorterAddr:input_type -> tistreampb.FetchRangeSorterAddrReq
	5, // 5: tistreampb.MetaService.TenantHasNewChange:output_type -> tistreampb.HasNewChangeResp
	6, // 6: tistreampb.MetaService.DispatcherHeartbeat:output_type -> tistreampb.DispatcherHeartbeatResp
	7, // 7: tistreampb.MetaService.SorterHeartbeat:output_type -> tistreampb.SorterHeartbeatResp
	8, // 8: tistreampb.MetaService.FetchSchemaRegistryAddr:output_type -> tistreampb.FetchSchemaRegistryAddrResp
	9, // 9: tistreampb.MetaService.FetchRangeSorterAddr:output_type -> tistreampb.FetchRangeSorterAddrResp
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_metaservicepb_proto_init() }
func file_metaservicepb_proto_init() {
	if File_metaservicepb_proto != nil {
		return
	}
	file_metapb_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metaservicepb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metaservicepb_proto_goTypes,
		DependencyIndexes: file_metaservicepb_proto_depIdxs,
	}.Build()
	File_metaservicepb_proto = out.File
	file_metaservicepb_proto_rawDesc = nil
	file_metaservicepb_proto_goTypes = nil
	file_metaservicepb_proto_depIdxs = nil
}
