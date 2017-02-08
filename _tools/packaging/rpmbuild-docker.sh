#!/bin/sh

set -e
set -x

docker run --rm -v "$PWD":/workspace -v "$PWD/rpmbuild":/rpmbuild astj-mackerel-packager-beta:$RPMBUILD_DOCKER_TAG $@
