// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package errors describe the low level errors of the Shelter system
package errors

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	NotFound = SystemError{
		Err: errors.New("Object not found"),
	}
)

// SystemError was created to encapsulate all low level errors, so we could store the file and line
// where did the error occurred. That really helps when we are trying to understand what's happening
type SystemError struct {
	Err  error
	File string
	Line int
}

// NewSystemError initialize an error storing the file and line where the problem occurred
func NewSystemError(err error) error {
	if _, ok := err.(SystemError); ok {
		// Trying to create a system error from another system error!
		return err
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	return SystemError{
		Err:  err,
		File: file,
		Line: line,
	}
}

// Error retrieve the description of the low level error
func (e SystemError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Err.Error())
}
