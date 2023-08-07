package codec

import (
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/encoding"

	"github.com/functionx/go-sdk/cosmos/types"
)

type (
	Codec interface {
		BinaryCodec
		JsonCodec
	}

	BinaryCodec interface {
		// Marshal returns binary encoding of v.
		Marshal(o ProtoMarshaler) ([]byte, error)
		// MustMarshal calls Marshal and panics if error is returned.
		MustMarshal(o ProtoMarshaler) []byte

		// Unmarshal parses the data encoded with Marshal method and stores the result
		// in the value pointed to by v.
		Unmarshal(bz []byte, ptr ProtoMarshaler) error
		// MustUnmarshal calls Unmarshal and panics if error is returned.
		MustUnmarshal(bz []byte, ptr ProtoMarshaler)

		// MarshalInterface is a helper method which will wrap `i` into `Any` for correct
		// binary interface (de)serialization.
		MarshalInterface(i proto.Message) ([]byte, error)
		// UnmarshalInterface is a helper method which will parse binary enoded data
		// into `Any` and unpack any into the `ptr`. It fails if the target interface type
		// is not registered in codec, or is not compatible with the serialized data
		UnmarshalInterface(bz []byte, ptr interface{}) error

		types.AnyUnpacker
	}

	JsonCodec interface {
		// MarshalJson returns Json encoding of v.
		MarshalJson(o proto.Message) ([]byte, error)
		// MustMarshalJson calls MarshalJson and panics if error is returned.
		MustMarshalJson(o proto.Message) []byte
		// MarshalInterfaceJson is a helper method which will wrap `i` into `Any` for correct
		// Json interface (de)serialization.
		MarshalInterfaceJson(i proto.Message) ([]byte, error)
		// UnmarshalInterfaceJson is a helper method which will parse Json enoded data
		// into `Any` and unpack any into the `ptr`. It fails if the target interface type
		// is not registered in codec, or is not compatible with the serialized data
		UnmarshalInterfaceJson(bz []byte, ptr interface{}) error

		// UnmarshalJson parses the data encoded with MarshalJson method and stores the result
		// in the value pointed to by v.
		UnmarshalJson(bz []byte, ptr proto.Message) error
		// MustUnmarshalJson calls Unmarshal and panics if error is returned.
		MustUnmarshalJson(bz []byte, ptr proto.Message)
	}

	// ProtoMarshaler defines an interface a type must implement to serialize itself
	// as a protocol buffer defined message.
	ProtoMarshaler interface {
		proto.Message // for Json serialization

		Marshal() ([]byte, error)
		MarshalTo(data []byte) (n int, err error)
		MarshalToSizedBuffer(dAtA []byte) (int, error)
		Size() int
		Unmarshal(data []byte) error
	}

	// GRPCCodecProvider is implemented by the Codec
	// implementations which return a gRPC encoding.Codec.
	// And it is used to decode requests and encode responses
	// passed through gRPC.
	GRPCCodecProvider interface {
		GRPCCodec() encoding.Codec
	}
)
