// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package messages manage the REST messages in a specific language
package messages

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSelect(t *testing.T) {
	languagePacks := LanguagePacks{
		Default: "en-US",
		Packs: []LanguagePack{
			{
				GenericName:  "en",
				SpecificName: "en-US",
			},
			{
				GenericName:  "pt",
				SpecificName: "pt-BR",
			},
		},
	}

	if languagePacks.Select("EN-US") == nil {
		t.Error("Not finding a valid specific language")
	}

	if languagePacks.Select("PT") == nil {
		t.Error("Not finding a valid generic language")
	}

	if languagePacks.Select("xx") != nil {
		t.Error("Finding an unknown language")
	}
}

func TestNames(t *testing.T) {
	languagePacks := LanguagePacks{
		Default: "en-US",
		Packs: []LanguagePack{
			{
				GenericName:  "en",
				SpecificName: "en-US",
			},
			{
				GenericName: "pt",
			},
		},
	}

	if languagePacks.Names() != "en-US,pt" {
		t.Error("Not building language names properly")
	}
}

func TestName(t *testing.T) {
	languagePack := LanguagePack{
		GenericName:  "en",
		SpecificName: "en-US",
	}

	if languagePack.Name() != "en-US" {
		t.Error("Not building specific language name properly")
	}

	languagePack = LanguagePack{
		GenericName: "en",
	}

	if languagePack.Name() != "en" {
		t.Error("Not building generic language name properly")
	}
}

func TestLoadConfig(t *testing.T) {
	if err := LoadConfig("UnknownFile.tmp"); err == nil {
		t.Error("Loading a file that does not exists!")
	}

	file, err := ioutil.TempFile(".", "shelter-lg-test")
	if err != nil {
		t.Fatal("Error creating test file")
	}

	fmt.Fprint(file, `{
    "default": "en-us",
    "packs": [
      {
        "GenericName": "en",
        "SpecificName": "en-us",
        "Messages": {
          "test-message": "This is a test!"
        }
      },
      {
        "GenericName": "pt",
        "SpecificName": "pt-br",
        "Messages": {
          "test-message": "This is another test!"
        }
      }
    ]
    }`)

	if err := LoadConfig(file.Name()); err != nil {
		t.Error("Not loading a good language file. Details:", err)
	}

	file.Close()
	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	file, err = ioutil.TempFile(".", "shelter-lg-test")
	if err != nil {
		t.Fatal("Error creating test file")
	}

	fmt.Fprint(file, `{
    "default": "en-us",
    "packs": [
      {
        "GenericName": "en",
        "SpecificName": "en-us",
        "Messages": {
          { "test-message": "This is a test!" }
        }
      }
    ]
    }`)

	if err := LoadConfig(file.Name()); err == nil {
		t.Error("Not detecting when language file doesn't have a valid format!")
	}

	file.Close()
	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	file, err = ioutil.TempFile(".", "shelter-lg-test")
	if err != nil {
		t.Fatal("Error creating test file")
	}

	fmt.Fprint(file, `{
    "default": "en-us",
    "packs": [
      {
        "GenericName": "en",
        "SpecificName": "en-zzzz",
        "Messages": {
          "test-message": "This is a test!"
        }
      }
    ]
    }`)

	if err := LoadConfig(file.Name()); err == nil {
		t.Error("Not detecting when language is invalid!")
	}

	file.Close()
	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	file, err = ioutil.TempFile(".", "shelter-lg-test")
	if err != nil {
		t.Fatal("Error creating test file")
	}

	fmt.Fprint(file, `{
    "default": "XXXX",
    "packs": [
      {
        "GenericName": "en",
        "SpecificName": "en-us",
        "Messages": {
          "test-message": "This is a test!"
        }
      },
      {
        "GenericName": "pt",
        "SpecificName": "pt-br",
        "Messages": {
          "test-message": "This is another test!"
        }
      }
    ]
    }`)

	if err := LoadConfig(file.Name()); err == nil {
		t.Error("Not detecting when default language is invalid!")
	}

	file.Close()
	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}

	file, err = ioutil.TempFile(".", "shelter-lg-test")
	if err != nil {
		t.Fatal("Error creating test file")
	}

	fmt.Fprint(file, `{
    "default": "en-br",
    "packs": [
      {
        "GenericName": "en",
        "SpecificName": "en-us",
        "Messages": {
          "test-message": "This is a test!"
        }
      },
      {
        "GenericName": "pt",
        "SpecificName": "pt-br",
        "Messages": {
          "test-message": "This is another test!"
        }
      }
    ]
    }`)

	if err := LoadConfig(file.Name()); err == nil {
		t.Error("Not detecting when default language is unknown!")
	}

	file.Close()
	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}
}

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
