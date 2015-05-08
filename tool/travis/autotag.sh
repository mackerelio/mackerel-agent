#!/bin/sh
set -ex

deploykey=~/.ssh/deploy.key

echo "
Host github.com
    StrictHostKeyChecking no
    IdentityFile $deploykey
" >> ~/.ssh/config
openssl aes-256-cbc -K $encrypted_e129501e13c5_key -iv $encrypted_e129501e13c5_iv -in tool/travis/mackerel-agent.pem.enc -out $deploykey -d
chmod 600 $deploykey
git config --global user.email "mackerel-developers@hatena.ne.jp"
git config --global user.name  "mackerel"
git remote set-url origin git@github.com:mackerelio/mackerel-agent-plugins.git
tool/autotag
