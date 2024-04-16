FROM ghcr.io/reubenmiller/c8y-shell
ARG SSH_PRIVATE_KEY
USER root
RUN apk update && apk add openssh openssl --no-cache \
    # add user to sudoers group
    && echo '%wheel ALL=(ALL) ALL' > /etc/sudoers.d/wheel \
    && adduser c8yuser wheel

USER c8yuser
COPY --chown=c8yuser . /home/c8yuser/c8y-tedge
RUN c8y extensions install /home/c8yuser/c8y-tedge \
    && mkdir -p /home/c8yuser/.ssh \
    && chmod 700 /home/c8yuser/.ssh \
    && echo "${SSH_PRIVATE_KEY}" | base64 -d > /home/c8yuser/.ssh/id_ed25519 \
    && chmod 600 /home/c8yuser/.ssh/id_ed25519 \
    && echo "eval \$(ssh-agent -s)" >> /home/c8yuser/.zshrc \
    && echo "ssh-add /home/c8yuser/.ssh/id_ed25519" >> /home/c8yuser/.zshrc \
    && printf "Host *\n    StrictHostKeyChecking no" > /home/c8yuser/.ssh/config \
    && chmod 600 /home/c8yuser/.ssh/config

ENV OPEN_WEBSITE=0
ENV CI=true

ENTRYPOINT [ "/bin/sleep", "infinity" ]
