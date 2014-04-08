#!/bin/sh

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

: ${GOPATH:?"Need to set GOPATH"}

usage() {
  echo "Usage: $1 <username> [--push]"
  exit 0
}

generate_certs() {
  echo ""
  echo "######################################"
  echo "#### Generating HTTPS private key ####"
  echo "######################################"
  echo ""

  # Generate a Private Key
  openssl genrsa -des3 -out key.pem 1024

  # Remove Passphrase from Key
  cp key.pem key.pem.org
  openssl rsa -in key.pem.org -out key.pem
  rm -f key.pem.org

  echo ""
  echo "######################################"
  echo "#### Generating HTTPS certificate ####"
  echo "######################################"
  echo ""

  # Generate a CSR (Certificate Signing Request)
  openssl req -new -key key.pem -out server.csr

  # Generating a Self-Signed Certificate
  openssl x509 -req -days 365 -in server.csr -signkey key.pem -out cert.pem
  rm -f server.csr
}

username=$1
if [ -z "$username" ]; then
  usage $0
fi

#if [ "$(id -u)" != "0" ]; then
#  echo "This script must be run as root" 1>&2
#  exit 1
#fi

workspace=`echo $GOPATH | cut -d: -f1`
cd $workspace/src/github.com/rafaeljusto/shelter

# Build main binary
go build shelter.go

# <src> must be the path to a file or directory relative
# to the source directory being built (also called the
# context of the build) or a remote file URL.
cd deploy/docker

rm -fr container
mkdir -p container/bin
mkdir -p container/etc/keys

mv ../../shelter container/bin/
cp entrypoint.sh container/bin/
cp ../../etc/shelter.conf.sample container/etc/shelter.conf
cp ../../etc/messages.conf container/etc/
cp -r ../../templates container/


# Generate certificates for container
cd container/etc/keys
generate_certs
cd ../../../

# Create container
sudo docker build --rm -t $username/shelter .

# Remove deploy data
rm -fr container

# Push the container to the index
if [ "$2" = "--push" ]; then
  sudo docker login
  sudo docker push $username/shelter
fi