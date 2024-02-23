// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: api/v1/covalyzer.proto

package v1

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

type Commit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sha string `protobuf:"bytes,1,opt,name=sha,proto3" json:"sha,omitempty"`
}

func (x *Commit) Reset() {
	*x = Commit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_covalyzer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_covalyzer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commit.ProtoReflect.Descriptor instead.
func (*Commit) Descriptor() ([]byte, []int) {
	return file_api_v1_covalyzer_proto_rawDescGZIP(), []int{0}
}

func (x *Commit) GetSha() string {
	if x != nil {
		return x.Sha
	}
	return ""
}

type Cover struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total float32 `protobuf:"fixed32,1,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *Cover) Reset() {
	*x = Cover{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_covalyzer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cover) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cover) ProtoMessage() {}

func (x *Cover) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_covalyzer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cover.ProtoReflect.Descriptor instead.
func (*Cover) Descriptor() ([]byte, []int) {
	return file_api_v1_covalyzer_proto_rawDescGZIP(), []int{1}
}

func (x *Cover) GetTotal() float32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type Coverage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Commit *Commit `protobuf:"bytes,1,opt,name=commit,proto3" json:"commit,omitempty"`
	Cover  *Cover  `protobuf:"bytes,2,opt,name=cover,proto3" json:"cover,omitempty"`
}

func (x *Coverage) Reset() {
	*x = Coverage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_covalyzer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Coverage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Coverage) ProtoMessage() {}

func (x *Coverage) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_covalyzer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Coverage.ProtoReflect.Descriptor instead.
func (*Coverage) Descriptor() ([]byte, []int) {
	return file_api_v1_covalyzer_proto_rawDescGZIP(), []int{2}
}

func (x *Coverage) GetCommit() *Commit {
	if x != nil {
		return x.Commit
	}
	return nil
}

func (x *Coverage) GetCover() *Cover {
	if x != nil {
		return x.Cover
	}
	return nil
}

type RepositoryCoverages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repository *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	Coverages  []*Coverage `protobuf:"bytes,2,rep,name=coverages,proto3" json:"coverages,omitempty"`
}

func (x *RepositoryCoverages) Reset() {
	*x = RepositoryCoverages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_covalyzer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RepositoryCoverages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RepositoryCoverages) ProtoMessage() {}

func (x *RepositoryCoverages) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_covalyzer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RepositoryCoverages.ProtoReflect.Descriptor instead.
func (*RepositoryCoverages) Descriptor() ([]byte, []int) {
	return file_api_v1_covalyzer_proto_rawDescGZIP(), []int{3}
}

func (x *RepositoryCoverages) GetRepository() *Repository {
	if x != nil {
		return x.Repository
	}
	return nil
}

func (x *RepositoryCoverages) GetCoverages() []*Coverage {
	if x != nil {
		return x.Coverages
	}
	return nil
}

type Covalyzer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repositories []*RepositoryCoverages `protobuf:"bytes,1,rep,name=repositories,proto3" json:"repositories,omitempty"`
}

func (x *Covalyzer) Reset() {
	*x = Covalyzer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_covalyzer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Covalyzer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Covalyzer) ProtoMessage() {}

func (x *Covalyzer) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_covalyzer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Covalyzer.ProtoReflect.Descriptor instead.
func (*Covalyzer) Descriptor() ([]byte, []int) {
	return file_api_v1_covalyzer_proto_rawDescGZIP(), []int{4}
}

func (x *Covalyzer) GetRepositories() []*RepositoryCoverages {
	if x != nil {
		return x.Repositories
	}
	return nil
}

var File_api_v1_covalyzer_proto protoreflect.FileDescriptor

var file_api_v1_covalyzer_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x76, 0x61, 0x6c, 0x79, 0x7a,
	0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x1a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1a, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x68, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x68,
	0x61, 0x22, 0x1d, 0x0a, 0x05, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x22, 0x57, 0x0a, 0x08, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x06, 0x63, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x12, 0x23, 0x0a, 0x05, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x76,
	0x65, 0x72, 0x52, 0x05, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x22, 0x79, 0x0a, 0x13, 0x52, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x73,
	0x12, 0x32, 0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x6f, 0x72, 0x79, 0x12, 0x2e, 0x0a, 0x09, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x52, 0x09, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x61, 0x67, 0x65, 0x73, 0x22, 0x4c, 0x0a, 0x09, 0x43, 0x6f, 0x76, 0x61, 0x6c, 0x79, 0x7a, 0x65,
	0x72, 0x12, 0x3f, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x43, 0x6f, 0x76, 0x65, 0x72,
	0x61, 0x67, 0x65, 0x73, 0x52, 0x0c, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x69,
	0x65, 0x73, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6e, 0x6f, 0x6b, 0x61, 0x6d, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x76, 0x61, 0x6c, 0x79,
	0x7a, 0x65, 0x72, 0x2d, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_covalyzer_proto_rawDescOnce sync.Once
	file_api_v1_covalyzer_proto_rawDescData = file_api_v1_covalyzer_proto_rawDesc
)

func file_api_v1_covalyzer_proto_rawDescGZIP() []byte {
	file_api_v1_covalyzer_proto_rawDescOnce.Do(func() {
		file_api_v1_covalyzer_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_covalyzer_proto_rawDescData)
	})
	return file_api_v1_covalyzer_proto_rawDescData
}

var file_api_v1_covalyzer_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_v1_covalyzer_proto_goTypes = []interface{}{
	(*Commit)(nil),              // 0: api.v1.Commit
	(*Cover)(nil),               // 1: api.v1.Cover
	(*Coverage)(nil),            // 2: api.v1.Coverage
	(*RepositoryCoverages)(nil), // 3: api.v1.RepositoryCoverages
	(*Covalyzer)(nil),           // 4: api.v1.Covalyzer
	(*Repository)(nil),          // 5: api.v1.Repository
}
var file_api_v1_covalyzer_proto_depIdxs = []int32{
	0, // 0: api.v1.Coverage.commit:type_name -> api.v1.Commit
	1, // 1: api.v1.Coverage.cover:type_name -> api.v1.Cover
	5, // 2: api.v1.RepositoryCoverages.repository:type_name -> api.v1.Repository
	2, // 3: api.v1.RepositoryCoverages.coverages:type_name -> api.v1.Coverage
	3, // 4: api.v1.Covalyzer.repositories:type_name -> api.v1.RepositoryCoverages
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_api_v1_covalyzer_proto_init() }
func file_api_v1_covalyzer_proto_init() {
	if File_api_v1_covalyzer_proto != nil {
		return
	}
	file_api_v1_config_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_v1_covalyzer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Commit); i {
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
		file_api_v1_covalyzer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cover); i {
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
		file_api_v1_covalyzer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Coverage); i {
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
		file_api_v1_covalyzer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RepositoryCoverages); i {
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
		file_api_v1_covalyzer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Covalyzer); i {
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
			RawDescriptor: file_api_v1_covalyzer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_covalyzer_proto_goTypes,
		DependencyIndexes: file_api_v1_covalyzer_proto_depIdxs,
		MessageInfos:      file_api_v1_covalyzer_proto_msgTypes,
	}.Build()
	File_api_v1_covalyzer_proto = out.File
	file_api_v1_covalyzer_proto_rawDesc = nil
	file_api_v1_covalyzer_proto_goTypes = nil
	file_api_v1_covalyzer_proto_depIdxs = nil
}