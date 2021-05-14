package sdk

import (
	"fmt"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"github.com/umbracle/go-web3/contract"
	"github.com/umbracle/go-web3/jsonrpc"
)

type ContractCaller struct {
	abis map[string]*callerGrp
}

type callerGrp struct {
	alias string
	items []*callerItem
}

type callerItem struct {
	c   *Caller
	abi *abi.ABI
}

type Caller struct {
	Signature  string
	DecodeHook func(v interface{}) interface{}
}

func (c *ContractCaller) AddCaller(name string, callers ...*Caller) error {
	if len(c.abis) == 0 {
		c.abis = map[string]*callerGrp{}
	}

	grp, ok := c.abis[name]
	if !ok {
		grp = &callerGrp{
			alias: name,
			items: []*callerItem{},
		}
		c.abis[name] = grp
	}
	for _, callr := range callers {
		abi, err := abi.NewABIFromList([]string{
			callr.Signature,
		})
		if err != nil {
			return err
		}
		grp.items = append(grp.items, &callerItem{
			c:   callr,
			abi: abi,
		})
	}
	return nil
}

func (c *ContractCaller) Call(alias string, addr web3.Address, provider *jsonrpc.Client) (interface{}, error) {
	grp, ok := c.abis[alias]
	if !ok {
		panic("Not found")
	}

	for _, item := range grp.items {
		// get the function
		var methodName string
		for k := range item.abi.Methods {
			methodName = k
		}

		fmt.Println(methodName)

		c1 := contract.NewContract(addr, item.abi, provider)
		vals, err := c1.Call(methodName, web3.Latest)
		if err != nil {
			fmt.Println("- not found -")
			continue
		}

		fmt.Println("-- vals --")
		fmt.Println(vals)

		ret := vals["0"]
		if item.c.DecodeHook != nil {
			ret = item.c.DecodeHook(ret)
		}
		return ret, nil
	}

	fmt.Println("- end empty handed -")
	return nil, fmt.Errorf("not found")
}
