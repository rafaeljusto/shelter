package messages

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
)

var (
	ShelterRESTLanguagePacks LanguagePacks // Store all possible languages
	ShelterRESTLanguagePack  *LanguagePack // Store the messages from the selected language
)

// List of possible errors that can occur in this language controller. There can be also
// other errors from low level drivers (like json unmarshall)
var (
	// When the default language is not found in the language pack this error is used. This
	// is a critical error because the ShelterRESTLanguagePack will be null and can cause
	// panic when the system try to use it
	ErrDefaultLanguageNotFound = errors.New("Default language not found in configuration file")
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

// Return the language name to be add to HTTP header response
func (l *LanguagePack) Name() string {
	if len(l.SpecificName) > 0 {
		return l.SpecificName
	}

	return l.GenericName
}

// Load the language packs from the configuration file
func LoadConfig(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &ShelterRESTLanguagePacks); err != nil {
		return err
	}

	// Load the default language pack
	ShelterRESTLanguagePack = ShelterRESTLanguagePacks.Select(ShelterRESTLanguagePacks.Default)
	if ShelterRESTLanguagePack == nil {
		return ErrDefaultLanguageNotFound
	}

	return nil
}
