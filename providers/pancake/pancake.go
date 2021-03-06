package pancake

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/umbracle/eth-indexer/indexer/proto"
	"github.com/umbracle/eth-indexer/sdk"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
)

func addrPtr(addr web3.Address) *web3.Address {
	return &addr
}

func trimRightZeros(b []byte) string {
	res := bytes.TrimRightFunc(b, func(r rune) bool {
		return r == 0
	})
	return string(res)
}

func cleanStr(s interface{}) string {
	i := s.(string)
	i = strings.Replace(i, "'", "", -1)
	return i
}

func Provider() *sdk.Provider {
	factoryAddr := web3.HexToAddress("0xBCfCcbde45cE874adCB698cC183deBcF17952812")
	routerAddr := web3.HexToAddress("0x5c69bee701ef814a2b6a3edd4b1652cb9cc5aa6f")

	return &sdk.Provider{
		Resources: schema,
		Filter: &sdk.FilterByAddr{
			// Catch all the events generated by the factory and router
			FromAddr: factoryAddr,
			ToAddr:   routerAddr,
		},
		Trackers: []*sdk.Tracker{
			{
				// pair-created
				Type: evntPairCreated,
				Handler: func(req *sdk.HandlerReq) {
					fmt.Println("_______ PAIR CREATED _________")
					fmt.Println(req.Evnt.Address, req.Evnt.TxHash)

					vals := req.Vals

					ecosystem := req.Get("ecosystem", "0")

					t0 := vals["token0"].(web3.Address)
					t1 := vals["token1"].(web3.Address)

					t0T := req.Get("token", t0)
					t1T := req.Get("token", t1)

					// add values from the ecosystem
					if t0T.IsNew() {
						ecosystem.Incr("numTokens")
					}
					if t1T.IsNew() {
						ecosystem.Incr("numTokens")
					}
					ecosystem.Incr("numPairs")

					// Add pair info
					pair := req.Get("pair", vals["pair"].(web3.Address).String())
					pair.Set("token0", t0.String())
					pair.Set("token1", t1.String())

					// Add token num Pairs
					t0T.Incr("numPairs")
					t1T.Incr("numPairs")
				},
			},
			{
				// mint event
				Type: evntMint,
				Handler: func(req *sdk.HandlerReq) {
					liquidityEvent(req, "mint")
				},
			},
			{
				// burn event
				Type: evntBurn,
				Handler: func(req *sdk.HandlerReq) {
					liquidityEvent(req, "burn")
				},
			},
			{
				// swap event
				Type: evntSwap,
				Handler: func(req *sdk.HandlerReq) {
					handleSwap(req)
				},
			},
		},
		Snapshots: map[string]*sdk.Snapshot2{
			"tokens_numPairs": {
				Table:     "token",
				Index:     []string{"numPairs"},
				SplitFunc: sdk.BlockSplitFunc(100),
			},
		},
	}
}

func handleSwap(req *sdk.HandlerReq) {
	swapEvent := &SwapEvent{}
	if err := sdk.DecodeEvent(&swapEvent, req.Vals); err != nil {
		panic(err)
	}

	ensemble := loadEnsemble(req, req.Evnt.Address)

	swap := req.Get("swap_event", sdk.UUID())
	swap.Set("pair", req.Evnt.Address)

	// Convert to decimals and add to the indexed event
	amount0In := ensemble.Token0.ToDecimals(swapEvent.Amount0In)
	swap.Set("amount0in", amount0In)

	amount1In := ensemble.Token1.ToDecimals(swapEvent.Amount1In)
	swap.Set("amount1In", amount1In)

	amount0Out := ensemble.Token0.ToDecimals(swapEvent.Amount0Out)
	swap.Set("amount0Out", amount0Out)

	amount1Out := ensemble.Token1.ToDecimals(swapEvent.Amount1Out)
	swap.Set("amount1Out", amount1Out)

	// Get total amounts of the operation
	amount0Total := amount0In.Add(amount0Out)
	amount1Total := amount1In.Add(amount1Out)

	// TODO: Store
	fmt.Println(amount0Total)
	fmt.Println(amount1Total)
}

func parseLogs(evnt *abi.Event, event proto.Event) (map[string]interface{}, error) {
	log, err := event.ToLog()
	if err != nil {
		panic(err)
	}
	vals, err := evnt.ParseLog(log)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

type Ensemble struct {
	Pair   *sdk.Obj2
	Token0 *token
	Token1 *token
}

type token struct {
	*sdk.Obj2
}

func (t *token) ToDecimals(i *big.Int) *sdk.Float {
	return new(sdk.Float).SetBigInt(i).DivUint(t.Get("decimals").(uint64))
}

func loadEnsemble(req *sdk.HandlerReq, addr string) *Ensemble {
	pair := req.Get("pair", addr)
	token0 := req.Get("token", pair.Get("token0"))
	token1 := req.Get("token", pair.Get("token1"))

	return &Ensemble{
		Pair:   pair,
		Token0: &token{token0},
		Token1: &token{token1},
	}
}

func handleSync(req *sdk.HandlerReq, ensemble *Ensemble, event proto.Event) {
	log, err := event.ToLog()
	if err != nil {
		panic(err)
	}
	vals, err := evntSync.ParseLog(log)
	if err != nil {
		panic(err)
	}

	// get reserve prices in decimal format
	reserve0 := vals["reserve0"].(*big.Int)
	reserve1 := vals["reserve1"].(*big.Int)

	reserve0Dec := new(sdk.Float).SetBigInt(reserve0).DivUint(ensemble.Token0.Get("decimals").(uint64))
	reserve1Dec := new(sdk.Float).SetBigInt(reserve1).DivUint(ensemble.Token1.Get("decimals").(uint64))

	ensemble.Pair.Set("token0Price", reserve0Dec.Div(reserve1Dec))
	ensemble.Pair.Set("token1Price", reserve1Dec.Div(reserve0Dec))
}

func liquidityEvent(req *sdk.HandlerReq, typ string) {
	obj := req.Get("liquidity_event", sdk.UUID())
	obj.Set("pair", req.Evnt.Address)
	obj.Set("eventType", typ)

	ensemble := loadEnsemble(req, req.Evnt.Address)

	amount0 := req.Vals["amount0"].(*big.Int)
	amount1 := req.Vals["amount1"].(*big.Int)

	amount0Dec := new(sdk.Float).SetBigInt(amount0).DivUint(ensemble.Token0.Get("decimals").(uint64))
	amount1Dec := new(sdk.Float).SetBigInt(amount1).DivUint(ensemble.Token1.Get("decimals").(uint64))

	obj.Set("amount0", amount0Dec)
	obj.Set("amount1", amount1Dec)

	// always do sync
	handleSync(req, ensemble, req.Action.Events[req.Indx-1])

	// transfer
	transferVals, err := parseLogs(evntTransfer, req.Action.Events[req.Indx-2])
	if err != nil {
		return
	}
	transferValue := transferVals["value"].(*big.Int)
	if typ == "mint" {
		// mint event
		ensemble.Pair.Add("totalSupply", new(sdk.Float).SetBigInt(transferValue))

	} else {
		// burn event
		ensemble.Pair.Sub("totalSupply", new(sdk.Float).SetBigInt(transferValue))
	}
}
