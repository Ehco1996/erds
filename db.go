package erds

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/tidwall/redcon"
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
	db := &DB{dict: make(map[string][]byte)}
	db.readAndLoadFromAof()
	return db
}

func (db *DB) execCmd(cmd redcon.Command) {
	f := string(cmd.Args[0])
	switch f {
	case "set":
		if len(cmd.Args) == 3 {
			db.set(string(cmd.Args[1]), cmd.Args[2])
		}
	case "del":
		if len(cmd.Args) == 2 {
			db.delete(string(cmd.Args[1]))
		}
	default:
		return
	}
}

func (db *DB) readAndLoadFromAof() {
	file, err := os.Open("erds.aof")
	defer file.Close()

	if err != nil {
		log.Printf("Can't open file. Error:'%s'", err.Error())
		return
	}
	// aof_parser.
	reader := newAofReader(file)
	for {
		resp, err := reader.readOneRESP()
		if err == io.EOF {
			break
		}
		cmd, err := redcon.Parse(resp)
		if err != nil {
			log.Fatal(err)
		}
		db.execCmd(cmd)
	}
}
