package answer

import (
	"log"
	"sync"
)

var lock sync.RWMutex

// TypeA ...
var TypeA map[string]string

func init() {
	TypeA = make(map[string]string, 0)
}

// AddTypeA ...
func AddTypeA(dom, ip string) {
	lock.Lock()
	defer lock.Unlock()
	TypeA[dom] = ip
}

// DelTypeA ...
func DelTypeA(dom string) {
	lock.Lock()
	defer lock.Unlock()
	delete(TypeA, dom)
}

// DumpTypeA ...
func DumpTypeA() {
	lock.RLock()
	defer lock.RUnlock()
	log.Println(TypeA)
}

// GetTypeA ...
func GetTypeA(name string) (string, bool) {
	lock.RLock()
	defer lock.RUnlock()
	v, ok := TypeA[name]
	return v, ok
}
