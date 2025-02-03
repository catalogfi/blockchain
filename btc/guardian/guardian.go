package guardian

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
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// A Guardian Client responsible for creating accounts, signing, and updating transactions
// with a guardian backend
//
// It has no state and it is responsibility of the caller to manage the state of the
// accounts.
type Client struct {
	url        string
	httpClient *http.Client
}

type Account struct {
	Address           string `json:"address"`
	GuardianPublicKey string `json:"guardian_public_key"`
	UserPublicKey     string `json:"user_public_key"`
}

type signTxResponse struct {
	TxHex string `json:"tx_hex"`
}

func NewClient(url string) *Client {
	return &Client{
		url:        url,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) CreateAccountIfNotExists(ctx context.Context, publicKey *btcec.PublicKey) (*Account, error) {
	account, err := c.GetAccount(ctx, publicKey)
	if err != nil {
		if strings.Contains(err.Error(), "account not found") {
			return c.CreateAccount(ctx, publicKey)
		}
		return nil, err
	}
	return account, nil
}

func (c *Client) CreateAccount(ctx context.Context, publicKey *btcec.PublicKey) (*Account, error) {
	// request body
	requestBody := map[string]string{
		"public_key": hex.EncodeToString(publicKey.SerializeCompressed()),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// make a http request
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/accounts", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("account creation failed to read error response: %w", err)
		}
		return nil, fmt.Errorf("account creation failed: %s", string(body))
	}

	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("account creation failed to decode response: %w", err)
	}
	return &account, nil
}

func (c *Client) GetGuardianPubkey(ctx context.Context) (*btcec.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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

func (c *Client) SignTransaction(ctx context.Context, accountAddress string, tx *wire.MsgTx, values []int64, mergeTxHexes []string) error {
	var buf = new(bytes.Buffer)
	if err := tx.Serialize(buf); err != nil {
		return err
	}
	txHex := hex.EncodeToString(buf.Bytes())
	signRequest := map[string]interface{}{
		"account_address": accountAddress,
		"tx_hex":          txHex,
		"values":          values,
		"merge_txs":       mergeTxHexes,
	}
	jsonData, err := json.Marshal(signRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/transactions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("sign transaction failed to read error response: %w", err)
		}
		return fmt.Errorf("sign transaction failed to decode response: %s", string(errMsg))
	}
	var signTxResp signTxResponse
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

func (c *Client) UpdateTransaction(
	ctx context.Context,
	accountAddr string,
	prevTxID chainhash.Hash,
	updatedTx *wire.MsgTx,
	mergeTxs []string,
	values []int64,
) error {
	var buf = new(bytes.Buffer)
	updatedTx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())

	type UpdateTxRequest struct {
		Address  string   `json:"address"`
		TxHex    string   `json:"tx_hex"`
		Values   []int64  `json:"values"`
		MergeTxs []string `json:"merge_txs"`
	}
	updateRequest := UpdateTxRequest{
		Address:  accountAddr,
		TxHex:    txHex,
		Values:   values,
		MergeTxs: mergeTxs,
	}
	jsonData, err := json.Marshal(updateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/transactions/%s", c.url, prevTxID.String()), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("update transaction failed to read error response: %w", err)
		}
		return fmt.Errorf("update transaction failed : %s", string(errMsg))
	}

	var updatedTxResp UpdateTxRequest
	err = json.NewDecoder(resp.Body).Decode(&updatedTxResp)
	if err != nil {
		return fmt.Errorf("update transaction failed to decode response: %w", err)
	}

	updatedTxBytes, err := hex.DecodeString(updatedTxResp.TxHex)
	if err != nil {
		return fmt.Errorf("update transaction failed to decode response: %w", err)
	}

	return updatedTx.Deserialize(strings.NewReader(string(updatedTxBytes)))
}

func (c *Client) GetAccount(ctx context.Context, publicKey *btcec.PublicKey) (*Account, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/accounts/%s", c.url, hex.EncodeToString(publicKey.SerializeCompressed())), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("get account failed to read error response: %w", err)
		}
		return nil, fmt.Errorf("get account failed: %s", string(errMsg))
	}
	var account Account
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("get account failed to decode response: %w", err)
	}

	return &account, nil
}

// Creates a guardian client and caches the account for further use
type GuardianClient struct {
	c       *Client
	account *Account
}

func NewGuardianClient(ctx context.Context, url string, userPubKey *btcec.PublicKey) *GuardianClient {
	client := NewClient(url)
	account, err := client.CreateAccountIfNotExists(ctx, userPubKey)
	if err != nil {
		panic(err)
	}
	return &GuardianClient{
		c:       client,
		account: account,
	}
}

func (gc *GuardianClient) GetGuardianPubkey() (*btcec.PublicKey, error) {
	return gc.c.GetGuardianPubkey(context.Background())
}

func (gc *GuardianClient) GetAccount() *Account {
	return gc.account
}

func (gc *GuardianClient) SignTransaction(ctx context.Context, tx *wire.MsgTx, values []int64, mergeTxHexes []string) error {
	return gc.c.SignTransaction(ctx, gc.account.Address, tx, values, mergeTxHexes)
}

func (gc *GuardianClient) UpdateTransaction(ctx context.Context, prevTxID chainhash.Hash, updatedTx *wire.MsgTx, mergeTxs []string, values []int64) error {
	return gc.c.UpdateTransaction(ctx, gc.account.Address, prevTxID, updatedTx, mergeTxs, values)
}
