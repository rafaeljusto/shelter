#!/bin/sh

echo ""
echo "#################################"
echo "####   Starting SSH server   ####"
echo "#################################"
echo ""

/usr/sbin/sshd

echo ""
echo "#################################"
echo "#### Starting Mongodb server ####"
echo "#################################"
echo ""

/usr/bin/mongod --fork --logpath /var/log/mongodb.log

echo ""
echo "#################################"
echo "#### Starting Shelter system ####"
echo "#################################"
echo ""

# This will run in foreground
/usr/shelter/bin/shelter --config /usr/shelter/etc/shelter.conf
