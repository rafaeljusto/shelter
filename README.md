shelter
=======

[![Build Status](https://travis-ci.org/rafaeljusto/shelter.png?branch=rest)](https://travis-ci.org/rafaeljusto/shelter)

Program created for registries to periodically validate and alert domains about DNS or
DNSSEC misconfiguration

installation
============

Shelter project depends on the following packages:
* code.google.com/p/go.tools/cmd/cover
* labix.org/v2/mgo
* github.com/miekg/dns
* code.google.com/p/go.net/idna

All the above packages can be installed using the command:

```
go get <package_name>
```

All system objects are persisted using a MongoDB database, this NoSQL database was choosen
because we can embed dependency objects, and avoid table relations in a relational
database scenario. Without the relations the code is cleaner and the database I/O is
faster. To install it check the webpage http://www.mongodb.org/

performance
===========

The domain persistence performance is listed in the table bellow, where "#" represents the
number of domains used in the operation. Each line was executed 5 times to calculate the
average time. Also, the "find" operation represents 3 object searchs in differents parts
of the range. For this results we are using 1000 concurrent agents (go routines) for
insert and remove operations, in an Intel Core i7-2600 3.40GHz with 8GiB of memory, using
the Ubuntu 12.10 (amd64) operational system.

```
 #       | Total           | Insert          | Find         | Remove
-----------------------------------------------------------------------------
 10      |      3.513014ms |      1.442383ms |    558.179us |      1.512109ms
 50      |      7.421929ms |      3.566502ms |    423.251us |      3.431936ms
 100     |     11.027912ms |      4.701412ms |    420.808us |      5.905437ms
 500     |      51.31504ms |     25.217248ms |     799.99us |     25.297463ms
 1000    |     96.370863ms |     46.482595ms |    1.00416ms |     48.883654ms
 5000    |      503.0236ms |    228.404845ms |   1.026406ms |    273.591873ms
 10000   |    1.026367147s |    484.794695ms |   1.475828ms |    540.096072ms
 50000   |    5.030043687s |    2.260250365s |  36.661348ms |    2.733131384s
 100000  |   10.211076193s |    4.500565689s |  36.303439ms |    5.674206572s
 500000  |    56.49850228s |   26.533725121s | 268.502266ms |   29.696274265s
 1000000 | 1m53.999031752s |   52.733993985s | 489.089955ms |  1m0.775947206s
 5000000 | 9m41.845527592s | 4m28.462709335s | 1.976624332s | 5m11.406193319s
```

Now the network layer that is responsable for sending, receiveing and analyzing the DNS
packages had the performance shown in the table bellow, where "#" represents the number of
domains used in the round of checks and QPS is the number of requests per second that the
system could process. For this results were used 400 concurrent queriers (go routines)
each one with an input wait list that can store 100 domains. The hardware used was an
Intel Core i7-2600 3.40GHz with 8GiB of memory, using the Ubuntu 12.10 (amd64) operational
system. All the requests were sent to a local DNS server.

```
 #       | Total            | QPS  | Memory (MB)
-----------------------------------------------------
 10      |       2.817312ms |   10 |           1.08
 50      |      14.766817ms |   50 |           1.24
 100     |      25.443762ms |  100 |           1.91
 500     |     1.059420731s |  500 |           4.31
 1000    |     1.054760148s | 1000 |           4.31
 5000    |     2.060903121s | 2500 |           8.07
 10000   |     5.402801345s | 2000 |           9.50
 50000   |    25.349122233s | 2000 |          24.69
 100000  |    49.294494104s | 2040 |          50.19
 500000  |   4m1.158914516s | 2074 |         390.90
 1000000 |  8m16.682795099s | 2016 |         488.53
 5000000 | 42m22.573865695s | 1966 |        4182.01
```

Bellow is the same test but for the root zone file.

```
 #       | Total            | QPS  | Memory (MB)
---------------------------------------------------
 341     |     7.154577009s |   48 |          12.18

Nameserver Status
-----------------
              OK:  92.45%
              UH:  0.16%
        QREFUSED:  0.05%
        NOTSYNCH:  5.11%
        SERVFAIL:  0.11%
            NOAA:  0.05%
         TIMEOUT:  1.63%
        CREFUSED:  0.16%
           ERROR:  0.27%

DS Status
---------
              OK:  82.50%
          SIGERR:  4.50%
           NOSIG:  5.50%
         TIMEOUT:  5.50%
           NOKEY:  2.00%
```

Finally, the full scan performance using a local DNS server that always answers with
authority, but takes a time to answer, because it signs the zones on-the-fly for the RRSIG
records. On the results bellow "#" represents the number of domains used in the round of
scan and DPS is the number of domains per second that the system could process. For this
results were used 400 concurrent queriers (go routines) each one with an input wait list
that can store 100 domains. The hardware used was an Intel Core i7-2600 3.40GHz with 8GiB
of memory, using the Ubuntu 12.10 (amd64) operational system. All the requests were sent
to a local DNS server.

```
 #       | Total             | DPS  | Memory (MB)
----------------------------------------------------
 10      |       4.727656ms  |   10 |           1.33
 50      |      21.922581ms  |   50 |           2.44
 100     |      38.108105ms  |  100 |           2.28
 500     |     1.063928826s  |  500 |           6.36
 1000    |     1.094325568s  | 1000 |           4.79
 5000    |     1.790971925s  | 5000 |           8.84
 10000   |     9.532850519s  | 1111 |           7.55
 50000   |    33.772397502s  | 1515 |           9.78
 100000  |   1m27.57901475s  | 1149 |          12.64
 500000  |  6m44.081205086s  | 1237 |           9.24
 1000000 | 13m32.216240243s  | 1231 |          19.04
 5000000 | 1h9m48.514012263s | 1193 |           9.25
```
