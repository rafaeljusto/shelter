package transcoder

import (
	"errors"
	"shelter/model"
	"shelter/net/http/rest/protocol"
)

// List of possible errors that can occur when calling methods from this object. Other
// erros can also occurs from low level layers
var (
	// Error returned when trying to merge to domain objects with different FQDNs
	ErrDomainsFQDNDontMatch = errors.New("Cannot merge domains with different FQDNs")
)

// Merge is used to merge a domain request object sent by the user into a domain object of
// the database. It can return errors related to merge problems that are problem caused by
// data format of the user input
func Merge(domain *model.Domain, domainRequest protocol.DomainRequest) (*model.Domain, error) {
	var err error
	if domainRequest.FQDN, err = model.NormalizeDomainName(domainRequest.FQDN); err != nil {
		return nil, err
	}

	if domain == nil {
		domain = new(model.Domain)
		domain.FQDN = domainRequest.FQDN

	} else {
		// Cannot merge domains with different FQDNs
		if domain.FQDN != domainRequest.FQDN {
			return nil, ErrDomainsFQDNDontMatch
		}
	}

	nameservers, err := toNameserversModel(domainRequest.Nameservers)
	if err != nil {
		return nil, err
	}

	for index, userNameserver := range nameservers {
		for _, nameserver := range domain.Nameservers {
			if nameserver.Host == userNameserver.Host {
				// Found the same nameserver in the user domain object, maybe the user updated the
				// IP addresses
				nameserver.IPv4 = otherNameserver.IPv4
				nameserver.IPv6 = otherNameserver.IPv6
				nameservers[index] = nameserver
				break
			}
		}
	}
	domain.Nameservers = nameservers

	// TODO: DSSet and emails
}

func ToDomainResponse(domain model.Domain) protocol.DomainResponse {
	// TODO
}
