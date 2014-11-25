// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// utils add features to make the test life easier
package utils

// NewString returns the pointer of the given string. Useful for unit tests of the protocol package,
// where we need to initialize string pointer inline
func NewString(value string) *string {
	return &value
}

// NewUint8 returns the pointer of the given int. Useful for unit tests of the protocol package,
// where we need to initialize int pointer inline
func NewUint8(value uint8) *uint8 {
	return &value
}

// NewUint16 returns the pointer of the given int. Useful for unit tests of the protocol package,
// where we need to initialize int pointer inline
func NewUint16(value uint16) *uint16 {
	return &value
}
