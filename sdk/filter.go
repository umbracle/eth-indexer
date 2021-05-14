package sdk

import "github.com/umbracle/go-web3"

// filters of data

// exact address filter
type FilterByAddr struct {
	FromAddr   web3.Address
	ToAddr     web3.Address
	StartBlock uint64
}
