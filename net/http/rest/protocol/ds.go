package protocol

import (
	"shelter/model"
)

// DS object used in the protocol to determinate what the user can update
type DSRequest struct {
	Keytag     uint16             `json:"keytag,omitempty"`     // DNSKEY's identification number
	Algorithm  model.DSAlgorithm  `json:"algorithm,omitempty"`  // DNSKEY's algorithm
	Digest     string             `json:"digest,omitempty"`     // Hash of the DNSKEY content
	DigestType model.DSDigestType `json:"digestType,omitempty"` // Hash type decided by user when generating the DS
}

// DS object used in the protocol to determinate what the user can see. The status was
// converted to text format for easy interpretation
type DSResponse struct {
	Keytag      uint16             `json:"keytag,omitempty"`      // DNSKEY's identification number
	Algorithm   model.DSAlgorithm  `json:"algorithm,omitempty"`   // DNSKEY's algorithm
	Digest      string             `json:"digest,omitempty"`      // Hash of the DNSKEY content
	DigestType  model.DSDigestType `json:"digestType,omitempty"`  // Hash type decided by user when generating the DS
	ExpiresAt   time.Time          `json:"expiresAt,omitempty"`   // DNSKEY's signature expiration date
	LastStatus  string             `json:"lastStatus,omitempty"`  // Result of the last configuration check
	LastCheckAt time.Time          `json:"lastCheckAt,omitempty"` // Time of the last configuration check
	LastOKAt    time.Time          `json:"lastOKAt,omitempty"`    // Last time that the DNSSEC configuration was OK
}
