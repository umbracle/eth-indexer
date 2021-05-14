package sdk

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/umbracle/go-web3"
)

type FieldType int

const (
	TypeAddress FieldType = iota + 1
	TypeUint
	TypeDecimal
)

func (f *Field) Decode(val string) (interface{}, error) {
	switch f.Type {
	case TypeAddress:
		return val, nil

	case TypeUint:
		v, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("failed to decode uint: %v", val)
		}
		return uint64(v), nil

	case TypeDecimal:
		f := new(Float)
		if !f.SetString(val) {
			return nil, fmt.Errorf("failed to decode float: %v", val)
		}
		return f, nil

	default:
		panic(fmt.Sprintf("Decode type not found %v", f.Type))
	}
}

func (f *Field) Encode(raw interface{}) (string, error) {
	if raw == nil {
		return "", fmt.Errorf("is empty")
	}
	switch f.Type {
	case TypeAddress:
		var val string
		switch obj := raw.(type) {
		case string:
			val = obj
		case web3.Address:
			val = obj.String()
		default:
			return "", fmt.Errorf("bad1 %s", raw)
		}
		return val, nil

	case TypeUint:
		var val string
		switch obj := raw.(type) {
		case *big.Int:
			val = obj.String()
		case uint64:
			val = strconv.Itoa(int(obj))
		default:
			return "", fmt.Errorf("%s bad 2 %s", f.Name, reflect.TypeOf(raw))
		}
		return val, nil

	case TypeDecimal:
		var val string
		switch obj := raw.(type) {
		case *Float:
			val = obj.String()
		default:
			return "", fmt.Errorf("bad 3: %v", raw)
		}
		return val, nil

	default:
		panic(fmt.Sprintf("Decode type not found %v", f.Type))
	}
}

type Field struct {
	Name        string
	References  *Reference
	ID          bool
	Static      bool
	Default     interface{}
	Type        FieldType
	Description string
}

type Reference struct {
	Table string
	Field string
}

type Schema2 struct {
	Fields []*Field
}

type Table struct {
	Name   string
	Fields []*Field
}

func (t *Table) getField(id string) *Field {
	for _, f := range t.Fields {
		if f.Name == id {
			return f
		}
	}
	return nil
}

func (t *Table) getIDS() []*Field {
	res := []*Field{}
	for _, j := range t.Fields {
		if j.ID {
			res = append(res, j)
		}
	}
	return res
}

type Schema struct {
	Tables []*Table
}

/*
func buildSchema(evnt *abi.Event) (*proto.Schema, error) {
	sch := &proto.Schema{
		Fields: []*proto.Field{},
	}
	for _, f := range evnt.Inputs.TupleElems() {
		var typ proto.Field_Type
		if f.Elem.Kind() == abi.KindUInt {
			typ = proto.Field_Uint
		}
		field := &proto.Field{
			Name: f.Name,
			Type: typ,
		}
		sch.Fields = append(sch.Fields, field)
	}
	return sch, nil
}
*/
