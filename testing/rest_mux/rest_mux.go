package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"shelter/config"
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/net/http/rest"
	"shelter/net/http/rest/check"
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
	createDomain(&mux)
	deleteDomain(&mux)

	utils.Println("SUCCESS!")
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
		utils.Fatalln("Error creting the HTTP request", err)
	}

	buildHTTPHeader(r, content)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		utils.Fatalln(fmt.Sprintf("Error creting a domain object. Expected HTTP status %d and got %d",
			http.StatusCreated, w.Code), nil)
	}
}

func deleteDomain(mux *rest.Mux) {
	r, err := http.NewRequest("DELETE", "/domain/example.com.br.", nil)
	if err != nil {
		utils.Fatalln("Error creting the HTTP request", err)
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
		r.Header.Add("Content-Type", check.SupportedContentType)

		hash := md5.New()
		hash.Write(content)
		hashBytes := hash.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)

		r.Header.Add("Content-MD5", hashBase64)
	}

	r.Header.Add("Date", time.Now().Format(time.RFC1123))

	stringToSign, err := check.BuildStringToSign(r, "1")
	if err != nil {
		utils.Fatalln("Error creating authorization", err)
	}

	signature := check.GenerateSignature(stringToSign, config.ShelterConfig.RESTServer.Secrets["1"])
	r.Header.Add("Authorization", fmt.Sprintf("%s %d:%s", check.SupportedNamespace, 1, signature))
}
