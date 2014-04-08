// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/model"
	"strings"
	"testing"
)

func TestDNSKEYRequestToDSModel(t *testing.T) {
	dnskeyRequest := DNSKEYRequest{
		Flags:     257,
		Algorithm: uint8(5),
		PublicKey: "AwEAAblaEaapG4inrQASY3HzwXwBaRSy5mkj7mZ30F+huI7zL8g0U7dv 7ufnSEQUlsC57OHoTBza+TQIv/mgQed8Fy4XGCGzYiHSYVYvGO9iWG3O 0voBYy/zv0z7ANfrA7Z3lY51CI6m/qoZUcDlNM0yTcJgilaKwUkLBHMA p9NJPuKVt8A7OHab00r2RDEVjiLWIIuTbz74gCXOVfAmvW07c8c=",
	}

	ds, err := dnskeyRequest.toDSModel("br.")
	if err != nil {
		t.Fatal(err)
	}

	if ds.Algorithm != model.DSAlgorithmRSASHA1 {
		t.Error("Not setting the algorithm properly")
	}

	if ds.Keytag != 41674 {
		t.Errorf("Not setting the keytag properly. Expected 41674 and got %d", ds.Keytag)
	}

	if ds.Digest != "6ec74914376b4f383ede3840088ae1d7bf13a19bfc51465cc2da57618889416a" {
		t.Errorf("Not converting to the correct digest. Expected "+
			"6ec74914376b4f383ede3840088ae1d7bf13a19bfc51465cc2da57618889416a and got %s",
			ds.Digest)
	}

	if ds.DigestType != model.DSDigestTypeSHA256 {
		t.Error("Not setting the digest type properly")
	}

	dnskeyRequest = DNSKEYRequest{
		Flags:     257,
		Algorithm: uint8(5),
		PublicKey: strings.Repeat("x", 65536) + "\uff00",
	}

	_, err = dnskeyRequest.toDSModel("br.")
	if err == nil {
		t.Error("Not detecting an invalid DNSKEY")
	}
}

func TestDNSKEYSRequestsToDSSetModel(t *testing.T) {
	dnskeysRequests := []DNSKEYRequest{
		{
			Flags:     257,
			Algorithm: uint8(5),
			PublicKey: "AwEAAblaEaapG4inrQASY3HzwXwBaRSy5mkj7mZ30F+huI7zL8g0U7dv 7ufnSEQUlsC57OHoTBza+TQIv/mgQed8Fy4XGCGzYiHSYVYvGO9iWG3O 0voBYy/zv0z7ANfrA7Z3lY51CI6m/qoZUcDlNM0yTcJgilaKwUkLBHMA p9NJPuKVt8A7OHab00r2RDEVjiLWIIuTbz74gCXOVfAmvW07c8c=",
		},
		{
			Flags:     256,
			Algorithm: uint8(5),
			PublicKey: "AwEAAfFQjspE7NgjAPclHrlyVFPRUHrU1p1U6POUXDpuIg8grg/s0lG1 8sjMkpxIvecIePLJw24gx48Ta9g0JJzPy35oGX5rYVJAu9BPqdUEuwIN ScTy3fPUhubvXP2fbyS6LeKNX/ZenihCD4HrViZehJmsKKv5fX8qx+RL 7NXCAAM1Xdet13cqR3LduW6wBzMiaQ==",
		},
	}

	_, err := dnskeysRequestsToDSSetModel("br.", dnskeysRequests)
	if err != nil {
		t.Error(err)
	}

	dnskeysRequests = []DNSKEYRequest{
		{
			Flags:     257,
			Algorithm: uint8(5),
			PublicKey: "AwEAAblaEaapG4inrQASY3HzwXwBaRSy5mkj7mZ30F+huI7zL8g0U7dv 7ufnSEQUlsC57OHoTBza+TQIv/mgQed8Fy4XGCGzYiHSYVYvGO9iWG3O 0voBYy/zv0z7ANfrA7Z3lY51CI6m/qoZUcDlNM0yTcJgilaKwUkLBHMA p9NJPuKVt8A7OHab00r2RDEVjiLWIIuTbz74gCXOVfAmvW07c8c=",
		},
		{
			Flags:     256,
			Algorithm: uint8(5),
			PublicKey: strings.Repeat("x", 65536) + "\uff00",
		},
	}

	_, err = dnskeysRequestsToDSSetModel("br.", dnskeysRequests)
	if err == nil {
		t.Error("Not verifying errors in DNSKEYS conversion")
	}
}
