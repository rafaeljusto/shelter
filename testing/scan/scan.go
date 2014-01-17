package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"net/mail"
	"runtime"
	"shelter/config"
	"shelter/dao"
	"shelter/database/mongodb"
	"shelter/model"
	"shelter/net/scan"
	"shelter/testing/utils"
	"time"
)

var (
	configFilePath string // Path for the configuration file with all the query parameters
	report         bool   // Flag to generate the scan performance report file
)

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

type ByteSize float64

// ScanQuerierTestConfigFile is a structure to store the test configuration file data
type ScanTestConfigFile struct {
	config.Config

	DNSServerPort int
	ReportFile    string
}

func init() {
	utils.TestName = "Scan"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanQuerier test")
	flag.BoolVar(&report, "report", false, "Report flag for ScanQuerier performance")
}

func main() {
	flag.Parse()

	var scanConfig ScanTestConfigFile
	err := utils.ReadConfigFile(configFilePath, &scanConfig)
	config.ShelterConfig = scanConfig.Config

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	database, err := mongodb.Open(scanConfig.Database.URI,
		scanConfig.Database.Name)
	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	// Remove all data before starting the test. This is necessary because maybe in the last
	// test there was an error and the data wasn't removed from the database
	domainDAO.RemoveAll()

	startDNSServer(scanConfig.DNSServerPort, scanConfig.Scan.UDPMaxSize)

	domainWithNoErrors(domainDAO)

	// Scan performance report is optional and only generated when the report file
	// path parameter is given
	if report {
		scanReport(domainDAO, scanConfig)
	}

	utils.Println("SUCCESS!")
}

func domainWithNoErrors(domainDAO dao.DomainDAO) {
	domain, dnskey, rrsig, lastCheckAt, lastOKAt := generateAndSaveDomain("br.", domainDAO)

	dns.HandleFunc("br.", func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		dnsResponseMessage := new(dns.Msg)
		defer w.WriteMsg(dnsResponseMessage)

		if dnsRequestMessage.Question[0].Qtype == dns.TypeSOA {
			dnsResponseMessage = &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Authoritative: true,
				},
				Question: dnsRequestMessage.Question,
				Answer: []dns.RR{
					&dns.SOA{
						Hdr: dns.RR_Header{
							Name:   "br.",
							Rrtype: dns.TypeSOA,
							Class:  dns.ClassINET,
							Ttl:    86400,
						},
						Ns:      "ns1.br.",
						Mbox:    "rafael.justo.net.br.",
						Serial:  2013112600,
						Refresh: 86400,
						Retry:   86400,
						Expire:  86400,
						Minttl:  900,
					},
				},
			}
			dnsResponseMessage.SetReply(dnsRequestMessage)

			w.WriteMsg(dnsResponseMessage)

		} else if dnsRequestMessage.Question[0].Qtype == dns.TypeDNSKEY {
			dnsResponseMessage = &dns.Msg{
				MsgHdr: dns.MsgHdr{
					Authoritative: true,
				},
				Question: dnsRequestMessage.Question,
				Answer: []dns.RR{
					dnskey,
					rrsig,
				},
			}
			dnsResponseMessage.SetReply(dnsRequestMessage)

			w.WriteMsg(dnsResponseMessage)
		}
	})

	scan.ScanDomains()

	domain, err := domainDAO.FindByFQDN(domain.FQDN)
	if err != nil {
		utils.Fatalln("Didn't find scanned domain", err)
	}

	for _, nameserver := range domain.Nameservers {
		if nameserver.LastStatus != model.NameserverStatusOK {
			utils.Fatalln(fmt.Sprintf("Fail to validate a supposedly well configured nameserver '%s'. Found status: %s",
				nameserver.Host, model.NameserverStatusToString(nameserver.LastStatus)), err)
		}

		if nameserver.LastCheckAt.Before(lastCheckAt) ||
			nameserver.LastCheckAt.Equal(lastCheckAt) {
			utils.Fatalln(fmt.Sprintf("Last check date was not updated in nameserver '%s'",
				nameserver.Host), nil)
		}

		if nameserver.LastOKAt.Before(lastOKAt) || nameserver.LastOKAt.Equal(lastOKAt) {
			utils.Fatalln(fmt.Sprintf("Last OK date was not updated in nameserver '%s'",
				nameserver.Host), nil)
		}
	}

	for _, ds := range domain.DSSet {
		if ds.LastStatus != model.DSStatusOK {
			utils.Fatalln(fmt.Sprintf("Fail to validate a supposedly well configured DS %d. "+
				"Found status: %s", ds.Keytag, model.DSStatusToString(ds.LastStatus)), err)
		}

		if ds.LastCheckAt.Before(lastCheckAt) || ds.LastCheckAt.Equal(lastCheckAt) {
			utils.Fatalln(fmt.Sprintf("Last check date was not updated in DS %d",
				ds.Keytag), nil)
		}

		if ds.LastOKAt.Before(lastOKAt) || ds.LastOKAt.Equal(lastOKAt) {
			utils.Fatalln(fmt.Sprintf("Last OK date was not updated in DS %d",
				ds.Keytag), nil)
		}
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		utils.Fatalln(fmt.Sprintf("Error removing domain %s", domain.FQDN), err)
	}
}

// Generates a report with the amount of time of a scan
func scanReport(domainDAO dao.DomainDAO, scanConfig ScanTestConfigFile) {
	report := " #       | Total            | DPS  | Memory (MB)\n" +
		"-----------------------------------------------------\n"

	// Report variables
	scale := []int{10, 50, 100, 500, 1000, 5000,
		10000, 50000, 100000, 500000, 1000000, 5000000}

	for _, numberOfItems := range scale {
		var domains []model.Domain
		for i := 0; i < numberOfItems; i++ {
			fqdn := fmt.Sprintf("domain%d.br.", i)
			domain, dnskey, rrsig, _, _ := generateAndSaveDomain(fqdn, domainDAO)

			dns.HandleFunc(fqdn, func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
				defer w.Close()

				dnsResponseMessage := new(dns.Msg)
				defer w.WriteMsg(dnsResponseMessage)

				if dnsRequestMessage.Question[0].Qtype == dns.TypeSOA {
					dnsResponseMessage = &dns.Msg{
						MsgHdr: dns.MsgHdr{
							Authoritative: true,
						},
						Question: dnsRequestMessage.Question,
						Answer: []dns.RR{
							&dns.SOA{
								Hdr: dns.RR_Header{
									Name:   fqdn,
									Rrtype: dns.TypeSOA,
									Class:  dns.ClassINET,
									Ttl:    86400,
								},
								Ns:      fmt.Sprintf("ns1.%s", fqdn),
								Mbox:    "rafael.justo.net.br.",
								Serial:  2013112600,
								Refresh: 86400,
								Retry:   86400,
								Expire:  86400,
								Minttl:  900,
							},
						},
					}
					dnsResponseMessage.SetReply(dnsRequestMessage)

					w.WriteMsg(dnsResponseMessage)

				} else if dnsRequestMessage.Question[0].Qtype == dns.TypeDNSKEY {
					dnsResponseMessage = &dns.Msg{
						MsgHdr: dns.MsgHdr{
							Authoritative: true,
						},
						Question: dnsRequestMessage.Question,
						Answer: []dns.RR{
							dnskey,
							rrsig,
						},
					}
					dnsResponseMessage.SetReply(dnsRequestMessage)

					w.WriteMsg(dnsResponseMessage)
				}
			})

			domains = append(domains, domain)
		}

		utils.Println(fmt.Sprintf("Generating report - scale %d", numberOfItems))
		totalDuration, domainsPerSecond := calculateScanDurations(scanConfig, domains)

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		report += fmt.Sprintf("% -8d | %16s | %4d | %14.2f\n",
			numberOfItems,
			time.Duration(int64(totalDuration)).String(),
			domainsPerSecond,
			float64(memStats.Alloc)/float64(MB),
		)

		for _, domain := range domains {
			if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
				utils.Fatalln(fmt.Sprintf("Error removing domain %s during report", domain.FQDN), err)
			}
		}
	}

	utils.WriteReport(scanConfig.ReportFile, report)
}

func calculateScanDurations(scanConfig ScanTestConfigFile,
	domains []model.Domain) (totalDuration time.Duration, domainsPerSecond int64) {

	beginTimer := time.Now()
	scan.ScanDomains()
	totalDuration = time.Since(beginTimer)

	totalDurationSeconds := int64(totalDuration / time.Second)
	if totalDurationSeconds > 0 {
		domainsPerSecond = int64(len(domains)) / totalDurationSeconds

	} else {
		domainsPerSecond = int64(len(domains))
	}

	return
}

// Function to mock a domain
func generateAndSaveDomain(fqdn string, domainDAO dao.DomainDAO) (
	model.Domain, *dns.DNSKEY, *dns.RRSIG, time.Time, time.Time,
) {
	dnskey, rrsig, err := utils.GenerateKeyAndSignZone(fqdn)
	if err != nil {
		utils.Fatalln("Error creating DNSSEC keys and signatures", err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

	domain := model.Domain{
		FQDN: fqdn,
		Nameservers: []model.Nameserver{
			{
				Host: fmt.Sprintf("ns1.%s", fqdn),
				IPv4: net.ParseIP("127.0.0.1"),
			},
		},
		DSSet: []model.DS{
			{
				Keytag:     dnskey.KeyTag(),
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}

	owner, _ := mail.ParseAddress("test@rafael.net.br")
	domain.Owners = []*mail.Address{owner}

	lastCheckAt := time.Now().Add(-72 * time.Hour)
	lastOKAt := lastCheckAt.Add(-24 * time.Hour)

	// Set all nameservers with error and the last check equal of the error check interval,
	// this will force the domain to be checked
	for index, _ := range domain.Nameservers {
		domain.Nameservers[index].LastCheckAt = lastCheckAt
		domain.Nameservers[index].LastOKAt = lastOKAt
		domain.Nameservers[index].LastStatus = model.NameserverStatusServerFailure
	}

	// Set all DS records with error and the last check equal of the error check interval,
	// this will force the domain to be checked
	for index, _ := range domain.DSSet {
		domain.DSSet[index].LastCheckAt = lastCheckAt
		domain.DSSet[index].LastOKAt = lastOKAt
		domain.DSSet[index].LastStatus = model.DSStatusTimeout
	}

	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln("Error saving domain for scan scenario", err)
	}

	return domain, dnskey, rrsig, lastCheckAt, lastOKAt
}

func startDNSServer(port int, udpMaxSize uint16) {
	// Change the querier DNS port for the scan
	scan.DNSPort = port

	server := dns.Server{
		Net:     "udp",
		Addr:    fmt.Sprintf("localhost:%d", port),
		UDPSize: int(udpMaxSize),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			utils.Fatalln("Error starting DNS test server", err)
		}
	}()

	// Wait the DNS server to start before testing
	time.Sleep(1 * time.Second)
}
