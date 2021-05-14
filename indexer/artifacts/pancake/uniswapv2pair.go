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

// UniswapV2Pair is a solidity contract
type UniswapV2Pair struct {
	c *contract.Contract
}

// DeployUniswapV2Pair deploys a new UniswapV2Pair contract
func DeployUniswapV2Pair(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiUniswapV2Pair, binUniswapV2Pair, args...)
}

// NewUniswapV2Pair creates a new instance of the contract at a specific address
func NewUniswapV2Pair(addr web3.Address, provider *jsonrpc.Client) *UniswapV2Pair {
	return &UniswapV2Pair{c: contract.NewContract(addr, abiUniswapV2Pair, provider)}
}

// Contract returns the contract object
func (a *UniswapV2Pair) Contract() *contract.Contract {
	return a.c
}

// calls

// DOMAINSEPARATOR calls the DOMAIN_SEPARATOR method in the solidity contract
func (a *UniswapV2Pair) DOMAINSEPARATOR(block ...web3.BlockNumber) (retval0 [32]byte, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("DOMAIN_SEPARATOR", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].([32]byte)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// MINIMUMLIQUIDITY calls the MINIMUM_LIQUIDITY method in the solidity contract
func (a *UniswapV2Pair) MINIMUMLIQUIDITY(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("MINIMUM_LIQUIDITY", web3.EncodeBlock(block...))
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

// PERMITTYPEHASH calls the PERMIT_TYPEHASH method in the solidity contract
func (a *UniswapV2Pair) PERMITTYPEHASH(block ...web3.BlockNumber) (retval0 [32]byte, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("PERMIT_TYPEHASH", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].([32]byte)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Allowance calls the allowance method in the solidity contract
func (a *UniswapV2Pair) Allowance(val0 web3.Address, val1 web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
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
func (a *UniswapV2Pair) BalanceOf(val0 web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
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
func (a *UniswapV2Pair) Decimals(block ...web3.BlockNumber) (retval0 uint8, err error) {
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

// Factory calls the factory method in the solidity contract
func (a *UniswapV2Pair) Factory(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("factory", web3.EncodeBlock(block...))
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

// GetReserves calls the getReserves method in the solidity contract
func (a *UniswapV2Pair) GetReserves(block ...web3.BlockNumber) (retval0 *big.Int, retval1 *big.Int, retval2 uint32, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("getReserves", web3.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["_reserve0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	retval1, ok = out["_reserve1"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 1")
		return
	}
	retval2, ok = out["_blockTimestampLast"].(uint32)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 2")
		return
	}

	return
}

// KLast calls the kLast method in the solidity contract
func (a *UniswapV2Pair) KLast(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("kLast", web3.EncodeBlock(block...))
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

// Name calls the name method in the solidity contract
func (a *UniswapV2Pair) Name(block ...web3.BlockNumber) (retval0 string, err error) {
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

// Nonces calls the nonces method in the solidity contract
func (a *UniswapV2Pair) Nonces(val0 web3.Address, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("nonces", web3.EncodeBlock(block...), val0)
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

// Price0CumulativeLast calls the price0CumulativeLast method in the solidity contract
func (a *UniswapV2Pair) Price0CumulativeLast(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("price0CumulativeLast", web3.EncodeBlock(block...))
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

// Price1CumulativeLast calls the price1CumulativeLast method in the solidity contract
func (a *UniswapV2Pair) Price1CumulativeLast(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("price1CumulativeLast", web3.EncodeBlock(block...))
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

// Symbol calls the symbol method in the solidity contract
func (a *UniswapV2Pair) Symbol(block ...web3.BlockNumber) (retval0 string, err error) {
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

// Token0 calls the token0 method in the solidity contract
func (a *UniswapV2Pair) Token0(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("token0", web3.EncodeBlock(block...))
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

// Token1 calls the token1 method in the solidity contract
func (a *UniswapV2Pair) Token1(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("token1", web3.EncodeBlock(block...))
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

// TotalSupply calls the totalSupply method in the solidity contract
func (a *UniswapV2Pair) TotalSupply(block ...web3.BlockNumber) (retval0 *big.Int, err error) {
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
func (a *UniswapV2Pair) Approve(spender web3.Address, value *big.Int) *contract.Txn {
	return a.c.Txn("approve", spender, value)
}

// Burn sends a burn transaction in the solidity contract
func (a *UniswapV2Pair) Burn(to web3.Address) *contract.Txn {
	return a.c.Txn("burn", to)
}

// Initialize sends a initialize transaction in the solidity contract
func (a *UniswapV2Pair) Initialize(token0 web3.Address, token1 web3.Address) *contract.Txn {
	return a.c.Txn("initialize", token0, token1)
}

// Mint sends a mint transaction in the solidity contract
func (a *UniswapV2Pair) Mint(to web3.Address) *contract.Txn {
	return a.c.Txn("mint", to)
}

// Permit sends a permit transaction in the solidity contract
func (a *UniswapV2Pair) Permit(owner web3.Address, spender web3.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) *contract.Txn {
	return a.c.Txn("permit", owner, spender, value, deadline, v, r, s)
}

// Skim sends a skim transaction in the solidity contract
func (a *UniswapV2Pair) Skim(to web3.Address) *contract.Txn {
	return a.c.Txn("skim", to)
}

// Swap sends a swap transaction in the solidity contract
func (a *UniswapV2Pair) Swap(amount0Out *big.Int, amount1Out *big.Int, to web3.Address, data []byte) *contract.Txn {
	return a.c.Txn("swap", amount0Out, amount1Out, to, data)
}

// Sync sends a sync transaction in the solidity contract
func (a *UniswapV2Pair) Sync() *contract.Txn {
	return a.c.Txn("sync")
}

// Transfer sends a transfer transaction in the solidity contract
func (a *UniswapV2Pair) Transfer(to web3.Address, value *big.Int) *contract.Txn {
	return a.c.Txn("transfer", to, value)
}

// TransferFrom sends a transferFrom transaction in the solidity contract
func (a *UniswapV2Pair) TransferFrom(from web3.Address, to web3.Address, value *big.Int) *contract.Txn {
	return a.c.Txn("transferFrom", from, to, value)
}
