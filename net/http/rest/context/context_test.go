// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package context

import (
	"errors"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"math"
	"net/http"
	"strings"
	"testing"
)

type ReaderWithError struct{}

func (r *ReaderWithError) Read(p []byte) (n int, err error) {
	n = 0
	err = errors.New("Just throwing an error for tests")
	return
}

func TestNewContext(t *testing.T) {
	r, err := http.NewRequest("", "", strings.NewReader("Test"))
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(context.RequestContent) != "Test" {
		t.Error("Not storing request body correctly")
	}

	r, err = http.NewRequest("", "", &ReaderWithError{})
	if err != nil {
		t.Fatal(err)
	}

	r.ContentLength = 100
	_, err = NewContext(r, nil)
	if err == nil {
		t.Error("Not detecting request content error")
	}
}

func TestJSONRequest(t *testing.T) {
	r, err := http.NewRequest("", "",
		strings.NewReader("{\"key\": \"value\"}"))
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	object := struct {
		Key string `json:"key"`
	}{
		Key: "",
	}

	if err := context.JSONRequest(&object); err != nil {
		t.Fatal(err)
	}

	if object.Key != "value" {
		t.Error("Not decoding a JSON object properly")
	}
}

func TestResponse(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	context.Response(http.StatusNotFound)
	if context.ResponseHTTPStatus != http.StatusNotFound {
		t.Error("Not setting the return status code properly")
	}
}

func TestMessageReponse(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	context.Language = &messages.LanguagePack{
		Messages: map[string]string{
			"key": "value",
		},
	}

	if err := context.MessageResponse(http.StatusOK,
		"key", "this% is% not% an% URI"); err == nil {

		t.Fatal("Not detecting when the ROID is not a valid URI")
	}

	if err := context.MessageResponse(http.StatusOK,
		"key", "/domain/example.com.br."); err != nil {

		t.Fatal("Not accepting a valid message response")
	}

	if context.ResponseHTTPStatus != http.StatusOK {
		t.Error("Not setting the return status code properly")
	}

	if string(context.ResponseContent) !=
		`{"id":"key","message":"value","links":[{"types":["related"],"href":"/domain/example.com.br."}]}` {

		t.Error("Not setting the return message properly")
	}

	if err := context.MessageResponse(http.StatusNotFound, "key", ""); err != nil {
		t.Fatal("Not accepting a valid message response")
	}

	if context.ResponseHTTPStatus != http.StatusNotFound {
		t.Error("Not setting the return status code properly")
	}

	if string(context.ResponseContent) != `{"id":"key","message":"value"}` {
		t.Error("Not setting the return message properly")
	}
}

func TestJSONReponse(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	object := struct {
		Key string `json:"key"`
	}{
		Key: "value",
	}
	if err := context.JSONResponse(http.StatusNotFound, object); err != nil {
		t.Fatal("Not creating a valid JSON")
	}

	if context.ResponseHTTPStatus != http.StatusNotFound {
		t.Error("Not setting the return status code properly")
	}

	if string(context.ResponseContent) != "{\"key\":\"value\"}" {
		t.Error("Not setting the return message properly")
	}

	if err := context.JSONResponse(http.StatusOK, math.NaN()); err == nil {
		t.Error("Not detecting strange JSON objects")
	}
}

func TestAddHeader(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	context.AddHeader("content-language", "en-us")
	if _, ok := context.HTTPHeader["Content-Language"]; ok {
		t.Error("Allowing fixed HTTP headers to be replaced")
	}

	context.AddHeader("ETag", "1")
	if value, ok := context.HTTPHeader["ETag"]; !ok || value != "1" {
		t.Error("Not storing HTTP custom header properly")
	}
}
