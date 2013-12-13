package protocol

// Domain object from the protocol used to determinate what the user can update
type DomainRequest struct {
	FQDN        string              `json:"-"`                     // Actual domain name
	Nameservers []NameserverRequest `json:"nameservers,omitempty"` // Nameservers that asnwer with authority for this domain
	DSSet       []DSRequest         `json:"dsset,omitempty"`       // Records for the DNS tree chain of trust
	Owners      []string            `json:"owners,omitempty"`      // E-mails that will be alerted on any problem
}

// Domain object from the protocol used to determinate what the user can see. The last
// modified field is not here because it is sent in HTTP header field as it is with
// revision (ETag)
type DomainResponse struct {
	FQDN        string               `json:"-"`                     // Actual domain name
	Nameservers []NameserverResponse `json:"nameservers,omitempty"` // Nameservers that asnwer with authority for this domain
	DSSet       []DSResponse         `json:"dsset,omitempty"`       // Records for the DNS tree chain of trust
	Owners      []string             `json:"owners,omitempty"`      // E-mails that will be alerted on any problem
}
