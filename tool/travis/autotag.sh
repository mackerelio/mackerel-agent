#!/bin/sh
set -ex

echo -e "Host github.com\n\tStrictHostKeyChecking no\nIdentityFile ~/.ssh/deploy.key\n" >> ~/.ssh/config
openssl aes-256-cbc -K $encrypted_e129501e13c5_key -iv $encrypted_e129501e13c5_iv -in tool/travis/mackerel-agent.pem.enc -out ../deploy.key -d
cp ../deploy.key ~/.ssh/
chmod 600 ~/.ssh/deploy.key
git config --global user.email "mackerel-developers@hatena.ne.jp"
git config --global user.name  "mackerel"
git remote set-url origin git@github.com:mackerelio/mackerel-agent-plugins.git
tool/autotag
