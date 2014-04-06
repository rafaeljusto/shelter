#!/bin/sh

sudo docker run -d -P --name shelter rafaeljusto/shelter
echo "REST server port:"
sudo docker port shelter 4344
echo "Web Client port:"
sudo docker port shelter 4444
