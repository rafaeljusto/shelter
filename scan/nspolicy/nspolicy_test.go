package nspolicy

import (
	"github.com/miekg/dns"
	"shelter/model"
	"testing"
)

func TestCNAMEPolicy(t *testing.T) {
	domain := &model.Domain{
		FQDN: "test.com.br",
	}

	domainNSPolicy := NewDomainNSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			&dns.SOA{
				Hdr: dns.RR_Header{
					Name:   "test.com.br",
					Rrtype: dns.TypeSOA,
				},
			},
			&dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   "test.com.br",
					Rrtype: dns.TypeCNAME,
				},
			},
		},
	}

	if domainNSPolicy.cnamePolicy(dnsResponseMessage) !=
		model.NameserverStatusCanonicalName {
		t.Error("Not verfying CNAME in apex rule")
	}

	dnsResponseMessage = &dns.Msg{
		Answer: []dns.RR{
			&dns.SOA{
				Hdr: dns.RR_Header{
					Name:   "test.com.br",
					Rrtype: dns.TypeSOA,
				},
			},
		},
	}

	if domainNSPolicy.cnamePolicy(dnsResponseMessage) ==
		model.NameserverStatusCanonicalName {
		t.Error("Returning CNAME in apex error when there's no CNAME")
	}
}

func TestRcodePolicy(t *testing.T) {
	domain := &model.Domain{}
	domainNSPolicy := NewDomainNSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode: dns.RcodeRefused,
		},
	}

	if domainNSPolicy.rcodePolicy(dnsResponseMessage) !=
		model.NameserverStatusQueryRefused {
		t.Error("Not verfying query refused in DNS response header")
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode: dns.RcodeServerFailure,
		},
	}

	if domainNSPolicy.rcodePolicy(dnsResponseMessage) !=
		model.NameserverStatusServerFailure {
		t.Error("Not verfying server failure in DNS response header")
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode: dns.RcodeNameError,
		},
	}

	if domainNSPolicy.rcodePolicy(dnsResponseMessage) !=
		model.NameserverStatusUnknownDomainName {
		t.Error("Not verfying UDN in DNS response header")
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode: dns.RcodeBadAlg,
		},
	}

	if domainNSPolicy.rcodePolicy(dnsResponseMessage) != model.NameserverStatusError {
		t.Error("Not verfying generic error in DNS response header")
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode: dns.RcodeSuccess,
		},
	}

	if domainNSPolicy.rcodePolicy(dnsResponseMessage) != model.NameserverStatusOK {
		t.Error("Not returning OK when there's no error in " +
			"DNS response header")
	}
}

func TestAuthorityPolicy(t *testing.T) {
	domain := &model.Domain{}
	domainNSPolicy := NewDomainNSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative: false,
		},
	}

	if domainNSPolicy.authorityPolicy(dnsResponseMessage) !=
		model.NameserverStatusNoAuthority {
		t.Error("Not verfying nameserver authority " +
			"in DNS response header")
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative: true,
		},
	}

	if domainNSPolicy.authorityPolicy(dnsResponseMessage) !=
		model.NameserverStatusOK {
		t.Error("Not returning OK when the nameserver has authority " +
			"in DNS response header")
	}
}
