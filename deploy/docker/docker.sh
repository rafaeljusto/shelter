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

# Create container
docker build -t shelter deploy/docker/Dockerfile 

# Push the container to the index
docker push $username/shelter
