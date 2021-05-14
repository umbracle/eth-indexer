package pancake

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umbracle/eth-indexer/indexer/artifacts/erc20"
	uniswap "github.com/umbracle/eth-indexer/indexer/artifacts/pancake"
	"github.com/umbracle/eth-indexer/indexer/proto"
	"github.com/umbracle/eth-indexer/sdk"
	protosdk "github.com/umbracle/eth-indexer/sdk/proto"
	"github.com/umbracle/eth-indexer/testutil"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/contract"
)

func bigNumberify(i uint64) *big.Int {
	b := new(big.Int).SetUint64(i)
	return b
}

func TestPancakeMint(t *testing.T) {
	srv := testutil.NewTestServer(t)

	h := &Harness{
		indexer: Provider(),
		t:       t,
		srv:     srv,
	}
	h.setup()

	token0Amount := expandTo18decimals(1)
	token1Amount := expandTo18decimals(4)

	h.transferToken(h.token0, h.pairAddr, token0Amount)
	h.transferToken(h.token1, h.pairAddr, token1Amount)

	//fmt.Println("-- pair --")
	//fmt.Println(h.pair)

	// expectedLiquidity := expandTo18decimals(2)

	// Process stuff
	diff := h.Process(h.SendTxn(h.pair.Mint(h.srv.Owner())))

	fmt.Println("-- diff --")
	fmt.Println(diff)

	//fmt.Println(expectedLiquidity)
	//fmt.Println(h.pair.TotalSupply())
	//fmt.Println(h.pair.BalanceOf(h.srv.Owner()))
}

func TestPancakeSwapToken0(t *testing.T) {
	srv := testutil.NewTestServer(t)

	h := &Harness{
		indexer: Provider(),
		t:       t,
		srv:     srv,
	}
	h.setup()

	token0Amount := expandTo18decimals(5)
	token1Amount := expandTo18decimals(10)
	h.addliquidity(token0Amount, token1Amount)

	swapAmount := expandTo18decimals(1)
	h.transferToken(h.token0, h.pairAddr, swapAmount)

	expectedOutputAmount, _ := new(big.Int).SetString("1662497915624478906", 10)
	diff := h.Process(h.SendTxn(h.pair.Swap(big.NewInt(0), expectedOutputAmount, h.srv.Owner(), []byte{})))

	fmt.Println("-- diff --")
	fmt.Println(diff)
	fmt.Println(h.pair.GetReserves())
}

var minimumLiquidity = big.NewInt(1000)

func TestPancakeSwapToken1(t *testing.T) {
	// TODO
}

func TestPancakeBurn(t *testing.T) {
	srv := testutil.NewTestServer(t)

	h := &Harness{
		t:   t,
		srv: srv,
	}
	h.setup()

	token0Amount := expandTo18decimals(3)
	token1Amount := expandTo18decimals(3)
	h.addliquidity(token0Amount, token1Amount)

	expectedLiquidity := expandTo18decimals(3)
	h.SendTxn(h.pair.Transfer(h.pairAddr, new(big.Int).Sub(expectedLiquidity, minimumLiquidity)))
	h.SendTxn(h.pair.Burn(srv.Owner()))

	fmt.Println(h.pair.BalanceOf(h.srv.Owner()))
	fmt.Println(h.pair.TotalSupply())
	fmt.Println(h.token0.BalanceOf(h.pairAddr))
	fmt.Println(h.token1.BalanceOf(h.pairAddr))
}

func TestPancakeFee(t *testing.T) {
	srv := testutil.NewTestServer(t)

	h := &Harness{
		t:   t,
		srv: srv,
	}
	h.setup()

	token0Amount := expandTo18decimals(1000)
	token1Amount := expandTo18decimals(1000)
	h.addliquidity(token0Amount, token1Amount)

	// swap
	swapAmount := expandTo18decimals(1)
	expectedOutputAmount := bigNumberify(996006981039903216)

	h.transferToken(h.token1, h.pairAddr, swapAmount)
	txn := h.SendTxn(h.pair.Swap(expectedOutputAmount, big.NewInt(0), srv.Owner(), []byte{}))

	data, err := json.Marshal(txn.Receipt().Logs)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

}

func expandTo18decimals(i uint64) *big.Int {
	num := strconv.Itoa(int(i))
	for j := 0; j < 18; j++ {
		num += "0"
	}
	b, ok := new(big.Int).SetString(num, 10)
	if !ok {
		panic("BUG")
	}
	return b
}

type Harness struct {
	indexer *sdk.Provider

	t   *testing.T
	srv *testutil.TestServer

	factory *uniswap.UniswapV2Factory
	token0  *erc20.ERC20
	token1  *erc20.ERC20

	pairAddr web3.Address
	pair     *uniswap.UniswapV2Pair
}

func (h *Harness) Process(txn *contract.Txn) []*protosdk.Diff {
	rr := txn.Receipt()

	fmt.Println("- process event -")
	fmt.Println(rr.Logs)

	events := []proto.Event{}
	for _, l := range rr.Logs {
		events = append(events, *proto.DecodeEvent(l))

		fmt.Println("-- log --")
		fmt.Println(l.Topics[0].String())
	}

	act := &sdk.Action{
		BlockNum: txn.Receipt().BlockNumber,
		Events:   events,
	}
	diff, err := h.indexer.Process(act)

	fmt.Println("-- err --")
	fmt.Println(err)

	return diff
}

func (h *Harness) SendTxn(txn *contract.Txn) *contract.Txn {
	txn.SetGasLimit(10000000)
	assert.NoError(h.t, txn.DoAndWait())
	return txn
}

func (h *Harness) setup() {
	h.indexer.Init()
	h.indexer.SetClient(h.srv.Provider())

	h.factory = h.deployFactory()
	h.token0 = h.deployErc20("Gold", "GLD")
	h.token1 = h.deployErc20("Silver", "SLV")
	h.pair = h.deployPair()
}

func (h *Harness) validateAddr(txn *web3.Receipt) {
	code, err := h.srv.Provider().Eth().GetCode(txn.ContractAddress)
	assert.NoError(h.t, err)

	if len(code) < 10 {
		fmt.Println(h.srv.Http())
		fmt.Println(txn.TransactionHash.String())
		fmt.Println(txn.ContractAddress)
		fmt.Println(code)
		panic("bad")
	}
}

func (h *Harness) transferToken(token *erc20.ERC20, to web3.Address, amount *big.Int) {
	h.SendTxn(token.Transfer(to, amount))
}

func (h *Harness) deployErc20(name, symbol string) *erc20.ERC20 {
	srv := h.srv

	txn := erc20.DeployERC20(srv.Provider(), srv.Owner()).AddArgs(name, symbol)
	h.SendTxn(txn)
	h.validateAddr(txn.Receipt())

	token := erc20.NewERC20(txn.Receipt().ContractAddress, srv.Provider())
	token.Contract().SetFrom(srv.Owner())

	// mint big amount
	h.SendTxn(token.Mint(srv.Owner(), expandTo18decimals(10000)))

	return token
}

func (h *Harness) deployFactory() *uniswap.UniswapV2Factory {
	srv := h.srv

	txn := uniswap.DeployUniswapV2Factory(srv.Provider(), srv.Owner()).AddArgs(srv.Owner())
	h.SendTxn(txn)
	h.validateAddr(txn.Receipt())

	factory := uniswap.NewUniswapV2Factory(txn.Receipt().ContractAddress, srv.Provider())
	factory.Contract().SetFrom(srv.Owner())

	return factory
}

func (h *Harness) deployPair() *uniswap.UniswapV2Pair {
	srv := h.srv

	txn := h.factory.CreatePair(h.token0.Contract().Addr(), h.token1.Contract().Addr())
	h.SendTxn(txn)

	// populate the stuff
	fmt.Println("----")
	fmt.Println(h.Process(txn))

	log := txn.Receipt().Logs[0]
	vals, err := evntPairCreated.ParseLog(log)
	if err != nil {
		h.t.Fatal(err)
	}

	pairAddr := vals["pair"].(web3.Address)
	h.pairAddr = pairAddr

	pair := uniswap.NewUniswapV2Pair(pairAddr, srv.Provider())
	pair.Contract().SetFrom(srv.Owner())

	code, _ := srv.Provider().Eth().GetCode(pairAddr)
	if len(code) <= 5 {
		h.t.Fatal("bad")
	}

	return pair
}

func (h *Harness) addliquidity(token0Amount, token1Amount *big.Int) *contract.Txn {
	h.transferToken(h.token0, h.pairAddr, token0Amount)
	h.transferToken(h.token1, h.pairAddr, token1Amount)

	return h.SendTxn(h.pair.Mint(h.srv.Owner()))
}
