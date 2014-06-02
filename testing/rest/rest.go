// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

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
	"runtime"
	"time"
)

var (
	configFilePath string // Path for the config file with the connection information
	report         bool   // Flag to generate the REST performance report file
	cpuProfile     bool   // Write profile about the CPU when executing the report
	goProfile      bool   // Write profile about the Go routines when executing the report
	memoryProfile  bool   // Write profile about the memory usage when executing the report
)

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

type ByteSize float64

// RESTTestConfigFile is a structure to store the test configuration file data
type RESTTestConfigFile struct {
	config.Config

	Report struct {
		File    string
		Profile struct {
			CPUFile        string
			GoRoutinesFile string
			MemoryFile     string
		}
	}
}

func init() {
	utils.TestName = "REST"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for REST test")
	flag.BoolVar(&report, "report", false, "Report flag for REST performance")
	flag.BoolVar(&cpuProfile, "cpuprofile", false, "Report flag to enable CPU profile")
	flag.BoolVar(&goProfile, "goprofile", false, "Report flag to enable Go routines profile")
	flag.BoolVar(&memoryProfile, "memprofile", false, "Report flag to enable memory profile")
}

func main() {
	flag.Parse()

	var restConfig RESTTestConfigFile
	err := utils.ReadConfigFile(configFilePath, &restConfig)
	config.ShelterConfig = restConfig.Config

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	database, databaseSession, err := mongodb.Open(
		config.ShelterConfig.Database.URIs,
		config.ShelterConfig.Database.Name,
		config.ShelterConfig.Database.Auth.Enabled,
		config.ShelterConfig.Database.Auth.Username,
		config.ShelterConfig.Database.Auth.Password,
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

	domainLifeCycle()

	// REST performance report is optional and only generated when the report file path parameter is
	// given
	if report {
		if cpuProfile {
			f := utils.StartCPUProfile(restConfig.Report.Profile.CPUFile)
			defer f()
		}

		if memoryProfile {
			f := utils.StartMemoryProfile(restConfig.Report.Profile.MemoryFile)
			defer f()
		}

		if goProfile {
			f := utils.StartGoRoutinesProfile(restConfig.Report.Profile.GoRoutinesFile)
			defer f()
		}

		restReport(restConfig)
	}

	utils.Println("SUCCESS!")
}

func domainLifeCycle() {
	data := []struct {
		method         string
		uri            string
		expectedStatus int
		content        string
		expectedBody   string
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
			expectedBody:   `{"fqdn":"example.com.br.","nameservers":[{"host":"ns1.example.com.br.","ipv4":"127.0.0.1","lastStatus":"NOTCHECKED","lastCheckAt":"0001-01-01T00:00:00Z","lastOKAt":"0001-01-01T00:00:00Z"}],"owners":[{"email":"admin2@example.com.br.","language":"en-US"}],"links":[{"types":["self"],"href":"/domain/example.com.br."}]}`,
		},
		{
			method:         "DELETE",
			uri:            "/domain/example.com.br.",
			expectedStatus: http.StatusNoContent,
		},
	}

	var client http.Client

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
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

		var responseContent []byte
		if response.ContentLength > 0 {
			responseContent, err = ioutil.ReadAll(response.Body)
			if err != nil {
				utils.Fatalln(fmt.Sprintf("Error reading response for method %s and URI %s",
					item.method, item.uri),
					err,
				)
			}
		}

		if response.StatusCode != item.expectedStatus {
			utils.Fatalln(fmt.Sprintf("Expected HTTP status %d and got %d for method %s and URI %s",
				item.expectedStatus, response.StatusCode, item.method, item.uri),
				errors.New(string(responseContent)),
			)
		} else if string(responseContent) != item.expectedBody {
			utils.Fatalln(fmt.Sprintf("Expected HTTP body [%s] and got [%s] for method %s and URI %s",
				item.expectedBody, string(responseContent), item.method, item.uri), nil)
		}
	}
}

func restReport(restConfig RESTTestConfigFile) {
	report := " #       | Operation | Total            | DPS  | Memory (MB)\n" +
		"---------------------------------------------------------------\n"

	// Report variables
	scale := []int{10, 50, 100, 500, 1000, 5000,
		10000, 50000, 100000, 500000, 1000000, 5000000}

	content := []byte(`{
      "Nameservers": [
        { "host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        { "email": "admin@example.com.br.", "language": "pt-br" }
      ]
    }`)

	url := ""
	if len(config.ShelterConfig.RESTServer.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.RESTServer.Listeners[0].IP,
			config.ShelterConfig.RESTServer.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	for _, numberOfItems := range scale {
		utils.Println(fmt.Sprintf("Generating report - scale %d", numberOfItems))

		// Generate domains
		report += restActionReport(numberOfItems, content, func(i int) (*http.Request, error) {
			return http.NewRequest("PUT", fmt.Sprintf("%s/domain/example%d.com.br.", url, i),
				bytes.NewReader(content))
		}, http.StatusCreated, "CREATE")

		// Retrieve domains
		report += restActionReport(numberOfItems, content, func(i int) (*http.Request, error) {
			return http.NewRequest("GET", fmt.Sprintf("%s/domain/example%d.com.br.", url, i), nil)
		}, http.StatusOK, "RETRIEVE")

		// Delete domains
		report += restActionReport(numberOfItems, content, func(i int) (*http.Request, error) {
			return http.NewRequest("DELETE", fmt.Sprintf("%s/domain/example%d.com.br.", url, i), nil)
		}, http.StatusNoContent, "DELETE")
	}

	utils.WriteReport(restConfig.Report.File, report)
}

func restActionReport(numberOfItems int,
	content []byte,
	requestBuilder func(int) (*http.Request, error),
	expectedStatus int,
	action string) string {

	var client http.Client

	totalDuration, domainsPerSecond := calculateRESTDurations(func() {
		for i := 0; i < numberOfItems; i++ {
			r, err := requestBuilder(i)
			if err != nil {
				utils.Fatalln("Error creating the HTTP request", err)
			}

			utils.BuildHTTPHeader(r, content)

			response, err := client.Do(r)
			if err != nil {
				utils.Fatalln(fmt.Sprintf("Error detected when sending request for action %s"+
					" and URI \"%s\"", action, r.URL.RequestURI()), err)
			}

			_, err = ioutil.ReadAll(response.Body)
			if err != nil {
				utils.Fatalln(fmt.Sprintf("Error detected when reading the response body for action %s"+
					" and URI \"%s\"", action, r.URL.RequestURI()), err)
			}
			response.Body.Close()

			if response.StatusCode != expectedStatus {
				utils.Fatalln(fmt.Sprintf("Error with the domain object in the action %s. "+
					"Expected HTTP status %d and got %d",
					action, expectedStatus, response.StatusCode), nil)
			}
		}
	}, numberOfItems)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return fmt.Sprintf("% -8d | %-9s | %16s | %4d | %14.2f\n",
		numberOfItems,
		action,
		time.Duration(int64(totalDuration)).String(),
		domainsPerSecond,
		float64(memStats.Alloc)/float64(MB),
	)
}

func calculateRESTDurations(action func(), numberOfDomains int) (totalDuration time.Duration, domainsPerSecond int64) {
	beginTimer := time.Now()
	action()
	totalDuration = time.Since(beginTimer)

	totalDurationSeconds := int64(totalDuration / time.Second)
	if totalDurationSeconds > 0 {
		domainsPerSecond = int64(numberOfDomains) / totalDurationSeconds

	} else {
		domainsPerSecond = int64(numberOfDomains)
	}

	return
}
