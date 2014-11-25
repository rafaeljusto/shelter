// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/testing/utils"
	"testing"
)

func TestDNSKEYNormalize(t *testing.T) {
	data := []struct {
		description string
		request     DNSKEYRequest
		expected    DNSKEYRequest
	}{
		{
			description: "it should spaces from the public key",
			request: DNSKEYRequest{
				PublicKey: utils.NewString("ABCD eFG HIJ LMN OpQ"),
			},
			expected: DNSKEYRequest{
				PublicKey: utils.NewString("ABCDeFGHIJLMNOpQ"),
			},
		},
	}

	for i, item := range data {
		item.request.Normalize()

		if !utils.CompareUint8(item.request.Algorithm, item.expected.Algorithm) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Algorithm,
				item.expected.Algorithm,
			)
		}

		if !utils.CompareUint16(item.request.Flags, item.expected.Flags) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Flags,
				item.expected.Flags,
			)
		}

		if !utils.CompareStrings(item.request.PublicKey, item.expected.PublicKey) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.PublicKey,
				item.expected.PublicKey,
			)
		}
	}
}

func TestDNSKEYValidate(t *testing.T) {
	data := []struct {
		description      string
		request          DNSKEYRequest
		expectedError    bool
		expectedMessages bool
	}{
		{
			description: "it should spaces from the public key",
			request: DNSKEYRequest{
				Flags:     utils.NewUint16(257),
				Algorithm: utils.NewUint8(5),
				PublicKey: utils.NewString(`
					AwEAAcPXy7f7HcMkvXjJhhpbku3d28zL/9AXwfCXgza4
					yj1I4aW6DZLOjUzz13x1AtkWLLWBb+0w88rqvEQ7FR1+
					HSTvVrmkjIPh9LsSjGuXsLdqGqm9sgb2DwzR4SDRBqZk
					O2hwCDy+/Q/FPIsRn/cKwU8bS9I5ASTQ+E+obGqdb/k7
					BVLGlKLpoe5GyqEqCfJO5Q==`),
			},
			expectedError:    false,
			expectedMessages: false,
		},
		{
			description: "it should detect when flags is not informed",
			request: DNSKEYRequest{
				Flags:     nil,
				Algorithm: utils.NewUint8(5),
				PublicKey: utils.NewString(`
					AwEAAcPXy7f7HcMkvXjJhhpbku3d28zL/9AXwfCXgza4
					yj1I4aW6DZLOjUzz13x1AtkWLLWBb+0w88rqvEQ7FR1+
					HSTvVrmkjIPh9LsSjGuXsLdqGqm9sgb2DwzR4SDRBqZk
					O2hwCDy+/Q/FPIsRn/cKwU8bS9I5ASTQ+E+obGqdb/k7
					BVLGlKLpoe5GyqEqCfJO5Q==`),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			description: "it should detect when algorithm is not informed",
			request: DNSKEYRequest{
				Flags:     utils.NewUint16(257),
				Algorithm: nil,
				PublicKey: utils.NewString(`
					AwEAAcPXy7f7HcMkvXjJhhpbku3d28zL/9AXwfCXgza4
					yj1I4aW6DZLOjUzz13x1AtkWLLWBb+0w88rqvEQ7FR1+
					HSTvVrmkjIPh9LsSjGuXsLdqGqm9sgb2DwzR4SDRBqZk
					O2hwCDy+/Q/FPIsRn/cKwU8bS9I5ASTQ+E+obGqdb/k7
					BVLGlKLpoe5GyqEqCfJO5Q==`),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			description: "it should detect when algorithm is invalid",
			request: DNSKEYRequest{
				Flags:     utils.NewUint16(257),
				Algorithm: utils.NewUint8(0),
				PublicKey: utils.NewString(`
					AwEAAcPXy7f7HcMkvXjJhhpbku3d28zL/9AXwfCXgza4
					yj1I4aW6DZLOjUzz13x1AtkWLLWBb+0w88rqvEQ7FR1+
					HSTvVrmkjIPh9LsSjGuXsLdqGqm9sgb2DwzR4SDRBqZk
					O2hwCDy+/Q/FPIsRn/cKwU8bS9I5ASTQ+E+obGqdb/k7
					BVLGlKLpoe5GyqEqCfJO5Q==`),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			description: "it should detect when public key is not informed",
			request: DNSKEYRequest{
				Flags:     utils.NewUint16(257),
				Algorithm: utils.NewUint8(5),
				PublicKey: nil,
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
