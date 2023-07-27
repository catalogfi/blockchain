This package contains all the golang bindings of the contracts defined in https://github.com/catalogfi/contracts-solidity/tree/main/contracts/instant

To generate the bindings, we need to 
1. Compile the contracts and get the ABI and bytecodes
```bash
## When you're in the root of the `contracts-solidity` repo 
solc --metadata-hash none --abi contracts/instant/InstantWalletFactory.sol --bin contracts/instant/InstantWalletFactory.sol -o build --include-path "$(pwd)/node_modules" --optimize --base-path .

## Generate files for the Entrypoint.sol contract 
solc --metadata-hash none --abi node_modules/@account-abstraction/contracts/core/EntryPoint.sol --bin node_modules/@account-abstraction/contracts/core/EntryPoint.sol -o build --optimize --base-path node_modules/@account-abstraction/contracts --overwrite
```
> You'll need to double check the version of `solc` is same to the compiler for deployment, otherwise you'll get different bytecode.
After running the above commands, you should have a folder called `build` which contains the `bin` and `abi` files for those contracts. We have copied the ones we're interested to this folder.
These files only need to be regenerated when we make changes to the contract code or we redepoly/update the contracts.

> If you're using some tools to deploy the contract (hardhat or truffle), you might need to manually copy the bytecode to the bin file.
since they're not always the same as `solc`

2. Generate go binding files 
```shell 
## when you're inside this folder (same as this README file)
abigen --abi artifacts/EntryPoint.abi --bin artifacts/EntryPoint.bin --pkg binding --type EntryPoint --out EntryPoint.go
abigen --abi artifacts/ERC1967Proxy.abi --bin artifacts/ERC1967Proxy.bin --pkg binding --type ERC1967Proxy --out ERC1967Proxy.go
abigen --abi artifacts/InstantWallet.abi --bin artifacts/InstantWallet.bin --pkg binding --type InstantWallet --out InstantWallet.go
abigen --abi artifacts/InstantWalletFactory.abi --bin artifacts/InstantWalletFactory.bin --pkg binding --type InstantWalletFactory --out InstatnWalletFactory.go
```

3. There might be some duplicated structs generated in different go files. There's currently no good way to avoid this, so I would 
recommend to remove them manually. (Hopefully we only have a few of this)

