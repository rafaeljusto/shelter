package check

import (
	"net/http"
	"testing"
	"time"
)

func TestHTTPIfModifiedSince(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("If-Modified-Since", "")
	if ok, err := HTTPIfModifiedSince(r, time.Now()); !ok || err != nil {
		t.Error("Not accepting when there's no HTTP If-Modified-Since header field")
	}

	r.Header.Set("If-Modified-Since", "2006-01-02T15:04:05Z07:00")
	if ok, err := HTTPIfModifiedSince(r, time.Now()); !ok || err == nil {
		t.Error("Accepting an invalid date format. Should only support RFC 1123")
	}

	r.Header.Set("If-Modified-Since", time.Now().Add(10*time.Second).Format(time.RFC1123))
	if ok, _ := HTTPIfModifiedSince(r, time.Now()); ok {
		t.Error("Not comparing properly when object was modified")
	}

	r.Header.Set("If-Modified-Since", time.Now().Format(time.RFC1123))
	if ok, _ := HTTPIfModifiedSince(r, time.Now().Add(1*time.Second)); !ok {
		t.Error("Not comparing properly when object was modified")
	}
}

func TestHTTPIfUnmodifiedSince(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	r.Header.Set("If-Unmodified-Since", "")
	if ok, err := HTTPIfUnmodifiedSince(r, time.Now()); !ok || err != nil {
		t.Error("Not accepting when there's no HTTP If-Unmodified-Since header field")
	}

	r.Header.Set("If-Unmodified-Since", "2006-01-02T15:04:05Z07:00")
	if ok, err := HTTPIfUnmodifiedSince(r, time.Now()); !ok || err == nil {
		t.Error("Accepting an invalid date format. Should only support RFC 1123")
	}

	r.Header.Set("If-Unmodified-Since", time.Now().Format(time.RFC1123))
	if ok, _ := HTTPIfUnmodifiedSince(r, time.Now().Add(1*time.Second)); ok {
		t.Error("Not comparing properly when object was modified")
	}

	r.Header.Set("If-Unmodified-Since", time.Now().Add(1*time.Second).Format(time.RFC1123))
	if ok, _ := HTTPIfUnmodifiedSince(r, time.Now()); !ok {
		t.Error("Not comparing properly when object was modified")
	}
}
