package sdk

import (
	"fmt"
	"strconv"

	"github.com/umbracle/eth-indexer/indexer/proto"
	protosdk "github.com/umbracle/eth-indexer/sdk/proto"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"github.com/umbracle/go-web3/jsonrpc"
)

var _ Backend = &Provider{}

type HandlerReq struct {
	*Snapshot
	Evnt *proto.Event
	Vals map[string]interface{}

	Indx   int
	Action *Action
}

type Handler func(HandlerReq)

type Tracker struct {
	Type    *abi.Event
	Handler func(*HandlerReq)
}

type ResourceInit func(addr web3.Address, provider *jsonrpc.Client, obj *Obj2) error

type Resource struct {
	Name      string
	Singleton bool
	Schema    *Table
	Init      ResourceInit
}

type Provider struct {
	Resources map[string]*Resource
	Snapshots map[string]*Snapshot2
	Trackers  []*Tracker
	Filter    *FilterByAddr

	// state resolver
	resolver StateResolver

	client *jsonrpc.Client

	// the list of all the schemas of this provider
	schemas map[string]*Table

	// indexers process the block
	indexers []indexer

	snap *Snapshot
}

func (p *Provider) GetFilter() *FilterByAddr {
	return p.Filter
}

func (p *Provider) SetStateResolver(resolver StateResolver) {
	p.resolver = resolver
}

func (p *Provider) SetClient(client *jsonrpc.Client) {
	p.client = client
}

func (p *Provider) initSchema(name string, addr string, obj *Obj2) error {
	initFn := p.Resources[name]
	if initFn.Init == nil {
		return nil
	}
	return initFn.Init(web3.HexToAddress(addr), p.client, obj)
}

func (p *Provider) addSchema(name string, sch *Table) {
	sch.Name = name

	// add extra fields into the schema
	sch.Fields = append(sch.Fields, extraFields...)
	p.schemas[name] = sch
}

type indexer interface {
	Process(act *Action, i *Snapshot) error
}

func (p *Provider) Init() error {
	p.schemas = map[string]*Table{}
	p.indexers = []indexer{}

	// build the schemas for the resources
	for name, c := range p.Resources {
		p.addSchema(name, c.Schema)
	}

	// parse the event trackers after the snapshots since we want
	// all the snapshot indexers at the end of the execution
	for _, t := range p.Trackers {
		p.indexers = append(p.indexers, &trackerIndexer22{tracker: t})
	}

	// build schemas for snapshots
	for name, def := range p.Snapshots {
		if err := p.buildSnapshot(name, def); err != nil {
			return err
		}
	}

	p.snap = newSnapshot()
	p.snap.provider = p
	p.snap.schemas = p.schemas
	return nil
}

type Snapshot2 struct {
	Table     string
	Index     []string
	SplitFunc SplitFunc
}

func (p *Provider) buildSnapshot(snapshotName string, snapshot *Snapshot2) error {
	// find the target table
	table, ok := p.schemas[snapshot.Table]
	if !ok {
		return fmt.Errorf("table '%s' not found", snapshot.Table)
	}

	snapshostSchema := &Table{
		Name:   snapshotName,
		Fields: []*Field{},
	}

	fmt.Println("-- snap schema --")
	fmt.Println(snapshostSchema)

	// get the ids of the table since those work as anchors
	idFields := table.getIDS()
	snapshostSchema.Fields = append(snapshostSchema.Fields, idFields...)

	// get the index fields
	for _, fieldName := range snapshot.Index {
		if field := table.getField(fieldName); field != nil {
			snapshostSchema.Fields = append(snapshostSchema.Fields, field)
		} else {
			return fmt.Errorf("field '%s' does not exists", fieldName)
		}
	}

	// append the split separator, its an ID as well in order to make
	// the index entry unique
	blockField := &Field{
		Name: "block",
		ID:   true,
		Type: TypeAddress,
	}
	snapshostSchema.Fields = append(snapshostSchema.Fields, blockField)

	p.addSchema(snapshotName, snapshostSchema)

	process := &snapIndexer22{
		snapshot:     snapshot,
		snapshotName: snapshotName,
		schema:       snapshostSchema,
	}
	p.indexers = append(p.indexers, process)
	return nil
}

func (p *Provider) Process(act *Action) ([]*protosdk.Diff, *ErrorEvent) {
	// start the snapshot
	p.snap.block = act.BlockNum

	// loop the indexers
	closeCh := make(chan struct{})
	go func() {
		defer func() {
			close(closeCh)
		}()
		for _, ii := range p.indexers {
			if err := ii.Process(act, p.snap); err != nil {
				p.snap.finish(&ErrorEvent{
					Type: ErrorEventGeneric,
					Err:  err,
				})
			}
		}
	}()

	<-closeCh

	// save the snapshot to generate the diffs
	diffs := p.snap.save()
	p.snap.reset()

	return diffs, p.snap.err
}

func (p *Provider) GetSchemas() GetSchemasResponse {
	resp := GetSchemasResponse{}
	for _, sch := range p.schemas {
		resp.Schemas = append(resp.Schemas, sch)
	}
	return resp
}

type SplitFunc func(block uint64) string

func BlockSplitFunc(blockSize uint64) SplitFunc {
	return func(block uint64) string {
		return strconv.Itoa(int(block / blockSize))
	}
}

func NoSplitFunc(block uint64) SplitFunc {
	// generates one snapshot per each block
	return BlockSplitFunc(1)
}

type trackerIndexer22 struct {
	tracker *Tracker
}

func (s *trackerIndexer22) Process(ac *Action, i *Snapshot) error {
	for indx, evnt := range ac.Events {
		if evnt.TopicID == s.tracker.Type.ID().String() {
			log, err := evnt.ToLog()
			if err != nil {
				return err
			}
			vals, err := s.tracker.Type.ParseLog(log)
			if err != nil {
				continue
			}
			req := &HandlerReq{
				Snapshot: i,
				Evnt:     &evnt,
				Vals:     vals,
				Action:   ac,
				Indx:     indx,
			}
			s.tracker.Handler(req)
		}
	}
	return nil
}

type snapIndexer22 struct {
	snapshot     *Snapshot2
	snapshotName string
	schema       *Table
}

func (s *snapIndexer22) Process(ac *Action, i *Snapshot) error {
	block := i.block

	for _, obj := range i.trackedObjs {
		if obj.table != s.snapshot.Table {
			return nil
		}

		indexColName := s.snapshot.Index[0]
		val, ok := obj.hasChanged(indexColName)
		if !ok {
			return nil
		}

		var numKey string
		if s.snapshot.SplitFunc == nil {
			// store this entry for sure so use the blcoka s index
			numKey = strconv.Itoa(int(block))
		} else {
			numKey = s.snapshot.SplitFunc(block)
		}

		kk := []interface{}{}
		for name, v := range obj.key {
			// we need to encode it not as string but as his own type
			val, err := s.schema.getField(name).Decode(v)
			if err != nil {
				return err
			}
			kk = append(kk, val)
		}

		kk = append(kk, numKey)

		index := i.Get(s.snapshotName, kk...)
		index.Set(indexColName, val)
	}
	return nil
}

// default fields for all the items
var extraFields = []*Field{}
