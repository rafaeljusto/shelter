package model

import (
	"net/mail"
)

// Domain stores all the necessary information for validating the DNS and DNSSEC. It also
// stores information to alert the domain's owners about the problems.
type Domain struct {
	Id          int            // Database identification
	FQDN        string         // Actual domain name
	Nameservers []Nameserver   // Nameservers that asnwer with authority for this domain
	DSSet       []DS           // Records for the DNS tree chain of trust
	Owners      []mail.Address // E-mails that will be alerted on any problem
}
