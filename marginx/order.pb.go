// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: marginx/v1/order.proto

package marginx

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// BOTH for query
type Direction int32

const (
	BOTH       Direction = 0
	BUY        Direction = 1
	SELL       Direction = 2
	MarketBuy  Direction = 3
	MarketSell Direction = 4
)

var Direction_name = map[int32]string{
	0: "BOTH",
	1: "BUY",
	2: "SELL",
	3: "MarketBuy",
	4: "MarketSell",
}

var Direction_value = map[string]int32{
	"BOTH":       0,
	"BUY":        1,
	"SELL":       2,
	"MarketBuy":  3,
	"MarketSell": 4,
}

func (x Direction) String() string {
	return proto.EnumName(Direction_name, int32(x))
}

func (Direction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2cb180b3e371b6e6, []int{0}
}

func init() {
	proto.RegisterEnum("fx.dex.v1.Direction", Direction_name, Direction_value)
}

func init() { proto.RegisterFile("marginx/v1/order.proto", fileDescriptor_2cb180b3e371b6e6) }

var fileDescriptor_2cb180b3e371b6e6 = []byte{
	// 244 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcb, 0x4d, 0x2c, 0x4a,
	0xcf, 0xcc, 0xab, 0xd0, 0x2f, 0x33, 0xd4, 0xcf, 0x2f, 0x4a, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0xe2, 0x4c, 0xab, 0xd0, 0x4b, 0x49, 0xad, 0xd0, 0x2b, 0x33, 0x94, 0x12, 0x49,
	0xcf, 0x4f, 0xcf, 0x07, 0x8b, 0xea, 0x83, 0x58, 0x10, 0x05, 0x5a, 0xbd, 0x8c, 0x5c, 0x9c, 0x2e,
	0x99, 0x45, 0xa9, 0xc9, 0x25, 0x99, 0xf9, 0x79, 0x42, 0x42, 0x5c, 0x2c, 0x4e, 0xfe, 0x21, 0x1e,
	0x02, 0x0c, 0x52, 0x1c, 0x5d, 0x73, 0x15, 0xc0, 0x6c, 0x21, 0x01, 0x2e, 0x66, 0xa7, 0xd0, 0x48,
	0x01, 0x46, 0x29, 0xf6, 0xae, 0xb9, 0x0a, 0x20, 0x26, 0x48, 0x55, 0xb0, 0xab, 0x8f, 0x8f, 0x00,
	0x13, 0x44, 0x15, 0x88, 0x2d, 0x24, 0xc3, 0xc5, 0xe9, 0x9b, 0x58, 0x94, 0x9d, 0x5a, 0xe2, 0x54,
	0x5a, 0x29, 0xc0, 0x2c, 0xc5, 0xdb, 0x35, 0x57, 0x01, 0x21, 0x20, 0x24, 0xc7, 0xc5, 0x05, 0xe1,
	0x04, 0xa7, 0xe6, 0xe4, 0x08, 0xb0, 0x48, 0xf1, 0x75, 0xcd, 0x55, 0x40, 0x12, 0x91, 0xe2, 0xe8,
	0x58, 0x2c, 0xc7, 0xb0, 0x62, 0x89, 0x1c, 0xa3, 0x93, 0xed, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e,
	0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72, 0x0c, 0x37,
	0x1e, 0xcb, 0x31, 0x44, 0x29, 0xa7, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea,
	0xa7, 0x95, 0xe6, 0x81, 0x1d, 0x5c, 0xa1, 0x9f, 0x9e, 0xaf, 0x5b, 0x9c, 0x92, 0xad, 0x0f, 0xf5,
	0x7e, 0x12, 0x1b, 0xd8, 0x57, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x47, 0xca, 0xc0, 0x7b,
	0x10, 0x01, 0x00, 0x00,
}
