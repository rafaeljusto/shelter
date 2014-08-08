// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package errors describe the low level errors of the Shelter system
package errors

import (
	"errors"
	"fmt"
	"github.com/rafaeljusto/shelter/log"
	"runtime"
)

var (
	NotFound = SystemError{
		Err: errors.New("Object not found"),
	}
)

// SystemError was created to encapsulate all low level errors, so we could log the error exactly in
// the file and line that it occur. That really helps when we are trying to understand what's
// happening
type SystemError struct {
	Err error
}

// NewSystemError initialize a system error logging the problem
func NewSystemError(err error) SystemError {
	_, file, line, ok := runtime.Caller(1)

	if !ok {
		file = "???"
		line = 0
	}

	log.Printf("[ERR] %s:%d: %s", file, line, err.Error())

	return SystemError{
		Err: err,
	}
}

func LogError(err error) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	log.Printf(fmt.Sprintf("[ERR] %s:%d: %s", file, line, err.Error()))
}

// Error retrieve the description of the low level error
func (e SystemError) Error() string {
	return e.Err.Error()
}
