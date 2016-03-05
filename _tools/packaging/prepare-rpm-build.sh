#!/bin/sh
set -e
set -x

MACKEREL_AGENT_NAME=${MACKEREL_AGENT_NAME:-mackerel-agent}

orig_dir="packaging/rpm"
build_dir="packaging/rpm-build"

cp mackerel-agent.sample.conf "$orig_dir/src/mackerel-agent.conf"
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
