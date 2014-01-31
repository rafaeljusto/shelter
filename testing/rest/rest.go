package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/net/http/rest"
	"github.com/rafaeljusto/shelter/net/http/rest/check"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
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

// ScanQuerierTestConfigFile is a structure to store the test configuration file data
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
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}
	defer databaseSession.Close()

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// If there was some problem in the last test, there could be some data in the
	// database, so let's clear it to don't affect this test. We avoid checking the error,
	// because if the collection does not exist yet, it will be created in the first
	// insert
	domainDAO.RemoveAll()

	createMessagesFile()
	defer removeMessagesFile()
	startRESTServer()

	// REST performance report is optional and only generated when the report file path parameter is
	// given
	if report {
		if cpuProfile {
			profileFile, err := os.Create(restConfig.Report.Profile.CPUFile)
			if err != nil {
				utils.Fatalln("Error creating CPU profile file", err)
			}

			if err := pprof.StartCPUProfile(profileFile); err != nil {
				utils.Fatalln("Error starting CPU profile file", err)
			}

			defer pprof.StopCPUProfile()
		}

		defer func() {
			if memoryProfile {
				runtime.GC()

				profileFile, err := os.Create(restConfig.Report.Profile.MemoryFile)
				if err != nil {
					utils.Fatalln("Error creating memory profile file", err)
				}

				if err := pprof.Lookup("heap").WriteTo(profileFile, 1); err != nil {
					utils.Fatalln("Error writing to memory profile file", err)
				}
				profileFile.Close()
			}

			if goProfile {
				profileFile, err := os.Create(restConfig.Report.Profile.GoRoutinesFile)
				if err != nil {
					utils.Fatalln("Error creating Go routines profile file", err)
				}

				if err := pprof.Lookup("goroutine").WriteTo(profileFile, 2); err != nil {
					utils.Fatalln("Error writing to Go routines profile file", err)
				}
				profileFile.Close()
			}
		}()

		restReport(restConfig)
	}

	utils.Println("SUCCESS!")
}

func restReport(restConfig RESTTestConfigFile) {
	report := " #       | Operation | Total            | DPS  | Memory (MB)\n" +
		"---------------------------------------------------------------\n"

	// Report variables
	scale := []int{10, 50, 100, 500, 1000, 5000,
		10000, 50000, 100000, 500000, 1000000, 5000000}

	content := []byte(`{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns2.example.com.br.", "ipv6": "::1" }
      ],
      "Owners": [
        "admin@example.com.br."
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

			buildHTTPHeader(r, content)

			response, err := client.Do(r)
			if err != nil {
				utils.Println(fmt.Sprintf("Error \"%s\" detected when sending request for action %s, "+
					"ignoring it...", err.Error(), action))
				continue
			}

			ioutil.ReadAll(response.Body)
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

func createMessagesFile() {
	languagePacks := messages.LanguagePacks{
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

	messagePath := filepath.Join(
		config.ShelterConfig.BasePath,
		config.ShelterConfig.RESTServer.LanguageConfigPath,
	)

	file, err := os.Create(messagePath)
	if err != nil {
		utils.Fatalln("Error creating messages file", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(languagePacks); err != nil {
		utils.Fatalln("Error encoding messages structure", err)
	}
}

func removeMessagesFile() {
	messagePath := filepath.Join(
		config.ShelterConfig.BasePath,
		config.ShelterConfig.RESTServer.LanguageConfigPath,
	)

	// We don't care if the file doesn't exists anymore, so we ignore the returned error
	os.Remove(messagePath)
}

func startRESTServer() {
	listeners, err := rest.Listen()
	if err != nil {
		utils.Fatalln("Error listening to interfaces", err)
	}

	if err := rest.Start(listeners); err != nil {
		utils.Fatalln("Error starting the REST server", err)
	}

	// Wait the REST server to start before testing
	time.Sleep(1 * time.Second)
}
