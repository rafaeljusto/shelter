// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package errors describe the low level errors of the Shelter system
package errors

import (
	"bytes"
	"errors"
	shelterLog "github.com/rafaeljusto/shelter/log"
	"log"
	"strings"
	"testing"
)

func TestNewSystemError(t *testing.T) {
	err := NewSystemError(errors.New("Error!"))

	if err.Error() != "Error!" {
		t.Error("Not storing low level error correctly")
	}
}

func TestNewSystemErrorLog(t *testing.T) {
	var buffer bytes.Buffer
	shelterLog.Logger = log.New(&buffer, "", log.Lshortfile)

	NewSystemError(errors.New("Something went wrong!"))

	got := strings.TrimSpace(buffer.String())
	expected := "Something went wrong!"
	if !strings.HasSuffix(got, expected) {
		t.Errorf("Not logging correctly when there's a system error. "+
			"Expected suffix '%s' and got '%s'", expected, got)
	}
}

func TestLogError(t *testing.T) {
	var buffer bytes.Buffer
	shelterLog.Logger = log.New(&buffer, "", log.Lshortfile)

	LogError(errors.New("Something went wrong!"))

	got := strings.TrimSpace(buffer.String())
	expected := "Something went wrong!"
	if !strings.HasSuffix(got, expected) {
		t.Errorf("Not logging correctly when there's a system error. "+
			"Expected suffix '%s' and got '%s'", expected, got)
	}
}
