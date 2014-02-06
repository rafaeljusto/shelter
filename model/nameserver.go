package model

import (
	"net"
	"strings"
	"time"
)

// List of possible nameserver status
const (
	NameserverStatusNotChecked        = iota // Nameserver not checked yet
	NameserverStatusOK                       // DNS configuration for this nameserver is OK
	NameserverStatusTimeout                  // Network timeout while trying to reach the nameserver
	NameserverStatusNoAuthority              // Nameserver does not have authority for this domain
	NameserverStatusUnknownDomainName        // Domain does not exists for this nameserver
	NameserverStatusUnknownHost              // Could not resolve nameserver (no glue)
	NameserverStatusServerFailure            // Nameserver configuration problem
	NameserverStatusQueryRefused             // DNS request rejected
	NameserverStatusConnectionRefused        // Connection refused by firewall or nameserver
	NameserverStatusCanonicalName            // Domain name is a link in the zone APEX
	NameserverStatusNotSynchronized          // Nameservers of this domain have a different version of the zone files
	NameserverStatusError                    // Generic error found in the nameserver
)

// NameserverStatus is a number that represents one of the possible nameserver status
// listed in the constant group above
type NameserverStatus int

// Convert the nameserver status enum to text for printing in reports or debugging
func NameserverStatusToString(status NameserverStatus) string {
	switch status {
	case NameserverStatusNotChecked:
		return "NOTCHECKED"
	case NameserverStatusOK:
		return "OK"
	case NameserverStatusTimeout:
		return "TIMEOUT"
	case NameserverStatusNoAuthority:
		return "NOAA"
	case NameserverStatusUnknownDomainName:
		return "UDN"
	case NameserverStatusUnknownHost:
		return "UH"
	case NameserverStatusServerFailure:
		return "SERVFAIL"
	case NameserverStatusQueryRefused:
		return "QREFUSED"
	case NameserverStatusConnectionRefused:
		return "CREFUSED"
	case NameserverStatusCanonicalName:
		return "CNAME"
	case NameserverStatusNotSynchronized:
		return "NOTSYNCH"
	case NameserverStatusError:
		return "ERROR"
	}

	return ""
}

// Nameserver store the information necessary to send the requests for a specific host and
// store the results of this requests
type Nameserver struct {
	Host        string           // Nameserver's name
	IPv4        net.IP           // Host's IPv4 (optional when don't need glue)
	IPv6        net.IP           // Host's IPv6 (optional)
	LastStatus  NameserverStatus // Result of the last configuration check
	LastCheckAt time.Time        // Time of the last configuration check
	LastOKAt    time.Time        // Last time that the DNS configuration was OK
}

// Method to check if the nameserver needs glue for a given domain name. A namerserver
// needs glue when the name of the domain is inside the nameserver (example: domain
// test.com.br and nameserver ns1.tes.com.br)
func (n Nameserver) NeedsGlue(fqdn string) bool {
	return strings.HasSuffix(n.Host, fqdn) ||
		strings.HasSuffix(n.Host, fqdn+".") ||
		strings.HasSuffix(n.Host+".", fqdn)
}

// ChangeStatus is a easy way to change the status of a nameserver because it also updates
// the last check date
func (n *Nameserver) ChangeStatus(status NameserverStatus) {
	n.LastStatus = status
	n.LastCheckAt = time.Now()

	if status == NameserverStatusOK {
		n.LastOKAt = n.LastCheckAt
	}
}
