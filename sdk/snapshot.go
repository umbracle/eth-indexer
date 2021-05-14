package sdk

import (
	"encoding/hex"
	"fmt"
	"runtime"

	protosdk "github.com/umbracle/eth-indexer/sdk/proto"
)

type objErr interface {
	finish(*ErrorEvent)
}

type Obj2 struct {
	schema *Table
	// TODO: provider interface to call finish
	objErr  objErr
	created bool
	id      []byte
	table   string
	key     map[string]string
	vals    map[string]string
	changes map[string]string
}

func (o *Obj2) IsNew() bool {
	return o.created
}

func (o *Obj2) getField(name string) *Field {
	for _, f := range o.schema.Fields {
		if f.Name == name {
			return f
		}
	}
	o.objErr.finish(&ErrorEvent{
		Type: ErrorEventFieldNotFound,
		Err:  fmt.Errorf("%s %s", o.table, name),
	})
	return nil
}

func (o *Obj2) isChanged() bool {
	return len(o.changes) != 0
}

func (o *Obj2) hasChanged(key string) (string, bool) {
	val, ok := o.changes[key]
	return val, ok
}

func (o *Obj2) Set(key string, val interface{}) {
	// convert val to a string representation
	field := o.getField(key)
	valStr, err := field.Encode(val)
	if err != nil {
		o.objErr.finish(&ErrorEvent{
			Type: ErrorEventSetEncode,
			Err:  fmt.Errorf("failed to encode schema %s %s with %s: %v", o.table, key, val, err),
		})
	}

	if raw, ok := o.vals[key]; ok {
		if raw != val {
			o.changes[key] = valStr
		}
	} else {
		o.changes[key] = valStr
	}
}

func (o *Obj2) Get(key string) interface{} {
	v, _ := o.GetOk(key)
	return v
}

func (o *Obj2) GetOk(key string) (interface{}, bool) {
	// make sure first the value exists
	field := o.getField(key)

	// try in changes first
	valStr, ok := o.changes[key]
	if !ok {
		// use the state value
		valStr, ok = o.vals[key]
		if !ok {
			return nil, false
		}
	}
	// convert val to his type
	val, err := field.Decode(valStr)
	if err != nil {
		o.objErr.finish(&ErrorEvent{
			Type: ErrorEventGetDecode,
			Err:  fmt.Errorf("failed to decode %s with %s: %v", key, val, err),
		})
	}
	return val, true
}

func (o *Obj2) expect(key string, fieldType FieldType) {
	field := o.getField(key)
	if field.Type != fieldType {
		o.objErr.finish(&ErrorEvent{
			Type: ErrorEventFieldBadType,
			Err:  fmt.Errorf("field %s expected %d but found %d", key, fieldType, field.Type),
		})
	}
}

func (o *Obj2) Incr(key string) {
	o.Add(key, uint64(1))
}

func (o *Obj2) Sub(key string, v interface{}) {
	var val interface{}

	switch obj := v.(type) {
	case *Float:
		o.expect(key, TypeDecimal)
		val = o.Get(key).(*Float).Sub(obj)

	case uint64:
		o.expect(key, TypeUint)
		val = o.Get(key).(uint64) - obj
	}
	o.Set(key, val)
}

func (o *Obj2) Add(key string, v interface{}) {
	var val interface{}

	switch obj := v.(type) {
	case *Float:
		o.expect(key, TypeDecimal)
		val = o.Get(key).(*Float).Add(obj)

	case uint64:
		o.expect(key, TypeUint)
		val = o.Get(key).(uint64) + obj

	default:
		panic("Not expected")
	}
	o.Set(key, val)
}

/*
func (o *Obj2) getNum(key string) *big.Int {
	current := new(big.Int)
	past, ok := o.getOk(key)
	if ok {
		current.SetString(past, 10)
	}
	return current
}

func (o *Obj2) div(key string, v *big.Int) {
	val := o.getNum(key)
	val.Div(val, v)
	o.set(key, val.String())
}

func (o *Obj2) mul(key string, v *big.Int) {
	val := o.getNum(key)
	val.Mul(val, v)
	o.set(key, val.String())
}
*/

func (o *Obj2) Copy() *Obj2 {
	oo := new(Obj2)
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

func (s *Snapshot) save() []*protosdk.Diff {
	diffs := []*protosdk.Diff{}

	for _, obj := range s.trackedObjs {
		if obj.isChanged() {
			obj2 := obj.Copy()

			diff := &protosdk.Diff{
				Table:    obj.table,
				Keys:     obj.key,
				Creation: obj.created,
				Vals:     obj.changes,
			}
			diffs = append(diffs, diff)

			for k, v := range obj.changes {
				obj2.vals[k] = v
			}

			// reset the object
			obj2.changes = map[string]string{}
			obj2.created = false

			s.inmemStore.add(string(obj2.id), obj2)
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

func (s *Snapshot) create(table string, schema *Table, id string, keys map[string]string) *Obj2 {
	obj := &Obj2{
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

func (s *Snapshot) getOk(table string, id string) (*Obj2, bool) {
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

// Snapshot is a wrapper around the snapshot that does safe checks on the types
type Snapshot struct {
	provider    *Provider
	err         *ErrorEvent
	block       uint64
	schemas     map[string]*Table
	inmemStore  *inmemStore
	trackedObjs map[string]*Obj2
}

func (s *Snapshot) reset() {
	s.trackedObjs = map[string]*Obj2{}
}

func newSnapshot() *Snapshot {
	return &Snapshot{
		inmemStore:  newInmemStore(),
		trackedObjs: map[string]*Obj2{},
	}
}

func (s *Snapshot) decodeKeyAccess(tableName string, keyRaw ...interface{}) (map[string]string, string) {
	// validation
	table, ok := s.schemas[tableName]
	if !ok {
		s.finish(&ErrorEvent{
			Type: ErrorEventSchemaNotFound,
			Err:  fmt.Errorf(tableName),
		})
	}
	// get ids in order
	idFields := []*Field{}
	for _, field := range table.Fields {
		if field.ID {
			idFields = append(idFields, field)
		}
	}
	if len(idFields) != len(keyRaw) {
		s.finish(&ErrorEvent{
			Type: ErrorEventIncorrectIdFields,
		})
	}

	vals := []string{}
	keysMap := map[string]string{}
	for indx, field := range idFields {
		valStr, err := field.Encode(keyRaw[indx])
		if err != nil {
			s.finish(&ErrorEvent{
				Type: ErrorEventGeneric,
				Err:  err,
			})
		}
		vals = append(vals, valStr)
		keysMap[field.Name] = valStr
	}

	id := hex.EncodeToString(buildIndex(tableName, vals))
	return keysMap, id
}

func (s *Snapshot) GetOk(tableName string, keyRaw ...interface{}) (*Obj2, bool) {
	_, idStr := s.decodeKeyAccess(tableName, keyRaw...)

	return s.getOk(tableName, idStr)
}

func (s *Snapshot) hasError() bool {
	return s.err != nil
}

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

func (s *Snapshot) finish(evnt *ErrorEvent) {
	s.err = evnt
	runtime.Goexit()
}

func (s *Snapshot) Get(tableName string, keyRaw ...interface{}) *Obj2 {
	keysMap, idStr := s.decodeKeyAccess(tableName, keyRaw...)

	obj, ok := s.getOk(tableName, idStr)
	if ok {
		return obj
	}

	var dataObj *Obj
	if s.provider.resolver != nil {
		// not found, try to search it still on the resolver (if any)
		raw, err := s.provider.resolver.GetObj2(tableName, keysMap)
		if err != nil {
			s.finish(&ErrorEvent{
				Type: ErrorEventRecoverObject,
				Err:  err,
			})
		}
		dataObj = raw
	}

	table := s.schemas[tableName]

	if dataObj == nil {
		// create a new object
		obj = s.create(tableName, table, idStr, keysMap)
	} else {
		// derive the object
		obj = &Obj2{
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

	// pass a reference so that the object can call errors
	obj.objErr = s

	// initialize the default values
	for _, field := range table.Fields {
		if field.Default != nil {
			// set will already handle the string conversion
			obj.Set(field.Name, field.Default)
		}
	}

	idFields := table.getIDS()

	if len(idFields) == 1 {
		// initialize only if there is one id since this is meant to be used
		// only for external contracts that need to call information
		idField := idFields[0]

		keyID, err := idField.Encode(keyRaw[0])
		if err != nil {
			panic(err)
		}
		// call init function
		if err := s.provider.initSchema(tableName, keyID, obj); err != nil {
			s.finish(&ErrorEvent{
				Type: ErrorEventContractInit,
				Err:  err,
			})
		}
	}

	return obj
}
