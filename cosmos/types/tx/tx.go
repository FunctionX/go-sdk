package tx

import (
	"fmt"

	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
	"github.com/functionx/go-sdk/cosmos/types"
)

const MaxGasWanted = uint64((1 << 63) - 1)

var _ types.Tx = &Tx{}

// GetMsgs implements the GetMsgs method on sdk.Tx.
func (t *Tx) GetMsgs() []types.Msg {
	if t == nil || t.Body == nil {
		return nil
	}

	msgs := t.Body.Messages
	res, err := GetMsgs(msgs, "transaction")
	if err != nil {
		panic(err)
	}
	return res
}

// ValidateBasic implements the ValidateBasic method on sdk.Tx.
func (t *Tx) ValidateBasic() error {
	if t == nil {
		return fmt.Errorf("bad Tx")
	}

	body := t.Body
	if body == nil {
		return fmt.Errorf("missing TxBody")
	}

	authInfo := t.AuthInfo
	if authInfo == nil {
		return fmt.Errorf("missing AuthInfo")
	}

	fee := authInfo.Fee
	if fee == nil {
		return fmt.Errorf("missing fee")
	}

	if fee.GasLimit > MaxGasWanted {
		return fmt.Errorf("invalid gas supplied; %d > %d", fee.GasLimit, MaxGasWanted)
	}

	if fee.Amount.IsAnyNil() {
		return fmt.Errorf("invalid fee provided: null")
	}

	if fee.Amount.IsAnyNegative() {
		return fmt.Errorf("invalid fee provided: %s", fee.Amount)
	}

	if fee.Payer != "" {
		_, err := types.AccAddressFromBech32(fee.Payer)
		if err != nil {
			return fmt.Errorf("invalid fee payer address: %w", err)
		}
	}

	sigs := t.Signatures

	if len(sigs) == 0 {
		return fmt.Errorf("no signers")
	}

	if len(sigs) != len(t.GetSigners()) {
		return fmt.Errorf("wrong number of signers; expected %d, got %d", len(t.GetSigners()), len(sigs))
	}

	return nil
}

// GetSigners retrieves all the signers of a tx.
// This includes all unique signers of the messages (in order),
// as well as the FeePayer (if specified and not already included).
func (t *Tx) GetSigners() []types.AccAddress {
	var signers []types.AccAddress
	seen := map[string]bool{}

	for _, msg := range t.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	// ensure any specified fee payer is included in the required signers (at the end)
	feePayer := t.AuthInfo.Fee.Payer
	if feePayer != "" && !seen[feePayer] {
		payerAddr := types.MustAccAddressFromBech32(feePayer)
		signers = append(signers, payerAddr)
		seen[feePayer] = true
	}

	return signers
}

func (t *Tx) GetGas() uint64 {
	return t.AuthInfo.Fee.GasLimit
}

func (t *Tx) GetFee() types.Coins {
	return t.AuthInfo.Fee.Amount
}

func (t *Tx) FeePayer() types.AccAddress {
	feePayer := t.AuthInfo.Fee.Payer
	if feePayer != "" {
		return types.MustAccAddressFromBech32(feePayer)
	}
	// use first signer as default if no payer specified
	return t.GetSigners()[0]
}

func (t *Tx) FeeGranter() types.AccAddress {
	feePayer := t.AuthInfo.Fee.Granter
	if feePayer != "" {
		return types.MustAccAddressFromBech32(feePayer)
	}
	return types.AccAddress{}
}

func (t *Tx) GetMemo() string {
	return t.Body.Memo
}

func GetMsgs(anys []*types.Any, name string) ([]types.Msg, error) {
	msgs := make([]types.Msg, len(anys))
	for i, a := range anys {
		cached := a.GetCachedValue()
		if cached == nil {
			return nil, fmt.Errorf("any cached value is nil, %s messages must be correctly packed any values", name)
		}
		msgs[i] = cached.(types.Msg)
	}
	return msgs, nil
}

func (m *SignerInfo) GetPubKey() cosmoscrypto.PubKey {
	return m.PublicKey.GetCachedValue().(cosmoscrypto.PubKey)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *TxBody) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, message := range m.Messages {
		var msg types.Msg
		if err := unpacker.UnpackAny(message, &msg); err != nil {
			return err
		}
	}
	return nil
}

// UnpackInterfaces implements the UnpackInterfaceMessages.UnpackInterfaces method
func (m *AuthInfo) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, signerInfo := range m.SignerInfos {
		if err := signerInfo.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *SignerInfo) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	return unpacker.UnpackAny(m.PublicKey, new(cosmoscrypto.PubKey))
}
