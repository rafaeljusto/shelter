#!/bin/sh

sudo -s <<EOF
echo "deb http://dl.bintray.com/rafaeljusto/deb ./" >> /etc/apt/sources.list
if [ $? -eq  0 ]; then
  apt-get update && apt-get install -y --force-yes mongodb shelter
fi
EOF