convert_for_alternative() {
  target_dir=$1
  agent_name=$2
  if [ "$agent_name" != "mackerel-agent" ]; then
    # rename and replace files for alternative build. ex. mackerel-agent-stage, mackerel-agent-kcps
    for f in $(find $target_dir -type f); do
      cat $f | (rm $f; sed "s/mackerel-agent/$agent_name/g" > $f)
      if expr "$f" : '.*mackerel-agent' > /dev/null; then
        dest=$(echo $f | sed "s/mackerel-agent/$agent_name/")
        mv "$f" "$dest"
      fi
    done
  fi
}
