// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: syzoj.service.problem.kafka.proto

package kafka // import "github.com/syzoj/syzoj-ng-go/service/problem/kafka"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ProblemEvent struct {
	ProblemId *string `protobuf:"bytes,1,opt,name=problem_id,json=problemId" json:"problem_id,omitempty"`
	// Types that are valid to be assigned to Event:
	//	*ProblemEvent_Submission
	Event                isProblemEvent_Event `protobuf_oneof:"event"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ProblemEvent) Reset()         { *m = ProblemEvent{} }
func (m *ProblemEvent) String() string { return proto.CompactTextString(m) }
func (*ProblemEvent) ProtoMessage()    {}
func (*ProblemEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_syzoj_service_problem_kafka_e523ec2f3887079d, []int{0}
}
func (m *ProblemEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProblemEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProblemEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ProblemEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemEvent.Merge(dst, src)
}
func (m *ProblemEvent) XXX_Size() int {
	return m.Size()
}
func (m *ProblemEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemEvent proto.InternalMessageInfo

type isProblemEvent_Event interface {
	isProblemEvent_Event()
	MarshalTo([]byte) (int, error)
	Size() int
}

type ProblemEvent_Submission struct {
	Submission *ProblemSubmissionEvent `protobuf:"bytes,16,opt,name=submission,oneof"`
}

func (*ProblemEvent_Submission) isProblemEvent_Event() {}

func (m *ProblemEvent) GetEvent() isProblemEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *ProblemEvent) GetProblemId() string {
	if m != nil && m.ProblemId != nil {
		return *m.ProblemId
	}
	return ""
}

func (m *ProblemEvent) GetSubmission() *ProblemSubmissionEvent {
	if x, ok := m.GetEvent().(*ProblemEvent_Submission); ok {
		return x.Submission
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ProblemEvent) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ProblemEvent_OneofMarshaler, _ProblemEvent_OneofUnmarshaler, _ProblemEvent_OneofSizer, []interface{}{
		(*ProblemEvent_Submission)(nil),
	}
}

func _ProblemEvent_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ProblemEvent)
	// event
	switch x := m.Event.(type) {
	case *ProblemEvent_Submission:
		_ = b.EncodeVarint(16<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Submission); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ProblemEvent.Event has unexpected type %T", x)
	}
	return nil
}

func _ProblemEvent_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ProblemEvent)
	switch tag {
	case 16: // event.submission
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ProblemSubmissionEvent)
		err := b.DecodeMessage(msg)
		m.Event = &ProblemEvent_Submission{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ProblemEvent_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ProblemEvent)
	// event
	switch x := m.Event.(type) {
	case *ProblemEvent_Submission:
		s := proto.Size(x.Submission)
		n += 2 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ProblemSubmissionEvent struct {
	SubmissionId         *string  `protobuf:"bytes,1,opt,name=submission_id,json=submissionId" json:"submission_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProblemSubmissionEvent) Reset()         { *m = ProblemSubmissionEvent{} }
func (m *ProblemSubmissionEvent) String() string { return proto.CompactTextString(m) }
func (*ProblemSubmissionEvent) ProtoMessage()    {}
func (*ProblemSubmissionEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_syzoj_service_problem_kafka_e523ec2f3887079d, []int{1}
}
func (m *ProblemSubmissionEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProblemSubmissionEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProblemSubmissionEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ProblemSubmissionEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemSubmissionEvent.Merge(dst, src)
}
func (m *ProblemSubmissionEvent) XXX_Size() int {
	return m.Size()
}
func (m *ProblemSubmissionEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemSubmissionEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemSubmissionEvent proto.InternalMessageInfo

func (m *ProblemSubmissionEvent) GetSubmissionId() string {
	if m != nil && m.SubmissionId != nil {
		return *m.SubmissionId
	}
	return ""
}

func init() {
	proto.RegisterType((*ProblemEvent)(nil), "syzoj.service.problem.kafka.ProblemEvent")
	proto.RegisterType((*ProblemSubmissionEvent)(nil), "syzoj.service.problem.kafka.ProblemSubmissionEvent")
}
func (m *ProblemEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProblemEvent) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ProblemId != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintSyzojServiceProblemKafka(dAtA, i, uint64(len(*m.ProblemId)))
		i += copy(dAtA[i:], *m.ProblemId)
	}
	if m.Event != nil {
		nn1, err := m.Event.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *ProblemEvent_Submission) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Submission != nil {
		dAtA[i] = 0x82
		i++
		dAtA[i] = 0x1
		i++
		i = encodeVarintSyzojServiceProblemKafka(dAtA, i, uint64(m.Submission.Size()))
		n2, err := m.Submission.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}
func (m *ProblemSubmissionEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProblemSubmissionEvent) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SubmissionId != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintSyzojServiceProblemKafka(dAtA, i, uint64(len(*m.SubmissionId)))
		i += copy(dAtA[i:], *m.SubmissionId)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintSyzojServiceProblemKafka(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ProblemEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ProblemId != nil {
		l = len(*m.ProblemId)
		n += 1 + l + sovSyzojServiceProblemKafka(uint64(l))
	}
	if m.Event != nil {
		n += m.Event.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ProblemEvent_Submission) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Submission != nil {
		l = m.Submission.Size()
		n += 2 + l + sovSyzojServiceProblemKafka(uint64(l))
	}
	return n
}
func (m *ProblemSubmissionEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SubmissionId != nil {
		l = len(*m.SubmissionId)
		n += 1 + l + sovSyzojServiceProblemKafka(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovSyzojServiceProblemKafka(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozSyzojServiceProblemKafka(x uint64) (n int) {
	return sovSyzojServiceProblemKafka(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProblemEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSyzojServiceProblemKafka
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProblemEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProblemEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProblemId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyzojServiceProblemKafka
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSyzojServiceProblemKafka
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.ProblemId = &s
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Submission", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyzojServiceProblemKafka
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSyzojServiceProblemKafka
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &ProblemSubmissionEvent{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Event = &ProblemEvent_Submission{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSyzojServiceProblemKafka(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSyzojServiceProblemKafka
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ProblemSubmissionEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSyzojServiceProblemKafka
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProblemSubmissionEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProblemSubmissionEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubmissionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSyzojServiceProblemKafka
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSyzojServiceProblemKafka
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.SubmissionId = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSyzojServiceProblemKafka(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthSyzojServiceProblemKafka
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSyzojServiceProblemKafka(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSyzojServiceProblemKafka
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSyzojServiceProblemKafka
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSyzojServiceProblemKafka
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthSyzojServiceProblemKafka
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowSyzojServiceProblemKafka
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipSyzojServiceProblemKafka(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthSyzojServiceProblemKafka = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSyzojServiceProblemKafka   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("syzoj.service.problem.kafka.proto", fileDescriptor_syzoj_service_problem_kafka_e523ec2f3887079d)
}

var fileDescriptor_syzoj_service_problem_kafka_e523ec2f3887079d = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2c, 0xae, 0xac, 0xca,
	0xcf, 0xd2, 0x2b, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x4f, 0xca, 0x49,
	0xcd, 0xd5, 0xcb, 0x4e, 0x4c, 0xcb, 0x4e, 0x04, 0xf1, 0x4a, 0xf2, 0x85, 0xa4, 0xf1, 0x28, 0x51,
	0xea, 0x65, 0xe4, 0xe2, 0x09, 0x80, 0x88, 0xb8, 0x96, 0xa5, 0xe6, 0x95, 0x08, 0xc9, 0x72, 0x71,
	0x41, 0x55, 0xc4, 0x67, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x71, 0x42, 0x45, 0x3c,
	0x53, 0x84, 0x42, 0xb9, 0xb8, 0x8a, 0x4b, 0x93, 0x72, 0x33, 0x8b, 0x8b, 0x33, 0xf3, 0xf3, 0x24,
	0x04, 0x14, 0x18, 0x35, 0xb8, 0x8d, 0x8c, 0xf5, 0xf0, 0x39, 0x02, 0x6a, 0x7a, 0x30, 0x5c, 0x17,
	0xd8, 0x1e, 0x0f, 0x86, 0x20, 0x24, 0x83, 0x9c, 0xd8, 0xb9, 0x58, 0x53, 0x41, 0xc2, 0x4a, 0xb6,
	0x5c, 0x62, 0xd8, 0x35, 0x08, 0x29, 0x73, 0xf1, 0x22, 0x34, 0x20, 0xdc, 0xc6, 0x83, 0x10, 0xf4,
	0x4c, 0x71, 0x72, 0x38, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18,
	0xa3, 0x8c, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xc1, 0x4e, 0x84,
	0x90, 0xba, 0x79, 0xe9, 0xba, 0xe9, 0xf9, 0xfa, 0x50, 0xe7, 0xea, 0x43, 0x9d, 0xab, 0x0f, 0x76,
	0x2e, 0x20, 0x00, 0x00, 0xff, 0xff, 0x98, 0x69, 0x38, 0x87, 0x51, 0x01, 0x00, 0x00,
}
