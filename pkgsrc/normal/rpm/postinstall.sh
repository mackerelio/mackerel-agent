#!/bin/sh
systemctl --no-reload preset mackerel-agent.service
systemctl enable mackerel-agent.service
