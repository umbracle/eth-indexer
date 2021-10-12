package main

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/umbracle/eth-indexer/schema"
	"github.com/umbracle/eth-indexer/sdk"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func encode(x starlark.Value) (interface{}, error) {
	switch obj := x.(type) {
	case starlark.String:
		return obj.String(), nil

	case starlark.Int:
		val, _ := obj.Int64()
		return val, nil

	default:
		return nil, fmt.Errorf("error")
	}
}

type objWrapper struct {
	obj *sdk.Obj
}

func (o *objWrapper) Module() *starlarkstruct.Module {
	return &starlarkstruct.Module{
		Name: o.obj.Table(),
		Members: starlark.StringDict{
			"set": starlark.NewBuiltin("Set", o.objSet),
		},
	}
}

func (o *objWrapper) objSet(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key string
	var valRaw starlark.Value
	if err := starlark.UnpackPositionalArgs("abs", args, kwargs, 1, &key, &valRaw); err != nil {
		return nil, err
	}

	fmt.Println("-- set --")
	fmt.Println(key, valRaw)

	val, err := encode(valRaw)
	if err != nil {
		return nil, err
	}
	o.obj.Set(key, val)
	return starlark.False, nil
}

type Indexer struct {
	snap *sdk.Snapshot
}

func (c *Indexer) get(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var table string
	var idRaw starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &table, &idRaw); err != nil {
		return nil, err
	}

	id, err := encode(idRaw)
	if err != nil {
		return nil, err
	}
	obj, err := c.snap.Get(table, id)
	if err != nil {
		return nil, err
	}
	wrapper := &objWrapper{
		obj: obj,
	}
	return wrapper.Module(), nil
}

func (c *Indexer) index(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	fmt.Println("xxxx")

	var table starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &table); err != nil {
		return nil, err
	}

	fmt.Println("-- table --")
	fmt.Println(reflect.TypeOf(table))

	return starlark.False, nil
}

func (c *Indexer) snapshot(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.False, nil
}

func (c *Indexer) Module() *starlarkstruct.Module {
	var Module = &starlarkstruct.Module{
		Name: "indexer",
		Members: starlark.StringDict{
			"get":      starlark.NewBuiltin("indexer.get", c.get),
			"index":    starlark.NewBuiltin("indexer.index", c.index),
			"snapshot": starlark.NewBuiltin("indexer.snapshot", c.snapshot),
		},
	}
	return Module
}

type Schema struct {
	tables []*schema.Table
}

func (s *Schema) add(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name string
	var fieldsRaw starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &name, &fieldsRaw); err != nil {
		return nil, err
	}

	fields, ok := fieldsRaw.(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("not a list")
	}

	iter := fields.Iterate()
	defer iter.Done()

	table := &schema.Table{
		Name:   name,
		Fields: []*schema.Field{},
	}
	var elem starlark.Value
	for i := 0; iter.Next(&elem); i++ {
		dict, ok := elem.(*starlark.Dict)
		if !ok {
			return nil, fmt.Errorf("not found")
		}

		decodeItem := func(name string) (interface{}, error) {
			val, ok, err := dict.Get(starlark.String(name))
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, fmt.Errorf("not found")
			}
			switch obj := val.(type) {
			case starlark.String:
				return obj.GoString(), nil
			default:
				return nil, fmt.Errorf("type not found %s", reflect.TypeOf(val).String())
			}
		}

		val, err := decodeItem("name")
		if err != nil {
			return nil, err
		}
		table.Fields = append(table.Fields, &schema.Field{
			Name: val.(string),
			Type: schema.TypeString,
		})
	}
	s.tables = append(s.tables, table)

	return starlark.False, nil
}

func (s *Schema) Module() *starlarkstruct.Module {
	var Module = &starlarkstruct.Module{
		Name: "schema",
		Members: starlark.StringDict{
			"add":   starlark.NewBuiltin("schema.add", s.add),
			"int64": starlark.String("a"),
		},
	}
	return Module
}

func main() {

	snap := sdk.NewSnapshot(sdk.NewMockStateResolver())
	snap.AddSchema("token", &schema.Table{
		Name: "token",
		Fields: []*schema.Field{
			{
				Name: "id",
				ID:   true,
				Type: schema.TypeAddress,
			},
			{
				Name: "val",
				Type: schema.TypeString,
			},
		},
	})

	ctxt := &Indexer{
		snap: snap,
	}

	sch := &Schema{}

	// Execute Starlark program in a file.
	thread := &starlark.Thread{
		Name: "my thread",
		Load: func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
			if module == "indexer.star" {
				return starlark.StringDict{"indexer": ctxt.Module()}, nil
			}
			if module == "schema.star" {
				return starlark.StringDict{"schema": sch.Module()}, nil
			}
			filename := filepath.Join(filepath.Dir(thread.CallFrame(0).Pos.Filename()), module)
			return starlark.ExecFile(thread, filename, nil, nil)
		},
	}
	globals, err := starlark.ExecFile(thread, "example.star", nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("-- globals --")
	fmt.Println(sch.tables)
	fmt.Println(globals)

	/*
		// Retrieve a module global.
		fibonacci := globals["fibonacci"]

		// Call Starlark function from Go.
		v, err := starlark.Call(thread, fibonacci, starlark.Tuple{starlark.MakeInt(10)}, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("fibonacci(10) = %v\n", v)
	*/
	// os.Exit(Run(os.Args[1:]))
}

/*
// Run starts the cli
func Run(args []string) int {
	commands := command.Commands()

	cli := &cli.CLI{
		Name:     "eth-indexer",
		Args:     args,
		Commands: commands,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
*/
