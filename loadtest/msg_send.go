package loadtest

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/gogo/protobuf/proto"
	"github.com/informalsystems/tm-load-test/pkg/loadtest"

	"github.com/functionx/go-sdk/cosmos"
	"github.com/functionx/go-sdk/cosmos/bank"
	cosmostypes "github.com/functionx/go-sdk/cosmos/types"
)

var (
	_ loadtest.ClientFactory = (*MsgSendClientFactory)(nil)
	_ loadtest.Client        = (*MsgSendClientFactory)(nil)
)

type MsgSendClientFactory struct {
	*BaseInfo
	denom string
}

func NewMsgSendClientFactory(baseInfo *BaseInfo, denom string) *MsgSendClientFactory {
	baseInfo.GasLimit = 100000
	baseInfo.GasPrice = cosmostypes.NewDecCoinFromDec(baseInfo.GetDenom(), sdkmath.LegacyNewDecWithPrec(1, 1))
	return &MsgSendClientFactory{BaseInfo: baseInfo, denom: denom}
}

func (c *MsgSendClientFactory) Name() string {
	return "msg_send"
}

func (c *MsgSendClientFactory) ValidateConfig(cfg loadtest.Config) error {
	return nil
}

func (c *MsgSendClientFactory) NewClient(cfg loadtest.Config) (loadtest.Client, error) {
	c.Memo = fmt.Sprintf("msg_send_%d", cfg.Rate)
	return c, nil
}

func (c *MsgSendClientFactory) GenerateTx() ([]byte, error) {
	account := c.Accounts.NextAccount()
	msgs := []cosmostypes.Msg{&bank.MsgSend{
		FromAddress: account.Address,
		ToAddress:   account.Address,
		Amount:      cosmostypes.NewCoins(cosmostypes.NewCoin(c.denom, sdkmath.NewInt(1))),
	}}
	if c.Accounts.IsFistAccount() {
		c.GasLimit--
	}
	txRaw, err := cosmos.BuildTxV1(
		c.ChainID, account.Sequence, account.AccountNumber, account.PrivateKey,
		msgs, c.GasPrice, c.GasLimit, c.Memo, 0,
	)
	if err != nil {
		return nil, err
	}
	// account.Sequence++
	txRawData, err := proto.Marshal(txRaw)
	if err != nil {
		return nil, err
	}
	// if account.AccountNumber%1000 == 0 {
	// 	fmt.Println("txHash: ", fmt.Sprintf("%X", cosmoscrypto.Sum(txRawData)))
	// }
	return txRawData, nil
}
