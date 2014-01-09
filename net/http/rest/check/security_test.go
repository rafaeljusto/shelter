package check

import (
	"encoding/base64"
	"net/http"
	"strings"
	"testing"
)

func TestBuildStringToSign(t *testing.T) {
	r, err := http.NewRequest(
		"GET",
		"http://test.com/domain?key2=value2&key1=value1",
		strings.NewReader("test"),
	)
	if err != nil {
		t.Fatal("Error creating the request. Details:", err)
	}

	_, err = BuildStringToSign(r, "1")
	if err == nil {
		t.Error("Building string to sign with a content without MD5 hash")
	}

	r.Header.Set("Content-MD5", "ZDhlOGZjYTJkYzBmODk2ZmQ3Y2I0Y2IwMDMxYmEyNDkgIC0K")

	_, err = BuildStringToSign(r, "1")
	if err == nil {
		t.Error("Building string to sign with a content without type defined")
	}

	r.Header.Set("Content-Type", "text/plain;v=5")

	_, err = BuildStringToSign(r, "1")
	if err == nil {
		t.Error("Building string to sign without date")
	}

	r.Header.Set("Date", "Mon, 02 Jan 2006 15:04:05 MST")

	stringToSign, err := BuildStringToSign(r, "1")
	if err != nil {
		t.Fatal(err)
	}

	if stringToSign != "GET\nZDhlOGZjYTJkYzBmODk2ZmQ3Y2I0Y2IwMDMxYmEyNDkgIC0K\ntext/plain\nMon, 02 Jan 2006 15:04:05 MST\n1\n/domain\nkey1=value1&key2=value2" {
		t.Error("Generating the wrong string to sign signature")
	}
}

func TestGenerateSignature(t *testing.T) {
	// Tests cases from RFC 2202. We are not testing the base64 enconding method, maybe we should
	// change this to compare generated base64 strings

	if GenerateSignature("Hi There", string([]byte{0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B, 0x0B})) != base64.StdEncoding.EncodeToString([]byte{0xb6, 0x17, 0x31, 0x86, 0x55, 0x05, 0x72, 0x64, 0xe2, 0x8b, 0xc0, 0xb6, 0xfb, 0x37, 0x8c, 0x8e, 0xf1, 0x46, 0xbe, 0x00}) {
		t.Error("Generating wrong signature")
	}

	if GenerateSignature("what do ya want for nothing?", "Jefe") != base64.StdEncoding.EncodeToString([]byte{0xef, 0xfc, 0xdf, 0x6a, 0xe5, 0xeb, 0x2f, 0xa2, 0xd2, 0x74, 0x16, 0xd5, 0xf1, 0x84, 0xdf, 0x9c, 0x25, 0x9a, 0x7c, 0x79}) {
		t.Error("Generating wrong signature")
	}
}
