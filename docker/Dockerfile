FROM neo4j:latest
ENV NEO4J_AUTH=none

# copy the required files into the installation directory
COPY opt/yaml-graph/* /opt/yaml-graph/

# create a ymlgraph user and group
RUN groupadd -r ymlgraph && useradd --no-log-init -r -m -g ymlgraph ymlgraph

# change ownership and permissions of yaml-graph installation directory
RUN chmod -R 755 /opt/yaml-graph
RUN chown -R ymlgraph:ymlgraph /opt/yaml-graph

# create symlink in /usr/bin to the yaml-graph binary
RUN ln -s /opt/yaml-graph/yaml-graph /usr/bin/yaml-graph

# as ymlgraph execute the entrypoint script, starting neo4j in the background
USER ymlgraph
ENTRYPOINT ["bash", "-c", "/opt/yaml-graph/docker-entrypoint.sh; cd ~; clear; cat /opt/yaml-graph/logo.txt; bash"]