package notification

import (
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"net/smtp"
)

func Notify() {
	defer func() {
		// Something went really wrong while notifying the owners. Log the error stacktrace and move out
		if r := recover(); r != nil {
			log.Println("Panic detected while notifying the owners. Details:", r)
		}
	}()

	database, databaseSession, err := mongodb.Open(
		config.ShelterConfig.Database.URI,
		config.ShelterConfig.Database.Name,
	)

	if err != nil {
		log.Println("Error while initializing database. Details:", err)
		return
	}
	defer databaseSession.Close()

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	domains, err := domainDAO.FindAllToBeNotified(
		config.ShelterConfig.Notification.NameserverErrorAlertDays,
		config.ShelterConfig.Notification.NameserverTimeoutAlertDays,
		config.ShelterConfig.Notification.DSErrorAlertDays,
		config.ShelterConfig.Notification.DSTimeoutAlertDays,
	)

	if err != nil {
		log.Println("Error retrieving domain to notify. Details:", err)
		return
	}

	for _, domain := range domains {
		// TODO: Is better do this async
		notifyDomain(domain)
	}
}

func notifyDomain(domain model.Domain) {
	from := config.ShelterConfig.Notification.From

	var emails []string
	for _, email := range domain.Owners {
		emails = append(emails, email.String())
	}

	server := fmt.Sprintf("%s:%d",
		config.ShelterConfig.Notification.SMTPServer.Server,
		config.ShelterConfig.Notification.SMTPServer.Port,
	)

	// TODO: Build template message
	msg := []byte{}

	switch config.ShelterConfig.Notification.SMTPServer.Auth.Type {
	case config.AuthenticationTypePlain:
		auth := smtp.PlainAuth("",
			config.ShelterConfig.Notification.SMTPServer.Auth.Username,
			config.ShelterConfig.Notification.SMTPServer.Auth.Password,
			config.ShelterConfig.Notification.SMTPServer.Server,
		)
		smtp.SendMail(server, auth, from, emails, msg)

	case config.AuthenticationTypeCRAMMD5Auth:
		auth := smtp.CRAMMD5Auth(
			config.ShelterConfig.Notification.SMTPServer.Auth.Username,
			config.ShelterConfig.Notification.SMTPServer.Auth.Password,
		)
		smtp.SendMail(server, auth, from, emails, msg)

	default:
		smtp.SendMail(server, nil, from, emails, msg)
	}
}
