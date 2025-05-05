FROM docker.elastic.co/beats/filebeat:8.18.0

# Copy config
COPY filebeat.yml /usr/share/filebeat/filebeat.yml

USER root
RUN chown root:filebeat /usr/share/filebeat/filebeat.yml

ENTRYPOINT ["/usr/share/filebeat/filebeat", "-e", "--strict.perms=false"]