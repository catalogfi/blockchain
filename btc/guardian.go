package btc

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"go.uber.org/zap"
)

func NewGaurdianScript(guardian, user *btcec.PublicKey) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_IF).
		AddInt64(2).
		AddData(guardian.SerializeCompressed()).
		AddData(user.SerializeCompressed()).
		AddInt64(2).
		AddOp(txscript.OP_CHECKMULTISIG).
		AddOp(txscript.OP_ELSE).
		AddInt64(144 * 180).
		AddOp(txscript.OP_CHECKLOCKTIMEVERIFY).
		AddOp(txscript.OP_DROP).
		AddData(user.SerializeCompressed()).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_ENDIF).
		Script()
}

type GuardianClient struct {
	PublicKey  *btcec.PublicKey
	Network    *chaincfg.Params
	Address    btcutil.Address
	URL        string
	Client     *http.Client
	Script     []byte
	Account    Account
	Cache      map[string]*wire.MsgTx
	spendUTXOs []wire.OutPoint
	LatestHash chainhash.Hash
}

type Account struct {
	Address           string `json:"address"`
	GuardianPublicKey string `json:"guardian_public_key"`
	UserPublicKey     string `json:"user_public_key"`
}

type SignTxResponse struct {
	TxHex string `json:"tx_hex"`
}

type UpdateTxRequest struct {
	Address  string   `json:"address"`
	TxHex    string   `json:"tx_hex"`
	Values   []int64  `json:"values"`
	MergeTxs []string `json:"merge_txs"`
}

func NewGuardianClient(url string, network *chaincfg.Params, publicKey *btcec.PublicKey) (*GuardianClient, error) {
	client := &http.Client{}
	publicKeyHex := hex.EncodeToString(publicKey.SerializeCompressed())

	createAccountRequest := map[string]string{
		"public_key": publicKeyHex,
	}
	jsonData, err := json.Marshal(createAccountRequest)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(fmt.Sprintf("%s/accounts", url), "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create account: %s", resp.Status)
	}

	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, err
	}

	addr, err := btcutil.DecodeAddress(account.Address, network)
	if err != nil {
		return nil, err
	}

	gClient := &GuardianClient{
		PublicKey: publicKey,
		Address:   addr,
		Network:   network,
		URL:       url,
		Client:    client,
		Account:   account,
		Cache:     map[string]*wire.MsgTx{},
	}

	guardianPubKey, err := gClient.GetGuardianPubkey()
	if err != nil {
		return nil, err
	}

	script, err := NewGaurdianScript(guardianPubKey, publicKey)
	if err != nil {
		return nil, err
	}
	gClient.Script = script
	return gClient, nil
}

func (gc *GuardianClient) SignTx(tx *wire.MsgTx, values []int64) error {
	txInt, ok := gc.Cache[tx.TxHash().String()]
	if !ok {
		hcO, newOutpoints := hasCommonOutpoints(gc.spendUTXOs, tx)
		gc.spendUTXOs = append(gc.spendUTXOs, newOutpoints...)
		if !hcO {
			if err := gc.signTx(tx, values); err != nil {
				return err
			}
		} else {
			if err := gc.updateTx(gc.LatestHash, tx, nil, values); err != nil {
				return err
			}
		}
		gc.Cache[tx.TxHash().String()] = tx
		gc.LatestHash = tx.TxHash()
		return nil
	}
	var buf bytes.Buffer
	if err := txInt.Serialize(&buf); err != nil {
		return err
	}
	return tx.Deserialize(&buf)
}

func getOutpoints(tx *wire.MsgTx) []wire.OutPoint {
	ops := make([]wire.OutPoint, len(tx.TxIn))
	for i, txin := range tx.TxIn {
		ops[i] = txin.PreviousOutPoint
	}
	return ops
}

func hasCommonOutpoints(ops []wire.OutPoint, tx *wire.MsgTx) (bool, []wire.OutPoint) {
	txOps := getOutpoints(tx)
	oldOps := map[string]bool{}
	newOps := []wire.OutPoint{}
	hasDuplicates := false
	for _, op := range ops {
		oldOps[op.String()] = true
	}
	for _, op := range txOps {
		if !oldOps[op.String()] {
			newOps = append(newOps, op)
		} else {
			hasDuplicates = hasDuplicates || true
		}
	}
	return hasDuplicates, newOps
}

func (gc *GuardianClient) signTx(tx *wire.MsgTx, values []int64) error {
	var buf = new(bytes.Buffer)
	if err := tx.Serialize(buf); err != nil {
		return err
	}
	txHex := hex.EncodeToString(buf.Bytes())
	signRequest := map[string]interface{}{
		"address": gc.Account.Address,
		"tx_hex":  txHex,
		"values":  values,
	}
	jsonData, err := json.Marshal(signRequest)
	if err != nil {
		return err
	}
	fmt.Println("txHex", txHex)

	resp, err := gc.Client.Post(fmt.Sprintf("%s/transactions", gc.URL), "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response: %v", err)
		}
		return fmt.Errorf("failed to sign transaction: %s", errMsg)
	}

	var signTxResp SignTxResponse
	err = json.NewDecoder(resp.Body).Decode(&signTxResp)
	if err != nil {
		return err
	}

	signedTxBytes, err := hex.DecodeString(signTxResp.TxHex)
	if err != nil {
		return err
	}

	return tx.Deserialize(strings.NewReader(string(signedTxBytes)))
}

func (gc *GuardianClient) updateTx(prevTxID chainhash.Hash, updatedTx *wire.MsgTx, mergeTxs []*wire.MsgTx, values []int64) error {
	var buf = new(bytes.Buffer)
	updatedTx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())

	mergeTxHexes := make([]string, len(mergeTxs))
	for i, tx := range mergeTxs {
		var mergeBuf = new(bytes.Buffer)
		tx.Serialize(mergeBuf)
		mergeTxHexes[i] = hex.EncodeToString(mergeBuf.Bytes())
	}

	updateRequest := UpdateTxRequest{
		Address:  gc.Account.Address,
		TxHex:    txHex,
		Values:   values,
		MergeTxs: mergeTxHexes,
	}
	jsonData, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/transactions/%s", gc.URL, prevTxID.String()), strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := gc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to update transaction: %s", errMsg)
	}

	var updatedTxResp SignTxResponse
	err = json.NewDecoder(resp.Body).Decode(&updatedTxResp)
	if err != nil {
		return err
	}

	updatedTxBytes, err := hex.DecodeString(updatedTxResp.TxHex)
	if err != nil {
		return err
	}

	return updatedTx.Deserialize(strings.NewReader(string(updatedTxBytes)))
}

func (gc *GuardianClient) GetAccounts() ([]string, error) {
	resp, err := gc.Client.Get(fmt.Sprintf("%s/accounts", gc.URL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get accounts: %s", resp.Status)
	}

	var accounts []string
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (gc *GuardianClient) GetAccount(address btcutil.Address) (*Account, error) {
	resp, err := gc.Client.Get(fmt.Sprintf("%s/accounts/%s", gc.URL, address.EncodeAddress()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get account: %s", resp.Status)
	}

	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (gc *GuardianClient) GetTransaction(txid *chainhash.Hash) (*wire.MsgTx, error) {
	resp, err := gc.Client.Get(fmt.Sprintf("%s/transactions/%s", gc.URL, txid.String()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get transaction: %s", resp.Status)
	}

	var txHexes []string
	err = json.NewDecoder(resp.Body).Decode(&txHexes)
	if err != nil {
		return nil, err
	}

	if len(txHexes) == 0 {
		return nil, fmt.Errorf("no transaction found")
	}

	txBytes, err := hex.DecodeString(txHexes[0])
	if err != nil {
		return nil, err
	}

	var tx wire.MsgTx
	err = tx.Deserialize(strings.NewReader(string(txBytes)))
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (gc *GuardianClient) GetTransactions() ([]*wire.MsgTx, error) {
	resp, err := gc.Client.Get(fmt.Sprintf("%s/transactions", gc.URL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get transactions: %s", resp.Status)
	}

	var txHexes []string
	err = json.NewDecoder(resp.Body).Decode(&txHexes)
	if err != nil {
		return nil, err
	}

	var txs []*wire.MsgTx
	for _, txHex := range txHexes {
		txBytes, err := hex.DecodeString(txHex)
		if err != nil {
			return nil, err
		}

		var tx wire.MsgTx
		err = tx.Deserialize(strings.NewReader(string(txBytes)))
		if err != nil {
			return nil, err
		}

		txs = append(txs, &tx)
	}

	return txs, nil
}

func (gc *GuardianClient) GetGuardianPubkey() (*btcec.PublicKey, error) {
	resp, err := gc.Client.Get(gc.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get guardian pubkey: %s", resp.Status)
	}

	var guardianResp struct {
		PublicKey string `json:"public_key"`
	}
	err = json.NewDecoder(resp.Body).Decode(&guardianResp)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := hex.DecodeString(guardianResp.PublicKey)
	if err != nil {
		return nil, err
	}

	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

// GuardianWallet implements the Wallet interface for guardian-based transactions
type GuardianWallet struct {
	privateKey *btcec.PrivateKey
	client     *GuardianClient
	wallet     Wallet
}

func NewGuardianWallet(url string, network *chaincfg.Params, privateKey *btcec.PrivateKey, indexer IndexerClient, feeEstimator FeeEstimator, feeLevel FeeLevel) (Wallet, error) {
	client, err := NewGuardianClient(url, network, privateKey.PubKey())
	if err != nil {
		return nil, err
	}
	wallet, err := NewSimpleWallet(privateKey, network, indexer, feeEstimator, feeLevel)
	if err != nil {
		return nil, err
	}
	return &GuardianWallet{
		privateKey: privateKey,
		client:     client,
		wallet:     wallet,
	}, nil
}

func (gw *GuardianWallet) Address() btcutil.Address {
	return gw.client.Address
}

func (gw *GuardianWallet) SignCoverUTXOs(tx *wire.MsgTx, utxos UTXOs, startingIdx int) error {
	if len(utxos) == 0 || len(utxos) == startingIdx {
		return nil
	}
	if startingIdx != 0 {
		fmt.Println(len(utxos), startingIdx)
		return fmt.Errorf("guardian wallet only supports spends from itself")
	}

	// for p2wsh, we only need to add the signature and pubkey
	witness := [][]byte{
		{},
		{},
		AddSignatureSegwitOp,
		{0x01},
		gw.client.Script,
	}
	idx := startingIdx
	values := make([]int64, len(utxos))
	for i, utxo := range utxos {
		fetcher := txscript.NewCannedPrevOutputFetcher(gw.client.Script, utxo.Amount)
		err := signTx(tx, fetcher, utxo.Amount, idx, witness, gw.client.Script, nil, txscript.SigHashAll, gw.privateKey)
		if err != nil {
			return err
		}
		values[i] = utxo.Amount
		idx++
	}
	if err := gw.client.SignTx(tx, values); err != nil {
		return err
	}
	return nil
}

func (gw *GuardianWallet) GenerateSACP(ctx context.Context, spendReq SpendRequest, to btcutil.Address) ([]byte, error) {
	return gw.wallet.GenerateSACP(ctx, spendReq, to)
}

func (gw *GuardianWallet) SignSACPTx(tx *wire.MsgTx, idx int, amount int64, leaf txscript.TapLeaf, scriptAddr btcutil.Address, witness [][]byte) ([][]byte, error) {
	return gw.wallet.SignSACPTx(tx, idx, amount, leaf, scriptAddr, witness)
}

func (gw *GuardianWallet) Send(ctx context.Context, sendReq []SendRequest, spendReq []SpendRequest, sacps [][]byte) (string, error) {
	panic("Should never be called")
}

func (gw *GuardianWallet) Status(ctx context.Context, id string) (Transaction, bool, error) {
	panic("Should never be called")
}

func (gw *GuardianWallet) CoverUTXOSpendWeight() int {
	return GuardianSpendWeight
}

type GuardianSpendBatch struct {
	SpendRequests []SpendRequest
	SACPs         [][]byte
}

type GuardianSendBatch struct {
	SendRequests []SendRequest
}

type GuardianCache interface {
	AddSendBatch(GuardianSendBatch) error
	AddSpendBatch(GuardianSpendBatch) error
}

type MultiBatcherWallet struct {
	SendBatcher  BatcherWallet
	SpendBatcher BatcherWallet
}

func NewMultiBatcherWallet(url string, privateKey *secp256k1.PrivateKey, indexer IndexerClient, feeEstimator FeeEstimator, chainParams *chaincfg.Params, sendCache, spendCache Cache, logger *zap.Logger, opts ...func(*batcherWallet) error) (BatcherWallet, error) {
	sendWallet := &batcherWallet{
		indexer:      indexer,
		privateKey:   privateKey,
		cache:        sendCache,
		logger:       logger,
		feeEstimator: feeEstimator,
		chainParams:  chainParams,
		opts:         defaultBatcherOptions(),
	}
	for _, opt := range opts {
		err := opt(sendWallet)
		if err != nil {
			return nil, err
		}
	}

	spendWallet := &batcherWallet{
		indexer:      indexer,
		privateKey:   privateKey,
		cache:        spendCache,
		logger:       logger,
		feeEstimator: feeEstimator,
		chainParams:  chainParams,
		opts:         defaultBatcherOptions(),
	}
	for _, opt := range opts {
		err := opt(spendWallet)
		if err != nil {
			return nil, err
		}
	}

	gWallet, err := NewGuardianWallet(url, chainParams, privateKey, indexer, feeEstimator, sendWallet.opts.TxOptions.FeeLevel)
	if err != nil {
		return nil, err
	}
	// sWallet, err := NewSimpleWallet(privateKey, chainParams, indexer, feeEstimator, sendWallet.opts.TxOptions.FeeLevel)
	// if err != nil {
	// 	return nil, err
	// }
	sendWallet.sw = gWallet
	spendWallet.sw = gWallet
	return &MultiBatcherWallet{SpendBatcher: spendWallet, SendBatcher: sendWallet}, nil
}

func (mbw *MultiBatcherWallet) Send(ctx context.Context, sendReq []SendRequest, spendReq []SpendRequest, sacps [][]byte) (string, error) {
	reqID := ""
	if len(sendReq) > 0 {
		sendReqID, err := mbw.SendBatcher.Send(ctx, sendReq, nil, nil)
		if err != nil {
			return "", err
		}
		reqID += "A" + sendReqID
	}
	if len(spendReq) > 0 || len(sacps) > 0 {
		spendReqID, err := mbw.SpendBatcher.Send(ctx, nil, spendReq, sacps)
		if err != nil {
			return "", err
		}
		if len(reqID) != 0 {
			reqID += ","
		}
		reqID += "B" + spendReqID
	}
	if reqID == "" {
		return "", fmt.Errorf("no requests provided")
	}
	return reqID, nil
}

func (mbw *MultiBatcherWallet) Status(ctx context.Context, id string) (Transaction, bool, error) {
	reqs := strings.Split(id, ",")
	if len(reqs) == 2 {
		sendTx, sendPending, err := mbw.SendBatcher.Status(ctx, reqs[0][1:])
		if err != nil {
			return sendTx, sendPending, err
		}
		spendTx, spendPending, err := mbw.SpendBatcher.Status(ctx, reqs[1][1:])
		if err != nil {
			return spendTx, spendPending, err
		}
		return sendTx, sendPending || spendPending, err
	} else if len(reqs) == 1 {
		fmt.Println(reqs[0])
		fmt.Println(reqs[0][1:])
		if reqs[0][0] == 'A' {
			return mbw.SendBatcher.Status(ctx, reqs[0][1:])
		} else {
			return mbw.SpendBatcher.Status(ctx, reqs[0][1:])
		}
	}
	return Transaction{}, false, fmt.Errorf("invalid status")
}

func (mbw *MultiBatcherWallet) GenerateSACP(ctx context.Context, spendReq SpendRequest, to btcutil.Address) ([]byte, error) {
	return mbw.SpendBatcher.GenerateSACP(ctx, spendReq, to)
}

func (mbw *MultiBatcherWallet) SignSACPTx(tx *wire.MsgTx, idx int, amount int64, leaf txscript.TapLeaf, scriptAddr btcutil.Address, witness [][]byte) ([][]byte, error) {
	return mbw.SpendBatcher.SignSACPTx(tx, idx, amount, leaf, scriptAddr, witness)
}

func (mbw *MultiBatcherWallet) Address() btcutil.Address {
	return mbw.SpendBatcher.Address()
}

func (mbw *MultiBatcherWallet) Start(ctx context.Context) error {
	if err := mbw.SendBatcher.Start(ctx); err != nil {
		return fmt.Errorf("failed to start send batcher: %w", err)
	}
	if err := mbw.SpendBatcher.Start(ctx); err != nil {
		return fmt.Errorf("failed to start spend batcher: %w", err)
	}
	return nil
}

func (mbw *MultiBatcherWallet) Stop() error {
	sendErr := mbw.SendBatcher.Stop()
	spendErr := mbw.SpendBatcher.Stop()

	if sendErr != nil && spendErr != nil {
		return fmt.Errorf("failed to stop both batchers: send error: %v, spend error: %v", sendErr, spendErr)
	}
	if sendErr != nil {
		return fmt.Errorf("failed to stop send batcher: %w", sendErr)
	}
	if spendErr != nil {
		return fmt.Errorf("failed to stop spend batcher: %w", spendErr)
	}
	return nil
}

func (mbw *MultiBatcherWallet) Restart(ctx context.Context) error {
	if err := mbw.Stop(); err != nil {
		return fmt.Errorf("failed to stop during restart: %w", err)
	}
	if err := mbw.Start(ctx); err != nil {
		return fmt.Errorf("failed to start during restart: %w", err)
	}
	return nil
}

func (mbw *MultiBatcherWallet) CoverUTXOSpendWeight() int {
	return mbw.SpendBatcher.CoverUTXOSpendWeight()
}

func (mbw *MultiBatcherWallet) SignCoverUTXOs(tx *wire.MsgTx, utxos UTXOs, startingIdx int) error {
	return mbw.SpendBatcher.SignCoverUTXOs(tx, utxos, startingIdx)
}
