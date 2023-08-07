package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/functionx/go-sdk/bech32"
)

// address is a common interface for different types of addresses used by the SDK
type address interface {
	Equals(address) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)
}

// Ensure that different address types implement the interface
var (
	_ address = AccAddress{}
)

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress struct {
	bytes  []byte
	prefix string
}

func NewAccAddress(bz []byte, prefix string) AccAddress {
	return AccAddress{bytes: bz, prefix: prefix}
}

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address string, prefix ...string) (addr AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return AccAddress{}, fmt.Errorf("empty address string is not allowed")
	}
	hrp, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return AccAddress{}, err
	}

	if len(prefix) > 0 && hrp != prefix[0] {
		return AccAddress{}, fmt.Errorf("invalid Bech32 prefix; expected %s, got %s", prefix, hrp)
	}
	return NewAccAddress(bz, hrp), nil
}

// MustAccAddressFromBech32 calls AccAddressFromBech32 and panics on error.
func MustAccAddressFromBech32(address string) AccAddress {
	addr, err := AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}
	return addr
}

// WithPrefix returns an AccAddress with the given prefix instead of the default one.
func (a AccAddress) WithPrefix(prefix string) AccAddress {
	a.prefix = prefix
	return a
}

// GetPrefix returns the prefix of the AccAddress.
func (a AccAddress) GetPrefix() string {
	return a.prefix
}

// Equals Returns boolean for whether two AccAddresses are Equal
func (a AccAddress) Equals(aa2 address) bool {
	if a.Empty() && aa2.Empty() {
		return true
	}
	return bytes.Equal(a.Bytes(), aa2.Bytes())
}

// Empty Returns boolean for whether an AccAddress is empty
func (a AccAddress) Empty() bool {
	return len(a.bytes) == 0
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (a AccAddress) Marshal() ([]byte, error) {
	return a.bytes, nil
}

// MarshalJSON marshals to JSON using Bech32.
func (a AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// Bytes returns the raw address bytes.
func (a AccAddress) Bytes() []byte {
	return a.bytes
}

// String implements the Stringer interface.
func (a AccAddress) String() string {
	if a.Empty() {
		return ""
	}
	if a.prefix == "" {
		panic("prefix should be set")
	}
	bech32Addr, err := bech32.ConvertAndEncode(a.prefix, a.bytes)
	if err != nil {
		panic(err)
	}
	return bech32Addr
}

// Format implements the fmt.Formatter interface.
func (a AccAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(a.String()))
	case 'p':
		_, _ = s.Write([]byte(fmt.Sprintf("%p", a)))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%X", a.bytes)))
	}
}
