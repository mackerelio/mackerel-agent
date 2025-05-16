#!/bin/sh
systemctl --no-reload preset mackerel-agent-kcps.service
systemctl enable mackerel-agent-kcps.service
