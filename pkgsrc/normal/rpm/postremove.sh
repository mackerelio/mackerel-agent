#!/bin/sh

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload
fi

