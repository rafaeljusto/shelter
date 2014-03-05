package notification

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/log"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/mail/notification/protocol"
	"net/smtp"
	"regexp"
)

// List of possible errors that can occur when calling functions from this file. Other
// erros can also occurs from low level layers
var (
	// Template needed to send an e-mail not found
	ErrTemplateNotFound = errors.New("Template not found")
)

var (
	// Regular expression to detect two or more line breaks. This is necessary because after the
	// template execution, the template controllers are removed, but the line breaks that come with
	// them are not. For that reason we remove them using the strategy to transform two or more
	// consecutivilly line breaks into one. If the user really wants to put a white line, he can add a
	// space to the blank line.
	extraSpaces = regexp.MustCompile("(\n){2,}")
)

// Notify is responsable for selecting the domains that should be notified in the system.
// It will send alert e-mails for each owner of a domain
func Notify() {
	defer func() {
		// Something went really wrong while notifying the owners. Log the error stacktrace
		// and move out
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
			append(emailsPerLanguage[owner.Language], owner.Email.String())
	}

	server := fmt.Sprintf("%s:%d",
		config.ShelterConfig.Notification.SMTPServer.Server,
		config.ShelterConfig.Notification.SMTPServer.Port,
	)

	for language, emails := range emailsPerLanguage {
		t := getTemplate(language)
		if t == nil {
			return ErrTemplateNotFound
		}

		domainMail := protocol.Domain{
			Domain: *domain,
			From:   config.ShelterConfig.Notification.From,
			To:     emails,
		}

		var msg bytes.Buffer
		if err := t.ExecuteTemplate(&msg, "notification", domainMail); err != nil {
			return err
		}

		// Remove extra new lines that can appear because of the template execution. Special lines used
		// for controlling the templates are removed but the new lines are left behind. We don't remove
		// the first one, because is the division beteween the e-mail header and body
		firstOcurrence := true
		msgBytes := extraSpaces.ReplaceAllFunc(msg.Bytes(), func(match []byte) []byte {
			if firstOcurrence {
				firstOcurrence = false
				return match
			}

			return []byte("\n")
		})

		switch config.ShelterConfig.Notification.SMTPServer.Auth.Type {
		case config.AuthenticationTypePlain:
			auth := smtp.PlainAuth("",
				config.ShelterConfig.Notification.SMTPServer.Auth.Username,
				config.ShelterConfig.Notification.SMTPServer.Auth.Password,
				config.ShelterConfig.Notification.SMTPServer.Server,
			)

			if err := smtp.SendMail(server, auth, from, emails, msgBytes); err != nil {
				return err
			}

		case config.AuthenticationTypeCRAMMD5Auth:
			auth := smtp.CRAMMD5Auth(
				config.ShelterConfig.Notification.SMTPServer.Auth.Username,
				config.ShelterConfig.Notification.SMTPServer.Auth.Password,
			)

			if err := smtp.SendMail(server, auth, from, emails, msgBytes); err != nil {
				return err
			}

		default:
			if err := smtp.SendMail(server, nil, from, emails, msgBytes); err != nil {
				return err
			}
		}
	}

	return nil
}
