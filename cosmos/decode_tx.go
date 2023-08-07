package cosmos

import (
	"github.com/functionx/go-sdk/cosmos/codec"
	"github.com/functionx/go-sdk/cosmos/types/tx"
)

func DecodeTx(cdc codec.BinaryCodec, txBytes []byte) (ttx *tx.Tx, err error) {
	// cdc, ok := marshal.(*codec.ProtoCodec)
	// if !ok {
	// 	return nil, errors.New("invalid proto codec")
	// }

	var raw tx.TxRaw
	// reject all unknown proto fields in the root TxRaw
	// if err = unknownproto.RejectUnknownFieldsStrict(txBytes, &raw, cdc.InterfaceRegistry()); err != nil {
	// 	return nil, err
	// }

	if err = cdc.Unmarshal(txBytes, &raw); err != nil {
		return nil, err
	}

	var body tx.TxBody
	if err = cdc.Unmarshal(raw.BodyBytes, &body); err != nil {
		return nil, err
	}

	var authInfo tx.AuthInfo
	// reject all unknown proto fields in AuthInfo
	// if err = unknownproto.RejectUnknownFieldsStrict(raw.AuthInfoBytes, &authInfo, cdc.InterfaceRegistry()); err != nil {
	// 	return nil, err
	// }
	if err = cdc.Unmarshal(raw.AuthInfoBytes, &authInfo); err != nil {
		return nil, err
	}
	return &tx.Tx{
		Body:       &body,
		AuthInfo:   &authInfo,
		Signatures: raw.Signatures,
	}, nil
}
