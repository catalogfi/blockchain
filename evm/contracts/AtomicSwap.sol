// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.18;

/**
 * @dev Interface of the ERC20 standard as defined in the EIP.
 */
interface IERC20 {
    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );

    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `to`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address to, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(
        address owner,
        address spender
    ) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `from` to `to` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);
}

/**
 * @author  Catalog
 * @title   HTLC smart contract for atomic swaps
 * @notice  Any signer can create an order to serve as one of either halfs of an cross chain
 *          atomic swap.
 * @dev     The contracts can be used to create an order to serve as the the commitment for two
 *          types of users :
 *          Initiator functions: 1. initate
 *                               2. refund
 *          Redeemer funtions: 1. redeem
 */

contract AtomicSwap {
    IERC20 public immutable token;

    struct Order {
        address redeemer;
        address initiator;
        uint256 expiry;
        uint256 initiatedAt;
        uint256 amount;
        bool isFulfilled;
    }
    mapping(bytes32 => Order) public atomicSwapOrders;

    event Redeemed(
        bytes32 indexed orderId,
        bytes32 indexed secrectHash,
        bytes secret
    );
    event Initiated(
        bytes32 indexed orderId,
        bytes32 indexed secretHash,
        uint256 initiatedAt,
        uint256 amount
    );
    event Refunded(bytes32 indexed orderId);

    /**
     * @notice  .
     * @dev     provides checks to ensure
     *              1. redeemer is not null address
     *              2. redeemer is not same as the refunder
     *              3. expiry is greater than current block number
     *              4. amount is not zero
     * @param   redeemer  public address of the reedeem
     * @param   intiator  public address of the initator
     * @param   expiry  expiry in period for the htlc order
     * @param   amount  amount of tokens to trade
     */
    modifier checkSafe(
        address redeemer,
        address intiator,
        uint256 expiry,
        uint256 amount
    ) {
        require(redeemer != address(0), "AtomicSwap: invalid redeemer address");
        require(
            intiator != redeemer,
            "AtomicSwap: redeemer and initiator cannot be the same"
        );
        require(expiry > 0, "AtomicSwap: expiry should be greater than zero");
        require(amount > 0, "AtomicSwap: amount cannot be zero");
        _;
    }

    constructor(address _token) {
        token = IERC20(_token);
    }

    /**
     * @notice  Signers can create an order with order params
     * @dev     Secret used to generate secret hash for iniatiation should be generated randomly
     *          and sha256 hash should be used to support hashing methods on other non-evm chains.
     *          Signers cannot generate orders with same secret hash or override an existing order.
     * @param   _redeemer  public address of the redeemer
     * @param   _expiry  expiry in period for the htlc order
     * @param   _amount  amount of tokens to trade
     * @param   _secretHash  sha256 hash of the secret used for redemtion
     */
    function initiate(
        address _redeemer,
        uint256 _expiry,
        uint256 _amount,
        bytes32 _secretHash
    ) external checkSafe(_redeemer, msg.sender, _expiry, _amount) {
        bytes32 OrderId = sha256(abi.encode(_secretHash, msg.sender));
        Order memory order = atomicSwapOrders[OrderId];
        require(order.redeemer == address(0x0), "AtomicSwap: duplicate order");
        Order memory newOrder = Order({
            redeemer: _redeemer,
            initiator: msg.sender,
            expiry: _expiry,
            initiatedAt: block.number,
            amount: _amount,
            isFulfilled: false
        });
        atomicSwapOrders[OrderId] = newOrder;
        emit Initiated(
            OrderId,
            _secretHash,
            newOrder.initiatedAt,
            newOrder.amount
        );
        token.transferFrom(msg.sender, address(this), newOrder.amount);
    }

    /**
     * @notice  Signers with correct secret to an order's secret hash can redeem to claim the locked
     *          token
     * @dev     Signers are not allowed to redeem an order with wrong secret or redeem the same order
     *          multiple times
     * @param   _orderId  orderIds if the htlc order
     * @param   _secret  secret used to redeem an order
     */
    function redeem(bytes32 _orderId, bytes calldata _secret) external {
        Order storage order = atomicSwapOrders[_orderId];
        require(
            order.redeemer != address(0x0),
            "AtomicSwap: order not initated"
        );
        require(!order.isFulfilled, "AtomicSwap: order already fulfilled");
        bytes32 secretHash = sha256(_secret);
        require(
            sha256(abi.encode(secretHash, order.initiator)) == _orderId,
            "AtomicSwap: invalid secret"
        );
        order.isFulfilled = true;
        emit Redeemed(_orderId, secretHash, _secret);
        token.transfer(order.redeemer, order.amount);
    }

    /**
     * @notice  Signers can refund the locked assets after expiry block number
     * @dev     Signers cannot refund the an order before epiry block number or refund the same order
     *          multiple times
     * @param   _orderId  orderId of the htlc order
     */
    function refund(bytes32 _orderId) external {
        Order storage order = atomicSwapOrders[_orderId];
        require(
            order.redeemer != address(0x0),
            "AtomicSwap: order not initated"
        );
        require(!order.isFulfilled, "AtomicSwap: order already fulfilled");
        require(
            order.initiatedAt + order.expiry < block.number,
            "AtomicSwap: order not expired"
        );
        order.isFulfilled = true;
        emit Refunded(_orderId);
        token.transfer(order.initiator, order.amount);
    }
}
