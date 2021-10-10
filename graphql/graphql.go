package graphql

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/umbracle/eth-indexer/schema"
	"github.com/umbracle/eth-indexer/state"
	"github.com/umbracle/go-web3"
)

type Server struct {
	resolver state.State
}

type tuple struct {
	obj   *graphql.Object
	table *schema.Table
}

func (s *Server) Register(sch *schema.Schema) {
	var objs []*tuple
	for _, table := range sch.Tables {
		graphqlFields := graphql.Fields{}
		for _, f := range table.Fields {
			graphqlFields[f.Name] = &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					obj := p.Source.(*state.Obj)
					return obj.Data[p.Info.FieldName], nil
				},
			}
		}
		obj := graphql.NewObject(graphql.ObjectConfig{
			Name:   strings.Title(table.Name),
			Fields: graphqlFields,
		})
		objs = append(objs, &tuple{
			obj:   obj,
			table: table,
		})
	}

	queryFields := graphql.Fields{}
	for _, obj := range objs {
		// simple resolve object
		simpleObj := &graphql.Field{
			Type: obj.obj,
			Args: graphql.FieldConfigArgument{
				"address": &graphql.ArgumentConfig{
					Description: "Address of the pair",
					Type:        CustomAddressType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				addr := p.Args["address"].(*web3.Address)
				tableName := strings.TrimRight(p.Info.FieldName, "s")

				obj, err := s.resolver.GetObj(tableName, map[string]string{"address": addr.String()})
				if err != nil {
					return nil, err
				}
				return obj, nil
			},
		}

		// build using the Table stuff
		listArgs := graphql.FieldConfigArgument{
			"first": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
			"skip": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
			"orderBy": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"orderDirection": &graphql.ArgumentConfig{
				DefaultValue: "asc",
				Type:         graphql.String,
			},
		}

		listObj := &graphql.Field{
			Type: graphql.NewList(obj.obj),
			Args: listArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tableName := strings.TrimRight(p.Info.FieldName, "s")

				query := &state.Query{
					Table:   tableName,
					First:   uint64(p.Args["first"].(int)),
					Skip:    uint64(p.Args["skip"].(int)),
					OrderBy: p.Args["orderBy"].(string),
					Order:   p.Args["orderDirection"].(string),
				}
				objs, err := s.resolver.GetObjs(query)
				if err != nil {
					return nil, err
				}
				return objs, nil
			},
		}

		queryFields[obj.table.Name] = simpleObj
		queryFields[obj.table.Name+"s"] = listObj
	}

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: queryFields,
	})

	schemaConfig := graphql.SchemaConfig{
		Query: queryType,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query

	query := `
	query {
		swap_events (first: 10, where: {pair: {eq: "foo"}}) {
			id
		}
	}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		fmt.Println(r.Errors[0])
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
}
