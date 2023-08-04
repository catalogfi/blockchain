package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type AccountAbstractionClient struct {
	network Network
	conn    *ethclient.Client
}

func NewAccountAbstractionClient(network Network, url string) (*AccountAbstractionClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &AccountAbstractionClient{
		network: network,
		conn:    client,
	}, nil
}

func (aac *AccountAbstractionClient) EstimateUserOperationGas(ctx context.Context, userOp userop.UserOperation) (*gas.GasEstimates, error) {
	userOp = aac.fillDummyData(userOp)
	entryPoint := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")

	params := &struct {
		Sender               string `json:"sender"`
		Nonce                string `json:"nonce"`
		InitCode             string `json:"initCode"`
		CallData             string `json:"callData"`
		CallGasLimit         string `json:"callGasLimit"`
		VerificationGasLimit string `json:"verificationGasLimit"`
		PreVerificationGas   string `json:"preVerificationGas"`
		MaxFeePerGas         string `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
		PaymasterAndData     string `json:"paymasterAndData"`
		Signature            string `json:"signature"`
	}{
		Sender:               userOp.Sender.String(),
		Nonce:                hexutil.EncodeBig(userOp.Nonce),
		InitCode:             hexutil.Encode(userOp.InitCode),
		CallData:             hexutil.Encode(userOp.CallData),
		CallGasLimit:         hexutil.EncodeBig(userOp.CallGasLimit),
		VerificationGasLimit: hexutil.EncodeBig(userOp.VerificationGasLimit),
		PreVerificationGas:   hexutil.EncodeBig(userOp.PreVerificationGas),
		MaxFeePerGas:         hexutil.EncodeBig(userOp.MaxFeePerGas),
		MaxPriorityFeePerGas: hexutil.EncodeBig(userOp.MaxPriorityFeePerGas),
		PaymasterAndData:     hexutil.Encode(userOp.PaymasterAndData),
		Signature:            hexutil.Encode(userOp.Signature),
	}
	var res gas.GasEstimates
	if err := aac.conn.Client().CallContext(ctx, &res, "eth_estimateUserOperationGas", params, entryPoint); err != nil {
		return nil, err
	}
	return &res, nil
}

func (aac *AccountAbstractionClient) SendUserOperation(ctx context.Context, userOp userop.UserOperation) error {
	return aac.conn.Client().CallContext(ctx, nil, "eth_sendUserOperation", userOp)
}

func (aac *AccountAbstractionClient) FeeEstimation(ctx context.Context, userOp userop.UserOperation) (userop.UserOperation, error) {
	// Update some data with dummy value
	userOp.CallGasLimit = big.NewInt(0)
	userOp.VerificationGasLimit = big.NewInt(0)
	userOp.PreVerificationGas = big.NewInt(0)
	userOp.MaxFeePerGas = big.NewInt(0)
	userOp.MaxPriorityFeePerGas = big.NewInt(0)
	userOp.Signature = make([]byte, crypto.SignatureLength*2)
	// TODO : Update the dummy data when we start using a paymaster
	userOp.PaymasterAndData = nil

	// Use the bundler API
	gasEstimate, err := aac.EstimateUserOperationGas(ctx, userOp)
	if err != nil {
		return userop.UserOperation{}, err
	}
	userOp.CallGasLimit = gasEstimate.CallGasLimit
	userOp.VerificationGasLimit = gasEstimate.VerificationGas
	userOp.PreVerificationGas = gasEstimate.PreVerificationGas

	// Estimate maxFee and maxGasTip with some overhead
	header, err := aac.conn.HeaderByNumber(ctx, nil)
	if err != nil {
		return userop.UserOperation{}, err
	}
	baseFee := big.NewInt(header.BaseFee.Int64() * 5 / 4)

	gasTip, err := aac.conn.SuggestGasTipCap(ctx)
	if err != nil {
		return userop.UserOperation{}, err
	}
	gasTip = big.NewInt(gasTip.Int64() * 11 / 10)
	maxFee := big.NewInt(0).Add(baseFee, gasTip)

	userOp.MaxFeePerGas = maxFee
	userOp.MaxPriorityFeePerGas = gasTip
	return userOp, nil
}

func (aac *AccountAbstractionClient) fillDummyData(userOp userop.UserOperation) userop.UserOperation {
	// Update some data with dummy value
	userOp.CallGasLimit = big.NewInt(0)
	userOp.VerificationGasLimit = big.NewInt(0)
	userOp.PreVerificationGas = big.NewInt(0)
	userOp.MaxFeePerGas = big.NewInt(0)
	userOp.MaxPriorityFeePerGas = big.NewInt(0)
	userOp.Signature = make([]byte, crypto.SignatureLength*2)
	return userOp
}
