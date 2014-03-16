package protocol

import (
	"errors"
	"fmt"
	"github.com/rafaeljusto/shelter/model"
)

// List of possible errors that can occur when calling methods from this object. Other
// erros can also occurs from low level layers
var (
	// Error returned when trying to merge to domain objects with different FQDNs
	ErrDomainsFQDNDontMatch = errors.New("Cannot merge domains with different FQDNs")
)

// Domain object from the protocol used to determinate what the user can update
type DomainRequest struct {
	FQDN        string              `json:"-"`                     // Actual domain name
	Nameservers []NameserverRequest `json:"nameservers,omitempty"` // Nameservers that asnwer with authority for this domain
	DSSet       []DSRequest         `json:"dsset,omitempty"`       // Records for the DNS tree chain of trust
	DNSKEYS     []DNSKEYRequest     `json:"dnskeys,omitempty"`     // Records that can be converted into DS records
	Owners      []OwnerRequest      `json:"owners,omitempty"`      // E-mails that will be alerted on any problem
}

// Merge is used to merge a domain request object sent by the user into a domain object of
// the database. It can return errors related to merge problems that are problem caused by
// data format of the user input
func Merge(domain model.Domain, domainRequest DomainRequest) (model.Domain, error) {
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

	dnskeysDSSet, err := dnskeysRequestsToDSSetModel(domain.FQDN, domainRequest.DNSKEYS)
	if err != nil {
		return domain, err
	}
	dsSet = append(dsSet, dnskeysDSSet...)

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
	domain.Owners, err = toOwnersModel(domainRequest.Owners)
	if err != nil {
		return domain, err
	}

	return domain, nil
}

// Domain object from the protocol used to determinate what the user can see. The last
// modified field is not here because it is sent in HTTP header field as it is with
// revision (ETag)
type DomainResponse struct {
	FQDN        string               `json:"fqdn"`                  // Actual domain name
	Nameservers []NameserverResponse `json:"nameservers,omitempty"` // Nameservers that asnwer with authority for this domain
	DSSet       []DSResponse         `json:"dsset,omitempty"`       // Records for the DNS tree chain of trust
	Owners      []OwnerResponse      `json:"owners,omitempty"`      // E-mails that will be alerted on any problem
	Links       []Link               `json:"links,omitempty"`       // Links to manipulate object
}

// Convert the domain system object to a limited information user format. We have a persisted flag
// to known when the object exists in our database or not to choose when we need to add the object
// links or not
func ToDomainResponse(domain model.Domain, persisted bool) DomainResponse {
	var links []Link

	// We don't add links when the object doesn't exist in the system yet
	if persisted {
		// We should add more links here for system navigation. For example, we could add links
		// for object update, delete, list, etc. But I did not found yet in IANA list the
		// correct link type to be used. Also, the URI is hard coded, I didn't have any idea on
		// how can we do this dynamically yet. We cannot get the URI from the handler because we
		// are going to have a cross-reference problem
		links = append(links, Link{
			Types: []LinkType{LinkTypeSelf},
			HRef:  fmt.Sprintf("/domain/%s", domain.FQDN),
		})
	}

	return DomainResponse{
		FQDN:        domain.FQDN,
		Nameservers: toNameserversResponse(domain.Nameservers),
		DSSet:       toDSSetResponse(domain.DSSet),
		Owners:      toOwnersResponse(domain.Owners),
		Links:       links,
	}
}
