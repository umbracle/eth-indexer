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

// UniswapRouter is a solidity contract
type UniswapRouter struct {
	c *contract.Contract
}

// DeployUniswapRouter deploys a new UniswapRouter contract
func DeployUniswapRouter(provider *jsonrpc.Client, from web3.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiUniswapRouter, binUniswapRouter, args...)
}

// NewUniswapRouter creates a new instance of the contract at a specific address
func NewUniswapRouter(addr web3.Address, provider *jsonrpc.Client) *UniswapRouter {
	return &UniswapRouter{c: contract.NewContract(addr, abiUniswapRouter, provider)}
}

// Contract returns the contract object
func (a *UniswapRouter) Contract() *contract.Contract {
	return a.c
}

// calls

// WETH calls the WETH method in the solidity contract
func (a *UniswapRouter) WETH(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("WETH", web3.EncodeBlock(block...))
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

// Factory calls the factory method in the solidity contract
func (a *UniswapRouter) Factory(block ...web3.BlockNumber) (retval0 web3.Address, err error) {
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

// GetAmountIn calls the getAmountIn method in the solidity contract
func (a *UniswapRouter) GetAmountIn(amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("getAmountIn", web3.EncodeBlock(block...), amountOut, reserveIn, reserveOut)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["amountIn"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// GetAmountOut calls the getAmountOut method in the solidity contract
func (a *UniswapRouter) GetAmountOut(amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("getAmountOut", web3.EncodeBlock(block...), amountIn, reserveIn, reserveOut)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["amountOut"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// GetAmountsIn calls the getAmountsIn method in the solidity contract
func (a *UniswapRouter) GetAmountsIn(amountOut *big.Int, path []web3.Address, block ...web3.BlockNumber) (retval0 []*big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("getAmountsIn", web3.EncodeBlock(block...), amountOut, path)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["amounts"].([]*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// GetAmountsOut calls the getAmountsOut method in the solidity contract
func (a *UniswapRouter) GetAmountsOut(amountIn *big.Int, path []web3.Address, block ...web3.BlockNumber) (retval0 []*big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("getAmountsOut", web3.EncodeBlock(block...), amountIn, path)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["amounts"].([]*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}

// Quote calls the quote method in the solidity contract
func (a *UniswapRouter) Quote(amountA *big.Int, reserveA *big.Int, reserveB *big.Int, block ...web3.BlockNumber) (retval0 *big.Int, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = a.c.Call("quote", web3.EncodeBlock(block...), amountA, reserveA, reserveB)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["amountB"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	
	return
}


// txns

// AddLiquidity sends a addLiquidity transaction in the solidity contract
func (a *UniswapRouter) AddLiquidity(tokenA web3.Address, tokenB web3.Address, amountADesired *big.Int, amountBDesired *big.Int, amountAMin *big.Int, amountBMin *big.Int, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("addLiquidity", tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, to, deadline)
}

// AddLiquidityETH sends a addLiquidityETH transaction in the solidity contract
func (a *UniswapRouter) AddLiquidityETH(token web3.Address, amountTokenDesired *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("addLiquidityETH", token, amountTokenDesired, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidity sends a removeLiquidity transaction in the solidity contract
func (a *UniswapRouter) RemoveLiquidity(tokenA web3.Address, tokenB web3.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("removeLiquidity", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidityETH sends a removeLiquidityETH transaction in the solidity contract
func (a *UniswapRouter) RemoveLiquidityETH(token web3.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("removeLiquidityETH", token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHSupportingFeeOnTransferTokens sends a removeLiquidityETHSupportingFeeOnTransferTokens transaction in the solidity contract
func (a *UniswapRouter) RemoveLiquidityETHSupportingFeeOnTransferTokens(token web3.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("removeLiquidityETHSupportingFeeOnTransferTokens", token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHWithPermit sends a removeLiquidityETHWithPermit transaction in the solidity contract
func (a *UniswapRouter) RemoveLiquidityETHWithPermit(token web3.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to web3.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) *contract.Txn {
	return a.c.Txn("removeLiquidityETHWithPermit", token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens sends a removeLiquidityETHWithPermitSupportingFeeOnTransferTokens transaction in the solidity contract
func (a *UniswapRouter) RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(token web3.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to web3.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) *contract.Txn {
	return a.c.Txn("removeLiquidityETHWithPermitSupportingFeeOnTransferTokens", token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityWithPermit sends a removeLiquidityWithPermit transaction in the solidity contract
func (a *UniswapRouter) RemoveLiquidityWithPermit(tokenA web3.Address, tokenB web3.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to web3.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) *contract.Txn {
	return a.c.Txn("removeLiquidityWithPermit", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline, approveMax, v, r, s)
}

// SwapETHForExactTokens sends a swapETHForExactTokens transaction in the solidity contract
func (a *UniswapRouter) SwapETHForExactTokens(amountOut *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapETHForExactTokens", amountOut, path, to, deadline)
}

// SwapExactETHForTokens sends a swapExactETHForTokens transaction in the solidity contract
func (a *UniswapRouter) SwapExactETHForTokens(amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapExactETHForTokens", amountOutMin, path, to, deadline)
}

// SwapExactETHForTokensSupportingFeeOnTransferTokens sends a swapExactETHForTokensSupportingFeeOnTransferTokens transaction in the solidity contract
func (a *UniswapRouter) SwapExactETHForTokensSupportingFeeOnTransferTokens(amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapExactETHForTokensSupportingFeeOnTransferTokens", amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH sends a swapExactTokensForETH transaction in the solidity contract
func (a *UniswapRouter) SwapExactTokensForETH(amountIn *big.Int, amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapExactTokensForETH", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETHSupportingFeeOnTransferTokens sends a swapExactTokensForETHSupportingFeeOnTransferTokens transaction in the solidity contract
func (a *UniswapRouter) SwapExactTokensForETHSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapExactTokensForETHSupportingFeeOnTransferTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens sends a swapExactTokensForTokens transaction in the solidity contract
func (a *UniswapRouter) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapExactTokensForTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokensSupportingFeeOnTransferTokens sends a swapExactTokensForTokensSupportingFeeOnTransferTokens transaction in the solidity contract
func (a *UniswapRouter) SwapExactTokensForTokensSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapExactTokensForTokensSupportingFeeOnTransferTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapTokensForExactETH sends a swapTokensForExactETH transaction in the solidity contract
func (a *UniswapRouter) SwapTokensForExactETH(amountOut *big.Int, amountInMax *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapTokensForExactETH", amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactTokens sends a swapTokensForExactTokens transaction in the solidity contract
func (a *UniswapRouter) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []web3.Address, to web3.Address, deadline *big.Int) *contract.Txn {
	return a.c.Txn("swapTokensForExactTokens", amountOut, amountInMax, path, to, deadline)
}
