#!/bin/sh
set -e

if [ "$1" = "remove" ] && [ -d /run/systemd/system ]; then
    systemctl stop mackerel-agent-kcps.service
fi
