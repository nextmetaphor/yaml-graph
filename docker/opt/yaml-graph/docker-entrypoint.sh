#!/usr/bin/env bash

# run local neo4j in the background
nohup /docker-entrypoint.sh neo4j &

# set the prompt
echo 'export PS1="\[\033[1;34m\]yaml-graph \[\033[00m\]$ "' >> /home/ymlgraph/.bashrc