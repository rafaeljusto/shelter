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
		log.Fatalln("DomainDAO integration test:", err)
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
