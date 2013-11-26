package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"net"
	"shelter/model"
	"shelter/scan"
	"time"
)

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Config file path is a mandatory parameter
	ErrConfigFileUndefined = errors.New("Config file path undefined")
)

var (
	configFilePath string // Path for the configuration file with all the query parameters
)

// Define some scan important variables for the test enviroment, this values indicates the
// size of the channel, the number of concurrently queriers and the UDP max package size
// for firewall problems
const (
	domainsBufferSize = 10
	numberOfQueriers  = 5
	udpMaxSize        = 4096
)

// ScanQuerierTestConfigFile is a structure to store the test configuration file data
type ScanQuerierTestConfigFile struct {
	Server struct {
		Port int
	}
}

func init() {
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanQuerier test")
}

func main() {
	flag.Parse()

	configFile, err := readConfigFile()
	if err == ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		fatalln("Error reading configuration file", err)
	}

	startDNSServer(configFile.Server.Port)
	domainWithNoDNSErrors()
	domainWithNoDNSSECErrors()

	println("SUCCESS!")
}

func domainWithNoDNSErrors() {
	domainsToQueryChannel := make(chan *model.Domain, domainsBufferSize)
	domainsToQueryChannel <- &model.Domain{
		FQDN: "br.",
		Nameservers: []model.Nameserver{
			{
				Host: "ns1.br",
				IPv4: net.ParseIP("127.0.0.1"),
			},
		},
	}
	domainsToQueryChannel <- nil // Poison pill

	dns.HandleFunc("br.", func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		dnsResponseMessage := &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Id:               dnsRequestMessage.Id,
				Response:         true,
				Opcode:           dns.OpcodeQuery,
				Authoritative:    true,
				RecursionDesired: dnsRequestMessage.RecursionDesired,
				CheckingDisabled: dnsRequestMessage.CheckingDisabled,
				Rcode:            dns.RcodeSuccess,
			},
			Question: dnsRequestMessage.Question,
			Answer: []dns.RR{
				&dns.SOA{
					Hdr: dns.RR_Header{
						Name:   "br.",
						Rrtype: dns.TypeSOA,
					},
					Serial: 2013112600,
				},
			},
		}

		w.WriteMsg(dnsResponseMessage)
	})

	domains := runScan(domainsToQueryChannel)
	for _, domain := range domains {
		if domain.FQDN != "br." ||
			domain.Nameservers[0].LastStatus != model.NameserverStatusOK {
			fatalln("Error checking a well configured DNS domain", nil)
		}
	}

	dns.HandleRemove("br.")
}

func domainWithNoDNSSECErrors() {
	dnskey, rrsig, err := generateKeyAndSignZone("br.")
	if err != nil {
		fatalln("Error creating DNSSEC keys and signatures", err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domainsToQueryChannel := make(chan *model.Domain, domainsBufferSize)
	domainsToQueryChannel <- &model.Domain{
		FQDN: "br.",
		Nameservers: []model.Nameserver{
			{
				Host: "ns1.br",
				IPv4: net.ParseIP("127.0.0.1"),
			},
		},
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}
	domainsToQueryChannel <- nil // Poison pill

	dns.HandleFunc("br.", func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		dnsResponseMessage := new(dns.Msg)
		defer w.WriteMsg(dnsResponseMessage)

		if dnsRequestMessage.Question[0].Qtype == dns.TypeSOA {
			dnsResponseMessage = &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Id:               dnsRequestMessage.Id,
					Response:         true,
					Opcode:           dns.OpcodeQuery,
					Authoritative:    true,
					RecursionDesired: dnsRequestMessage.RecursionDesired,
					CheckingDisabled: dnsRequestMessage.CheckingDisabled,
					Rcode:            dns.RcodeSuccess,
				},
				Question: dnsRequestMessage.Question,
				Answer: []dns.RR{
					&dns.SOA{
						Hdr: dns.RR_Header{
							Name:   "br.",
							Rrtype: dns.TypeSOA,
						},
						Serial: 2013112600,
					},
				},
			}

			w.WriteMsg(dnsResponseMessage)

		} else if dnsRequestMessage.Question[0].Qtype == dns.TypeDNSKEY {
			dnsResponseMessage = &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Id:               dnsRequestMessage.Id,
					Response:         true,
					Opcode:           dns.OpcodeQuery,
					Authoritative:    true,
					RecursionDesired: dnsRequestMessage.RecursionDesired,
					CheckingDisabled: dnsRequestMessage.CheckingDisabled,
					Rcode:            dns.RcodeSuccess,
				},
				Question: dnsRequestMessage.Question,
				Answer: []dns.RR{
					dnskey,
					rrsig,
				},
			}

		}
	})

	domains := runScan(domainsToQueryChannel)
	for _, domain := range domains {
		if domain.FQDN != "br." ||
			domain.DSSet[0].LastStatus != model.DSStatusOK {
			fatalln("Error checking a well configured DNSSEC domain", nil)
		}
	}

	dns.HandleRemove("br.")
}

// Method responsable to configure and start scan injector for tests
func runScan(domainsToQueryChannel chan *model.Domain) []*model.Domain {
	var scanQuerierDispacther scan.QuerierDispatcher

	domainsToSaveChannel := scanQuerierDispacther.Start(domainsToQueryChannel, domainsBufferSize,
		numberOfQueriers, udpMaxSize)

	var domains []*model.Domain

	for {
		exit := false

		select {
		case domain := <-domainsToSaveChannel:
			// Detect the poison pills
			if domain == nil {
				exit = true

			} else {
				domains = append(domains, domain)
			}
		}

		if exit {
			break
		}
	}

	return domains
}

func generateKeyAndSignZone(zone string) (*dns.DNSKEY, *dns.RRSIG, error) {
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

func startDNSServer(port int) {
	// Change the querier DNS port for the scan
	scan.DNSPort = port

	server := dns.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		UDPSize: udpMaxSize,
	}

	go func() {
		server.ListenAndServe()
	}()
}

// Function to read the configuration file
func readConfigFile() (ScanQuerierTestConfigFile, error) {
	var configFile ScanQuerierTestConfigFile

	// Config file path is a mandatory program parameter
	if len(configFilePath) == 0 {
		return configFile, ErrConfigFileUndefined
	}

	confBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return configFile, err
	}

	if err := json.Unmarshal(confBytes, &configFile); err != nil {
		return configFile, err
	}

	return configFile, nil
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints a simple message without ending the test
func println(message string) {
	message = fmt.Sprintf("ScanQuerier integration test: %s", message)
	log.Println(message)
}

// Function only to add the test name before the log message. This is useful when you have
// many tests running and logging in the same file, like in a continuous deployment
// scenario. Prints an error message and ends the test
func fatalln(message string, err error) {
	message = fmt.Sprintf("ScanQuerier integration test: %s", message)
	if err != nil {
		message = fmt.Sprintf("%s. Details: %s", message, err.Error())
	}

	log.Fatalln(message)
}
