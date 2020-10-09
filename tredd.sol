pragma solidity ^0.7.2; // TODO: determine the best (lowest?) version number that works here.
pragma experimental ABIEncoderV2; // This is needed to compile the ProofStep[] params of refund().

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
  // TODO: How to specify the type of the desired token? (It's not "bytes.") Should we stick with just ETH to start with?
  bytes public mTokenType;

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
  constructor(address seller,
              uint amount,
              bytes memory tokenType,
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
  function reclaim() public {
    require (msg.sender == mBuyer);
    require (block.timestamp >= mRevealDeadline);
    require (!mRevealed);
    // TODO: transfer the balance in this contract to the buyer
    selfdestruct(msg.sender);
  }

  event evDecryptionKey(bytes32 decryptionKey);

  // The seller reveals the decryption key.
  function reveal(bytes32 decryptionKey) public {
    require (msg.sender == mSeller);
    require (block.timestamp < mRevealDeadline);
    require (!mRevealed);
    // TODO: require the seller to add (or have already added) the requested collateral

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
      bytes32 subkey = sha256(abi.encodePacked(mDecryptionKey, index, i));
      for (uint32 j = 0; j < 32 && i*32+j < chunk.length; j++) {
        output[i*32+j] = chunk[i*32+j] ^ subkey[j];
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
    require (block.timestamp < mRefundDeadline);
    require (mRevealed);

    //  1. Verify cipherProof w.r.t. Hash(index || cipherChunk) and mCipherRoot
    require (checkProof(cipherProof, sha256(abi.encodePacked(index, cipherChunk)), mCipherRoot)); // TODO: check abi.encodePacked(index, cipherChunk) exactly matches Go impl.

    //  2. Verify clearProof w.r.t. Hash(index || clearChunk) (given as clearHash) and mClearRoot
    require (checkProof(clearProof, clearHash, mClearRoot));

    //  3. Show Hash(index || decrypt(cipherChunk)) != Hash(index || clearChunk)
    require (sha256(abi.encodePacked(index, decrypt(cipherChunk, index))) != clearHash);

    //  4. TODO: Transfer the balance in this contract to the buyer.

    selfdestruct(msg.sender);
  }

  // The seller claims payment (and reclaims collateral) after the refund deadline.
  function claimPayment() public {
    require (msg.sender == mSeller);
    require (block.timestamp >= mRefundDeadline);
    // TODO: transfer the balance in this contract to the seller.
    selfdestruct(msg.sender);
  }
}
