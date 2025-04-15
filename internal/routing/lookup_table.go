package routing

import "sync"

type LookupTable struct {
	addressees map[string]string
	mutex      *sync.RWMutex
}

func NewLookupTable() *LookupTable {
	return &LookupTable{
		addressees: make(map[string]string),
		mutex:      &sync.RWMutex{},
	}
}

func (lt *LookupTable) AddressLookup(hostname string) (bool, string) {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()

	serverAddress := lt.addressees[hostname]
	if serverAddress == "" {
		return false, ""
	}

	return true, serverAddress
}

func (lt *LookupTable) UpdateLookupTable(table map[string]string) {
	lt.mutex.Lock()
	defer lt.mutex.Unlock()

	lt.addressees = table
}
