package pancake

import (
	"math/big"
	"strconv"

	"github.com/umbracle/eth-indexer/sdk"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
)

var tokenCaller *sdk.ContractCaller

func init() {
	// token contract caller
	tokenCaller := &sdk.ContractCaller{}

	// name
	tokenCaller.AddCaller("name", &sdk.Caller{
		Signature: "function name() returns (string)",
	}, &sdk.Caller{
		Signature: "function name() returns (bytes32)",
		DecodeHook: func(v interface{}) interface{} {
			raw := v.([32]byte)
			return trimRightZeros(raw[:])
		},
	}, &sdk.Caller{
		Signature: "function getName() returns (string)",
	})

	// symbol
	tokenCaller.AddCaller("symbol", &sdk.Caller{
		Signature: "function symbol() returns (string)",
	}, &sdk.Caller{
		Signature: "function symbol() returns (bytes32)",
		DecodeHook: func(v interface{}) interface{} {
			raw := v.([32]byte)
			return trimRightZeros(raw[:])
		},
	}, &sdk.Caller{
		Signature: "function getSymbol() returns (string)",
	})

	// decimals
	tokenCaller.AddCaller("decimals", &sdk.Caller{
		Signature: "function decimals() returns (uint8)",
		DecodeHook: func(v interface{}) interface{} {
			return strconv.Itoa(int(v.(uint8)))
		},
	}, &sdk.Caller{
		Signature: "function decimals() returns (uint256)",
		DecodeHook: func(v interface{}) interface{} {
			return v.(*big.Int).String()
		},
	})
}

var schema = map[string]*sdk.Resource{
	"ecosystem": {
		Schema: &sdk.Table{
			Fields: []*sdk.Field{
				{
					Name: "id",
					Type: sdk.TypeAddress,
					ID:   true,
				},
				{
					Name:    "numPairs",
					Type:    sdk.TypeUint,
					Default: uint64(0),
				},
				{
					Name:    "numTokens",
					Type:    sdk.TypeUint,
					Default: uint64(0),
				},
				{
					Name: "totalLiquidity",
					Type: sdk.TypeUint,
				},
				{
					Name: "totalVolume",
					Type: sdk.TypeUint,
				},
			},
		},
	},
	"pair": {
		Schema: &sdk.Table{
			Fields: []*sdk.Field{
				{
					Name: "address",
					Type: sdk.TypeAddress,
					ID:   true,
				},
				{
					Name:   "token0",
					Type:   sdk.TypeAddress,
					Static: true,
				},
				{
					Name:   "token1",
					Type:   sdk.TypeAddress,
					Static: true,
				},
				{
					Name:    "totalSupply",
					Type:    sdk.TypeDecimal,
					Default: sdk.Float0,
				},
				{
					Name:    "reserve0",
					Type:    sdk.TypeDecimal,
					Default: sdk.Float0,
				},
				{
					Name:    "reserve1",
					Type:    sdk.TypeDecimal,
					Default: sdk.Float0,
				},
				{
					Name:    "token0Price",
					Type:    sdk.TypeDecimal,
					Default: sdk.Float0,
				},
				{
					Name:    "token1Price",
					Type:    sdk.TypeDecimal,
					Default: sdk.Float0,
				},
				// list of events
				{
					Name:    "numSwapEvents",
					Type:    sdk.TypeUint,
					Default: uint64(0),
				},
				{
					Name:    "numMintEvents",
					Type:    sdk.TypeUint,
					Default: uint64(0),
				},
				{
					Name:    "numBurnEvents",
					Type:    sdk.TypeUint,
					Default: uint64(0),
				},
			},
		},
	},
	"token": {
		Schema: &sdk.Table{
			Fields: []*sdk.Field{
				{
					Name: "address",
					Type: sdk.TypeAddress,
					ID:   true,
				},
				{
					Name:   "name",
					Type:   sdk.TypeAddress,
					Static: true,
				},
				{
					Name:   "symbol",
					Type:   sdk.TypeAddress,
					Static: true,
				},
				{
					Name:   "decimals",
					Type:   sdk.TypeUint,
					Static: true,
				},
				{
					Name:    "numPairs",
					Type:    sdk.TypeUint,
					Default: uint64(0),
				},
			},
		},
		Init: func(addr web3.Address, provider *jsonrpc.Client, obj *sdk.Obj2) error {
			name, err := tokenCaller.Call("name", addr, provider)
			if err != nil {
				name = "empty"
			}

			decimals, err := tokenCaller.Call("decimals", addr, provider)
			if err != nil {
				// by default, TODO: Filter this only if not found.
				// contract caller should handle all the issues with the jsonrpc endpoint
				decimals = "18"
			}

			symbol, err := tokenCaller.Call("symbol", addr, provider)
			if err != nil {
				symbol = "empty"
			}

			// insert as int
			num, err := strconv.Atoi(decimals.(string))
			if err != nil {
				panic(err)
			}

			obj.Set("name", cleanStr(name))
			obj.Set("decimals", uint64(num))
			obj.Set("symbol", cleanStr(symbol))

			return nil
		},
	},
	"liquidity_event": {
		Schema: &sdk.Table{
			Fields: []*sdk.Field{
				{
					Name: "id",
					Type: sdk.TypeAddress,
					ID:   true,
				},
				{
					Name: "pair",
					Type: sdk.TypeAddress,
				},
				{
					Name: "eventType",
					Type: sdk.TypeAddress,
				},
				{
					Name: "amount0",
					Type: sdk.TypeDecimal,
				},
				{
					Name: "amount1",
					Type: sdk.TypeDecimal,
				},
			},
		},
	},
	"swap_event": {
		Schema: &sdk.Table{
			Fields: []*sdk.Field{
				{
					Name: "id",
					Type: sdk.TypeAddress,
					ID:   true,
				},
				{
					// references to pair
					Name: "pair",
					Type: sdk.TypeAddress,
				},
				{
					Name: "senderaddr",
					Type: sdk.TypeAddress,
				},
				{
					Name: "toaddr",
					Type: sdk.TypeAddress,
				},
				{
					Name: "amount0in",
					Type: sdk.TypeDecimal,
				},
				{
					Name: "amount1In",
					Type: sdk.TypeDecimal,
				},
				{
					Name: "amount0Out",
					Type: sdk.TypeDecimal,
				},
				{
					Name: "amount1Out",
					Type: sdk.TypeDecimal,
				},
			},
		},
	},
}
