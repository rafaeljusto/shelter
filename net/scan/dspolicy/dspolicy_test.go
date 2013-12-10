package dspolicy

import (
	"github.com/miekg/dns"
	"shelter/model"
	"shelter/testing/utils"
	"testing"
	"time"
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

	err := myErr{
		timeout: true,
	}

	if domainDSPolicy.CheckNetworkError(err) ||
		domain.DSSet[0].LastStatus != model.DSStatusTimeout {
		t.Error("Not detecting DNSSEC network timeout")
	}

	err = myErr{
		err:     "lookup",
		timeout: false,
	}

	if domainDSPolicy.CheckNetworkError(err) ||
		domain.DSSet[0].LastStatus != model.DSStatusDNSError {
		t.Error("Not detecting DNS network errors")
	}
}

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
	dnskey, rrsig, err := utils.GenerateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			dnskey,
			rrsig,
		},
	}

	if !domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusOK {
		t.Error("Not accepting a valid DS")
	}
}

func TestDNSSECPolicyMissingKey(t *testing.T) {
	dnskey, _, err := utils.GenerateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusNoKey {
		t.Error("Not detecting missing key")
	}
}

func TestDNSSECPolicyNoSEPKey(t *testing.T) {
	dnskey, rrsig, err :=
		utils.GenerateKeyAndSignZoneWithNoSEPKey("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			dnskey,
			rrsig,
		},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusNoSEP {
		t.Error("Not detecting no SEP key")
	}
}

func TestDNSSECPolicyMissingSignature(t *testing.T) {
	dnskey, _, err := utils.GenerateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			dnskey,
		},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusNoSignature {
		t.Error("Not detecting missing signature")
	}
}

func TestDNSSECPolicyExpiredSignature(t *testing.T) {
	dnskey, rrsig, err := utils.GenerateKeyAndSignZoneWithExpiredSignature("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			dnskey,
			rrsig,
		},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusExpiredSignature {
		t.Error("Not detecting expired signature")
	}
}

func TestDNSSECPolicySignatureError(t *testing.T) {
	dnskey, rrsig, err := utils.GenerateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	// Changing signature timers turns it invalid
	rrsig.Expiration = uint32(time.Now().Unix())

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			dnskey,
			rrsig,
		},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusSignatureError {
		t.Error("Not detecting signature validation errors")
	}
}

func TestDNSSECPolicyWrongDSDigest(t *testing.T) {
	dnskey, rrsig, err := utils.GenerateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     "A" + ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			dnskey,
			rrsig,
		},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusNoKey {
		t.Error("Not detecting DS digest inconsistency")
	}
}
