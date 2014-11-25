// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/testing/utils"
	"net"
	"testing"
)

func TestNameserverNormalize(t *testing.T) {
	data := []struct {
		description string
		request     NameserverRequest
		expected    NameserverRequest
	}{
		{
			description: "it should normalize host with different case and spaces",
			request: NameserverRequest{
				Host: utils.NewString("   Ns1.EXAmPlE.cOM   "),
			},
			expected: NameserverRequest{
				Host: utils.NewString("ns1.example.com."),
			},
		},
		{
			description: "it should normalize host with IDNA",
			request: NameserverRequest{
				Host: utils.NewString("ns1.itaú.com.br"),
			},
			expected: NameserverRequest{
				Host: utils.NewString("ns1.xn--ita-boa.com.br."),
			},
		},
	}

	for i, item := range data {
		item.expected.Normalize()

		if !utils.CompareStrings(item.request.Host, item.expected.Host) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Host,
				item.expected.Host,
			)
		}

		if !utils.CompareIPs(item.request.IPv4, item.expected.IPv4) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.IPv4,
				item.expected.IPv4,
			)
		}

		if !utils.CompareIPs(item.request.IPv6, item.expected.IPv6) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.IPv6,
				item.expected.IPv6,
			)
		}
	}
}

func TestNameserverValidate(t *testing.T) {
	data := []struct {
		description      string
		request          NameserverRequest
		expectedError    bool
		expectedMessages bool
	}{
		{
			request: NameserverRequest{
				Host: utils.NewString("ns1.example.com."),
				IPV4: utils.NewString("127.0.0.1"),
				IPV6: utils.NewString("::1"),
			},
			expectedError:    false,
			expectedMessages: false,
		},
		{
			request: NameserverRequest{
				Host: nil,
				IPV4: utils.NewString("127.0.0.1"),
				IPV6: utils.NewString("::1"),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			request: NameserverRequest{
				Host: utils.NewString("ns1.example.com."),
				IPV4: utils.NewString("not valid IP"),
				IPV6: utils.NewString("::1"),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			request: NameserverRequest{
				Host: utils.NewString("ns1.example.com."),
				IPV4: utils.NewString("127.0.0.1"),
				IPV6: utils.NewString("not valid IP"),
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
