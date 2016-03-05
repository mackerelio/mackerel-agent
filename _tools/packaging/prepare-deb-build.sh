#!/bin/sh
set -e
set -x

MACKEREL_AGENT_NAME=${MACKEREL_AGENT_NAME:-mackerel-agent}
MACKEREL_AGENT_VERSION=${MACKEREL_AGENT_VERSION:-0.0.0}

orig_dir="packaging/deb"
build_dir="packaging/deb-build"

cp mackerel-agent.sample.conf   "$orig_dir/debian/mackerel-agent.conf"
rm -rf "$build_dir"
cp -r "$orig_dir" "$build_dir"

if [ "$MACKEREL_AGENT_NAME" != "mackerel-agent" ]; then
  for filename in $(find $build_dir -type f); do
    perl -i -pe "s/mackerel-agent/$MACKEREL_AGENT_NAME/g" $filename
    if expr "$filename" : '.*mackerel-agent' > /dev/null; then
      destfile=$(echo $filename | sed "s/mackerel-agent/$MACKEREL_AGENT_NAME/")
      mv "$filename" "$destfile"
    fi
  done
fi
cp "build/$MACKEREL_AGENT_NAME" "$build_dir/debian/$MACKEREL_AGENT_NAME.bin"
cp packaging/dummy-empty.tar.gz "packaging/${MACKEREL_AGENT_NAME}_$MACKEREL_AGENT_VERSION.orig.tar.gz"
