#!/bin/sh
/usr/bin/mongod --quiet &
/usr/shelter/bin/shelter --config /usr/shelter/etc/shelter.conf
