// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"testing"
)

func TestDomainOrderByFieldFromString(t *testing.T) {
	if _, ok := DomainOrderByFieldFromString("xxx"); ok {
		t.Error("Accepting an invalid order by field")
	}

	if field, ok := DomainOrderByFieldFromString("  FQDN  "); !ok || field != DomainOrderByFieldFQDN {
		t.Error("Not accepting a valid order by field FQDN")
	}

	if field, ok := DomainOrderByFieldFromString("lastModified"); !ok || field != DomainOrderByFieldLastModifiedAt {
		t.Error("Not accepting a valid order by field LastModified")
	}
}

func TestDomainOrderByFieldToString(t *testing.T) {
	if field := DomainOrderByFieldToString(DomainOrderByField(9999)); len(field) > 0 {
		t.Error("Not returning empty string when is an unknown order by field")
	}

	if field := DomainOrderByFieldToString(DomainOrderByFieldFQDN); field != "fqdn" {
		t.Error("Not returning the correct order by field for FQDN")
	}

	if field := DomainOrderByFieldToString(DomainOrderByFieldLastModifiedAt); field != "lastmodified" {
		t.Error("Not returning the correct order by field for LastModified")
	}
}
