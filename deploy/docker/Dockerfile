# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

# Shelter
#
# VERSION 2.0

# use the ubuntu base image provided by dotCloud
FROM    ubuntu:latest

MAINTAINER Rafael Dantas Justo, adm@rafael.net.br

# Ensure that the ubuntu package repository is up to date
RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN echo "deb http://dl.bintray.com/rafaeljusto/deb ./" >> /etc/apt/sources.list

# Add 10gen official apt source to the sources list
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | tee /etc/apt/sources.list.d/10gen.list

# Install MongoDB, SSH server and Shelter
RUN apt-get update
ENV "DEBIAN_FRONTEND" "noninteractive"
RUN apt-get install -y --force-yes mongodb-10gen openssh-server shelter
ENV "DEBIAN_FRONTEND" ""

# Create the MongoDB, SSH data directory and Shelter configuration file
RUN mkdir -p /data/db
RUN mkdir /var/run/sshd
RUN echo 'root:abc123' | chpasswd
RUN cp /usr/shelter/etc/shelter.conf.unix.sample /usr/shelter/etc/shelter.conf

# Copy necessary files to container
ADD entrypoint.sh /usr/shelter/bin/entrypoint.sh

# Expose REST port and web client port to the outside world
EXPOSE 22 4443 4444

# Run services when the container starts
ENTRYPOINT ["usr/shelter/bin/entrypoint.sh"]
