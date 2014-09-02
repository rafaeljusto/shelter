// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package errors describe the low level errors of the Shelter system
package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestNewSystemError(t *testing.T) {
	err := NewSystemError(errors.New("Something went wrong!"))

	sysErr, ok := err.(SystemError)
	if !ok {
		t.Fatal("Not creating a system error")
	}

	if !strings.HasSuffix(sysErr.File, "system_error_test.go") {
		t.Error("Not storing the file where the error occurred")
	}

	if !strings.HasSuffix(sysErr.Error(), "Something went wrong!") {
		t.Error("Not storing low level error correctly")
	}
}
