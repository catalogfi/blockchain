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

There're a few commonly used fee estimator for you to give you an estimate of the current network. Beware these are 
mostly for mainnet, as min relay fees usually is enough for testnet and regnet. 
> The result is in `sats/vB`

Example:

```go
    // Mempool uses the mempool's public fee estimation api (https://mempool.space/docs/api/rest#get-recommended-fees). 
    // It returns currently suggested fees for new transactions. It's safe for concurrent use and you can define a 
    // duration for how long you want to cache the result to avoid spamming the requests.
    estimator := btc.NewMempoolFeeEstimator(&chaincfg.MainNetParams, btc.MempoolFeeAPI, 15*time.Second)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    }
    
    // Blockstream uses the blockstream api and returns a fee estimation basing on the past blocks. 
    estimator := btc.NewBlockstreamFeeEstimator(&chaincfg.MainNetParams, btc.BlockstreamAPI, 15*time.Second)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    } 
    
    // If you know the exact fee rate you want, you can use the FixedFeeEstimator which will always returns the provided 
    // fee rate. 
    estimator := btc.NewFixedFeeEstimator(10)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    }
```

### Build a bitcoin transaction 

One of the important use case for this library is to build a bitcoin transaction. We use the `BuildTransaction` function 
to do this most of time. Don't be scared by the number of parameters of this function. These parameters can be divided 
into three categories. 

1. General 
- `network` is the network where this tx is built on. 
- `feeRate` is a minimum feeRate you want to use for the tx. Usually the actual tx will be a little over the feeRate due 
  to the estimation will always use the upperbound. 
2. Inputs 
- `inputs` are the utxos you want to spend in the tx, this means they are guaranteed to be included in the tx. You'll 
  need to provide the estimate size of these utxos. Use the default value `NewRawInputs()` if you don't have any utxo 
  want to be specified
- `utxos` is the available utxos we can add to the inputs if the `inputs` amount is not enough to cover the output. 
  It simply adds the utxo with the given order from the list. Use the `nil` value if you don't have this. 
- `sizeUpdater` describes how much sizes each of the `utxos` will add to the tx. It assumes all the `utxos` are coming 
  from the same address. There are predefined `sizeUpdate` for you to grab and use. i.e. `P2pkhUpdater` and `P2wpkhUpdater`
  Use the `nil` value if `utxos` is nil. 
3. Outputs
- `recipients` defines who will receive the funds, like the `inputs`, all the recipients will be guaranteed to receive
  the provided amount. 
- `changeAddr` is where you want to send the change to. Usually this will be the sender's address and the change will be
  sent back to him.

Some examples

1. Transfer 0.1 btc to `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`
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

2. Spend an utxo and send all the money to `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`

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

3. Redeem an htcl and send  0.1 btc to `1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa` and rest to sender

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

- `MultisigScript`: Returns a 2-of-2 multisig script used by the [Guardian](https://docs.catalog.fi/catalog-accounts/instant-wallet/guardian) component.
- `HtlcScript`: Returns a HTLC script as described in [BIP-199](https://github.com/bitcoin/bips/blob/e643d247c8bc086745f3031cdee0899803edea2f/bip-0199.mediawiki#L22).

[tests-url]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml
[tests-badge]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml/badge.svg?branch=master
