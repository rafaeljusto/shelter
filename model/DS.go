package model

import (
	"time"
)

// List of possible DS algorithms (RFC 4034 - A.1, RFC 5155, RFC 5702, RFC 5933 and RFC
// 6605). Only algorithms used for signing were listed here
const (
	DSAlgorithmDSASHA1      DSAlgorithm = 3   // DSA/SHA-1 [DSA]
	DSAlgorithmRSASHA1      DSAlgorithm = 5   // RSA/SHA-1 [RSASHA1]
	DSAlgorithmRSASHA1NSEC3 DSAlgorithm = 7   // RSA/SHA1-NSEC3 [RSASHA1-NSEC3]
	DSAlgorithmRSASHA256    DSAlgorithm = 8   // RSA/SHA-256 [RSASHA256]
	DSAlgorithmRSASHA512    DSAlgorithm = 10  // RSA/SHA-512 [RSASHA512]
	DSAlgorithmGOST         DSAlgorithm = 12  // GOST R 34.10-2001
	DSAlgorithmECDSASHA256  DSAlgorithm = 13  // ECDSA/SHA-256 - Elliptic Curve Digital Signature
	DSAlgorithmECDSASHA384  DSAlgorithm = 14  // ECDSA/SHA-384 - Elliptic Curve Digital Signature
	DSAlgorithmPrivateDNS   DSAlgorithm = 253 // Private [PRIVATEDNS]
	DSAlgorithmPrivateOID   DSAlgorithm = 254 // Private [PRIVATEOID]
)

// DSAlgorithm is a number that represents one of the possible DS algorithms listed in the
// constant group above
type DSAlgorithm uint8

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
	Keytag     uint16      // DNSKEY's identification number
	Algorithm  DSAlgorithm // DNSKEY's algorithm
	Digest     string      // Hash of the DNSKEY content
	LastStatus DSStatus    // Result of the last configuration check
	LastCheck  time.Time   // Time of the last configuration check
	LastOK     time.Time   // Last time that the DNSSEC configuration was OK
}
