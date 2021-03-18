#!/bin/bash

set -eu

if [ "$(id -u)" -ne 0 ]; then
    echo "Please run with sudo"
    exit 1
fi

service_name=$1
if [ -z "$service_name" ]; then
    echo "Usage: $0 SERVICE-NAME"    
    exit 1
fi

cp "$PWD/$service_name" "/etc/systemd/system/$service_name"

systemctl daemon-reload

systemctl start "$service_name"
systemctl enable "$service_name"
