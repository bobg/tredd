//go:generate abigen --abi erc20.abi --pkg contract --type ERC20 --out erc20.go
//go:generate abigen --abi tredd.abi --bin tredd.bin --pkg contract --type Tredd --out tredd.go
package contract
