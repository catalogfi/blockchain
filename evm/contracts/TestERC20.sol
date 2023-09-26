// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.18;
import {ERC20} from "./ERC20.sol";

contract TestToken is ERC20 {
    uint8 internal _decimals;

    constructor(
        string memory name_,
        string memory symbol_,
        uint8 decimals_,
        uint256 totalSupply_,
        address totalSupplyRecipient
    ) ERC20(name_, symbol_) {
        _decimals = decimals_;
        ERC20._mint(totalSupplyRecipient, totalSupply_);
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
    }
}

contract TestERC20 is TestToken {
    constructor() TestToken("ABC", "ABC", 6, 1e27, msg.sender) {}
}
