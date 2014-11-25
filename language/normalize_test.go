// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package language keep track of all IANA languages
package language

import (
	"testing"
)

func TestNormalizeLanguage(t *testing.T) {
	if NormalizeLanguage("") != "" {
		t.Error("Not normalizing correctly the empty language")
	}

	if NormalizeLanguage("  Pt  ") != "pt" {
		t.Error("Not normalizing correctly the language name")
	}

	if NormalizeLanguage("  Pt  -  bR ") != "pt-BR" {
		t.Error("Not normalizing correctly the language name with country code")
	}

	if NormalizeLanguage("  Pt  -  bR - zzzz") != "pt-BR" {
		t.Error("Not ignoring extra fields")
	}
}
