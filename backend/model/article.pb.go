// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: backend/model/article.proto

package model

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

type Articles struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ListOfArticles []*Article `protobuf:"bytes,1,rep,name=listOfArticles,proto3" json:"listOfArticles,omitempty"`
}

func (x *Articles) Reset() {
	*x = Articles{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_model_article_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Articles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Articles) ProtoMessage() {}

func (x *Articles) ProtoReflect() protoreflect.Message {
	mi := &file_backend_model_article_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Articles.ProtoReflect.Descriptor instead.
func (*Articles) Descriptor() ([]byte, []int) {
	return file_backend_model_article_proto_rawDescGZIP(), []int{0}
}

func (x *Articles) GetListOfArticles() []*Article {
	if x != nil {
		return x.ListOfArticles
	}
	return nil
}

type Article struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid              string  `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Title             string  `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Url               string  `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Tags              string  `protobuf:"bytes,4,opt,name=tags,proto3" json:"tags,omitempty"`
	FriendlyUrl       string  `protobuf:"bytes,5,opt,name=friendly_url,json=friendlyUrl,proto3" json:"friendly_url,omitempty"`
	CreationTimestamp int64   `protobuf:"varint,6,opt,name=creation_timestamp,json=creationTimestamp,proto3" json:"creation_timestamp,omitempty"`
	EditTimestamp     *int64  `protobuf:"varint,7,opt,name=edit_timestamp,json=editTimestamp,proto3,oneof" json:"edit_timestamp,omitempty"`
	MetaDescription   *string `protobuf:"bytes,8,opt,name=meta_description,json=metaDescription,proto3,oneof" json:"meta_description,omitempty"`
	Published         bool    `protobuf:"varint,9,opt,name=published,proto3" json:"published,omitempty"`
	IntegrityHash     string  `protobuf:"bytes,10,opt,name=integrity_hash,json=integrityHash,proto3" json:"integrity_hash,omitempty"`
}

func (x *Article) Reset() {
	*x = Article{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_model_article_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Article) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Article) ProtoMessage() {}

func (x *Article) ProtoReflect() protoreflect.Message {
	mi := &file_backend_model_article_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Article.ProtoReflect.Descriptor instead.
func (*Article) Descriptor() ([]byte, []int) {
	return file_backend_model_article_proto_rawDescGZIP(), []int{1}
}

func (x *Article) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Article) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Article) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Article) GetTags() string {
	if x != nil {
		return x.Tags
	}
	return ""
}

func (x *Article) GetFriendlyUrl() string {
	if x != nil {
		return x.FriendlyUrl
	}
	return ""
}

func (x *Article) GetCreationTimestamp() int64 {
	if x != nil {
		return x.CreationTimestamp
	}
	return 0
}

func (x *Article) GetEditTimestamp() int64 {
	if x != nil && x.EditTimestamp != nil {
		return *x.EditTimestamp
	}
	return 0
}

func (x *Article) GetMetaDescription() string {
	if x != nil && x.MetaDescription != nil {
		return *x.MetaDescription
	}
	return ""
}

func (x *Article) GetPublished() bool {
	if x != nil {
		return x.Published
	}
	return false
}

func (x *Article) GetIntegrityHash() string {
	if x != nil {
		return x.IntegrityHash
	}
	return ""
}

var File_backend_model_article_proto protoreflect.FileDescriptor

var file_backend_model_article_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f,
	0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a,
	0x08, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x12, 0x30, 0x0a, 0x0e, 0x6c, 0x69, 0x73,
	0x74, 0x4f, 0x66, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x08, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x0e, 0x6c, 0x69, 0x73,
	0x74, 0x4f, 0x66, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x22, 0xf4, 0x02, 0x0a, 0x07,
	0x41, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x72, 0x69, 0x65, 0x6e,
	0x64, 0x6c, 0x79, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x66,
	0x72, 0x69, 0x65, 0x6e, 0x64, 0x6c, 0x79, 0x55, 0x72, 0x6c, 0x12, 0x2d, 0x0a, 0x12, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x2a, 0x0a, 0x0e, 0x65, 0x64, 0x69,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x03, 0x48, 0x00, 0x52, 0x0d, 0x65, 0x64, 0x69, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x88, 0x01, 0x01, 0x12, 0x2e, 0x0a, 0x10, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x0f, 0x6d, 0x65, 0x74, 0x61, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73,
	0x68, 0x65, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x69, 0x74, 0x79,
	0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x69, 0x6e, 0x74,
	0x65, 0x67, 0x72, 0x69, 0x74, 0x79, 0x48, 0x61, 0x73, 0x68, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x65,
	0x64, 0x69, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x13, 0x0a,
	0x11, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x42, 0x11, 0x5a, 0x0f, 0x2e, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_backend_model_article_proto_rawDescOnce sync.Once
	file_backend_model_article_proto_rawDescData = file_backend_model_article_proto_rawDesc
)

func file_backend_model_article_proto_rawDescGZIP() []byte {
	file_backend_model_article_proto_rawDescOnce.Do(func() {
		file_backend_model_article_proto_rawDescData = protoimpl.X.CompressGZIP(file_backend_model_article_proto_rawDescData)
	})
	return file_backend_model_article_proto_rawDescData
}

var file_backend_model_article_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_backend_model_article_proto_goTypes = []interface{}{
	(*Articles)(nil), // 0: Articles
	(*Article)(nil),  // 1: Article
}
var file_backend_model_article_proto_depIdxs = []int32{
	1, // 0: Articles.listOfArticles:type_name -> Article
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_backend_model_article_proto_init() }
func file_backend_model_article_proto_init() {
	if File_backend_model_article_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_backend_model_article_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Articles); i {
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
		file_backend_model_article_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Article); i {
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
	file_backend_model_article_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_backend_model_article_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_backend_model_article_proto_goTypes,
		DependencyIndexes: file_backend_model_article_proto_depIdxs,
		MessageInfos:      file_backend_model_article_proto_msgTypes,
	}.Build()
	File_backend_model_article_proto = out.File
	file_backend_model_article_proto_rawDesc = nil
	file_backend_model_article_proto_goTypes = nil
	file_backend_model_article_proto_depIdxs = nil
}
