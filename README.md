shelter
=======

Program created for registries to periodically validate and alert domains about DNS or
DNSSEC misconfiguration

installation
============

Shelter project depends on the following packages:
* labix.org/v2/mgo
* github.com/miekg/dns

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
system. The input data (domains) was generated from the root zone.

```
 #       | Total            | QPS
----------------------------------
 10      |     4.867428051s |    2
 50      |    10.573376628s |    5
 100     |     7.375709482s |   14
 500     |      9.60334264s |   55
 1000    |     8.682962307s |  125
 5000    |    10.913498722s |  500

Nameserver Status
-----------------
            NOAA:  0.01%
              OK:  84.69%
        NOTSYNCH:  14.88%
         TIMEOUT:  0.32%
        QREFUSED:  0.02%
        CREFUSED:  0.02%
              UH:  0.02%
           ERROR:  0.03%
        SERVFAIL:  0.01%

DS Status
---------
              OK:  95.90%
         TIMEOUT:  2.91%
           NOSIG:  0.46%
          SIGERR:  0.54%
           NOKEY:  0.20%
```
