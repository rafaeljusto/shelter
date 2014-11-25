// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package messages manage the REST messages in a specific language
package messages

import (
	"encoding/json"
	"fmt"
	"github.com/rafaeljusto/shelter/errors"
	"github.com/rafaeljusto/shelter/language"
	"github.com/rafaeljusto/shelter/normalize"
	"io/ioutil"
	"strings"
	"sync"
)

var (
	ShelterRESTLanguagePacks LanguagePacks // Store all possible languages
	ShelterRESTLanguagePack  *LanguagePack // Store the messages from the selected language
)

var (
	// List of possible languages that can be used by the owner. This is an array of strings
	// instead of constant values because it depends on the e-mail template files created by
	// the system administrator.
	languages []string

	// As we are going to add and read the array of languages on-the-fly we need a lock
	// mechanism to allow concurrent access
	languagesLock sync.RWMutex
)

// List of possible errors that can occur in this language controller. There can be also
// other errors from low level drivers (like json unmarshall)
var (
	// When the default language is not found in the language pack this error is used. This
	// is a critical error because the ShelterRESTLanguagePack will be null and can cause
	// panic when the system try to use it
	ErrDefaultLanguageNotFound = fmt.Errorf("Default language not found in configuration file")

	// All loaded languages are checked to see if they match with the predefined list of
	// IANA languages, if not this error is returned to alert the administrator
	ErrInvalidLanguage = fmt.Errorf("Language is not valid")
)

// Structure responsable for storing all messages from the REST server in different idioms
// for a flexible internationalization
type LanguagePacks struct {
	Default string         // Language used for returning messages
	Packs   []LanguagePack // List of possible languages
}

// Select the language that is going to be used in the REST server messages based on the
// Language HTTP header
func (l *LanguagePacks) Select(language string) *LanguagePack {
	language = strings.ToLower(language)

	// Try searching for the specific language name first
	for _, pack := range l.Packs {
		if strings.ToLower(pack.SpecificName) == language {
			return &pack
		}
	}

	// Now we look for the first generic language name that we found
	for _, pack := range l.Packs {
		if strings.ToLower(pack.GenericName) == language {
			ShelterRESTLanguagePack = &pack
			return &pack
		}
	}

	return nil
}

// Return the languages supported by Shelter REST server
func (l *LanguagePacks) Names() string {
	var names []string
	for _, pack := range l.Packs {
		names = append(names, pack.Name())
	}
	return strings.Join(names, ",")
}

// LanguagePack defines a structure for a specific message
type LanguagePack struct {
	GenericName  string            // Language acronym (eg. pt)
	SpecificName string            // Language acronym with region (eg. pt-br)
	Messages     map[string]string // List of messages with identifiers
}

// Return the language name to be add to HTTP header response. We are going to normalize
// all the names in lower case
func (l *LanguagePack) Name() string {
	if len(l.SpecificName) > 0 {
		return normalize.NormalizeLanguage(l.SpecificName)
	}

	return normalize.NormalizeLanguage(l.GenericName)
}

// Load the language packs from the configuration file
func LoadConfig(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.NewSystemError(err)
	}

	if err := json.Unmarshal(bytes, &ShelterRESTLanguagePacks); err != nil {
		return errors.NewSystemError(err)
	}

	// Check language formats. They should follow IANA language types

	if !language.IsValidLanguage(ShelterRESTLanguagePacks.Default) {
		return errors.NewSystemError(ErrInvalidLanguage)
	}

	for _, l := range ShelterRESTLanguagePacks.Packs {
		if !language.IsValidLanguage(l.Name()) {
			return errors.NewSystemError(ErrInvalidLanguage)
		}
	}

	// Load the default language pack
	ShelterRESTLanguagePack = ShelterRESTLanguagePacks.Select(ShelterRESTLanguagePacks.Default)
	if ShelterRESTLanguagePack == nil {
		return errors.NewSystemError(ErrDefaultLanguageNotFound)
	}

	// TODO: Should we check if the configuration file languages are the same of the REST
	// messages file? Don't known what can we gain with this restriction, maybe system
	// language consistency.

	return nil
}

// AddLanguage is a safe way to add a supported language for the owner
func AddLanguage(l string) bool {
	if !language.IsValidLanguage(l) {
		// Error returned when a language input doesn't respect the language/region types
		// defined by IANA
		return false
	}

	languagesLock.Lock()
	defer languagesLock.Unlock()

	// Normalize all languages
	l = strings.ToLower(l)

	// Try to find duplicated value, we don't use the LanguageExists method because of the
	// language mutex
	for _, existingLanguage := range languages {
		if existingLanguage == l {
			// Ignore when we already have the language
			return true
		}
	}

	languages = append(languages, l)
	return true
}

// LanguageExists checks if a current language is supported
func LanguageExists(language string) bool {
	languagesLock.RLock()
	defer languagesLock.RUnlock()

	// Normalize all languages
	language = strings.ToLower(language)

	for _, existingLanguage := range languages {
		if existingLanguage == language {
			return true
		}
	}

	return false
}
