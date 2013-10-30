package model

import (
	"net"
	"time"
)

// List of possible nameserver status
const (
	NameserverStatusOK = iota
	NameserverStatusTimeout
	NameserverStatusNoAuthority
	NameserverStatusUnknownDomainName
	NameserverStatusUnknownHost
	NameserverStatusFail
	NameserverStatusQueryRefused
	NameserverStatusConnectionRefused
	NameserverStatusCanonicalName
	NameserverStatusNotSynchronized
)

var NameserverStatus int

// Nameserver store the information necessary to send the requests for a specific host and
// store the results of this requests
type Nameserver struct {
	Host       string
	IPv4       net.IP
	IPv6       net.IP
	LastStatus NameserverStatus
	LastCheck  time.Time
	LastOK     time.Time
}
