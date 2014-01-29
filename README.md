shelter
=======

[![Build Status](https://travis-ci.org/rafaeljusto/shelter.png?branch=rest)](https://travis-ci.org/rafaeljusto/shelter)

Program created for registries to periodically validate and alert domains about DNS or
DNSSEC misconfiguration

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

All system objects are persisted using a MongoDB database, this NoSQL database was choosen
because we can embed dependency objects, and avoid table relations in a relational
database scenario. Without the relations the code is cleaner and the database I/O is
faster. To install it check the webpage http://www.mongodb.org/
