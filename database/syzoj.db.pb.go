// Code generated by protoc-gen-go. DO NOT EDIT.
// source: syzoj.db.proto

package database

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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
	Id                   *UserRef             `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	UserName             *string              `protobuf:"bytes,2,opt,name=user_name,json=userName" json:"user_name,omitempty"`
	Auth                 *model.UserAuth      `protobuf:"bytes,3,opt,name=auth" json:"auth,omitempty"`
	RegisterTime         *timestamp.Timestamp `protobuf:"bytes,4,opt,name=register_time,json=registerTime" json:"register_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
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

func (m *User) GetRegisterTime() *timestamp.Timestamp {
	if m != nil {
		return m.RegisterTime
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
	Id                   *ProblemRef          `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	User                 *UserRef             `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	CreateTime           *timestamp.Timestamp `protobuf:"bytes,4,opt,name=create_time,json=createTime" json:"create_time,omitempty"`
	Problem              *model.Problem       `protobuf:"bytes,5,opt,name=problem" json:"problem,omitempty"`
	Title                *string              `protobuf:"bytes,6,opt,name=title" json:"title,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
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

func (m *Problem) GetUser() UserRef {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

func (m *Problem) GetCreateTime() *timestamp.Timestamp {
	if m != nil {
		return m.CreateTime
	}
	return nil
}

func (m *Problem) GetProblem() *model.Problem {
	if m != nil {
		return m.Problem
	}
	return nil
}

func (m *Problem) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

// A problem displayed in public.
type ProblemEntry struct {
	Id                   *ProblemEntryRef `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Title                *string          `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	Problem              *ProblemRef      `protobuf:"bytes,3,opt,name=problem" json:"problem,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ProblemEntry) Reset()         { *m = ProblemEntry{} }
func (m *ProblemEntry) String() string { return proto.CompactTextString(m) }
func (*ProblemEntry) ProtoMessage()    {}
func (*ProblemEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{3}
}

func (m *ProblemEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemEntry.Unmarshal(m, b)
}
func (m *ProblemEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemEntry.Marshal(b, m, deterministic)
}
func (m *ProblemEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemEntry.Merge(m, src)
}
func (m *ProblemEntry) XXX_Size() int {
	return xxx_messageInfo_ProblemEntry.Size(m)
}
func (m *ProblemEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemEntry.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemEntry proto.InternalMessageInfo

func (m *ProblemEntry) GetId() ProblemEntryRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *ProblemEntry) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *ProblemEntry) GetProblem() ProblemRef {
	if m != nil && m.Problem != nil {
		return *m.Problem
	}
	return ""
}

type Submission struct {
	Id                   *SubmissionRef `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	ProblemJudger        *string        `protobuf:"bytes,2,opt,name=problem_judger,json=problemJudger" json:"problem_judger,omitempty"`
	User                 *UserRef       `protobuf:"bytes,3,opt,name=user" json:"user,omitempty"`
	Data                 *any.Any       `protobuf:"bytes,16,opt,name=data" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Submission) Reset()         { *m = Submission{} }
func (m *Submission) String() string { return proto.CompactTextString(m) }
func (*Submission) ProtoMessage()    {}
func (*Submission) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{4}
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

func (m *Submission) GetProblemJudger() string {
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

func (m *Submission) GetData() *any.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

type Judger struct {
	Id                   *JudgerRef `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Token                *string    `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Judger) Reset()         { *m = Judger{} }
func (m *Judger) String() string { return proto.CompactTextString(m) }
func (*Judger) ProtoMessage()    {}
func (*Judger) Descriptor() ([]byte, []int) {
	return fileDescriptor_549dafe8664ea10e, []int{5}
}

func (m *Judger) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Judger.Unmarshal(m, b)
}
func (m *Judger) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Judger.Marshal(b, m, deterministic)
}
func (m *Judger) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Judger.Merge(m, src)
}
func (m *Judger) XXX_Size() int {
	return xxx_messageInfo_Judger.Size(m)
}
func (m *Judger) XXX_DiscardUnknown() {
	xxx_messageInfo_Judger.DiscardUnknown(m)
}

var xxx_messageInfo_Judger proto.InternalMessageInfo

func (m *Judger) GetId() JudgerRef {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Judger) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "syzoj.db.User")
	proto.RegisterType((*Device)(nil), "syzoj.db.Device")
	proto.RegisterType((*Problem)(nil), "syzoj.db.Problem")
	proto.RegisterType((*ProblemEntry)(nil), "syzoj.db.ProblemEntry")
	proto.RegisterType((*Submission)(nil), "syzoj.db.Submission")
	proto.RegisterType((*Judger)(nil), "syzoj.db.Judger")
}

func init() { proto.RegisterFile("syzoj.db.proto", fileDescriptor_549dafe8664ea10e) }

var fileDescriptor_549dafe8664ea10e = []byte{
	// 459 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xdf, 0x6a, 0x13, 0x4f,
	0x14, 0xc7, 0x99, 0xed, 0x36, 0x69, 0x4e, 0x9a, 0xf0, 0xfb, 0x0d, 0x15, 0xd7, 0xdc, 0x34, 0x2c,
	0x16, 0x23, 0xd2, 0x09, 0xf8, 0x06, 0x2d, 0x15, 0xd4, 0x0b, 0x91, 0x51, 0x6f, 0xbc, 0x30, 0xcc,
	0x66, 0x4f, 0x36, 0x53, 0xb3, 0x33, 0x61, 0x77, 0x56, 0x88, 0x8f, 0xe0, 0x45, 0xdf, 0xc5, 0x97,
	0x10, 0x7c, 0x2b, 0x99, 0x3f, 0x1b, 0xd3, 0xc6, 0x22, 0xde, 0x0c, 0x9c, 0x73, 0xbe, 0xe7, 0x7b,
	0x3e, 0xe7, 0x0c, 0x0c, 0xeb, 0xcd, 0x57, 0x7d, 0xcd, 0xf2, 0x8c, 0xad, 0x2b, 0x6d, 0x34, 0x3d,
	0x6a, 0xe3, 0xd1, 0x20, 0xcf, 0x4a, 0x9d, 0xe3, 0xca, 0x17, 0x46, 0xff, 0xfb, 0xc2, 0x6e, 0xea,
	0xb4, 0xd0, 0xba, 0x58, 0xe1, 0xd4, 0x45, 0x59, 0xb3, 0x98, 0x1a, 0x59, 0x62, 0x6d, 0x44, 0xb9,
	0x0e, 0x82, 0x47, 0x77, 0x05, 0x42, 0x6d, 0x7c, 0x29, 0xfd, 0x4e, 0x20, 0xfe, 0x50, 0x63, 0x45,
	0x87, 0x10, 0xc9, 0x3c, 0x21, 0x63, 0x32, 0xe9, 0xf1, 0x48, 0xe6, 0xf4, 0x31, 0xf4, 0x9a, 0x1a,
	0xab, 0x99, 0x12, 0x25, 0x26, 0x91, 0x4d, 0x5f, 0x76, 0xbf, 0xfd, 0xbc, 0x21, 0xd1, 0x98, 0xf0,
	0x23, 0x5b, 0x79, 0x23, 0x4a, 0xa4, 0x4f, 0x21, 0x16, 0x8d, 0x59, 0x26, 0x07, 0x63, 0x32, 0xe9,
	0x3f, 0x7f, 0xc0, 0x76, 0xe1, 0xac, 0xed, 0x45, 0x63, 0x96, 0xdc, 0x49, 0xe8, 0x4b, 0x18, 0x54,
	0x58, 0xc8, 0xda, 0x60, 0x35, 0xb3, 0x80, 0x49, 0xec, 0x7a, 0x46, 0xcc, 0xc3, 0xb1, 0x16, 0x8e,
	0xbd, 0x6f, 0xe9, 0xc3, 0xc0, 0x84, 0xf0, 0xe3, 0xb6, 0xd3, 0xd6, 0xd2, 0x05, 0x74, 0xae, 0xf0,
	0x8b, 0x9c, 0xe3, 0x1e, 0xf4, 0x29, 0xc4, 0x16, 0x2d, 0xf0, 0xf6, 0x6d, 0x7b, 0x87, 0xba, 0x14,
	0x77, 0x2f, 0x7d, 0x06, 0xb1, 0x54, 0x0b, 0x1d, 0x78, 0x1f, 0xde, 0xe2, 0xf5, 0x9e, 0xaf, 0xd4,
	0x42, 0x73, 0x27, 0x4a, 0x7f, 0x10, 0xe8, 0xbe, 0xad, 0x74, 0xb6, 0xc2, 0xf2, 0xde, 0x49, 0x07,
	0xf7, 0x4d, 0xba, 0x82, 0xfe, 0xbc, 0x42, 0x61, 0xf0, 0x9f, 0x97, 0x05, 0xdf, 0x67, 0x2b, 0x94,
	0x41, 0x77, 0xed, 0x09, 0x92, 0x43, 0xe7, 0x70, 0x72, 0x0b, 0x39, 0xd0, 0xf1, 0x56, 0x44, 0x4f,
	0xe0, 0xd0, 0x48, 0xb3, 0xc2, 0xa4, 0xe3, 0x48, 0x7d, 0x90, 0x7e, 0x82, 0xe3, 0xa0, 0x7c, 0xa1,
	0x4c, 0xb5, 0xd9, 0x5b, 0x66, 0xdb, 0x15, 0xed, 0x74, 0xd1, 0xc9, 0xef, 0xd9, 0x7e, 0xcb, 0xa1,
	0x25, 0xec, 0xd1, 0x36, 0xbb, 0x9d, 0x9a, 0xde, 0x10, 0x80, 0x77, 0x4d, 0x56, 0xca, 0xba, 0x96,
	0x5a, 0xed, 0xd9, 0x9f, 0xc1, 0x30, 0x28, 0x67, 0xd7, 0x4d, 0x5e, 0xb4, 0xff, 0xc3, 0x07, 0x21,
	0xfb, 0xda, 0x25, 0xff, 0x7e, 0xd2, 0x09, 0xc4, 0xb9, 0x30, 0x22, 0xf9, 0x2f, 0x5c, 0xe2, 0xee,
	0x2d, 0x2f, 0xd4, 0x86, 0x3b, 0x45, 0xca, 0xa0, 0x13, 0x4c, 0xff, 0xb4, 0xaa, 0xfe, 0x8c, 0x6a,
	0xbb, 0xaa, 0x0d, 0x2e, 0x9f, 0x7c, 0x3c, 0x2b, 0xa4, 0x59, 0x36, 0x19, 0x9b, 0xeb, 0x72, 0xea,
	0x2e, 0xec, 0xdf, 0x73, 0x55, 0x9c, 0x17, 0x7a, 0x6a, 0x4d, 0x33, 0x51, 0xe3, 0xaf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xfe, 0x41, 0x3b, 0x90, 0xa7, 0x03, 0x00, 0x00,
}
