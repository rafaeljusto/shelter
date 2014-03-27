package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"net/http"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "ClientDomain"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for client domain test")
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

	database, databaseSession, err := mongodb.Open(
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}
	defer databaseSession.Close()

	// If there was some problem in the last test, there could be some data in the database,
	// so let's clear it to don't affect this test. We avoid checking the error, because if
	// the collection does not exist yet, it will be created in the first insert
	utils.ClearDatabase(database)

	finishRESTServer := utils.StartRESTServer()
	defer finishRESTServer()
	utils.StartWebClient()

	domainLifeCycle()

	utils.Println("SUCCESS!")
}

func domainLifeCycle() {
	data := []struct {
		method         string
		uri            string
		expectedStatus int
		content        string
	}{
		{
			method:         "PUT",
			uri:            "/domain/example.com.br.",
			expectedStatus: http.StatusCreated,
			content: `{
  "Nameservers": [
    { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
    { "host": "ns2.example.com.br.", "ipv6": "::1" }
  ],
  "Owners": [
    { "email": "admin@example.com.br.", "language": "pt-br" }
  ]
}`,
		},
		{
			method:         "PUT",
			uri:            "/domain/example.com.br.",
			expectedStatus: http.StatusNoContent,
			content: `{
  "Nameservers": [
    { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" }
  ],
  "Owners": [
    { "email": "admin2@example.com.br.", "language": "en-us" }
  ]
}`,
		},
		{
			method:         "GET",
			uri:            "/domain/example.com.br.",
			expectedStatus: http.StatusOK,
			content:        "",
		},
		{
			method:         "DELETE",
			uri:            "/domain/example.com.br.",
			expectedStatus: http.StatusNoContent,
			content:        "",
		},
	}

	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	for _, item := range data {
		var r *http.Request
		var err error

		if len(item.content) > 0 {
			r, err = http.NewRequest(item.method, fmt.Sprintf("%s%s", url, item.uri),
				bytes.NewReader([]byte(item.content)))

		} else {
			r, err = http.NewRequest(item.method, fmt.Sprintf("%s%s", url, item.uri), nil)
		}

		if err != nil {
			utils.Fatalln("Error creating the HTTP request", err)
		}

		utils.BuildHTTPHeader(r, []byte(item.content))

		response, err := client.Do(r)
		if err != nil {
			utils.Fatalln("Error sending request", err)
		}

		if response.StatusCode != item.expectedStatus {
			responseContent, err := ioutil.ReadAll(response.Body)

			if err == nil {
				utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method %s and URI %s",
					item.expectedStatus, response.StatusCode, item.method, item.uri),
					errors.New(string(responseContent)),
				)
			} else {
				utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method %s and URI %s",
					item.expectedStatus, response.StatusCode, item.method, item.uri), nil)
			}
		}
	}
}
