FROM alpine:3.22
ARG TARGETOS
ARG TARGETARCH
COPY build/changelog-cli_${TARGETOS}_${TARGETARCH} /usr/local/bin/changelog-cli
RUN chmod +x /usr/local/bin/changelog-cli
ENTRYPOINT ["/usr/local/bin/changelog-cli"]
