// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/rafaeljusto/shelter/config"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	basePath             = "/usr/shelter"
	configFilePath       = basePath + "/etc/shelter.conf"
	sampleConfigFilePath = basePath + "/etc/shelter.conf.unix.sample"
)

var (
	scanModule         = false
	restModule         = false
	webClientModule    = false
	notificationModule = false
)

var (
	hostnameOrIPInput = regexp.MustCompile("[0-9A-Za-z\\-\\.\\:]")
	hostnameInput     = regexp.MustCompile("[0-9A-Za-z\\-\\.]")
	alphaNumericInput = regexp.MustCompile("[0-9A-Za-z]")
	numericInput      = regexp.MustCompile("[0-9]")
	ipRangeInput      = regexp.MustCompile("[0-9a-fA-F\\:\\./]")
	secretAlphabet    = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")
)

type Option struct {
	Value    string
	Selected bool
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Initialize configuration file

		if err := config.LoadConfig(sampleConfigFilePath); err != nil {
			log.Fatalln(err)
			return
		}

		if !readEnabledModules() ||
			!readDatabaseParameters() ||
			!readRESTParameters() ||
			!readWebClientParameters() ||
			!readNotificationParameters() {

			return
		}

		jsonConfig, err := json.MarshalIndent(config.ShelterConfig, " ", " ")
		if err != nil {
			log.Fatalln(err)
			return
		}

		if err := ioutil.WriteFile(configFilePath, jsonConfig, 0664); err != nil {
			log.Fatalln(err)
			return
		}

	} else {
		// Update current configuration file
	}

	cmd := exec.Command("initctl", "stop", "shelter")
	cmd.Run() // We don't care about stop errors, because maybe the process wasn't there

	cmd = exec.Command("initctl", "start", "shelter")
	if err := cmd.Run(); err != nil {
		log.Println("Error starting shelter. Details:", err)
	}
}

func readEnabledModules() bool {
	options := []Option{
		{Value: "Scan", Selected: true},
		{Value: "REST", Selected: true},
		{Value: "Web Client", Selected: true},
		{Value: "Notification", Selected: true},
	}

	title := "Modules"
	description := "Please select the modules that are going to be enabled:"

	options, continueProcessing :=
		manageInputOptionsScreen(title, description, options, nil)

	if !continueProcessing {
		return false
	}

	for index, option := range options {
		switch index {
		case 0:
			if option.Selected {
				scanModule = true
			}

		case 1:
			if option.Selected {
				restModule = true
			}

		case 2:
			if option.Selected {
				webClientModule = true
			}

		case 3:
			if option.Selected {
				notificationModule = true
			}
		}
	}

	config.ShelterConfig.Scan.Enabled = scanModule
	config.ShelterConfig.RESTServer.Enabled = restModule
	config.ShelterConfig.WebClient.Enabled = webClientModule
	config.ShelterConfig.Notification.Enabled = notificationModule

	return true
}

func readDatabaseParameters() bool {
	var address, name string
	var port int
	var continueProcessing bool

	if restModule || notificationModule || scanModule {
		address, continueProcessing = readDatabaseHost()
		if !continueProcessing {
			return false
		}

		port, continueProcessing = readDatabasePort()
		if !continueProcessing {
			return false
		}

		name, continueProcessing = readDatabaseName()
		if !continueProcessing {
			return false
		}
	}

	config.ShelterConfig.Database.URI = fmt.Sprintf("%s:%d", address, port)
	config.ShelterConfig.Database.Name = name
	return true
}

func readDatabaseHost() (string, bool) {
	host := "localhost_________________________________________"
	title := "Database Configurations"
	description := "Please inform the host (IP or domain) of MongoDB database:"

	return manageInputTextScreen(title, description, host, hostnameOrIPInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "Host cannot be empty"
			}

			return true, ""
		})
}

func readDatabasePort() (int, bool) {
	port := "27017"
	title := "Database Configurations"
	description := "Please inform the port for the MongoDB database:"
	return manageInputNumberScreen(title, description, port)
}

func readDatabaseName() (string, bool) {
	name := "shelter___________________________________________"
	title := "Database Configurations"
	description := "Please inform the name of the MongoDB database:"
	return manageInputTextScreen(title, description, name, alphaNumericInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "Database name cannot be empty"
			}

			return true, ""
		})
}

func readRESTParameters() bool {
	return readRESTListeners() && readRESTACL() && readRESTSecret()
}

func readRESTListeners() bool {
	var addresses []string
	var port int
	var useTLS, generateCerts bool
	var continueProcessing bool

	if restModule || webClientModule {
		port, addresses, continueProcessing = readRESTPortAddresses()
		if !continueProcessing {
			return false
		}
	}

	if restModule {
		useTLS, generateCerts, continueProcessing = readRESTTLS()
		if !continueProcessing {
			return false
		}

		if generateCerts {
			hostname, continueProcessing := readRESTCertsParams()
			if !continueProcessing {
				return false
			}

			cert, key := generateCertificates("rest", hostname)
			config.ShelterConfig.RESTServer.TLS.CertificatePath = "etc/keys/" + cert
			config.ShelterConfig.RESTServer.TLS.PrivateKeyPath = "etc/keys/" + key
		}
	}

	config.ShelterConfig.RESTServer.Listeners = []struct {
		IP   string
		Port int
		TLS  bool
	}{}

	for _, address := range addresses {
		config.ShelterConfig.RESTServer.Listeners = append(config.ShelterConfig.RESTServer.Listeners, struct {
			IP   string
			Port int
			TLS  bool
		}{
			IP:   address,
			Port: port,
			TLS:  useTLS,
		})
	}

	return true
}

func readRESTPortAddresses() (int, []string, bool) {
	port := "4443_"

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return 0, nil, true
	}

	var options []Option
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err)
			return 0, nil, true
		}

		for _, a := range addrs {
			ip := a.String()
			ip = ip[:strings.Index(ip, "/")]

			// Avoid fe80::/10 prefix because is used only for link local
			if strings.HasPrefix(ip, "fe80") {
				continue
			}

			options = append(options, Option{
				Value:    fmt.Sprintf("%s", ip),
				Selected: true,
			})
		}
	}

	title := "REST Configurations"
	description := "Please inform the port number and select the IP addresses that the REST server will listen:"

	port, options, continueProcessing :=
		manageInputTextOptionsScreen(title, description, port, numericInput, func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "You must inform a port number"
			}

			return true, ""
		}, options, nil)

	if !continueProcessing {
		return 0, nil, false
	}

	portConverted, _ := strconv.Atoi(port)

	var selectedAddresses []string
	for _, option := range options {
		if option.Selected {
			selectedAddresses = append(selectedAddresses, option.Value)
		}
	}

	return portConverted, selectedAddresses, true
}

func readRESTTLS() (useTLS, generateCerts, continueProcessing bool) {
	options := []Option{
		{Value: "Use TLS on interfaces (HTTPS)", Selected: true},
		{Value: "Generate self-signed certificates automatically (valid for 1 year)", Selected: true},
	}

	title := "REST Configurations"
	description := "Please select the following TLS options:"

	selectedOptions, continueProcessing :=
		manageInputOptionsScreen(title, description, options,
			func(options []Option, optionIndex int) []Option {
				// Automatically certificates generation cannot exist without TLS
				if optionIndex == 0 && !options[0].Selected {
					options[1].Selected = false

				} else if optionIndex > 0 && options[optionIndex].Selected {
					options[0].Selected = true
				}

				return options
			})

	if !continueProcessing {
		return
	}

	for index, option := range selectedOptions {
		switch index {
		case 0:
			if option.Selected {
				useTLS = true
			}

		case 1:
			if option.Selected {
				generateCerts = true
			}
		}
	}

	return
}

func readRESTCertsParams() (string, bool) {
	host := "localhost_________________________________________"
	title := "REST Configurations"
	description := "Please inform the hostname of the certificate:"
	return manageInputTextScreen(title, description, host, hostnameInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "Certificate hostname cannot be empty"
			}

			return true, ""
		})
}

func readRESTACL() bool {
	acl := "127.0.0.1/8_______________________________________"
	title := "REST Configurations"
	description := "Please inform IP ranges that will have access to " +
		"the REST server (separeted by comma):"

	acl, continueProcessing :=
		manageInputTextScreen(title, description, acl, ipRangeInput,
			func(input string) (bool, string) {
				if len(input) == 0 {
					return false, "ACL cannot be empty"
				}

				aclParts := strings.Split(input, ",")
				for _, aclPart := range aclParts {
					aclPart = strings.TrimSpace(aclPart)
					if _, _, err := net.ParseCIDR(aclPart); err != nil {
						return false, "IP range " + aclPart + " is invalid"
					}
				}

				return true, ""
			})

	if !continueProcessing {
		return false
	}

	config.ShelterConfig.RESTServer.ACL = []string{}
	aclParts := strings.Split(acl, ",")
	for _, aclPart := range aclParts {
		aclPart = strings.TrimSpace(aclPart)
		config.ShelterConfig.RESTServer.ACL =
			append(config.ShelterConfig.RESTServer.ACL, aclPart)
	}

	return true
}

func readRESTSecret() bool {
	keyId, generateAutomatically, continueProcessing := readRESTSecretId()

	if !continueProcessing {
		return false
	}

	var secret string
	if generateAutomatically {
		for i := 0; i < 30; i++ {
			randNumber, err := rand.Int(rand.Reader, big.NewInt(int64(len(secretAlphabet))))
			if err != nil {
				log.Println("Error generating random numbers. Details:", err)
				return false
			}

			secret += string(secretAlphabet[randNumber.Int64()])
		}

	} else {
		secret, continueProcessing = readRESTSecretContent(keyId)
		if !continueProcessing {
			return false
		}
	}

	config.ShelterConfig.RESTServer.Secrets[keyId] = secret
	return true
}

func readRESTSecretId() (keyId string, generateSecret bool, continueProcessing bool) {
	keyId = "key01_______________"
	options := []Option{
		{Value: "Generate shared secret automatically", Selected: true},
	}

	title := "REST Configurations"
	description := "Please inform the shared secret identification:"

	keyId, options, continueProcessing =
		manageInputTextOptionsScreen(title, description, keyId, alphaNumericInput,
			func(input string) (bool, string) {
				if len(input) == 0 {
					return false, "Certificate hostname cannot be empty"
				}

				return true, ""
			}, options, nil)

	return keyId, options[0].Selected, continueProcessing
}

func readRESTSecretContent(keyId string) (string, bool) {
	secret := "__________________________________________________"
	title := "REST Configurations"
	description := "Please inform the shared secret for " + keyId + ":"

	return manageInputTextScreen(title, description, secret, alphaNumericInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "Certificate hostname cannot be empty"
			}

			return true, ""
		})
}

func readWebClientParameters() bool {
	var addresses []string
	var port int
	var useTLS, generateCerts, useRESTCerts bool
	var continueProcessing bool

	if webClientModule {
		port, addresses, continueProcessing = readWebClientPortAddresses()
		if !continueProcessing {
			return false
		}

		useTLS, generateCerts, useRESTCerts, continueProcessing = readWebClientTLS()
		if !continueProcessing {
			return false
		}

		if generateCerts && !useRESTCerts {
			hostname, continueProcessing := readWebClientCertsParams()
			if !continueProcessing {
				return false
			}

			cert, key := generateCertificates("webclient", hostname)
			config.ShelterConfig.WebClient.TLS.CertificatePath = "etc/keys/" + cert
			config.ShelterConfig.WebClient.TLS.PrivateKeyPath = "etc/keys/" + key

		} else if !generateCerts && useRESTCerts {
			config.ShelterConfig.WebClient.TLS.CertificatePath =
				config.ShelterConfig.RESTServer.TLS.CertificatePath
			config.ShelterConfig.WebClient.TLS.PrivateKeyPath =
				config.ShelterConfig.RESTServer.TLS.PrivateKeyPath
		}
	}

	config.ShelterConfig.WebClient.Listeners = []struct {
		IP   string
		Port int
		TLS  bool
	}{}

	for _, address := range addresses {
		config.ShelterConfig.WebClient.Listeners =
			append(config.ShelterConfig.WebClient.Listeners, struct {
				IP   string
				Port int
				TLS  bool
			}{
				IP:   address,
				Port: port,
				TLS:  useTLS,
			})
	}

	return true
}

func readWebClientPortAddresses() (int, []string, bool) {
	port := "4444_"

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return 0, nil, true
	}

	var options []Option
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err)
			return 0, nil, true
		}

		for _, a := range addrs {
			ip := a.String()
			ip = ip[:strings.Index(ip, "/")]

			// Avoid fe80::/10 prefix because is used only for link local
			if strings.HasPrefix(ip, "fe80") {
				continue
			}

			options = append(options, Option{
				Value:    fmt.Sprintf("%s", ip),
				Selected: true,
			})
		}
	}

	title := "Web Client Configurations"
	description := "Please inform the port number and select the IP addresses " +
		"that the Web Client will listen:"

	port, options, continueProcessing :=
		manageInputTextOptionsScreen(title, description, port, numericInput,
			func(input string) (bool, string) {
				if len(input) == 0 {
					return false, "You must inform a port number"
				}

				return true, ""
			}, options, nil)

	if !continueProcessing {
		return 0, nil, false
	}

	portConverted, _ := strconv.Atoi(port)

	var selectedAddresses []string
	for _, option := range options {
		if option.Selected {
			selectedAddresses = append(selectedAddresses, option.Value)
		}
	}
	return portConverted, selectedAddresses, true
}

func readWebClientTLS() (useTLS, generateCerts, useRESTCerts, continueProcessing bool) {
	options := []Option{
		{Value: "Use TLS on interfaces (HTTPS)", Selected: true},
		{Value: "Generate self-signed certificates automatically (valid for 1 year)", Selected: true},
	}

	if len(config.ShelterConfig.RESTServer.Listeners) > 0 &&
		config.ShelterConfig.RESTServer.Listeners[0].TLS {
		options = append(options, Option{
			Value:    "Use same certificate from REST server",
			Selected: true,
		})
		options[1].Selected = false
	}

	title := "Web Client Configurations"
	description := "Please select the following TLS options:"

	options, continueProcessing =
		manageInputOptionsScreen(title, description, options,
			func(options []Option, optionIndex int) []Option {
				// Automatically certificates generation and use REST certificate cannot exist
				// without TLS

				if optionIndex == 0 && !options[0].Selected {
					options[1].Selected = false

					if len(options) == 3 {
						options[2].Selected = false
					}

				} else if optionIndex > 0 && options[optionIndex].Selected {
					options[0].Selected = true

					if optionIndex == 1 && len(options) == 3 {
						options[2].Selected = false
					} else if optionIndex == 2 {
						options[1].Selected = false
					}
				}

				return options
			})

	if !continueProcessing {
		return
	}

	for index, option := range options {
		switch index {
		case 0:
			if option.Selected {
				useTLS = true
			}

		case 1:
			if option.Selected {
				generateCerts = true
			}

		case 2:
			if option.Selected {
				useRESTCerts = true
			}
		}
	}

	return
}

func readWebClientCertsParams() (string, bool) {
	host := "localhost_________________________________________"
	title := "Web Client Configurations"
	description := "Please inform the hostname of the certificate:"
	return manageInputTextScreen(title, description, host, hostnameInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "Certificate hostname cannot be empty"
			}

			return true, ""
		})
}

func readNotificationParameters() bool {
	var smtpServer, authenticationType, username, password string
	var port int
	var continueProcessing bool

	if notificationModule {
		smtpServer, continueProcessing = readNotificationSMTPServer()
		if !continueProcessing {
			return false
		}

		port, continueProcessing = readNotificationPort()
		if !continueProcessing {
			return false
		}

		authenticationType, continueProcessing = readNotificationAuthType()
		if !continueProcessing {
			return false
		}

		if authenticationType != "None" {
			username, continueProcessing = readNotificationUsername()
			if !continueProcessing {
				return false
			}

			password, continueProcessing = readNotificationPassword()
			if !continueProcessing {
				return false
			}
		}
	}

	config.ShelterConfig.Notification.SMTPServer.Server = smtpServer
	config.ShelterConfig.Notification.SMTPServer.Port = port

	switch authenticationType {
	case "None":
		config.ShelterConfig.Notification.SMTPServer.Auth.Type = config.AuthenticationTypeNone
	case "Plain":
		config.ShelterConfig.Notification.SMTPServer.Auth.Type = config.AuthenticationTypePlain
	case "CRAM-MD5":
		config.ShelterConfig.Notification.SMTPServer.Auth.Type = config.AuthenticationTypeCRAMMD5Auth
	}

	config.ShelterConfig.Notification.SMTPServer.Auth.Username = username
	config.ShelterConfig.Notification.SMTPServer.Auth.Password = password
	return true
}

func readNotificationSMTPServer() (string, bool) {
	host := "smtp.gmail.com.___________________________________"
	title := "Notification Configurations"
	description := "Please inform the hostname or IP address of the SMTP server:"
	return manageInputTextScreen(title, description, host, hostnameOrIPInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "Hostname or IP field cannot be empty"
			}

			return true, ""
		})
}

func readNotificationPort() (int, bool) {
	port := "587_"
	title := "Notification Configurations"
	description := "Please inform the port for the SMTP server:"
	return manageInputNumberScreen(title, description, port)
}

func readNotificationAuthType() (string, bool) {
	options := []Option{
		{Value: "None", Selected: false},
		{Value: "Plain", Selected: true},
		{Value: "CRAM-MD5", Selected: false},
	}

	title := "Notification Configurations"
	description := "Please inform the SMTP authentication type:"

	options, continueProcessing :=
		manageInputOptionsScreen(title, description, options,
			func(options []Option, optionIndex int) []Option {
				for index, _ := range options {
					if options[optionIndex].Selected && index != optionIndex {
						options[index].Selected = false

					} else if !options[optionIndex].Selected && index == optionIndex {
						options[index].Selected = true
					}
				}

				return options
			})

	if !continueProcessing {
		return "", false
	}

	for _, option := range options {
		if option.Selected {
			return option.Value, true
		}
	}

	return "", false
}

func readNotificationUsername() (string, bool) {
	host := "anonymous_________________________________________"
	title := "Notification Configurations"
	description := "Please inform the SMTP username:"
	return manageInputTextScreen(title, description, host, alphaNumericInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "SMTP username cannot be empty"
			}

			return true, ""
		})
}

func readNotificationPassword() (string, bool) {
	host := "__________________________________________________"
	title := "Notification Configurations"
	description := "Please inform the SMTP password:"
	return manageInputTextScreen(title, description, host, alphaNumericInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "SMTP password cannot be empty"
			}

			return true, ""
		})
}

func readInput(inputsDraw func(), inputsAction func(termbox.Event) bool) bool {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	draw(inputsDraw)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				return false
			}

			if !inputsAction(ev) {
				return true
			}

		case termbox.EventResize:
			draw(inputsDraw)
		}
	}
}

func draw(inputsDraw func()) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	windowWidth, windowsHeight := termbox.Size()

	// Background
	for i := 0; i < windowWidth; i++ {
		for j := 0; j < windowsHeight; j++ {
			termbox.SetCell(i, j, 0x0, termbox.ColorWhite, termbox.ColorBlue)
		}
	}

	// Edges
	termbox.SetCell(0, 0, 0x250C, termbox.ColorWhite, termbox.ColorBlue)
	termbox.SetCell(windowWidth-1, 0, 0x2510, termbox.ColorWhite, termbox.ColorBlue)
	termbox.SetCell(0, windowsHeight-1, 0x2514, termbox.ColorWhite, termbox.ColorBlue)
	termbox.SetCell(windowWidth-1, windowsHeight-1, 0x2518, termbox.ColorWhite, termbox.ColorBlue)

	// Logo
	termbox.SetCell((windowWidth/2)-7, 2, 0xff33, termbox.ColorWhite, termbox.ColorBlue) // S
	termbox.SetCell((windowWidth/2)-5, 2, 0xff28, termbox.ColorWhite, termbox.ColorBlue) // H
	termbox.SetCell((windowWidth/2)-3, 2, 0xff25, termbox.ColorWhite, termbox.ColorBlue) // E
	termbox.SetCell((windowWidth/2)-1, 2, 0xff2c, termbox.ColorWhite, termbox.ColorBlue) // L
	termbox.SetCell((windowWidth/2)+1, 2, 0xff34, termbox.ColorWhite, termbox.ColorBlue) // T
	termbox.SetCell((windowWidth/2)+3, 2, 0xff25, termbox.ColorWhite, termbox.ColorBlue) // E
	termbox.SetCell((windowWidth/2)+5, 2, 0xff32, termbox.ColorWhite, termbox.ColorBlue) // R

	// Footer tip
	writeText("Press ESC to quit", windowWidth-20, windowsHeight-2)

	inputsDraw()

	termbox.Flush()
}

func writeTitle(text string, x, y int) {
	xTmp := x
	for _, character := range text {
		termbox.SetCell(xTmp, y, rune(character), termbox.ColorWhite, termbox.ColorBlue)
		xTmp += 1
	}

	xTmp = x
	for xTmp <= len(text)+1 {
		termbox.SetCell(xTmp, y+1, 0x2015, termbox.ColorWhite, termbox.ColorBlue)
		xTmp += 1
	}
}

func writeText(text string, x, y int) {
	for _, character := range text {
		termbox.SetCell(x, y, rune(character), termbox.ColorWhite, termbox.ColorBlue)
		x += 1
	}
}

func writeOptions(options []Option, x, y int) {
	for index, option := range options {
		writeText("[ ] "+option.Value, x, y+index)
	}
}

func manageInputOptionsScreen(
	title, description string,
	options []Option,
	checkConsistency func([]Option, int) []Option,
) ([]Option, bool) {

	overOption := -1

	inputsDraw := func() {
		writeTitle(title, 2, 4)
		writeText(description, 2, 7)
		writeOptions(options, 2, 9)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Continue", 2, windowsHeight-2)

		for index, option := range options {
			if option.Selected {
				termbox.SetCell(3, 9+index, 0x221a, termbox.ColorYellow, termbox.ColorBlue)
			}
		}

		if overOption > -1 {
			termbox.SetCell(2, 9+overOption, rune('['), termbox.ColorYellow, termbox.ColorBlue)
			termbox.SetCell(4, 9+overOption, rune(']'), termbox.ColorYellow, termbox.ColorBlue)
		}
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyTab:
			// Move to the next option

			overOption += 1
			if overOption >= len(options) {
				overOption = 0
			}

		case termbox.KeySpace:
			// Select the option

			var optionIndex int
			for index, option := range options {
				if index == overOption {
					optionIndex = index
					options[index].Selected = !option.Selected
					break
				}
			}

			if checkConsistency != nil {
				options = checkConsistency(options, optionIndex)
			}

		case termbox.KeyEnter:
			// Finish reading inputs
			return false
		}

		draw(inputsDraw)
		return true
	}

	if !readInput(inputsDraw, restInputsAction) {
		return nil, false
	}

	return options, true
}

func manageInputTextScreen(
	title, description, input string,
	allowedInput *regexp.Regexp,
	validate func(string) (bool, string),
) (string, bool) {

	inputPosition := 0

	inputsDraw := func() {
		writeTitle(title, 2, 4)
		writeText(description, 2, 7)
		writeText(input, 2, 9)

		if inputPosition < len(input) {
			termbox.SetCell(2+inputPosition, 9, rune(input[inputPosition]), termbox.ColorWhite, termbox.ColorYellow)
		}

		_, windowsHeight := termbox.Size()
		writeText("[ENTER] Continue", 2, windowsHeight-2)
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			inputPosition -= 1
			if inputPosition < 0 {
				inputPosition = 0
			}

			input = input[:inputPosition] + "_" + input[inputPosition+1:]

		case termbox.KeyDelete:
			if inputPosition < len(input) {
				input = input[:inputPosition] + input[inputPosition+1:] + "_"
			}

		case termbox.KeyEnter:
			if validate != nil {
				inputTmp := strings.Replace(input, "_", "", -1)
				valid, msg := validate(inputTmp)

				if !valid {
					draw(func() {
						inputsDraw()
						writeText("ERROR: "+msg, 2, 11)
					})

					return true
				}
			}

			// Finish reading inputs
			return false

		default:
			if allowedInput.MatchString(string(ev.Ch)) &&
				inputPosition < len(input) {

				input = input[:inputPosition] + string(ev.Ch) + input[inputPosition+1:]

				inputPosition += 1
				if inputPosition > len(input) {
					inputPosition = len(input)
				}
			}
		}

		draw(inputsDraw)
		return true
	}

	if !readInput(inputsDraw, restInputsAction) {
		return "", false
	}

	return strings.Replace(input, "_", "", -1), true
}

func manageInputNumberScreen(title, description string, number string) (int, bool) {
	number, continueProcessing :=
		manageInputTextScreen(title, description, number, numericInput,
			func(input string) (bool, string) {
				if len(input) == 0 {
					return false, "You must inform a number"
				}

				return true, ""
			})

	if !continueProcessing {
		return 0, false
	}

	numberConverted, _ := strconv.Atoi(number)
	return numberConverted, true
}

func manageInputTextOptionsScreen(
	title, description string,
	input string,
	allowedInput *regexp.Regexp,
	validate func(string) (bool, string),
	options []Option,
	checkConsistency func([]Option, int) []Option,
) (string, []Option, bool) {

	inputPosition := 0
	overOption := 0

	inputsDraw := func() {
		writeTitle(title, 2, 4)
		writeText(description, 2, 7)
		writeText(input, 2, 9)
		writeOptions(options, 2, 11)

		_, windowsHeight := termbox.Size()
		writeText("[TAB] Move over options", 2, windowsHeight-4)
		writeText("[SPACE] Select an option", 2, windowsHeight-3)
		writeText("[ENTER] Continue", 2, windowsHeight-2)

		for index, option := range options {
			if option.Selected {
				termbox.SetCell(3, 11+index, 0x221a, termbox.ColorYellow, termbox.ColorBlue)
			}
		}

		if inputPosition < len(input) && overOption == 0 {
			termbox.SetCell(2+inputPosition, 9, rune(input[inputPosition]), termbox.ColorWhite, termbox.ColorYellow)
		}

		if overOption > 0 {
			termbox.SetCell(2, 10+overOption, rune('['), termbox.ColorYellow, termbox.ColorBlue)
			termbox.SetCell(4, 10+overOption, rune(']'), termbox.ColorYellow, termbox.ColorBlue)
		}
	}

	restInputsAction := func(ev termbox.Event) bool {
		switch ev.Key {
		case termbox.KeyTab:
			// Move to the next option

			overOption += 1
			if overOption >= len(options)+1 {
				overOption = 0
			}

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if overOption != 0 {
				break
			}

			inputPosition -= 1
			if inputPosition < 0 {
				inputPosition = 0
			}

			input = input[:inputPosition] + "_" + input[inputPosition+1:]

		case termbox.KeyDelete:
			if overOption != 0 {
				break
			}

			if inputPosition < len(input) {
				input = input[:inputPosition] + input[inputPosition+1:] + "_"
			}

		case termbox.KeySpace:
			if overOption == 0 {
				break
			}

			// Select the option

			var optionIndex int
			for index, option := range options {
				if index == overOption-1 {
					optionIndex = index
					options[index].Selected = !option.Selected
					break
				}
			}

			if checkConsistency != nil {
				options = checkConsistency(options, optionIndex)
			}

		case termbox.KeyEnter:
			if validate != nil {
				inputTmp := strings.Replace(input, "_", "", -1)
				valid, msg := validate(inputTmp)

				if !valid {
					draw(func() {
						inputsDraw()
						writeText("ERROR: "+msg, 2, 12+len(options))
					})

					return true
				}
			}

			// Finish reading inputs
			return false

		default:
			if overOption != 0 {
				break
			}

			if allowedInput.MatchString(string(ev.Ch)) &&
				inputPosition < len(input) {

				input = input[:inputPosition] + string(ev.Ch) + input[inputPosition+1:]

				inputPosition += 1
				if inputPosition > len(input) {
					inputPosition = len(input)
				}
			}
		}

		draw(inputsDraw)
		return true
	}

	if !readInput(inputsDraw, restInputsAction) {
		return "", nil, false
	}

	return strings.Replace(input, "_", "", -1), options, true
}

func generateCertificates(prefix, hostname string) (cert, key string) {
	err := os.MkdirAll(basePath+"/etc/keys", os.ModeDir|0600)
	if err != nil {
		log.Println("Error creating certificates directory. Details:", err)
		return
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("Error creating certificates. Details:", err)
		return
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	// end of ASN.1 time
	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)
	if notAfter.After(endOfTime) {
		notAfter = endOfTime
	}

	name := pkix.Name{
		CommonName:         hostname,
		Organization:       []string{"Shelter project"},
		OrganizationalUnit: []string{"TI"},
		SerialNumber:       "1",
	}

	template := x509.Certificate{
		Version:               1,
		SerialNumber:          new(big.Int).SetInt64(1),
		Issuer:                name,
		Subject:               name,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(hostname); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, hostname)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template,
		&template, &priv.PublicKey, priv)
	if err != nil {
		log.Println("Error creating certificates. Details:", err)
		return
	}

	cert = prefix + "-cert.pem"
	certOut, err := os.Create(basePath + "/etc/keys/" + cert)
	if err != nil {
		log.Println("Error creating certificates. Details:", err)
		return
	}

	pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	certOut.Close()

	key = prefix + "-key.pem"
	keyOut, err := os.OpenFile(basePath+"/etc/keys/"+key,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println("Error creating certificates. Details:", err)
		return
	}

	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	keyOut.Close()

	return
}

func moveFile(dst, orig string) error {
	file, err := os.Open(orig)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(dstFile, file); err != nil {
		return err
	}

	if err := os.Remove(orig); err != nil {
		return err
	}

	return nil
}
