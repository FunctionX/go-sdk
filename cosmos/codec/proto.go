package codec

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	gogoproto "github.com/gogo/protobuf/proto"
	legacyproto "github.com/golang/protobuf/proto" //nolint:staticcheck
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"

	"github.com/functionx/go-sdk/cosmos/types"
)

// ProtoCodecMarshaler defines an interface for codecs that utilize Protobuf for both
// binary and Json encoding.
type ProtoCodecMarshaler interface {
	Codec
	InterfaceRegistry() types.InterfaceRegistry
}

// ProtoCodec defines a codec that utilizes Protobuf for both binary and Json
// encoding.
type ProtoCodec struct {
	interfaceRegistry types.InterfaceRegistry
}

var (
	_ Codec               = &ProtoCodec{}
	_ ProtoCodecMarshaler = &ProtoCodec{}
)

// NewProtoCodec returns a reference to a new ProtoCodec
func NewProtoCodec(interfaceRegistry types.InterfaceRegistry) *ProtoCodec {
	return &ProtoCodec{interfaceRegistry: interfaceRegistry}
}

// Marshal implements BinaryMarshaler.Marshal method.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.MarshalInterface
func (pc *ProtoCodec) Marshal(o ProtoMarshaler) ([]byte, error) {
	// Size() check can catch the typed nil value.
	if o == nil || o.Size() == 0 {
		// return empty bytes instead of nil, because nil has special meaning in places like store.Set
		return []byte{}, nil
	}
	return o.Marshal()
}

// MustMarshal implements BinaryMarshaler.MustMarshal method.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.MarshalInterface
func (pc *ProtoCodec) MustMarshal(o ProtoMarshaler) []byte {
	bz, err := pc.Marshal(o)
	if err != nil {
		panic(err)
	}

	return bz
}

// MarshalLengthPrefixed implements BinaryMarshaler.MarshalLengthPrefixed method.
func (pc *ProtoCodec) MarshalLengthPrefixed(o ProtoMarshaler) ([]byte, error) {
	bz, err := pc.Marshal(o)
	if err != nil {
		return nil, err
	}

	var sizeBuf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(sizeBuf[:], uint64(o.Size()))
	return append(sizeBuf[:n], bz...), nil
}

// MustMarshalLengthPrefixed implements BinaryMarshaler.MustMarshalLengthPrefixed method.
func (pc *ProtoCodec) MustMarshalLengthPrefixed(o ProtoMarshaler) []byte {
	bz, err := pc.MarshalLengthPrefixed(o)
	if err != nil {
		panic(err)
	}

	return bz
}

// Unmarshal implements BinaryMarshaler.Unmarshal method.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.UnmarshalInterface
func (pc *ProtoCodec) Unmarshal(bz []byte, ptr ProtoMarshaler) error {
	err := ptr.Unmarshal(bz)
	if err != nil {
		return err
	}
	err = types.UnpackInterfaces(ptr, pc.interfaceRegistry)
	if err != nil {
		return err
	}
	return nil
}

// MustUnmarshal implements BinaryMarshaler.MustUnmarshal method.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.UnmarshalInterface
func (pc *ProtoCodec) MustUnmarshal(bz []byte, ptr ProtoMarshaler) {
	if err := pc.Unmarshal(bz, ptr); err != nil {
		panic(err)
	}
}

// UnmarshalLengthPrefixed implements BinaryMarshaler.UnmarshalLengthPrefixed method.
func (pc *ProtoCodec) UnmarshalLengthPrefixed(bz []byte, ptr ProtoMarshaler) error {
	size, n := binary.Uvarint(bz)
	if n < 0 {
		return fmt.Errorf("invalid number of bytes read from length-prefixed encoding: %d", n)
	}

	if size > uint64(len(bz)-n) {
		return fmt.Errorf("not enough bytes to read; want: %v, got: %v", size, len(bz)-n)
	} else if size < uint64(len(bz)-n) {
		return fmt.Errorf("too many bytes to read; want: %v, got: %v", size, len(bz)-n)
	}

	bz = bz[n:]
	return pc.Unmarshal(bz, ptr)
}

// MustUnmarshalLengthPrefixed implements BinaryMarshaler.MustUnmarshalLengthPrefixed method.
func (pc *ProtoCodec) MustUnmarshalLengthPrefixed(bz []byte, ptr ProtoMarshaler) {
	if err := pc.UnmarshalLengthPrefixed(bz, ptr); err != nil {
		panic(err)
	}
}

// MarshalJson implements JsonCodec.MarshalJson method,
// it marshals to Json using proto codec.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.MarshalInterfaceJson
func (pc *ProtoCodec) MarshalJson(o gogoproto.Message) ([]byte, error) {
	m, ok := o.(ProtoMarshaler)
	if !ok {
		return nil, fmt.Errorf("cannot protobuf Json encode unsupported type: %T", o)
	}

	return ProtoMarshalJson(m, pc.interfaceRegistry)
}

// MustMarshalJson implements JsonCodec.MustMarshalJson method,
// it executes MarshalJson except it panics upon failure.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.MarshalInterfaceJson
func (pc *ProtoCodec) MustMarshalJson(o gogoproto.Message) []byte {
	bz, err := pc.MarshalJson(o)
	if err != nil {
		panic(err)
	}

	return bz
}

// UnmarshalJson implements JsonCodec.UnmarshalJson method,
// it unmarshals from Json using proto codec.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.UnmarshalInterfaceJson
func (pc *ProtoCodec) UnmarshalJson(bz []byte, ptr gogoproto.Message) error {
	m, ok := ptr.(ProtoMarshaler)
	if !ok {
		return fmt.Errorf("cannot protobuf Json decode unsupported type: %T", ptr)
	}

	unmarshaler := jsonpb.Unmarshaler{AnyResolver: pc.interfaceRegistry}
	err := unmarshaler.Unmarshal(strings.NewReader(string(bz)), m)
	if err != nil {
		return err
	}

	return types.UnpackInterfaces(ptr, pc.interfaceRegistry)
}

// MustUnmarshalJson implements JsonCodec.MustUnmarshalJson method,
// it executes UnmarshalJson except it panics upon failure.
// NOTE: this function must be used with a concrete type which
// implements proto.Message. For interface please use the codec.UnmarshalInterfaceJson
func (pc *ProtoCodec) MustUnmarshalJson(bz []byte, ptr gogoproto.Message) {
	if err := pc.UnmarshalJson(bz, ptr); err != nil {
		panic(err)
	}
}

// MarshalInterface is a convenience function for proto marshalling interfaces. It packs
// the provided value, which must be an interface, in an Any and then marshals it to bytes.
// NOTE: to marshal a concrete type, you should use Marshal instead
func (pc *ProtoCodec) MarshalInterface(i gogoproto.Message) ([]byte, error) {
	if i == nil {
		return nil, fmt.Errorf("can't marshal <nil> value")
	}
	anyType, err := types.NewAnyWithValue(i)
	if err != nil {
		return nil, err
	}

	return pc.Marshal(anyType)
}

// UnmarshalInterface is a convenience function for proto unmarshaling interfaces. It
// unmarshals an Any from bz bytes and then unpacks it to the `ptr`, which must
// be a pointer to a non empty interface with registered implementations.
// NOTE: to unmarshal a concrete type, you should use Unmarshal instead
//
// Example:
//
//	var x MyInterface
//	err := cdc.UnmarshalInterface(bz, &x)
func (pc *ProtoCodec) UnmarshalInterface(bz []byte, ptr interface{}) error {
	anyType := &types.Any{}
	err := pc.Unmarshal(bz, anyType)
	if err != nil {
		return err
	}

	return pc.UnpackAny(anyType, ptr)
}

// MarshalInterfaceJson is a convenience function for proto marshalling interfaces. It
// packs the provided value in an Any and then marshals it to bytes.
// NOTE: to marshal a concrete type, you should use MarshalJson instead
func (pc *ProtoCodec) MarshalInterfaceJson(x gogoproto.Message) ([]byte, error) {
	anyType, err := types.NewAnyWithValue(x)
	if err != nil {
		return nil, err
	}
	return pc.MarshalJson(anyType)
}

// UnmarshalInterfaceJson is a convenience function for proto unmarshaling interfaces.
// It unmarshals an Any from bz bytes and then unpacks it to the `iface`, which must
// be a pointer to a non empty interface, implementing proto.Message with registered implementations.
// NOTE: to unmarshal a concrete type, you should use UnmarshalJson instead
//
// Example:
//
//	var x MyInterface  // must implement proto.Message
//	err := cdc.UnmarshalInterfaceJson(&x, bz)
func (pc *ProtoCodec) UnmarshalInterfaceJson(bz []byte, iface interface{}) error {
	anyType := &types.Any{}
	err := pc.UnmarshalJson(bz, anyType)
	if err != nil {
		return err
	}
	return pc.UnpackAny(anyType, iface)
}

// UnpackAny implements AnyUnpacker.UnpackAny method,
// it unpacks the value in any to the interface pointer passed in as
// iface.
func (pc *ProtoCodec) UnpackAny(any *types.Any, iface interface{}) error {
	return pc.interfaceRegistry.UnpackAny(any, iface)
}

// InterfaceRegistry returns InterfaceRegistry
func (pc *ProtoCodec) InterfaceRegistry() types.InterfaceRegistry {
	return pc.interfaceRegistry
}

// GRPCCodec returns the gRPC Codec for this specific ProtoCodec
func (pc *ProtoCodec) GRPCCodec() encoding.Codec {
	return &grpcProtoCodec{cdc: pc}
}

// grpcProtoCodec is the implementation of the gRPC proto codec.
type grpcProtoCodec struct {
	cdc *ProtoCodec
}

func (g grpcProtoCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case proto.Message:
		return proto.Marshal(m)
	case ProtoMarshaler:
		return g.cdc.Marshal(m)
	case legacyproto.Message:
		return legacyproto.Marshal(m)
	default:
		return nil, fmt.Errorf("codec: unknown proto type: cannot marshal type %T", v)
	}
}

func (g grpcProtoCodec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case proto.Message:
		return proto.Unmarshal(data, m)
	case ProtoMarshaler:
		return g.cdc.Unmarshal(data, m)
	case legacyproto.Message:
		return legacyproto.Unmarshal(data, m)
	default:
		return fmt.Errorf("codec: unknown proto type: cannot unmarshal type %T", v)
	}
}

func (g grpcProtoCodec) Name() string {
	return "cosmos-sdk-grpc-codec"
}
