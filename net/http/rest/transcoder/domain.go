package transcoder

import (
	"errors"
	"net/mail"
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
func Merge(domain model.Domain, domainRequest protocol.DomainRequest) (model.Domain, error) {
	var err error
	if domainRequest.FQDN, err = model.NormalizeDomainName(domainRequest.FQDN); err != nil {
		return domain, err
	}

	// Detect when the domain object is empty, that is the case when we are creating a new
	// domain in the Shelter project
	if len(domain.FQDN) == 0 {
		domain.FQDN = domainRequest.FQDN

	} else {
		// Cannot merge domains with different FQDNs
		if domain.FQDN != domainRequest.FQDN {
			return domain, ErrDomainsFQDNDontMatch
		}
	}

	nameservers, err := toNameserversModel(domainRequest.Nameservers)
	if err != nil {
		return domain, err
	}

	for index, userNameserver := range nameservers {
		for _, nameserver := range domain.Nameservers {
			if nameserver.Host == userNameserver.Host {
				// Found the same nameserver in the user domain object, maybe the user updated the
				// IP addresses
				nameserver.IPv4 = userNameserver.IPv4
				nameserver.IPv6 = userNameserver.IPv6
				nameservers[index] = nameserver
				break
			}
		}
	}
	domain.Nameservers = nameservers

	dsSet, err := toDSSetModel(domainRequest.DSSet)
	if err != nil {
		return domain, err
	}

	for index, userDS := range dsSet {
		for _, ds := range domain.DSSet {
			if ds.Keytag == userDS.Keytag {
				// Found the same DS in the user domain object
				ds.Algorithm = userDS.Algorithm
				ds.Digest = userDS.Digest
				ds.DigestType = userDS.DigestType
				dsSet[index] = ds
				break
			}
		}
	}
	domain.DSSet = dsSet

	// We can replace the whole structure of the e-mail every time that a new UPDATE arrives
	// because there's no extra information in server side that we need to keep
	var owners []*mail.Address
	for _, owner := range domainRequest.Owners {
		email, err := mail.ParseAddress(owner)
		if err != nil {
			return domain, err
		}

		owners = append(owners, email)
	}
	domain.Owners = owners

	return domain, nil
}

func ToDomainResponse(domain model.Domain) protocol.DomainResponse {
	// TODO

	return protocol.DomainResponse{}
}
