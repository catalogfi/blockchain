package cosigner

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

// type Account struct {
// 	Address           string          `json:"address"`
// 	GuardianPublicKey string          `json:"guardian_public_key"`
// 	UserPublicKey     string          `json:"user_public_key"`
// 	SpentOutpoints    map[string]bool `json:"spent_outpoints"`
// }

type Transaction struct {
	Values              map[wire.OutPoint]int64           `json:"values"`
	MergeTransactions   map[string]*wire.MsgTx            `json:"merge_transactions"`
	MempoolTransactions map[string]*wire.MsgTx            `json:"mempool_transactions"`
	TxHashes            []string                          `json:"txhashes"`
	Spends              []map[wire.OutPoint]wire.TxOut    `json:"spends"`
	SpendLinks          []map[wire.OutPoint]wire.OutPoint `json:"spendlinks"`
	BackupTransactions  []map[string]*wire.MsgTx          `json:"backup_transactions"`
}

type Client struct {
	hc  *http.Client
	url string
}

func (client *Client) NewAccount(pub btcec.PublicKey) (string, error) {
	pubStr := hex.EncodeToString(pub.SerializeCompressed())
	var req = struct {
		PublicKey string `json:"public_key"`
	}{
		PublicKey: pubStr,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%v/accounts", client.url)
	response, err := client.hc.Post(path, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(response.Status)
	}

	var resp struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
		return "", err
	}
	return resp.Address, nil
}

func (client *Client) NewTransaction(address string, unsignedTx *wire.MsgTx) (*wire.MsgTx, error) {
	path := fmt.Sprintf("%v/account/%v/transactions", client.url, address)
	buff := bytes.NewBuffer(make([]byte, unsignedTx.SerializeSize()))
	if err := unsignedTx.Serialize(buff); err != nil {
		return nil, err
	}

	request := struct {
		UnsignedTx string `json:"unsigned_tx"`
	}{
		hex.EncodeToString(buff.Bytes()),
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := client.hc.Post(path, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(response.Status)
	}

	var resp struct {
		SignedTx string `json:"signed_tx"`
	}
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
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

func (client *Client) GetTransactions(address string) ([]Transaction, error) {
	path := fmt.Sprintf("%v/account/%v/transactions", client.url, address)
	response, err := client.hc.Get(path)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(response.Status)
	}
	var txs []Transaction
	if err := json.NewDecoder(response.Body).Decode(&txs); err != nil {
		return nil, err
	}
	return txs, nil
}

func (client *Client) GetLatestTransaction(address string) (*Transaction, error) {
	path := fmt.Sprintf("%v/account/%v/transaction/latest", client.url, address)
	response, err := client.hc.Get(path)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf(response.Status)
	}
	var tx Transaction
	if err := json.NewDecoder(response.Body).Decode(&tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (client *Client) GetTransactionByNonce(address string, nonce int) (*Transaction, error) {
	path := fmt.Sprintf("%v/account/%v/transaction/%v", client.url, address, nonce)
	response, err := client.hc.Get(path)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf(response.Status)
	}
	var tx Transaction
	if err := json.NewDecoder(response.Body).Decode(&tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (client *Client) UpdateTransaction(address string, newTx *wire.MsgTx, backupTxs []*wire.MsgTx) (*wire.MsgTx, error) {
	path := fmt.Sprintf("%v/account/%v/transaction/latest", client.url, address)

	backupMap := make(map[string]string)
	for _, tx := range backupTxs {
		buffer := bytes.NewBuffer([]byte{})
		if err := tx.Serialize(buffer); err != nil {
			return nil, err
		}
		backupMap[tx.TxHash().String()] = hex.EncodeToString(buffer.Bytes())
	}
	buff := bytes.NewBuffer(make([]byte, newTx.SerializeSize()))
	if err := newTx.Serialize(buff); err != nil {
		return nil, err
	}
	request := struct {
		BackupTxs  map[string]string `json:"backup_txs"`
		UnsignedTx string            `json:"unsigned_txF"`
	}{
		BackupTxs:  backupMap,
		UnsignedTx: hex.EncodeToString(buff.Bytes()),
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.hc.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	var resp struct {
		SignedTx string `json:"signed_tx"`
	}
	if err := json.NewDecoder(response.Body).Decode(&resp); err != nil {
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
