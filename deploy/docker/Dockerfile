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

# Install MongoDB
RUN apt-get update
RUN apt-get install -y mongodb-10gen openssh-server

# Create the MongoDB and SSH data directory
RUN mkdir -p /data/db
RUN mkdir /var/run/sshd
RUN echo 'root:abc123' | chpasswd

# Copy necessary files to container
ADD container /usr/shelter

# Create project structure for logging
RUN mkdir -p /usr/shelter/var/log

# Expose REST port and web client port to the outside world
EXPOSE 22 4344 4444

# Run services when the container starts
ENTRYPOINT ["usr/shelter/bin/entrypoint.sh"]