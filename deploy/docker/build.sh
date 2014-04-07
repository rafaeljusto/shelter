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
  # generate root key
  openssl genrsa -out root-key.pem 2048

  # generate root certificate
  openssl req -new -x509 -days 3650 -key root-key.pem -out root.pem

  # generate server and client keys
  openssl genrsa -out key.pem 2048
  openssl genrsa -out client-key.pem 2048

  # generate server and client certificate requests
  openssl req -new -key key.pem -out server.csr
  openssl req -new -key client-key.pem -out client.csr

  # generate server certificate
  mkdir -p demoCA
  touch demoCA/index.txt
  echo 01 > demoCA/serial
  openssl ca -out cert.pem -days 3650 -keyfile root-key.pem -cert root.pem -outdir . -infiles server.csr
  rm -rf demoCA
  rm -f server.csr

  # generate client certificate
  mkdir -p demoCA
  touch demoCA/index.txt
  echo 02 > demoCA/serial
  openssl ca -out client-cert.pem -days 3650 -keyfile root-key.pem -cert root.pem -outdir . -infiles client.csr
  rm -rf demoCA
  rm -f client.csr

  # Merge files and remove temporary ones
  cat client-cert.pem client-key.pem > client.pem
  rm -f client-key.pem client-cert.pem
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
go run ../../../../debian/generate_cert.go --host=localhost --ca
cd ../../../

# Create container
sudo docker login
sudo docker build --rm -t $username/shelter .

# Remove deploy data
rm -fr container

# Push the container to the index
if [ "$2" = "--push" ]; then
  sudo docker push $username/shelter
fi