// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package protocol describes the REST protocol
package protocol

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestPreciseTimeMarshalJSON(t *testing.T) {
	data := "2014-02-11T10:15:07.8766-02:00"
	exampleTime, err := time.Parse(time.RFC3339Nano, data)
	if err != nil {
		t.Fatal(err)
	}

	preciseTime := PreciseTime{Time: exampleTime}
	dataBytes, err := preciseTime.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if string(dataBytes) != fmt.Sprintf(`"%s"`, data) {
		t.Error("Marshal not working well")
	}

	dataBytes, err = json.Marshal(preciseTime)
	if err != nil {
		t.Fatal(err)
	}

	if string(dataBytes) != fmt.Sprintf(`"%s"`, data) {
		t.Error(fmt.Sprintf("Real marshal not working well. Expected: [%s], Received: [%s]",
			fmt.Sprintf(`"%s"`, data), string(dataBytes)))
	}
}

func TestPreciseTimeUnmarshalJSON(t *testing.T) {
	data := "2014-02-11T10:15:07.8766-02:00"
	dataBytes := []byte(fmt.Sprintf(`"%s"`, data))

	var preciseTime PreciseTime
	if err := preciseTime.UnmarshalJSON(dataBytes); err != nil {
		t.Fatal(err)
	}

	if preciseTime.Format(time.RFC3339Nano) != data {
		t.Error("Unmarshal not working well")
	}

	if err := json.Unmarshal(dataBytes, &preciseTime); err != nil {
		t.Error(err)
	}

	if preciseTime.Format(time.RFC3339Nano) != data {
		t.Error(fmt.Sprintf("Real unmarshal not working well. Expected: [%s], Received: [%s]",
			data, preciseTime.Format(time.RFC3339Nano)))
	}

	data = "2014-02-11T10:15:07.8766-02:00Z"
	dataBytes = []byte(fmt.Sprintf(`"%s"`, data))

	if err := json.Unmarshal(dataBytes, &preciseTime); err == nil {
		t.Error("Accepting an invalid JSON time value with invalid character")
	}

	data = "2000000000000014-02-11T10:15:07.8766-02:00"
	dataBytes = []byte(fmt.Sprintf(`"%s"`, data))

	if err := json.Unmarshal(dataBytes, &preciseTime); err == nil {
		t.Error("Accepting an invalid JSON time value to much in the future")
	}
}
