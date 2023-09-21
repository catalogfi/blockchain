# ⛓️ blockchain

[![Tests][tests-badge]][tests-url]
[![Go Coverage](https://github.com/catalogfi/blockchain/wiki/coverage.svg)](https://raw.githack.com/wiki/catalogfi/blockchain/coverage.html)

For any Catalog-related blockchain interactions, in Golang.
-  [x] Bitcoin
-  [ ] Ethereum

## Install

The `blockchain` package can be imported by running:
```shell
$ go get github.com/catalogfi/blockchain
```

## Bitcoin

### JSON-RPC

This client follows the standard bitcoind JSON-RPC interface.

Example:

```go
    // Initialise the client.
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
    // Initialise the client.
    logger, _ := zap.NewDevelopment()
    indexer := btc.NewElectrsIndexerClient(logger, host, btc.DefaultRetryInterval)
    
    addr := "xxxxxxx"
    utxos, err := indexer.GetUTXOs(context.Background(), addr)
    if err != nil {
        panic(err)
    }
```

### Fee estimator

Example:

```go
    // Mempool
    estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, 15*time.Second)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    }
    
    // Blockstream
    estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, 15*time.Second)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    } 
    
    // Fixed fee estimator
    estimator := btc.NewFixedFeeEstimator(10)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    }
```

### Bitcoin scripts

- `MultisigScript`: Returns a 2-of-2 multisig script used by the [Guardian](https://docs.catalog.fi/catalog-accounts/instant-wallet/guardian) component.
- `HtlcScript`: Returns a HTLC script as described in [BIP-199](https://github.com/bitcoin/bips/blob/e643d247c8bc086745f3031cdee0899803edea2f/bip-0199.mediawiki#L22).

[tests-url]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml
[tests-badge]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml/badge.svg?branch=master
