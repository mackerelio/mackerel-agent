#!/bin/sh
systemctl --no-reload preset mackerel-agent-stage.service
systemctl enable mackerel-agent-stage.service
