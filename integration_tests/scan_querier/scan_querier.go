package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"shelter/model"
	"shelter/net/scan"
	"strconv"
	"strings"
	"sync"
	"time"
)

// List of possible errors in this test. There can be also other errors from low level
// structures
var (
	// Config file path is a mandatory parameter
	ErrConfigFileUndefined = errors.New("Config file path undefined")
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
	flag.StringVar(&configFilePath, "config", "", "Configuration file for ScanQuerier test")
	flag.BoolVar(&report, "report", false, "Report flag for ScanQuerier performance")
	flag.BoolVar(&inputReport, "inputReport", false, "Input report flag for ScanQuerier performance")
}

func main() {
	flag.Parse()

	config, err := readConfigFile()
	if err == ErrConfigFileUndefined {
		fmt.Println(err.Error())
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return

	} else if err != nil {
		fatalln("Error reading configuration file", err)
	}

	startDNSServer(config.Server.Port, config.Scan.UDPMaxSize)

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

	println("SUCCESS!")
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
			fatalln("Error checking a well configured DNS domain", nil)
		}
	}

	dns.HandleRemove("br.")
}

func domainWithNoDNSSECErrors(config ScanQuerierTestConfigFile) {
	dnskey, rrsig, err := generateKeyAndSignZone("br.")
	if err != nil {
		fatalln("Error creating DNSSEC keys and signatures", err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

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

	domains := runScan(config, domainsToQueryChannel)
	for _, domain := range domains {
		if domain.FQDN != "br." ||
			domain.DSSet[0].LastStatus != model.DSStatusOK {
			fatalln("Error checking a well configured DNSSEC domain", nil)
		}
	}

	dns.HandleRemove("br.")
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
			fatalln("Error checking a timeout domain", nil)
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
			fatalln("Error checking a unknown host", nil)
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

	dnskey, rrsig, err := generateKeyAndSignZone(fqdn)
	if err != nil {
		fatalln("Error creating DNSSEC keys and signatures", err)
	}
	ds := dnskey.ToDS(int(model.DSDigestTypeSHA1))

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
						Algorithm:  convertKeyAlgorithm(dnskey.Algorithm),
						DigestType: model.DSDigestTypeSHA1,
						Digest:     ds.Digest,
					},
				},
			})
		}

		println(fmt.Sprintf("Generating report - scale %d", numberOfItems))
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

	// If we found a report file in the current path, rename it so we don't lose the old
	// data. We are going to use the modification date from the file. We also don't check
	// the errors because we really don't care
	if file, err := os.Open(config.Report.ReportFile); err == nil {
		newFilename := config.Report.ReportFile + ".old-"

		if fileStatus, err := file.Stat(); err == nil {
			newFilename += fileStatus.ModTime().Format("20060102150405")

		} else {
			// Did not find the modification date, so lets use now
			newFilename += time.Now().Format("20060102150405")
		}

		// We don't use defer because we want to rename it before the end of scope
		file.Close()

		os.Rename(config.Report.ReportFile, newFilename)
	}

	ioutil.WriteFile(config.Report.ReportFile, []byte(report), 0444)
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
		fatalln("Error while loading input data for report", err)
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

	// If we found a report file in the current path, rename it so we don't lose the old
	// data. We are going to use the modification date from the file. We also don't check
	// the errors because we really don't care
	if file, err := os.Open(config.Report.OutputFile); err == nil {
		newFilename := config.Report.OutputFile + ".old-"

		if fileStatus, err := file.Stat(); err == nil {
			newFilename += fileStatus.ModTime().Format("20060102150405")

		} else {
			// Did not find the modification date, so lets use now
			newFilename += time.Now().Format("20060102150405")
		}

		// We don't use defer because we want to rename it before the end of scope
		file.Close()

		os.Rename(config.Report.OutputFile, newFilename)
	}

	ioutil.WriteFile(config.Report.OutputFile, []byte(report), 0444)
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
			fatalln("Error starting DNS test server", err)
		}
	}()

	// Wait the DNS server to start before testing
	time.Sleep(1 * time.Second)
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
