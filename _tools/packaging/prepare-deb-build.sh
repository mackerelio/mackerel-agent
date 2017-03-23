#!/bin/sh
set -e
set -x

pwd=`dirname $0`
. "$pwd/common.sh"

MACKEREL_AGENT_NAME=${MACKEREL_AGENT_NAME:-mackerel-agent}
BUILD_DIRECTORY=${BUILD_DIRECTORY:-build}

orig_dir="packaging/deb"
build_dir="packaging/deb-build"

if [ "$BUILD_SYSTEMD" != "" ]; then
    orig_dir="packaging/deb-systemd"
fi

MACKEREL_AGENT_VERSION=$(grep -o -e "[0-9]\+.[0-9]\+.[0-9]\+-[0-9]" "$orig_dir/debian/changelog" | head -1 | sed 's/-.*$//')

rm -rf "$build_dir"
cp -r "$orig_dir" "$build_dir"
cp mackerel-agent.sample.conf   "$build_dir/debian/mackerel-agent.conf"

convert_for_alternative $build_dir $MACKEREL_AGENT_NAME
cp "${BUILD_DIRECTORY}/$MACKEREL_AGENT_NAME" "$build_dir/debian/$MACKEREL_AGENT_NAME.bin"
cp packaging/dummy-empty.tar.gz "packaging/${MACKEREL_AGENT_NAME}_$MACKEREL_AGENT_VERSION.orig.tar.gz"
