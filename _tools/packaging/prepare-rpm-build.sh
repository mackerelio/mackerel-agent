#!/bin/sh
set -e
set -x

pwd=`dirname $0`
. "$pwd/common.sh"

MACKEREL_AGENT_NAME=${MACKEREL_AGENT_NAME:-mackerel-agent}
spec_filename="mackerel-agent.spec"
if [ "$BUILD_SYSTEMD" != "" ]; then
    spec_filename="mackerel-agent-systemd.spec"
fi

orig_dir="packaging/rpm"
build_dir="packaging/rpm-build"

mkdir -p rpmbuild/RPMS/{noarch,x86_64}

rm -rf "$build_dir"
mkdir -p "$build_dir"
cp -r "$orig_dir/src" "$build_dir/src"
cp mackerel-agent.sample.conf "$build_dir/src/mackerel-agent.conf"
cp "$orig_dir/$spec_filename" "$build_dir/mackerel-agent.spec"

convert_for_alternative $build_dir $MACKEREL_AGENT_NAME
