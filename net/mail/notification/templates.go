// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package notification is the notification service
package notification

import (
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/model"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var (
	// Global variable used to store templates read from disk. We keep them in memory for a
	// faster notification
	templates map[string]*template.Template

	// Lock to avoid concurrent access on the templates map
	templatesLock sync.RWMutex
)

var (
	// Extension used in template files. Used for now as a variable for unit tests and integration
	// tests that can only predefine temporary filenames
	TemplateExtension string = ".tmpl"
)

func init() {
	templates = make(map[string]*template.Template)
}

// Load all templates from disk. The templates files should be in the templates path and
// must use as filename the language of the template (ex. pt-BR, pt, en-US, en). The
// language should follow the IANA format
func LoadTemplates() error {
	templatesPath := filepath.Join(
		config.ShelterConfig.BasePath,
		config.ShelterConfig.Notification.TemplatesPath,
	)

	// List all files of the directory to make it possible to load a case insensitive
	// language file
	filesInfo, err := ioutil.ReadDir(templatesPath)
	if err != nil {
		return err
	}

	// Languages from configuration file were already checked when it was loaded into memory
	for _, language := range config.ShelterConfig.Languages {
		filename := fmt.Sprintf("%s%s", language, TemplateExtension)
		templatePath := ""

		for _, fileInfo := range filesInfo {
			if fileInfo.IsDir() {
				continue
			}

			// We are listing all files in the directory to compare with the language file that
			// we want in a way that this could be case insensitive
			if strings.ToLower(fileInfo.Name()) == strings.ToLower(filename) {
				templatePath = filepath.Join(
					config.ShelterConfig.BasePath,
					config.ShelterConfig.Notification.TemplatesPath,
					fileInfo.Name(),
				)

				break
			}
		}

		templateContent, err := ioutil.ReadFile(templatePath)
		if err != nil {
			return err
		}

		t, err := template.New("notification").Funcs(template.FuncMap{
			"nsStatusEq":       nameserverStatusEquals,
			"dsStatusEq":       dsStatusEquals,
			"isNearExpiration": isNearExpirationDS,
		}).Parse(string(templateContent))

		if err != nil {
			return err
		}

		addTemplate(language, t)
	}

	return nil
}

// Safe way to add a template concurrently. In reallity we don't have concurrent problems
// while adding templates because there's only one synchronous function that add templates
// (LoadTemplates) and there's no read while we add them, but for consistency we are using
// it
func addTemplate(language string, t *template.Template) {
	language = model.NormalizeLanguage(language)

	templatesLock.Lock()
	defer templatesLock.Unlock()
	templates[language] = t
}

// While notifing we will use a specific template to send an e-mail for the owner. To
// allow concurrent access over the templates map we should use this function
func getTemplate(language string) *template.Template {
	language = model.NormalizeLanguage(language)

	templatesLock.RLock()
	defer templatesLock.RUnlock()
	return templates[language]
}

// Function created to clear the templates map, for now is used only for unit tests scenarios
func clearTemplates() {
	templatesLock.Lock()
	defer templatesLock.Unlock()
	templates = make(map[string]*template.Template)
}

// Auxiliary function for template that compares two nameserver status (case insensitive)
func nameserverStatusEquals(nameserverStatus model.NameserverStatus,
	expectedNameserverTextStatus string) bool {

	return strings.ToLower(model.NameserverStatusToString(nameserverStatus)) ==
		strings.TrimSpace(strings.ToLower(expectedNameserverTextStatus))
}

// Auxiliary function for template that compares two DS status (case insensitive)
func dsStatusEquals(dsStatus model.DSStatus, expectedDSTextStatus string) bool {
	return strings.ToLower(model.DSStatusToString(dsStatus)) ==
		strings.TrimSpace(strings.ToLower(expectedDSTextStatus))
}

// Auxiliary function for template that checks if a DS is near expiration or not
func isNearExpirationDS(ds model.DS) bool {
	// If the status of the DS record is OK, it was selected because the expiration is
	// near. If in the future we add some other notification over the well configured
	// DS we must change this logic
	return ds.LastStatus == model.DSStatusOK
}
