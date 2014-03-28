// utils - Features for make the test life easier
//
// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package utils

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/net/http/client"
	"github.com/rafaeljusto/shelter/net/http/rest"
	"github.com/rafaeljusto/shelter/net/http/rest/messages"
	"github.com/rafaeljusto/shelter/net/scan"
	"os"
	"path/filepath"
	"time"
)

func StartDNSServer(port int, udpMaxSize uint16) (server *dns.Server) {
	// Change the querier DNS port for the scan
	scan.DNSPort = port

	server = &dns.Server{
		Net:     "udp",
		Addr:    fmt.Sprintf("localhost:%d", port),
		UDPSize: int(udpMaxSize),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			Fatalln("Error starting DNS test server", err)
		}
	}()

	// Wait the DNS server to start before testing
	time.Sleep(1 * time.Second)
	return
}

func StartWebClient() {
	listeners, err := client.Listen()
	if err != nil {
		Fatalln("Error listening to interfaces", err)
	}

	if err := client.Start(listeners); err != nil {
		Fatalln("Error starting the WEB client", err)
	}

	// Wait the REST server to start before testing
	time.Sleep(1 * time.Second)
}

func StartRESTServer() func() {
	createMessagesFile()

	listeners, err := rest.Listen()
	if err != nil {
		Fatalln("Error listening to interfaces", err)
	}

	if err := rest.Start(listeners); err != nil {
		Fatalln("Error starting the REST server", err)
	}

	// Wait the REST server to start before testing
	time.Sleep(1 * time.Second)
	return removeMessagesFile
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
		Fatalln("Error creating messages file", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(languagePacks); err != nil {
		Fatalln("Error encoding messages structure", err)
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
