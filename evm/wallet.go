package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/catalogfi/blockchain"
	"github.com/catalogfi/blockchain/evm/bindings/contracts/htlc/gardenhtlc"
	"github.com/catalogfi/blockchain/evm/bindings/openzeppelin/contracts/token/ERC20/ierc20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
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

type TransactFunc func(*bind.TransactOpts) (*types.Transaction, error)

type Wallet interface {

	// Address returns the address of the wallet
	Address() common.Address

	// Client returns the blockchain client.
	Client() *ethclient.Client

	// Balance returns the ETH balance of the wallet address
	Balance(ctx context.Context, pending bool) (*big.Int, error)

	// TokenBalance returns the token balance of the wallet address. Token is assumed an ERC-20 token and retrieved from
	// the HTLC contract.
	TokenBalance(ctx context.Context, pending bool) (*big.Int, error)

	// Initiate an atomic swap.
	Initiate(ctx context.Context, htlc Htlc) (*types.Transaction, error)

	// Redeem an atomic swap.
	Redeem(ctx context.Context, htlc Htlc, secret []byte) (*types.Transaction, error)

	// Refund an atomic swap.
	Refund(ctx context.Context, htlc Htlc) (*types.Transaction, error)

	// InstantRefund an atomic swap
	InstantRefund(ctx context.Context, htlc Htlc, sig []byte) (*types.Transaction, error)
}

type wallet struct {
	options Options
	key     *ecdsa.PrivateKey
	client  *ethclient.Client

	mu           *sync.Mutex
	addr         common.Address
	htlc         *gardenhtlc.GardenHTLC
	token        *ierc20.IERC20
	transactOpts *bind.TransactOpts
}

func NewWallet(options Options, key *ecdsa.PrivateKey, client *ethclient.Client) (Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()
	callOpts := &bind.CallOpts{Context: ctx}
	addr := crypto.PubkeyToAddress(key.PublicKey)

	// Initialise bindings.
	htlc, err := gardenhtlc.NewGardenHTLC(options.SwapAddr, client)
	if err != nil {
		return nil, err
	}
	tokenAddr, err := htlc.Token(callOpts)
	if err != nil {
		return nil, err
	}
	erc20, err := ierc20.NewIERC20(tokenAddr, client)
	if err != nil {
		return nil, err
	}

	// Make sure the chain ID matches our expectation, so we know we are on the right chain.
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	if options.ChainID.Cmp(chainID) != 0 {
		return nil, fmt.Errorf("wrong chain ID, expect %v, got %v", options.ChainID, chainID)
	}

	// Initialise the transactor
	nonce, err := client.PendingNonceAt(ctx, crypto.PubkeyToAddress(key.PublicKey))
	if err != nil {
		return nil, err
	}
	transactor, err := bind.NewKeyedTransactorWithChainID(key, options.ChainID)
	if err != nil {
		return nil, err
	}
	transactor.Nonce = big.NewInt(int64(nonce))

	wal := &wallet{
		options: options,
		key:     key,
		client:  client,

		mu:           new(sync.Mutex),
		addr:         addr,
		htlc:         htlc,
		token:        erc20,
		transactOpts: transactor,
	}

	// Check token allowance against the token contract
	if err := wal.allowanceCheck(); err != nil {
		return nil, err
	}

	return wal, nil
}

func (wallet *wallet) Address() common.Address {
	return wallet.addr
}

func (wallet *wallet) Client() *ethclient.Client {
	return wallet.client
}

func (wallet *wallet) Balance(ctx context.Context, pending bool) (*big.Int, error) {
	if pending {
		return wallet.client.PendingBalanceAt(ctx, wallet.addr)
	}
	return wallet.client.BalanceAt(ctx, wallet.addr, nil)
}

func (wallet *wallet) TokenBalance(ctx context.Context, pending bool) (*big.Int, error) {
	callOpts := &bind.CallOpts{
		Pending: pending,
		Context: ctx,
	}
	return wallet.token.BalanceOf(callOpts, wallet.addr)
}

func (wallet *wallet) Initiate(ctx context.Context, htlc Htlc) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return wallet.htlc.Initiate(opts, htlc.Redeemer, htlc.Expiry, htlc.Amount, htlc.SecretHash)
	}
	return wallet.transact(ctx, f)
}

func (wallet *wallet) Redeem(ctx context.Context, htlc Htlc, secret []byte) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return wallet.htlc.Redeem(opts, htlc.ID, secret)
	}
	return wallet.transact(ctx, f)
}

func (wallet *wallet) Refund(ctx context.Context, htlc Htlc) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return wallet.htlc.Refund(opts, htlc.ID)
	}
	return wallet.transact(ctx, f)
}

func (wallet *wallet) InstantRefund(ctx context.Context, htlc Htlc, sig []byte) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	// Generate signature if redeemerSig is nil
	if sig == nil {
		domain, err := wallet.htlc.Eip712Domain(&bind.CallOpts{Context: ctx})
		if err != nil {
			return nil, fmt.Errorf("failed to get EIP-712 domain: %w", err)
		}
		td := apitypes.TypedData{
			Domain: apitypes.TypedDataDomain{
				Name:              domain.Name,
				Version:           domain.Version,
				ChainId:           math.NewHexOrDecimal256(wallet.options.ChainID.Int64()),
				VerifyingContract: wallet.options.SwapAddr.String(),
			},
			Message: map[string]interface{}{
				"orderId": htlc.ID,
			},
			PrimaryType: "Refund",
			Types: apitypes.Types{
				"EIP712Domain": {
					{Name: "name", Type: "string"},
					{Name: "version", Type: "string"},
					{Name: "chainId", Type: "uint256"},
					{Name: "verifyingContract", Type: "address"},
				},
				"Refund": {
					{Name: "orderId", Type: "bytes32"},
				},
			},
		}
		domainSeparator, err := td.HashStruct("EIP712Domain", td.Domain.Map())
		if err != nil {
			return nil, fmt.Errorf("failed to hash domain: %w", err)
		}
		typedDataHash, err := td.HashStruct(td.PrimaryType, td.Message)
		if err != nil {
			return nil, fmt.Errorf("failed to hash message: %w", err)
		}
		rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
		digest := crypto.Keccak256Hash(rawData)
		signature, err := crypto.Sign(digest.Bytes(), wallet.key)
		if err != nil {
			return nil, fmt.Errorf("failed to sign digest: %w", err)
		}

		// Adjust V value (last byte) to conform to Ethereum's signature format
		signature[64] += 27
		sig = signature
	}

	f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return wallet.htlc.InstantRefund(opts, htlc.ID, sig)
	}
	return wallet.transact(ctx, f)
}

func (wallet *wallet) allowanceCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), wallet.options.Timeout*2)
	defer cancel()
	callOpts := &bind.CallOpts{Context: ctx}

	// Check we have enough allowance for the swap contract
	allowance, err := wallet.token.Allowance(callOpts, wallet.addr, wallet.options.SwapAddr)
	if err != nil {
		return err
	}
	totalSupply, err := wallet.token.TotalSupply(callOpts)
	if err != nil {
		return err
	}

	// Do a large approval when the allowance is low, we should only need to do this once.
	if allowance.Cmp(totalSupply) == -1 {
		data := make([]byte, 32)
		for i := 0; i < 32; i++ {
			data[i] = 0xff
		}
		max := big.NewInt(0).SetBytes(data)
		f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return wallet.token.Approve(opts, wallet.options.SwapAddr, max)
		}
		tx, err := wallet.transact(ctx, f)
		if err != nil {
			return err
		}

		// Wait for the tx to be mined and check receipt status
		receipt, err := bind.WaitMined(ctx, wallet.client, tx)
		if err != nil {
			return err
		}
		if receipt.Status == 0 {
			return fmt.Errorf("tx reverted, hash = %v", receipt.TxHash.Hex())
		}
	}
	return nil
}

// transact runs a transaction with the given function and automatically retries the transaction if nonce is incorrect.
// you need to hold the wallet lock when calling this function.
func (wallet *wallet) transact(ctx context.Context, f TransactFunc) (*types.Transaction, error) {
	wallet.transactOpts.Context = ctx
	for {
		tx, err := f(wallet.transactOpts)
		if err != nil {
			// If nonce is incorrect
			if strings.Contains(err.Error(), "nonce too low") || strings.Contains(err.Error(), "tx doesn't have the correct nonce") {
				nonce, err := wallet.client.PendingNonceAt(ctx, wallet.addr)
				if err != nil {
					return nil, err
				}
				wallet.transactOpts.Nonce = big.NewInt(int64(nonce))
				continue
			}

			// Return other errors immediately without retrying
			return nil, err
		}
		wallet.transactOpts.Nonce = big.NewInt(wallet.transactOpts.Nonce.Int64() + 1)
		return tx, nil
	}
}
