// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package model describes the objects of the system
package model

import (
	"testing"
)

func TestAddLanguage(t *testing.T) {
	if AddLanguage("zzz") {
		t.Error("Not detecting an invalid language")
	}

	if !AddLanguage("pt-br") {
		t.Error("Not allowing a valid language")
	}

	if !AddLanguage("pt-br") {
		t.Error("Repeating same add language is not being ignored")
	}
}

func TestLanguageExists(t *testing.T) {
	if !AddLanguage("pt-br") {
		t.Fatal("Not allowing a valid language")
	}

	if LanguageExists("zzz") {
		t.Error("Found a language that does not exists")
	}

	if !LanguageExists("pt-br") {
		t.Error("Did not find a language that exists")
	}
}
