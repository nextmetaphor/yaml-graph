FROM neo4j:latest
ENV NEO4J_AUTH=none
COPY src/yaml-graph /opt/yaml-graph/yaml-graph
COPY logo.txt /opt/yaml-graph/logo.txt
COPY docker-entrypoint.sh /opt/yaml-graph/docker-entrypoint.sh
RUN groupadd -r ymlgraph && useradd --no-log-init -r -m -g ymlgraph ymlgraph
RUN chmod 755 /opt/yaml-graph/docker-entrypoint.sh
RUN ln -s /opt/yaml-graph/yaml-graph /usr/bin/yaml-graph
USER ymlgraph
ENTRYPOINT ["bash", "-c", "/opt/yaml-graph/docker-entrypoint.sh; cd ~; clear; cat /opt/yaml-graph/logo.txt; bash"]