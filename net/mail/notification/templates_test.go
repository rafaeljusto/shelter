package notification

import (
	"github.com/rafaeljusto/shelter/config"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	config.ShelterConfig.BasePath = "."

	config.ShelterConfig.Notification.TemplatesPath = "idontexist"
	if err := LoadTemplates(); err == nil {
		t.Error("Not returnig error when the template path does not exist")
	}

	err := ioutil.WriteFile(
		"pt-BR.tmpl",
		[]byte("This is a test file"),
		os.ModePerm,
	)

	if err != nil {
		t.Fatal("Error creating test file")
	}

	config.ShelterConfig.Notification.TemplatesPath = "."

	config.ShelterConfig.Languages = []string{"pt-BR"}
	if err := LoadTemplates(); err != nil {
		t.Error(err)
	}

	config.ShelterConfig.Languages = []string{"en-US"}
	if err := LoadTemplates(); err == nil {
		t.Error("Not returnig error when a defined language doesn't have your " +
			"corresponding template file")
	}

	if err := os.Remove("pt-BR.tmpl"); err != nil {
		t.Fatal(err)
	}

	if err := os.Mkdir("pt-BR", os.ModePerm); err != nil {
		t.Fatal(err)
	}

	config.ShelterConfig.Languages = []string{"pt-BR"}
	if err := LoadTemplates(); err == nil && getTemplate("pt-BR") != nil {
		t.Error("Loading directory as template file")
	}

	if err := os.RemoveAll("pt-BR"); err != nil {
		t.Fatal(err)
	}
}

func TestAddTemplate(t *testing.T) {
	// TODO
}

func TestGetTemplate(t *testing.T) {
	// TODO
}
