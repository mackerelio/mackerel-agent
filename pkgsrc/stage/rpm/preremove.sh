#!/bin/sh

if [ -d /run/systemd/system ]; then
    systemctl --no-reload disable --now --no-warn mackerel-agent-stage.service
else
    systemctl --no-reload disable --no-warn mackerel-agent-stage.service
fi
