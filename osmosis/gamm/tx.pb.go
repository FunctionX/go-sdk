// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: osmosis/gamm/tx.proto

package gamm

import (
	fmt "fmt"
	github_com_functionx_go_sdk_cosmos_types "github.com/functionx/go-sdk/cosmos/types"
	types "github.com/functionx/go-sdk/cosmos/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ===================== MsgExitPool
type MsgExitPool struct {
	Sender        string                                       `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty" yaml:"sender"`
	PoolId        uint64                                       `protobuf:"varint,2,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty" yaml:"pool_id"`
	ShareInAmount github_com_functionx_go_sdk_cosmos_types.Int `protobuf:"bytes,3,opt,name=share_in_amount,json=shareInAmount,proto3,customtype=github.com/functionx/go-sdk/cosmos/types.Int" json:"share_in_amount" yaml:"share_in_amount"`
	TokenOutMins  []types.Coin                                 `protobuf:"bytes,4,rep,name=token_out_mins,json=tokenOutMins,proto3" json:"token_out_mins" yaml:"token_out_min_amounts"`
}

func (m *MsgExitPool) Reset()         { *m = MsgExitPool{} }
func (m *MsgExitPool) String() string { return proto.CompactTextString(m) }
func (*MsgExitPool) ProtoMessage()    {}
func (*MsgExitPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_e80d3998509184a8, []int{0}
}
func (m *MsgExitPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgExitPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgExitPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgExitPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgExitPool.Merge(m, src)
}
func (m *MsgExitPool) XXX_Size() int {
	return m.Size()
}
func (m *MsgExitPool) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgExitPool.DiscardUnknown(m)
}

var xxx_messageInfo_MsgExitPool proto.InternalMessageInfo

func (m *MsgExitPool) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *MsgExitPool) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *MsgExitPool) GetTokenOutMins() []types.Coin {
	if m != nil {
		return m.TokenOutMins
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgExitPool)(nil), "osmosis.gamm.v1beta1.MsgExitPool")
}

func init() { proto.RegisterFile("osmosis/gamm/tx.proto", fileDescriptor_e80d3998509184a8) }

var fileDescriptor_e80d3998509184a8 = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xcd, 0x4e, 0xea, 0x40,
	0x14, 0xc7, 0x5b, 0x20, 0xdc, 0xdc, 0x72, 0xe1, 0xc6, 0x06, 0x4d, 0x25, 0xa6, 0x25, 0x8d, 0x8b,
	0x1a, 0x75, 0x26, 0xe8, 0xce, 0x9d, 0x35, 0x2e, 0x58, 0xa0, 0xa6, 0x4b, 0x37, 0x4d, 0x3f, 0xc6,
	0x32, 0x81, 0xce, 0x21, 0xcc, 0x94, 0xc0, 0x5b, 0xf8, 0x50, 0x2e, 0x58, 0xb2, 0x34, 0x2e, 0x1a,
	0x03, 0x6f, 0xd0, 0x27, 0x30, 0xfd, 0x30, 0x91, 0x8d, 0xbb, 0x33, 0xff, 0xf3, 0x9b, 0xf3, 0x3f,
	0x1f, 0xca, 0x21, 0xf0, 0x18, 0x38, 0xe5, 0x38, 0xf2, 0xe2, 0x18, 0x8b, 0x25, 0x9a, 0xcd, 0x41,
	0x80, 0xda, 0xad, 0x64, 0x94, 0xcb, 0x68, 0x31, 0xf0, 0x89, 0xf0, 0x06, 0xbd, 0x6e, 0x04, 0x11,
	0x14, 0x00, 0xce, 0xa3, 0x92, 0xed, 0xe9, 0x41, 0x01, 0x63, 0xdf, 0xe3, 0x04, 0x57, 0x28, 0x0e,
	0x80, 0xb2, 0x32, 0x6f, 0xbe, 0xd5, 0x94, 0xd6, 0x88, 0x47, 0xf7, 0x4b, 0x2a, 0x9e, 0x00, 0xa6,
	0xea, 0x99, 0xd2, 0xe4, 0x84, 0x85, 0x64, 0xae, 0xc9, 0x7d, 0xd9, 0xfa, 0x6b, 0x1f, 0x64, 0xa9,
	0xd1, 0x5e, 0x79, 0xf1, 0xf4, 0xc6, 0x2c, 0x75, 0xd3, 0xa9, 0x00, 0xf5, 0x5c, 0xf9, 0x33, 0x03,
	0x98, 0xba, 0x34, 0xd4, 0x6a, 0x7d, 0xd9, 0x6a, 0xd8, 0x6a, 0x96, 0x1a, 0x9d, 0x92, 0xad, 0x12,
	0xa6, 0xd3, 0xcc, 0xa3, 0x61, 0xa8, 0x2e, 0x94, 0xff, 0x7c, 0xec, 0xcd, 0x89, 0x4b, 0x99, 0xeb,
	0xc5, 0x90, 0x30, 0xa1, 0xd5, 0x0b, 0x83, 0x87, 0x75, 0x6a, 0x48, 0x1f, 0xa9, 0x71, 0x11, 0x51,
	0x31, 0x4e, 0x7c, 0x14, 0x40, 0x8c, 0x5f, 0x12, 0x16, 0x08, 0x0a, 0x6c, 0x89, 0x23, 0xb8, 0xe4,
	0xe1, 0x04, 0x57, 0x43, 0x88, 0xd5, 0x8c, 0x70, 0x34, 0x64, 0x22, 0x4b, 0x8d, 0xa3, 0xaa, 0xa9,
	0xfd, 0xa2, 0xa6, 0xd3, 0x2e, 0x94, 0x21, 0xbb, 0x2d, 0xde, 0x2a, 0x51, 0x3a, 0x02, 0x26, 0x84,
	0xb9, 0x90, 0x08, 0x37, 0xa6, 0x8c, 0x6b, 0x8d, 0x7e, 0xdd, 0x6a, 0x5d, 0x1d, 0xa3, 0xb2, 0x26,
	0xca, 0x17, 0xf3, 0xbd, 0x43, 0x74, 0x07, 0x94, 0xd9, 0xa7, 0x79, 0x47, 0x59, 0x6a, 0x9c, 0x94,
	0x0e, 0x7b, 0xdf, 0x2b, 0x1b, 0x6e, 0x3a, 0xff, 0x0a, 0xfd, 0x31, 0x11, 0x23, 0xca, 0xb8, 0x6d,
	0xaf, 0xb7, 0xba, 0xbc, 0xd9, 0xea, 0xf2, 0xe7, 0x56, 0x97, 0x5f, 0x77, 0xba, 0xb4, 0xd9, 0xe9,
	0xd2, 0xfb, 0x4e, 0x97, 0x9e, 0xad, 0xdf, 0xe6, 0xfa, 0x79, 0x5f, 0xbf, 0x59, 0x5c, 0xe4, 0xfa,
	0x2b, 0x00, 0x00, 0xff, 0xff, 0x50, 0xf7, 0x3d, 0x5e, 0xf6, 0x01, 0x00, 0x00,
}

func (m *MsgExitPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgExitPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgExitPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TokenOutMins) > 0 {
		for iNdEx := len(m.TokenOutMins) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TokenOutMins[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	{
		size := m.ShareInAmount.Size()
		i -= size
		if _, err := m.ShareInAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.PoolId != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgExitPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.PoolId != 0 {
		n += 1 + sovTx(uint64(m.PoolId))
	}
	l = m.ShareInAmount.Size()
	n += 1 + l + sovTx(uint64(l))
	if len(m.TokenOutMins) > 0 {
		for _, e := range m.TokenOutMins {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgExitPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgExitPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgExitPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ShareInAmount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ShareInAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenOutMins", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenOutMins = append(m.TokenOutMins, types.Coin{})
			if err := m.TokenOutMins[len(m.TokenOutMins)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
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
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
