// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"testing"
	"time"
)

func TestDSChangeStatus(t *testing.T) {
	timeMark := time.Now()

	ds := DS{
		LastStatus:  DSStatusDNSError,
		LastCheckAt: timeMark,
	}

	// Wait for a while to make a difference between the first and the second
	// LastCheck fields
	time.Sleep(1 * time.Millisecond)

	ds.ChangeStatus(DSStatusOK)

	if ds.LastStatus != DSStatusOK {
		t.Error("ChangeStatus method did not change DS attribute")
	}

	if ds.LastCheckAt.Before(timeMark) || ds.LastCheckAt.Equal(timeMark) {
		t.Error("ChangeStatus method did not update the last check date")
	}
}

func TestDSStatusToString(t *testing.T) {
	if DSStatusToString(DSStatusNotChecked) != "NOTCHECKED" {
		t.Error("DS status NOTCHECKED not converting correctly to string")
	}

	if DSStatusToString(DSStatusOK) != "OK" {
		t.Error("DS status OK not converting correctly to string")
	}

	if DSStatusToString(DSStatusTimeout) != "TIMEOUT" {
		t.Error("DS status TIMEOUT not converting correctly to string")
	}

	if DSStatusToString(DSStatusNoSignature) != "NOSIG" {
		t.Error("DS status NOSIG not converting correctly to string")
	}

	if DSStatusToString(DSStatusExpiredSignature) != "EXPSIG" {
		t.Error("DS status EXPSIG not converting correctly to string")
	}

	if DSStatusToString(DSStatusNoKey) != "NOKEY" {
		t.Error("DS status NOKEY not converting correctly to string")
	}

	if DSStatusToString(DSStatusNoSEP) != "NOSEP" {
		t.Error("DS status NOSEP not converting correctly to string")
	}

	if DSStatusToString(DSStatusSignatureError) != "SIGERR" {
		t.Error("DS status SIGERR not converting correctly to string")
	}

	if DSStatusToString(DSStatusDNSError) != "DNSERR" {
		t.Error("DS status DNSERR not converting correctly to string")
	}

	if DSStatusToString(999999) != "" {
		t.Error("Unknown DS status associated to some existing status")
	}
}

func TestDSApplyDNSKEY(t *testing.T) {
	dnskeyRequest := DNSKEYRequest{
		FQDN:      "br.",
		Flags:     257,
		Algorithm: uint8(5),
		PublicKey: "AwEAAblaEaapG4inrQASY3HzwXwBaRSy5mkj7mZ30F+huI7zL8g0U7dv 7ufnSEQUlsC57OHoTBza+TQIv/mgQed8Fy4XGCGzYiHSYVYvGO9iWG3O 0voBYy/zv0z7ANfrA7Z3lY51CI6m/qoZUcDlNM0yTcJgilaKwUkLBHMA p9NJPuKVt8A7OHab00r2RDEVjiLWIIuTbz74gCXOVfAmvW07c8c=",
	}

	var ds DS
	if !ds.ApplyDNSKEY(dnskeyRequest) {
		t.Fatal("Error converting DNSKEY request to DS model")
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
		FQDN:      "br.",
		Flags:     257,
		Algorithm: uint8(5),
		PublicKey: strings.Repeat("x", 65536) + "\uff00",
	}

	if ds.ApplyDNSKEY(dnskeyRequest) {
		t.Error("Not detecting an invalid DNSKEY")
	}
}

func TestDSApply(t *testing.T) {
	dsRequest := DSRequest{
		Keytag:     41674,
		Algorithm:  5,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 1,
	}

	var ds DS
	if !ds.Apply(dsRequest) {
		t.Fatal("Error applying a valid DS request")
	}

	if ds.Keytag != 41674 {
		t.Error("Not keeping keytag in conversion")
	}

	if ds.Algorithm != model.DSAlgorithmRSASHA1 {
		t.Error("Not converting DS algorithm correctly")
	}

	if ds.Digest != "eaa0978f38879db70a53f9ff1acf21d046a98b5c" {
		t.Error("Not keeping digest in conversion")
	}

	if ds.DigestType != model.DSDigestTypeSHA1 {
		t.Error("Not converting DS digest type correctly")
	}

	dsRequest = DSRequest{
		Keytag:     nil,
		Algorithm:  0,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 1,
	}

	if ds.Apply(dsRequest) {
		t.Error("Allowing an invalid DS keytag")
	}
}

func TestDSProtocol(t *testing.T) {
	now := time.Now()

	ds := model.DS{
		Keytag:      41674,
		Algorithm:   model.DSAlgorithmRSASHA1,
		Digest:      "eaa0978f38879db70a53f9ff1acf21d046a98b5c",
		DigestType:  model.DSDigestTypeSHA1,
		LastStatus:  model.DSStatusOK,
		LastCheckAt: now,
		LastOKAt:    now,
	}

	dsResponse := ds.Protocol()

	if dsResponse.Keytag != 41674 {
		t.Error("Fail to convert keytag")
	}

	if dsResponse.Algorithm != 5 {
		t.Error("Fail to convert algorithm")
	}

	if dsResponse.Digest != "eaa0978f38879db70a53f9ff1acf21d046a98b5c" {
		t.Error("Fail to convert digest")
	}

	if dsResponse.DigestType != 1 {
		t.Error("Fail to convert digest type")
	}

	if dsResponse.LastStatus != model.DSStatusToString(model.DSStatusOK) {
		t.Error("Fail to convert last status")
	}

	if dsResponse.LastCheckAt.Unix() != now.Unix() ||
		dsResponse.LastOKAt.Unix() != now.Unix() {

		t.Error("Fail to convert dates")
	}
}

func TestDSSetApplyDNSKEYs(t *testing.T) {
	dnskeysRequests := []DNSKEYRequest{
		{
			FQDN:      "br.",
			Flags:     257,
			Algorithm: uint8(5),
			PublicKey: "AwEAAblaEaapG4inrQASY3HzwXwBaRSy5mkj7mZ30F+huI7zL8g0U7dv 7ufnSEQUlsC57OHoTBza+TQIv/mgQed8Fy4XGCGzYiHSYVYvGO9iWG3O 0voBYy/zv0z7ANfrA7Z3lY51CI6m/qoZUcDlNM0yTcJgilaKwUkLBHMA p9NJPuKVt8A7OHab00r2RDEVjiLWIIuTbz74gCXOVfAmvW07c8c=",
		},
		{
			FQDN:      "br.",
			Flags:     256,
			Algorithm: uint8(5),
			PublicKey: "AwEAAfFQjspE7NgjAPclHrlyVFPRUHrU1p1U6POUXDpuIg8grg/s0lG1 8sjMkpxIvecIePLJw24gx48Ta9g0JJzPy35oGX5rYVJAu9BPqdUEuwIN ScTy3fPUhubvXP2fbyS6LeKNX/ZenihCD4HrViZehJmsKKv5fX8qx+RL 7NXCAAM1Xdet13cqR3LduW6wBzMiaQ==",
		},
	}

	var dsset DSSet
	if !dsset.ApplyDNSKEYs(dnskeysRequest) {
		t.Error("Not converting DNSKEYs to DS set correctly")
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

	if dsset.ApplyDNSKEYs(dnskeysRequest) {
		t.Error("Not verifying errors in DNSKEYS conversion")
	}
}

func TestDSSetApply(t *testing.T) {
	dsSetRequest := []DSRequest{
		{
			Keytag:     41674,
			Algorithm:  5,
			Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
			DigestType: 1,
		},
		{
			Keytag:     45966,
			Algorithm:  7,
			Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
			DigestType: 1,
		},
	}

	var dsset DSSet
	if !dsset.Apply(dsSetRequest) {
		t.Error("Not accepting a valid DS set ")
	}

	dsSetRequest = []DSRequest{
		{
			Keytag:     41674,
			Algorithm:  5,
			Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
			DigestType: 1,
		},
		{
			Keytag:     0,
			Algorithm:  5,
			Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
			DigestType: 1,
		},
	}

	if dsset.Apply(dsSetRequest) {
		t.Error("Not verifying errors in DS set conversion")
	}
}

func TestDSSetProtocol(t *testing.T) {
	now := time.Now()

	dsSet := []model.DS{
		{
			Keytag:      41674,
			Algorithm:   model.DSAlgorithmRSASHA1,
			Digest:      "eaa0978f38879db70a53f9ff1acf21d046a98b5c",
			DigestType:  model.DSDigestTypeSHA1,
			LastStatus:  model.DSStatusOK,
			LastCheckAt: now,
			LastOKAt:    now,
		},
		{
			Keytag:      45966,
			Algorithm:   model.DSAlgorithmRSASHA1NSEC3,
			Digest:      "b7c0bde8f3c90e573b956b14a14caf5001a3e841",
			DigestType:  model.DSDigestTypeSHA1,
			LastStatus:  model.DSStatusTimeout,
			LastCheckAt: now,
		},
	}

	dsSetResponse := dsSet.Protocol()
	if len(dsSetResponse) != 2 {
		t.Error("Fail to convert a DS set")
	}
}
