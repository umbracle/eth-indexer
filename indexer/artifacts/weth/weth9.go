package weth

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

// WETH9 is a solidity contract
type WETH9 struct {
	c *contract.Contract
}

// DeployWETH9 deploys a new WETH9 contract
func DeployWETH9(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiWETH9, binWETH9, args...)
}

// NewWETH9 creates a new instance of the contract at a specific address
func NewWETH9(addr web3.Address, provider *jsonrpc.Client) *WETH9 {
	return &WETH9{c: contract.NewContract(addr, abiWETH9, provider)}
}

// Contract returns the contract object
func (a *WETH9) Contract() *contract.Contract {
	return a.c
}

// calls

// Allowance calls the allowance method in the solidity contract
func (a *WETH9) Allowance(val0 web3.Address, val1 web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("allowance", web3.EncodeBlock(block...), val0, val1)
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

// BalanceOf calls the balanceOf method in the solidity contract
func (a *WETH9) BalanceOf(val0 web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("balanceOf", web3.EncodeBlock(block...), val0)
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

// Decimals calls the decimals method in the solidity contract
func (a *WETH9) Decimals(block ...web3.BlockNumber) (retval0 uint8, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("decimals", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(uint8)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// Name calls the name method in the solidity contract
func (a *WETH9) Name(block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("name", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// Symbol calls the symbol method in the solidity contract
func (a *WETH9) Symbol(block ...web3.BlockNumber) (retval0 string, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("symbol", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(string)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// TotalSupply calls the totalSupply method in the solidity contract
func (a *WETH9) TotalSupply(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("totalSupply", web3.EncodeBlock(block...))
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


// txns

// Approve sends a approve transaction in the solidity contract
func (a *WETH9) Approve(guy web3.Address, wad *big.Int) *contract.Txn {
	return a.c.Txn("approve", guy, wad)
}

// Deposit sends a deposit transaction in the solidity contract
func (a *WETH9) Deposit() *contract.Txn {
	return a.c.Txn("deposit")
}

// Transfer sends a transfer transaction in the solidity contract
func (a *WETH9) Transfer(dst web3.Address, wad *big.Int) *contract.Txn {
	return a.c.Txn("transfer", dst, wad)
}

// TransferFrom sends a transferFrom transaction in the solidity contract
func (a *WETH9) TransferFrom(src web3.Address, dst web3.Address, wad *big.Int) *contract.Txn {
	return a.c.Txn("transferFrom", src, dst, wad)
}

// Withdraw sends a withdraw transaction in the solidity contract
func (a *WETH9) Withdraw(wad *big.Int) *contract.Txn {
	return a.c.Txn("withdraw", wad)
}
