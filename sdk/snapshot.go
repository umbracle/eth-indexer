package sdk

import (
	"encoding/hex"
	"fmt"

	"github.com/umbracle/eth-indexer/schema"
	"github.com/umbracle/eth-indexer/state"
)

type objErr interface {
	finish(*ErrorEvent)
}

type Obj struct {
	schema *schema.Table
	// TODO: provider interface to call finish
	// objErr  objErr
	created bool
	id      []byte
	table   string
	key     map[string]string
	vals    map[string]string
	changes map[string]string
}

func (o *Obj) Table() string {
	return o.table
}

func (o *Obj) Keys() map[string]string {
	return o.key
}

func (o *Obj) IsNew() bool {
	return o.created
}

func (o *Obj) getField(name string) (*schema.Field, bool) {
	for _, f := range o.schema.Fields {
		if f.Name == name {
			return f, true
		}
	}
	return nil, false
}

func (o *Obj) isChanged() bool {
	return len(o.changes) != 0
}

func (o *Obj) hasChanged(key string) (string, bool) {
	val, ok := o.changes[key]
	return val, ok
}

func (o *Obj) Set(key string, val interface{}) error {
	// convert val to a string representation
	field, ok := o.getField(key)
	if !ok {
		return fmt.Errorf("field not found '%s'", key)
	}
	valStr, err := field.Encode(val)
	if err != nil {
		return fmt.Errorf("failed to encode schema %s %s with %s: %v", o.table, key, val, err)
	}

	if raw, ok := o.vals[key]; ok {
		if raw != val {
			o.changes[key] = valStr
		}
	} else {
		o.changes[key] = valStr
	}
	return nil
}

func (o *Obj) Get(key string) (interface{}, error) {
	// make sure first the value exists
	field, ok := o.getField(key)
	if !ok {
		return nil, fmt.Errorf("field '%s' not found", key)
	}

	// try in changes first
	valStr, ok := o.changes[key]
	if !ok {
		// use the state value
		valStr, ok = o.vals[key]
		if !ok {
			panic("bad")
		}
	}
	// convert val to his type
	val, err := field.Decode(valStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode")
		/*
			o.objErr.finish(&ErrorEvent{
				Type: ErrorEventGetDecode,
				Err:  fmt.Errorf("failed to decode %s with %s: %v", key, val, err),
			})
		*/
	}
	return val, nil
}

func (o *Obj) expect(key string, fieldType schema.FieldType) error {
	field, ok := o.getField(key)
	if !ok {
		return fmt.Errorf("field '%s' not found", key)
	}
	if field.Type != fieldType {
		return fmt.Errorf("field %s expected %d but found %d", key, fieldType, field.Type)
	}
	return nil
}

func (o *Obj) Incr(key string) error {
	return o.Add(key, uint64(1))
}

func (o *Obj) Sub(key string, v interface{}) error {
	var val interface{}

	var err error
	switch obj := v.(type) {
	case *schema.Float:
		if err := o.expect(key, schema.TypeDecimal); err != nil {
			return err
		}
		val, err = o.Get(key)
		if err != nil {
			return err
		}
		val.(*schema.Float).Sub(obj)

	case uint64:
		if err := o.expect(key, schema.TypeUint); err != nil {
			return err
		}
		val, err = o.Get(key)
		if err != nil {
			return err
		}
		val = val.(uint64) - obj
	}
	o.Set(key, val)
	return nil
}

func (o *Obj) Add(key string, v interface{}) error {
	var val interface{}

	var err error
	switch obj := v.(type) {
	case *schema.Float:
		if err := o.expect(key, schema.TypeDecimal); err != nil {
			return err
		}
		val, err = o.Get(key)
		if err != nil {
			return err
		}
		val.(*schema.Float).Add(obj)

	case uint64:
		if err := o.expect(key, schema.TypeUint); err != nil {
			return err
		}
		val, err = o.Get(key)
		if err != nil {
			return err
		}
		val = val.(uint64) + obj

	default:
		panic("Not expected")
	}
	o.Set(key, val)
	return nil
}

func (o *Obj) Copy() *Obj {
	oo := new(Obj)
	*oo = *o

	oo.key = map[string]string{}
	for k, v := range o.key {
		oo.key[k] = v
	}

	oo.vals = map[string]string{}
	for k, v := range o.vals {
		oo.vals[k] = v
	}
	return oo
}

func (s *Snapshot) save() []*schema.Diff {
	diffs := []*schema.Diff{}

	for _, obj := range s.trackedObjs {
		if obj.isChanged() {
			obj := obj.Copy()

			diff := &schema.Diff{
				Table:    obj.table,
				Keys:     obj.key,
				Creation: obj.created,
				Vals:     obj.changes,
			}
			diffs = append(diffs, diff)

			for k, v := range obj.changes {
				obj.vals[k] = v
			}

			// reset the object
			obj.changes = map[string]string{}
			obj.created = false

			s.inmemStore.add(string(obj.id), obj)
			obj.changes = map[string]string{} // reset changes in parent
		}
	}
	return diffs
}

func buildIndex(table string, keys []string) []byte {
	res := table
	for _, j := range keys {
		res += j
	}
	return []byte(res)
}

func (s *Snapshot) create(table string, schema *schema.Table, id string, keys map[string]string) *Obj {
	obj := &Obj{
		schema:  schema,
		created: true,
		id:      []byte(id),
		table:   table,
		key:     keys,
		vals:    map[string]string{},
		changes: map[string]string{},
	}
	s.trackedObjs[id] = obj
	return obj
}

func (s *Snapshot) getOk(table string, id string) (*Obj, bool) {
	// check tracked objects
	if obj, ok := s.trackedObjs[id]; ok {
		return obj, true
	}

	// check inmemory cache
	obj, ok := s.inmemStore.get(id)
	if ok {
		s.trackedObjs[id] = obj
		return obj, true
	}
	return nil, false
}

// StateResolver is an interface used to resolve state information
type StateResolver interface {
	GetObj(table string, keys map[string]string) (*state.Obj, error)
}

func NewMockStateResolver() *MockStateResolver {
	return &MockStateResolver{}
}

type MockStateResolver struct {
}

func (m *MockStateResolver) GetObj(table string, keys map[string]string) (*state.Obj, error) {
	return nil, nil
}

// Snapshot is a wrapper around the snapshot that does safe checks on the types
type Snapshot struct {
	// provider    *Provider
	// err         *ErrorEvent
	// block       uint64
	schemas     map[string]*schema.Table
	inmemStore  *inmemStore
	trackedObjs map[string]*Obj
	// errCh       chan struct{}
	resolver StateResolver
}

func (s *Snapshot) reset() {
	s.trackedObjs = map[string]*Obj{}
}

func (s *Snapshot) AddSchema(name string, t *schema.Table) {
	if s.schemas == nil {
		s.schemas = map[string]*schema.Table{}
	}
	s.schemas[name] = t
}

func NewSnapshot(resolver StateResolver) *Snapshot {
	return &Snapshot{
		inmemStore:  newInmemStore(),
		trackedObjs: map[string]*Obj{},
		// errCh:       make(chan struct{}),
		schemas:  map[string]*schema.Table{},
		resolver: resolver,
	}
}

func (s *Snapshot) decodeKeyAccess(tableName string, keyRaw ...interface{}) (map[string]string, string, error) {
	// validation
	table, ok := s.schemas[tableName]
	if !ok {
		return nil, "", fmt.Errorf("table name not found %s", tableName)
	}
	// get ids in order
	idFields := []*schema.Field{}
	for _, field := range table.Fields {
		if field.ID {
			idFields = append(idFields, field)
		}
	}
	if len(idFields) != len(keyRaw) {
		return nil, "", fmt.Errorf("bad size")
	}

	vals := []string{}
	keysMap := map[string]string{}
	for indx, field := range idFields {
		valStr, err := field.Encode(keyRaw[indx])
		if err != nil {
			return nil, "", err
		}
		vals = append(vals, valStr)
		keysMap[field.Name] = valStr
	}

	id := hex.EncodeToString(buildIndex(tableName, vals))
	return keysMap, id, nil
}

/*
func (s *Snapshot) GetOk(tableName string, keyRaw ...interface{}) (*Obj, bool, error) {
	_, idStr, err := s.decodeKeyAccess(tableName, keyRaw...)
	if err != nil {
		return nil, false, err
	}
	return s.getOk(tableName, idStr)
}
*/

/*
func (s *Snapshot) hasError() bool {
	return s.err != nil
}
*/

type ErrorEvent struct {
	Type        string
	Description string
	Err         error
}

const (
	ErrorEventSetEncode         = "ErrorSetEncode"
	ErrorEventGetDecode         = "ErrorGetDecode"
	ErrorEventContractInit      = "ErrorContractInit"
	ErrorEventRecoverObject     = "ErrorRecoverObject"
	ErrorEventFieldNotFound     = "ErrorFieldNotFound"
	ErrorEventFieldBadType      = "ErrorFieldBadType"
	ErrorEventSchemaNotFound    = "ErrorSchemaNotFound"
	ErrorEventIncorrectIdFields = "ErrorIncorrectIdFields"
	ErrorEventGeneric           = "ErrorEventGeneric"
)

/*
func (s *Snapshot) ErrCh() <-chan struct{} {
	return s.errCh
}
*/

/*
func (s *Snapshot) finish(evnt *ErrorEvent) {
	s.err = evnt
	close(s.errCh)

	fmt.Println(evnt.Description)
	fmt.Println(evnt.Err)

	panic("x")
	// runtime.Goexit()
}
*/

func (s *Snapshot) Get(tableName string, keyRaw ...interface{}) (*Obj, error) {
	keysMap, idStr, err := s.decodeKeyAccess(tableName, keyRaw...)
	if err != nil {
		return nil, err
	}

	obj, ok := s.getOk(tableName, idStr)
	if ok {
		return obj, nil
	}

	var dataObj *state.Obj

	// not found, try to search it on the resolver (if any)
	dataObj, err = s.resolver.GetObj(tableName, keysMap)
	if err != nil {
		/*
			s.finish(&ErrorEvent{
				Type: ErrorEventRecoverObject,
				Err:  err,
			})
		*/
		return nil, fmt.Errorf("bad 1")
	}

	table := s.schemas[tableName]

	if dataObj == nil {
		// create a new object
		obj = s.create(tableName, table, idStr, keysMap)
	} else {
		// derive the object
		obj = &Obj{
			schema:  table,
			id:      []byte(idStr),
			table:   tableName,
			key:     keysMap,
			vals:    map[string]string{},
			changes: map[string]string{},
		}
		// move all the non key values
		for k, v := range dataObj.Data {
			if _, ok := keysMap[k]; !ok {
				obj.vals[k] = v
			}
		}

		// add the object to cache
		s.inmemStore.add(idStr, obj.Copy())
	}

	// initialize the default values
	for _, field := range table.Fields {
		if field.Default != nil {
			// set will already handle the string conversion
			obj.Set(field.Name, field.Default)
		}
	}
	return obj, nil
}
