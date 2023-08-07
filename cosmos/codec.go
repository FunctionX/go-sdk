package cosmos

import (
	"github.com/functionx/go-sdk/cosmos/auth"
	"github.com/functionx/go-sdk/cosmos/bank"
	"github.com/functionx/go-sdk/cosmos/codec"
	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
	"github.com/functionx/go-sdk/cosmos/crypto/ethsecp256k1"
	"github.com/functionx/go-sdk/cosmos/crypto/secp256k1"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
)

func NewProtoCodec() *codec.ProtoCodec {
	registry := types.NewInterfaceRegistry()
	RegisterInterfaces(registry)
	return codec.NewProtoCodec(registry)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.base.v1beta1.Msg", (*types.Msg)(nil))
	registry.RegisterImplementations((*types.Msg)(nil), &bank.MsgSend{}, &bank.MsgMultiSend{})

	registry.RegisterInterface("cosmos.auth.v1beta1.AccountI", (*auth.AccountI)(nil))
	registry.RegisterImplementations((*auth.AccountI)(nil), &auth.BaseAccount{})

	registry.RegisterInterface("cosmos.tx.v1beta1.Tx", (*types.Tx)(nil))
	registry.RegisterImplementations((*types.Tx)(nil), &tx.Tx{})

	registry.RegisterInterface("cosmos.crypto.PubKey", (*cosmoscrypto.PubKey)(nil))
	registry.RegisterImplementations((*cosmoscrypto.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*cosmoscrypto.PubKey)(nil), &ethsecp256k1.PubKey{})

	registry.RegisterInterface("cosmos.crypto.PrivKey", (*cosmoscrypto.PrivKey)(nil))
	registry.RegisterImplementations((*cosmoscrypto.PrivKey)(nil), &secp256k1.PrivKey{})
	registry.RegisterImplementations((*cosmoscrypto.PrivKey)(nil), &ethsecp256k1.PrivKey{})
}
