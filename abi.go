package tredd

import "github.com/ethereum/go-ethereum/accounts/abi"

var (
	uint64ABIType = mustABIType("uint64")
	bytesABIType  = mustABIType("bytes")
	byte32ABIType = mustABIType("byte32")
)

var prefixedChunkArgTypes = abi.Arguments{
	{Type: uint64ABIType},
	{Type: bytesABIType},
}

var prefixedHashArgTypes = abi.Arguments{
	{Type: uint64ABIType},
	{Type: byte32ABIType},
}

var cryptArgTypes = abi.Arguments{
	{Type: byte32ABIType},
	{Type: uint64ABIType},
	{Type: uint64ABIType},
}

func mustABIType(name string) abi.Type {
	typ, err := abi.NewType(name, "", nil)
	if err != nil {
		panic(err)
	}
	return typ
}
