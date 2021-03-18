#!/bin/bash
set -eu

service_name=$1

echo "stop the service to avoid systemd trying it to restart"
systemctl start "$service_name"
echo

echo "stop again using docker"
docker stop "$service_name"
echo

echo "delete the current $service_name image"
docker rm "$service_name"
echo 

echo "install new controller image"
./run.sh
echo

echo "start the service to avoid systemd trying it to restart"
systemctl start "$service_name"
echo

echo "check that controller is running" 
docker ps
echo
