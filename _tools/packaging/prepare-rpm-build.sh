#!/bin/sh
set -e
set -x

pwd=`dirname $0`
. "$pwd/common.sh"

MACKEREL_AGENT_NAME=${MACKEREL_AGENT_NAME:-mackerel-agent}

orig_dir="packaging/rpm"
build_dir="packaging/rpm-build"

cp mackerel-agent.sample.conf "$orig_dir/src/mackerel-agent.conf"
rm -rf "$build_dir"
cp -r "$orig_dir" "$build_dir"

convert_for_alternative $build_dir $MACKEREL_AGENT_NAME
