package context

import (
	"net/http"
	"shelter/net/http/rest/language"
	"strings"
	"testing"
)

func TestNewShelterRESTContext(t *testing.T) {
	r, err := http.NewRequest("", "", strings.NewReader("Test"))
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewShelterRESTContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(context.RequestContent) != "Test" {
		t.Error("Not storing request body correctly")
	}
}

func TestJSONRequest(t *testing.T) {
	r, err := http.NewRequest("", "",
		strings.NewReader("{\"key\": \"value\"}"))
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewShelterRESTContext(r, nil)
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

	context, err := NewShelterRESTContext(r, nil)
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

	context, err := NewShelterRESTContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	context.Language = &language.LanguagePack{
		Messages: map[string]string{
			"key": "value",
		},
	}

	context.MessageResponse(http.StatusNotFound, "key")

	if context.ResponseHTTPStatus != http.StatusNotFound {
		t.Error("Not setting the return status code properly")
	}

	if string(context.ResponseContent) != "value" {
		t.Error("Not setting the return message properly")
	}
}

func TestJSONReponse(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewShelterRESTContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	object := struct {
		Key string `json:"key"`
	}{
		Key: "value",
	}
	context.JSONResponse(http.StatusNotFound, object)

	if context.ResponseHTTPStatus != http.StatusNotFound {
		t.Error("Not setting the return status code properly")
	}

	if string(context.ResponseContent) != "{\"key\":\"value\"}" {
		t.Error("Not setting the return message properly")
	}
}

func TestAddHeader(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	context, err := NewShelterRESTContext(r, nil)
	if err != nil {
		t.Fatal(err)
	}

	context.AddHeader("Content-Type", "text/plain")
	if _, ok := context.HTTPHeader["Content-Type"]; ok {
		t.Error("Allowing fixed HTTP headers to be replaced")
	}

	context.AddHeader("ETag", "1")
	if value, ok := context.HTTPHeader["ETag"]; !ok || value != "1" {
		t.Error("Not storing HTTP custom header properly")
	}
}
