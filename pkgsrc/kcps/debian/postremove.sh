#!/bin/sh
set -e

if [ "$1" = "purge" ]; then
  rm -f /etc/mackerel-agent/mackerel-agent-kcps.conf
  rm -f /var/lib/mackerel-agent-kcps/id
fi
rm -f /var/run/mackerel-agent-kcps.pid

if [ -d /run/systemd/system ]; then
  systemctl --system daemon-reload
fi
