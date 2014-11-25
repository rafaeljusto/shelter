// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"github.com/rafaeljusto/shelter/testing/utils"
	"testing"
)

func TestOwnerNormalize(t *testing.T) {
	data := []struct {
		description string
		request     OwnerRequest
		expected    OwnerRequest
	}{
		{
			description: "it should normalize the language",
			request: OwnerRequest{
				Language: utils.NewString("   PT   "),
			},
			expected: OwnerRequest{
				Language: utils.NewString("pt"),
			},
		},
		{
			description: "it should normalize the language with specific region",
			request: OwnerRequest{
				Language: utils.NewString("   PT   -   br   "),
			},
			expected: OwnerRequest{
				Language: utils.NewString("pt-BR"),
			},
		},
	}

	for i, item := range data {
		item.request.Normalize()

		if !utils.CompareStrings(item.request.Email, item.expected.Email) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Email,
				item.expected.Email,
			)
		}

		if !utils.CompareStrings(item.request.Language, item.expected.Language) {
			t.Errorf(
				"Item %d, “%s”: mismatch results. Expecting '%v'; found '%v'",
				i,
				item.description,
				item.request.Language,
				item.expected.Language,
			)
		}
	}
}

func TestOwnerValidate(t *testing.T) {
	data := []struct {
		description      string
		request          OwnerRequest
		expectedError    bool
		expectedMessages bool
	}{
		{
			description: "it should validate correctly",
			request: OwnerRequest{
				Email:    utils.NewString("test@example.com"),
				Language: utils.NewString("pt-BR"),
			},
			expectedError:    false,
			expectedMessages: false,
		},
		{
			description: "it should alert when there's an invalid e-mail",
			request: OwnerRequest{
				Email:    utils.NewString("testexamplecom"),
				Language: utils.NewString("pt-BR"),
			},
			expectedError:    false,
			expectedMessages: true,
		},
		{
			description: "it should alert when there's an invalid language",
			request: OwnerRequest{
				Email:    utils.NewString("test@example.com"),
				Language: utils.NewString("xxx"),
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
