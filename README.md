# ⛓️ blockchain

[![Tests][tests-badge]][tests-url]
[![Go Coverage](https://github.com/catalogfi/blockchain/wiki/coverage.svg)](https://raw.githack.com/wiki/catalogfi/blockchain/coverage.html)

For any Catalog-related blockchain interactions, in Golang.
-  [x] Bitcoin
-  [x] Ethereum

## Install

To import the `blockchain` package:
```shell
$ go get github.com/catalogfi/blockchain
```

## Bitcoin

### JSON-RPC

This client follows the standard bitcoind JSON-RPC interface.

Example:

```go
    // Initialize the client.
    config := &rpcclient.ConnConfig{
        Params:       chaincfg.RegressionNetParams.Name,
        Host:         "0.0.0.0:18443",
        User:         user,
        Pass:         password,
        HTTPPostMode: true,
        DisableTLS:   true,
    }
    client := btc.NewClient(config)
    
    // Get the latest block.
    height, hash, err := client.LatestBlock()
    if err != nil {
    	panic(err)
    }
```

### Indexer

The indexer client follows the [electrs indexer API](https://github.com/blockstream/esplora/blob/master/API.md).

Example:

```go
    // Initialize the client.
    logger, _ := zap.NewDevelopment()
    indexer := btc.NewElectrsIndexerClient(logger, host, btc.DefaultRetryInterval)
    
    addr := "xxxxxxx"
    utxos, err := indexer.GetUTXOs(context.Background(), addr)
    if err != nil {
        panic(err)
    }
```

### Fee estimator

The fee estimator estimates network fees using various APIs.
> The result is in `sats/vB`

Example:

```go
    // Mempool API
    estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, 15*time.Second)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    }
    
    // Blockstream API
    estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, 15*time.Second)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    } 
    
    // Fixed fee rate
    estimator := btc.NewFixedFeeEstimator(10)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    }
```

### Build a bitcoin transaction 

We use the `BuildTransaction` function with these parameters:

1. General 
- `network`: The network for the transaction.
- `feeRate`: Minimum fee rate; actual fee may be higher.
  
2. Inputs
- `inputs`: UTXOs to spend. Use NewRawInputs() if unspecified.
- `utxos`: Additional UTXOs if needed. Use `nil` if none.
- `sizeUpdater`: Describes UTXO size impact. It assumes all UTXOs come from the same address. Use predefined updaters like `P2pkhUpdater` and `P2wpkhUpdater`, or `nil` if `utxos` is empty.
  
3. Outputs
- `recipients`: Fund recipients with specified amounts.
- `changeAddr`: Address for the change, usually the sender's address.

**Examples:**

1. Transfer 0.1 btc
```go
    // Fetch available utxos of the address 
    utxos, err := indexer.GetUTXOs(ctx, sender)
	if err != nil {
	    panic(err)	
    }
    recipients := []btc.Recipient{
        {
            To:     "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
            Amount: 1e7,
        },
    }
	
    transaction, err := btc.BuildTransaction(&chaincfg.MainNetParams, 20, btc.NewRawInputs(), utxos, btc.P2pkhUpdater, recipients, sender)
    Expect(err).To(BeNil())
```

2. Spend a UTXO and send all to a recipient

```go
    rawInputs := btc.RawInputs{
        VIN:        utxos,
        BaseSize:   txsizes.RedeemP2PKHSigScriptSize * len(utxos),
        SegwitSize: 0,
    }
	sender := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	
    transaction, err := btc.BuildTransaction(&chaincfg.MainNetParams, 20, rawInputs, nil, nil, nil, sender)
    Expect(err).To(BeNil())
```

3. Redeem an HTLC and send to a recipient and sender

```go
    htlcUtxos := []btc.UTXO{
        {
            TxID:   txid,
            Vout:   vout,
            Amount: amount,
        },
    }
    rawInputs := btc.RawInputs{
        VIN:        htlcUtxos,
        BaseSize:   0,
        SegwitSize: btc.RedeemHtlcRefundSigScriptSize,
    }
    recipients := []btc.Recipient{
        {
            To:     "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
            Amount: 1e7,
        },
    }
	
    transaction, err := btc.BuildTransaction(&chaincfg.MainNetParams, 20, rawInputs, nil, nil, recipients, sender)
    Expect(err).To(BeNil())
```

### Bitcoin scripts

- `MultisigScript`: 2-of-2 multisig script for the [Guardian](https://docs.catalog.fi/catalog-accounts/instant-wallet/guardian) component.
- `HtlcScript`: HTLC script as described in [BIP-199](https://github.com/bitcoin/bips/blob/e643d247c8bc086745f3031cdee0899803edea2f/bip-0199.mediawiki#L22).

[tests-url]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml
[tests-badge]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml/badge.svg?branch=master
