#/bin/sh
prof=${1:-".profile.cov"}
echo "mode: count" > $prof

cleanup() {
  if [ $tmpprof != "" && -f $tmpprof ]; then
    rm -f $tmpprof
  fi
  exit
}
trap cleanup INT QUIT TERM EXIT

for dir in $(find . -not -path '*/[._]*' -type d); do
  if ls $dir/*.go > /dev/null 2>&1; then
    tmpprof=$dir/profile.tmp
    go test -covermode=count -coverprofile=$tmpprof $dir
    if [ -f $tmpprof ]; then
      cat $tmpprof | tail -n +2 >> $prof
      rm $tmpprof
    fi
  fi
done
