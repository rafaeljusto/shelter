package check

import (
	"io/ioutil"
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

func TestHTTPAcceptCharset(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Accept-Charset", "")
	if !HTTPAcceptCharset(r) {
		t.Error("Not accepting when there's no " +
			"HTTP Accept Charset header field")
	}

	r.Header.Set("Accept-Charset", "iso-8859-1")
	if HTTPAcceptCharset(r) {
		t.Error("Accepting an unsupported language")
	}

	r.Header.Set("Accept-Charset", "iso-8859-1,    UTF-8  ")
	if !HTTPAcceptCharset(r) {
		t.Error("Not accepting a supported charset with " +
			"different case and spaces")
	}

	r.Header.Set("Accept-Charset", "iso-8859-1,utf-8;q=0.7")
	if !HTTPAcceptCharset(r) {
		t.Error("Not accepting a supported charset with options")
	}

	r.Header.Set("Accept-Charset", "*")
	if !HTTPAcceptCharset(r) {
		t.Error("Not accepting a wildcard charset")
	}
}

func TestHTTPContentType(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Content-Type", "")
	if !HTTPContentType(r) {
		t.Error("Not accepting when there's no " +
			"HTTP Content Type header field")
	}

	r.Header.Set("Content-Type", "text/html")
	if HTTPContentType(r) {
		t.Error("Accepting an unsupported content type")
	}

	r.Header.Set("Content-Type", "  "+strings.ToUpper(SupportedContentType)+"   ")
	if !HTTPContentType(r) {
		t.Error("Not accepting a supported content type with " +
			"different case and spaces")
	}

	r.Header.Set("Content-Type", SupportedContentType+";v=2")
	if !HTTPContentType(r) {
		t.Error("Not accepting a supported content type with options")
	}
}

func TestHTTPContentMD5(t *testing.T) {
	r, err := http.NewRequest("", "", strings.NewReader("Check Integrity!"))
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	var context context.ShelterRESTContext
	if context.RequestContent, err = ioutil.ReadAll(r.Body); err != nil {
		t.Fatal(err)
	}

	r.Header.Set("Content-MD5", "")
	if !HTTPContentMD5(r, &context) {
		t.Error("Not accepting when there's no " +
			"HTTP Content MD5 header field")
	}

	r.Header.Set("Content-MD5", "bGlmZSBvZiBicmlhbg==")
	if HTTPContentMD5(r, &context) {
		t.Error("Accepting an invalid content MD5")
	}

	r.Header.Set("Content-MD5", "   nwqq6b6ua/tTDk7B5M184w==  ")
	if !HTTPContentMD5(r, &context) {
		t.Error("Not accepting a valid content MD5 with spaces")
	}
}
