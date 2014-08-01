// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/model"
	"testing"
	"time"
)

func TestToDSModel(t *testing.T) {
	dsRequest := DSRequest{
		Keytag:     41674,
		Algorithm:  5,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 1,
	}

	ds, err := dsRequest.toDSModel()
	if err != nil {
		t.Fatal(err)
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
		Keytag:     41674,
		Algorithm:  0,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 1,
	}

	ds, err = dsRequest.toDSModel()
	if err == nil {
		t.Error("Allowing an invalid DS algorithm")
	}

	dsRequest = DSRequest{
		Keytag:     41674,
		Algorithm:  5,
		Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
		DigestType: 6,
	}

	ds, err = dsRequest.toDSModel()
	if err == nil {
		t.Error("Allowing an invalid DS digest type")
	}
}

func TestToDSSetModel(t *testing.T) {
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

	_, err := toDSSetModel(dsSetRequest)
	if err != nil {
		t.Error(err)
	}

	dsSetRequest = []DSRequest{
		{
			Keytag:     41674,
			Algorithm:  5,
			Digest:     "EAA0978F38879DB70A53F9FF1ACF21D046A98B5C",
			DigestType: 1,
		},
		{
			Keytag:     45966,
			Algorithm:  0,
			Digest:     "B7C0BDE8F3C90E573B956B14A14CAF5001A3E841",
			DigestType: 1,
		},
	}

	_, err = toDSSetModel(dsSetRequest)
	if err == nil {
		t.Error("Not verifying errors in DS set conversion")
	}
}

func TestToDSResponse(t *testing.T) {
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

	dsResponse := toDSResponse(ds)

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

func TestToDSSetResponse(t *testing.T) {
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

	dsSetResponse := toDSSetResponse(dsSet)
	if len(dsSetResponse) != 2 {
		t.Error("Fail to convert a DS set")
	}
}
