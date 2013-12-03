#!/bin/sh
wget http://www.internic.net/domain/root.zone
cat root.zone | egrep -v "RRSIG|NSEC|DNSKEY|SOA" | awk '{ printf $1 " " $4 " "; if ($4 == "DS") printf $5 " " $6 " " $7 " " $8; else printf $5; printf "\n" }' > scan_querier.input
rm -f root.zone
