package model

import (
	"time"
)

// List of possible DS status
const (
	DSStatusOK               = iota // DNSSEC configuration for this DS is OK
	DSStatusTimeout                 // Network timeout while trying to retrieve the DNSKEY
	DSStatusNoSignature             // No RRSIG records found for the related DNSKEY
	DSStatusExpiredSignature        // At least one RRSIG record was expired
	DSStatusNoKey                   // No DNSKEY was found with the keytag of the DS
	DSStatusNoSEP                   // DNSKEY related to DS does not have the bit SEP on
	DSStatusSignatureError          // Error while checking DNSKEY signatures
	DSStatusNoDNSSEC                // Domain is not configured with DNSSEC
	DSStatusDNSError                // DNS error (check nameserver status)
)

// DSStatus is a number that represents one of the possible DS status listed in the
// constant group above
type DSStatus int

// DS store the information necessary to validate if a domain is configured correctly with
// DNSSEC, and it also stores the results of the validations. When the hosts have multiple
// DNSSEC problems, the worst problem (using a priority algorithm) will be stored in the
// DS
type DS struct {
	Id        int       // Database identification
	Keytag    uint16    // DNSKEY's identification number
	Algorithm uint8     // DNSKEY's algorithm
	Digest    string    // Hash of the DNSKEY content
	Status    DSStatus  // Result of the last configuration check
	LastCheck time.Time // Time of the last configuration check
	LastOK    time.Time // Last time that the DNSSEC configuration was OK
}
