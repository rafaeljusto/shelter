package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"shelter/config"
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/net/http/rest"
	"shelter/net/http/rest/check"
	"shelter/net/http/rest/messages"
	"shelter/testing/utils"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "RESTMux"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTMux test")
}

func main() {
	flag.Parse()

	err := utils.ReadConfigFile(configFilePath, &config.ShelterConfig)

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	database, err := mongodb.Open(
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// If there was some problem in the last test, there could be some data in the
	// database, so let's clear it to don't affect this test. We avoid checking the error,
	// because if the collection does not exist yet, it will be created in the first
	// insert
	domainDAO.RemoveAll()

	var mux rest.Mux
	checkHeaders(&mux)
	createDomain(&mux)
	checkWrongACL(mux)
	deleteDomain(&mux)

	utils.Println("SUCCESS!")
}

func checkHeaders(mux *rest.Mux) {
	r, err := http.NewRequest("GET", "/xxx/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusServiceUnavailable {
		utils.Fatalln("Not returning HTTP Service Unavailable when the URI is not registered", nil)
	}

	data := []struct {
		Header             string
		HeaderValue        string
		ExpectedHTTPStatus int
	}{
		{
			Header:             "Accept-Language",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusNotAcceptable,
		},
		{
			Header:             "Accept",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusNotAcceptable,
		},
		{
			Header:             "Accept-Charset",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusNotAcceptable,
		},
		{
			Header:             "Content-Type",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Content-Type",
			HeaderValue:        "",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Content-MD5",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Content-MD5",
			HeaderValue:        "",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Date",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Date",
			HeaderValue:        "",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Authorization",
			HeaderValue:        "xxx",
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Authorization",
			HeaderValue:        fmt.Sprintf("%s %d:%s", check.SupportedNamespace, 999, "0PN5J17HBGZHT7JJ3X82"),
			ExpectedHTTPStatus: http.StatusBadRequest,
		},
		{
			Header:             "Authorization",
			HeaderValue:        fmt.Sprintf("%s %d:%s", check.SupportedNamespace, 1, "0PN5J17HBGZHT7JJ3X82"),
			ExpectedHTTPStatus: http.StatusUnauthorized,
		},
	}

	content := []byte(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "admin@example.com.br."
      ]
    }`)

	for _, item := range data {
		r, err = http.NewRequest("PUT", "/domain/example.com.br.",
			bytes.NewReader(content))
		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}

		buildHTTPHeader(r, content)
		r.Header.Set(item.Header, item.HeaderValue)

		w = httptest.NewRecorder()
		mux.ServeHTTP(w, r)

		if w.Code != item.ExpectedHTTPStatus {
			utils.Fatalln(fmt.Sprintf("For HTTP header %s with the value '%s' we were "+
				"expecting the HTTP status %d but got %d",
				item.Header, item.HeaderValue, item.ExpectedHTTPStatus, w.Code), nil)
		}
	}
}

func createDomain(mux *rest.Mux) {
	content := []byte(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "admin@example.com.br."
      ]
    }`)

	r, err := http.NewRequest("PUT", "/domain/example.com.br.",
		bytes.NewReader(content))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	buildHTTPHeader(r, content)

	messages.ShelterRESTLanguagePacks = messages.LanguagePacks{
		Default: "en-us",
		Packs: []messages.LanguagePack{
			{
				GenericName:  "en",
				SpecificName: "en-us",
			},
			{
				GenericName:  "pt",
				SpecificName: "pt-br",
			},
		},
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		utils.Fatalln(fmt.Sprintf("Error creting a domain object. Expected HTTP status %d and got %d",
			http.StatusCreated, w.Code), nil)
	}

	if len(w.Header().Get("Accept")) == 0 {
		utils.Fatalln("Not setting Accept HTTP header in response", nil)
	}

	if len(w.Header().Get("Accept-Language")) == 0 {
		utils.Fatalln("Not setting Accept-Language HTTP header in response", nil)
	}

	if len(w.Header().Get("Accept-Charset")) == 0 {
		utils.Fatalln("Not setting Accept-Charset HTTP header in response", nil)
	}

	if len(w.Header().Get("Date")) == 0 {
		utils.Fatalln("Not setting Date HTTP header in response", nil)
	}
}

func checkWrongACL(mux rest.Mux) {
	r, err := http.NewRequest("GET", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	r.RemoteAddr = "127.0.0.1"

	buildHTTPHeader(r, nil)

	_, cidr, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		utils.Fatalln("Error parsing CIDR", err)
	}
	mux.ACL = append(mux.ACL, cidr)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		utils.Fatalln("Not checking ACL", nil)
	}
}

func deleteDomain(mux *rest.Mux) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	buildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		utils.Fatalln("Error removing a domain object", nil)
	}
}

func buildHTTPHeader(r *http.Request, content []byte) {
	if r.ContentLength > 0 {
		r.Header.Set("Content-Type", check.SupportedContentType)

		hash := md5.New()
		hash.Write(content)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

		r.Header.Set("Content-MD5", hashBase64)
	}

	r.Header.Set("Date", time.Now().Format(time.RFC1123))

	stringToSign, err := check.BuildStringToSign(r, "1")
	if err != nil {
		utils.Fatalln("Error creating authorization", err)
	}

	signature := check.GenerateSignature(stringToSign, config.ShelterConfig.RESTServer.Secrets["1"])
	r.Header.Set("Authorization", fmt.Sprintf("%s %d:%s", check.SupportedNamespace, 1, signature))
}
