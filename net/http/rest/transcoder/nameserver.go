package transcoder

import (
	"net"
	"shelter/model"
	"shelter/net/http/rest/protocol"
)

// Convert a list of nameserver requests object into a list of nameserver model objects.
// Useful when merging domain object from the network with a domain object from the
// database. It can return errors related to the conversion of IP addresses and
// normalization of nameserver's hostname
func toNameserversModel(nameserversRequest []protocol.NameserverRequest) ([]model.Nameserver, error) {
	var nameservers []model.Nameserver
	for _, nameserverRequest := range nameserversRequest {
		nameserver, err := toNameserverModel(nameserverRequest)
		if err != nil {
			return nil, err
		}

		nameservers = append(nameservers, nameserver)
	}

	return nameservers
}

// Convert a nameserver request object into a nameserver model object. It can return
// errors related to the conversion of IP addresses and normalization of nameserver's
// hostname
func toNameserverModel(nameserverRequest protocol.NameserverRequest) (model.Nameserver, error) {
	var nameserver model.Nameserver

	host, err := model.NormalizeDomainName(nameserverRequest.Host)
	if err != nil {
		return nameserver, err
	}
	nameserver.Host = host

	if len(nameserverRequest.IPv4) > 0 {
		ipv4, err := net.ParseIP(nameserverRequest.IPv4)
		if err != nil {
			return nameserver, err
		}
		nameserver.IPv4 = ipv4
	}

	if len(nameserverRequest.IPv6) > 0 {
		ipv6, err := net.ParseIP(nameserverRequest.IPv6)
		if err != nil {
			return nameserver, err
		}
		nameserver.IPv6 = ipv6
	}

	return nameserver, nil
}
