package evm

import (
	"math/big"
	"time"

	"github.com/catalogfi/blockchain"
	"github.com/ethereum/go-ethereum/common"
)

type Options struct {
	ChainID  *big.Int
	SwapAddr common.Address
	Timeout  time.Duration
	L2       bool
}

func NewOptions(chain blockchain.EvmChain, contract common.Address, timeout time.Duration) Options {
	return Options{
		ChainID:  chain.ChainID(),
		SwapAddr: contract,
		Timeout:  timeout,
		L2:       chain.L2(),
	}
}

func (opts Options) WithChainID(id *big.Int) Options {
	opts.ChainID = id
	return opts
}

func (opts Options) WithSwapAddr(swapAddr common.Address) Options {
	opts.SwapAddr = swapAddr
	return opts
}

func (opts Options) WithTimeout(timeout time.Duration) Options {
	opts.Timeout = timeout
	return opts
}

func (opts Options) WithL2(l2 bool) Options {
	opts.L2 = l2
	return opts
}
