package dspolicy

import (
	"github.com/miekg/dns"
	"shelter/model"
	"testing"
)

func TestDNSHeaderPolicy(t *testing.T) {
	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     6726,
				Algorithm:  model.DSAlgorithmRSASHA1,
				DigestType: model.DSDigestTypeSHA1,
				Digest:     "56064EE6A01A9BAB7F347934D10E6AD9A4FD6DD0",
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode:         dns.RcodeRefused,
			Authoritative: true,
		},
	}

	if domainDSPolicy.dnsHeaderPolicy(dnsResponseMessage) {
		t.Error("Did not detect DNS error package before analyzing DNSSEC")
	}

	for _, ds := range domain.DSSet {
		if ds.LastStatus != model.DSStatusDNSError {
			t.Error("Did not set DNS error on the DS records")
		}
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode:         dns.RcodeSuccess,
			Authoritative: false,
		},
	}

	if domainDSPolicy.dnsHeaderPolicy(dnsResponseMessage) {
		t.Error("Did not detect authority problem before analyzing DNSSEC")
	}

	for _, ds := range domain.DSSet {
		if ds.LastStatus != model.DSStatusDNSError {
			t.Error("Did not set DNS error on the DS records")
		}
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode:         dns.RcodeSuccess,
			Authoritative: true,
		},
	}

	if !domainDSPolicy.dnsHeaderPolicy(dnsResponseMessage) {
		t.Error("Not allowing a goo DNS package to be analyzed")
	}
}

func TestDNSSECPolicy(t *testing.T) {
	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     6726,
				Algorithm:  model.DSAlgorithmRSASHA1,
				DigestType: model.DSDigestTypeSHA1,
				Digest:     "56064EE6A01A9BAB7F347934D10E6AD9A4FD6DD0",
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			&dns.DNSKEY{
				Hdr: dns.RR_Header{
					Name:   "com.br.",
					Rrtype: dns.TypeDNSKEY,
					Class:  dns.ClassINET,
					Ttl:    21600,
				},
				Flags:     257,
				Protocol:  3,
				Algorithm: dns.RSASHA1NSEC3SHA1,
				PublicKey: "AwEAAagYGg2mmhHIs7W65vXX2Q+/7T8C4exh0i6oXU+GP+m4CkEBe5s6 geOpb95WOujmN0/f03AHQ2SICBoifwWrkWpwGuGrQmBKh9TPJRzM9ccu YhCkIyKqjDEBYoZHs8sNusTHoqLjzBTLXPom0oI4zG6N3b9+qVeUOihU 3Q7te97D",
			},
			&dns.RRSIG{
				Hdr: dns.RR_Header{
					Name:   "com.br.",
					Rrtype: dns.TypeRRSIG,
					Class:  dns.ClassINET,
					Ttl:    21600,
				},
				TypeCovered: dns.TypeDNSKEY,
				Algorithm:   dns.RSASHA1NSEC3SHA1,
				Labels:      2,
				OrigTtl:     21600,
				Expiration:  1385719200,
				Inception:   1385114400,
				KeyTag:      6726,
				SignerName:  "com.br.",
				Signature:   "KXUjzrq0VCctQ/qkc24quNca7kyrl0ZKkdIpcbzGLmgICE5gIDIDtP9+ Qwb/bemYWjxrQ6lCtPpgWn/KTCwZLgPdScSoD+Bj0RCkUgnfqXs6u3mn p3We+QAO+XhFbWIt3QbjYZauEA6BkFqPnhfIYOL/Thmnva/F59N6zs/E HFc=",
			},
		},
	}

	if !domainDSPolicy.dnssecPolicy(dnsResponseMessage) {
		t.Error("Not accepting a valid DS")
	}
}
