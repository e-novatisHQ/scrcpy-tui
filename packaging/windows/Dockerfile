FROM debian:trixie-slim@sha256:a99cfc517144bc59b1978475ec53b46ecabec7e43635402ee5b77cc54cd1b20a
RUN apt-get update && apt-get install -y --no-install-recommends \
    wixl=0.106+repack-1 msitools=0.106+repack-1 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /build
ENTRYPOINT ["wixl"]
