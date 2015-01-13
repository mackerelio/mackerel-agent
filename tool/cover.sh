#/bin/sh
# Calculate coverage
# run this and you can see coverage reports as following
#   % go tool cover -html .profile.cov
set -e

prof=${1:-".profile.cov"}
echo "mode: count" > $prof

cleanup() {
  if [ $tmpprof != "" ] && [ -f $tmpprof ]; then
    rm -f $tmpprof
  fi
  exit
}
trap cleanup INT QUIT TERM EXIT

# cmd/cover doesn't support -coverprofile for multiple packages, so run tests in each dir.
# ref. https://github.com/mattn/goveralls/issues/20
gopath1=$(echo $GOPATH | cut -d: -f1)
for pkg in $(go list ./...); do
  tmpprof=$gopath1/src/$pkg/profile.tmp
  go test -covermode=count -coverprofile=$tmpprof $pkg
  if [ -f $tmpprof ]; then
    cat $tmpprof | tail -n +2 >> $prof
    rm $tmpprof
  fi
done
