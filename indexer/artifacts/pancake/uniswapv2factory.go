package uniswap

import (
	"fmt"
	"math/big"

	web3 "github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
)

var (
	_ = big.NewInt
)

// UniswapV2Factory is a solidity contract
type UniswapV2Factory struct {
	c *contract.Contract
}

// DeployUniswapV2Factory deploys a new UniswapV2Factory contract
func DeployUniswapV2Factory(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiUniswapV2Factory, binUniswapV2Factory, args...)
}

// NewUniswapV2Factory creates a new instance of the contract at a specific address
func NewUniswapV2Factory(addr web3.Address, provider *jsonrpc.Client) *UniswapV2Factory {
	return &UniswapV2Factory{c: contract.NewContract(addr, abiUniswapV2Factory, provider)}
}

// Contract returns the contract object
func (a *UniswapV2Factory) Contract() *contract.Contract {
	return a.c
}

// calls

// AllPairs calls the allPairs method in the solidity contract
func (a *UniswapV2Factory) AllPairs(val0 *big.Int, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("allPairs", web3.EncodeBlock(block...), val0)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(web3.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// AllPairsLength calls the allPairsLength method in the solidity contract
func (a *UniswapV2Factory) AllPairsLength(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("allPairsLength", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// FeeTo calls the feeTo method in the solidity contract
func (a *UniswapV2Factory) FeeTo(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("feeTo", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(web3.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// FeeToSetter calls the feeToSetter method in the solidity contract
func (a *UniswapV2Factory) FeeToSetter(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("feeToSetter", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(web3.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// GetPair calls the getPair method in the solidity contract
func (a *UniswapV2Factory) GetPair(val0 web3.Address, val1 web3.Address, block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("getPair", web3.EncodeBlock(block...), val0, val1)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(web3.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}


// txns

// CreatePair sends a createPair transaction in the solidity contract
func (a *UniswapV2Factory) CreatePair(tokenA web3.Address, tokenB web3.Address) *contract.Txn {
	return a.c.Txn("createPair", tokenA, tokenB)
}

// SetFeeTo sends a setFeeTo transaction in the solidity contract
func (a *UniswapV2Factory) SetFeeTo(feeTo web3.Address) *contract.Txn {
	return a.c.Txn("setFeeTo", feeTo)
}

// SetFeeToSetter sends a setFeeToSetter transaction in the solidity contract
func (a *UniswapV2Factory) SetFeeToSetter(feeToSetter web3.Address) *contract.Txn {
	return a.c.Txn("setFeeToSetter", feeToSetter)
}
