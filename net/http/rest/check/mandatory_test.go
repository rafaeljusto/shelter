package check

import (
	"net/http"
	"shelter/net/http/rest/context"
	"shelter/net/http/rest/language"
	"strings"
	"testing"
)

func TestHTTPAccept(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Accept", "")
	if !HTTPAccept(r) {
		t.Error("Not accepting when there's no HTTP Accept header field")
	}

	r.Header.Set("Accept", "text/html")
	if HTTPAccept(r) {
		t.Error("Accepting a not supported format")
	}

	r.Header.Set("Accept", "text/html,  "+
		strings.ToUpper(SupportedContentType)+"   ")
	if !HTTPAccept(r) {
		t.Error("Not accepting a supported format with " +
			"different case and spaces")
	}

	r.Header.Set("Accept", SupportedContentType+";version=2")
	if !HTTPAccept(r) {
		t.Error("Not accepting a supported format with " +
			"options")
	}

	r.Header.Set("Accept", "*")
	if !HTTPAccept(r) {
		t.Error("Not accepting wildcard format")
	}
}

func TestHTTPAcceptLanguage(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	var context context.ShelterRESTContext

	language.ShelterRESTLanguagePacks = language.LanguagePacks{
		Default: "en-US",
		Packs: []language.LanguagePack{
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

	r.Header.Set("Accept-Language", "")
	if !HTTPAcceptLanguage(r, &context) {
		t.Error("Not accepting when there's no " +
			"HTTP Accept Language header field")
	}

	r.Header.Set("Accept-Language", "de")
	if HTTPAcceptLanguage(r, &context) {
		t.Error("Accepting an unsupported language")
	}

	r.Header.Set("Accept-Language", "de,    PT-BR  ")
	if !HTTPAcceptLanguage(r, &context) {
		t.Error("Not accepting a supported language with " +
			"different case and spaces")
	}

	r.Header.Set("Accept-Language", "de,pt;q=0.4")
	if !HTTPAcceptLanguage(r, &context) {
		t.Error("Not accepting a supported generic language with options")
	}

	r.Header.Set("Accept-Language", "*")
	if !HTTPAcceptLanguage(r, &context) {
		t.Error("Not accepting a wildcard language")
	}
}
