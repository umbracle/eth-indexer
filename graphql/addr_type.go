package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/umbracle/go-web3"
)

func newWeb3Addr(addr string) *web3.Address {
	var add web3.Address
	if err := add.UnmarshalText([]byte(addr)); err != nil {
		return nil
	}
	return &add
}

var CustomAddressType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "CustomAddressType",
	Description: "The `CustomAddressType` scalar type represents a web3 Address Object.",

	// Serialize serializes `CustomID` to string.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case web3.Address:
			return value.String()
		case *web3.Address:
			v := *value
			return v.String()
		default:
			return nil
		}
	},

	// ParseValue parses GraphQL variables from `string` to `CustomID`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			return newWeb3Addr(value)
		case *string:
			return newWeb3Addr(*value)
		default:
			return nil
		}
	},

	// ParseLiteral parses GraphQL AST value to `CustomID`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return newWeb3Addr(valueAST.Value)
		default:
			return nil
		}
	},
})
