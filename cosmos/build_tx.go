package cosmos

import (
	"github.com/cosmos/gogoproto/proto"

	"github.com/functionx/go-sdk/cosmos/auth"
	cosmoscrypto "github.com/functionx/go-sdk/cosmos/crypto"
	"github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/cosmos/types/tx"
	"github.com/functionx/go-sdk/cosmos/types/tx/signing"
)

const (
	DefGasLimit      int64   = 200000
	DefGasAdjustment float64 = 1.5
)

type simulationClient interface {
	EstimatingGas(raw *tx.TxRaw) (*types.GasInfo, error)
}

type buildTxClient interface {
	simulationClient
	GetChainId() (string, error)
	QueryAccount(address string) (auth.AccountI, error)
	GetGasPrices() (types.DecCoins, error)
	GetAddressPrefix() (string, error)
}

func BuildTx(cli buildTxClient, privKey cosmoscrypto.PrivKey, msgs []types.Msg) (*tx.TxRaw, error) {
	prefix, err := cli.GetAddressPrefix()
	if err != nil {
		return nil, err
	}
	account, err := cli.QueryAccount(types.NewAccAddress(privKey.PubKey().Address(), prefix).String())
	if err != nil {
		return nil, err
	}

	txBodyBytes, err := BuildTxBody(msgs, "", 0)
	if err != nil {
		return nil, err
	}

	pubAny, err := types.NewAnyWithValue(privKey.PubKey())
	if err != nil {
		return nil, err
	}
	gasPrices, err := cli.GetGasPrices()
	if err != nil {
		return nil, err
	}
	var gasPrice types.DecCoin
	if len(gasPrices) > 0 {
		gasPrice = gasPrices[0]
	}
	authInfo := NewAuthInfo(pubAny, account.GetSequence(), uint64(DefGasLimit), gasPrice)
	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}

	chainId, err := cli.GetChainId()
	if err != nil {
		return nil, err
	}
	signDoc := BuildSignDoc(chainId, txBodyBytes, txAuthInfoBytes, account.GetAccountNumber())

	gasInfo, err := EstimatingGas(cli, privKey, signDoc)
	if err != nil {
		return nil, err
	}
	authInfo.Fee.GasLimit = gasInfo.GasUsed * uint64(DefGasAdjustment*100) / 100
	authInfo.Fee.Amount = types.NewCoins(types.NewCoin(gasPrice.Denom, gasPrice.Amount.MulInt64(int64(authInfo.Fee.GasLimit)).Ceil().RoundInt()))

	signDoc.AuthInfoBytes, err = proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	sign, err := Sign(privKey, signDoc)
	if err != nil {
		return nil, err
	}
	return &tx.TxRaw{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}, nil
}

func BuildTxV1(chainId string, sequence, accountNumber uint64, privKey cosmoscrypto.PrivKey, msgs []types.Msg, gasPrice types.DecCoin, gasLimit int64, memo string, timeout uint64) (*tx.TxRaw, error) {
	txBodyBytes, err := BuildTxBody(msgs, memo, timeout)
	if err != nil {
		return nil, err
	}

	pubAny, err := types.NewAnyWithValue(privKey.PubKey())
	if err != nil {
		return nil, err
	}

	authInfo := NewAuthInfo(pubAny, sequence, uint64(gasLimit), gasPrice)
	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	signDoc := &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: accountNumber,
	}
	sign, err := Sign(privKey, signDoc)
	if err != nil {
		return nil, err
	}
	return &tx.TxRaw{
		BodyBytes:     signDoc.BodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}, nil
}

func BuildTxV2(chainId string, sequence, accountNumber uint64, privKey cosmoscrypto.PrivKey, txBytes []byte, gasPrice types.DecCoin, gasLimit uint64) (*tx.TxRaw, error) {
	pubAny, err := types.NewAnyWithValue(privKey.PubKey())
	if err != nil {
		return nil, err
	}
	authInfo := NewAuthInfoWithDecCoin(pubAny, sequence, gasLimit, gasPrice)
	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}
	signDoc := &tx.SignDoc{
		BodyBytes:     txBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: accountNumber,
	}
	sign, err := Sign(privKey, signDoc)
	if err != nil {
		return nil, err
	}
	return &tx.TxRaw{
		BodyBytes:     txBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	}, nil
}

func BuildTxBody(msgs []types.Msg, memo string, timeout uint64) ([]byte, error) {
	txBodyMessage := make([]*types.Any, 0)
	for i := 0; i < len(msgs); i++ {
		msgAnyValue, err := types.NewAnyWithValue(msgs[i])
		if err != nil {
			return nil, err
		}
		txBodyMessage = append(txBodyMessage, msgAnyValue)
	}

	txBody := &tx.TxBody{
		Messages:                    txBodyMessage,
		Memo:                        memo,
		TimeoutHeight:               timeout,
		ExtensionOptions:            nil,
		NonCriticalExtensionOptions: nil,
	}
	txBodyBytes, err := proto.Marshal(txBody)
	if err != nil {
		return nil, err
	}
	return txBodyBytes, nil
}

func NewAuthInfo(pubAny *types.Any, sequence, gasLimit uint64, gasPrice types.DecCoin) *tx.AuthInfo {
	authInfo := &tx.AuthInfo{
		SignerInfos: []*tx.SignerInfo{
			{
				PublicKey: pubAny,
				ModeInfo: &tx.ModeInfo{
					Sum: &tx.ModeInfo_Single_{
						Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
					},
				},
				Sequence: sequence,
			},
		},
		Fee: &tx.Fee{
			Amount:   types.NewCoins(types.NewCoin(gasPrice.Denom, gasPrice.Amount.MulInt64(int64(gasLimit)).Ceil().RoundInt())),
			GasLimit: gasLimit,
			Payer:    "",
			Granter:  "",
		},
	}
	return authInfo
}

func NewAuthInfoWithDecCoin(pubAny *types.Any, sequence, gasLimit uint64, gasPrice types.DecCoin) *tx.AuthInfo {
	authInfo := &tx.AuthInfo{
		SignerInfos: []*tx.SignerInfo{
			{
				PublicKey: pubAny,
				ModeInfo: &tx.ModeInfo{
					Sum: &tx.ModeInfo_Single_{
						Single: &tx.ModeInfo_Single{Mode: signing.SignMode_SIGN_MODE_DIRECT},
					},
				},
				Sequence: sequence,
			},
		},
		Fee: &tx.Fee{
			Amount:   types.NewCoins(types.NewCoin(gasPrice.Denom, gasPrice.Amount.MulInt64(int64(gasLimit)).TruncateInt())),
			GasLimit: gasLimit,
			Payer:    "",
			Granter:  "",
		},
	}
	return authInfo
}

func Sign(privKey cosmoscrypto.PrivKey, signDoc *tx.SignDoc) ([]byte, error) {
	signatures, err := proto.Marshal(signDoc)
	if err != nil {
		return nil, err
	}
	sign, err := privKey.Sign(signatures)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

func BuildSignDoc(chainId string, txBodyBytes, txAuthInfoBytes []byte, accountNumber uint64) *tx.SignDoc {
	signDoc := &tx.SignDoc{
		BodyBytes:     txBodyBytes,
		AuthInfoBytes: txAuthInfoBytes,
		ChainId:       chainId,
		AccountNumber: accountNumber,
	}
	return signDoc
}

func EstimatingGas(cli simulationClient, privKey cosmoscrypto.PrivKey, signDoc *tx.SignDoc) (*types.GasInfo, error) {
	sign, err := Sign(privKey, signDoc)
	if err != nil {
		return nil, err
	}
	gasInfo, err := cli.EstimatingGas(&tx.TxRaw{
		BodyBytes:     signDoc.BodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	})
	if err != nil {
		return nil, err
	}

	return gasInfo, nil
}

type TransactOpts struct {
	ChainId       string               `yaml:"chain_id" mapstructure:"chain_id"`
	PrivKey       cosmoscrypto.PrivKey `yaml:"-" mapstructure:"-"`
	Sequence      uint64               `yaml:"-" mapstructure:"-"`
	AccountNumber uint64               `yaml:"account_number" mapstructure:"account_number"`
	AddressPrefix string               `yaml:"address_prefix" mapstructure:"address_prefix"`

	GasPrice        types.DecCoin    `yaml:"gas_price" mapstructure:"gas_price"`
	GasLimit        uint64           `yaml:"gas_limit" mapstructure:"gas_limit"`
	GasAdjustment   float64          `yaml:"gas_adjustment" mapstructure:"gas_adjustment"`
	TxBroadcastMode tx.BroadcastMode `yaml:"tx_broadcast_mode" mapstructure:"tx_broadcast_mode"`

	TimeoutHeight uint64 `yaml:"timeout_height" mapstructure:"timeout_height"`
	Memo          string `yaml:"memo" mapstructure:"memo"`
}

func (opts TransactOpts) Sign(signDoc *tx.SignDoc) ([]byte, error) {
	sign, err := Sign(opts.PrivKey, signDoc)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

func (opts TransactOpts) BuildTx(msgs []types.Msg) (*tx.SignDoc, error) {
	txBodyBytes, err := BuildTxBody(msgs, opts.Memo, opts.TimeoutHeight)
	if err != nil {
		return nil, err
	}

	pubAny, err := types.NewAnyWithValue(opts.PrivKey.PubKey())
	if err != nil {
		return nil, err
	}

	authInfo := NewAuthInfo(pubAny, opts.Sequence, opts.GasLimit, opts.GasPrice)
	txAuthInfoBytes, err := proto.Marshal(authInfo)
	if err != nil {
		return nil, err
	}

	signDoc := BuildSignDoc(opts.ChainId, txBodyBytes, txAuthInfoBytes, opts.AccountNumber)
	return signDoc, nil
}

func (opts TransactOpts) EstimatingGas(cli simulationClient, signDoc *tx.SignDoc) (*types.GasInfo, error) {
	sign, err := opts.Sign(signDoc)
	if err != nil {
		return nil, err
	}
	gasInfo, err := cli.EstimatingGas(&tx.TxRaw{
		BodyBytes:     signDoc.BodyBytes,
		AuthInfoBytes: signDoc.AuthInfoBytes,
		Signatures:    [][]byte{sign},
	})
	if err != nil {
		return nil, err
	}

	return gasInfo, nil
}
