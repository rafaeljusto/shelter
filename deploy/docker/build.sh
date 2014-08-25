#!/bin/sh

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

: ${GOPATH:?"Need to set GOPATH"}

usage() {
  echo "Usage: $1 <username> [--push]"
  exit 0
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

# <src> must be the path to a file or directory relative
# to the source directory being built (also called the
# context of the build) or a remote file URL.
cd deploy/docker

# Create container
sudo docker build --rm -t $username/shelter .

# Push the container to the index
if [ "$2" = "--push" ]; then
  sudo docker login
  sudo docker push $username/shelter
fi