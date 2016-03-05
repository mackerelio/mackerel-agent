convert_for_alternative() {
  target_dir=$1
  agent_name=$2
  if [ "$agent_name" != "mackerel-agent" ]; then
    # rename and replace files for alternative build. ex. mackerel-agent-stage, mackerel-agent-kcps
    for filename in $(find $target_dir -type f); do
      perl -i -pe "s/mackerel-agent/$agent_name/g" $filename
      if expr "$filename" : '.*mackerel-agent' > /dev/null; then
        destfile=$(echo $filename | sed "s/mackerel-agent/$agent_name/")
        mv "$filename" "$destfile"
      fi
    done
  fi
}
