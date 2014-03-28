// scan - Scan service
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

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

	// This test is really fast, but there's a little chance to fail when the epoch from the
	// created object is different from the current epoch. This will happen if we are
	// creating the object in the end of a second. After some tests using shell scripts a
	// saw that this scenario is not frequent, so I will stay with this strategy
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
	querierCache.hosts = make(map[string]*hostCache)

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

	_, err = querierCache.Get("localhost")
	if err != nil {
		t.Fatal("Not recovering correctly from local cache")
	}

	h, exists := querierCache.hosts["localhost"]

	if !exists {
		t.Fatal("Not storing results into cache")
	}

	if len(h.addresses) != len(addresses) {
		t.Error("Storing different data from the returned one")
	}

	// The tests bellow are really fast, but there's a little chance to fail when the epoch
	// from the created object is different from the current epoch. This will happen if we
	// are creating the object in the end of a second. After some tests using shell scripts
	// a saw that this scenario is not frequent, so I will stay with this strategy

	h.lastEpoch = time.Now().Unix()
	h.queriesPerSecond = MaxQPSPerHost + 1
	h.timeouts = 0
	querierCache.hosts["localhost"] = h

	addresses, err = querierCache.Get("localhost")
	if err != ErrHostQPSExceeded {
		t.Error("Not returning error when maximum QPS per host is exceeded")
	}

	h.lastEpoch = 0
	h.queriesPerSecond = 0
	h.timeouts = maxTimeoutsPerHost + 1
	querierCache.hosts["localhost"] = h

	addresses, err = querierCache.Get("localhost")
	if err != ErrHostTimeout {
		t.Error("Not returning error when maximum timeouts in the host is exceeded")
	}

	_, err = querierCache.Get("abc123idontexist321cba.com.br")
	if err == nil {
		t.Error("Resolving an unknown name")
	}
}

func TestQuerierCacheTimeout(t *testing.T) {
	querierCache.hosts = make(map[string]*hostCache)

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

func TestQuerierCacheQuery(t *testing.T) {
	querierCache.hosts = make(map[string]*hostCache)

	// The tests bellow are really fast, but there's a little chance to fail when the epoch
	// from the created object is different from the current epoch. This will happen if we
	// are creating the object in the end of a second. After some tests using shell scripts
	// a saw that this scenario is not frequent, so I will stay with this strategy

	querierCache.Query("localhost")

	if _, exists := querierCache.hosts["localhost"]; exists {
		t.Error("Creating cache entry when alerting about a query")
	}

	_, err := querierCache.Get("localhost")
	if err != nil {
		t.Fatal("Not resolving a valid nameserver")
	}

	querierCache.Query("localhost")

	h, exists := querierCache.hosts["localhost"]

	if !exists {
		t.Fatal("Not creating cache entries when necessary")
	}

	if h.lastEpoch != time.Now().Unix() {
		t.Error("Not setting epoch correctly")
	}

	if h.queriesPerSecond != 1 {
		t.Error("Not counting QPS correctly")
	}

	querierCache.Query("localhost")

	h, exists = querierCache.hosts["localhost"]

	if !exists {
		t.Fatal("Removing cache entry in an awkward moment")
	}

	if h.queriesPerSecond != 2 {
		t.Error("Not counting QPS correctly")
	}

	// Forcing epoch change
	time.Sleep(1 * time.Second)

	querierCache.Query("localhost")

	if !exists {
		t.Fatal("Removing cache entry in an awkward moment")
	}

	if h.lastEpoch != time.Now().Unix() {
		t.Error("Not replacing epoch correctly")
	}

	if h.queriesPerSecond != 1 {
		t.Error("Not counting QPS correctly")
	}
}

func TestQuerierCacheClear(t *testing.T) {
	querierCache.hosts = make(map[string]*hostCache)

	_, err := querierCache.Get("localhost")
	if err != nil {
		t.Fatal("Not resolving a valid nameserver")
	}

	querierCache.Clear()

	if querierCache.hosts == nil || len(querierCache.hosts) > 0 {
		t.Error("Not clearing the cache correctly")
	}
}
