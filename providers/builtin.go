package providers

import (
	"github.com/umbracle/eth-indexer/providers/hashmask"
	"github.com/umbracle/eth-indexer/providers/pancake"
	"github.com/umbracle/eth-indexer/sdk"
)

type ProviderFactory func() *sdk.Provider

var BuiltinProviders = map[string]ProviderFactory{
	"pancake":  pancake.Provider,
	"hashmask": hashmask.Provider,
}
