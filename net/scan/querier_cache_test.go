package scan

import (
	"testing"
	"time"
)

func TestMaxTimeoutsPerHost(t *testing.T) {
	h := hostCache{
		timeouts: maxTimeoutsPerHost + 1,
	}

	if !h.timeoutsPerHostExceeded() {
		t.Error("Not checking when maximum number of timeouts exceeded")
	}

	h = hostCache{
		timeouts: maxTimeoutsPerHost - 1,
	}

	if h.timeoutsPerHostExceeded() {
		t.Error("Alerting maximum number of timeouts when it's not")
	}
}

func TestQueriesPerSecondExceeded(t *testing.T) {
	h := hostCache{
		lastEpoch:        time.Now().Unix(),
		queriesPerSecond: MaxQPSPerHost + 1,
	}

	// This test is really fast, but there's a little chance to fail when the epoch from the created
	// object is different from the current epoch. This will happen if we are creating the object in
	// the end of a second. After some tests using shell scripts a saw that this scenario is not
	// frequent, so I will stay with this strategy
	if !h.queriesPerSecondExceeded() {
		t.Error("Not checking when QPS per host exceeded")
	}

	h = hostCache{
		lastEpoch:        time.Now().Unix(),
		queriesPerSecond: MaxQPSPerHost - 1,
	}

	if h.queriesPerSecondExceeded() {
		t.Error("Alerting maximum QPS per host when it's not")
	}

	h = hostCache{
		lastEpoch:        time.Now().Unix(),
		queriesPerSecond: MaxQPSPerHost + 1,
	}

	MaxQPSPerHost = 0
	if h.queriesPerSecondExceeded() {
		t.Error("Not working with disabled QPS per host feature")
	}
	MaxQPSPerHost = 500
}

func TestQuerierCacheGet(t *testing.T) {
	addresses, err := querierCache.Get("localhost")
	if err != nil {
		t.Fatal("Not resolving a valid nameserver")
	}

	// Localhost can have IPv4 and IPv6 addresses
	if len(addresses) == 0 || len(addresses) > 2 {
		t.Fatal("Something wrong with the number of IP address obtained from localhost")
	}

	if addresses[0].String() != "127.0.0.1" &&
		addresses[0].String() != "::1" {
		t.Error("Not resolving correctly the localhost")
	}

	if len(addresses) == 2 {
		if addresses[1].String() != "127.0.0.1" &&
			addresses[1].String() != "::1" {
			t.Error("Not resolving correctly the localhost")
		}
	}

	h, exists := querierCache.hosts["localhost"]

	if !exists {
		t.Fatal("Not storing results into cache")
	}

	if len(h.addresses) != len(addresses) {
		t.Error("Storing different data from the returned one")
	}
}

func TestQuerierCacheTimeout(t *testing.T) {
	_, err := querierCache.Get("localhost")
	if err != nil {
		t.Fatal("Not resolving a valid nameserver")
	}

	querierCache.Timeout("localhost")

	h, exists := querierCache.hosts["localhost"]

	if !exists {
		t.Fatal("Not storing results into cache")
	}

	if h.timeouts != 1 {
		t.Error("Not working well with timeouts counter")
	}
}
