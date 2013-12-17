package protocol

import (
	"errors"
	"net"
	"shelter/model"
	"time"
)

// List of possible errors that can occur when calling methods from this object. Other
// erros can also occurs from low level layers
var (
	// Error returned when trying to convert an invalid IP address
	ErrInvalidIP = errors.New("IP address is not in a valid format")
)

// Nameserver object used in the protocol to determinate what the user can update
type NameserverRequest struct {
	Host string `json:"host,omitempty"` // Nameserver's name
	IPv4 string `json:"ipv4,omitempty"` // Host's IPv4 (optional when don't need glue)
	IPv6 string `json:"ipv6,omitempty"` // Host's IPv6 (optional)
}

// Namerserver object used in the protocol to determinate what the user can see. The
// status was converted to text format for easy interpretation
type NameserverResponse struct {
	Host        string    `json:"host,omitempty"`        // Nameserver's name
	IPv4        string    `json:"ipv4,omitempty"`        // Host's IPv4 (optional when don't need glue)
	IPv6        string    `json:"ipv6,omitempty"`        // Host's IPv6 (optional)
	LastStatus  string    `json:"lastStatus,omitempty"`  // Result of the last configuration check
	LastCheckAt time.Time `json:"lastCheckAt,omitempty"` // Time of the last configuration check
	LastOKAt    time.Time `json:"lastOKAt,omitempty"`    // Last time that the DNS configuration was OK
}

// Convert a list of nameserver requests objects into a list of nameserver model objects.
// Useful when merging domain object from the network with a domain object from the
// database. It can return errors related to the conversion of IP addresses and
// normalization of nameserver's hostname
func toNameserversModel(nameserversRequest []NameserverRequest) ([]model.Nameserver, error) {
	var nameservers []model.Nameserver
	for _, nameserverRequest := range nameserversRequest {
		nameserver, err := toNameserverModel(nameserverRequest)
		if err != nil {
			return nil, err
		}

		nameservers = append(nameservers, nameserver)
	}

	return nameservers, nil
}

// Convert a nameserver request object into a nameserver model object. It can return
// errors related to the conversion of IP addresses and normalization of nameserver's
// hostname
func toNameserverModel(nameserverRequest NameserverRequest) (model.Nameserver, error) {
	var nameserver model.Nameserver

	host, err := model.NormalizeDomainName(nameserverRequest.Host)
	if err != nil {
		return nameserver, err
	}
	nameserver.Host = host

	if len(nameserverRequest.IPv4) > 0 {
		ipv4 := net.ParseIP(nameserverRequest.IPv4)
		if ipv4 == nil {
			return nameserver, ErrInvalidIP
		}
		nameserver.IPv4 = ipv4
	}

	if len(nameserverRequest.IPv6) > 0 {
		ipv6 := net.ParseIP(nameserverRequest.IPv6)
		if ipv6 == nil {
			return nameserver, ErrInvalidIP
		}
		nameserver.IPv6 = ipv6
	}

	return nameserver, nil
}
