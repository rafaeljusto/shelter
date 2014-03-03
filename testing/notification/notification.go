package main

import (
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/mail/notification"
	"github.com/rafaeljusto/shelter/testing/utils"
	"io/ioutil"
	"net"
	"net/mail"
	"time"
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

	messageChannel, errorChannel, err := utils.StartMailServer(config.ShelterConfig.Notification.SMTPServer.Port)
	if err != nil {
		utils.Fatalln("Error starting the mail server", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	simpleNotification(domainDAO, messageChannel, errorChannel)

	utils.Println("SUCCESS!")
}

func simpleNotification(domainDAO dao.DomainDAO,
	messageChannel chan *mail.Message, errorChannel chan error) {

	generateAndSaveDomain("example.com.br.", domainDAO)

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
		if message.Header.Get("From") != "shelter@example.com.br." {
			utils.Fatalln("E-mail from header is different", nil)
		}

		if message.Header.Get("To") != "[<test@rafael.net.br>]" {
			utils.Fatalln("E-mail to header is different", nil)
		}

		if message.Header.Get("Subject") != "Problema de configuracao com o dominio example.com.br." {
			utils.Fatalln("E-mail subject header is different", nil)
		}

		body, err := ioutil.ReadAll(message.Body)
		if err != nil {
			utils.Fatalln("Error reading e-mail body", err)
		}

		expectedBody := `
Prezado Sr./Sra.,

Durante a validação periódica de domínio, um problema de configuração foi detectado com o
domínio example.com.br..


  
  * Servidor DNS ns1.example.com.br. gerou um erro interno enquanto recebia a
    requisição DNS. Por favor verifique os logs para detectar e resolver o problema.

  




Atenciosamente,
LACTLD
.
`

		if string(body) != expectedBody {
			// TODO: Not validating for now
			// utils.Fatalln(fmt.Sprintf("E-mail body is different from what we expected. "+
			// 	"Expected [%s], but found [%s]", expectedBody, body), nil)
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
func generateAndSaveDomain(fqdn string, domainDAO dao.DomainDAO) {
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
				Language: "pt-BR",
			},
		},
	}

	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln(fmt.Sprintf("Fail to save domain %s", domain.FQDN), err)
	}
}
