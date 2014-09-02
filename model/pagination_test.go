// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"testing"
)

func TestOrderByDirectionFromString(t *testing.T) {
	if _, ok := OrderByDirectionFromString("xxx"); ok {
		t.Error("Accepting an invalid order by direction")
	}

	if direction, ok := OrderByDirectionFromString("  ASC  "); !ok || direction != OrderByDirectionAscending {
		t.Error("Not accepting a valid order by direction ASC")
	}

	if direction, ok := OrderByDirectionFromString("desc"); !ok || direction != OrderByDirectionDescending {
		t.Error("Not accepting a valid order by direction DESC")
	}
}

func TestOrderByDirectionToString(t *testing.T) {
	if direction := OrderByDirectionToString(OrderByDirection(9999)); len(direction) > 0 {
		t.Error("Not returning empty string when is an unknown order by direction")
	}

	if direction := OrderByDirectionToString(OrderByDirectionAscending); direction != "asc" {
		t.Error("Not returning the correct order by direction for ASC")
	}

	if direction := OrderByDirectionToString(OrderByDirectionDescending); direction != "desc" {
		t.Error("Not returning the correct order by direction for DESC")
	}
}
