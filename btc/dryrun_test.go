package btc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/waddrmgr"
	"go.uber.org/zap"
)

type dryRunTestIndexer struct {
	txs map[string]Transaction
}

func (m *dryRunTestIndexer) GetAddressTxs(context.Context, btcutil.Address, string) ([]Transaction, error) {
	return nil, nil
}
func (m *dryRunTestIndexer) GetUTXOs(context.Context, btcutil.Address) (UTXOs, error) {
	return nil, nil
}
func (m *dryRunTestIndexer) GetUTXOsForAmount(context.Context, btcutil.Address, int64) (UTXOs, int64, error) {
	return nil, 0, nil
}
func (m *dryRunTestIndexer) GetTipBlockHeight(context.Context) (uint64, error) { return 0, nil }
func (m *dryRunTestIndexer) GetTx(_ context.Context, txid string) (Transaction, error) {
	return m.txs[txid], nil
}
func (m *dryRunTestIndexer) GetTxHex(context.Context, string) (string, error) { return "", nil }
func (m *dryRunTestIndexer) SubmitTx(context.Context, *wire.MsgTx) error      { return nil }
func (m *dryRunTestIndexer) FeeEstimate(context.Context) (FeeSuggestion, error) {
	return FeeSuggestion{}, nil
}

func TestClassifyRejectReason(t *testing.T) {
	tests := []struct {
		name   string
		reason string
		want   rejectCategory
	}{
		{name: "missing or spent", reason: "bad-txns-inputs-missingorspent", want: rejectMissingOrSpent},
		{name: "script verify", reason: "mandatory-script-verify-flag-failed (Signature must be zero for failed CHECK(MULTI)SIG operation)", want: rejectScriptVerify},
		{name: "adds unconfirmed", reason: "replacement-adds-unconfirmed, replacement eb9396e92324ef6f0d382eedd6774b12d0b8a17ce09e224d4e0d7d4ec9c54405 adds unconfirmed input, idx 4", want: rejectAddsUnconfirmed},
		{name: "scrap fee", reason: "insufficient fee", want: rejectScrap},
		{name: "scrap dust", reason: "dust", want: rejectScrap},
		{name: "scrap empty", reason: "", want: rejectScrap},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := classifyRejectReason(tc.reason); got != tc.want {
				t.Fatalf("classifyRejectReason(%q) = %v, want %v", tc.reason, got, tc.want)
			}
		})
	}
}

func TestBuildOutpointOwnerMapIncludesSpendsAndSACPs(t *testing.T) {
	sacpOutpoint := newTestOutPoint(t, "0000000000000000000000000000000000000000000000000000000000000002", 1)
	sacpBytes := mustSerializeTx(t, testTxWithInputs(sacpOutpoint))

	requests := []BatcherRequest{
		{
			ID: "req-spend",
			Spends: []SpendRequest{{
				Utxos: UTXOs{{TxID: "0000000000000000000000000000000000000000000000000000000000000001", Vout: 0, Amount: 5_000}},
			}},
		},
		{
			ID:    "req-sacp",
			SACPs: [][]byte{sacpBytes},
		},
	}

	owners, err := buildOutpointOwnerMap(requests)
	if err != nil {
		t.Fatalf("buildOutpointOwnerMap: %v", err)
	}

	if owners["0000000000000000000000000000000000000000000000000000000000000001:0"] != "req-spend" {
		t.Fatalf("expected spend outpoint to map to req-spend, got %q", owners["0000000000000000000000000000000000000000000000000000000000000001:0"])
	}
	if owners[outpointKey(sacpOutpoint)] != "req-sacp" {
		t.Fatalf("expected sacp outpoint to map to req-sacp, got %q", owners[outpointKey(sacpOutpoint)])
	}
}

func TestFindSpentInputOwnerReturnsMatchingRequest(t *testing.T) {
	rpc := newGetTxOutClient(t, map[string]bool{
		"0000000000000000000000000000000000000000000000000000000000000001:0": true,
	})
	wallet := &batcherWallet{rpc: rpc, logger: zap.NewNop()}
	tx := testTxWithInputs(
		newTestOutPoint(t, "0000000000000000000000000000000000000000000000000000000000000001", 0),
		newTestOutPoint(t, "0000000000000000000000000000000000000000000000000000000000000009", 0),
	)
	requests := []BatcherRequest{{
		ID: "spent-request",
		Spends: []SpendRequest{{
			Utxos: UTXOs{{TxID: "0000000000000000000000000000000000000000000000000000000000000001", Vout: 0, Amount: 7_000}},
		}},
	}}

	requestID, found, err := wallet.findSpentInputOwner(context.Background(), tx, requests)
	if err != nil {
		t.Fatalf("findSpentInputOwner: %v", err)
	}
	if !found || requestID != "spent-request" {
		t.Fatalf("findSpentInputOwner = (%q, %v), want (spent-request, true)", requestID, found)
	}
}

func TestFindSpentInputOwnerIgnoresWalletOnlyInputs(t *testing.T) {
	rpc := newGetTxOutClient(t, map[string]bool{
		"0000000000000000000000000000000000000000000000000000000000000009:0": true,
	})
	wallet := &batcherWallet{rpc: rpc, logger: zap.NewNop()}
	tx := testTxWithInputs(
		newTestOutPoint(t, "0000000000000000000000000000000000000000000000000000000000000001", 0),
		newTestOutPoint(t, "0000000000000000000000000000000000000000000000000000000000000009", 0),
	)
	requests := []BatcherRequest{{
		ID: "tracked-request",
		Spends: []SpendRequest{{
			Utxos: UTXOs{{TxID: "0000000000000000000000000000000000000000000000000000000000000001", Vout: 0, Amount: 7_000}},
		}},
	}}

	requestID, found, err := wallet.findSpentInputOwner(context.Background(), tx, requests)
	if err != nil {
		t.Fatalf("findSpentInputOwner: %v", err)
	}
	if found || requestID != "" {
		t.Fatalf("findSpentInputOwner = (%q, %v), want (\"\", false)", requestID, found)
	}
}

func TestFindBadScriptInputOwnerReturnsMatchingRequest(t *testing.T) {
	chainParams := &chaincfg.RegressionNetParams
	signerKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("new private key: %v", err)
	}
	recipientKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("new private key: %v", err)
	}

	signerAddr, err := PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, signerKey.PubKey())
	if err != nil {
		t.Fatalf("signer address: %v", err)
	}
	recipientAddr, err := PublicKeyAddress(chainParams, waddrmgr.WitnessPubKey, recipientKey.PubKey())
	if err != nil {
		t.Fatalf("recipient address: %v", err)
	}
	recipientScript, err := txscript.PayToAddrScript(recipientAddr)
	if err != nil {
		t.Fatalf("recipient script: %v", err)
	}

	outpoint := newTestOutPoint(t, "0000000000000000000000000000000000000000000000000000000000000003", 0)
	tx := wire.NewMsgTx(DefaultTxVersion)
	tx.AddTxIn(wire.NewTxIn(&outpoint, nil, nil))
	tx.TxIn[0].Witness = wire.TxWitness{[]byte{0x01}, []byte{0x02}}
	tx.AddTxOut(wire.NewTxOut(4_000, recipientScript))

	requests := []BatcherRequest{{
		ID: "bad-script-request",
		Spends: []SpendRequest{{
			ScriptAddress: signerAddr,
			Script:        mustPayToAddrScript(t, signerAddr),
			HashType:      txscript.SigHashAll,
			Utxos:         UTXOs{{TxID: outpoint.Hash.String(), Vout: outpoint.Index, Amount: 5_000}},
		}},
	}}

	wallet := &batcherWallet{
		indexer: &dryRunTestIndexer{},
		logger:  zap.NewNop(),
	}

	requestID, found, err := wallet.findBadScriptInputOwner(context.Background(), tx, requests)
	if err != nil {
		t.Fatalf("findBadScriptInputOwner: %v", err)
	}
	if !found || requestID != "bad-script-request" {
		t.Fatalf("findBadScriptInputOwner = (%q, %v), want (bad-script-request, true)", requestID, found)
	}
}

func newGetTxOutClient(t *testing.T, spent map[string]bool) BitcoinClient {
	t.Helper()

	client := NewBitcoinClient("", "", "http://rpc.invalid")
	client.httpClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var req struct {
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}

		if req.Method != "gettxout" {
			t.Fatalf("unexpected rpc method: %s", req.Method)
		}

		var txid string
		var vout uint32
		if err := json.Unmarshal(req.Params[0], &txid); err != nil {
			t.Fatalf("decode txid: %v", err)
		}
		if err := json.Unmarshal(req.Params[1], &vout); err != nil {
			t.Fatalf("decode vout: %v", err)
		}

		resp := struct {
			Result any   `json:"result"`
			Error  any   `json:"error"`
			ID     int64 `json:"id"`
		}{
			Error: nil,
			ID:    1,
		}

		key := txid + ":" + itoa(vout)
		if spent[key] {
			resp.Result = nil
		} else {
			resp.Result = btcjson.GetTxOutResult{
				Value: 0.0001,
			}
		}

		raw, err := json.Marshal(resp)
		if err != nil {
			return nil, err
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(raw)),
			Header:     make(http.Header),
		}, nil
	})}
	return client
}

func newTestOutPoint(t *testing.T, txid string, vout uint32) wire.OutPoint {
	t.Helper()
	hash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		t.Fatalf("new hash: %v", err)
	}
	return wire.OutPoint{Hash: *hash, Index: vout}
}

func testTxWithInputs(outpoints ...wire.OutPoint) *wire.MsgTx {
	tx := wire.NewMsgTx(DefaultTxVersion)
	for _, outpoint := range outpoints {
		tx.AddTxIn(wire.NewTxIn(&outpoint, nil, nil))
	}
	tx.AddTxOut(wire.NewTxOut(1_000, []byte{txscript.OP_TRUE}))
	return tx
}

func mustSerializeTx(t *testing.T, tx *wire.MsgTx) []byte {
	t.Helper()
	raw, err := GetTxRawBytes(tx)
	if err != nil {
		t.Fatalf("serialize tx: %v", err)
	}
	return raw
}

func mustPayToAddrScript(t *testing.T, addr btcutil.Address) []byte {
	t.Helper()
	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		t.Fatalf("pay to addr script: %v", err)
	}
	return script
}

func itoa(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
