// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cloudstore.proto

/*
Package cloudstore is a generated protocol buffer package.

It is generated from these files:
	cloudstore.proto

It has these top-level messages:
	Chunk
	UploadResponse
	GetRequest
	DeleteRequest
	DeleteResponse
*/
package cloudstore

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type UploadStatusCode int32

const (
	UploadStatusCode_unknown UploadStatusCode = 0
	UploadStatusCode_ok      UploadStatusCode = 1
	UploadStatusCode_failed  UploadStatusCode = 2
)

var UploadStatusCode_name = map[int32]string{
	0: "unknown",
	1: "ok",
	2: "failed",
}
var UploadStatusCode_value = map[string]int32{
	"unknown": 0,
	"ok":      1,
	"failed":  2,
}

func (x UploadStatusCode) String() string {
	return proto.EnumName(UploadStatusCode_name, int32(x))
}
func (UploadStatusCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Chunk struct {
	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (m *Chunk) Reset()                    { *m = Chunk{} }
func (m *Chunk) String() string            { return proto.CompactTextString(m) }
func (*Chunk) ProtoMessage()               {}
func (*Chunk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Chunk) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type UploadResponse struct {
	Message string           `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
	Url     string           `protobuf:"bytes,2,opt,name=url" json:"url,omitempty"`
	Code    UploadStatusCode `protobuf:"varint,3,opt,name=code,enum=cloudstore.UploadStatusCode" json:"code,omitempty"`
}

func (m *UploadResponse) Reset()                    { *m = UploadResponse{} }
func (m *UploadResponse) String() string            { return proto.CompactTextString(m) }
func (*UploadResponse) ProtoMessage()               {}
func (*UploadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UploadResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *UploadResponse) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *UploadResponse) GetCode() UploadStatusCode {
	if m != nil {
		return m.Code
	}
	return UploadStatusCode_unknown
}

type GetRequest struct {
	Filename string `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetRequest) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

type DeleteRequest struct {
	Filename string `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DeleteRequest) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

type DeleteResponse struct {
	Filename string `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
}

func (m *DeleteResponse) Reset()                    { *m = DeleteResponse{} }
func (m *DeleteResponse) String() string            { return proto.CompactTextString(m) }
func (*DeleteResponse) ProtoMessage()               {}
func (*DeleteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DeleteResponse) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func init() {
	proto.RegisterType((*Chunk)(nil), "cloudstore.Chunk")
	proto.RegisterType((*UploadResponse)(nil), "cloudstore.UploadResponse")
	proto.RegisterType((*GetRequest)(nil), "cloudstore.GetRequest")
	proto.RegisterType((*DeleteRequest)(nil), "cloudstore.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "cloudstore.DeleteResponse")
	proto.RegisterEnum("cloudstore.UploadStatusCode", UploadStatusCode_name, UploadStatusCode_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CloudStorageService service

type CloudStorageServiceClient interface {
	Store(ctx context.Context, opts ...grpc.CallOption) (CloudStorageService_StoreClient, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
}

type cloudStorageServiceClient struct {
	cc *grpc.ClientConn
}

func NewCloudStorageServiceClient(cc *grpc.ClientConn) CloudStorageServiceClient {
	return &cloudStorageServiceClient{cc}
}

func (c *cloudStorageServiceClient) Store(ctx context.Context, opts ...grpc.CallOption) (CloudStorageService_StoreClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_CloudStorageService_serviceDesc.Streams[0], c.cc, "/cloudstore.CloudStorageService/Store", opts...)
	if err != nil {
		return nil, err
	}
	x := &cloudStorageServiceStoreClient{stream}
	return x, nil
}

type CloudStorageService_StoreClient interface {
	Send(*Chunk) error
	CloseAndRecv() (*UploadResponse, error)
	grpc.ClientStream
}

type cloudStorageServiceStoreClient struct {
	grpc.ClientStream
}

func (x *cloudStorageServiceStoreClient) Send(m *Chunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cloudStorageServiceStoreClient) CloseAndRecv() (*UploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cloudStorageServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := grpc.Invoke(ctx, "/cloudstore.CloudStorageService/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CloudStorageService service

type CloudStorageServiceServer interface {
	Store(CloudStorageService_StoreServer) error
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
}

func RegisterCloudStorageServiceServer(s *grpc.Server, srv CloudStorageServiceServer) {
	s.RegisterService(&_CloudStorageService_serviceDesc, srv)
}

func _CloudStorageService_Store_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CloudStorageServiceServer).Store(&cloudStorageServiceStoreServer{stream})
}

type CloudStorageService_StoreServer interface {
	SendAndClose(*UploadResponse) error
	Recv() (*Chunk, error)
	grpc.ServerStream
}

type cloudStorageServiceStoreServer struct {
	grpc.ServerStream
}

func (x *cloudStorageServiceStoreServer) SendAndClose(m *UploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cloudStorageServiceStoreServer) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CloudStorageService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudStorageServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cloudstore.CloudStorageService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudStorageServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CloudStorageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cloudstore.CloudStorageService",
	HandlerType: (*CloudStorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Delete",
			Handler:    _CloudStorageService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Store",
			Handler:       _CloudStorageService_Store_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "cloudstore.proto",
}

func init() { proto.RegisterFile("cloudstore.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 316 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x4f, 0xc2, 0x40,
	0x10, 0xc5, 0x29, 0x48, 0xd1, 0x51, 0x49, 0x1d, 0x2f, 0x95, 0x18, 0x83, 0x3d, 0x35, 0x6a, 0xc0,
	0xc0, 0xcd, 0x9b, 0x62, 0xe2, 0xbd, 0xc4, 0x8b, 0xb7, 0xa5, 0x1d, 0xa0, 0x69, 0xd9, 0xc1, 0xfd,
	0xa3, 0x5f, 0xc4, 0x0f, 0x6c, 0x5a, 0xac, 0x50, 0x4d, 0x8c, 0xb7, 0x7d, 0x2f, 0xbf, 0x9d, 0x7d,
	0x6f, 0xb2, 0xe0, 0xc5, 0x39, 0xdb, 0x44, 0x1b, 0x56, 0x34, 0x58, 0x2b, 0x36, 0x8c, 0xb0, 0x75,
	0x82, 0x4b, 0x68, 0x4f, 0x96, 0x56, 0x66, 0xe8, 0x43, 0x27, 0x66, 0x69, 0x48, 0x1a, 0xdf, 0xe9,
	0x3b, 0xe1, 0x51, 0x54, 0xc9, 0x40, 0x42, 0xf7, 0x79, 0x9d, 0xb3, 0x48, 0x22, 0xd2, 0x6b, 0x96,
	0x9a, 0x0a, 0x76, 0x45, 0x5a, 0x8b, 0x05, 0x95, 0xec, 0x41, 0x54, 0x49, 0xf4, 0xa0, 0x65, 0x55,
	0xee, 0x37, 0x4b, 0xb7, 0x38, 0xe2, 0x2d, 0xec, 0xc5, 0x9c, 0x90, 0xdf, 0xea, 0x3b, 0x61, 0x77,
	0x74, 0x3e, 0xd8, 0x49, 0xb3, 0x99, 0x3a, 0x35, 0xc2, 0x58, 0x3d, 0xe1, 0x84, 0xa2, 0x92, 0x0c,
	0x42, 0x80, 0x27, 0x32, 0x11, 0xbd, 0x5a, 0xd2, 0x06, 0x7b, 0xb0, 0x3f, 0x4f, 0x73, 0x92, 0x62,
	0x55, 0x3d, 0xf6, 0xad, 0x83, 0x6b, 0x38, 0x7e, 0xa4, 0x9c, 0x0c, 0xfd, 0x07, 0xbe, 0x81, 0x6e,
	0x05, 0x7f, 0xd5, 0xf8, 0x83, 0xbe, 0x1a, 0x83, 0xf7, 0x33, 0x1e, 0x1e, 0x42, 0xc7, 0xca, 0x4c,
	0xf2, 0xbb, 0xf4, 0x1a, 0xe8, 0x42, 0x93, 0x33, 0xcf, 0x41, 0x00, 0x77, 0x2e, 0xd2, 0x9c, 0x12,
	0xaf, 0x39, 0xfa, 0x70, 0xe0, 0x74, 0x52, 0xf4, 0x9b, 0x1a, 0x56, 0x62, 0x41, 0x53, 0x52, 0x6f,
	0x69, 0x4c, 0x78, 0x07, 0xed, 0xc2, 0x21, 0x3c, 0xd9, 0xad, 0x5f, 0xee, 0xbd, 0xd7, 0xfb, 0xbd,
	0x91, 0x2a, 0x60, 0xd0, 0x08, 0x1d, 0xbc, 0x07, 0x77, 0x13, 0x1b, 0xcf, 0x76, 0xc9, 0x5a, 0xef,
	0xfa, 0x90, 0x7a, 0xcb, 0xa0, 0xf1, 0xd0, 0x7f, 0xb9, 0x58, 0xa4, 0x66, 0x69, 0x67, 0x83, 0x98,
	0x57, 0x43, 0x65, 0x35, 0x49, 0xa1, 0xb3, 0xe1, 0xf6, 0xca, 0xcc, 0x2d, 0x3f, 0xc6, 0xf8, 0x33,
	0x00, 0x00, 0xff, 0xff, 0x27, 0xa5, 0x2b, 0xaf, 0x2c, 0x02, 0x00, 0x00,
}
