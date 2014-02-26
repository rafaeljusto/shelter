package notification

import (
	"errors"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/model"
	"io/ioutil"
	"path/filepath"
	"sync"
	"text/template"
)

// List of possible errors that can occur in this language controller. There can be also
// other errors from low level drivers (like unknown templates path)
var (
	// All template names are checked to see if they match with the predefined list of IANA
	// languages, if not this error is returned to alert the administrator
	ErrInvalidTemplateName = errors.New("Template name does not follow a valid language")
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
	templatesPath := filepath.Join(
		config.ShelterConfig.BasePath,
		config.ShelterConfig.Notification.TemplatesPath,
	)

	filesInfo, err := ioutil.ReadDir(templatesPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range filesInfo {
		if fileInfo.IsDir() {
			continue
		}

		if !model.IsValidLanguage(fileInfo.Name()) {
			return ErrInvalidTemplateName
		}

		templatePath := filepath.Join(
			config.ShelterConfig.BasePath,
			config.ShelterConfig.Notification.TemplatesPath,
			fileInfo.Name(),
		)

		t, err := template.New("notification").ParseFiles(templatePath)
		if err != nil {
			return err
		}

		addTemplate(fileInfo.Name(), t)
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
