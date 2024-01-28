#!/bin/sh
set -e

# TODO: This whole script could be replaced by creating a new tedge cert
# command to create a csr and device key
# Example: tedge create --csr device.csr --device-id gateway_abcdef

DEVICE_ID="$1"
DEVICE_KEY_PATH=$(tedge config get device.key_path)

SUDO=
if [ "$(id -u)" != 0 ]; then
    if command -V sudo >/dev/null 2>&1; then
        SUDO=sudo
    fi
fi

if ! command -V openssl >/dev/null 2>&1; then
    echo "Missing dependency: openssl" >&2
    exit 1
fi

if [ ! -f "$DEVICE_KEY_PATH" ]; then
    echo "Creating device private key" >&2
    $SUDO openssl genrsa -out "$DEVICE_KEY_PATH" 2048
    $SUDO chown mosquitto:root "$DEVICE_KEY_PATH"
    $SUDO chmod 600 "$DEVICE_KEY_PATH"
else
    echo "Using existing private key" >&2
fi

# Create Certificate Signing Request (CSR)
DEVICE_CSR=$(
    $SUDO openssl req \
        -key "$DEVICE_KEY_PATH" \
        -new \
        -subj "/O=thin-edge/OU=Test\ Device/CN=${DEVICE_ID}"
)

echo "$DEVICE_CSR"
printf '{"csr":"%s"}' "$DEVICE_CSR"
