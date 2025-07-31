package cosigner

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

type Transaction struct {
	Values              map[string]int64  `json:"values"`
	MergeTransactions   map[string]string `json:"merge_transactions"`
	MempoolTransactions map[string]string `json:"mempool_transactions"`
	TxHashes            []string          `json:"txhashes"`
	Spends              map[string]string `json:"spends"`
	SpendLinks          map[string]string `json:"spendlinks"`
	BackupTransactions  map[string]string `json:"backup_transactions"`
}

type Client struct {
	hc  *http.Client
	url string
}

func NewClient(url string) *Client {
	return &Client{
		hc:  &http.Client{},
		url: url,
	}
}

// CosignerPub queries the public key of the cosigner server.
func (client *Client) CosignerPub() (*btcec.PublicKey, error) {
	path := fmt.Sprintf("%v", client.url)
	var resp struct {
		Pub string `json:"public_key"`
	}
	if err := client.do("GET", path, http.StatusOK, nil, &resp); err != nil {
		return nil, err
	}
	raw, err := hex.DecodeString(resp.Pub)
	if err != nil {
		return nil, err
	}
	return btcec.ParsePubKey(raw)
}

// NewAccount creates a new account with given public key. It returns the address of the new account.
func (client *Client) NewAccount(pub *btcec.PublicKey) (string, error) {
	pubStr := hex.EncodeToString(pub.SerializeCompressed())
	var req = struct {
		PublicKey string `json:"public_key"`
	}{
		PublicKey: pubStr,
	}
	var resp struct {
		Address string `json:"address"`
	}
	path := fmt.Sprintf("%v/accounts", client.url)
	if err := client.do("POST", path, http.StatusCreated, req, &resp); err != nil {
		return "", err
	}
	return resp.Address, nil
}

// NewTransaction posts a new transaction of the given cosigner account. It returns the transaction with cosigner's
// signature.
func (client *Client) NewTransaction(address string, unsignedTx *wire.MsgTx) (*wire.MsgTx, error) {
	path := fmt.Sprintf("%v/account/%v/transactions", client.url, address)
	buff := bytes.NewBuffer(make([]byte, 0, unsignedTx.SerializeSize()))
	if err := unsignedTx.Serialize(buff); err != nil {
		return nil, err
	}
	request := struct {
		UnsignedTx string `json:"unsigned_tx"`
	}{
		hex.EncodeToString(buff.Bytes()),
	}
	var resp struct {
		SignedTx string `json:"signed_tx"`
	}
	if err := client.do("POST", path, http.StatusCreated, request, &resp); err != nil {
		return nil, err
	}

	signedTxBytes, err := hex.DecodeString(resp.SignedTx)
	if err != nil {
		return nil, err
	}
	tx, err := btcutil.NewTxFromBytes(signedTxBytes)
	if err != nil {
		return nil, err
	}
	return tx.MsgTx(), nil
}

// GetTransactions returns the current transaction status of the address.
func (client *Client) GetTransactions(address string) ([]Transaction, error) {
	path := fmt.Sprintf("%v/account/%v/transactions", client.url, address)
	var txs []Transaction
	err := client.do("GET", path, http.StatusOK, nil, &txs)
	return txs, err
}

func (client *Client) GetLatestTransaction(address string) (*Transaction, error) {
	path := fmt.Sprintf("%v/account/%v/transaction/latest", client.url, address)
	// todo : handle  status code not found
	var tx Transaction
	err := client.do("GET", path, http.StatusOK, nil, &tx)
	return &tx, err
}

func (client *Client) GetTransactionByNonce(address string, nonce int) (*Transaction, error) {
	path := fmt.Sprintf("%v/account/%v/transaction/%v", client.url, address, nonce)
	var tx Transaction
	err := client.do("GET", path, http.StatusOK, nil, &tx)
	return &tx, err
}

func (client *Client) UpdateTransaction(address string, newTx *wire.MsgTx, backupTxs map[string]*wire.MsgTx) (*wire.MsgTx, error) {
	path := fmt.Sprintf("%v/account/%v/transaction/latest", client.url, address)
	backupMap := make(map[string]string)
	for txid, tx := range backupTxs {
		buffer := bytes.NewBuffer([]byte{})
		if err := tx.Serialize(buffer); err != nil {
			return nil, err
		}
		backupMap[txid] = hex.EncodeToString(buffer.Bytes())
	}
	buff := bytes.NewBuffer(make([]byte, 0, newTx.SerializeSize()))
	if err := newTx.Serialize(buff); err != nil {
		return nil, err
	}
	request := struct {
		BackupTxs  map[string]string `json:"backup_txs"`
		UnsignedTx string            `json:"unsigned_tx"`
	}{
		BackupTxs:  backupMap,
		UnsignedTx: hex.EncodeToString(buff.Bytes()),
	}
	var resp struct {
		SignedTx string `json:"signed_tx"`
	}
	if err := client.do("PUT", path, http.StatusOK, request, &resp); err != nil {
		return nil, err
	}

	signedTxBytes, err := hex.DecodeString(resp.SignedTx)
	if err != nil {
		return nil, err
	}
	tx, err := btcutil.NewTxFromBytes(signedTxBytes)
	if err != nil {
		return nil, err
	}
	return tx.MsgTx(), nil
}

func (client *Client) NewMergeTx(tx *wire.MsgTx) error {
	path := fmt.Sprintf("%v/transactions", client.url)
	buff := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buff); err != nil {
		return err
	}

	request := struct {
		TxHex string `json:"txhex"`
	}{
		hex.EncodeToString(buff.Bytes()),
	}
	return client.do("POST", path, http.StatusCreated, request, nil)
}

func (client *Client) do(method, path string, code int, request, resp any) error {
	// Construct the request
	var body io.Reader
	if request != nil {
		data, err := json.Marshal(request)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return err
	}
	switch method {
	case "POST", "PUT":
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := client.hc.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != code {
		msg, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("code [%v] msg [%v]", response.Status, string(msg))
	}
	if resp != nil {
		return json.NewDecoder(response.Body).Decode(&resp)
	}
	return nil
}
