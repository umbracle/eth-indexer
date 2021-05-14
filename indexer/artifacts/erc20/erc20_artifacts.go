package erc20

import (
	"encoding/hex"
	"fmt"

	"github.com/umbracle/go-web3/abi"
)

var abiERC20 *abi.ABI

// ERC20Abi returns the abi of the ERC20 contract
func ERC20Abi() *abi.ABI {
	return abiERC20
}

var binERC20 []byte

// ERC20Bin returns the bin of the ERC20 contract
func ERC20Bin() []byte {
	return binERC20
}

func init() {
	var err error
	abiERC20, err = abi.NewABI(abiERC20Str)
	if err != nil {
		panic(fmt.Errorf("cannot parse ERC20 abi: %v", err))
	}
	if len(binERC20Str) != 0 {
		binERC20, err = hex.DecodeString(binERC20Str[2:])
		if err != nil {
			panic(fmt.Errorf("cannot parse ERC20 bin: %v", err))
		}
	}
}

var binERC20Str = "0x60806040523480156200001157600080fd5b50604051620022213803806200222183398181016040528101906200003791906200024f565b818160006200004b6200012560201b60201c565b9050806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3508160049080519060200190620001019291906200012d565b5080600590805190602001906200011a9291906200012d565b5050505050620003f3565b600033905090565b8280546200013b906200035f565b90600052602060002090601f0160209004810192826200015f5760008555620001ab565b82601f106200017a57805160ff1916838001178555620001ab565b82800160010185558215620001ab579182015b82811115620001aa5782518255916020019190600101906200018d565b5b509050620001ba9190620001be565b5090565b5b80821115620001d9576000816000905550600101620001bf565b5090565b6000620001f4620001ee84620002f6565b620002c2565b9050828152602081018484840111156200020d57600080fd5b6200021a84828562000329565b509392505050565b600082601f8301126200023457600080fd5b815162000246848260208601620001dd565b91505092915050565b600080604083850312156200026357600080fd5b600083015167ffffffffffffffff8111156200027e57600080fd5b6200028c8582860162000222565b925050602083015167ffffffffffffffff811115620002aa57600080fd5b620002b88582860162000222565b9150509250929050565b6000604051905081810181811067ffffffffffffffff82111715620002ec57620002eb620003c4565b5b8060405250919050565b600067ffffffffffffffff821115620003145762000313620003c4565b5b601f19601f8301169050602081019050919050565b60005b83811015620003495780820151818401526020810190506200032c565b8381111562000359576000848401525b50505050565b600060028204905060018216806200037857607f821691505b602082108114156200038f576200038e62000395565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611e1e80620004036000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063715018a611610097578063a457c2d711610066578063a457c2d71461029d578063a9059cbb146102cd578063dd62ed3e146102fd578063f2fde38b1461032d57610100565b8063715018a61461023b5780638da5cb5b1461024557806395d89b41146102635780639dc29fac1461028157610100565b8063313ce567116100d3578063313ce567146101a157806339509351146101bf57806340c10f19146101ef57806370a082311461020b57610100565b806306fdde0314610105578063095ea7b31461012357806318160ddd1461015357806323b872dd14610171575b600080fd5b61010d610349565b60405161011a9190611a13565b60405180910390f35b61013d600480360381019061013891906114b0565b6103db565b60405161014a91906119f8565b60405180910390f35b61015b6103f9565b6040516101689190611bb5565b60405180910390f35b61018b60048036038101906101869190611461565b610403565b60405161019891906119f8565b60405180910390f35b6101a9610504565b6040516101b69190611bd0565b60405180910390f35b6101d960048036038101906101d491906114b0565b61050d565b6040516101e691906119f8565b60405180910390f35b610209600480360381019061020491906114b0565b6105b9565b005b610225600480360381019061022091906113fc565b610643565b6040516102329190611bb5565b60405180910390f35b61024361068c565b005b61024d6107c6565b60405161025a91906119dd565b60405180910390f35b61026b6107ef565b6040516102789190611a13565b60405180910390f35b61029b600480360381019061029691906114b0565b610881565b005b6102b760048036038101906102b291906114b0565b61090b565b6040516102c491906119f8565b60405180910390f35b6102e760048036038101906102e291906114b0565b6109ff565b6040516102f491906119f8565b60405180910390f35b61031760048036038101906103129190611425565b610a1d565b6040516103249190611bb5565b60405180910390f35b610347600480360381019061034291906113fc565b610aa4565b005b60606004805461035890611d19565b80601f016020809104026020016040519081016040528092919081815260200182805461038490611d19565b80156103d15780601f106103a6576101008083540402835291602001916103d1565b820191906000526020600020905b8154815290600101906020018083116103b457829003601f168201915b5050505050905090565b60006103ef6103e8610c4d565b8484610c55565b6001905092915050565b6000600354905090565b6000610410848484610e20565b6000600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600061045b610c4d565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050828110156104db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d290611ad5565b60405180910390fd5b6104f8856104e7610c4d565b85846104f39190611c5d565b610c55565b60019150509392505050565b60006012905090565b60006105af61051a610c4d565b848460026000610528610c4d565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546105aa9190611c07565b610c55565b6001905092915050565b6105c1610c4d565b73ffffffffffffffffffffffffffffffffffffffff166105df6107c6565b73ffffffffffffffffffffffffffffffffffffffff1614610635576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161062c90611af5565b60405180910390fd5b61063f82826110a2565b5050565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b610694610c4d565b73ffffffffffffffffffffffffffffffffffffffff166106b26107c6565b73ffffffffffffffffffffffffffffffffffffffff1614610708576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106ff90611af5565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff1660008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a360008060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6060600580546107fe90611d19565b80601f016020809104026020016040519081016040528092919081815260200182805461082a90611d19565b80156108775780601f1061084c57610100808354040283529160200191610877565b820191906000526020600020905b81548152906001019060200180831161085a57829003601f168201915b5050505050905090565b610889610c4d565b73ffffffffffffffffffffffffffffffffffffffff166108a76107c6565b73ffffffffffffffffffffffffffffffffffffffff16146108fd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108f490611af5565b60405180910390fd5b61090782826111f7565b5050565b6000806002600061091a610c4d565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050828110156109d7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109ce90611b75565b60405180910390fd5b6109f46109e2610c4d565b8585846109ef9190611c5d565b610c55565b600191505092915050565b6000610a13610a0c610c4d565b8484610e20565b6001905092915050565b6000600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b610aac610c4d565b73ffffffffffffffffffffffffffffffffffffffff16610aca6107c6565b73ffffffffffffffffffffffffffffffffffffffff1614610b20576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b1790611af5565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610b90576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b8790611a75565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff1660008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610cc5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cbc90611b55565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610d35576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d2c90611a95565b60405180910390fd5b80600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610e139190611bb5565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610e90576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e8790611b35565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610f00576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ef790611a35565b60405180910390fd5b610f0b8383836113cd565b6000600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610f92576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f8990611ab5565b60405180910390fd5b8181610f9e9190611c5d565b600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546110309190611c07565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516110949190611bb5565b60405180910390a350505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611112576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161110990611b95565b60405180910390fd5b61111e600083836113cd565b80600360008282546111309190611c07565b9250508190555080600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546111869190611c07565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516111eb9190611bb5565b60405180910390a35050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611267576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161125e90611b15565b60405180910390fd5b611273826000836113cd565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050818110156112fa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112f190611a55565b60405180910390fd5b81816113069190611c5d565b600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816003600082825461135b9190611c5d565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516113c09190611bb5565b60405180910390a3505050565b505050565b6000813590506113e181611dba565b92915050565b6000813590506113f681611dd1565b92915050565b60006020828403121561140e57600080fd5b600061141c848285016113d2565b91505092915050565b6000806040838503121561143857600080fd5b6000611446858286016113d2565b9250506020611457858286016113d2565b9150509250929050565b60008060006060848603121561147657600080fd5b6000611484868287016113d2565b9350506020611495868287016113d2565b92505060406114a6868287016113e7565b9150509250925092565b600080604083850312156114c357600080fd5b60006114d1858286016113d2565b92505060206114e2858286016113e7565b9150509250929050565b6114f581611c91565b82525050565b61150481611ca3565b82525050565b600061151582611beb565b61151f8185611bf6565b935061152f818560208601611ce6565b61153881611da9565b840191505092915050565b6000611550602383611bf6565b91507f45524332303a207472616e7366657220746f20746865207a65726f206164647260008301527f65737300000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b60006115b6602283611bf6565b91507f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60008301527f63650000000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b600061161c602683611bf6565b91507f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008301527f64647265737300000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611682602283611bf6565b91507f45524332303a20617070726f766520746f20746865207a65726f20616464726560008301527f73730000000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b60006116e8602683611bf6565b91507f45524332303a207472616e7366657220616d6f756e742065786365656473206260008301527f616c616e636500000000000000000000000000000000000000000000000000006020830152604082019050919050565b600061174e602883611bf6565b91507f45524332303a207472616e7366657220616d6f756e742065786365656473206160008301527f6c6c6f77616e63650000000000000000000000000000000000000000000000006020830152604082019050919050565b60006117b4602083611bf6565b91507f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726000830152602082019050919050565b60006117f4602183611bf6565b91507f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360008301527f73000000000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b600061185a602583611bf6565b91507f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008301527f64726573730000000000000000000000000000000000000000000000000000006020830152604082019050919050565b60006118c0602483611bf6565b91507f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008301527f72657373000000000000000000000000000000000000000000000000000000006020830152604082019050919050565b6000611926602583611bf6565b91507f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008301527f207a65726f0000000000000000000000000000000000000000000000000000006020830152604082019050919050565b600061198c601f83611bf6565b91507f45524332303a206d696e7420746f20746865207a65726f2061646472657373006000830152602082019050919050565b6119c881611ccf565b82525050565b6119d781611cd9565b82525050565b60006020820190506119f260008301846114ec565b92915050565b6000602082019050611a0d60008301846114fb565b92915050565b60006020820190508181036000830152611a2d818461150a565b905092915050565b60006020820190508181036000830152611a4e81611543565b9050919050565b60006020820190508181036000830152611a6e816115a9565b9050919050565b60006020820190508181036000830152611a8e8161160f565b9050919050565b60006020820190508181036000830152611aae81611675565b9050919050565b60006020820190508181036000830152611ace816116db565b9050919050565b60006020820190508181036000830152611aee81611741565b9050919050565b60006020820190508181036000830152611b0e816117a7565b9050919050565b60006020820190508181036000830152611b2e816117e7565b9050919050565b60006020820190508181036000830152611b4e8161184d565b9050919050565b60006020820190508181036000830152611b6e816118b3565b9050919050565b60006020820190508181036000830152611b8e81611919565b9050919050565b60006020820190508181036000830152611bae8161197f565b9050919050565b6000602082019050611bca60008301846119bf565b92915050565b6000602082019050611be560008301846119ce565b92915050565b600081519050919050565b600082825260208201905092915050565b6000611c1282611ccf565b9150611c1d83611ccf565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611c5257611c51611d4b565b5b828201905092915050565b6000611c6882611ccf565b9150611c7383611ccf565b925082821015611c8657611c85611d4b565b5b828203905092915050565b6000611c9c82611caf565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b60005b83811015611d04578082015181840152602081019050611ce9565b83811115611d13576000848401525b50505050565b60006002820490506001821680611d3157607f821691505b60208210811415611d4557611d44611d7a565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000601f19601f8301169050919050565b611dc381611c91565b8114611dce57600080fd5b50565b611dda81611ccf565b8114611de557600080fd5b5056fea26469706673582212206ec672399ba153cd248cd866fdd37bd6590ab7e21f4c14918bade2abefc29f8b64736f6c63430008000033"

var abiERC20Str = `[
        {
            "inputs": [
                {
                    "internalType": "string",
                    "name": "name",
                    "type": "string"
                },
                {
                    "internalType": "string",
                    "name": "symbol",
                    "type": "string"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "constructor"
        },
        {
            "anonymous": false,
            "inputs": [
                {
                    "indexed": true,
                    "internalType": "address",
                    "name": "owner",
                    "type": "address"
                },
                {
                    "indexed": true,
                    "internalType": "address",
                    "name": "spender",
                    "type": "address"
                },
                {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "value",
                    "type": "uint256"
                }
            ],
            "name": "Approval",
            "type": "event"
        },
        {
            "anonymous": false,
            "inputs": [
                {
                    "indexed": true,
                    "internalType": "address",
                    "name": "previousOwner",
                    "type": "address"
                },
                {
                    "indexed": true,
                    "internalType": "address",
                    "name": "newOwner",
                    "type": "address"
                }
            ],
            "name": "OwnershipTransferred",
            "type": "event"
        },
        {
            "anonymous": false,
            "inputs": [
                {
                    "indexed": true,
                    "internalType": "address",
                    "name": "from",
                    "type": "address"
                },
                {
                    "indexed": true,
                    "internalType": "address",
                    "name": "to",
                    "type": "address"
                },
                {
                    "indexed": false,
                    "internalType": "uint256",
                    "name": "value",
                    "type": "uint256"
                }
            ],
            "name": "Transfer",
            "type": "event"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "owner",
                    "type": "address"
                },
                {
                    "internalType": "address",
                    "name": "spender",
                    "type": "address"
                }
            ],
            "name": "allowance",
            "outputs": [
                {
                    "internalType": "uint256",
                    "name": "",
                    "type": "uint256"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "spender",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "amount",
                    "type": "uint256"
                }
            ],
            "name": "approve",
            "outputs": [
                {
                    "internalType": "bool",
                    "name": "",
                    "type": "bool"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "account",
                    "type": "address"
                }
            ],
            "name": "balanceOf",
            "outputs": [
                {
                    "internalType": "uint256",
                    "name": "",
                    "type": "uint256"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "account",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "amount",
                    "type": "uint256"
                }
            ],
            "name": "burn",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [],
            "name": "decimals",
            "outputs": [
                {
                    "internalType": "uint8",
                    "name": "",
                    "type": "uint8"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "spender",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "subtractedValue",
                    "type": "uint256"
                }
            ],
            "name": "decreaseAllowance",
            "outputs": [
                {
                    "internalType": "bool",
                    "name": "",
                    "type": "bool"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "spender",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "addedValue",
                    "type": "uint256"
                }
            ],
            "name": "increaseAllowance",
            "outputs": [
                {
                    "internalType": "bool",
                    "name": "",
                    "type": "bool"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "account",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "amount",
                    "type": "uint256"
                }
            ],
            "name": "mint",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [],
            "name": "name",
            "outputs": [
                {
                    "internalType": "string",
                    "name": "",
                    "type": "string"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [],
            "name": "owner",
            "outputs": [
                {
                    "internalType": "address",
                    "name": "",
                    "type": "address"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [],
            "name": "renounceOwnership",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [],
            "name": "symbol",
            "outputs": [
                {
                    "internalType": "string",
                    "name": "",
                    "type": "string"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [],
            "name": "totalSupply",
            "outputs": [
                {
                    "internalType": "uint256",
                    "name": "",
                    "type": "uint256"
                }
            ],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "recipient",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "amount",
                    "type": "uint256"
                }
            ],
            "name": "transfer",
            "outputs": [
                {
                    "internalType": "bool",
                    "name": "",
                    "type": "bool"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "sender",
                    "type": "address"
                },
                {
                    "internalType": "address",
                    "name": "recipient",
                    "type": "address"
                },
                {
                    "internalType": "uint256",
                    "name": "amount",
                    "type": "uint256"
                }
            ],
            "name": "transferFrom",
            "outputs": [
                {
                    "internalType": "bool",
                    "name": "",
                    "type": "bool"
                }
            ],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        {
            "inputs": [
                {
                    "internalType": "address",
                    "name": "newOwner",
                    "type": "address"
                }
            ],
            "name": "transferOwnership",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        }
    ]`
