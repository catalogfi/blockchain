package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// ContractAddresses returns the deployed contract address on the given network. The returned addresses will be in the
// order of Entrypoint, WalletFactory and WalletImplementation. `-1` version means the latest version.
func ContractAddresses(network Network, version int) (common.Address, common.Address, common.Address, error) {
	switch network {
	case NetworkSepolia:
		switch version {
		case -1, 1:
			return common.HexToAddress("0xE3d938fe81f85Bc3dC924E03B4206f3e58A0eED9"),
				common.HexToAddress("0xA86B9a9c7C67644Be13B1C66EAcD7300EAE6Bcf6"),
				common.HexToAddress("0x864ACB78a5C3EE95B3a4BDf3F0cdB77F2911e620"), nil
		default:
			return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("unknown version")
		}
	default:
		return common.Address{}, common.Address{}, common.Address{}, fmt.Errorf("unknown network")
	}
}
