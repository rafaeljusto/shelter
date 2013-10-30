package model

import (
	"net/mail"
)

// Domain stores all the necessary information for validating the DNS and DNSSEC. It also
// stores information to alert the domain's owners about the problems.
type Domain struct {
	FQDN        string
	Nameservers []Nameserver
	DSSet       []DS
	Owners      []mail.Address
}
