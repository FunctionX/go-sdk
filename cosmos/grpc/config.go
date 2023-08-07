package grpc

import (
	"time"

	"github.com/pkg/errors"
)

type Config struct {
	ChainId       string        `yaml:"chain_id" mapstructure:"chain_id"`
	RpcUrl        string        `yaml:"rpc_url" mapstructure:"rpc_url"`
	Timeout       time.Duration `yaml:"timeout" mapstructure:"timeout"`
	AddressPrefix string        `yaml:"address_prefix" mapstructure:"address_prefix"`
}

func NewDefConfig() Config {
	return Config{
		Timeout: 10 * time.Second,
	}
}

func (c Config) Check() error {
	if c.ChainId == "" {
		return errors.New("chain id can not be empty")
	}
	if c.RpcUrl == "" {
		return errors.New("rpc url can not be empty")
	}
	if c.Timeout <= 0 || c.Timeout > 600*time.Second {
		return errors.New("timeout is invalid, should be in (0, 600)s")
	}
	if c.AddressPrefix == "" {
		return errors.New("address prefix can not be empty")
	}
	return nil
}
