// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package notification is the notification service
package notification

import (
	"io/ioutil"
	"os"
	"testing"
	"text/template"
	"time"

	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/model"
)

func TestLoadTemplates(t *testing.T) {
	clearTemplates()
	config.ShelterConfig.BasePath = "."

	// Change extension to empty because we can only set the prefix of temporary files
	TemplateExtension = ""

	config.ShelterConfig.Notification.TemplatesPath = "idontexist"
	if err := LoadTemplates(); err == nil {
		t.Error("Not returnig error when the template path does not exist")
	}

	file, err := ioutil.TempFile(".", "shelter-nf-test")
	if err != nil {
		t.Fatal("Error creating test file")
	}

	config.ShelterConfig.Notification.TemplatesPath = "."
	config.ShelterConfig.Languages = []string{file.Name()}
	if err := LoadTemplates(); err != nil {
		t.Error(err)
	}

	config.ShelterConfig.Languages = []string{file.Name() + "idontexist"}
	if err := LoadTemplates(); err == nil {
		t.Error("Not returnig error when a defined language doesn't have your " +
			"corresponding template file")
	}

	file.Close()
	if err := os.Remove(file.Name()); err != nil {
		t.Fatal(err)
	}
}

func TestLoadTemplatesWithDirectory(t *testing.T) {
	clearTemplates()

	dirName, err := ioutil.TempDir(".", "shelter-nf-test")
	if err != nil {
		t.Fatal(err)
	}

	config.ShelterConfig.BasePath = "."
	config.ShelterConfig.Languages = []string{dirName}

	LoadTemplates()

	if getTemplate(dirName) != nil {
		t.Error("Loading directory as template file")
	}

	if err := os.RemoveAll(dirName); err != nil {
		t.Fatal(err)
	}
}

func TestAddAndGetTemplate(t *testing.T) {
	clearTemplates()

	addTemplate("pt-br", new(template.Template))
	addTemplate("en-US", new(template.Template))

	if len(templates) != 2 {
		t.Error("Not adding templates correctly")
	}

	if getTemplate("pt-BR") == nil {
		t.Error("Not retrieving the template properly")
	}

	if getTemplate("zzzz") != nil {
		t.Error("Retrieving an unknown template")
	}
}

func TestClearTemplates(t *testing.T) {
	clearTemplates()

	addTemplate("pt-br", new(template.Template))
	addTemplate("en-US", new(template.Template))

	clearTemplates()

	if len(templates) != 0 {
		t.Error("Not cleaning templates correctly")
	}
}

func TestNameserverStatusEquals(t *testing.T) {
	if !nameserverStatusEquals(model.NameserverStatusNotChecked, "NOTcHEcKED   ") {
		t.Error("Not comparing correctly when namserver status are equal")
	}

	if nameserverStatusEquals(model.NameserverStatusCanonicalName, "ZZZ") {
		t.Error("Not returnig false when status are different")
	}
}

func TestDSStatusEquals(t *testing.T) {
	if !dsStatusEquals(model.DSStatusNoKey, "noKey   ") {
		t.Error("Not comparing correctly when DS status are equal")
	}

	if dsStatusEquals(model.DSStatusSignatureError, "ZZZ") {
		t.Error("Not returnig false when status are different")
	}
}

func TestIsNearExpirationDS(t *testing.T) {
	config.ShelterConfig.Scan.VerificationIntervals.MaxExpirationAlertDays = 2

	if !isNearExpirationDS(model.DS{
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}) {
		t.Error("Not detecting when DS is near expiration")
	}

	config.ShelterConfig.Scan.VerificationIntervals.MaxExpirationAlertDays = 1

	if isNearExpirationDS(model.DS{
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}) {
		t.Error("Returning near expiration is wrong scenarios")
	}
}

func TestFQDNToUnicode(t *testing.T) {
	data := []struct {
		fqdn     string
		expected string
	}{
		{fqdn: "example.com", expected: "example.com"},
		{fqdn: "xn--mgbaal8b0b9b2b.icom.museum", expected: "افغانستا.icom.museum"},
		{fqdn: "xn--romnia-yta.icom.museum", expected: "românia.icom.museum"},
		{fqdn: "xn--!.com", expected: "xn--!.com"},
	}

	for i, item := range data {
		fqdn := fqdnToUnicode(item.fqdn)
		if fqdn != item.expected {
			t.Errorf("item %d: expected FQDN “%s” and got “%s”", i, item.expected, fqdn)
		}
	}
}

func TestNormalizeEmailHeader(t *testing.T) {
	data := []struct {
		value    string
		expected string
	}{
		{
			value:    "u understand the example.",
			expected: "=?UTF-8?B?dSB1bmRlcnN0YW5kIHRoZSBleGFtcGxlLg==?=",
		},
		{
			value:    "Problema de configuración con el dominio sebastián.cl",
			expected: "=?UTF-8?B?UHJvYmxlbWEgZGUgY29uZmlndXJhY2nDs24gY29uIGVsIGRvbWluaW8gc2ViYXN0acOhbi5jbA==?=",
		},
	}

	for i, item := range data {
		value := normalizeEmailHeader(item.value)
		if value != item.expected {
			t.Errorf("item %d: expected value “%s” and got “%s”", i, item.expected, value)
		}
	}
}
