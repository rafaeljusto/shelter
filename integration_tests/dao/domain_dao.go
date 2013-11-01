package main

import (
	"log"
	"net"
	"net/mail"
	"shelter/dao"
	"shelter/model"
)

// This test objective is to verify the domain data persistence. The strategy is to insert
// and search for the information. Check for insert/update consistency (updates don't
// create a new element) and if the object id is set on creation

func main() {
	domain := newDomain()

	var domainDAO dao.DomainDAO
	if err := domainDAO.Save(&domain); err != nil {
		log.Fatalln("DomainDAO integration test: Couldn't save domain in database.", err)
	}

	if domainRetrieved, err := domainDAO.FindByFQDN(domain.FQDN); err != nil {
		log.Fatalln("DomainDAO integration test: Couldn't find domain in database.", err)

	} else if !compareDomains(domain, domainRetrieved) {
		log.Fatalln("DomainDAO integration test: Domain in being persisted wrongly")
	}

	if err := domainDAO.RemoveByFQDN(domain.FQDN); err != nil {
		log.Fatalln("DomainDAO integration test: Error while trying to remove a domain.", err)
	}

	if _, err := domainDAO.FindByFQDN(domain.FQDN); err == nil {
		log.Fatalln("DomainDAO integration test: Domain was not removed from database")
	}

	log.Println("DomainDAO integration test: SUCCESS!")
}

// Function to mock a domain object
func newDomain() model.Domain {
	var domain model.Domain
	domain.FQDN = "rafael.net.br"

	domain.Nameservers = []model.Nameserver{
		{
			Host: "ns1.rafael.net.br",
			IPv4: net.ParseIP("127.0.0.1"),
			IPv6: net.ParseIP("::1"),
		},
		{
			Host: "ns2.rafael.net.br",
			IPv4: net.ParseIP("127.0.0.2"),
		},
	}

	domain.DSSet = []model.DS{
		{
			Keytag:    1234,
			Algorithm: model.DSAlgorithmRSASHA1,
			Digest:    "A790A11EA430A85DA77245F091891F73AA740483",
		},
	}

	owner, _ := mail.ParseAddress("test@rafael.net.br")
	domain.Owners = []*mail.Address{owner}

	return domain
}

// Function to compare if two domains are equal, cannot use operator == because of the
// slices inside the domain object
func compareDomains(d1, d2 model.Domain) bool {
	if d1.Id != d2.Id || d1.FQDN != d2.FQDN {
		return false
	}

	if len(d1.Nameservers) != len(d2.Nameservers) {
		return false
	}

	for i := 0; i < len(d1.Nameservers); i++ {
		// Cannot compare the nameservers directly with operator == because of the
		// pointers for IP addresses
		if d1.Nameservers[i].Host != d2.Nameservers[i].Host ||
			d1.Nameservers[i].IPv4.String() != d2.Nameservers[i].IPv4.String() ||
			d1.Nameservers[i].IPv6.String() != d2.Nameservers[i].IPv6.String() ||
			d1.Nameservers[i].LastStatus != d2.Nameservers[i].LastStatus ||
			d1.Nameservers[i].LastCheck != d2.Nameservers[i].LastCheck ||
			d1.Nameservers[i].LastOK != d2.Nameservers[i].LastOK {
			return false
		}
	}

	if len(d1.DSSet) != len(d2.DSSet) {
		return false
	}

	for i := 0; i < len(d1.DSSet); i++ {
		if d1.DSSet[i] != d2.DSSet[i] {
			return false
		}
	}

	if len(d1.Owners) != len(d2.Owners) {
		return false
	}

	for i := 0; i < len(d1.Owners); i++ {
		if d1.Owners[i].String() != d2.Owners[i].String() {
			return false
		}
	}

	return true
}
