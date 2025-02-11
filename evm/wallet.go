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
	ChainID   *big.Int
	SwapAddrs []common.Address
	Timeout   time.Duration
	L2        bool
}

func NewOptions(chain blockchain.EvmChain, contracts []common.Address, timeout time.Duration) Options {
	return Options{
		ChainID:   chain.ChainID(),
		SwapAddrs: contracts,
		Timeout:   timeout,
		L2:        chain.L2(),
	}
}

func (opts Options) WithChainID(id *big.Int) Options {
	opts.ChainID = id
	return opts
}

func (opts Options) WithSwapAddrs(swapAddrs []common.Address) Options {
	opts.SwapAddrs = swapAddrs
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
	addr    common.Address
	htlcs   map[common.Address]*gardenhtlc.GardenHTLC
	tokens  map[common.Address]*ierc20.IERC20

	mu           *sync.Mutex
	transactOpts *bind.TransactOpts
	approved     bool
}

func NewWallet(options Options, key *ecdsa.PrivateKey, client *ethclient.Client) (Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()
	callOpts := &bind.CallOpts{Context: ctx}
	addr := crypto.PubkeyToAddress(key.PublicKey)

	// Make sure the chain ID matches our expectation, so we know we are on the right chain.
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	if options.ChainID.Cmp(chainID) != 0 {
		return nil, fmt.Errorf("wrong chain ID, expect %v, got %v", options.ChainID, chainID)
	}

	// Initialise contract bindings.
	htlcs := make(map[common.Address]*gardenhtlc.GardenHTLC)
	tokens := make(map[common.Address]*ierc20.IERC20)
	for _, swapAddr := range options.SwapAddrs {
		var err error
		htlcs[swapAddr], err = gardenhtlc.NewGardenHTLC(swapAddr, client)
		if err != nil {
			return nil, err
		}
		tokenAddr, err := htlcs[swapAddr].Token(callOpts)
		if err != nil {
			return nil, err
		}
		tokens[swapAddr], err = ierc20.NewIERC20(tokenAddr, client)
		if err != nil {
			return nil, err
		}
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
		htlcs:        htlcs,
		tokens:       tokens,
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

func (wallet *wallet) Initiate(ctx context.Context, htlc Htlc) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	if err := wallet.allowanceCheck(); err != nil {
		return nil, err
	}
	contract, ok := wallet.htlcs[htlc.Contract]
	if !ok {
		return nil, fmt.Errorf("unknown contract %v", htlc.Contract.Hex())
	}
	f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Initiate(opts, htlc.Redeemer, htlc.Expiry, htlc.Amount, htlc.SecretHash)
	}
	return wallet.transact(ctx, f)
}

func (wallet *wallet) Redeem(ctx context.Context, htlc Htlc, secret []byte) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	contract, ok := wallet.htlcs[htlc.Contract]
	if !ok {
		return nil, fmt.Errorf("unknown contract %v", htlc.Contract.Hex())
	}
	f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Redeem(opts, htlc.ID, secret)
	}
	return wallet.transact(ctx, f)
}

func (wallet *wallet) Refund(ctx context.Context, htlc Htlc) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	contract, ok := wallet.htlcs[htlc.Contract]
	if !ok {
		return nil, fmt.Errorf("unknown contract %v", htlc.Contract.Hex())
	}
	f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return contract.Refund(opts, htlc.ID)
	}
	return wallet.transact(ctx, f)
}

func (wallet *wallet) InstantRefund(ctx context.Context, htlc Htlc, sig []byte) (*types.Transaction, error) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	contract, ok := wallet.htlcs[htlc.Contract]
	if !ok {
		return nil, fmt.Errorf("unknown contract %v", htlc.Contract.Hex())
	}

	// Generate signature if redeemerSig is nil
	if sig == nil {
		domain, err := contract.Eip712Domain(&bind.CallOpts{Context: ctx})
		if err != nil {
			return nil, fmt.Errorf("failed to get EIP-712 domain: %w", err)
		}
		td := apitypes.TypedData{
			Domain: apitypes.TypedDataDomain{
				Name:              domain.Name,
				Version:           domain.Version,
				ChainId:           math.NewHexOrDecimal256(wallet.options.ChainID.Int64()),
				VerifyingContract: htlc.Contract.String(),
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
		return contract.InstantRefund(opts, htlc.ID, sig)
	}
	return wallet.transact(ctx, f)
}

// allowanceCheck checks if the allowance for the swap contract is sufficient. It does a large approval when the
// allowance is low. The function should be called with the wallet lock held.
func (wallet *wallet) allowanceCheck() error {
	// Skip the allowance check if we have approved before
	if wallet.approved {
		return nil
	}

	for htlcAddr, token := range wallet.tokens {
		ctx, cancel := context.WithTimeout(context.Background(), wallet.options.Timeout)
		defer cancel()
		callOpts := &bind.CallOpts{Context: ctx}

		// Check if we have enough allowance for the swap contract
		allowance, err := token.Allowance(callOpts, wallet.addr, htlcAddr)
		if err != nil {
			return err
		}
		totalSupply, err := token.TotalSupply(callOpts)
		if err != nil {
			return err
		}

		// Do a large approval when the allowance is low, we should only need to do this once.
		if allowance.Cmp(totalSupply) == -1 {
			// If balance is 0, check later
			bal, err := wallet.Balance(ctx, true)
			if err != nil {
				return err
			}
			if bal.Uint64() == 0 {
				return nil
			}

			data := make([]byte, 32)
			for i := 0; i < 32; i++ {
				data[i] = 0xff
			}
			max := big.NewInt(0).SetBytes(data)
			f := func(opts *bind.TransactOpts) (*types.Transaction, error) {
				return token.Approve(opts, htlcAddr, max)
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
	}

	wallet.approved = true
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
