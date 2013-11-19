package nspolicy

import (
	"github.com/miekg/dns"
	"net"
	"regexp"
	"shelter/model"
	"shelter/scan/dnsutils"
)

var (
	nsPolicies = []func(DomainNSPolicy, *dns.Msg) model.NameserverStatus{
		DomainNSPolicy.cnamePolicy,
		DomainNSPolicy.rcodePolicy,
		DomainNSPolicy.authorityPolicy,
		DomainNSPolicy.versionPolicy,
	}
)

type DomainNSPolicy struct {
	domain     *model.Domain
	soaVersion uint32 // Variable used to check if all nameservers have the same zone
}

func NewDomainNSPolicy(domain *model.Domain) DomainNSPolicy {
	return DomainNSPolicy{
		domain:     domain,
		soaVersion: 0,
	}
}

func (d DomainNSPolicy) CheckNetworkError(err error) model.NameserverStatus {
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

func (d DomainNSPolicy) Run(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	for _, policy := range nsPolicies {
		if status := policy(d, dnsResponseMessage); status != model.NameserverStatusOK {
			return status
		}
	}

	return model.NameserverStatusOK
}

func (d DomainNSPolicy) cnamePolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	for _, rr := range dnsResponseMessage.Answer {
		if rr.Header().Name == d.domain.FQDN && rr.Header().Rrtype == dns.TypeCNAME {
			return model.NameserverStatusCanonicalName
		}
	}

	return model.NameserverStatusOK
}

func (d DomainNSPolicy) rcodePolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	switch dnsResponseMessage.MsgHdr.Rcode {
	case dns.RcodeRefused:
		return model.NameserverStatusQueryRefused

	case dns.RcodeServerFailure:
		return model.NameserverStatusServerFailure

	case dns.RcodeNameError:
		return model.NameserverStatusUnknownDomainName

	default:
		if dnsResponseMessage.MsgHdr.Rcode != dns.RcodeSuccess {
			return model.NameserverStatusError
		}
	}

	return model.NameserverStatusOK
}

func (d DomainNSPolicy) authorityPolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
	if !dnsResponseMessage.Authoritative {
		return model.NameserverStatusNoAuthority
	}

	return model.NameserverStatusOK
}

func (d DomainNSPolicy) versionPolicy(dnsResponseMessage *dns.Msg) model.NameserverStatus {
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
