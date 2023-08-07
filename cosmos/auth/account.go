package auth

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/functionx/go-sdk/cosmos/codec"
	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
	"github.com/functionx/go-sdk/cosmos/types"
)

var (
	_ AccountI                      = (*BaseAccount)(nil)
	_ types.UnpackInterfacesMessage = (*BaseAccount)(nil)
)

// GetAddress - Implements types.AccountI.
func (acc *BaseAccount) GetAddress() types.AccAddress {
	return types.MustAccAddressFromBech32(acc.Address)
}

// SetAddress - Implements types.AccountI.
func (acc *BaseAccount) SetAddress(addr types.AccAddress) error {
	if len(acc.Address) != 0 {
		return fmt.Errorf("cannot override BaseAccount address")
	}

	acc.Address = addr.String()
	return nil
}

// GetPubKey - Implements types.AccountI.
func (acc *BaseAccount) GetPubKey() (pk cosmoscrypto.PubKey) {
	if acc.PubKey == nil {
		return nil
	}
	content, ok := acc.PubKey.GetCachedValue().(cosmoscrypto.PubKey)
	if !ok {
		return nil
	}
	return content
}

// SetPubKey - Implements types.AccountI.
func (acc *BaseAccount) SetPubKey(pubKey cosmoscrypto.PubKey) error {
	if pubKey == nil {
		acc.PubKey = nil
		return nil
	}
	anyMsg, err := types.NewAnyWithValue(pubKey)
	if err != nil {
		return err
	}
	acc.PubKey = anyMsg
	return nil
}

// GetAccountNumber - Implements AccountI
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements AccountI
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements types.AccountI.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements types.AccountI.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

func (acc *BaseAccount) String() string {
	bz, err := codec.ProtoMarshalJson(acc, nil)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (acc *BaseAccount) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	if acc.PubKey == nil {
		return nil
	}
	var pubKey cosmoscrypto.PubKey
	return unpacker.UnpackAny(acc.PubKey, &pubKey)
}

// AccountI is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements AccountI.
type AccountI interface {
	proto.Message

	GetAddress() types.AccAddress
	SetAddress(types.AccAddress) error // errors if already set.

	GetPubKey() cosmoscrypto.PubKey // can return nil.
	SetPubKey(cosmoscrypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	String() string
}
