shelter
=======

[![Build Status](https://travis-ci.org/rafaeljusto/shelter.png?branch=rest)](https://travis-ci.org/rafaeljusto/shelter)

System created for registries to periodically validate and alert domains about DNS or
DNSSEC misconfiguration.

The idea started at the roundtable "DNSSEC: cooperation for the uptake of regional
initiatives" in LACTLD Tech Workshop occured in Panama City, Panama on September of 2013.
The roundtable was leaded by Hugo Salgado (.cl); Rafael Dantas Justo (.br); Robert Martin-
Legene (PCH). Many others participants from other registries of Latin America also
contributed with ideas for the project.


installation
============

Shelter project depends on the following packages:
* code.google.com/p/go.net/idna
* code.google.com/p/go.tools/cmd/cover
* github.com/miekg/dns
* labix.org/v2/mgo

All the above packages can be installed using the command:

```
go get -u <package_name>
```

All system objects are persisted using a MongoDB database.
To install it check the webpage http://www.mongodb.org/
