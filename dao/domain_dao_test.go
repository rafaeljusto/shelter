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
		t.Error("Not returning the correct order by field for FQDN")
	}
}

func TestDomainDAOOrderByDirectionFromString(t *testing.T) {
	if _, err := DomainDAOOrderByDirectionFromString("xxx"); err == nil {
		t.Error("Accepting an invalid order by direction")
	}

	if direction, err := DomainDAOOrderByDirectionFromString("  ASC  "); err != nil || direction != DomainDAOOrderByDirectionAscending {
		t.Error("Not accepting a valid order by direction ASC")
	}

	if direction, err := DomainDAOOrderByDirectionFromString("desc"); err != nil || direction != DomainDAOOrderByDirectionDescending {
		t.Error("Not accepting a valid order by direction DESC")
	}
}

func TestDomainDAOOrderByDirectionToString(t *testing.T) {
	if direction := DomainDAOOrderByDirectionToString(DomainDAOOrderByDirection(9999)); len(direction) > 0 {
		t.Error("Not returning empty string when is an unknown order by direction")
	}

	if direction := DomainDAOOrderByDirectionToString(DomainDAOOrderByDirectionAscending); direction != "asc" {
		t.Error("Not returning the correct order by direction for ASC")
	}

	if direction := DomainDAOOrderByDirectionToString(DomainDAOOrderByDirectionDescending); direction != "desc" {
		t.Error("Not returning the correct order by direction for DESC")
	}
}
