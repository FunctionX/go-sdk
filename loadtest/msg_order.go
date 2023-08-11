package loadtest

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/informalsystems/tm-load-test/pkg/loadtest"

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

func (c *MsgOrderClientFactory) ValidateConfig(_ loadtest.Config) error {
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
	return c.BuildTx(account, msgs)
}
