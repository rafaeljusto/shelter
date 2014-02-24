package model

import (
	"testing"
)

func TestLanguageTypeExists(t *testing.T) {
	if LanguageTypeExists("zz") {
		t.Error("Finding an unknown language")
	}

	if !LanguageTypeExists("  Pt ") {
		t.Error("Not finding an existing language with different case and spaces")
	}
}

func TestRegionTypeExists(t *testing.T) {
	if RegionTypeExists("zzz") {
		t.Error("Finding an unknown region")
	}

	if !RegionTypeExists("  bR ") {
		t.Error("Not finding an existing region with different case and spaces")
	}
}

func TestLanguageIsValid(t *testing.T) {
	if LanguageIsValid("zz") {
		t.Error("Validating language that does not exist")
	}

	if !LanguageIsValid(" pT ") {
		t.Error("Not validating language that is valid with cases and spaces")
	}

	if LanguageIsValid("pt-zzz") {
		t.Error("Validating language with region that does not exist")
	}

	if !LanguageIsValid(" pT-Br ") {
		t.Error("Not validating language with region that is valid with cases and spaces")
	}

	if LanguageIsValid("pt-br-zzz") {
		t.Error("Validating language with invalid format")
	}
}
