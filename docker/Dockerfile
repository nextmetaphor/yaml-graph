FROM neo4j:latest

ENV NEO4J_AUTH=none

# create a ymlgraph user and group
RUN useradd --no-log-init -r -m ymlgraph

# copy the required files into the installation directory
COPY --chown=ymlgraph:ymlgraph ./utils/* /opt/yaml-graph/

# create symlink in /usr/bin to the yaml-graph binary
RUN ln -s /opt/yaml-graph/yaml-graph /usr/bin/yaml-graph

# as ymlgraph execute the entrypoint script, starting neo4j in the background
USER ymlgraph

WORKDIR /home/ymlgraph

ENTRYPOINT ["tini", "--", "/opt/yaml-graph/docker-entrypoint.sh"]