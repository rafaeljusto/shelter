package dspolicy

import (
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/model"
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

func TestDSNetworkError(t *testing.T) {
	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     6726,
				Algorithm:  model.DSAlgorithmRSASHA1,
				DigestType: model.DSDigestTypeSHA1,
				Digest:     "56064EE6A01A9BAB7F347934D10E6AD9A4FD6DD0",
				LastStatus: model.DSStatusOK,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	if !domainDSPolicy.CheckNetworkError(nil) ||
		domain.DSSet[0].LastStatus != model.DSStatusOK {
		t.Error("Not detecting DNSSEC network without problem")
	}

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

func TestRunPolicies(t *testing.T) {
	dnskey, rrsig, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	if domainDSPolicy.Run(nil) {
		t.Error("Allowed an empty message in DS policies")
	}

	dnsResponseMessage := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode:         dns.RcodeSuccess,
			Authoritative: true,
		},
		Answer: []dns.RR{
			dnskey,
			rrsig,
		},
	}

	if !domainDSPolicy.Run(dnsResponseMessage) {
		t.Error("Not running DS policies and checking a valid DS")
	}

	dnsResponseMessage = &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Rcode:         dns.RcodeSuccess,
			Authoritative: false,
		},
		Answer: []dns.RR{
			dnskey,
			rrsig,
		},
	}

	if domainDSPolicy.Run(dnsResponseMessage) {
		t.Error("Not running all DS policies and not exiting when a invalid policy is detected")
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
	dnskey, rrsig, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
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
	dnskey, _, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	otherDNSKEY, otherRRSIG, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			otherDNSKEY,
			otherRRSIG,
		},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusNoKey {
		t.Error("Not detecting missing key")
	}
}

func TestDNSSECPolicyNoSEPKey(t *testing.T) {
	dnskey, rrsig, err :=
		generateKeyAndSignZoneWithNoSEPKey("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
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
	dnskey, _, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	otherDNSKEY, otherRRSIG, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	domainDSPolicy := NewDomainDSPolicy(domain)

	dnsResponseMessage := &dns.Msg{
		Answer: []dns.RR{
			otherDNSKEY,
			otherRRSIG,
			dnskey,
		},
	}

	if domainDSPolicy.dnssecPolicy(dnsResponseMessage) ||
		domain.DSSet[0].LastStatus != model.DSStatusNoSignature {
		t.Error("Not detecting missing signature")
	}
}

func TestDNSSECPolicyExpiredSignature(t *testing.T) {
	dnskey, rrsig, err := generateKeyAndSignZoneWithExpiredSignature("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
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
	dnskey, rrsig, err := generateKeyAndSignZone("test.br.")
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
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
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
	dnskey, rrsig, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := &model.Domain{
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
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

func TestSelectDNSKEY(t *testing.T) {
	dnskey, rrsig, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}

	otherDNSKEY, otherRRSIG, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}

	domainDSPolicy := NewDomainDSPolicy(new(model.Domain))

	dnskeys := []dns.RR{
		rrsig,
		dnskey,
		otherDNSKEY,
	}

	if domainDSPolicy.selectDNSKEY(dnskeys, dnskey.KeyTag()) == nil {
		t.Error("Not selecting a DNSKEY that exists")
	}

	dnskeys = []dns.RR{
		rrsig,
		otherDNSKEY,
		otherRRSIG,
	}

	if domainDSPolicy.selectDNSKEY(dnskeys, dnskey.KeyTag()) != nil {
		t.Error("Selecting a DNSKEY that don't exists")
	}
}

func TestSelectRRSIG(t *testing.T) {
	dnskey, rrsig, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}

	otherDNSKEY, otherRRSIG, err := generateKeyAndSignZone("test.br.")
	if err != nil {
		t.Fatal(err)
	}

	domainDSPolicy := NewDomainDSPolicy(new(model.Domain))

	rrsigs := []dns.RR{
		dnskey,
		rrsig,
		otherDNSKEY,
	}

	if domainDSPolicy.selectRRSIG(rrsigs, dnskey.KeyTag()) == nil {
		t.Error("Not selecting a RRSIG that exists")
	}

	rrsigs = []dns.RR{
		dnskey,
		otherDNSKEY,
		otherRRSIG,
	}

	if domainDSPolicy.selectRRSIG(rrsigs, dnskey.KeyTag()) != nil {
		t.Error("Selecting a RRSIG that don't exists")
	}
}

func generateKeyAndSignZone(zone string) (*dns.DNSKEY, *dns.RRSIG, error) {
	var globalErr error

	// When creating a lot of keys in a small amount of time, sometimes the systems fails to
	// generate or sign the key. For that reason we try at least 3 times of failure before
	// returning the error. Only this method has this feature because the other ones are not
	// used in performance reports
	for i := 0; i < 3; i++ {
		dnskey := &dns.DNSKEY{
			Hdr: dns.RR_Header{
				Name:   zone,
				Rrtype: dns.TypeDNSKEY,
			},
			Flags:     257,
			Protocol:  3,
			Algorithm: dns.RSASHA1NSEC3SHA1,
		}

		privateKey, err := dnskey.Generate(1024)
		if err != nil {
			globalErr = err
			continue
		}

		rrsig := &dns.RRSIG{
			Hdr: dns.RR_Header{
				Name:   zone,
				Rrtype: dns.TypeRRSIG,
			},
			TypeCovered: dns.TypeDNSKEY,
			Algorithm:   dnskey.Algorithm,
			Expiration:  uint32(time.Now().Add(10 * time.Second).Unix()),
			Inception:   uint32(time.Now().Unix()),
			KeyTag:      dnskey.KeyTag(),
			SignerName:  zone,
		}

		if err := rrsig.Sign(privateKey, []dns.RR{dnskey}); err != nil {
			globalErr = err
			continue
		}

		return dnskey, rrsig, nil
	}

	return nil, nil, globalErr
}

func generateKeyAndSignZoneWithNoSEPKey(zone string) (*dns.DNSKEY, *dns.RRSIG, error) {
	dnskey := &dns.DNSKEY{
		Hdr: dns.RR_Header{
			Name:   zone,
			Rrtype: dns.TypeDNSKEY,
		},
		Flags:     256,
		Protocol:  3,
		Algorithm: dns.RSASHA1NSEC3SHA1,
	}

	privateKey, err := dnskey.Generate(1024)
	if err != nil {
		return nil, nil, err
	}

	rrsig := &dns.RRSIG{
		Hdr: dns.RR_Header{
			Name:   zone,
			Rrtype: dns.TypeRRSIG,
		},
		TypeCovered: dns.TypeDNSKEY,
		Algorithm:   dnskey.Algorithm,
		Expiration:  uint32(time.Now().Add(10 * time.Second).Unix()),
		Inception:   uint32(time.Now().Unix()),
		KeyTag:      dnskey.KeyTag(),
		SignerName:  zone,
	}

	if err := rrsig.Sign(privateKey, []dns.RR{dnskey}); err != nil {
		return nil, nil, err
	}

	return dnskey, rrsig, nil
}

func generateKeyAndSignZoneWithExpiredSignature(zone string) (*dns.DNSKEY, *dns.RRSIG, error) {
	dnskey := &dns.DNSKEY{
		Hdr: dns.RR_Header{
			Name:   zone,
			Rrtype: dns.TypeDNSKEY,
		},
		Flags:     257,
		Protocol:  3,
		Algorithm: dns.RSASHA1NSEC3SHA1,
	}

	privateKey, err := dnskey.Generate(1024)
	if err != nil {
		return nil, nil, err
	}

	rrsig := &dns.RRSIG{
		Hdr: dns.RR_Header{
			Name:   zone,
			Rrtype: dns.TypeRRSIG,
		},
		TypeCovered: dns.TypeDNSKEY,
		Algorithm:   dnskey.Algorithm,
		Expiration:  uint32(time.Now().Add(-2 * time.Second).Unix()),
		Inception:   uint32(time.Now().Add(-5 * time.Second).Unix()),
		KeyTag:      dnskey.KeyTag(),
		SignerName:  zone,
	}

	if err := rrsig.Sign(privateKey, []dns.RR{dnskey}); err != nil {
		return nil, nil, err
	}

	return dnskey, rrsig, nil
}

func convertKeyAlgorithm(algorithm uint8) model.DSAlgorithm {
	switch algorithm {
	case dns.RSAMD5:
		return model.DSAlgorithmRSAMD5
	case dns.DH:
		return model.DSAlgorithmDH
	case dns.DSA:
		return model.DSAlgorithmDSASHA1
	case dns.ECC:
		return model.DSAlgorithmECC
	case dns.RSASHA1:
		return model.DSAlgorithmRSASHA1
	case dns.DSANSEC3SHA1:
		return model.DSAlgorithmDSASHA1NSEC3
	case dns.RSASHA1NSEC3SHA1:
		return model.DSAlgorithmRSASHA1NSEC3
	case dns.RSASHA256:
		return model.DSAlgorithmRSASHA256
	case dns.RSASHA512:
		return model.DSAlgorithmRSASHA512
	case dns.ECCGOST:
		return model.DSAlgorithmECCGOST
	case dns.ECDSAP256SHA256:
		return model.DSAlgorithmECDSASHA256
	case dns.ECDSAP384SHA384:
		return model.DSAlgorithmECDSASHA384
	case dns.INDIRECT:
		return model.DSAlgorithmIndirect
	case dns.PRIVATEDNS:
		return model.DSAlgorithmPrivateDNS
	case dns.PRIVATEOID:
		return model.DSAlgorithmPrivateOID
	}

	return model.DSAlgorithmRSASHA1
}
