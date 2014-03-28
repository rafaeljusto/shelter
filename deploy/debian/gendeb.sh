#!/bin/sh

# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

pack_name="shelter"
version="0.1"
vendor="Rafael Dantas Justo"
maintainer="Rafael Dantas Justo <adm@rafael.net.br>"
url="http://github.com/rafaeljusto/shelter"
license="GPL"
description="System that checks periodically DNS servers for DNS and DNSSEC misconfigurations"

install_path=/usr/local/shelter
tmp_dir=/tmp/shelter
project_root=$tmp_dir$install_path

workspace=`echo $GOPATH | cut -d: -f1`
workspace=$workspace/src/github.com/rafaeljusto/shelter

# recompiling everything
current_dir=`pwd`
cd $workspace
go build shelter.go
cd $current_dir

if [ -f $pack_name*.deb ]; then
  # remove old deb
  rm $pack_name*.deb
fi

if [ -d $tmp_dir ]; then
  rm -rf $tmp_dir
fi

mkdir -p $tmp_dir$install_path
mkdir -p $tmp_dir$install_path/bin
mkdir -p $tmp_dir$install_path/var/log

cp -r $workspace/etc $project_root/
cp -r $workspace/templates $project_root/
mv $workspace/shelter $project_root/bin/

fpm -s dir -t deb \
  --exclude=.git -n $pack_name -v $version --vendor "$vendor" \
  --maintainer "$maintainer" --url $url --license $license --description "$description" \
  --config-files usr/local/shelter/etc/shelter.conf \
  --deb-upstart $workspace/deploy/debian/shelter.upstart \
  --deb-user root --deb-group root \
  --prefix / -C $tmp_dir usr/local/shelter
