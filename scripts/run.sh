#!/bin/bash

set -eu

docker run --init \
    --detach \
    --restart unless-stopped \
    -e NEATO_ROBOT_SERIALNUMBER=123123 \
    -e NEATO_ROBOT_SECRET=12334
    -p 8102:8080 \
    --name neato-http\
    mrbuk/neato-http:0.5
