package scan

import (
	"net"
	"sync"
)

var (
	// Global variable used by all queriers (go routines) to access the cache
	querierCache QuerierCache
)

func init() {
	querierCache = QuerierCache{
		hosts: make(map[string][]net.IP),
	}
}

// QuerierCache was created to make the name resolution faster. Many domains use ISP
// nameservers, so if we cache the nameserver addres we are speeding up many domains scan
type QuerierCache struct {
	hosts      map[string][]net.IP // key-value structure that store the addresses
	hostsMutex sync.RWMutex        // Lock to allow concurrent access
}

// Method used to retrieve addresses of a given nameserver, if the address does not exist
// in the local cache the system will lookup for the domain and will store the result
func (q QuerierCache) Get(name string) ([]net.IP, error) {
	q.hostsMutex.RLock()
	addresses, found := q.hosts[name]
	q.hostsMutex.RUnlock()

	if found {
		return addresses, nil
	}

	// Not found in cache, lets discover the address of this name sending DNS requests
	addresses, err := net.LookupIP(name)
	if err != nil {
		return nil, err
	}

	q.hostsMutex.Lock()
	q.hosts[name] = addresses
	q.hostsMutex.Unlock()

	return addresses, nil
}
