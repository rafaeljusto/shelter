package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/config"
	"github.com/rafaeljusto/shelter/dao"
	"github.com/rafaeljusto/shelter/database/mongodb"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/scan"
	"github.com/rafaeljusto/shelter/testing/utils"
	"net"
	"net/mail"
	"runtime"
	"time"
)

var (
	configFilePath string     // Path for the configuration file with all the query parameters
	report         bool       // Flag to generate the scan performance report file
	cpuProfile     bool       // Write profile about the CPU when executing the report
	goProfile      bool       // Write profile about the Go routines when executing the report
	memoryProfile  bool       // Write profile about the memory usage when executing the report
	server         dns.Server // DNS server used to simulate DNS requests
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
	Report        struct {
		File    string
		Profile struct {
			CPUFile        string
			GoRoutinesFile string
			MemoryFile     string
		}
	}
}

func init() {
	utils.TestName = "Scan"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanQuerier test")
	flag.BoolVar(&report, "report", false, "Report flag for Scan performance")
	flag.BoolVar(&cpuProfile, "cpuprofile", false, "Report flag to enable CPU profile")
	flag.BoolVar(&goProfile, "goprofile", false, "Report flag to enable Go routines profile")
	flag.BoolVar(&memoryProfile, "memprofile", false, "Report flag to enable memory profile")
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

	database, databaseSession, err := mongodb.Open(scanConfig.Database.URI,
		scanConfig.Database.Name)
	if err != nil {
		utils.Fatalln("Error connecting the database", err)
	}
	defer databaseSession.Close()

	// Remove all data before starting the test. This is necessary because maybe in the last
	// test there was an error and the data wasn't removed from the database
	utils.ClearDatabase(database)

	startDNSServer(scanConfig.DNSServerPort, scanConfig.Scan.UDPMaxSize)

	domainDAO := dao.DomainDAO{
		Database: database,
	}

	domainWithNoErrors(domainDAO)
	domainWithNoErrorsOnTheFly()

	// Scan performance report is optional and only generated when the report file
	// path parameter is given
	if report {
		if cpuProfile {
			f := utils.StartCPUProfile(scanConfig.Report.Profile.CPUFile)
			defer f()
		}

		if memoryProfile {
			f := utils.StartMemoryProfile(scanConfig.Report.Profile.MemoryFile)
			defer f()
		}

		if goProfile {
			f := utils.StartGoRoutinesProfile(scanConfig.Report.Profile.GoRoutinesFile)
			defer f()
		}

		scanDAO := dao.ScanDAO{
			Database: database,
		}

		scanReport(domainDAO, scanDAO, scanConfig)
	}

	utils.Println("SUCCESS!")
}

func domainWithNoErrors(domainDAO dao.DomainDAO) {
	domain, dnskey, rrsig, lastCheckAt, lastOKAt := generateSignAndSaveDomain("br.", domainDAO)

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
				nameserver.Host, model.NameserverStatusToString(nameserver.LastStatus)), nil)
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
				"Found status: %s", ds.Keytag, model.DSStatusToString(ds.LastStatus)), nil)
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

func domainWithNoErrorsOnTheFly() {
	domain, dnskey, rrsig, lastCheckAt, lastOKAt := generateAndSignDomain("br.")

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

	scan.ScanDomain(&domain)

	for _, nameserver := range domain.Nameservers {
		if nameserver.LastStatus != model.NameserverStatusOK {
			utils.Fatalln(fmt.Sprintf("Fail to validate a supposedly well configured nameserver '%s'. Found status: %s",
				nameserver.Host, model.NameserverStatusToString(nameserver.LastStatus)), nil)
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
				"Found status: %s", ds.Keytag, model.DSStatusToString(ds.LastStatus)), nil)
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
}

type ReportHandler struct {
	DNSKEY     *dns.DNSKEY
	PrivateKey dns.PrivateKey
}

func (r ReportHandler) ServeDNS(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
	dnsResponseMessage := new(dns.Msg)
	defer w.WriteMsg(dnsResponseMessage)

	fqdn := dnsRequestMessage.Question[0].Name

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
		rrsig, err := utils.SignKey(fqdn, r.DNSKEY, r.PrivateKey)
		if err != nil {
			utils.Fatalln(fmt.Sprintf("Error signing zone %s", fqdn), err)
		}

		dnsResponseMessage = &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Authoritative: true,
			},
			Question: dnsRequestMessage.Question,
			Answer: []dns.RR{
				r.DNSKEY,
				rrsig,
			},
		}
		dnsResponseMessage.SetReply(dnsRequestMessage)

		w.WriteMsg(dnsResponseMessage)
	}
}

// Generates a report with the amount of time of a scan
func scanReport(domainDAO dao.DomainDAO, scanDAO dao.ScanDAO, scanConfig ScanTestConfigFile) {
	report := " #       | Total            | DPS  | Memory (MB)\n" +
		"-----------------------------------------------------\n"

	// Report variables
	scale := []int{10, 50, 100, 500, 1000, 5000,
		10000, 50000, 100000, 500000, 1000000, 5000000}

	dnskey, privateKey, err := utils.GenerateKey()
	if err != nil {
		utils.Fatalln("Error generating DNSKEY", err)
	}

	reportHandler := ReportHandler{
		DNSKEY:     dnskey,
		PrivateKey: privateKey,
	}

	server.Handler = reportHandler
	dns.DefaultServeMux = nil

	for _, numberOfItems := range scale {
		utils.Println(fmt.Sprintf("Generating report - scale %d", numberOfItems))

		for i := 0; i < numberOfItems; i++ {
			if i%1000 == 0 {
				utils.PrintProgress("building scenario", (i*100)/numberOfItems)
			}

			fqdn := fmt.Sprintf("domain%d.br.", i)
			generateAndSaveDomain(fqdn, domainDAO, dnskey)
		}

		utils.PrintProgress("building scenario", 100)
		totalDuration, domainsPerSecond := calculateScanDurations(numberOfItems, scanDAO)

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		report += fmt.Sprintf("% -8d | %16s | %4d | %14.2f\n",
			numberOfItems,
			time.Duration(int64(totalDuration)).String(),
			domainsPerSecond,
			float64(memStats.Alloc)/float64(MB),
		)

		if err := domainDAO.RemoveAll(); err != nil {
			// When the result set is too big to remove, we got a timeout error from the
			// connection, but it's ok
			//utils.Fatalln("Error removing domains generated for report", err)
		}
	}

	utils.WriteReport(scanConfig.Report.File, report)
}

func calculateScanDurations(numberOfDomains int, scanDAO dao.ScanDAO) (
	totalDuration time.Duration, domainsPerSecond int64,
) {

	beginTimer := time.Now()
	scan.ScanDomains()
	totalDuration = time.Since(beginTimer)

	totalDurationSeconds := int64(totalDuration / time.Second)
	if totalDurationSeconds > 0 {
		domainsPerSecond = int64(numberOfDomains) / totalDurationSeconds

	} else {
		domainsPerSecond = int64(numberOfDomains)
	}

	// As we are running a lot of scans at the same time, and the scan information unique
	// key is the start time of the scan, we must clear the database to avoid log messages
	// of scan insert errors
	scanDAO.RemoveAll()

	return
}

func generateAndSignDomain(fqdn string) (
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

	return domain, dnskey, rrsig, lastCheckAt, lastOKAt
}

// Function to mock a domain
func generateSignAndSaveDomain(fqdn string, domainDAO dao.DomainDAO) (
	model.Domain, *dns.DNSKEY, *dns.RRSIG, time.Time, time.Time,
) {
	domain, dnskey, rrsig, lastCheckAt, lastOKAt := generateAndSignDomain(fqdn)

	if err := domainDAO.Save(&domain); err != nil {
		utils.Fatalln("Error saving domain", err)
	}

	return domain, dnskey, rrsig, lastCheckAt, lastOKAt
}

// Function to mock a domain
func generateAndSaveDomain(fqdn string, domainDAO dao.DomainDAO, dnskey *dns.DNSKEY) {
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
		utils.Fatalln(fmt.Sprintf("Fail to save domain %s", domain.FQDN), err)
	}
}

func startDNSServer(port int, udpMaxSize uint16) {
	// Change the querier DNS port for the scan
	scan.DNSPort = port

	server = dns.Server{
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
