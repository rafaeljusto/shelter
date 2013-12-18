package model

import (
	"code.google.com/p/go.net/idna"
	"strings"
)

// Normalize the domain name to have always the same mask. The following rules applied
// are, all in lower case, no spaces in the edges (spaces in the middle are going to be
// detected by other validation mechanisms), dot at the end of the name and the ASCII
// format (for IDNA domains)
func NormalizeDomainName(domainName string) (string, error) {
	domainName = strings.ToLower(domainName)
	domainName = strings.TrimSpace(domainName)

	// We will not add the dot for empty domain name so that become root zone ("."), because
	// maybe the user forgot to inform the domain name and we don't want to make it easy for
	// him to pass in the validation verifications
	if len(domainName) > 0 && !strings.HasSuffix(domainName, ".") {
		domainName += "."
	}

	// We always manage domains in ASCII format, so we convert unicode domains, that use
	// accents according to IDNA rules
	var err error
	domainName, err = idna.ToASCII(domainName)
	return domainName, err
}

// We will store the digest always in lower case. According to RFC 4034 (section 5.3):
// "Digest MUST be represented as a sequence of case-insensitive hexadecimal digits.
// Whitespace is allowed within the hexadecimal text.", so there's no problem changing the
// letter case
func NormalizeDSDigest(digest string) string {
	digest = strings.ToLower(digest)
	digest = strings.TrimSpace(digest)
	return digest
}
