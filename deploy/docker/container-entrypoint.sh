#!/bin/sh
/usr/bin/mongod --fork --logpath /var/log/mongodb.log
/usr/shelter/bin/shelter --config /usr/shelter/etc/shelter.conf
