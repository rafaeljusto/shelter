# Copyright 2014 Rafael Dantas Justo. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

# Shelter
#
# VERSION 1.0

# use the ubuntu base image provided by dotCloud
FROM    ubuntu:latest

MAINTAINER Rafael Dantas Justo, adm@rafael.net.br

# Ensure that the ubuntu package repository is up to date
RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get update

# Add 10gen official apt source to the sources list
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
RUN echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | tee /etc/apt/sources.list.d/10gen.list

# Hack for initctl not being available in Ubuntu
#RUN dpkg-divert --local --rename --add /sbin/initctl
#RUN ln -s /bin/true /sbin/initctl

# Install MongoDB
RUN apt-get update
RUN apt-get install -y mongodb-10gen

# Create the MongoDB data directory
RUN mkdir -p /data/db

# Create project structure
RUN mkdir -p /usr/shelter/bin
RUN mkdir -p /usr/shelter/templates
RUN mkdir -p /usr/shelter/var/log
RUN mkdir -p /usr/shelter/etc/keys

# Download the necessary files
ADD container/bin/shelter /usr/shelter/bin/shelter
ADD container/bin/container-entrypoint.sh /usr/shelter/bin/container-entrypoint.sh
ADD container/templates /usr/shelter/templates
ADD container/etc/shelter.conf /usr/shelter/etc/shelter.conf
ADD container/etc/messages.conf /usr/shelter/etc/messages.conf
ADD container/etc/keys/key.pem /usr/shelter/etc/keys/key.pem
ADD container/etc/keys/cert.pem /usr/shelter/etc/keys/cert.pem

# Expose REST port and web client port to the outside world
EXPOSE 4443
EXPOSE 4444

# Run services when the container starts
ENTRYPOINT ["usr/shelter/bin/container-entrypoint.sh"]
