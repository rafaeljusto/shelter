package notification

import (
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/model"
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

// Load all templates from disk. The templates files should be in the templates path and
// must use as filename the language of the template (ex. pt-BR, pt, en-US, en). The
// language should follow the IANA format
func LoadTemplates() error {
	// Languages from configuration file were already checked when it was loaded into memory
	for _, language := range config.ShelterConfig.Languages {
		templatePath := filepath.Join(
			config.ShelterConfig.BasePath,
			config.ShelterConfig.Notification.TemplatesPath,
			fmt.Sprintf("%s.tmpl", language),
		)

		t, err := template.New("notification").Funcs(template.FuncMap{
			"nsStatusEq": func(nameserverStatus model.NameserverStatus,
				expectedNameserverTextStatus string) bool {

				return strings.ToLower(model.NameserverStatusToString(nameserverStatus)) ==
					strings.TrimSpace(strings.ToLower(expectedNameserverTextStatus))
			},
			"dsStatusEq": func(dsStatus model.DSStatus, expectedDSTextStatus string) bool {
				return strings.ToLower(model.DSStatusToString(dsStatus)) ==
					strings.TrimSpace(strings.ToLower(expectedDSTextStatus))
			},
			"isExpired": func(ds model.DS) bool {
				// If the status of the DS record is OK, it was selected because the expiration is
				// near. If in the future we add some other notification over the well configured
				// DS we must change this logic
				return ds.LastStatus == model.DSStatusOK
			},
		}).ParseFiles(templatePath)

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
	templatesLock.Lock()
	defer templatesLock.Unlock()
	templates[language] = t
}

// While notifing we will use a specific template to send an e-mail for the owner. To
// allow concurrent access over the templates map we should use this function
func getTemplate(language string) *template.Template {
	templatesLock.RLock()
	defer templatesLock.RUnlock()
	return templates[language]
}
