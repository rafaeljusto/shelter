// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/mail"
	"os"
	"time"

	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/mail/notification"
	"github.com/rafaeljusto/shelter/testing/utils"
)

var (
	configFilePath string // Path for the config file with the connection information
)

func init() {
	utils.TestName = "Notification"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for Notification test")
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

	messageChannel, errorChannel, err := utils.StartMailServer(config.ShelterConfig.Notification.SMTPServer.Port)
	if err != nil {
		utils.Fatalln("Error starting the mail server", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	templateName := createTemplateFile()
	defer removeTemplateFile(templateName)

	simpleNotification(domainDAO, templateName, messageChannel, errorChannel)

	utils.Println("SUCCESS!")
}

func createTemplateFile() string {
	f, err := ioutil.TempFile(".", "shelter-nf-test-template")
	if err != nil {
		utils.Fatalln("Error creating template file", err)
	}
	defer f.Close()

	_, err = f.WriteString(`{{$domain := .}}

From: {{.From}}
To: {{.To}}
Subject: {{normalizeEmailHeader (printf "%s %s" "Misconfiguration on domain" (fqdnToUnicode $domain.FQDN))}}


Dear Sir/Madam,

During our periodically domain verification, a configuration problem was detected with the
domain {{$domain.FQDN}}.

{{range $nameserver := $domain.Nameservers}}
  {{if nsStatusEq $nameserver.LastStatus "TIMEOUT"}}
  * Nameserver {{$nameserver.Host}} isn't answering the DNS requests.
    Please check your firewalls and DNS server and make sure that the service is up and
    the port 53 via UDP and TCP are allowed.

  {{else if nsStatusEq $nameserver.LastStatus "NOAA"}}
  * Nameserver {{$nameserver.Host}} don't have authority over the domain
    {{$domain.FQDN}}. Please check your nameserver configuration.

  {{else if nsStatusEq $nameserver.LastStatus "UDN"}}
  * Nameserver {{$nameserver.Host}} don't have data about the domain {{$domain.FQDN}}.

  {{else if nsStatusEq $nameserver.LastStatus "UH"}}
  * Nameserver {{$nameserver.Host}} couldn't be resolved. The authoritative DNS server
    could not be found.

  {{else if nsStatusEq $nameserver.LastStatus "SERVFAIL"}}
  * Nameserver {{$nameserver.Host}} got an internal error while receiving the DNS request.
    Please check the DNS server log to detect and solve the problem.

  {{else if nsStatusEq $nameserver.LastStatus "QREFUSED"}}
  * Nameserver {{$nameserver.Host}} refused to answer the DNS query. This is probably
    occuring because of an ACL. Authority nameservers cannot restrict requests for
    specific clients, please review the DNS server configuration.

  {{else if nsStatusEq $nameserver.LastStatus "CREFUSED"}}
  * Nameserver {{$nameserver.Host}} DNS query connection was refused. This is probably
    occuring because of firewall rule. Firewalls should allow port 53 in TCP and UDP
    protocols.

  {{else if nsStatusEq $nameserver.LastStatus "CNAME"}}
  * Nameserver {{$nameserver.Host}} has a CNAME in the zone APEX. According to RFC 1034 -
    section 3.6.2 and RFC 1912 - section 2.4 the CNAME record cannot exist with other
    resource record with the same name in the zone. As the SOA record is mandatory in the
    zone APEX, the CNAME cannot exist in it.

  {{else if nsStatusEq $nameserver.LastStatus "NOTSYNCH"}}
  * Nameserver {{$nameserver.Host}} is not synchronized with other nameservers of the
    domain {{$domain.FQDN}}. Check out the serial of the SOA records on each nameserver's zone.

  {{else if nsStatusEq $nameserver.LastStatus "ERROR"}}
  * Nameserver {{$nameserver.Host}} got an unexpected error.

  {{end}}
{{end}}

{{range $ds := $domain.DSSet}}
  {{if dsStatusEq $ds.LastStatus "TIMEOUT"}}
  * DS with keytag {{$ds.Keytag}} isn't answering the DNS requests.
    Please check your firewalls and DNS server and make sure that the service is up and
    the port 53 via UDP and TCP are allowed. Also, verify if your network supports
    fragmented UDP packagaes and UDP packages above 512 bytes (check EDSN0 for more
    information).

  {{else if dsStatusEq $ds.LastStatus "NOSIG"}}
  * DS with keytag {{$ds.Keytag}} references a DNSKEY record that don't have a RRSIG
    record (signature). Please sign the zone file with the DNSKEY record.

  {{else if dsStatusEq $ds.LastStatus "EXPSIG"}}
  * DS with keytag {{$ds.Keytag}} references a DNSKEY record with a expired signature.
    Please, resign the zone as soon as possible.

  {{else if dsStatusEq $ds.LastStatus "NOKEY"}}
  * DS with keytag {{$ds.Keytag}} references a DNSKEY record that does not exist in the
    zone

  {{else if dsStatusEq $ds.LastStatus "NOSEP"}}
  * DS with keytag {{$ds.Keytag}} references a DNSKEY that is not a security entry point.
    Some recursive DNS servers could invalidate the chain of trust for that reason.
    Please use a DNSKEY record with the bit SEP on.

  {{else if dsStatusEq $ds.LastStatus "SIGERR"}}
  * DS with keytag {{$ds.Keytag}} references a DNSKEY that have an invalid signature.
    Please resign your zone to fix this problem.

  {{else if dsStatusEq $ds.LastStatus "DNSERR"}}
  * DS with keytag {{$ds.Keytag}} could not be verified due to a problem on the
    nameservers.

  {{else if isNearExpiration $ds}}
  * DS with keytag {{$ds.Keytag}} references a DNSKEY with signatures that are near the
    expiration date. Please resign the zone before it expires to avoid DNS problems.

  {{end}}
{{end}}

Best regards,
LACTLD`)

	if err != nil {
		utils.Fatalln("Could not write to template file", err)
	}

	config.ShelterConfig.Languages = append(config.ShelterConfig.Languages, f.Name())
	return f.Name()
}

func removeTemplateFile(templateName string) {
	if err := os.Remove(templateName); err != nil {
		utils.Fatalln("Error removing template file", err)
	}
}

func simpleNotification(domainDAO dao.DomainDAO, templateName string,
	messageChannel chan *mail.Message, errorChannel chan error) {

	generateAndSaveDomain("example.com.br.", domainDAO, templateName)

	notification.TemplateExtension = ""
	if err := notification.LoadTemplates(); err != nil {
		utils.Fatalln("Error loading templates", err)
	}

	notification.Notify()

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- true
	}()

	select {
	case message := <-messageChannel:
		if message.Header.Get("From") != "shelter@example.com.br" {
			utils.Fatalln(fmt.Sprintf("E-mail from header is different. Expected "+
				"shelter@example.com.br but found %s", message.Header.Get("From")), nil)
		}

		if message.Header.Get("To") != "test@rafael.net.br" {
			utils.Fatalln("E-mail to header is different", nil)
		}

		if message.Header.Get("Subject") != "=?UTF-8?B?TWlzY29uZmlndXJhdGlvbiBvbiBkb21haW4gZXhhbXBsZS5jb20uYnIu?=" {
			utils.Fatalln("E-mail subject header is different", nil)
		}

		body, err := ioutil.ReadAll(message.Body)
		if err != nil {
			utils.Fatalln("Error reading e-mail body", err)
		}

		expectedBody := "Dear Sir/Madam,\r\n" +
			"\r\n" +
			"During our periodically domain verification, a configuration problem was detected with the\r\n" +
			"domain example.com.br..\r\n" +
			"\r\n" +
			"  * Nameserver ns1.example.com.br. got an internal error while receiving the DNS request.\r\n" +
			"    Please check the DNS server log to detect and solve the problem.\r\n" +
			"\r\n" +
			"Best regards,\r\n" +
			"LACTLD\r\n" +
			".\r\n"

		if string(body) != expectedBody {
			utils.Fatalln(fmt.Sprintf("E-mail body is different from what we expected. "+
				"Expected [%s], but found [%s]", expectedBody, body), nil)
		}

	case err := <-errorChannel:
		utils.Fatalln("Error receiving message", err)

	case <-timeout:
		utils.Fatalln("No mail sent", nil)
	}

	if err := domainDAO.RemoveByFQDN("example.com.br."); err != nil {
		utils.Fatalln("Error removing domain", err)
	}
}

// Function to mock a domain
func generateAndSaveDomain(fqdn string, domainDAO dao.DomainDAO, language string) {
	lastOKAt := time.Now().Add(time.Duration(-config.ShelterConfig.Notification.NameserverErrorAlertDays*24) * time.Hour)
	owner, _ := mail.ParseAddress("test@rafael.net.br")

	domain := model.Domain{
		FQDN: fqdn,
		Nameservers: []model.Nameserver{
			{
				Host:       fmt.Sprintf("ns1.%s", fqdn),
				IPv4:       net.ParseIP("127.0.0.1"),
				LastStatus: model.NameserverStatusServerFailure,
				LastOKAt:   lastOKAt,
			},
		},
		Owners: []model.Owner{
			{
				Email:    owner,
				Language: language,
			},
		},
	}

	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln(fmt.Sprintf("Fail to save domain %s", domain.FQDN), err)
	}
}
