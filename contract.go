package tedd

import (
	"fmt"

	"github.com/chain/txvm/protocol/txbuilder/standard"
	"github.com/chain/txvm/protocol/txvm/asm"
)

const senderCollectsFmt = `
                     #  con stack                                arg stack                           log                                 notes
                     #  ---------                                ---------                           ---                                 -----
                     #  refundDeadline buyer payment+collateral                                                                          
"" put "" put        #  refundDeadline buyer payment+collateral  "" ""                                                                   
put                  #  refundDeadline buyer                     "" "" payment+collateral                                                
1 tuple put          #  refundDeadline                           "" "" payment+collateral {buyer}                                        
1 put                #  refundDeadline                           "" "" payment+collateral {buyer} 1                                      
x'%x' contract call  #  refundDeadline                                                               {'O', seed, outputID}               
0 timerange          #                                                                               ... {'R', seed, refundDeadline, 0}  
`

var (
	senderCollectsStr = fmt.Sprintf(senderCollectsFmt, standard.PayToMultisigProg2)
	senderCollects    = mustAssemble(senderCollectsStr)
)

const merkleCheckStr = `
                    #  con stack                                              arg stack  log  notes
                    #  ---------                                              ---------  ---  -----
                    #  wanthash {hash isleft hash isleft ...} leaf                            A merkle proof is a list of hash/isleft pairs, supplied here as a flat tuple. When checking, pairs are consumed from right to left. Note, this is the opposite of the order in github.com/bobg/merkle.Proof.
x'00' swap          #  wanthash {hash isleft hash isleft ...} leaf x'00'                      
cat sha256          #  wanthash {hash isleft hash isleft ...} gothash                         leaf hashes are made by prepending x'00' to the leaf value
swap                #  wanthash gothash {hash isleft hash isleft ...}                         
untuple             #  wanthash gothash hash isleft hash isleft ... 2N                        
$loop               #                                                                         
dup 0 eq            #  wanthash gothash hash isleft hash isleft ... 2N 2N==0                  
jumpif:$check       #  wanthash gothash hash isleft hash isleft ... 2N                        
dup 1 add           #  wanthash gothash hash isleft hash isleft ... 2N 2N+1                   
roll                #  wanthash hash isleft hash isleft ... 2N gothash                        
3 roll              #  wanthash hash isleft isleft ... 2N gothash hash                        
3 roll              #  wanthash hash isleft ... 2N gothash hash isleft                        
not jumpif:$dohash  #  wanthash hash isleft ... 2N gothash hash                               
swap                #  wanthash hash isleft ... 2N hash gothash                               
$dohash             #                                                                         
cat                 #  wanthash hash isleft ... 2N combined                                   
x'01' swap          #  wanthash hash isleft ... 2N x'01' combined                             interior hashes are made by prepending x'01' to the concatenated left and right subhashes
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

var merkleCheck = mustAssemble(merkleCheckStr)

func mustAssemble(s string) []byte {
	prog, err := asm.Assemble(s)
	if err != nil {
		panic(err)
	}
	return prog
}
