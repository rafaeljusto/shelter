package dspolicy

import (
	"github.com/miekg/dns"
	"net"
	"shelter/model"
	"shelter/scan/dnsutils"
	"time"
)

var (
	dsPolicies = []func(DomainDSPolicy, *dns.Msg) bool{
		DomainDSPolicy.dnsHeaderPolicy,
		DomainDSPolicy.recordsMatchPolicy,
	}
)

type DomainDSPolicy struct {
	domain *model.Domain
}

func NewDomainDSPolicy(domain *model.Domain) DomainDSPolicy {
	return DomainDSPolicy{
		domain: domain,
	}
}

func (d DomainDSPolicy) CheckNetworkError(err error) bool {
	if err == nil {
		return true
	}

	// We can have timeouts only for DNSSEC requests, because usually the response is bigger
	// and firewalls are not configured for big UDP packages, or for DNS over TCP
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		for index, _ := range d.domain.DSSet {
			d.domain.DSSet[index].ChangeStatus(model.DSStatusTimeout)
		}
		return false
	}

	// Other types of network errors are not a specific problem of DNSSEC configuration, so
	// let's just set a status for the user to fix the DNS configuration to make the DNSSEC
	// configuration check possible
	for index, _ := range d.domain.DSSet {
		d.domain.DSSet[index].ChangeStatus(model.DSStatusDNSError)
	}
	return false
}

func (d DomainDSPolicy) Run(dnsResponseMessage *dns.Msg) bool {
	for _, policy := range dsPolicies {
		if !policy(d, dnsResponseMessage) {
			return false
		}
	}

	return true
}

func (d DomainDSPolicy) dnsHeaderPolicy(dnsResponseMessage *dns.Msg) bool {
	if dnsResponseMessage.Rcode == dns.RcodeSuccess &&
		dnsResponseMessage.MsgHdr.Authoritative {
		return true
	}

	// Authority errors are not a specific problem of DNSSEC configuration, so let's just
	// set a status for the user to fix the DNS configuration to make the DNSSEC
	// configuration check possible
	for index, _ := range d.domain.DSSet {
		d.domain.DSSet[index].ChangeStatus(model.DSStatusDNSError)
	}
	return false
}

func (d DomainDSPolicy) recordsMatchPolicy(dnsResponseMessage *dns.Msg) bool {
	// Get all DNSSEC public keys
	dnskeys := dnsutils.FilterRRs(dnsResponseMessage.Answer, dns.TypeDNSKEY)

	// Get all signatures from the DNS response
	rrsigs := dnsutils.FilterRRs(dnsResponseMessage.Answer, dns.TypeRRSIG)

	success := true
	for index, ds := range d.domain.DSSet {
		// Find the DNSSEC public key related to the DS
		var selectedDNSKEY *dns.DNSKEY
		for _, rr := range dnskeys {
			dnskey, ok := rr.(*dns.DNSKEY)
			if !ok {
				continue
			}

			if dnskey.KeyTag() == ds.Keytag {
				selectedDNSKEY = dnskey
				break
			}
		}

		// Find the signature of the DNSSEC key that signed the keyset
		var selectedRRSIG *dns.RRSIG
		for _, rr := range rrsigs {
			rrsig, ok := rr.(*dns.RRSIG)
			if !ok {
				continue
			}

			if rrsig.KeyTag == ds.Keytag {
				selectedRRSIG = rrsig
				break
			}
		}

		if selectedDNSKEY == nil {
			d.domain.DSSet[index].ChangeStatus(model.DSStatusNoKey)
			success = false
			continue
		}

		if selectedRRSIG == nil {
			d.domain.DSSet[index].ChangeStatus(model.DSStatusNoSignature)
			success = false
			continue
		}

		// Check if the DNSSEC key related to the DS has the security entry point. Check RFCs
		// 3755 and 4034
		if (selectedDNSKEY.Flags & (1 << 15)) == 0 {
			d.domain.DSSet[index].ChangeStatus(model.DSStatusNoSEP)
			success = false
			continue
		}

		// We store the DS expiration date to alert clients whenever an expiration date is
		// near. There's no status in DS to define a near expiration state, because this
		// isn't a problem
		d.domain.DSSet[index].ExpiresAt = time.Unix(int64(selectedRRSIG.Expiration), 0)

		// Check signature expiration
		if ds.ExpiresAt.Before(time.Now()) {
			d.domain.DSSet[index].ChangeStatus(model.DSStatusExpiredSignature)
			success = false
			continue
		}

		// Check signature consistency
		if err := selectedRRSIG.Verify(selectedDNSKEY, dnskeys); err != nil {
			d.domain.DSSet[index].ChangeStatus(model.DSStatusSignatureError)
			success = false
			continue
		}

		// Check DNSKEY hash is the same of the DS digest
		if selectedDNSKEY.ToDS(int(ds.DigestType)).Digest != ds.Digest {
			d.domain.DSSet[index].ChangeStatus(model.DSStatusNoKey)
			success = false
			continue
		}

		d.domain.DSSet[index].ChangeStatus(model.DSStatusOK)
	}

	return success
}
