# blockchain
[![Tests][tests-badge]][tests-url] 
The blockchain repo provides things we need to interact with different blockchains 
-  [x] bitcoin 
-  [ ] ethereum 

## Install 
The `blockchain` package can be imported directly to your golang project by 
```shell
$ go get github.com/catalogfi/blockchain
```

## Bitcoin 
- json-rpc client

This follows the standard bitcoind json-rpc 
```go
    // Initialise the client 
    config := &rpcclient.ConnConfig{
        Params:       chaincfg.RegressionNetParams.Name,
        Host:         "0.0.0.0:18443",
        User:         user,
        Pass:         password,
        HTTPPostMode: true,
        DisableTLS:   true,
    }
    client := btc.NewClient(config) 
    
    // Get the latest block
    height, hash, err := client.LatestBlock()
    if err != nil {
    	panic(err)
    }
```

- Indexer client

This follow the electrs indexer API, more details can be found [here](https://github.com/blockstream/esplora/blob/master/API.md)
```go
    // Initialise the client
    logger, _ := zap.NewDevelopment()
    indexer := btc.NewElectrsIndexerClient(logger, host, btc.DefaultRetryInterval)
    
    addr := "xxxxxxx"
    utxos, err := indexer.GetUTXOs(context.Background(), addr)
    if err != nil {
        panic(err)
    }
```

- FeeEstimator

There're a few different implementations for the FeeEstimator 

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
	
	// Fix fee estimator 
	estimator := btc.NewFixFeeEstimator(10)
    fees, err := estimator.FeeSuggestion()
    if err != nil {
        panic(err)
    } 
```

- Bitcoin script 

There're some bitcoin scripts we focus in this package, since they're fundamental to other projects we build in catalog.
The `MultisigScript` will return a 2-out-2 multisig script and `HtlcScript` will return a HTLC script which is described 
in [BIP-199](https://github.com/bitcoin/bips/blob/e643d247c8bc086745f3031cdee0899803edea2f/bip-0199.mediawiki#L22). It 
also has helper function to spend these scripts. 




[tests-url]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml
[tests-badge]: https://github.com/catalogfi/blockchain/actions/workflows/test.yml/badge.svg?branch=master