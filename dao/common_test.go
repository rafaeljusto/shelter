// dao - Objects persistence layer
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package dao

import (
	"testing"
)

func TestDAOOrderByDirectionFromString(t *testing.T) {
	if _, err := DAOOrderByDirectionFromString("xxx"); err == nil {
		t.Error("Accepting an invalid order by direction")
	}

	if direction, err := DAOOrderByDirectionFromString("  ASC  "); err != nil || direction != DAOOrderByDirectionAscending {
		t.Error("Not accepting a valid order by direction ASC")
	}

	if direction, err := DAOOrderByDirectionFromString("desc"); err != nil || direction != DAOOrderByDirectionDescending {
		t.Error("Not accepting a valid order by direction DESC")
	}
}

func TestDAOOrderByDirectionToString(t *testing.T) {
	if direction := DAOOrderByDirectionToString(DAOOrderByDirection(9999)); len(direction) > 0 {
		t.Error("Not returning empty string when is an unknown order by direction")
	}

	if direction := DAOOrderByDirectionToString(DAOOrderByDirectionAscending); direction != "asc" {
		t.Error("Not returning the correct order by direction for ASC")
	}

	if direction := DAOOrderByDirectionToString(DAOOrderByDirectionDescending); direction != "desc" {
		t.Error("Not returning the correct order by direction for DESC")
	}
}
