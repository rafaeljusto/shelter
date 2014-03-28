// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package protocol

import (
	"errors"
	"regexp"
	"time"
)

var (
	// When the JSON time doesn't have quotes or doesn't have a valid sequence of characters
	// this error will be throw
	ErrPreciseTimeFormat = errors.New("JSON time has an invalid format")
)

var (
	// Flexible date verification for RFC3339Nano with quotes
	isJSONTime = regexp.MustCompile(`"[0-9]+\-[0-9]+\-[0-9]+(?i:T)[0-9]+\:[0-9]+\:[0-9]+(\.[0-9]+)?(((\-)?[0-9]+\:[0-9]+)|(?i:Z))?"`)
)

// PreciseTime structure created to add milliseconds to the string representation of the
// dates on the system when exposing to the user via JSON protocol
type PreciseTime struct {
	time.Time
}

// Implementing encoding/json interface to transform a time into string representation
func (p *PreciseTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.Format(time.RFC3339Nano) + `"`), nil
}

// Implementing encoding/json interface to transform a string into a time object
func (p *PreciseTime) UnmarshalJSON(data []byte) error {
	// At least the quotes must exist
	if !isJSONTime.MatchString(string(data)) {
		return ErrPreciseTimeFormat
	}

	t, err := time.Parse(time.RFC3339Nano, string(data[1:len(data)-1]))
	if err != nil {
		return err
	}

	*p = PreciseTime{t}
	return nil
}
