package crypto

import (
	"crypto/sha256"

	"github.com/gogo/protobuf/proto"

	"github.com/functionx/go-sdk/tendermint/bytes"
)

type Algorithm string

// PubKey defines a public key and extends proto.Message.
type PubKey interface {
	proto.Message

	Address() Address
	Bytes() []byte
	VerifySignature(msg []byte, sig []byte) bool
	Equals(PubKey) bool
	Type() Algorithm
}

// PrivKey defines a private key and extends proto.Message.
type PrivKey interface {
	proto.Message
	Bytes() []byte
	Sign(msg []byte) ([]byte, error)
	PubKey() PubKey
	Equals(PrivKey) bool
	Type() Algorithm
}

type Address bytes.HexBytes

func (a Address) Bytes() []byte {
	return a
}

func AddressHash(bz []byte) Address {
	return Sum(bz)[:20]
}

func Sum(bz []byte) []byte {
	h := sha256.Sum256(bz)
	return h[:]
}
