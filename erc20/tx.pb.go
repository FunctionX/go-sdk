// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: erc20/v1/tx.proto

package erc20

import (
	context "context"
	fmt "fmt"
	github_com_functionx_go_sdk_cosmos_types "github.com/functionx/go-sdk/cosmos/types"
	_ "github.com/functionx/go-sdk/cosmos/types/msgservice"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// MsgConvertERC20 defines a Msg to convert an ERC20 token to a Cosmos SDK coin.
type MsgConvertERC20 struct {
	// ERC20 token contract address registered on erc20 bridge
	ContractAddress string `protobuf:"bytes,1,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
	// amount of ERC20 tokens to mint
	Amount github_com_functionx_go_sdk_cosmos_types.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/functionx/go-sdk/cosmos/types.Int" json:"amount"`
	// bech32 address to receive SDK coins.
	Receiver string `protobuf:"bytes,3,opt,name=receiver,proto3" json:"receiver,omitempty"`
	// sender hex address from the owner of the given ERC20 tokens
	Sender string `protobuf:"bytes,4,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (m *MsgConvertERC20) Reset()         { *m = MsgConvertERC20{} }
func (m *MsgConvertERC20) String() string { return proto.CompactTextString(m) }
func (*MsgConvertERC20) ProtoMessage()    {}
func (*MsgConvertERC20) Descriptor() ([]byte, []int) {
	return fileDescriptor_c099a28ada2f5869, []int{0}
}
func (m *MsgConvertERC20) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgConvertERC20) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgConvertERC20.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgConvertERC20) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgConvertERC20.Merge(m, src)
}
func (m *MsgConvertERC20) XXX_Size() int {
	return m.Size()
}
func (m *MsgConvertERC20) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgConvertERC20.DiscardUnknown(m)
}

var xxx_messageInfo_MsgConvertERC20 proto.InternalMessageInfo

func (m *MsgConvertERC20) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

func (m *MsgConvertERC20) GetReceiver() string {
	if m != nil {
		return m.Receiver
	}
	return ""
}

func (m *MsgConvertERC20) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

// MsgConvertERC20Response returns no fields
type MsgConvertERC20Response struct {
}

func (m *MsgConvertERC20Response) Reset()         { *m = MsgConvertERC20Response{} }
func (m *MsgConvertERC20Response) String() string { return proto.CompactTextString(m) }
func (*MsgConvertERC20Response) ProtoMessage()    {}
func (*MsgConvertERC20Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_c099a28ada2f5869, []int{1}
}
func (m *MsgConvertERC20Response) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgConvertERC20Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgConvertERC20Response.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgConvertERC20Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgConvertERC20Response.Merge(m, src)
}
func (m *MsgConvertERC20Response) XXX_Size() int {
	return m.Size()
}
func (m *MsgConvertERC20Response) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgConvertERC20Response.DiscardUnknown(m)
}

var xxx_messageInfo_MsgConvertERC20Response proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgConvertERC20)(nil), "erc20.v1.MsgConvertERC20")
	proto.RegisterType((*MsgConvertERC20Response)(nil), "erc20.v1.MsgConvertERC20Response")
}

func init() { proto.RegisterFile("erc20/v1/tx.proto", fileDescriptor_c099a28ada2f5869) }

var fileDescriptor_c099a28ada2f5869 = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xb1, 0x4e, 0x02, 0x41,
	0x10, 0x86, 0xef, 0xc4, 0x10, 0xdc, 0x98, 0xa0, 0x17, 0x23, 0xc7, 0x15, 0x8b, 0x50, 0x69, 0xa2,
	0xbb, 0x80, 0x76, 0x56, 0x42, 0x2c, 0x4c, 0xa0, 0x39, 0x3b, 0x1b, 0x03, 0x7b, 0xeb, 0x4a, 0xcc,
	0xed, 0x90, 0x9d, 0xe5, 0x82, 0x6f, 0xe1, 0x13, 0x59, 0x53, 0x52, 0x1a, 0x0b, 0x62, 0xe0, 0x45,
	0x0c, 0x0b, 0x98, 0x48, 0xa2, 0xdd, 0xcc, 0x37, 0xff, 0xfe, 0x3b, 0xf3, 0x93, 0x43, 0x69, 0x44,
	0xb3, 0xce, 0xb3, 0x06, 0xb7, 0x63, 0x36, 0x34, 0x60, 0x21, 0x28, 0x38, 0xc4, 0xb2, 0x46, 0x54,
	0x12, 0x80, 0x29, 0x20, 0x4f, 0x51, 0x2d, 0x15, 0x29, 0xaa, 0x95, 0x24, 0x3a, 0x52, 0xa0, 0xc0,
	0x95, 0x7c, 0x59, 0xad, 0x68, 0xed, 0xdd, 0x27, 0xc5, 0x2e, 0xaa, 0x36, 0xe8, 0x4c, 0x1a, 0x7b,
	0x1b, 0xb7, 0x9b, 0xf5, 0xe0, 0x8c, 0x1c, 0x08, 0xd0, 0xd6, 0xf4, 0x84, 0x7d, 0xec, 0x25, 0x89,
	0x91, 0x88, 0xa1, 0x7f, 0xe2, 0x9f, 0xee, 0xc5, 0xc5, 0x0d, 0xbf, 0x59, 0xe1, 0xa0, 0x43, 0xf2,
	0xbd, 0x14, 0x46, 0xda, 0x86, 0x3b, 0x4b, 0x41, 0xeb, 0x6a, 0x32, 0xab, 0x78, 0x9f, 0xb3, 0xca,
	0xb9, 0x1a, 0xd8, 0xe7, 0x51, 0x9f, 0x09, 0x48, 0xf9, 0xd3, 0x48, 0x0b, 0x3b, 0x00, 0x3d, 0xe6,
	0x0a, 0x2e, 0x30, 0x79, 0xe1, 0xeb, 0x0d, 0xed, 0xeb, 0x50, 0x22, 0xbb, 0xd3, 0x36, 0x5e, 0x7b,
	0x04, 0x11, 0x29, 0x18, 0x29, 0xe4, 0x20, 0x93, 0x26, 0xcc, 0xb9, 0x0f, 0x7f, 0xfa, 0xe0, 0x98,
	0xe4, 0x51, 0xea, 0x44, 0x9a, 0x70, 0xd7, 0x4d, 0xd6, 0x5d, 0xad, 0x4c, 0x4a, 0x5b, 0xfb, 0xc7,
	0x12, 0x87, 0xa0, 0x51, 0x36, 0xef, 0x49, 0xae, 0x8b, 0x2a, 0xe8, 0x90, 0xfd, 0x5f, 0xe7, 0x95,
	0xd9, 0x26, 0x2c, 0xb6, 0xf5, 0x32, 0xaa, 0xfe, 0x39, 0xda, 0x98, 0xb6, 0xae, 0x27, 0x73, 0xea,
	0x4f, 0xe7, 0xd4, 0xff, 0x9a, 0x53, 0xff, 0x6d, 0x41, 0xbd, 0xe9, 0x82, 0x7a, 0x1f, 0x0b, 0xea,
	0x3d, 0x54, 0xff, 0xbb, 0xd9, 0xf9, 0xf6, 0xf3, 0x2e, 0xf4, 0xcb, 0xef, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xcc, 0x8e, 0xd4, 0x86, 0xc2, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// ConvertERC20 mints a Cosmos coin representation of the ERC20 token contract
	// that is registered on the token mapping.
	ConvertERC20(ctx context.Context, in *MsgConvertERC20, opts ...grpc.CallOption) (*MsgConvertERC20Response, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) ConvertERC20(ctx context.Context, in *MsgConvertERC20, opts ...grpc.CallOption) (*MsgConvertERC20Response, error) {
	out := new(MsgConvertERC20Response)
	err := c.cc.Invoke(ctx, "/erc20.v1.Msg/ConvertERC20", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// ConvertERC20 mints a Cosmos coin representation of the ERC20 token contract
	// that is registered on the token mapping.
	ConvertERC20(context.Context, *MsgConvertERC20) (*MsgConvertERC20Response, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) ConvertERC20(ctx context.Context, req *MsgConvertERC20) (*MsgConvertERC20Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConvertERC20 not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_ConvertERC20_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgConvertERC20)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ConvertERC20(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/erc20.v1.Msg/ConvertERC20",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ConvertERC20(ctx, req.(*MsgConvertERC20))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "erc20.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConvertERC20",
			Handler:    _Msg_ConvertERC20_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "erc20/v1/tx.proto",
}

func (m *MsgConvertERC20) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgConvertERC20) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgConvertERC20) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Receiver) > 0 {
		i -= len(m.Receiver)
		copy(dAtA[i:], m.Receiver)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Receiver)))
		i--
		dAtA[i] = 0x1a
	}
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.ContractAddress) > 0 {
		i -= len(m.ContractAddress)
		copy(dAtA[i:], m.ContractAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ContractAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgConvertERC20Response) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgConvertERC20Response) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgConvertERC20Response) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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
func (m *MsgConvertERC20) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ContractAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovTx(uint64(l))
	l = len(m.Receiver)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgConvertERC20Response) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgConvertERC20) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgConvertERC20: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgConvertERC20: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ContractAddress", wireType)
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
			m.ContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
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
			m.Receiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
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
func (m *MsgConvertERC20Response) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgConvertERC20Response: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgConvertERC20Response: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
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
