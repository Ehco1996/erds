package erds

import (
	"sync"
)

// DB Global Hash DB
type DB struct {
	dict map[string][]byte
	lock sync.RWMutex
}

// get get k from dict,use read lock to make it threadsafe
func (db *DB) get(k string) ([]byte, bool) {
	db.lock.RLock()
	val, ok := db.dict[k]
	db.lock.RUnlock()
	return val, ok
}

// set get k from dict
func (db *DB) set(k string, val []byte) {
	db.lock.Lock()
	db.dict[k] = val
	db.lock.Unlock()
}

// delete get k from dict
func (db *DB) delete(k string) bool {
	db.lock.Lock()
	_, ok := db.dict[k]
	delete(db.dict, k)
	db.lock.Unlock()
	return ok
}

func initDb() *DB {
	return &DB{dict: make(map[string][]byte)}
}
