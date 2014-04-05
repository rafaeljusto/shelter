#!/bin/sh

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

: ${GOPATH:?"Need to set GOPATH"}

usage() {
  echo "Usage: $1 <username>"
  exit 0
}

username=$1
if [ -z "$username" ]; then
  usage $0
fi

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
cp container-entrypoint.sh container/bin/
cp ../../etc/shelter.conf.sample container/etc/shelter.conf
cp ../../etc/messages.conf container/etc/
cp -r ../../templates container/


# Generate certificates for container
cd container/etc/keys
go run $workspace/src/github.com/rafaeljusto/shelter/deploy/debian/generate_cert.go --host=localhost
cd ../../../

# Create container
sudo docker build -t shelter . 

# Remove deploy data 
rm -fr container

# Push the container to the index
docker push $username/shelter
