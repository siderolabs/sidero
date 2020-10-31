// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SystemInformation struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Manufacturer         string   `protobuf:"bytes,2,opt,name=manufacturer,proto3" json:"manufacturer,omitempty"`
	ProductName          string   `protobuf:"bytes,3,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	Version              string   `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	SerialNumber         string   `protobuf:"bytes,5,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	SkuNumber            string   `protobuf:"bytes,6,opt,name=sku_number,json=skuNumber,proto3" json:"sku_number,omitempty"`
	Family               string   `protobuf:"bytes,7,opt,name=family,proto3" json:"family,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemInformation) Reset()         { *m = SystemInformation{} }
func (m *SystemInformation) String() string { return proto.CompactTextString(m) }
func (*SystemInformation) ProtoMessage()    {}
func (*SystemInformation) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *SystemInformation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemInformation.Unmarshal(m, b)
}

func (m *SystemInformation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemInformation.Marshal(b, m, deterministic)
}

func (m *SystemInformation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemInformation.Merge(m, src)
}

func (m *SystemInformation) XXX_Size() int {
	return xxx_messageInfo_SystemInformation.Size(m)
}

func (m *SystemInformation) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemInformation.DiscardUnknown(m)
}

var xxx_messageInfo_SystemInformation proto.InternalMessageInfo

func (m *SystemInformation) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *SystemInformation) GetManufacturer() string {
	if m != nil {
		return m.Manufacturer
	}
	return ""
}

func (m *SystemInformation) GetProductName() string {
	if m != nil {
		return m.ProductName
	}
	return ""
}

func (m *SystemInformation) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *SystemInformation) GetSerialNumber() string {
	if m != nil {
		return m.SerialNumber
	}
	return ""
}

func (m *SystemInformation) GetSkuNumber() string {
	if m != nil {
		return m.SkuNumber
	}
	return ""
}

func (m *SystemInformation) GetFamily() string {
	if m != nil {
		return m.Family
	}
	return ""
}

type CPU struct {
	Manufacturer         string   `protobuf:"bytes,1,opt,name=manufacturer,proto3" json:"manufacturer,omitempty"`
	Version              string   `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CPU) Reset()         { *m = CPU{} }
func (m *CPU) String() string { return proto.CompactTextString(m) }
func (*CPU) ProtoMessage()    {}
func (*CPU) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *CPU) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CPU.Unmarshal(m, b)
}

func (m *CPU) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CPU.Marshal(b, m, deterministic)
}

func (m *CPU) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CPU.Merge(m, src)
}

func (m *CPU) XXX_Size() int {
	return xxx_messageInfo_CPU.Size(m)
}

func (m *CPU) XXX_DiscardUnknown() {
	xxx_messageInfo_CPU.DiscardUnknown(m)
}

var xxx_messageInfo_CPU proto.InternalMessageInfo

func (m *CPU) GetManufacturer() string {
	if m != nil {
		return m.Manufacturer
	}
	return ""
}

func (m *CPU) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type CreateServerRequest struct {
	SystemInformation    *SystemInformation `protobuf:"bytes,1,opt,name=system_information,json=systemInformation,proto3" json:"system_information,omitempty"`
	Cpu                  *CPU               `protobuf:"bytes,2,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Hostname             string             `protobuf:"bytes,3,opt,name=hostname,proto3" json:"hostname,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *CreateServerRequest) Reset()         { *m = CreateServerRequest{} }
func (m *CreateServerRequest) String() string { return proto.CompactTextString(m) }
func (*CreateServerRequest) ProtoMessage()    {}
func (*CreateServerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *CreateServerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateServerRequest.Unmarshal(m, b)
}

func (m *CreateServerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateServerRequest.Marshal(b, m, deterministic)
}

func (m *CreateServerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateServerRequest.Merge(m, src)
}

func (m *CreateServerRequest) XXX_Size() int {
	return xxx_messageInfo_CreateServerRequest.Size(m)
}

func (m *CreateServerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateServerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateServerRequest proto.InternalMessageInfo

func (m *CreateServerRequest) GetSystemInformation() *SystemInformation {
	if m != nil {
		return m.SystemInformation
	}
	return nil
}

func (m *CreateServerRequest) GetCpu() *CPU {
	if m != nil {
		return m.Cpu
	}
	return nil
}

func (m *CreateServerRequest) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

type Address struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *Address) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Address.Unmarshal(m, b)
}

func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Address.Marshal(b, m, deterministic)
}

func (m *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(m, src)
}

func (m *Address) XXX_Size() int {
	return xxx_messageInfo_Address.Size(m)
}

func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type CreateServerResponse struct {
	Wipe                 bool     `protobuf:"varint,1,opt,name=wipe,proto3" json:"wipe,omitempty"`
	InsecureWipe         bool     `protobuf:"varint,2,opt,name=insecure_wipe,json=insecureWipe,proto3" json:"insecure_wipe,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateServerResponse) Reset()         { *m = CreateServerResponse{} }
func (m *CreateServerResponse) String() string { return proto.CompactTextString(m) }
func (*CreateServerResponse) ProtoMessage()    {}
func (*CreateServerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *CreateServerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateServerResponse.Unmarshal(m, b)
}

func (m *CreateServerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateServerResponse.Marshal(b, m, deterministic)
}

func (m *CreateServerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateServerResponse.Merge(m, src)
}

func (m *CreateServerResponse) XXX_Size() int {
	return xxx_messageInfo_CreateServerResponse.Size(m)
}

func (m *CreateServerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateServerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateServerResponse proto.InternalMessageInfo

func (m *CreateServerResponse) GetWipe() bool {
	if m != nil {
		return m.Wipe
	}
	return false
}

func (m *CreateServerResponse) GetInsecureWipe() bool {
	if m != nil {
		return m.InsecureWipe
	}
	return false
}

type MarkServerAsWipedRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MarkServerAsWipedRequest) Reset()         { *m = MarkServerAsWipedRequest{} }
func (m *MarkServerAsWipedRequest) String() string { return proto.CompactTextString(m) }
func (*MarkServerAsWipedRequest) ProtoMessage()    {}
func (*MarkServerAsWipedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *MarkServerAsWipedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MarkServerAsWipedRequest.Unmarshal(m, b)
}

func (m *MarkServerAsWipedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MarkServerAsWipedRequest.Marshal(b, m, deterministic)
}

func (m *MarkServerAsWipedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarkServerAsWipedRequest.Merge(m, src)
}

func (m *MarkServerAsWipedRequest) XXX_Size() int {
	return xxx_messageInfo_MarkServerAsWipedRequest.Size(m)
}

func (m *MarkServerAsWipedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MarkServerAsWipedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MarkServerAsWipedRequest proto.InternalMessageInfo

func (m *MarkServerAsWipedRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type MarkServerAsWipedResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MarkServerAsWipedResponse) Reset()         { *m = MarkServerAsWipedResponse{} }
func (m *MarkServerAsWipedResponse) String() string { return proto.CompactTextString(m) }
func (*MarkServerAsWipedResponse) ProtoMessage()    {}
func (*MarkServerAsWipedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *MarkServerAsWipedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MarkServerAsWipedResponse.Unmarshal(m, b)
}

func (m *MarkServerAsWipedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MarkServerAsWipedResponse.Marshal(b, m, deterministic)
}

func (m *MarkServerAsWipedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarkServerAsWipedResponse.Merge(m, src)
}

func (m *MarkServerAsWipedResponse) XXX_Size() int {
	return xxx_messageInfo_MarkServerAsWipedResponse.Size(m)
}

func (m *MarkServerAsWipedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MarkServerAsWipedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MarkServerAsWipedResponse proto.InternalMessageInfo

type ReconcileServerAddressesRequest struct {
	Uuid                 string     `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Address              []*Address `protobuf:"bytes,2,rep,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ReconcileServerAddressesRequest) Reset()         { *m = ReconcileServerAddressesRequest{} }
func (m *ReconcileServerAddressesRequest) String() string { return proto.CompactTextString(m) }
func (*ReconcileServerAddressesRequest) ProtoMessage()    {}
func (*ReconcileServerAddressesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{7}
}

func (m *ReconcileServerAddressesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReconcileServerAddressesRequest.Unmarshal(m, b)
}

func (m *ReconcileServerAddressesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReconcileServerAddressesRequest.Marshal(b, m, deterministic)
}

func (m *ReconcileServerAddressesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReconcileServerAddressesRequest.Merge(m, src)
}

func (m *ReconcileServerAddressesRequest) XXX_Size() int {
	return xxx_messageInfo_ReconcileServerAddressesRequest.Size(m)
}

func (m *ReconcileServerAddressesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReconcileServerAddressesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReconcileServerAddressesRequest proto.InternalMessageInfo

func (m *ReconcileServerAddressesRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *ReconcileServerAddressesRequest) GetAddress() []*Address {
	if m != nil {
		return m.Address
	}
	return nil
}

type ReconcileServerAddressesResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReconcileServerAddressesResponse) Reset()         { *m = ReconcileServerAddressesResponse{} }
func (m *ReconcileServerAddressesResponse) String() string { return proto.CompactTextString(m) }
func (*ReconcileServerAddressesResponse) ProtoMessage()    {}
func (*ReconcileServerAddressesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{8}
}

func (m *ReconcileServerAddressesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReconcileServerAddressesResponse.Unmarshal(m, b)
}

func (m *ReconcileServerAddressesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReconcileServerAddressesResponse.Marshal(b, m, deterministic)
}

func (m *ReconcileServerAddressesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReconcileServerAddressesResponse.Merge(m, src)
}

func (m *ReconcileServerAddressesResponse) XXX_Size() int {
	return xxx_messageInfo_ReconcileServerAddressesResponse.Size(m)
}

func (m *ReconcileServerAddressesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReconcileServerAddressesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReconcileServerAddressesResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SystemInformation)(nil), "api.SystemInformation")
	proto.RegisterType((*CPU)(nil), "api.CPU")
	proto.RegisterType((*CreateServerRequest)(nil), "api.CreateServerRequest")
	proto.RegisterType((*Address)(nil), "api.Address")
	proto.RegisterType((*CreateServerResponse)(nil), "api.CreateServerResponse")
	proto.RegisterType((*MarkServerAsWipedRequest)(nil), "api.MarkServerAsWipedRequest")
	proto.RegisterType((*MarkServerAsWipedResponse)(nil), "api.MarkServerAsWipedResponse")
	proto.RegisterType((*ReconcileServerAddressesRequest)(nil), "api.ReconcileServerAddressesRequest")
	proto.RegisterType((*ReconcileServerAddressesResponse)(nil), "api.ReconcileServerAddressesResponse")
}

func init() {
	proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c)
}

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 541 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xcd, 0x6e, 0xd3, 0x4c,
	0x14, 0x55, 0x92, 0x36, 0x3f, 0x37, 0xf9, 0x16, 0x99, 0x0f, 0x55, 0x6e, 0x50, 0xa1, 0x98, 0x1f,
	0x75, 0x93, 0x58, 0x0a, 0x0b, 0xd6, 0x21, 0x62, 0x51, 0x21, 0x4a, 0xe5, 0xaa, 0x42, 0x42, 0x42,
	0xd1, 0xc4, 0xbe, 0x49, 0x47, 0xb1, 0x67, 0x86, 0xf9, 0x29, 0xca, 0x23, 0xb0, 0xe7, 0x11, 0x79,
	0x10, 0xe4, 0xf1, 0x98, 0x26, 0x4a, 0x42, 0x77, 0xbe, 0xe7, 0xdc, 0x39, 0x73, 0xcf, 0xb9, 0x23,
	0x43, 0x87, 0x4a, 0x36, 0x92, 0x4a, 0x18, 0x41, 0x1a, 0x54, 0xb2, 0xf0, 0x77, 0x0d, 0xfa, 0x37,
	0x6b, 0x6d, 0x30, 0xbf, 0xe4, 0x0b, 0xa1, 0x72, 0x6a, 0x98, 0xe0, 0x84, 0xc0, 0x91, 0xb5, 0x2c,
	0x0d, 0x6a, 0xe7, 0xb5, 0x8b, 0x4e, 0xec, 0xbe, 0x49, 0x08, 0xbd, 0x9c, 0x72, 0xbb, 0xa0, 0x89,
	0xb1, 0x0a, 0x55, 0x50, 0x77, 0xdc, 0x16, 0x46, 0x5e, 0x40, 0x4f, 0x2a, 0x91, 0xda, 0xc4, 0xcc,
	0x38, 0xcd, 0x31, 0x68, 0xb8, 0x9e, 0xae, 0xc7, 0xae, 0x68, 0x8e, 0x24, 0x80, 0xd6, 0x3d, 0x2a,
	0xcd, 0x04, 0x0f, 0x8e, 0x1c, 0x5b, 0x95, 0xe4, 0x25, 0xfc, 0xa7, 0x51, 0x31, 0x9a, 0xcd, 0xb8,
	0xcd, 0xe7, 0xa8, 0x82, 0xe3, 0xf2, 0x86, 0x12, 0xbc, 0x72, 0x18, 0x39, 0x03, 0xd0, 0x2b, 0x5b,
	0x75, 0x34, 0x5d, 0x47, 0x47, 0xaf, 0xac, 0xa7, 0x4f, 0xa0, 0xb9, 0xa0, 0x39, 0xcb, 0xd6, 0x41,
	0xcb, 0x51, 0xbe, 0x0a, 0xa7, 0xd0, 0x98, 0x5e, 0xdf, 0xee, 0x78, 0xa8, 0xed, 0xf1, 0xb0, 0x31,
	0x60, 0x7d, 0x6b, 0xc0, 0xf0, 0x57, 0x0d, 0xfe, 0x9f, 0x2a, 0xa4, 0x06, 0x6f, 0x50, 0xdd, 0xa3,
	0x8a, 0xf1, 0xbb, 0x45, 0x6d, 0xc8, 0x07, 0x20, 0xda, 0x45, 0x38, 0x63, 0x0f, 0x19, 0x3a, 0xed,
	0xee, 0xf8, 0x64, 0x54, 0x04, 0xbe, 0x93, 0x70, 0xdc, 0xd7, 0x3b, 0xa1, 0x0f, 0xa0, 0x91, 0x48,
	0xeb, 0x2e, 0xed, 0x8e, 0xdb, 0xee, 0xdc, 0xf4, 0xfa, 0x36, 0x2e, 0x40, 0x32, 0x80, 0xf6, 0x9d,
	0xd0, 0x66, 0x23, 0xd4, 0xbf, 0x75, 0xf8, 0x0e, 0x5a, 0x93, 0x34, 0x55, 0xa8, 0x75, 0xb1, 0x37,
	0xb3, 0x96, 0x58, 0xed, 0xad, 0xf8, 0x2e, 0xfc, 0xd0, 0x92, 0xae, 0xfc, 0xf8, 0x32, 0xfc, 0x0c,
	0x4f, 0xb6, 0xed, 0x68, 0x29, 0xb8, 0xc6, 0x42, 0xe5, 0x07, 0xf3, 0x2a, 0xed, 0xd8, 0x7d, 0x17,
	0xcb, 0x61, 0x5c, 0x63, 0x62, 0x15, 0xce, 0x1c, 0x59, 0x77, 0x64, 0xaf, 0x02, 0xbf, 0x30, 0x89,
	0xe1, 0x08, 0x82, 0x4f, 0x54, 0xad, 0x4a, 0xb9, 0x89, 0x2e, 0xb0, 0xb4, 0x0a, 0x69, 0xcf, 0x93,
	0x0a, 0x9f, 0xc2, 0xe9, 0x9e, 0xfe, 0x72, 0x8a, 0xf0, 0x1b, 0x3c, 0x8f, 0x31, 0x11, 0x3c, 0x61,
	0x99, 0x1f, 0xd0, 0xbb, 0x44, 0xfd, 0x0f, 0x4d, 0xf2, 0x66, 0xd3, 0x6e, 0xe3, 0xa2, 0x3b, 0xee,
	0xb9, 0x24, 0xfd, 0xd9, 0x07, 0xf3, 0x21, 0x9c, 0x1f, 0x96, 0x2f, 0x47, 0x18, 0xff, 0xac, 0xc3,
	0xf1, 0x64, 0x89, 0xdc, 0x90, 0x29, 0xf4, 0x36, 0xa3, 0x22, 0x41, 0xb9, 0x9e, 0xdd, 0xc7, 0x30,
	0x38, 0xdd, 0xc3, 0xf8, 0x5c, 0x63, 0xe8, 0xef, 0xd8, 0x25, 0x67, 0xae, 0xff, 0x50, 0x6c, 0x83,
	0x67, 0x87, 0x68, 0xaf, 0xb9, 0x84, 0xe0, 0x90, 0x0d, 0xf2, 0xca, 0x9d, 0x7d, 0x24, 0xc4, 0xc1,
	0xeb, 0x47, 0xba, 0xca, 0x8b, 0xde, 0x7f, 0xfc, 0x7a, 0xb9, 0x64, 0xe6, 0xce, 0xce, 0x47, 0x89,
	0xc8, 0x23, 0x43, 0x33, 0xa1, 0x87, 0xe5, 0x1b, 0xd6, 0x91, 0x66, 0x29, 0x2a, 0x11, 0x51, 0x29,
	0xa3, 0x1c, 0x0d, 0xcd, 0x86, 0x89, 0xe0, 0x46, 0x89, 0x2c, 0x43, 0x35, 0xcc, 0x29, 0xa7, 0x4b,
	0x54, 0x11, 0xe3, 0x06, 0x15, 0xa7, 0x59, 0x44, 0x25, 0x9b, 0x37, 0xdd, 0x1f, 0xe8, 0xed, 0x9f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x76, 0x56, 0x70, 0x65, 0x8e, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConnInterface
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AgentClient interface {
	CreateServer(ctx context.Context, in *CreateServerRequest, opts ...grpc.CallOption) (*CreateServerResponse, error)
	MarkServerAsWiped(ctx context.Context, in *MarkServerAsWipedRequest, opts ...grpc.CallOption) (*MarkServerAsWipedResponse, error)
	ReconcileServerAddresses(ctx context.Context, in *ReconcileServerAddressesRequest, opts ...grpc.CallOption) (*ReconcileServerAddressesResponse, error)
}

type agentClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentClient(cc grpc.ClientConnInterface) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) CreateServer(ctx context.Context, in *CreateServerRequest, opts ...grpc.CallOption) (*CreateServerResponse, error) {
	out := new(CreateServerResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/CreateServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) MarkServerAsWiped(ctx context.Context, in *MarkServerAsWipedRequest, opts ...grpc.CallOption) (*MarkServerAsWipedResponse, error) {
	out := new(MarkServerAsWipedResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/MarkServerAsWiped", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) ReconcileServerAddresses(ctx context.Context, in *ReconcileServerAddressesRequest, opts ...grpc.CallOption) (*ReconcileServerAddressesResponse, error) {
	out := new(ReconcileServerAddressesResponse)
	err := c.cc.Invoke(ctx, "/api.Agent/ReconcileServerAddresses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServer is the server API for Agent service.
type AgentServer interface {
	CreateServer(context.Context, *CreateServerRequest) (*CreateServerResponse, error)
	MarkServerAsWiped(context.Context, *MarkServerAsWipedRequest) (*MarkServerAsWipedResponse, error)
	ReconcileServerAddresses(context.Context, *ReconcileServerAddressesRequest) (*ReconcileServerAddressesResponse, error)
}

// UnimplementedAgentServer can be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct {
}

func (*UnimplementedAgentServer) CreateServer(ctx context.Context, req *CreateServerRequest) (*CreateServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateServer not implemented")
}

func (*UnimplementedAgentServer) MarkServerAsWiped(ctx context.Context, req *MarkServerAsWipedRequest) (*MarkServerAsWipedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkServerAsWiped not implemented")
}

func (*UnimplementedAgentServer) ReconcileServerAddresses(ctx context.Context, req *ReconcileServerAddressesRequest) (*ReconcileServerAddressesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReconcileServerAddresses not implemented")
}

func RegisterAgentServer(s *grpc.Server, srv AgentServer) {
	s.RegisterService(&_Agent_serviceDesc, srv)
}

func _Agent_CreateServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).CreateServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/CreateServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).CreateServer(ctx, req.(*CreateServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_MarkServerAsWiped_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkServerAsWipedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).MarkServerAsWiped(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/MarkServerAsWiped",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).MarkServerAsWiped(ctx, req.(*MarkServerAsWipedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_ReconcileServerAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReconcileServerAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).ReconcileServerAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Agent/ReconcileServerAddresses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).ReconcileServerAddresses(ctx, req.(*ReconcileServerAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Agent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateServer",
			Handler:    _Agent_CreateServer_Handler,
		},
		{
			MethodName: "MarkServerAsWiped",
			Handler:    _Agent_MarkServerAsWiped_Handler,
		},
		{
			MethodName: "ReconcileServerAddresses",
			Handler:    _Agent_ReconcileServerAddresses_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
