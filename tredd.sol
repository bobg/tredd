pragma solidity ^0.7.2; // TODO: determine the best (lowest?) version number that works here.
pragma experimental ABIEncoderV2; // This is needed to compile the ProofStep[] params of refund().

// This interface is copied from
// https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v3.0.0/contracts/token/ERC20/IERC20.sol.
interface ERC20 {
    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `recipient`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(address owner, address spender) external view returns (uint256);

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
     * @dev Moves `amount` tokens from `sender` to `recipient` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

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
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

// A Tredd contract represents a single exchange of payment for data.
// It is deployed on-chain by the buyer,
// who funds it with the proposed payment
// and specifies the conditions a seller must meet in order to collect payment.
//
// Before deploying this contract,
// the buyer contacts the seller out of band to receive an encrypted copy of the data.
// A hash computed from the encrypted data is included when the buyer deploys the contract.
//
// The seller must supply the decryption key by a given deadline
// (otherwise the buyer can reclaim the funds in the contract).
// The buyer then has until a later deadline to claim a refund
// if it turns out that the decryption key was wrong,
// or the decrypted content is.
// After that deadline, the seller may collect payment from the contract.
contract Tredd {
  address public mBuyer;
  address public mSeller;

  // The type of the token the buyer proposes to pay with.
  ERC20 public mTokenType;

  // The amount that the buyer proposes to pay.
  // The contract must be funded with this amount.
  uint public mAmount;

  // The amount of collateral that the buyer requires the seller to add
  // when revealing the decryption key.
  // This is denominated in the same units as mAmount (namely, mTokenType).
  // The buyer collects this collateral in case of a successful refund claim,
  // and acts as a penalty for sellers supplying bad data
  // (who might otherwise be tempted to grief the network because it's free).
  uint public mCollateral;

  // Merkle root hash of the cleartext chunks of the desired content.
  bytes32 public mClearRoot;

  // Merkle root hash of the ciphertext chunks received.
  bytes32 public mCipherRoot;

  // Seller must reveal the decryption key by this time.
  int64 public mRevealDeadline;

  // Buyer must claim a refund by this time.
  int64 public mRefundDeadline;

  // The seller supplies this.
  bytes32 public mDecryptionKey;

  // False until the decryption key is revealed.
  bool public mRevealed;

  // Constructor.
  constructor(address seller,
              address tokenType,
              uint amount,
              uint collateral,
              bytes32 clearRoot,
              bytes32 cipherRoot,
              int64 revealDeadline,
              int64 refundDeadline) {
    mBuyer = msg.sender;
    mSeller = seller;
    mTokenType = ERC20(tokenType);
    mAmount = amount;
    mCollateral = collateral;
    mClearRoot = clearRoot;
    mCipherRoot = cipherRoot;
    mRevealDeadline = revealDeadline;
    mRefundDeadline = refundDeadline;
    mRevealed = false;
  }

  event evPaid();

  // The buyer adds their payment.
  // The buyer must first have approved the token transfer.
  // This emits the "paid" event,
  // which the seller watches for,
  // so the buyer should not pay into this contract by other means.
  function pay() public {
    require (msg.sender == mBuyer);
    require (block.timestamp < uint(mRevealDeadline));

    uint balance = mTokenType.balanceOf(address(this));
    require (balance < mAmount);
    require (mTokenType.transferFrom(mBuyer, address(this), mAmount - balance));

    emit evPaid();
  }

  // The reveal deadline has passed without reveal being called.
  // The buyer cancels the contract, reclaiming any payment made.
  function cancel() public {
    require (msg.sender == mBuyer);
    require (block.timestamp >= uint(mRevealDeadline));
    require (!mRevealed);

    uint balance = mTokenType.balanceOf(address(this));
    if (balance > 0) {
      mTokenType.transfer(mBuyer, balance);
    }

    selfdestruct(msg.sender);
  }

  event evDecryptionKey(bytes32 decryptionKey);

  // The seller reveals the decryption key.
  // Before calling this,
  // the seller must approve a transfer of the required collateral,
  // and should verify that the buyer has made their payment.
  function reveal(bytes32 decryptionKey) public {
    require (msg.sender == mSeller);
    require (block.timestamp < uint(mRevealDeadline));
    require (!mRevealed);

    // Add collateral.
    require (mTokenType.transferFrom(mSeller, address(this), mCollateral));

    mDecryptionKey = decryptionKey;
    mRevealed = true;

    emit evDecryptionKey(decryptionKey);
  }

  struct ProofStep {
    bytes h;
    bool left;
  }

  function checkProof(ProofStep[] memory steps, bytes32 leaf, bytes32 want) internal pure returns (bool) {
    bytes1 leafPrefix = '\x00';
    bytes1 interiorPrefix = '\x01';

    bytes32 got = sha256(abi.encodePacked(leafPrefix, leaf));

    for (uint32 i = 0; i < steps.length; i++) {
      ProofStep memory step = steps[i];
      if (step.left) {
        got = sha256(abi.encodePacked(interiorPrefix, step.h, got));
      } else {
        got = sha256(abi.encodePacked(interiorPrefix, got, step.h));
      }
    }

    return got == want;
  }

  function decrypt(bytes memory chunk, uint64 index) internal view returns (bytes memory) {
    bytes memory output = new bytes(chunk.length);
    for (uint64 i = 0; i*32 < chunk.length; i++) {
      uint64 pos = i*32;
      bytes32 subkey = sha256(abi.encodePacked(mDecryptionKey, index, i));
      for (uint32 j = 0; j < 32 && pos+j < chunk.length; j++) {
        output[pos+j] = chunk[pos+j] ^ subkey[j];
      }
    }
    return output;
  }

  // The buyer claims a refund by proving a chunk is wrong.
  // Args:
  //   - index: index of the chunk being proven wrong
  //   - cipherChunk: encrypted version of the chunk
  //   - clearHash: the expected value of Hash(index || clearChunk)
  //   - cipherProof: Merkle proof that cipherChunk was delivered before payment was proposed
  //   - clearProof: Merkle proof that the expected value of Hash(index || clearChunk) is in mClearRoot
  function refund(uint64 index,
                  bytes memory cipherChunk,
                  bytes32 clearHash,
                  ProofStep[] memory cipherProof,
                  ProofStep[] memory clearProof) public {
    require (msg.sender == mBuyer);
    require (block.timestamp < uint(mRefundDeadline));
    require (mRevealed);

    // 1. Verify cipherProof w.r.t. Hash(index || cipherChunk) and mCipherRoot
    require (checkProof(cipherProof, sha256(abi.encodePacked(index, cipherChunk)), mCipherRoot)); // TODO: check abi.encodePacked(index, cipherChunk) exactly matches Go impl.

    // 2. Verify clearProof w.r.t. Hash(index || clearChunk) (given as clearHash) and mClearRoot
    require (checkProof(clearProof, clearHash, mClearRoot));

    // 3. Show Hash(index || decrypt(cipherChunk)) != Hash(index || clearChunk)
    require (sha256(abi.encodePacked(index, decrypt(cipherChunk, index))) != clearHash);

    // 4. Transfer the balance in this contract to the buyer.
    require (mTokenType.transfer(mBuyer, mAmount+mCollateral));

    // 5. Self destruct.
    selfdestruct(msg.sender);
  }

  // The seller claims payment (and reclaims collateral) after the refund deadline.
  function claimPayment() public {
    require (msg.sender == mSeller);
    require (block.timestamp >= uint(mRefundDeadline));
    require (mTokenType.transfer(mSeller, mAmount+mCollateral));

    selfdestruct(msg.sender);
  }
}
