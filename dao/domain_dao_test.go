// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package dao

import (
	"testing"
)

func TestDomainDAOOrderByFieldFromString(t *testing.T) {
	if _, err := DomainDAOOrderByFieldFromString("xxx"); err == nil {
		t.Error("Accepting an invalid order by field")
	}

	if field, err := DomainDAOOrderByFieldFromString("  FQDN  "); err != nil || field != DomainDAOOrderByFieldFQDN {
		t.Error("Not accepting a valid order by field FQDN")
	}

	if field, err := DomainDAOOrderByFieldFromString("lastModified"); err != nil || field != DomainDAOOrderByFieldLastModifiedAt {
		t.Error("Not accepting a valid order by field LastModified")
	}
}

func TestDomainDAOOrderByFieldToString(t *testing.T) {
	if field := DomainDAOOrderByFieldToString(DomainDAOOrderByField(9999)); len(field) > 0 {
		t.Error("Not returning empty string when is an unknown order by field")
	}

	if field := DomainDAOOrderByFieldToString(DomainDAOOrderByFieldFQDN); field != "fqdn" {
		t.Error("Not returning the correct order by field for FQDN")
	}

	if field := DomainDAOOrderByFieldToString(DomainDAOOrderByFieldLastModifiedAt); field != "lastmodified" {
		t.Error("Not returning the correct order by field for LastModified")
	}
}
