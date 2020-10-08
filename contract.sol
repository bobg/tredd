pragma solidity ^0.7.2; // TODO: determine the best (lowest?) version number that works here.

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

  // The amount that the buyer proposes to pay.
  // The contract must be funded with this amount.
  uint public mAmount;

  // The type of the token the buyer proposes to pay with.
  // TODO: How to specify the type of the desired token? Should we stick with just ETH to start with?
  TODO public mTokenType;

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
  uint public mRevealDeadline;

  // Buyer must claim a refund by this time.
  uint public mRefundDeadline;

  // The seller supplies this.
  bytes32 public mDecryptionKey;

  // False until the decryption key is revealed.
  bool public mRevealed;

  // Constructor.
  function Tredd(address seller,
                 uint amount,
                 TODO tokenType,
                 uint collateral,
                 bytes32 clearRoot,
                 bytes32 cipherRoot,
                 uint revealDeadline,
                 uint refundDeadline) {
    mBuyer = msg.sender;
    mSeller = seller;
    mAmount = amount;
    mTokenType = tokenType;
    mCollateral = collateral;
    mClearRoot = clearRoot;
    mCipherRoot = cipherRoot;
    mRevealDeadline = revealDeadline;
    mRefundDeadline = refundDeadline;
    mRevealed = false;
  }

  // The buyer may reclaim payment if it's after the reveal deadline and no decryption key has been revealed.
  function reclaim() {
    require (now >= mRevealDeadline);
    require (!mRevealed);
    // TODO: transfer the balance in this contract to the buyer and destroy the contract.
  }

  event evDecryptionKey(bytes32 decryptionKey);

  // The seller reveals the decryption key.
  function reveal(bytes32 decryptionKey) {
    require (now < mRevealDeadline);
    require (!mRevealed);
    // TODO: require the seller to add (or have already added) the requested collateral

    mDecryptionKey = decryptionKey;
    mRevealed = true;

    evDecryptionKey(decryptionKey);
  }

  // The buyer claims a refund by proving a chunk is wrong.
  // Args:
  //   - index: index of the chunk being proven wrong
  //   - cipherChunk: encrypted version of the chunk
  //   - clearHash: the expected value of Hash(index || clearChunk)
  //   - cipherProof: Merkle proof that cipherChunk was delivered before payment was proposed
  //   - clearProof: Merkle proof that the expected value of Hash(index || clearChunk) is in mClearRoot
  function refund(uint index,
                  bytes cipherChunk,
                  bytes32 clearHash,
                  TODO cipherProof, // How to express a Merkle proof?
                  TODO clearProof) {
    require (now < mRefundDeadline);
    require (mRevealed);

    // TODO: implement as follows:
    //  1. Verify cipherProof w.r.t. Hash(index || cipherChunk) and mCipherRoot
    //  2. Verify clearProof w.r.t. Hash(index || clearChunk) and mClearRoot
    //  3. Show Hash(index || decrypt(cipherChunk)) != Hash(index || clearChunk)
    //  4. Transfer the balance in this contract to the buyer and destroy the contract.
  }

  // The seller claims payment (and reclaims collateral) after the refund deadline.
  function claimPayment() {
    require (now >= mRefundDeadline);
    // TODO: transfer the balance in this contract to the seller and destroy the contract.
  }
}
