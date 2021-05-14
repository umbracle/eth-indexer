package sdk

import lru "github.com/hashicorp/golang-lru"

type inmemStore struct {
	cache *lru.Cache
	data  map[string]*Obj2
}

func newInmemStore() *inmemStore {
	cache, _ := lru.New(10000)
	return &inmemStore{
		cache: cache,
	}
}

func (i *inmemStore) get(k string) (*Obj2, bool) {
	v, ok := i.cache.Get(k)
	if !ok {
		return nil, false
	}
	return v.(*Obj2), true
}

func (i *inmemStore) add(k string, val *Obj2) {
	// it is assumed that this is locked
	i.cache.Add(k, val)
}
