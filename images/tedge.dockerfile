FROM ghcr.io/thin-edge/tedge-demo-main-systemd:20240319.1412
ARG SSH_PUBLIC_KEY
RUN apt-get install iproute2 \
    && systemctl enable ssh \
    && echo "${SSH_PUBLIC_KEY}" | base64 -d >> "/root/.ssh/authorized_keys"
