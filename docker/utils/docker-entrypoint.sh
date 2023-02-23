#!/bin/bash -eu

# set the prompt
echo 'export PS1="\[\033[1;34m\]yaml-graph \[\033[00m\]$ "' >> /home/ymlgraph/.bashrc

cat /opt/yaml-graph/logo.txt; yaml-graph version; echo;

# run local neo4j
/startup/docker-entrypoint.sh neo4j