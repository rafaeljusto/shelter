// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package notification is the notification service
package notification

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/mail/notification/protocol"
	"github.com/rafaeljusto/shelter/secret"
)

// List of possible errors that can occur when calling functions from this file. Other
// erros can also occurs from low level layers
var (
	// Template needed to send an e-mail not found
	ErrTemplateNotFound = errors.New("Template not found")
)

var (
	// Regular expression to detect two or more line breaks. This is necessary because after
	// the template execution, the template controllers are removed, but the line breaks
	// that come with them are not. For that reason we remove them using the strategy to
	// transform three or more consecutivilly line breaks into one. We ignore white spaces
	// beteween the line breaks
	extraSpaces = regexp.MustCompile("(( )*\n){3,}")
)

// Notify is responsable for selecting the domains that should be notified in the system.
// It will send alert e-mails for each owner of a domain
func Notify() {
	defer func() {
		// Something went really wrong while notifying the owners. Log the error stacktrace
		// and move out
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("Panic detected while notifying the owners. Details: %v\n%s", r, buf)
		}
	}()

	log.Info("Start notification job")
	defer func() {
		log.Info("End notification job")
	}()

	log.Debugf("Initializing database with the parameters: URIS - %v | Name - %s | Auth - %t | Username - %s",
		config.ShelterConfig.Database.URIs,
		config.ShelterConfig.Database.Name,
		config.ShelterConfig.Database.Auth.Enabled,
		config.ShelterConfig.Database.Auth.Username,
	)

	database, databaseSession, err := mongodb.Open(
		config.ShelterConfig.Database.URIs,
		config.ShelterConfig.Database.Name,
		config.ShelterConfig.Database.Auth.Enabled,
		config.ShelterConfig.Database.Auth.Username,
		config.ShelterConfig.Database.Auth.Password,
	)

	if err != nil {
		log.Println("Error while initializing database. Details:", err)
		return
	}
	defer databaseSession.Close()

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	domainChannel, err := domainDAO.FindAllAsyncToBeNotified(
		config.ShelterConfig.Notification.NameserverErrorAlertDays,
		config.ShelterConfig.Notification.NameserverTimeoutAlertDays,
		config.ShelterConfig.Notification.DSErrorAlertDays,
		config.ShelterConfig.Notification.DSTimeoutAlertDays,

		// TODO: Should we move this configuration parameter to a place were both modules can
		// access it. This sounds better for configuration deployment
		config.ShelterConfig.Scan.VerificationIntervals.MaxExpirationAlertDays,
	)

	if err != nil {
		log.Println("Error retrieving domains to notify. Details:", err)
		return
	}

	// Dispatch the asynchronous part of the method
	for {
		// Get domain from the database (one-by-one)
		domainResult := <-domainChannel

		// Detect errors while retrieving a specific domain. We are not going to stop all the
		// process when only one domain got an error
		if domainResult.Error != nil {
			log.Println("Error retrieving domain to notify. Details:", domainResult.Error)
			continue
		}

		// Problem detected while retrieving a domain or we don't have domains anymore
		if domainResult.Error != nil || domainResult.Domain == nil {
			break
		}

		if err := notifyDomain(domainResult.Domain); err != nil {
			log.Println("Error notifying a domain. Details:", err)
		}
	}
}

// Function used to notify a single domain. It can return error if there's a problem while
// filling the template or sending the e-mail
func notifyDomain(domain *model.Domain) error {
	from := config.ShelterConfig.Notification.From

	emailsPerLanguage := make(map[string][]string)
	for _, owner := range domain.Owners {
		emailsPerLanguage[owner.Language] =
			append(emailsPerLanguage[owner.Language], owner.Email.Address)
	}

	if len(emailsPerLanguage) == 0 {
		log.Infof("There's no owner to notify domain %s", domain.FQDN)
	}

	server := fmt.Sprintf("%s:%d",
		config.ShelterConfig.Notification.SMTPServer.Server,
		config.ShelterConfig.Notification.SMTPServer.Port,
	)

	password := config.ShelterConfig.Notification.SMTPServer.Auth.Password
	if len(password) > 0 {
		var err error
		password, err = secret.Decrypt(password)
		if err != nil {
			return err
		}
	}

	for language, emails := range emailsPerLanguage {
		t := getTemplate(language)
		if t == nil {
			return ErrTemplateNotFound
		}

		domainMail := protocol.Domain{
			Domain: *domain,
			From:   config.ShelterConfig.Notification.From,
			To:     strings.Join(emails, ","),
			Date:   FormatDate(time.Now()),
		}

		var msg bytes.Buffer
		if err := t.ExecuteTemplate(&msg, "notification", domainMail); err != nil {
			return err
		}

		// Remove extra new lines that can appear because of the template execution. Special
		// lines used for controlling the templates are removed but the new lines are left
		// behind
		msgBytes := bytes.TrimSpace(msg.Bytes())
		msgBytes = extraSpaces.ReplaceAll(msgBytes, []byte("\n\n"))

		switch config.ShelterConfig.Notification.SMTPServer.Auth.Type {
		case config.AuthenticationTypePlain:
			log.Debugf("Sending notification for domain %s to %v via server %s with plain authentication",
				domain.FQDN, emails, server)

			auth := smtp.PlainAuth("",
				config.ShelterConfig.Notification.SMTPServer.Auth.Username,
				password,
				config.ShelterConfig.Notification.SMTPServer.Server,
			)

			if err := smtp.SendMail(server, auth, from, emails, msgBytes); err != nil {
				return err
			}

		case config.AuthenticationTypeCRAMMD5Auth:
			log.Debugf("Sending notification for domain %s to %v via server %s with CRAM MD5 authentication",
				domain.FQDN, emails, server)

			auth := smtp.CRAMMD5Auth(
				config.ShelterConfig.Notification.SMTPServer.Auth.Username,
				password,
			)

			if err := smtp.SendMail(server, auth, from, emails, msgBytes); err != nil {
				return err
			}

		default:
			log.Debugf("Sending notification for domain %s to %v via server %s without authentication",
				domain.FQDN, emails, server)

			if err := smtp.SendMail(server, nil, from, emails, msgBytes); err != nil {
				return err
			}
		}
	}

	return nil
}

// FormatDate returns a compliant RFC5322 datetime
func FormatDate(datetime time.Time) string {
    return date.Format(datetime.RFC1123Z)
}

