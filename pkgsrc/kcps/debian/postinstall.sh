#!/bin/sh

set -e
 
case "$1" in
 configure)
   if [ ! -e /etc/mackerel-agent-kcps/mackerel-agent-kcps.conf ]; then
     cp /etc/mackerel-agent-kcps/mackerel-agent-kcps.conf.example /etc/mackerel-agent-kcps/mackerel-agent-kcps.conf
   fi

   if [ -d /run/systemd/system ]; then
     systemctl --system daemon-reload
   fi
#  systemctl enable mackerel-agent-kcps.service
 ;;
 abort-upgrade|abort-remove|abort-deconfigure)
   exit 0
 ;;
 *)
   echo "postinst called with unknown argument \`$1'" >&2
   exit 1
 ;;
esac

