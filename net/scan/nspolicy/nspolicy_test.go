package nspolicy

import (
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/model"
	"testing"
)

// myErr was created only to test possible network errors
type myErr struct {
	err     string
	timeout bool
}

func (e myErr) Error() string {
	return e.err
}

func (e myErr) Timeout() bool {
	return e.timeout
}

func (e myErr) Temporary() bool {
	return true
}

func TestNSNetworkError(t *testing.T) {
	domain := &model.Domain{}
	domainNSPolicy := NewDomainNSPolicy(domain)

	err := myErr{
		timeout: true,
	}

	if domainNSPolicy.CheckNetworkError(err) != model.NameserverStatusTimeout {
		t.Error("Not detecting network timeout")
	}

	err = myErr{
		err:     "lookup",
		timeout: false,
	}

	if domainNSPolicy.CheckNetworkError(err) != model.NameserverStatusUnknownHost {
		t.Error("Not detecting when it couldn't resolve the nameserver")
	}

	err = myErr{
		err:     "connection refused",
		timeout: false,
	}

	if domainNSPolicy.CheckNetworkError(err) != model.NameserverStatusConnectionRefused {
		t.Error("Not detecting when the connection was refused")
	}

	err = myErr{
		err:     "generic",
		timeout: false,
	}

	if domainNSPolicy.CheckNetworkError(err) != model.NameserverStatusError {
		t.Error("Not detecting a network generic error")
	}

	if domainNSPolicy.CheckNetworkError(nil) != model.NameserverStatusOK {
		t.Error("Reporting error when everything was OK")
	}
}

func TestRunPolicies(t *testing.T) {
	domain := &model.Domain{
		FQDN: "test.com.br",
	}

	domainNSPolicy := NewDomainNSPolicy(domain)

	if domainNSPolicy.Run(nil) != model.NameserverStatusError {
		t.Error("Allowed an empty message in nameserver policies")
	}

	dnsResponseMessage := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative: true,
			Rcode:         dns.RcodeSuccess,
		},
		Answer: []dns.RR{
			&dns.SOA{
				Hdr: dns.RR_Header{
					Name:   "test.com.br",
					Rrtype: dns.TypeSOA,
				},
			},
		},
	}

	if domainNSPolicy.Run(dnsResponseMessage) != model.NameserverStatusOK {
		t.Error("Not running NS policies and checking a valid nameserver")
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative: false,
			Rcode:         dns.RcodeSuccess,
		},
		Answer: []dns.RR{
			&dns.SOA{
				Hdr: dns.RR_Header{
					Name:   "test.com.br",
					Rrtype: dns.TypeSOA,
				},
			},
		},
	}

	if domainNSPolicy.Run(dnsResponseMessage) == model.NameserverStatusOK {
		t.Error("Not running all nameserver policies and not exiting when a invalid policy is detected")
	}
}

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

func TestSOAPolicy(t *testing.T) {
	domain := &model.Domain{}
	domainNSPolicy := NewDomainNSPolicy(domain)

	dnsResponseMessage := &dns.Msg{}

	if domainNSPolicy.soaPolicy(dnsResponseMessage) !=
		model.NameserverStatusUnknownDomainName {
		t.Error("Not detecting when there's no SOA record")
	}

	dnsResponseMessage = &dns.Msg{
		Answer: []dns.RR{
			&dns.SOA{
				Hdr: dns.RR_Header{
					Rrtype: dns.TypeSOA,
				},
				Serial: 1,
			},
		},
	}

	if domainNSPolicy.soaPolicy(dnsResponseMessage) !=
		model.NameserverStatusOK {
		t.Error("Returning version problems only with one " +
			"nameserver check")
	}

	if domainNSPolicy.soaPolicy(dnsResponseMessage) !=
		model.NameserverStatusOK {
		t.Error("Returning version problems when the version is OK")
	}

	dnsResponseMessage = &dns.Msg{
		Answer: []dns.RR{
			&dns.SOA{
				Hdr: dns.RR_Header{
					Rrtype: dns.TypeSOA,
				},
				Serial: 2,
			},
		},
	}

	if domainNSPolicy.soaPolicy(dnsResponseMessage) !=
		model.NameserverStatusNotSynchronized {
		t.Error("Not detecting different versions of the same zone file")
	}
}
