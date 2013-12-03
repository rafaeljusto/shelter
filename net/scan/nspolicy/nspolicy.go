package nspolicy

import (
	"github.com/miekg/dns"
	"net"
	"regexp"
	"shelter/model"
	"shelter/net/scan/dnsutils"
)

var (
	// List of all nameserver policies that are going to be executed in the order defined
	// here. The order is important because the policies depends on each other, assuming
	// that something was already verified
	nsPolicies = []func(*DomainNSPolicy, *dns.Msg) model.NameserverStatus{
		(*DomainNSPolicy).cnamePolicy,
		(*DomainNSPolicy).rcodePolicy,
		(*DomainNSPolicy).authorityPolicy,
		(*DomainNSPolicy).soaPolicy,
	}
)

// DomainNSPolicy store the domain object and the version of the DNS zone. This is
// necessary because we need to check the DNS zone version on each nameserver and detect
// if they are different
type DomainNSPolicy struct {
	domain     *model.Domain // Domain object is used for glue validations
	soaVersion uint32        // Variable used to check if all nameservers have the same zone
}

// This function initialize a DomainNSPolicy object, it was created to force the
// programmer to initialize the domain object, so we don't need to check if domain is nil
// inside each method. Maybe there's a better approach (think about)
func NewDomainNSPolicy(domain *model.Domain) DomainNSPolicy {
	return DomainNSPolicy{
		domain:     domain,
		soaVersion: 0,
	}
}

// When there's a error while sending a nameserver request over the network, this method
// is responsable for detecting any usual problems. There's also some unknown problems
// that are going to be treated as DNS error, for now we are not logging the generic
// error, but maybe is a good idea if occurs to often
func (d *DomainNSPolicy) CheckNetworkError(err error) model.NameserverStatus {
	if err == nil {
		return model.NameserverStatusOK
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return model.NameserverStatusTimeout
	}

	// Some erros we don't have a specific structure to detect it, so we are going to catch
	// the problem analyzing the error message

	if match, _ := regexp.MatchString(".*lookup.*", err.Error()); match {
		return model.NameserverStatusUnknownHost
	}

	if match, _ := regexp.MatchString(".*connection refused.*", err.Error()); match {
		return model.NameserverStatusConnectionRefused
	}

	return model.NameserverStatusError
}

// Method responsable for running all nameserver policies. It will return the nameserver
// status of the first error that occurred
func (d *DomainNSPolicy) Run(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	// Something went really wrong, because if we got here there was no network error and it
	// should have a DNS response message, but as a safety check we don't allow to continue
	if dnsResponseMessage == nil {
		return model.NameserverStatusError
	}

	for _, policy := range nsPolicies {
		if status := policy(d, dnsResponseMessage); status != model.NameserverStatusOK {
			return status
		}
	}

	return model.NameserverStatusOK
}

// CNAME policy is responsable to check if there's no CNAME in the top level of the zone.
// According to the RFC CNAME resource record cannot exist with another resource record
// with the same name, as SOA resource record is mandatory in the top of the zone, CNAME
// cannot exist in the apex
func (d *DomainNSPolicy) cnamePolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	for _, rr := range dnsResponseMessage.Answer {
		if rr.Header().Name == d.domain.FQDN && rr.Header().Rrtype == dns.TypeCNAME {
			return model.NameserverStatusCanonicalName
		}
	}

	return model.NameserverStatusOK
}

// DNS response message has the return code that indicates if something went wrong, in
// this method we check it before proceding with other analyzes
func (d *DomainNSPolicy) rcodePolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	switch dnsResponseMessage.MsgHdr.Rcode {
	case dns.RcodeSuccess:
		// Everything is OK with the DNS response message. In Go every switch case has a
		// automatic break, so we don't need to do nothing here

	case dns.RcodeRefused:
		return model.NameserverStatusQueryRefused

	case dns.RcodeServerFailure:
		return model.NameserverStatusServerFailure

	case dns.RcodeNameError:
		return model.NameserverStatusUnknownDomainName

	default:
		return model.NameserverStatusError
	}

	return model.NameserverStatusOK
}

// Check if the nameserver owns the domain or not
func (d *DomainNSPolicy) authorityPolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	if !dnsResponseMessage.Authoritative {
		return model.NameserverStatusNoAuthority
	}

	return model.NameserverStatusOK
}

// Check if the zones between the nameservers have all the same version and also garantee
// that the zones have the SOA record, that indicates the authority for the zone
func (d *DomainNSPolicy) soaPolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	var soaRR *dns.SOA

	// Check if the SOA resource record exists in the response
	if rr := dnsutils.FilterFirstRR(dnsResponseMessage.Answer, dns.TypeSOA); rr != nil {
		soaRR, _ = rr.(*dns.SOA)

	} else {
		return model.NameserverStatusUnknownDomainName
	}

	// Find the SOA record to check the version of the zone. If it's the first nameserver
	// that we are checking (version 0) we don't need to compare

	if d.soaVersion == 0 {
		d.soaVersion = soaRR.Serial

	} else if d.soaVersion != soaRR.Serial {
		return model.NameserverStatusNotSynchronized
	}

	return model.NameserverStatusOK
}
