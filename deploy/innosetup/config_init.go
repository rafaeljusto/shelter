// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/deploy"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	basePath             = "c:\\Program Files\\shelter"
	configFilePath       = basePath + "\\conf\\shelter.conf"
	sampleConfigFilePath = basePath + "\\conf\\shelter.conf.windows.sample"
)

var (
	scanModule         = false
	restModule         = false
	webClientModule    = false
	notificationModule = false
)

var (
	secretAlphabet = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")
)

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

	// TODO: (Re)start service

	fmt.Println("==========================================================================")
	fmt.Printf("Edit advanced configurations on %s\n", configFilePath)

	if webClientModule && len(config.ShelterConfig.WebClient.Listeners) > 0 {
		ln := config.ShelterConfig.WebClient.Listeners[0]
		url := ""

		if ln.TLS {
			url = fmt.Sprintf("https://%s:%d", ln.IP, ln.Port)
		} else {
			url = fmt.Sprintf("http://%s:%d", ln.IP, ln.Port)
		}

		fmt.Printf("Check the web client on %s\n", url)
	}

	fmt.Println("==========================================================================")
}

func readEnabledModules() bool {
	options := []deploy.Option{
		{Value: "Scan", Selected: true},
		{Value: "REST", Selected: true},
		{Value: "Web Client", Selected: true},
		{Value: "Notification", Selected: true},
	}

	title := "Modules"
	description := "Please select the modules that are going to be enabled:"

	options, continueProcessing :=
		deploy.ManageInputOptionsScreen(title, description, options, nil)

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

	return deploy.ManageInputTextScreen(title, description, host, deploy.HostnameOrIPInput,
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
	return deploy.ManageInputNumberScreen(title, description, port)
}

func readDatabaseName() (string, bool) {
	name := "shelter___________________________________________"
	title := "Database Configurations"
	description := "Please inform the name of the MongoDB database:"
	return deploy.ManageInputTextScreen(title, description, name, deploy.AlphaNumericInput,
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
	var useTLS bool
	var continueProcessing bool

	if restModule || webClientModule {
		port, addresses, continueProcessing = readRESTPortAddresses()
		if !continueProcessing {
			return false
		}
	}

	if restModule {
		useTLS, continueProcessing = readRESTTLS()
		if !continueProcessing {
			return false
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

	var options []deploy.Option
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

			options = append(options, deploy.Option{
				Value:    fmt.Sprintf("%s", ip),
				Selected: true,
			})
		}
	}

	title := "REST Configurations"
	description := "Please inform the port number and select the IP addresses that the REST server will listen:"

	port, options, continueProcessing :=
		deploy.ManageInputTextOptionsScreen(title, description, port, deploy.NumericInput,
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

func readRESTTLS() (useTLS, continueProcessing bool) {
	options := []deploy.Option{
		{Value: "Use TLS on interfaces (HTTPS)", Selected: true},
	}

	title := "REST Configurations"
	description := "Please select the following TLS options:"

	selectedOptions, continueProcessing :=
		deploy.ManageInputOptionsScreen(title, description, options, nil)

	if !continueProcessing {
		return
	}

	for index, option := range selectedOptions {
		switch index {
		case 0:
			if option.Selected {
				useTLS = true
			}
		}
	}

	return
}

func readRESTACL() bool {
	acl := "127.0.0.1/8_______________________________________"
	title := "REST Configurations"
	description := "Please inform IP ranges that will have access to " +
		"the REST server (separeted by comma):"

	acl, continueProcessing :=
		deploy.ManageInputTextScreen(title, description, acl, deploy.IPRangeInput,
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
	options := []deploy.Option{
		{Value: "Generate shared secret automatically", Selected: true},
	}

	title := "REST Configurations"
	description := "Please inform the shared secret identification:"

	keyId, options, continueProcessing =
		deploy.ManageInputTextOptionsScreen(title, description, keyId, deploy.AlphaNumericInput,
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

	return deploy.ManageInputTextScreen(title, description, secret, deploy.AlphaNumericInput,
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
	var useTLS bool
	var continueProcessing bool

	if webClientModule {
		port, addresses, continueProcessing = readWebClientPortAddresses()
		if !continueProcessing {
			return false
		}

		useTLS, continueProcessing = readWebClientTLS()
		if !continueProcessing {
			return false
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

	var options []deploy.Option
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

			options = append(options, deploy.Option{
				Value:    fmt.Sprintf("%s", ip),
				Selected: true,
			})
		}
	}

	title := "Web Client Configurations"
	description := "Please inform the port number and select the IP addresses " +
		"that the Web Client will listen:"

	port, options, continueProcessing :=
		deploy.ManageInputTextOptionsScreen(title, description, port, deploy.NumericInput,
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

func readWebClientTLS() (useTLS, continueProcessing bool) {
	options := []deploy.Option{
		{Value: "Use TLS on interfaces (HTTPS)", Selected: true},
	}

	title := "Web Client Configurations"
	description := "Please select the following TLS options:"

	options, continueProcessing =
		deploy.ManageInputOptionsScreen(title, description, options, nil)

	if !continueProcessing {
		return
	}

	for index, option := range options {
		switch index {
		case 0:
			if option.Selected {
				useTLS = true
			}
		}
	}

	return
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
	return deploy.ManageInputTextScreen(title, description, host, deploy.HostnameOrIPInput,
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
	return deploy.ManageInputNumberScreen(title, description, port)
}

func readNotificationAuthType() (string, bool) {
	options := []deploy.Option{
		{Value: "None", Selected: false},
		{Value: "Plain", Selected: true},
		{Value: "CRAM-MD5", Selected: false},
	}

	title := "Notification Configurations"
	description := "Please inform the SMTP authentication type:"

	options, continueProcessing :=
		deploy.ManageInputOptionsScreen(title, description, options,
			func(options []deploy.Option, optionIndex int) []deploy.Option {
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
	return deploy.ManageInputTextScreen(title, description, host, deploy.AlphaNumericInput,
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
	return deploy.ManageInputTextScreen(title, description, host, deploy.AlphaNumericInput,
		func(input string) (bool, string) {
			if len(input) == 0 {
				return false, "SMTP password cannot be empty"
			}

			return true, ""
		})
}
