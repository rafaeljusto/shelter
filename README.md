shelter
=======

[![Build Status](https://travis-ci.org/rafaeljusto/shelter.png?branch=master)](https://travis-ci.org/rafaeljusto/shelter)
[![GoDoc](https://godoc.org/github.com/rafaeljusto/shelter?status.png)](https://godoc.org/github.com/rafaeljusto/shelter)
[![Download](https://api.bintray.com/packages/rafaeljusto/deb/shelter/images/download.png) ](https://bintray.com/rafaeljusto/deb/shelter/_latestVersion)

System created for registries to periodically validate and alert domains about DNS or
DNSSEC misconfiguration.

The idea started at the roundtable "DNSSEC: cooperation for the uptake of regional
initiatives" in LACTLD Tech Workshop occured in Panama City, Panama on September of 2013.
The roundtable was leaded by Hugo Salgado (.cl); Rafael Dantas Justo (.br); Robert Martin-
Legene (PCH). Many others participants from other registries of Latin America also
contributed with ideas for the project.

For more information check the [Wiki](https://github.com/rafaeljusto/shelter/wiki).

features
--------

* Automatically detect DNS/DNSSEC configuration problems of the registered domains
* Automatically sends e-mails notifying domain's owners of the configuration problems
* System can be deployed on registry or provider back-end infrastructure, not letting
critical data to spread to other networks
* Uses REST architecture to allow a distributted system and easy integration with other
softwares
* Multi-language support for notification's e-mails that can be distinct for each domain's
owner
* Built-in web client to manage domains easily without the necessity to develop a REST
client
* IDN support for domains' names
* Optimized scan strategy to verify all registered domains configurations
* On-the-fly domain verification interface
* Allow a cluster of MongoDB servers for data persistency

installing
----------

The install information can be found [here](https://github.com/rafaeljusto/shelter/wiki/Install-and-configure).

building
--------

The Shelter project was developed using the [Go language](http://golang.org/)

The objects are persisted using a MongoDB database.
To install it check the webpage http://www.mongodb.org/

Also, to easy run the project tests you will need the following:
* Python3 - http://www.python.org/
* Karma and dependencies - http://karma-runner.github.io
  * npm install -g karma
  * npm install -g karma-cli
  * npm install -g karma-jasmine
  * npm install -g karma-phantomjs-launcher
  * npm install -g karma-ng-html2js-preprocessor

Remember that the project directory should respect the path bellow, because the source
code dependencies can make references to this structure.

```
<GOPATH>/src/github.com/rafaeljusto/shelter
```

You can automatically retrieve the project with the desired structure using the following
command:

```
go get -u github.com/rafaeljusto/shelter
```

And finally, to build the project, just run the following command on the project root:

```
go build shelter.go
```

Optionally you can run the tests executing the following command on the project root:

```
python3 test-all.py
```

deploying
---------

To deploy the project you will need one or more programs bellow, depending on the operational system
that you choose, and if you want a CI enviroment.

* FPM - https://github.com/jordansissel/fpm (Debian packages)
* ~~Termbox - https://github.com/nsf/termbox-go (Debian packages)~~
* Docker - http://docker.io/
* ~~OpenSSL - https://www.openssl.org/ (used in Docker build)~~
* Inno Setup - http://www.jrsoftware.org/isinfo.php (Windows)

All necessary scripts to generate the packages are under the deploy folder in the project
root.
