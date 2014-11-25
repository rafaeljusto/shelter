// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/testing/utils"
	"testing"
	"time"
)

func TestValidDSAlgorithm(t *testing.T) {
	if !IsValidDSAlgorithm(1) || !IsValidDSAlgorithm(8) || !IsValidDSAlgorithm(254) {
		t.Error("Not accepting valid DS algorithms")
	}

	if IsValidDSAlgorithm(0) || IsValidDSAlgorithm(255) {
		t.Error("Accepting invalid DS algorithms")
	}
}

func TestValidDSDigestType(t *testing.T) {
	if !IsValidDSDigestType(0) || !IsValidDSDigestType(3) || !IsValidDSDigestType(5) {
		t.Error("Not accepting valid DS digest type")
	}

	if IsValidDSDigestType(6) {
		t.Error("Accepting invalid digest type")
	}
}

func TestDSNormalize(t *testing.T) {
	data := []struct {
		description string
		request     DSRequest
		expected    DSRequest
	}{
		{
			description: "it should normalize DS digest",
			request: DSRequest{
				Digest: utils.NewString("   EaA0978f38879db70A53f9FF1ACF21D046a98B5C   "),
			},
			expected: DSRequest{
				Digest: utils.NewString("eaa0978f38879db70a53f9ff1acf21d046a98b5c"),
			},
		},
	}

	for i, item := range data {
		item.request.Normalize()

		if !utils.CompareUint16(item.request.Keytag, item.expected.Keytag) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Keytag,
				item.expected.Keytag,
			)
		}

		if !utils.CompareUint8(item.request.Algorithm, item.expected.Algorithm) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Algorithm,
				item.expected.Algorithm,
			)
		}

		if !utils.CompareStrings(item.request.Digest, item.expected.Digest) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Digest,
				item.expected.Digest,
			)
		}

		if !utils.CompareUint8(item.request.DigestType, item.expected.DigestType) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.DigestType,
				item.expected.DigestType,
			)
		}
	}
}

func TestDSValidate(t *testing.T) {
	data := []struct {
		description      string
		request          DSRequest
		expectedError    bool
		expectedMessages bool
	}{
		{
			description: "it should accept a valid DS",
			request: DSRequest{
				Keytag:     utils.NewUint16(41895),
				Algorithm:  utils.NewUint8(5),
				Digest:     utils.NewString("eaa0978f38879db70a53f9ff1acf21d046a98b5c"),
				DigestType: utils.NewUint8(1),
			},
			expectedError:    false,
			expectedMessages: false,
		},
		{
			description: "it should alert when keytag is nil",
			request: DSRequest{
				Keytag:     nil,
				Algorithm:  utils.NewUint8(5),
				Digest:     utils.NewString("eaa0978f38879db70a53f9ff1acf21d046a98b5c"),
				DigestType: utils.NewUint8(1),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			description: "it should alert when the DS has an invalid algorithm",
			request: DSRequest{
				Keytag:     utils.NewUint16(41895),
				Algorithm:  utils.NewUint8(0),
				Digest:     utils.NewString("eaa0978f38879db70a53f9ff1acf21d046a98b5c"),
				DigestType: utils.NewUint8(1),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			description: "it should alert when the DS has an invalid digest type",
			request: DSRequest{
				Keytag:     utils.NewUint16(41895),
				Algorithm:  utils.NewUint8(5),
				Digest:     utils.NewString("eaa0978f38879db70a53f9ff1acf21d046a98b5c"),
				DigestType: utils.NewUint8(0),
			},
			expectedError:    false,
			expectedMessages: true,
		},
	}

	for i, item := range data {
		messages, err := item.request.Validate()
		if err != nil && !item.expectedError {
			t.Errorf("Item %d, “%s”: did not expect error '%s'", i, item.description, err)

		} else if err == nil && item.expectedError {
			t.Errorf("Item %d, “%s”: expected error", i, item.description)
		}

		if messages != nil && !item.expectedMessages {
			t.Errorf("Item %d, “%s”: did not expect messages '%v'", i, item.description, messages)

		} else if messages == nil && item.expectedMessages {
			t.Errorf("Item %d, “%s”: expected messages", i, item.description)
		}
	}
}
