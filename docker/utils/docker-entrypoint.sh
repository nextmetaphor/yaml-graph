#!/usr/bin/env bash

cd ~; clear; cat /opt/yaml-graph/logo.txt; yaml-graph version; echo;

# run local neo4j in the background
nohup /startup/docker-entrypoint.sh neo4j &> /logs/startup.out &
