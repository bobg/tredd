package tedd

import (
	"fmt"

	"github.com/chain/txvm/protocol/txbuilder/standard"
	"github.com/chain/txvm/protocol/txvm"
	"github.com/chain/txvm/protocol/txvm/asm"
)

// The Tedd contract, which holds funds in escrow and reveals the decryption key.
// The "seller reveals key" section is inlined as an argument to yield.
// The redemption clauses are inlined indirectly as part of that.
// The buyer constructs a partial transaction that calls this contract once, with its initial args.
// That partial contract ends with the yielded contract on the stack.
// The seller completes the transaction by adding another call of the contract with its own args.
// Since the buyer unlocks funds to use in this contract but does not publish the transaction,
// they should secure those funds by using a signature program checking the transaction log for the expected values (as described in the "log" comment column here).
const teddContractFmt = `
	                #  con stack                                          arg stack                                                         log                              notes
	                #  ---------                                          ---------                                                         ---                              -----
	                #                                                     payment clearRoot cipherRoot buyer refundDeadline revealDeadline
	0               #  0                                                  payment clearRoot cipherRoot buyer refundDeadline revealDeadline
	get             #  0 revealDeadline                                   payment clearRoot cipherRoot buyer refundDeadline
	timerange       #                                                     payment clearRoot cipherRoot buyer refundDeadline                 {'R', seed, 0, revealDeadline}   Seller must publish this transaction before revealDeadline.
	get dup log     #  refundDeadline                                     payment clearRoot cipherRoot buyer                                ... {'L', seed, refundDeadline}
	get dup log     #  refundDeadline buyer                               payment clearRoot cipherRoot                                      ... {'L', seed, buyer}
	get dup log     #  refundDeadline buyer cipherRoot                    payment clearRoot                                                 ... {'L', seed, cipherRoot}
	get dup log     #  refundDeadline buyer cipherRoot clearRoot          payment                                                           ... {'L', seed, clearRoot}
	get amount log  #  refundDeadline buyer cipherRoot clearRoot payment                                                                    ... {'L', seed, paymentAmount}
	assetid log     #  refundDeadline buyer cipherRoot clearRoot payment                                                                    ... {'L', seed, paymentAssetID}
	anchor log      #  refundDeadline buyer cipherRoot clearRoot payment                                                                    ... {'L', seed, paymentAnchor}   Buyer can recognize when seller has completed and published this transaction by observing this unique log entry.
	x'%x' yield     #                                                                                                                                                        "Seller reveals key" program goes here.
`

var (
	teddContractSrc  = fmt.Sprintf(teddContractFmt, sellerRevealsKeyProg)
	teddContractProg = mustAssemble(teddContractSrc)
	teddContractSeed = txvm.ContractSeed(teddContractProg)
)

// The second, post-yield phase of the Tedd contract,
// in which the seller reveals the decryption key.
const sellerRevealsKeyFmt = `
	                   #  con stack                                                                                              arg stack              log                            notes
	                   #  ---------                                                                                              ---------              ---                            -----
	                   #  refundDeadline buyer cipherRoot clearRoot payment                                                      seller key collateral
	amount 2 mul swap  #  refundDeadline buyer cipherRoot clearRoot 2*paymentAmount payment                                      seller key collateral
	get merge          #  refundDeadline buyer cipherRoot clearRoot 2*paymentAmount payment+collateral                           seller key                                            This will fail if collateral is of the wrong asset type.
	anchor log         #  refundDeadline buyer cipherRoot clearRoot 2*paymentAmount payment+collateral                           seller key             ... {'L', seed, mergedAnchor}  This is in addition to the log entries from "buyer proposes payment," which is the same transaction.
	amount swap        #  refundDeadline buyer cipherRoot clearRoot 2*paymentAmount payment+collateralAmount payment+collateral  seller key
	4 bury             #  refundDeadline buyer payment+collateral cipherRoot clearRoot 2*paymentAmount payment+collateralAmount  seller key
	eq verify          #  refundDeadline buyer payment+collateral cipherRoot clearRoot                                           seller key                                            Seller must put up collateral equal to the buyer's amount.
	get dup log        #  refundDeadline buyer payment+collateral cipherRoot clearRoot key                                       seller                 ... {'L', seed, key}
	get dup log        #  refundDeadline buyer payment+collateral cipherRoot clearRoot key seller                                                       ... {'L', seed, seller}
	x'%x' output       #                                                                                                                                ... {'O', seed, outputid}      Redemption program goes here.
`

var (
	sellerRevealsKeySrc  = fmt.Sprintf(sellerRevealsKeyFmt, redemptionProg)
	sellerRevealsKeyProg = mustAssemble(sellerRevealsKeySrc)
)

// Dispatcher for the redemption phase of the Tedd contract:
// either the seller claims payment (by putting a 0 on the stack)
// or the buyer claims a refund (by putting proofs, etc. followed by a 1 on the stack).
const redemptionFmt = `
	                     #  con stack                                                                         arg stack																								 notes
	                     #  ---------                                                                         ---------																								 -----
	                     #  refundDeadline buyer payment+collateral cipherRoot clearRoot key seller           ...args... selector
	get                  #  refundDeadline buyer payment+collateral cipherRoot clearRoot key seller selector  ...args...
	1 eq jumpif:$refund  #  refundDeadline buyer payment+collateral cipherRoot clearRoot key seller
	5 bury               #  refundDeadline seller buyer payment+collateral cipherRoot clearRoot key
	drop drop drop       #  refundDeadline seller buyer payment+collateral
	swap drop            #  refundDeadline seller payment+collateral
	x'%x' exec           #																																																																						 "Seller claims payment" program goes here.
	jump:$end            #
	$refund              #  refundDeadline buyer payment+collateral cipherRoot clearRoot key seller           cipherproof clearproof wantclearhash cipherchunk prefix
	drop                 #  refundDeadline buyer payment+collateral cipherRoot clearRoot key                  cipherproof clearproof wantclearhash cipherchunk prefix
	x'%x' exec           #																																																																						 "Buyer claims refund" program goes here.
	$end                 #
`

var (
	redemptionSrc  = fmt.Sprintf(redemptionFmt, sellerClaimsPaymentProg, buyerClaimsRefundProg)
	redemptionProg = mustAssemble(redemptionSrc)
)

// One of two Tedd-contract redemption clauses.
const sellerClaimsPaymentFmt = `
	                     #  con stack                                 arg stack                                    log                                 notes
	                     #  ---------                                 ---------                                    ---                                 -----
	                     #  refundDeadline seller payment+collateral
	splitzero put        #  refundDeadline seller payment+collateral  zeroval
	"" put "" put        #  refundDeadline seller payment+collateral  zeroval "" ""
	put                  #  refundDeadline seller                     zeroval "" "" payment+collateral
	1 tuple put          #  refundDeadline                            zeroval "" "" payment+collateral {seller}
	1 put                #  refundDeadline                            zeroval "" "" payment+collateral {seller} 1
	x'%x' contract call  #  refundDeadline                            zeroval                                      {'O', seed, outputID}
	0 timerange          #                                            zeroval                                      ... {'R', seed, refundDeadline, 0}  Check that the refund deadline has passed.
`

var (
	sellerClaimsPaymentSrc  = fmt.Sprintf(sellerClaimsPaymentFmt, standard.PayToMultisigProg2)
	sellerClaimsPaymentProg = mustAssemble(sellerClaimsPaymentSrc)
)

// One of two Tedd-contract redemption clauses.
const buyerClaimsRefundFmt = `
	                     #  con stack                                                                                                                        arg stack                                                log                                 notes
	                     #  ---------                                                                                                                        ---------                                                ---                                 -----
	                     #  refundDeadline buyer payment+collateral cipherRoot clearRoot key                                                                 cipherproof clearproof wantclearhash cipherchunk prefix
	get dup              #  refundDeadline buyer payment+collateral cipherRoot clearRoot key prefix prefix                                                   cipherproof clearproof wantclearhash cipherchunk
	2 bury               #  refundDeadline buyer payment+collateral cipherRoot clearRoot prefix key prefix                                                   cipherproof clearproof wantclearhash cipherchunk
	int                  #  refundDeadline buyer payment+collateral cipherRoot clearRoot prefix key prefixnum                                                cipherproof clearproof wantclearhash cipherchunk
	get dup              #  refundDeadline buyer payment+collateral cipherRoot clearRoot prefix key prefixnum cipherchunk cipherchunk                        cipherproof clearproof wantclearhash
	5 bury               #  refundDeadline buyer payment+collateral cipherRoot cipherchunk clearRoot prefix key prefixnum cipherchunk                        cipherproof clearproof wantclearhash
	x'%x' exec           #  refundDeadline buyer payment+collateral cipherRoot cipherchunk clearRoot prefix clearchunk                                       cipherproof clearproof wantclearhash                                                         The decrypt subroutine goes here.
	x'00'                #  refundDeadline buyer payment+collateral cipherRoot cipherchunk clearRoot prefix clearchunk x'00'                                 cipherproof clearproof wantclearhash
	2 roll dup dup       #  refundDeadline buyer payment+collateral cipherRoot cipherchunk clearRoot clearchunk x'00' prefix prefix prefix                   cipherproof clearproof wantclearhash
	4 bury               #  refundDeadline buyer payment+collateral cipherRoot cipherchunk clearRoot prefix clearchunk x'00' prefix prefix                   cipherproof clearproof wantclearhash
	6 bury               #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix clearchunk x'00' prefix                   cipherproof clearproof wantclearhash
	cat                  #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix clearchunk x'00'+prefix                   cipherproof clearproof wantclearhash
	swap cat             #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix x'00'+prefix+clearchunk                   cipherproof clearproof wantclearhash
	sha256               #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix gotclearhash                              cipherproof clearproof wantclearhash
	get dup              #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix gotclearhash wantclearhash wantclearhash  cipherproof clearproof
	2 bury               #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix wantclearhash gotclearhash wantclearhash  cipherproof clearproof
	eq not verify        #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix wantclearhash                             cipherproof clearproof                                                                       Show hash(decrypt(key, chunk)) != wantclearhash.
	cat                  #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot prefix+wantclearhash                             cipherproof clearproof
	get swap             #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk clearRoot clearproof prefix+wantclearhash                  cipherproof
	x'%x' exec           #  refundDeadline buyer payment+collateral cipherRoot prefix cipherchunk                                                            cipherproof                                                                                  Check merkle proof subroutine goes here. This shows that wantclearhash is the right clear hash for the given prefix.
	cat                  #  refundDeadline buyer payment+collateral cipherRoot prefix+cipherchunk                                                            cipherproof
	get swap             #  refundDeadline buyer payment+collateral cipherRoot cipherproof prefix+cipherchunk
	x'%x' exec           #  refundDeadline buyer payment+collateral                                                                                                                                                                                       Check merkle proof subroutine goes here again. This shows that cipherchunk, with the same prefix, is the right chunk.
	splitzero put        #  refundDeadline buyer payment+collateral                                                                                          zeroval
	"" put "" put        #  refundDeadline buyer payment+collateral                                                                                          zeroval "" ""
	put                  #  refundDeadline buyer                                                                                                             zeroval "" "" payment+collateral
	1 tuple put          #  refundDeadline                                                                                                                   zeroval "" "" payment+collateral {buyer}
	1 put                #  refundDeadline                                                                                                                   zeroval "" "" payment+collateral {buyer} 1
	x'%x' contract call  #  refundDeadline                                                                                                                   zeroval                                                  {'O', seed, outputID}
	0 swap timerange     #                                                                                                                                   zeroval                                                  ... {'L', seed, 0, refundDeadline}
`

var (
	buyerClaimsRefundSrc  = fmt.Sprintf(buyerClaimsRefundFmt, decryptProg, merkleCheckProg, merkleCheckProg, standard.PayToMultisigProg2)
	buyerClaimsRefundProg = mustAssemble(buyerClaimsRefundSrc)
)

// Subroutine to check a Merkle proof + leaf value against an expected root hash.
const merkleCheckSrc = `
	                    #  con stack																							notes
	                    #  ---------																							-----
	                    #  wanthash {hash isleft hash isleft ...} leaf						A merkle proof is a list of hash/isleft pairs, supplied here as a flat tuple. When checking, pairs are consumed from right to left. Note, this is the opposite of the order in github.com/bobg/merkle.Proof.
	x'00' swap          #  wanthash {hash isleft hash isleft ...} leaf x'00'
	cat sha256          #  wanthash {hash isleft hash isleft ...} gothash					Leaf hashes are made by prepending x'00' to the leaf value. See merkle.LeafHash.
	swap                #  wanthash gothash {hash isleft hash isleft ...}
	untuple             #  wanthash gothash hash isleft hash isleft ... 2N
	$loop               #
	dup 0 eq            #  wanthash gothash hash isleft hash isleft ... 2N 2N==0
	jumpif:$check       #  wanthash gothash hash isleft hash isleft ... 2N
	                    #
	dup 1 add           #  wanthash gothash hash isleft hash isleft ... 2N 2N+1
	roll                #  wanthash hash isleft hash isleft ... 2N gothash
	3 roll              #  wanthash hash isleft isleft ... 2N gothash hash
	3 roll              #  wanthash hash isleft ... 2N gothash hash isleft
	not jumpif:$dohash  #  wanthash hash isleft ... 2N gothash hash
	swap                #  wanthash hash isleft ... 2N hash gothash
	$dohash             #
	cat                 #  wanthash hash isleft ... 2N combined
	x'01' swap          #  wanthash hash isleft ... 2N x'01' combined							Interior hashes are made by prepending x'01' to the concatenated left and right subhashes.
	cat sha256          #  wanthash hash isleft ... 2N gothash
	swap                #  wanthash hash isleft ... gothash 2N
	2 sub               #  wanthash hash isleft ... gothash 2N-2
	dup                 #  wanthash hash isleft ... gothash 2N-2 2N-2
	2 roll              #  wanthash hash isleft ... 2N-2 2N-2 gothash
	swap                #  wanthash hash isleft ... 2N-2 gothash 2N-2
	1 add               #  wanthash hash isleft ... 2N-2 gothash 2N-1
	bury                #  wanthash gothash hash isleft ... 2N-2
	jump:$loop          #
	$check              #  wanthash gothash 0
	drop                #  wanthash gothash
	eq verify           #
`

var merkleCheckProg = mustAssemble(merkleCheckSrc)

// Subroutine to decrypt a chunk by xoring with hashes derived from a key.
const decryptSrc = `
	                       #  con stack
	                       #  ---------
	                       #  key index msg
	0 ''                   #  key index msg 0 ''
	2 roll                 #  key index 0 '' msg
	$loop                  #
	dup len dup            #  key index subindex output msg msglen msglen
	0 eq                   #  key index subindex output msg msglen msglen==0
	jumpif:$cleanup1       #  key index subindex output msg msglen
	5 roll                 #  index subindex output msg msglen key
	dup                    #  index subindex output msg msglen key key
	6 bury                 #  key index subindex output msg msglen key
	5 roll                 #  key subindex output msg msglen key index
	dup                    #  key subindex output msg msglen key index index
	6 bury                 #  key index subindex output msg msglen key index
	encode                 #  key index subindex output msg msglen key indexstr
	cat                    #  key index subindex output msg msglen key+indexstr
	4 roll                 #  key index output msg msglen key+indexstr subindex
	dup                    #  key index output msg msglen key+indexstr subindex subindex
	5 bury                 #  key index subindex output msg msglen key+indexstr subindex
	encode                 #  key index subindex output msg msglen key+indexstr subindexstr
	cat                    #  key index subindex output msg msglen key+indexstr+subindexstr
	sha256                 #  key index subindex output msg msglen subkey
	swap                   #  key index subindex output msg subkey msglen
	dup                    #  key index subindex output msg subkey msglen msglen
	32                     #  key index subindex output msg subkey msglen msglen 32
	lt                     #  key index subindex output msg subkey msglen msglen<32
	jumpif:$finalsubchunk  #  key index subindex output msg subkey msglen
	2 roll                 #  key index subindex output subkey msglen msg
	dup                    #  key index subindex output subkey msglen msg msg
	0 32                   #  key index subindex output subkey msglen msg msg 0 32
	slice                  #  key index subindex output subkey msglen msg msg[:32]
	3 roll                 #  key index subindex output msglen msg msg[:32] subkey
	bitxor                 #  key index subindex output msglen msg decryptedsubchunk
	3 roll                 #  key index subindex msglen msg decryptedsubchunk output
	swap cat               #  key index subindex msglen msg output'
	2 bury                 #  key index subindex output' msglen msg
	32                     #  key index subindex output' msglen msg 32
	2 roll                 #  key index subindex output' msg 32 msglen
	slice                  #  key index subindex output' msg[32:]
	2 roll                 #  key index output' msg[32:] subindex
	1 add                  #  key index output' msg[32:] subindex+1
	2 bury                 #  key index subindex+1 output' msg[32:]
	jump:$loop             #
	$finalsubchunk         #  key index subindex output msg subkey msglen
	0 swap                 #  key index subindex output msg subkey 0 msglen
	slice                  #  key index subindex output msg subkey[:msglen]
	bitxor                 #  key index subindex output decryptedfinalsubchunk
	cat                    #  key index subindex output'
	jump:$cleanup2         #
	$cleanup1              #  key index subindex output msg msglen
	drop drop              #  key index subindex output
	$cleanup2              #  key index subindex output
	3 bury                 #  output key index subindex
	drop drop drop         #  output
`

var decryptProg = mustAssemble(decryptSrc)

func mustAssemble(s string) []byte {
	prog, err := asm.Assemble(s)
	if err != nil {
		panic(err)
	}
	return prog
}
