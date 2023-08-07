package loadtest

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/gogo/protobuf/proto"
	"github.com/informalsystems/tm-load-test/pkg/loadtest"

	"github.com/functionx/go-sdk/cosmos"
	cosmostypes "github.com/functionx/go-sdk/cosmos/types"
	"github.com/functionx/go-sdk/marginx"
)

var (
	_ loadtest.ClientFactory = (*MsgOrderClientFactory)(nil)
	_ loadtest.Client        = (*MsgOrderClientFactory)(nil)
)

type MsgOrderClientFactory struct {
	*BaseInfo
	pairId string
}

func NewMsgOrderClientFactory(baseInfo *BaseInfo) *MsgOrderClientFactory {
	baseInfo.GasLimit = 50000000
	baseInfo.GasPrice = cosmostypes.NewDecCoinFromDec(baseInfo.GetDenom(), sdkmath.LegacyNewDecWithPrec(1, 6))
	return &MsgOrderClientFactory{BaseInfo: baseInfo, pairId: "TSLA:USDT"}
}

func (c *MsgOrderClientFactory) Name() string {
	return "msg_order"
}

func (c *MsgOrderClientFactory) ValidateConfig(cfg loadtest.Config) error {
	return nil
}

func (c *MsgOrderClientFactory) NewClient(cfg loadtest.Config) (loadtest.Client, error) {
	c.Memo = fmt.Sprintf("msg_send_%d", cfg.Rate)
	return c, nil
}

func (c *MsgOrderClientFactory) GenerateTx() ([]byte, error) {
	account := c.Accounts.NextAccount()
	msgs := []cosmostypes.Msg{
		&marginx.MsgCreateOrder{
			Owner:        account.Address,
			PairId:       c.pairId,
			Direction:    marginx.BUY,
			Price:        sdkmath.LegacyNewDecWithPrec(600, 12),
			BaseQuantity: sdkmath.LegacyNewDecWithPrec(100, 15),
			Leverage:     1,
		},
		// &marginx.MsgCreateOrder{
		// 	Owner:        account.Address,
		// 	PairId:       c.pairId,
		// 	Direction:    marginx.BUY,
		// 	Price:        sdkmath.LegacyNewDecWithPrec(600, 12),
		// 	BaseQuantity: sdkmath.LegacyNewDecWithPrec(20, 15),
		// 	Leverage:     1,
		// },
		// &marginx.MsgCreateOrder{
		// 	Owner:        account.Address,
		// 	PairId:       c.pairId,
		// 	Direction:    marginx.BUY,
		// 	Price:        sdkmath.LegacyNewDecWithPrec(600, 12),
		// 	BaseQuantity: sdkmath.LegacyNewDecWithPrec(20, 15),
		// 	Leverage:     1,
		// },
		// &marginx.MsgCreateOrder{
		// 	Owner:        account.Address,
		// 	PairId:       c.pairId,
		// 	Direction:    marginx.BUY,
		// 	Price:        sdkmath.LegacyNewDecWithPrec(600, 12),
		// 	BaseQuantity: sdkmath.LegacyNewDecWithPrec(20, 15),
		// 	Leverage:     1,
		// },
		//
		&marginx.MsgCreateOrder{
			Owner:        account.Address,
			PairId:       c.pairId,
			Direction:    marginx.SELL,
			Price:        sdkmath.LegacyNewDecWithPrec(600, 12),
			BaseQuantity: sdkmath.LegacyNewDecWithPrec(100, 15),
			Leverage:     1,
		},
	}
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
