shelter
=======

Program created for registries to periodically validate and alert domains about DNS or
DNSSEC misconfiguration

installation
============

The shelter projects depends on the following packages:
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
number of domains used in the operation. Each line was executed 5 times to calculate
the average time. Also, the "find" operation represents 3 object searchs in differents
parts of the range. For this results we are using 1000 concurrent agents (go routines) for
insert and remove operations.

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
