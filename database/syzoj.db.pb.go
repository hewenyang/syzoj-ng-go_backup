// Code generated by protoc-gen-go. DO NOT EDIT.
// source: syzoj.db.proto

package database

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/syzoj/protoc-gen-gotype/gotype"
	_ "github.com/syzoj/syzoj-ng-go/database/protoc-gen-dbmodel/dbmodel"
	model "github.com/syzoj/syzoj-ng-go/model"
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

type User struct {
	Id                   *UserRef        `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	UserName             *string         `protobuf:"bytes,2,opt,name=user_name,json=userName" json:"user_name,omitempty"`
	Auth                 *model.UserAuth `protobuf:"bytes,3,opt,name=auth" json:"auth,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{0}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() UserRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *User) GetUserName() string {
	if m != nil && m.UserName != nil {
		return *m.UserName
	}
	return ""
}

func (m *User) GetAuth() *model.UserAuth {
	if m != nil {
		return m.Auth
	}
	return nil
}

type Device struct {
	Id                   *DeviceRef        `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	User                 *UserRef          `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
	Info                 *model.DeviceInfo `protobuf:"bytes,3,opt,name=info" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Device) Reset()         { *m = Device{} }
func (m *Device) String() string { return proto.CompactTextString(m) }
func (*Device) ProtoMessage()    {}
func (*Device) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{1}
}

func (m *Device) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Device.Unmarshal(m, b)
}
func (m *Device) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Device.Marshal(b, m, deterministic)
}
func (m *Device) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Device.Merge(m, src)
}
func (m *Device) XXX_Size() int {
	return xxx_messageInfo_Device.Size(m)
}
func (m *Device) XXX_DiscardUnknown() {
	xxx_messageInfo_Device.DiscardUnknown(m)
}

var xxx_messageInfo_Device proto.InternalMessageInfo

func (m *Device) GetId() DeviceRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Device) GetUser() UserRef {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

func (m *Device) GetInfo() *model.DeviceInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

type Problem struct {
	Id                   *ProblemRef `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Title                *string     `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	User                 *UserRef    `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Problem) Reset()         { *m = Problem{} }
func (m *Problem) String() string { return proto.CompactTextString(m) }
func (*Problem) ProtoMessage()    {}
func (*Problem) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{2}
}

func (m *Problem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Problem.Unmarshal(m, b)
}
func (m *Problem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Problem.Marshal(b, m, deterministic)
}
func (m *Problem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Problem.Merge(m, src)
}
func (m *Problem) XXX_Size() int {
	return xxx_messageInfo_Problem.Size(m)
}
func (m *Problem) XXX_DiscardUnknown() {
	xxx_messageInfo_Problem.DiscardUnknown(m)
}

var xxx_messageInfo_Problem proto.InternalMessageInfo

func (m *Problem) GetId() ProblemRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Problem) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *Problem) GetUser() UserRef {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

type ProblemSource struct {
	Id                   *ProblemSourceRef    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Source               *model.ProblemSource `protobuf:"bytes,16,opt,name=source" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ProblemSource) Reset()         { *m = ProblemSource{} }
func (m *ProblemSource) String() string { return proto.CompactTextString(m) }
func (*ProblemSource) ProtoMessage()    {}
func (*ProblemSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{3}
}

func (m *ProblemSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemSource.Unmarshal(m, b)
}
func (m *ProblemSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemSource.Marshal(b, m, deterministic)
}
func (m *ProblemSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemSource.Merge(m, src)
}
func (m *ProblemSource) XXX_Size() int {
	return xxx_messageInfo_ProblemSource.Size(m)
}
func (m *ProblemSource) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemSource.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemSource proto.InternalMessageInfo

func (m *ProblemSource) GetId() ProblemSourceRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *ProblemSource) GetSource() *model.ProblemSource {
	if m != nil {
		return m.Source
	}
	return nil
}

type ProblemJudger struct {
	Id                   *ProblemJudgerRef `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Problem              *ProblemRef       `protobuf:"bytes,2,opt,name=problem" json:"problem,omitempty"`
	User                 *UserRef          `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	Type                 *string           `protobuf:"bytes,4,opt,name=type" json:"type,omitempty"`
	Data                 *model.Any        `protobuf:"bytes,16,opt,name=data" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ProblemJudger) Reset()         { *m = ProblemJudger{} }
func (m *ProblemJudger) String() string { return proto.CompactTextString(m) }
func (*ProblemJudger) ProtoMessage()    {}
func (*ProblemJudger) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{4}
}

func (m *ProblemJudger) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemJudger.Unmarshal(m, b)
}
func (m *ProblemJudger) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemJudger.Marshal(b, m, deterministic)
}
func (m *ProblemJudger) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemJudger.Merge(m, src)
}
func (m *ProblemJudger) XXX_Size() int {
	return xxx_messageInfo_ProblemJudger.Size(m)
}
func (m *ProblemJudger) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemJudger.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemJudger proto.InternalMessageInfo

func (m *ProblemJudger) GetId() ProblemJudgerRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *ProblemJudger) GetProblem() ProblemRef {
	if m != nil && m.Problem != nil {
		return *m.Problem
	}
	return ""
}

func (m *ProblemJudger) GetUser() UserRef {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

func (m *ProblemJudger) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func (m *ProblemJudger) GetData() *model.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

type ProblemStatement struct {
	Id                   *ProblemStatementRef    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Problem              *ProblemRef             `protobuf:"bytes,2,opt,name=problem" json:"problem,omitempty"`
	User                 *UserRef                `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	Data                 *model.ProblemStatement `protobuf:"bytes,16,opt,name=data" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *ProblemStatement) Reset()         { *m = ProblemStatement{} }
func (m *ProblemStatement) String() string { return proto.CompactTextString(m) }
func (*ProblemStatement) ProtoMessage()    {}
func (*ProblemStatement) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{5}
}

func (m *ProblemStatement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemStatement.Unmarshal(m, b)
}
func (m *ProblemStatement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemStatement.Marshal(b, m, deterministic)
}
func (m *ProblemStatement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemStatement.Merge(m, src)
}
func (m *ProblemStatement) XXX_Size() int {
	return xxx_messageInfo_ProblemStatement.Size(m)
}
func (m *ProblemStatement) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemStatement.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemStatement proto.InternalMessageInfo

func (m *ProblemStatement) GetId() ProblemStatementRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *ProblemStatement) GetProblem() ProblemRef {
	if m != nil && m.Problem != nil {
		return *m.Problem
	}
	return ""
}

func (m *ProblemStatement) GetUser() UserRef {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

func (m *ProblemStatement) GetData() *model.ProblemStatement {
	if m != nil {
		return m.Data
	}
	return nil
}

type Submission struct {
	Id                   *SubmissionRef    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	ProblemJudger        *ProblemJudgerRef `protobuf:"bytes,2,opt,name=problem_judger,json=problemJudger" json:"problem_judger,omitempty"`
	User                 *UserRef          `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	Data                 *model.Any        `protobuf:"bytes,16,opt,name=data" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Submission) Reset()         { *m = Submission{} }
func (m *Submission) String() string { return proto.CompactTextString(m) }
func (*Submission) ProtoMessage()    {}
func (*Submission) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{6}
}

func (m *Submission) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Submission.Unmarshal(m, b)
}
func (m *Submission) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Submission.Marshal(b, m, deterministic)
}
func (m *Submission) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Submission.Merge(m, src)
}
func (m *Submission) XXX_Size() int {
	return xxx_messageInfo_Submission.Size(m)
}
func (m *Submission) XXX_DiscardUnknown() {
	xxx_messageInfo_Submission.DiscardUnknown(m)
}

var xxx_messageInfo_Submission proto.InternalMessageInfo

func (m *Submission) GetId() SubmissionRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Submission) GetProblemJudger() ProblemJudgerRef {
	if m != nil && m.ProblemJudger != nil {
		return *m.ProblemJudger
	}
	return ""
}

func (m *Submission) GetUser() UserRef {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

func (m *Submission) GetData() *model.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*User)(nil), "syzoj.db.User")
	proto.RegisterType((*Device)(nil), "syzoj.db.Device")
	proto.RegisterType((*Problem)(nil), "syzoj.db.Problem")
	proto.RegisterType((*ProblemSource)(nil), "syzoj.db.ProblemSource")
	proto.RegisterType((*ProblemJudger)(nil), "syzoj.db.ProblemJudger")
	proto.RegisterType((*ProblemStatement)(nil), "syzoj.db.ProblemStatement")
	proto.RegisterType((*Submission)(nil), "syzoj.db.Submission")
}

func init() { proto.RegisterFile("syzoj.db.proto", fileDescriptor_549dafe8664ea10e) }

var fileDescriptor_549dafe8664ea10e = []byte{
	// 501 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x92, 0xc1, 0x6e, 0xd3, 0x30,
	0x1c, 0xc6, 0x95, 0x2e, 0xac, 0xeb, 0x7f, 0xa4, 0xea, 0xcc, 0xd0, 0xa2, 0xc2, 0xa0, 0x8b, 0x86,
	0xd6, 0x09, 0xad, 0x13, 0x13, 0xe2, 0x82, 0x84, 0xb4, 0x09, 0x0e, 0xe3, 0x30, 0x41, 0xa6, 0x1e,
	0xc2, 0x65, 0x38, 0x8d, 0x9b, 0x7a, 0x6a, 0xe2, 0x2a, 0x71, 0x26, 0x95, 0x47, 0xc8, 0x61, 0xf7,
	0x3c, 0x07, 0x27, 0x78, 0x03, 0xde, 0x0a, 0xc5, 0x76, 0x42, 0x3c, 0x6d, 0x1a, 0x97, 0x5d, 0xaa,
	0xfa, 0xff, 0x7d, 0xf5, 0xf7, 0xfb, 0xfe, 0x2e, 0x74, 0xd3, 0xe5, 0x0f, 0x76, 0x39, 0x0a, 0xfc,
	0xd1, 0x22, 0x61, 0x9c, 0xa1, 0xb5, 0xea, 0xdc, 0x7f, 0x1c, 0x32, 0xbe, 0x5c, 0x10, 0x39, 0xef,
	0x5b, 0x81, 0x1f, 0xb1, 0x80, 0xcc, 0xd5, 0x71, 0x43, 0xda, 0x1a, 0x23, 0xe7, 0xda, 0x00, 0x73,
	0x9c, 0x92, 0x04, 0x3d, 0x83, 0x16, 0x0d, 0x6c, 0x63, 0x60, 0x0c, 0x3b, 0x27, 0xeb, 0x79, 0xe1,
	0xb5, 0xcb, 0xa9, 0x4b, 0xa6, 0x6e, 0x8b, 0x06, 0xe8, 0x03, 0x74, 0xb2, 0x94, 0x24, 0x17, 0x31,
	0x8e, 0x88, 0xdd, 0x12, 0x9e, 0x9d, 0xfc, 0xcf, 0xb5, 0xf1, 0xbc, 0x9e, 0x0e, 0xae, 0x70, 0x32,
	0x99, 0xe1, 0x64, 0xf8, 0xee, 0xed, 0xfe, 0x60, 0x7c, 0x76, 0xfa, 0x75, 0xfc, 0xc9, 0x5d, 0x2b,
	0xd5, 0x33, 0x1c, 0x11, 0xb4, 0x0f, 0x26, 0xce, 0xf8, 0xcc, 0x5e, 0x19, 0x18, 0xc3, 0xf5, 0xa3,
	0xa7, 0xa3, 0x26, 0x47, 0x99, 0x73, 0x9c, 0xf1, 0x99, 0x2b, 0x2c, 0x4e, 0x06, 0xab, 0x1f, 0xc9,
	0x15, 0x9d, 0x10, 0xb4, 0xdd, 0x20, 0xb2, 0xf2, 0xc2, 0xeb, 0xc8, 0x79, 0xc5, 0xf4, 0x12, 0xcc,
	0xf2, 0x7e, 0x85, 0xa3, 0x21, 0x0b, 0x01, 0xbd, 0x06, 0x93, 0xc6, 0x53, 0xa6, 0x42, 0xb7, 0xb4,
	0x50, 0x79, 0xd5, 0x69, 0x3c, 0x65, 0xae, 0x30, 0x39, 0xdf, 0xa1, 0xfd, 0x25, 0x61, 0xfe, 0x9c,
	0x44, 0xe8, 0x45, 0x23, 0xb7, 0x9b, 0x17, 0x1e, 0x28, 0xa1, 0x0a, 0xde, 0x84, 0x47, 0x9c, 0xf2,
	0xb9, 0x5a, 0x84, 0x2b, 0x0f, 0x35, 0xce, 0xca, 0x1d, 0x38, 0x0e, 0x05, 0x4b, 0x5d, 0x74, 0xce,
	0xb2, 0x64, 0x42, 0xd0, 0x6e, 0x23, 0x67, 0x33, 0x2f, 0xbc, 0x9e, 0x26, 0x57, 0x69, 0x47, 0xb0,
	0x9a, 0x8a, 0x81, 0xdd, 0x13, 0x3d, 0xfa, 0x5a, 0x0f, 0xfd, 0x27, 0xca, 0xe9, 0xfc, 0x36, 0xea,
	0xac, 0xcf, 0x59, 0x10, 0x92, 0xe4, 0xce, 0x2c, 0x29, 0x57, 0x59, 0x43, 0x68, 0x2f, 0xe4, 0x5c,
	0x6d, 0xf5, 0x66, 0xfd, 0x4a, 0xbe, 0xb7, 0x2d, 0x42, 0x60, 0x96, 0xff, 0x43, 0xdb, 0x14, 0x3b,
	0x12, 0xdf, 0xd1, 0x2e, 0x98, 0x01, 0xe6, 0x58, 0x15, 0xe9, 0x69, 0x45, 0x8e, 0xe3, 0xa5, 0x2b,
	0x54, 0xe7, 0x97, 0x01, 0xf5, 0x26, 0x38, 0xe6, 0x24, 0x22, 0x31, 0x47, 0x7b, 0x0d, 0xfe, 0xad,
	0xbc, 0xf0, 0x9e, 0xdc, 0x74, 0x3c, 0x40, 0x85, 0x37, 0x1a, 0xee, 0xf6, 0xad, 0x7b, 0xaf, 0xe3,
	0x25, 0xfb, 0x4f, 0x03, 0xe0, 0x3c, 0xf3, 0x23, 0x9a, 0xa6, 0x94, 0xc5, 0x68, 0xa7, 0x41, 0xbd,
	0x91, 0x17, 0x9e, 0xf5, 0x4f, 0xab, 0x78, 0xdf, 0x43, 0x57, 0x01, 0x5d, 0x5c, 0x8a, 0xb7, 0x50,
	0xd8, 0xb7, 0x3f, 0x92, 0xb5, 0xd0, 0x5e, 0xf5, 0xde, 0x0a, 0xff, 0xb5, 0xf1, 0x93, 0xbd, 0x6f,
	0xaf, 0x42, 0xca, 0x67, 0x99, 0x3f, 0x9a, 0xb0, 0xe8, 0x50, 0x78, 0xe4, 0xe7, 0x41, 0x1c, 0x1e,
	0x84, 0xec, 0xb0, 0x34, 0xf9, 0x38, 0x25, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc1, 0xe8, 0x37,
	0x9b, 0x77, 0x04, 0x00, 0x00,
}
