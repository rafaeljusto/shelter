#!/bin/sh

sudo docker run -d -P rafaeljusto/shelter

sleep 1

# For this to work we assume that the running container is the last one from the PS list
container=`sudo docker ps | grep -v CONTAINER | awk 'END{ print $1 }'`

echo -n "SSH server: "
sudo docker port $container 22

echo -n "REST server: "
sudo docker port $container 4443

echo -n "Web Client: "
sudo docker port $container 4444
