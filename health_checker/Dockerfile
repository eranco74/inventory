FROM registry.access.redhat.com/ubi8/ubi-minimal:latest


COPY build/health_checker /usr/local/bin

ENTRYPOINT ["/usr/local/bin/health_checker"]