FROM ubuntu:26.04

ARG UBUNTU_MIRROR=http://archive.ubuntu.com/ubuntu
ARG IMAGE_VERSION=dev
ENV DEBIAN_FRONTEND=noninteractive \
    HOME=/home/pasturestack

RUN set -eux; \
    find /etc/apt -type f \( -name '*.list' -o -name '*.sources' \) -exec sed -i \
        -e "s|http://security.ubuntu.com/ubuntu|${UBUNTU_MIRROR}|g" \
        -e "s|http://archive.ubuntu.com/ubuntu|${UBUNTU_MIRROR}|g" {} +; \
    printf 'Acquire::Retries "5";\nAcquire::http::Timeout "30";\nAcquire::https::Timeout "30";\nAcquire::http::Pipeline-Depth "0";\n' > /etc/apt/apt.conf.d/80pasturestack-retries; \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        openssh-client && \
    rm -rf /var/lib/apt/lists/* /usr/bin/pebble

COPY pasturestack /usr/bin/
COPY LICENSE SECURITY.md ORIGIN.md COMPATIBILITY.md /licenses/
RUN groupadd --system --gid 10001 pasturestack && \
    useradd --system --uid 10001 --gid pasturestack --home-dir /home/pasturestack --create-home --shell /usr/sbin/nologin pasturestack && \
    chmod +x /usr/bin/pasturestack && \
    chown -R pasturestack:pasturestack /home/pasturestack /licenses

LABEL org.opencontainers.image.title="pasturestack-cli" \
      org.opencontainers.image.description="PastureStack compatible control-platform CLI" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.source="https://github.com/PastureStack/cli" \
      org.opencontainers.image.version="${IMAGE_VERSION}"

WORKDIR /home/pasturestack
USER 10001:10001
ENTRYPOINT ["pasturestack"]
CMD  ["--help"]
