FROM confluentinc/cp-kafka:7.2.1

USER root

WORKDIR /etc/kafka

RUN dnf install git --assumeyes

# Download JMX exporter
RUN wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.15.0/jmx_prometheus_javaagent-0.15.0.jar

COPY ./metrics ./metrics

CMD ["/etc/confluent/docker/run"]
