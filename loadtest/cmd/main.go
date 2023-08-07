package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/informalsystems/tm-load-test/pkg/loadtest"
	"github.com/spf13/cobra"

	loadtest2 "github.com/functionx/go-sdk/loadtest"
)

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErr(err)
	}
}

func newRootCmd() *cobra.Command {
	var cfg loadtest.Config
	rootCmd := &cobra.Command{
		Use:   "loadtest",
		Short: "Load test a Cosmos node",
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 2 {
				cfgData, err := os.ReadFile(args[1])
				if err != nil {
					return err
				}
				if err = json.Unmarshal(cfgData, &cfg); err != nil {
					return err
				}
			}
			keyDir := filepath.Join(os.ExpandEnv("$HOME"), "test_accounts")
			if len(args) == 3 {
				keyDir = args[2]
			}
			baseInfo, err := loadtest2.NewBaseInfo(args[0], keyDir)
			if err != nil {
				return err
			}
			fmt.Println("init accounts success", baseInfo.Accounts.Len())

			msgSendClientFactory := loadtest2.NewMsgSendClientFactory(baseInfo, baseInfo.GetDenom())
			if err = loadtest.RegisterClientFactory(msgSendClientFactory.Name(), msgSendClientFactory); err != nil {
				return err
			}

			msgOrderClientFactory := loadtest2.NewMsgOrderClientFactory(baseInfo)
			if err = loadtest.RegisterClientFactory(msgOrderClientFactory.Name(), msgOrderClientFactory); err != nil {
				return err
			}

			if err = cfg.Validate(); err != nil {
				return err
			}
			return loadtest.ExecuteStandalone(cfg)
		},
	}
	rootCmd.PersistentFlags().StringVar(&cfg.ClientFactory, "client-factory", "msg_send", "The identifier of the client factory to use for generating load testing transactions")
	rootCmd.PersistentFlags().IntVarP(&cfg.Connections, "connections", "c", 1, "The number of connections to open to each endpoint simultaneously")
	rootCmd.PersistentFlags().IntVarP(&cfg.Time, "time", "T", 60, "The duration (in seconds) for which to handle the load test")
	rootCmd.PersistentFlags().IntVarP(&cfg.SendPeriod, "send-period", "p", 1, "The period (in seconds) at which to send batches of transactions")
	rootCmd.PersistentFlags().IntVarP(&cfg.Rate, "rate", "r", 1000, "The number of transactions to generate each second on each connection, to each endpoint")
	rootCmd.PersistentFlags().IntVarP(&cfg.Size, "size", "s", 250, "The size of each transaction, in bytes - must be greater than 40")
	rootCmd.PersistentFlags().IntVarP(&cfg.Count, "count", "N", -1, "The maximum number of transactions to send - set to -1 to turn off this limit")
	rootCmd.PersistentFlags().StringVar(&cfg.BroadcastTxMethod, "broadcast-tx-method", "async", "The broadcast_tx method to use when submitting transactions - can be async, sync or commit")
	rootCmd.PersistentFlags().StringSliceVar(&cfg.Endpoints, "endpoints", []string{}, "A comma-separated list of URLs indicating Tendermint WebSockets RPC endpoints to which to connect")
	rootCmd.PersistentFlags().StringVar(&cfg.EndpointSelectMethod, "endpoint-select-method", loadtest.SelectSuppliedEndpoints, "The method by which to select endpoints")
	rootCmd.PersistentFlags().IntVar(&cfg.ExpectPeers, "expect-peers", 0, "The minimum number of peers to expect when crawling the P2P network from the specified endpoint(s) prior to waiting for workers to connect")
	rootCmd.PersistentFlags().IntVar(&cfg.MaxEndpoints, "max-endpoints", 0, "The maximum number of endpoints to use for testing, where 0 means unlimited")
	rootCmd.PersistentFlags().IntVar(&cfg.PeerConnectTimeout, "peer-connect-timeout", 600, "The number of seconds to wait for all required peers to connect if expect-peers > 0")
	rootCmd.PersistentFlags().IntVar(&cfg.MinConnectivity, "min-peer-connectivity", 0, "The minimum number of peers to which each peer must be connected before starting the load test")
	rootCmd.PersistentFlags().StringVar(&cfg.StatsOutputFile, "stats-output", "", "Where to store aggregate statistics (in CSV format) for the load test")
	return rootCmd
}
