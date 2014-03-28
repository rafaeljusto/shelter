// model - Description of the objects
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package model

import (
	"errors"
	"net/mail"
	"strings"
	"sync"
)

var (
	// Error returned when a language input doesn't respect the language/region types
	// defined by IANA
	ErrInvalidLanguage = errors.New("Language is not valid")
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

// Owner represents the responsable for the domain that can be alerted if any
// configuration problem is detected
type Owner struct {
	Email    *mail.Address // E-mail that will be alerted on any problem
	Language string        // Language used to send alerts
}

// AddLanguage is a safe way to add a supported language for the owner
func AddLanguage(language string) error {
	if !IsValidLanguage(language) {
		return ErrInvalidLanguage
	}

	languagesLock.Lock()
	defer languagesLock.Unlock()

	// Normalize all languages
	language = strings.ToLower(language)

	// Try to find duplicated value, we don't use the LanguageExists method because of the
	// language mutex
	for _, existingLanguage := range languages {
		if existingLanguage == language {
			// Ignore when we already have the language
			return nil
		}
	}

	languages = append(languages, language)
	return nil
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
