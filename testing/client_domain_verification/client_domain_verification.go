// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/protocol"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	configFilePath string // Path for the config file with the connection information
)

// ClientDomainTestConfigFile is a structure to store the test configuration file data
type ClientDomainVerificationTestConfigFile struct {
	config.Config
	DNSServerPort int
}

func init() {
	utils.TestName = "ClientDomainVerification"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ClientDomainVerification test")
}

func main() {
	flag.Parse()

	var clientConfig ClientDomainVerificationTestConfigFile
	err := utils.ReadConfigFile(configFilePath, &clientConfig)
	config.ShelterConfig = clientConfig.Config

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	utils.StartDNSServer(clientConfig.DNSServerPort, clientConfig.Scan.UDPMaxSize)
	finishRESTServer := utils.StartRESTServer()
	defer finishRESTServer()
	utils.StartWebClient()

	scanDomain()
	queryDomain()

	utils.Println("SUCCESS!")
}

func scanDomain() {
	dns.HandleFunc("example.com.br.", func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		dnsResponseMessage := &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Authoritative: true,
			},
			Question: dnsRequestMessage.Question,
			Answer: []dns.RR{
				&dns.SOA{
					Hdr: dns.RR_Header{
						Name:   "example.com.br.",
						Rrtype: dns.TypeSOA,
						Class:  dns.ClassINET,
						Ttl:    86400,
					},
					Ns:      "ns1.example.com.br.",
					Mbox:    "rafael.justo.net.br.",
					Serial:  2013112600,
					Refresh: 86400,
					Retry:   86400,
					Expire:  86400,
					Minttl:  900,
				},
			},
		}

		dnsResponseMessage.SetReply(dnsRequestMessage)
		w.WriteMsg(dnsResponseMessage)
	})

	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	content := `{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns2.example.com.br.", "ipv4": "127.0.0.1" }
      ]
    }`

	r, err := http.NewRequest("PUT",
		fmt.Sprintf("%s%s", url, "/domain/example.com.br./verification"),
		strings.NewReader(content))

	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, []byte(content))

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Fatalln("Error reading response content", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln("Error scanning domain", errors.New(string(responseContent)))
	}

	var domainResponse protocol.DomainResponse
	if err := json.Unmarshal(responseContent, &domainResponse); err != nil {
		utils.Fatalln("Error decoding domain response", err)
	}

	if len(domainResponse.Nameservers) != 2 {
		utils.Fatalln("Wrong number of nameservers", nil)
	}

	if domainResponse.Nameservers[0].LastStatus != model.NameserverStatusToString(model.NameserverStatusOK) {
		utils.Fatalln("Scan did not work for ns1", nil)
	}

	if domainResponse.Nameservers[1].LastStatus != model.NameserverStatusToString(model.NameserverStatusOK) {
		utils.Fatalln("Scan did not work for ns2", nil)
	}
}

func queryDomain() {
	dns.HandleFunc("example.com.br.", func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		if dnsRequestMessage.Question[0].Qtype == dns.TypeNS {
			dnsResponseMessage := &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Authoritative: true,
				},
				Question: dnsRequestMessage.Question,
				Answer: []dns.RR{
					&dns.NS{
						Hdr: dns.RR_Header{
							Name:   "example.com.br.",
							Rrtype: dns.TypeNS,
							Class:  dns.ClassINET,
							Ttl:    86400,
						},
						Ns: "a.dns.br.",
					},
					&dns.NS{
						Hdr: dns.RR_Header{
							Name:   "example.com.br.",
							Rrtype: dns.TypeNS,
							Class:  dns.ClassINET,
							Ttl:    86400,
						},
						Ns: "b.dns.br.",
					},
				},
			}

			dnsResponseMessage.SetReply(dnsRequestMessage)
			w.WriteMsg(dnsResponseMessage)

		} else if dnsRequestMessage.Question[0].Qtype == dns.TypeDNSKEY {
			// Empty response
			dnsResponseMessage := &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Authoritative: true,
				},
				Question: dnsRequestMessage.Question,
				Answer:   []dns.RR{},
			}

			dnsResponseMessage.SetReply(dnsRequestMessage)
			w.WriteMsg(dnsResponseMessage)
		}
	})

	var client http.Client

	url := ""
	if len(config.ShelterConfig.WebClient.Listeners) > 0 {
		url = fmt.Sprintf("http://%s:%d", config.ShelterConfig.WebClient.Listeners[0].IP,
			config.ShelterConfig.WebClient.Listeners[0].Port)
	}

	if len(url) == 0 {
		utils.Fatalln("There's no interface to connect to", nil)
	}

	r, err := http.NewRequest("GET",
		fmt.Sprintf("%s%s", url, "/domain/example.com.br./verification"), nil)

	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}

	utils.BuildHTTPHeader(r, nil)

	response, err := client.Do(r)
	if err != nil {
		utils.Fatalln("Error sending request", err)
	}

	responseContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Fatalln("Error reading response content", err)
	}

	if response.StatusCode != http.StatusOK {
		utils.Fatalln("Error querying domain", errors.New(string(responseContent)))
	}

	var domainResponse protocol.DomainResponse
	if err := json.Unmarshal(responseContent, &domainResponse); err != nil {
		utils.Fatalln("Error decoding domain response", err)
	}

	if len(domainResponse.Nameservers) != 2 {
		utils.Fatalln("Wrong number of nameservers", nil)
	}

	if len(domainResponse.DSSet) != 0 {
		utils.Fatalln("Wrong number of DS records", nil)
	}
}
