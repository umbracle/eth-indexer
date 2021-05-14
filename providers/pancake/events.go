package pancake

import (
	"math/big"

	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
)

var (
	factoryContract = web3.HexToAddress("0xBCfCcbde45cE874adCB698cC183deBcF17952812")
	routerContract  = web3.HexToAddress("0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F")
)

type PairCreatedEvent struct {
	Token0 web3.Address
	Token1 web3.Address
	Pair   web3.Address
}

var evntPairCreated = abi.MustNewEvent(`PairCreated(
	address indexed token0, 
	address indexed token1,
	address pair, 
	uint256
)`)

type MintEvent struct {
	Sender  web3.Address
	Amount0 *big.Int
	Amount1 *big.Int
}

var evntMint = abi.MustNewEvent(`Mint(
	address indexed sender, 
	uint256 amount0, 
	uint256 amount1
)`)

type BurnEvent struct {
	Sender  web3.Address
	Amount0 *big.Int
	Amount1 *big.Int
	To      web3.Address
}

var evntBurn = abi.MustNewEvent(`Burn(
	address indexed sender,
	uint256 amount0,
	uint256 amount1,
	address indexed to
)`)

type SwapEvent struct {
	Sender     web3.Address `mapstructure:"sender"`
	Amount0In  *big.Int     `mapstructure:"amount0In"`
	Amount1In  *big.Int     `mapstructure:"amount1In"`
	Amount0Out *big.Int     `mapstructure:"amount0Out"`
	Amount1Out *big.Int     `mapstructure:"amount1Out"`
	To         web3.Address `mapstructure:"to"`
}

var evntSwap = abi.MustNewEvent(`Swap(
	address indexed sender,
	uint256 amount0In,
	uint256 amount1In,
	uint256 amount0Out,
	uint256 amount1Out,
	address indexed to
)`)

type SyncEvent struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}

var evntSync = abi.MustNewEvent(`Sync(
	uint112 reserve0, 
	uint112 reserve1
)`)

var evntApproval = abi.MustNewEvent(`Approval(
	address indexed owner, 
	address indexed spender, 
	uint256 value
)`)

var evntTransfer = abi.MustNewEvent(`Transfer(
	address indexed from,
	address indexed to,
	uint256 value
)`)

// testing

var eth2Deposit = abi.MustNewEvent(`DepositEvent(
	bytes pubkey,
	bytes whitdrawalcred,
	bytes amount,
	bytes signature,
	bytes index
)`)
