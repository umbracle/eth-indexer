
# Eth-Indexer

Eth-indexer is an indexer of Ethereum events. Write your own extensions to specify how to index, store and track contracts or events. It includes support to index out of the box things like snapshots or aggregates.

It uses Postgresql as datastore and has support to expose the data in Graphql or API rest automatically.

## How does it work?

You can find examples on how to build your own plugin on the **providers** folder. There, we have examples of plugins to index the popular Pancakeswap and Hashmask contracts.

First, you need to create a function that returns an sdk.Provider object. This is a framework that providers helper methods and ergonomic methods for the developer to index data.

```
func Provider() *sdk.Provider {
	return &sdk.Provider{
    }
}
```

Next, we will be filling this framework with our custom extension. 

### Schema

The first thing we need to do is to specify the schema of your data, this is, what type of data do you want to store in the database. For example, with PancakeSwap we might want to store all the pairs created, the tokens of each pair and entries for each event in the contract: swap, mint or burn. We define the data types and the schema like this:

```
Schema: map[string]*sdk.Resource{
    "pair": {
        Schema: &sdk.Table{
            Fields: []*sdk.Field{
                {
                    Name: "address",
                    Type: sdk.TypeAddress,
                    ID:   true,
                },
                ...
            },
        },
    },
    "token": {
        Schema: &sdk.Table{
            Fields: []*sdk.Field{
                {
                    Name: "address",
                    Type: sdk.TypeAddress,
                    ID: true,
                },
                {
                    Name: "numPairs",
                    Type: sdk.TypeUint,
                    Default: uint64(0),
                },
                ...
            },
        },
    },
    "ecosystem": {
        Schema: &sdk.Table{
            Fields: []*sdk.Field{
				{
					Name:    "numPairs",
					Type:    sdk.TypeUint,
					Default: uint64(0),
				},
                ...
            },
        },
    },
}
```

Note that we can also define things that represent the whole Dapp context like the 'ecosystem' schema.

Upon starting, eth-indexer takes all the schemas in the extension (including the self-generated, more on this later), creates the tables in the datastore (or updates them) and creates the Graphql types and endpoints. Fields with ID=true are the primary key value for that data type. 

We recommend using random IDs for types like Transfer or Approval events.

### Filter

Now, we select which contracts we are interested in filtering.

```
Filter: &sdk.FilterByAddr{
	FromAddr: factoryAddr,
},
```

### Track

Once we have our filter to track events (ethereum) and the schemas to store the information (datastore) we need to write our custom logic to fill the data from one side to the other.

To do so, we write a custom callback that gets triggered for a specific event.

```
Trackers: []*sdk.Tracker{
	{
		Type: evntPairCreated,
		Handler: func(req *sdk.HandlerReq) {
            // apply logic
        },
    },
}
```

The object 'sdk.HandlerReq' provides many functions to interact with the event that generated the callback and the current state. This is an example from the PancakeSwap extension:

```
vals := req.Vals // event values
ecosystem := req.Get("ecosystem", "0")

t0 := vals["token0"].(web3.Address)
t1 := vals["token1"].(web3.Address)

t0T := req.Get("token", t0)
t1T := req.Get("token", t1)

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
```

This is the description of some of the functions:
- req.Get(\<schema name>, \<id>...) \<Obj>: Return the object from the schema with the given id. If the object is not fund, create it. There can be more than one ids for the object.
- \<Obj>.IsNew(): Whether the object has been created right now.
- \<Obj>.Set(key: string, val: \<any>): Set the value for 'key' in that object.
- \<Obj>.Incr(key): Increase the count in that key, only if it is a numeric type (int or float).

You can check the PancakeSwap extension to learn more about other functions and helper primitives.

### Snapshots

Note that using the previous primitives it is possible to build complex things like snapshots or aggregates in time during a specific period. However, writting that repetitive logic by hand is tedious and error prone. Thus, eth-indexer provides native support for snapshots:

```
Snapshots: map[string]*sdk.Snapshot2{
	"tokens_numPairs": {
		Table:     "token",
		Index:     []string{"numPairs"},
		SplitFunc: sdk.BlockSplitFunc(100),
	},
},
```

In this example, we want to have a snapshot every 100 blocks that has the number of pairs in which the token is included at that point in time. If SplitFund is empty, a snapshot is taken everytime there is a change in the index field.

Automatically, eth-indexer will create another schema for this data type. In this example, the schema looks like this:

```
Schema: &sdk.Table{
    Fields: []*sdk.Field{
		{
			Name: "address",
			Type:  sdk.Address,
		},
        {
            Name: "numPairs",
            Type: sdk.TypeUint,
        },
        {
            Name: "block",
            Type: sdk.TypeUint,
        },
    },
},
```

Note that it also includes in the new schema the ID fields of the 'token' table plus the indexed field.

In the future we want to expand the number of snaphots to include things like aggregates and snapshots for a certain period (i.e. number of pairs created from blocks x to y).

## Performance

It takes less than 30 seconds to compute all the hashmask events and around 4 hours to index 6 million PancakeSwap events with less than 1Gb of memory.
