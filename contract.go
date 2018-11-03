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
	senderCollectsSrc  = fmt.Sprintf(senderCollectsFmt, standard.PayToMultisigProg2)
	senderCollectsProg = mustAssemble(senderCollectsSrc)
)

const merkleCheckSrc = `
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

var merkleCheckProg = mustAssemble(merkleCheckSrc)

const decryptSrc = `
                       #  con stack                                                      arg stack  log  notes
                       #  ---------                                                      ---------  ---  -----
                       #  contract stack                                                 arg stack  log  notes
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
dup                    #  key index subindex output msg subkey msglen msglen                             
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
