// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package check verify REST policies
package check

import (
	"errors"
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"net/http"
	"strings"
	"testing"
	"time"
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

	context, err := context.NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	messages.ShelterRESTLanguagePacks = messages.LanguagePacks{
		Default: "en-US",
		Packs: []messages.LanguagePack{
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

	r.Header.Set("Content-Type", SupportedContentType+";v=2;charset=iso-8859-1")
	if HTTPContentType(r) {
		t.Error("Accepting a charset that the system is not prepared")
	}

	r.Header.Set("Content-Type", SupportedContentType+";v=2;charset=utf-8")
	if !HTTPContentType(r) {
		t.Error("Not accepting a valid charset")
	}
}

func TestHTTPContentMD5(t *testing.T) {
	r, err := http.NewRequest("", "", strings.NewReader("Check Integrity!"))
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	context, err := context.NewContext(r, nil)
	if err != nil {
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

func TestHTTPDate(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Date", "")
	if ok, err := HTTPDate(r); !ok || err != nil {
		t.Error("Not accepting when there's no " +
			"HTTP Date header field")
	}

	r.Header.Set("Date", "2006-01-02T15:04:05Z07:00")
	if ok, err := HTTPDate(r); ok || err == nil {
		t.Error("Accepting an invalid date format out " +
			"of the time frame. Should only support RFC 1123")
	}

	r.Header.Set("Date", "Mon, 02 Jan 2006 15:04:05 MST")
	if ok, _ := HTTPDate(r); ok {
		t.Error("Accepting a HTTP Date header " +
			"field outside the time frame")
	}

	timeInTimeFrame := time.
		Now().
		UTC().
		Add(-timeFrameDuration / 2).
		Format(time.RFC1123)

	r.Header.Set("Date", "   "+strings.ToUpper(timeInTimeFrame)+"  ")
	if ok, err := HTTPDate(r); !ok || err != nil {
		t.Error("Not accepting a valid HTTP Date header "+
			"field with spaces and case changes. Details:", err)
	}
}

func TestHTTPAuthorization(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("Authorization", "")
	if ok, err := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "abc123", nil
	}); ok || err == nil {
		t.Error("Accepting requests without HTTP Authorization header field")
	}

	r.Header.Set("Authorization", SupportedNamespace+" 1:bWq2s1WEIj+Ydj0vQ697zp+IXMU=")
	if ok, err := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "abc123", nil
	}); ok || err == nil {
		t.Error("Accepting an invalid HTTP Authorization header field")
	}

	timeInTimeFrame := time.
		Now().
		UTC().
		Add(-timeFrameDuration / 2).
		Format(time.RFC1123)
	r.Header.Set("Date", timeInTimeFrame)

	stringToSign, err := BuildStringToSign(r, "1")
	if err != nil {
		t.Fatal(err)
	}
	signature := GenerateSignature(stringToSign, "abc123")

	r.Header.Set("Authorization", SupportedNamespace+"X 1:"+signature)
	if ok, err := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "abc123", nil
	}); ok || err == nil {
		t.Error("Accepting a HTTP Authorization header field " +
			"with an invalid namespace")
	}

	r.Header.Set("Authorization", SupportedNamespace+"1:"+signature)
	if ok, err := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "abc123", nil
	}); ok || err == nil {
		t.Error("Accepting an invalid HTTP Authorization " +
			"header field with no spaces between namespace " +
			"and secret parts")
	}

	r.Header.Set("Authorization", SupportedNamespace+" 1"+signature)
	if ok, err := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "abc123", nil
	}); ok || err == nil {
		t.Error("Accepting an invalid HTTP Authorization " +
			"header field with no spaces between secret id " +
			"and secret")
	}

	r.Header.Set("Authorization", SupportedNamespace+" 2:"+signature)
	if ok, err := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "", errors.New("Not found")
	}); ok || err == nil {
		t.Error("Accepting an invalid HTTP Authorization " +
			"header field with an unknown secret id")
	}

	r.Header.Set("Authorization", SupportedNamespace+" 1:"+signature+"X")
	if ok, _ := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "abc123", nil
	}); ok {
		t.Error("Accepting an invalid HTTP Authorization " +
			"header field with a wrong signature")
	}

	r.Header.Set("Authorization", SupportedNamespace+" 1:"+signature)
	if ok, err := HTTPAuthorization(r, func(keyId string) (string, error) {
		return "abc123", nil
	}); !ok || err != nil {
		t.Error("Not accepting a valid HTTP Authorization "+
			"header field. Details:", err)
	}
}
