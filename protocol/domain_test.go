// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"testing"
)

func TestDomainNormalize(t *testing.T) {
	data := []struct {
		description string
		request     DomainRequest
		expected    DomainRequest
	}{}

	for i, item := range data {
		item.request.Normalize()
	}
}

func TestDomainValidate(t *testing.T) {
	data := []struct {
		description      string
		request          DomainRequest
		expectedError    bool
		expectedMessages bool
	}{}

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
