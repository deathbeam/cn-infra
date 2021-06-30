// Code generated by protoc-gen-go. DO NOT EDIT.
// source: model.proto

package testdata

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TestData struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TestData) Reset()         { *m = TestData{} }
func (m *TestData) String() string { return proto.CompactTextString(m) }
func (*TestData) ProtoMessage()    {}
func (*TestData) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{0}
}

func (m *TestData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestData.Unmarshal(m, b)
}
func (m *TestData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestData.Marshal(b, m, deterministic)
}
func (m *TestData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestData.Merge(m, src)
}
func (m *TestData) XXX_Size() int {
	return xxx_messageInfo_TestData.Size(m)
}
func (m *TestData) XXX_DiscardUnknown() {
	xxx_messageInfo_TestData.DiscardUnknown(m)
}

var xxx_messageInfo_TestData proto.InternalMessageInfo

func (m *TestData) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *TestData) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type TestNested struct {
	Name                 string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data                 *TestData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *TestNested) Reset()         { *m = TestNested{} }
func (m *TestNested) String() string { return proto.CompactTextString(m) }
func (*TestNested) ProtoMessage()    {}
func (*TestNested) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{1}
}

func (m *TestNested) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestNested.Unmarshal(m, b)
}
func (m *TestNested) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestNested.Marshal(b, m, deterministic)
}
func (m *TestNested) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestNested.Merge(m, src)
}
func (m *TestNested) XXX_Size() int {
	return xxx_messageInfo_TestNested.Size(m)
}
func (m *TestNested) XXX_DiscardUnknown() {
	xxx_messageInfo_TestNested.DiscardUnknown(m)
}

var xxx_messageInfo_TestNested proto.InternalMessageInfo

func (m *TestNested) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TestNested) GetData() *TestData {
	if m != nil {
		return m.Data
	}
	return nil
}

type TestSlice struct {
	Name                 string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data                 []*TestData `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TestSlice) Reset()         { *m = TestSlice{} }
func (m *TestSlice) String() string { return proto.CompactTextString(m) }
func (*TestSlice) ProtoMessage()    {}
func (*TestSlice) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c16552f9fdb66d8, []int{2}
}

func (m *TestSlice) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestSlice.Unmarshal(m, b)
}
func (m *TestSlice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestSlice.Marshal(b, m, deterministic)
}
func (m *TestSlice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestSlice.Merge(m, src)
}
func (m *TestSlice) XXX_Size() int {
	return xxx_messageInfo_TestSlice.Size(m)
}
func (m *TestSlice) XXX_DiscardUnknown() {
	xxx_messageInfo_TestSlice.DiscardUnknown(m)
}

var xxx_messageInfo_TestSlice proto.InternalMessageInfo

func (m *TestSlice) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TestSlice) GetData() []*TestData {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*TestData)(nil), "testdata.TestData")
	proto.RegisterType((*TestNested)(nil), "testdata.TestNested")
	proto.RegisterType((*TestSlice)(nil), "testdata.TestSlice")
}

func init() { proto.RegisterFile("model.proto", fileDescriptor_4c16552f9fdb66d8) }

var fileDescriptor_4c16552f9fdb66d8 = []byte{
	// 158 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xcd, 0x4f, 0x49,
	0xcd, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x28, 0x49, 0x2d, 0x2e, 0x49, 0x49, 0x2c,
	0x49, 0x54, 0x72, 0xe2, 0xe2, 0x08, 0x49, 0x2d, 0x2e, 0x71, 0x49, 0x2c, 0x49, 0x14, 0x92, 0xe2,
	0xe2, 0x28, 0x2d, 0x4e, 0x2d, 0xca, 0x4b, 0xcc, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c,
	0x82, 0xf3, 0x41, 0x72, 0x05, 0x89, 0xc5, 0xc5, 0xe5, 0xf9, 0x45, 0x29, 0x12, 0x4c, 0x10, 0x39,
	0x18, 0x5f, 0xc9, 0x83, 0x8b, 0x0b, 0x64, 0x86, 0x5f, 0x6a, 0x71, 0x49, 0x6a, 0x8a, 0x90, 0x10,
	0x17, 0x0b, 0x92, 0x09, 0x60, 0xb6, 0x90, 0x1a, 0x17, 0x0b, 0xc8, 0x36, 0xb0, 0x4e, 0x6e, 0x23,
	0x21, 0x3d, 0x98, 0xf5, 0x7a, 0x30, 0xbb, 0x83, 0xc0, 0xf2, 0x4a, 0xee, 0x5c, 0x9c, 0x20, 0x91,
	0xe0, 0x9c, 0xcc, 0xe4, 0x54, 0x02, 0x06, 0x31, 0xe3, 0x33, 0x28, 0x89, 0x0d, 0xec, 0x4f, 0x63,
	0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xed, 0x5d, 0xc2, 0xcc, 0xf6, 0x00, 0x00, 0x00,
}