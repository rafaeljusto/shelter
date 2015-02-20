// Copyright 2014 Rafael Dantas Justo. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/rafaeljusto/shelter/Godeps/_workspace/src/github.com/miekg/dns"
	"github.com/rafaeljusto/shelter/model"
	"github.com/rafaeljusto/shelter/net/scan"
	"github.com/rafaeljusto/shelter/testing/utils"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Input path is mandatory for performance tests
	ErrInputFileUndefined = errors.New("Input file path undefined")
	// Syntax error while parsing the input file. Check the format of the file in
	// readInputFile function comment
	ErrInputFileInvalidFormat = errors.New("Input file has an invalid format")
)

var (
	configFilePath string // Path for the configuration file with all the query parameters
	report         bool   // Flag to generate the scan querier performance report file
	inputReport    bool   // Flag to generate the report based on input or not
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
type ScanQuerierTestConfigFile struct {
	Server struct {
		Port int
	}

	Scan struct {
		NumberOfQueriers  int    // Number of concurrently queriers
		DomainsBufferSize int    // Size of the channel
		UDPMaxSize        uint16 // UDP max package size for firewall problems
		ConnectionRetries int    // Number of retries before setting timeout

		Timeouts struct {
			DialSeconds  time.Duration
			ReadSeconds  time.Duration
			WriteSeconds time.Duration
		}
	}

	Report struct {
		ReportFile string
		InputFile  string
		OutputFile string
	}
}

func init() {
	utils.TestName = "ScanQuerier"
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanQuerier test")
	flag.BoolVar(&report, "report", false, "Report flag for ScanQuerier performance")
	flag.BoolVar(&inputReport, "inputReport", false, "Input report flag for ScanQuerier performance")
}

func main() {
	flag.Parse()

	var config ScanQuerierTestConfigFile
	err := utils.ReadConfigFile(configFilePath, &config)

	if err == utils.ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		utils.Fatalln("Error reading configuration file", err)
	}

	utils.StartDNSServer(config.Server.Port, config.Scan.UDPMaxSize)

	domainWithNoDNSErrors(config)
	domainWithNoDNSSECErrors(config)
	domainDNSTimeout(config)
	domainDNSUnknownHost(config)

	// Scan querier performance report is optional and only generated when the report file
	// path parameter is given
	if report {
		scanQuerierReport(config)
	}

	if inputReport {
		inputScanReport(config)
	}

	utils.Println("SUCCESS!")
}

func domainWithNoDNSErrors(config ScanQuerierTestConfigFile) {
	domainsToQueryChannel := make(chan *model.Domain, config.Scan.DomainsBufferSize)
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
				Authoritative: true,
			},
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
	})

	domains := runScan(config, domainsToQueryChannel)
	for _, domain := range domains {
		if domain.FQDN != "br." ||
			domain.Nameservers[0].LastStatus != model.NameserverStatusOK {
			utils.Fatalln(fmt.Sprintf("Error checking a well configured DNS domain. "+
				"Expected FQDN 'br.' with status %d and got FQDN '%s' with status %d",
				model.NameserverStatusOK, domain.FQDN, domain.Nameservers[0].LastStatus), nil)
		}
	}
}

func domainWithNoDNSSECErrors(config ScanQuerierTestConfigFile) {
	dnskey, rrsig, err := utils.GenerateKSKAndSignZone("br.")
	if err != nil {
		utils.Fatalln("Error creating DNSSEC keys and signatures", err)
	}
	ds := dnskey.ToDS(uint8(model.DSDigestTypeSHA1))

	domainsToQueryChannel := make(chan *model.Domain, config.Scan.DomainsBufferSize)
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
				Algorithm:  utils.ConvertKeyAlgorithm(dnskey.Algorithm),
				DigestType: model.DSDigestTypeSHA1,
				Digest:     ds.Digest,
			},
		},
	}
	domainsToQueryChannel <- nil // Poison pill

	dns.HandleFunc("br.", func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		if dnsRequestMessage.Question[0].Qtype == dns.TypeSOA {
			dnsResponseMessage := &dns.Msg{
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
			dnsResponseMessage := &dns.Msg{
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

	domains := runScan(config, domainsToQueryChannel)
	for _, domain := range domains {
		if domain.FQDN != "br." ||
			domain.DSSet[0].LastStatus != model.DSStatusOK {
			utils.Fatalln(fmt.Sprintf("Error checking a well configured DNSSEC domain. "+
				"Expected FQDN 'br.' with status %d and got FQDN '%s' with status %d",
				model.DSStatusOK, domain.FQDN, domain.DSSet[0].LastStatus), nil)
		}
	}
}

func domainDNSTimeout(config ScanQuerierTestConfigFile) {
	domainsToQueryChannel := make(chan *model.Domain, config.Scan.DomainsBufferSize)
	domainsToQueryChannel <- &model.Domain{
		FQDN: "br.",
		Nameservers: []model.Nameserver{
			{
				Host: "google.com.",
			},
		},
	}
	domainsToQueryChannel <- nil // Poison pill

	domains := runScan(config, domainsToQueryChannel)
	for _, domain := range domains {
		if domain.Nameservers[0].LastStatus != model.NameserverStatusTimeout {
			utils.Fatalln("Error checking a timeout domain", nil)
		}
	}
}

func domainDNSUnknownHost(config ScanQuerierTestConfigFile) {
	domainsToQueryChannel := make(chan *model.Domain, config.Scan.DomainsBufferSize)
	domainsToQueryChannel <- &model.Domain{
		FQDN: "br.",
		Nameservers: []model.Nameserver{
			{
				Host: "br.br.",
			},
		},
	}
	domainsToQueryChannel <- nil // Poison pill

	domains := runScan(config, domainsToQueryChannel)
	for _, domain := range domains {
		if domain.Nameservers[0].LastStatus != model.NameserverStatusUnknownHost {
			utils.Fatalln(fmt.Sprintf("Error checking a unknown host. Expected status %d "+
				"and found status %d", model.NameserverStatusUnknownHost, domain.Nameservers[0].LastStatus), nil)
		}
	}
}

// Generates a report with the amount of time of a scan
func scanQuerierReport(config ScanQuerierTestConfigFile) {
	report := " #       | Total            | QPS  | Memory (MB)\n" +
		"-----------------------------------------------------\n"

	// Report variables
	scale := []int{10, 50, 100, 500, 1000, 5000,
		10000, 50000, 100000, 500000, 1000000, 5000000}

	fqdn := "domain.com.br."

	dnskey, rrsig, err := utils.GenerateKSKAndSignZone(fqdn)
	if err != nil {
		utils.Fatalln("Error creating DNSSEC keys and signatures", err)
	}
	ds := dnskey.ToDS(uint8(model.DSDigestTypeSHA1))

	dns.HandleFunc(fqdn, func(w dns.ResponseWriter, dnsRequestMessage *dns.Msg) {
		defer w.Close()

		if dnsRequestMessage.Question[0].Qtype == dns.TypeSOA {
			dnsResponseMessage := &dns.Msg{
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
						Ns:      "ns1." + fqdn,
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
			dnsResponseMessage := &dns.Msg{
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

	for _, numberOfItems := range scale {
		var domains []*model.Domain
		for i := 0; i < numberOfItems; i++ {
			// We create an object with different nameservers because we don't want to put the
			// nameserver in the query rate limit check
			domains = append(domains, &model.Domain{
				FQDN: fqdn,
				Nameservers: []model.Nameserver{
					{
						Host: fmt.Sprintf("ns%d.%s", i, fqdn),
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
			})
		}

		utils.Println(fmt.Sprintf("Generating report - scale %d", numberOfItems))
		totalDuration, queriesPerSecond, _, _ :=
			calculateScanQuerierDurations(config, domains)

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		report += fmt.Sprintf("% -8d | %16s | %4d | %14.2f\n",
			numberOfItems,
			time.Duration(int64(totalDuration)).String(),
			queriesPerSecond,
			float64(memStats.Alloc)/float64(MB),
		)
	}

	utils.WriteReport(config.Report.ReportFile, report)
}

// Generates a report with the result of a scan in the root zone file, it should be last
// last thing from the test, because it changes the DNS test port to the original one for
// real tests
func inputScanReport(config ScanQuerierTestConfigFile) {
	// Move back to default port, because we are going to query the world for real to check
	// querier performance
	scan.DNSPort = 53

	// As we are using the same domains repeatedly we should be careful about how many
	// requests we send to only one host to avoid abuses. This value should be beteween 5
	// and 10
	scan.MaxQPSPerHost = 5

	report := " #       | Total            | QPS  | Memory (MB)\n" +
		"---------------------------------------------------\n"

	domains, err := readInputFile(config.Report.InputFile)
	if err != nil {
		utils.Fatalln("Error while loading input data for report", err)
	}

	nameserverStatusCounter := 0
	nameserversStatus := make(map[model.NameserverStatus]int)

	dsStatusCounter := 0
	dsSetStatus := make(map[model.DSStatus]int)

	totalDuration, queriesPerSecond, nameserversStatus, dsSetStatus :=
		calculateScanQuerierDurations(config, domains)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	report += fmt.Sprintf("% -8d | %16s | %4d | %14.2f\n",
		len(domains),
		time.Duration(int64(totalDuration)).String(),
		queriesPerSecond,
		float64(memStats.Alloc)/float64(MB),
	)

	report += "\nNameserver Status\n" +
		"-----------------\n"
	for _, counter := range nameserversStatus {
		nameserverStatusCounter += counter
	}
	for status, counter := range nameserversStatus {
		report += fmt.Sprintf("%16s: % 3.2f%%\n",
			model.NameserverStatusToString(status),
			(float64(counter*100) / float64(nameserverStatusCounter)),
		)
	}

	report += "\nDS Status\n" +
		"---------\n"
	for _, counter := range dsSetStatus {
		dsStatusCounter += counter
	}
	for status, counter := range dsSetStatus {
		report += fmt.Sprintf("%16s: % 3.2f%%\n",
			model.DSStatusToString(status),
			(float64(counter*100) / float64(dsStatusCounter)),
		)
	}

	utils.WriteReport(config.Report.OutputFile, report)
}

func calculateScanQuerierDurations(config ScanQuerierTestConfigFile,
	domains []*model.Domain) (totalDuration time.Duration,
	queriesPerSecond int64, nameserversStatus map[model.NameserverStatus]int,
	dsSetStatus map[model.DSStatus]int) {

	domainsToQueryChannel := make(chan *model.Domain, config.Scan.DomainsBufferSize)
	go func() {
		for _, domain := range domains {
			domainsToQueryChannel <- domain
		}
		domainsToQueryChannel <- nil
	}()

	beginTimer := time.Now()
	results := runScan(config, domainsToQueryChannel)
	totalDuration = time.Since(beginTimer)

	totalDurationSeconds := int64(totalDuration / time.Second)
	if totalDurationSeconds > 0 {
		queriesPerSecond = int64(len(results)) / totalDurationSeconds

	} else {
		queriesPerSecond = int64(len(results))
	}

	nameserversStatus = make(map[model.NameserverStatus]int)
	dsSetStatus = make(map[model.DSStatus]int)

	for _, domain := range results {
		for _, nameserver := range domain.Nameservers {
			nameserversStatus[nameserver.LastStatus] += 1
		}

		for _, ds := range domain.DSSet {
			dsSetStatus[ds.LastStatus] += 1
		}
	}

	return
}

// Method responsable to configure and start scan injector for tests
func runScan(config ScanQuerierTestConfigFile,
	domainsToQueryChannel chan *model.Domain) []*model.Domain {

	dialTimeout := config.Scan.Timeouts.DialSeconds * time.Second
	readTimeout := config.Scan.Timeouts.ReadSeconds * time.Second
	writeTimeout := config.Scan.Timeouts.WriteSeconds * time.Second

	querierDispatcher := scan.NewQuerierDispatcher(
		config.Scan.NumberOfQueriers,
		config.Scan.DomainsBufferSize,
		config.Scan.UDPMaxSize,
		dialTimeout,
		readTimeout,
		writeTimeout,
		config.Scan.ConnectionRetries,
	)

	// Go routines group control created, but not used for this tests, as we are simulating
	// a collector receiver
	var scanGroup sync.WaitGroup

	domainsToSaveChannel := querierDispatcher.Start(&scanGroup, domainsToQueryChannel)

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

// Function to read input data file for performance tests. The file must use the following
// format:
//
//  <zonename1> <type1> <data1>
//  <zonename2> <type2> <data2>
//  ...
//  <zonenameN> <typeN> <dataN>
//
// Where type can be NS, A, AAAA or DS. All types, except for DS, will have only one field
// in data, for DS we will have four fields. For example:
//
// br.       NS   a.dns.br.
// br.       NS   b.dns.br.
// br.       NS   c.dns.br.
// br.       NS   d.dns.br.
// br.       NS   e.dns.br.
// br.       NS   f.dns.br.
// br.       DS   41674 5 1 EAA0978F38879DB70A53F9FF1ACF21D046A98B5C
// a.dns.br. A    200.160.0.10
// a.dns.br. AAAA 2001:12ff:0:0:0:0:0:10
// b.dns.br. A    200.189.41.10
// c.dns.br. A    200.192.233.10
// d.dns.br. A    200.219.154.10
// d.dns.br. AAAA 2001:12f8:4:0:0:0:0:10
// e.dns.br. A    200.229.248.10
// e.dns.br. AAAA 2001:12f8:1:0:0:0:0:10
// f.dns.br. A    200.219.159.10
//
func readInputFile(inputFilePath string) ([]*model.Domain, error) {
	// Input file path is necessary when we want to run a performance test, because in this
	// file we have real DNS authoritative servers
	if len(inputFilePath) == 0 {
		return nil, ErrInputFileUndefined
	}

	file, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	domainsInfo := make(map[string]*model.Domain)
	nameserversInfo := make(map[string]model.Nameserver)

	// Read line by line
	for scanner.Scan() {
		inputParts := strings.Fields(scanner.Text())
		if len(inputParts) < 3 {
			return nil, ErrInputFileInvalidFormat
		}

		zone, rrType := strings.ToLower(inputParts[0]), strings.ToUpper(inputParts[1])

		if rrType == "NS" {
			domain := domainsInfo[zone]
			if domain == nil {
				domain = &model.Domain{
					FQDN: zone,
				}
			}

			domain.Nameservers = append(domain.Nameservers, model.Nameserver{
				Host: strings.ToLower(inputParts[2]),
			})

			domainsInfo[zone] = domain

		} else if rrType == "DS" {
			domain := domainsInfo[zone]
			if domain == nil {
				domain = &model.Domain{
					FQDN: zone,
				}
			}

			if len(inputParts) < 6 {
				return nil, ErrInputFileInvalidFormat
			}

			keytag, err := strconv.Atoi(inputParts[2])
			if err != nil {
				return nil, ErrInputFileInvalidFormat
			}

			algorithm, err := strconv.Atoi(inputParts[3])
			if err != nil {
				return nil, ErrInputFileInvalidFormat
			}

			digestType, err := strconv.Atoi(inputParts[4])
			if err != nil {
				return nil, ErrInputFileInvalidFormat
			}

			domain.DSSet = append(domain.DSSet, model.DS{
				Keytag:     uint16(keytag),
				Algorithm:  model.DSAlgorithm(algorithm),
				DigestType: model.DSDigestType(digestType),
				Digest:     strings.ToUpper(inputParts[5]),
			})

			domainsInfo[zone] = domain

		} else if rrType == "A" {
			nameserver := nameserversInfo[zone]
			nameserver.Host = strings.ToLower(zone)
			nameserver.IPv4 = net.ParseIP(inputParts[2])
			nameserversInfo[zone] = nameserver

		} else if rrType == "AAAA" {
			nameserver := nameserversInfo[zone]
			nameserver.Host = strings.ToLower(zone)
			nameserver.IPv6 = net.ParseIP(inputParts[2])
			nameserversInfo[zone] = nameserver

		} else {
			return nil, ErrInputFileInvalidFormat
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var domains []*model.Domain
	for _, domain := range domainsInfo {
		for index, nameserver := range domain.Nameservers {
			if nameserverGlue, found := nameserversInfo[nameserver.Host]; found {
				domain.Nameservers[index] = nameserverGlue
			}
		}
		domains = append(domains, domain)
	}

	return domains, nil
}
