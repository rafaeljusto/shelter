// model - Description of the objects
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

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

func TestIsValidLanguage(t *testing.T) {
	if IsValidLanguage("zz") {
		t.Error("Validating language that does not exist")
	}

	if !IsValidLanguage(" pT ") {
		t.Error("Not validating language that is valid with cases and spaces")
	}

	if IsValidLanguage("pt-zzz") {
		t.Error("Validating language with region that does not exist")
	}

	if !IsValidLanguage(" pT-Br ") {
		t.Error("Not validating language with region that is valid with cases and spaces")
	}

	if IsValidLanguage("pt-br-zzz") {
		t.Error("Validating language with invalid format")
	}
}
