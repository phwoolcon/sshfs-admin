# Build
FROM golang:alpine3.10 as builder

ARG UTILS_BASE="https://raw.githubusercontent.com/phwoolcon/docker-utils/master"
WORKDIR /app
RUN wget ${UTILS_BASE}/alpine/pick-mirror -O /usr/local/bin/pick-mirror; \
    chmod +x /usr/local/bin/*; \
    pick-mirror v3.10; \
    apk update; apk upgrade; \
    apk add --no-cache bash coreutils make upx;
COPY . /app
RUN make build;

# Release
FROM alpine:3.10
ARG UTILS_BASE="https://raw.githubusercontent.com/phwoolcon/docker-utils/master"
RUN wget ${UTILS_BASE}/alpine/pick-mirror -O /usr/local/bin/pick-mirror; \
    chmod +x /usr/local/bin/*; \
    pick-mirror v3.10; \
    apk update; apk upgrade; \
    apk add --no-cache bash coreutils openssh;
COPY --from=builder /app/build/sshfs-admin-linux-x64 /app/sshfs-admin-linux-x64
COPY web /app/web
COPY scripts /app/scripts
WORKDIR /app
VOLUME ["/data", "/data/sshfs/etc/root.ssh"]
EXPOSE 8000 8443
ENTRYPOINT ["/app/scripts/entrypoint"]
