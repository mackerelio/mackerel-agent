#!/bin/sh
set -e

if [ "$1" = "purge" ]; then
  rm -f /etc/mackerel-agent/mackerel-agent.conf
  rm -f /var/lib/mackerel-agent/id
fi
rm -f /var/run/mackerel-agent.pid

if [ -d /run/systemd/system ]; then
  systemctl --system daemon-reload
fi
