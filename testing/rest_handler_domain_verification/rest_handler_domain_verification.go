// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/rafaeljusto/handy"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
)

var (
	configFilePath string // Path for the config file with the connection information
)

// RESTHandlerDomainTestConfigFile is a structure to store the test configuration file data
type RESTHandlerDomainVerificationTestConfigFile struct {
	config.Config
	DNSServerPort int
}

func init() {
	utils.TestName = "RESTHandlerDomainVerification"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for RESTHandlerDomainVerification test")
}

func main() {
	flag.Parse()

	var restConfig RESTHandlerDomainVerificationTestConfigFile
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
		restConfig.Database.URIs,
		restConfig.Database.Name,
		false, "", "",
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

	utils.StartDNSServer(restConfig.DNSServerPort, restConfig.Scan.UDPMaxSize)

	scanDomain()
	scanPersistedDomain(domainDAO)
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

	mux := handy.NewHandy()

	h := new(handler.DomainVerificationHandler)
	mux.Handle("/domain/{fqdn}/verification", func() handy.Handler {
		return h
	})

	requestContent := `{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns2.example.com.br.", "ipv4": "127.0.0.1" }
      ]
    }`

	r, err := http.NewRequest("PUT", "/domain/example.com.br./verification",
		strings.NewReader(requestContent))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, []byte(requestContent))

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	responseContent, err := ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln(fmt.Sprintf("Error scanning domain. "+
			"Expected %d and got %d", http.StatusOK, w.Code),
			errors.New(string(responseContent)))
	}

	if len(h.Response.Nameservers) != 2 {
		utils.Fatalln("Wrong number of nameservers", nil)
	}

	if h.Response.Nameservers[0].LastStatus !=
		model.NameserverStatusToString(model.NameserverStatusOK) {

		utils.Fatalln("Scan did not work for ns1", nil)
	}

	if h.Response.Nameservers[1].LastStatus !=
		model.NameserverStatusToString(model.NameserverStatusOK) {

		utils.Fatalln("Scan did not work for ns2", nil)
	}
}

func scanPersistedDomain(domainDAO dao.DomainDAO) {
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

	mux := handy.NewHandy()

	h := new(handler.DomainVerificationHandler)
	mux.Handle("/domain/{fqdn}/verification", func() handy.Handler {
		return h
	})

	requestContent := `{
      "Nameservers": [
        { "Host": "ns1.example.com.br.", "ipv4": "127.0.0.1" },
        { "Host": "ns2.example.com.br.", "ipv4": "127.0.0.1" }
      ]
    }`

	domain := model.Domain{
		FQDN: "example.com.br.",
		Nameservers: []model.Nameserver{
			{
				Host: "ns1.example.com.br.",
				IPv4: net.ParseIP("127.0.0.1"),
			},
			{
				Host: "ns2.example.com.br.",
				IPv4: net.ParseIP("127.0.0.1"),
			},
		},
	}

	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln("Error saving domain", err)
	}

	r, err := http.NewRequest("PUT", "/domain/example.com.br./verification",
		strings.NewReader(requestContent))
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, []byte(requestContent))

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	responseContent, err := ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln(fmt.Sprintf("Error scanning domain. "+
			"Expected %d and got %d", http.StatusOK, w.Code),
			errors.New(string(responseContent)))
	}

	domain, err = domainDAO.FindByFQDN(domain.FQDN)
	if err != nil {
		utils.Fatalln("Error retrieving the domain", err)
	}

	if len(domain.Nameservers) != 2 ||
		domain.Nameservers[0].LastStatus != model.NameserverStatusOK ||
		domain.Nameservers[1].LastStatus != model.NameserverStatusOK {

		utils.Fatalln("Not updating domain on a scan", nil)
	}

	// Now we are going to test when the database domain is different from the scanned domain

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		utils.Fatalln("Error removing domain", err)
	}

	domain = model.Domain{
		FQDN: "example.com.br.",
		Nameservers: []model.Nameserver{
			{
				Host: "ns1.example.com.br.",
				IPv4: net.ParseIP("127.0.0.1"),
			},
			{
				Host: "ns3.example.com.br.",
				IPv4: net.ParseIP("127.0.0.1"),
			},
		},
	}

	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln("Error saving domain", err)
	}

	mux.ServeHTTP(w, r)

	responseContent, err = ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln(fmt.Sprintf("Error scanning domain. "+
			"Expected %d and got %d", http.StatusOK, w.Code),
			errors.New(string(responseContent)))
	}

	domain, err = domainDAO.FindByFQDN(domain.FQDN)
	if err != nil {
		utils.Fatalln("Error retrieving the domain", err)
	}

	if len(domain.Nameservers) != 2 ||
		domain.Nameservers[0].LastStatus != model.NameserverStatusNotChecked ||
		domain.Nameservers[1].LastStatus != model.NameserverStatusNotChecked {

		utils.Fatalln("Updating domain on a scan when it shouldn't", nil)
	}
}

func queryDomain() {
	dns.HandleFunc("example.com.br.", func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		dnsResponseMessage := &dns.Msg{
			MsgHdr:   dns.MsgHdr{},
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
	})

	mux := handy.NewHandy()

	h := new(handler.DomainVerificationHandler)
	mux.Handle("/domain/{fqdn}/verification", func() handy.Handler {
		return h
	})

	r, err := http.NewRequest("GET", "/domain/example.com.br./verification", nil)
	if err != nil {
		utils.Fatalln("Error creating the HTTP request", err)
	}
	utils.BuildHTTPHeader(r, nil)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	responseContent, err := ioutil.ReadAll(w.Body)
	if err != nil {
		utils.Fatalln("Error reading response body", err)
	}

	if w.Code != http.StatusOK {
		utils.Fatalln(fmt.Sprintf("Error scanning domain. "+
			"Expected %d and got %d", http.StatusOK, w.Code),
			errors.New(string(responseContent)))
	}

	if len(h.Response.Nameservers) != 2 {
		utils.Fatalln("Wrong number of nameservers", nil)
	}
}
