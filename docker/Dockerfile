FROM neo4j:5.5.0-community

ENV NEO4J_AUTH=none

# create a ymlgraph user and group
RUN useradd --no-log-init -r -m ymlgraph && \
    echo 'export PS1="\[\033[1;34m\]yaml-graph \[\033[00m\]$ "' >> /home/ymlgraph/.bashrc

# copy the required files into the installation directory
COPY --chown=ymlgraph:ymlgraph ./utils/* /opt/yaml-graph/

# create symlink in /usr/bin to the yaml-graph binary
RUN ln -s /opt/yaml-graph/yaml-graph /usr/bin/yaml-graph && chmod u+x /opt/yaml-graph/docker-entrypoint.sh

USER ymlgraph

WORKDIR /home/ymlgraph

# as ymlgraph execute the entrypoint script, starting neo4j in the background
ENTRYPOINT ["bash", "-c", "/opt/yaml-graph/docker-entrypoint.sh; bash"]