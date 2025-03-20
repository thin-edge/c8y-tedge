// Template to a random device with realistic thin-edge.io meta data
//
// Usage:
//   # Create 10 random devices, replace the "10" with however many devices you want to create
//   c8y util repeat 10 | c8y devices create --template ./device.jsonnet
//
//   # Delete the devices
//   c8y devices list --query "has(c8y_IsDummy)" --includeAll | c8y devices delete
//
local nameSuffix = _.Hex(12);

# List of realistic devices, where one will be selected (later on)
local devices = [
    //
    // Raspberry Pi 4
    //
    {
        "name": "rpi4-" + nameSuffix,
        "c8y_Hardware": {
            "model": "Raspberry Pi 4 Model B Rev 1.1",
            "revision": "c03111",
            "serialNumber": "10000000" + _.Hex(8),
        },
        "device_OS": {
            "arch": "aarch64",
            "displayName": "Poky (Yocto Project Reference Distro) %s (kirkstone)" % $.device_OS.version,
            "family": "GNU/Linux",
            "hostname": $.name,
            "kernel": "#1 SMP PREEMPT Tue Mar 19 17:41:59 UTC 2024",
            "version": "4.0.25"
        },
        // Change the agent version to keep it realistic (but keep the any other static fields)
        "c8y_Agent"+: {
            "version": "1.3.1"
        },
    },

    //
    // Raspberry Pi 5
    //
    {
        "name": "rpi5-" + nameSuffix,
        "c8y_Hardware": {
            "model": "Raspberry Pi 5 Model B Rev 1.0",
            "revision": "d04170",
            "serialNumber": _.Hex(16),
        },
        "device_OS": {
            "arch": "aarch64",
            "displayName": "Poky (Yocto Project Reference Distro) %s (scarthgap)" % $.device_OS.version,
            "family": "GNU/Linux",
            "hostname": $.name,
            "kernel": "#1 SMP PREEMPT Tue Mar 19 17:41:59 UTC 2024",
            "version": "5.0.7"
        },
        "c8y_Agent"+: {
            "version": "1.4.2"
        },
    },

    //
    // StarFive VisionFive 2
    //
    {
        "name": "starfive-" + nameSuffix,
        "c8y_Hardware": {
            "model": "StarFive VisionFive 2",
            "revision": "a01882",
            "serialNumber": _.Hex(16),
        },
        "device_OS": {
            "arch": "riscv4",
            "displayName": "Poky (Yocto Project Reference Distro) %s (scarthgap)" % $.device_OS.version,
            "family": "GNU/Linux",
            "hostname": $.name,
            "kernel": "#1 SMP PREEMPT Tue Mar 19 17:41:59 UTC 2024",
            "version": "5.0.7"
        },
        "c8y_Agent"+: {
            "version": "1.4.2"
        },
    },
    //
    // Jetson Nano
    //
    {
        "name": "jetson-" + nameSuffix,
        "c8y_Hardware": {
            "model": "Jetson Nano",
            "revision": "18d7d",
            "serialNumber": _.Hex(16),
        },
        "device_OS": {
            "arch": "aarch64",
            "displayName": "Poky (Yocto Project Reference Distro) %s (scarthgap)" % $.device_OS.version,
            "family": "GNU/Linux",
            "hostname": $.name,
            "kernel": "#1 SMP PREEMPT Tue Mar 19 17:41:59 UTC 2024",
            "version": "5.0.7"
        },
        "c8y_Agent"+: {
            "version": "1.3.1"
        },
    },
];

# Randomly select one of the above device types
local device = devices[_.Int(0, std.length(devices))];

# Static values (then the device Value is merged at the end)
{
    // Marker to indicate that these devices are fake
    "c8y_IsDummy": {},

    "c8y_Agent": {
        "name": "thin-edge.io",
        "url": "https://thin-edge.io",
        "version": "1.4.2"
    },
    "c8y_Firmware": {
        "name": "core-image-tedge-rauc",
        "version": "20250305.1329"
    },
    "c8y_IsDevice": {},
    "c8y_RequiredAvailability": {
        "responseInterval": 1440
    },
    "c8y_SoftwareList": [
        {
            "name": "apt",
            "softwareType": "apt",
            "url": "",
            "version": "2.6.1-r0"
        },
        {
            "name": "avahi-daemon",
            "softwareType": "apt",
            "url": "",
            "version": "0.8-r0"
        },
        {
            "name": "base-files",
            "softwareType": "apt",
            "url": "",
            "version": "3.0.14-r0"
        },
        {
            "name": "base-passwd",
            "softwareType": "apt",
            "url": "",
            "version": "3.6.3-r0"
        },
        {
            "name": "bash",
            "softwareType": "apt",
            "url": "",
            "version": "5.2.21-r0"
        },
        {
            "name": "bluez-firmware-rpidistro-bcm43430a1-hcd",
            "softwareType": "apt",
            "url": "",
            "version": "1.2-9+rpt30+78d6a07730-r0"
        },
        {
            "name": "bluez-firmware-rpidistro-bcm43430b0-hcd",
            "softwareType": "apt",
            "url": "",
            "version": "1.2-9+rpt30+78d6a07730-r0"
        },
        {
            "name": "bluez-firmware-rpidistro-bcm4345c0-hcd",
            "softwareType": "apt",
            "url": "",
            "version": "1.2-9+rpt30+78d6a07730-r0"
        },
        {
            "name": "bluez-firmware-rpidistro-bcm4345c5-hcd",
            "softwareType": "apt",
            "url": "",
            "version": "1.2-9+rpt30+78d6a07730-r0"
        },
        {
            "name": "bluez-firmware-rpidistro-cypress-license",
            "softwareType": "apt",
            "url": "",
            "version": "1.2-9+rpt30+78d6a07730-r0"
        },
        {
            "name": "bluez5",
            "softwareType": "apt",
            "url": "",
            "version": "5.72-r0"
        },
        {
            "name": "bridge-utils",
            "softwareType": "apt",
            "url": "",
            "version": "1.7.1-r0"
        },
        {
            "name": "busybox",
            "softwareType": "apt",
            "url": "",
            "version": "1.36.1-r0"
        },
        {
            "name": "busybox-syslog",
            "softwareType": "apt",
            "url": "",
            "version": "1.36.1-r0"
        },
        {
            "name": "busybox-udhcpc",
            "softwareType": "apt",
            "url": "",
            "version": "1.36.1-r0"
        },
        {
            "name": "ca-certificates",
            "softwareType": "apt",
            "url": "",
            "version": "20211016-r0"
        },
        {
            "name": "cjson",
            "softwareType": "apt",
            "url": "",
            "version": "1.7.18-r0"
        },
        {
            "name": "collectd",
            "softwareType": "apt",
            "url": "",
            "version": "5.12.0-r0"
        },
        {
            "name": "containerd-opencontainers",
            "softwareType": "apt",
            "url": "",
            "version": "v2.0.0-beta.0+git0+b1624c3628-r0"
        },
        {
            "name": "db",
            "softwareType": "apt",
            "url": "",
            "version": "1:5.3.28-r0"
        },
        {
            "name": "dbus-1",
            "softwareType": "apt",
            "url": "",
            "version": "1.14.10-r0"
        },
        {
            "name": "dbus-common",
            "softwareType": "apt",
            "url": "",
            "version": "1.14.10-r0"
        },
        {
            "name": "dbus-tools",
            "softwareType": "apt",
            "url": "",
            "version": "1.14.10-r0"
        },
        {
            "name": "dnsmasq",
            "softwareType": "apt",
            "url": "",
            "version": "2.90-r0"
        },
        {
            "name": "docker-compose",
            "softwareType": "apt",
            "url": "",
            "version": "v2.26.0-r0"
        },
        {
            "name": "docker-moby",
            "softwareType": "apt",
            "url": "",
            "version": "25.0.3+gitf417435e5f6216828dec57958c490c4f8bae4f980+f417435e5f_67e0588f1d-r0"
        },
        {
            "name": "docker-moby-cli",
            "softwareType": "apt",
            "url": "",
            "version": "25.0.3+gitf417435e5f6216828dec57958c490c4f8bae4f980+f417435e5f_67e0588f1d-r0"
        },
        {
            "name": "dosfstools",
            "softwareType": "apt",
            "url": "",
            "version": "4.2-r0"
        },
        {
            "name": "dpkg",
            "softwareType": "apt",
            "url": "",
            "version": "1.22.0-r0"
        },
        {
            "name": "dpkg-start-stop",
            "softwareType": "apt",
            "url": "",
            "version": "1.22.0-r0"
        },
        {
            "name": "dropbear",
            "softwareType": "apt",
            "url": "",
            "version": "2022.83-r0"
        },
        {
            "name": "e2fsprogs",
            "softwareType": "apt",
            "url": "",
            "version": "1.47.0-r0"
        },
        {
            "name": "e2fsprogs-badblocks",
            "softwareType": "apt",
            "url": "",
            "version": "1.47.0-r0"
        },
        {
            "name": "e2fsprogs-dumpe2fs",
            "softwareType": "apt",
            "url": "",
            "version": "1.47.0-r0"
        },
        {
            "name": "e2fsprogs-e2fsck",
            "softwareType": "apt",
            "url": "",
            "version": "1.47.0-r0"
        },
        {
            "name": "e2fsprogs-mke2fs",
            "softwareType": "apt",
            "url": "",
            "version": "1.47.0-r0"
        },
        {
            "name": "gnupg",
            "softwareType": "apt",
            "url": "",
            "version": "2.4.5-r0"
        },
        {
            "name": "gnupg-gpg",
            "softwareType": "apt",
            "url": "",
            "version": "2.4.5-r0"
        },
        {
            "name": "hdparm",
            "softwareType": "apt",
            "url": "",
            "version": "9.65-r0"
        },
        {
            "name": "i2c-tools",
            "softwareType": "apt",
            "url": "",
            "version": "4.3-r0"
        },
        {
            "name": "iptables",
            "softwareType": "apt",
            "url": "",
            "version": "1.8.10-r0"
        },
        {
            "name": "iw",
            "softwareType": "apt",
            "url": "",
            "version": "6.7-r0"
        },
        {
            "name": "jq",
            "softwareType": "apt",
            "url": "",
            "version": "1.7.1-r0"
        },
        {
            "name": "kbd",
            "softwareType": "apt",
            "url": "",
            "version": "2.6.4-r0"
        },
        {
            "name": "kbd-consolefonts",
            "softwareType": "apt",
            "url": "",
            "version": "2.6.4-r0"
        },
        {
            "name": "kbd-keymaps",
            "softwareType": "apt",
            "url": "",
            "version": "2.6.4-r0"
        },
        {
            "name": "kbd-keymaps-pine",
            "softwareType": "apt",
            "url": "",
            "version": "2.6.4-r0"
        },
        {
            "name": "keymaps",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "kmod",
            "softwareType": "apt",
            "url": "",
            "version": "31-r0"
        },
        {
            "name": "ldconfig",
            "softwareType": "apt",
            "url": "",
            "version": "2.39+git0+662516aca8-r0"
        },
        {
            "name": "linux-firmware-rpidistro-bcm43430",
            "softwareType": "apt",
            "url": "",
            "version": "20230625-2+rpt30+4b356e134e-r0"
        },
        {
            "name": "linux-firmware-rpidistro-bcm43436",
            "softwareType": "apt",
            "url": "",
            "version": "20230625-2+rpt30+4b356e134e-r0"
        },
        {
            "name": "linux-firmware-rpidistro-bcm43436s",
            "softwareType": "apt",
            "url": "",
            "version": "20230625-2+rpt30+4b356e134e-r0"
        },
        {
            "name": "linux-firmware-rpidistro-bcm43455",
            "softwareType": "apt",
            "url": "",
            "version": "20230625-2+rpt30+4b356e134e-r0"
        },
        {
            "name": "linux-firmware-rpidistro-bcm43456",
            "softwareType": "apt",
            "url": "",
            "version": "20230625-2+rpt30+4b356e134e-r0"
        },
        {
            "name": "linux-firmware-rpidistro-license",
            "softwareType": "apt",
            "url": "",
            "version": "20230625-2+rpt30+4b356e134e-r0"
        },
        {
            "name": "lz4",
            "softwareType": "apt",
            "url": "",
            "version": "1:1.9.4-r0"
        },
        {
            "name": "mobile-broadband-provider-info",
            "softwareType": "apt",
            "url": "",
            "version": "1:20240407-r0"
        },
        {
            "name": "monit",
            "softwareType": "apt",
            "url": "",
            "version": "5.33.0-r0"
        },
        {
            "name": "mosquitto",
            "softwareType": "apt",
            "url": "",
            "version": "2.0.20-r0"
        },
        {
            "name": "ncurses-terminfo-base",
            "softwareType": "apt",
            "url": "",
            "version": "6.4-r0"
        },
        {
            "name": "neard",
            "softwareType": "apt",
            "url": "",
            "version": "0.19-r0"
        },
        {
            "name": "netbase",
            "softwareType": "apt",
            "url": "",
            "version": "1:6.4-r0"
        },
        {
            "name": "nettle",
            "softwareType": "apt",
            "url": "",
            "version": "3.9.1-r0"
        },
        {
            "name": "networkmanager",
            "softwareType": "apt",
            "url": "",
            "version": "1.46.0-r0"
        },
        {
            "name": "networkmanager-daemon",
            "softwareType": "apt",
            "url": "",
            "version": "1.46.0-r0"
        },
        {
            "name": "networkmanager-locale-en-gb",
            "softwareType": "apt",
            "url": "",
            "version": "1.46.0-r0"
        },
        {
            "name": "networkmanager-nmcli",
            "softwareType": "apt",
            "url": "",
            "version": "1.46.0-r0"
        },
        {
            "name": "networkmanager-wifi",
            "softwareType": "apt",
            "url": "",
            "version": "1.46.0-r0"
        },
        {
            "name": "nftables",
            "softwareType": "apt",
            "url": "",
            "version": "1.0.9-r0"
        },
        {
            "name": "nspr",
            "softwareType": "apt",
            "url": "",
            "version": "4.35-r0"
        },
        {
            "name": "nss",
            "softwareType": "apt",
            "url": "",
            "version": "3.98-r0"
        },
        {
            "name": "ofono",
            "softwareType": "apt",
            "url": "",
            "version": "2.4-r0"
        },
        {
            "name": "openssh-sftp-server",
            "softwareType": "apt",
            "url": "",
            "version": "9.6p1-r0"
        },
        {
            "name": "openssl",
            "softwareType": "apt",
            "url": "",
            "version": "3.2.4-r0"
        },
        {
            "name": "openssl-bin",
            "softwareType": "apt",
            "url": "",
            "version": "3.2.4-r0"
        },
        {
            "name": "openssl-conf",
            "softwareType": "apt",
            "url": "",
            "version": "3.2.4-r0"
        },
        {
            "name": "openssl-ossl-module-legacy",
            "softwareType": "apt",
            "url": "",
            "version": "3.2.4-r0"
        },
        {
            "name": "os-release",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "pam-plugin-deny",
            "softwareType": "apt",
            "url": "",
            "version": "1.5.3-r0"
        },
        {
            "name": "pam-plugin-permit",
            "softwareType": "apt",
            "url": "",
            "version": "1.5.3-r0"
        },
        {
            "name": "pam-plugin-umask",
            "softwareType": "apt",
            "url": "",
            "version": "1.5.3-r0"
        },
        {
            "name": "pam-plugin-unix",
            "softwareType": "apt",
            "url": "",
            "version": "1.5.3-r0"
        },
        {
            "name": "pam-plugin-warn",
            "softwareType": "apt",
            "url": "",
            "version": "1.5.3-r0"
        },
        {
            "name": "parted",
            "softwareType": "apt",
            "url": "",
            "version": "3.6-r0"
        },
        {
            "name": "pciutils",
            "softwareType": "apt",
            "url": "",
            "version": "3.11.1-r0"
        },
        {
            "name": "pciutils-ids",
            "softwareType": "apt",
            "url": "",
            "version": "3.11.1-r0"
        },
        {
            "name": "perl",
            "softwareType": "apt",
            "url": "",
            "version": "5.38.2-r0"
        },
        {
            "name": "perl-module-config-heavy",
            "softwareType": "apt",
            "url": "",
            "version": "5.38.2-r0"
        },
        {
            "name": "pi-bluetooth",
            "softwareType": "apt",
            "url": "",
            "version": "0.1.17-r0"
        },
        {
            "name": "pinentry",
            "softwareType": "apt",
            "url": "",
            "version": "1.2.1-r0"
        },
        {
            "name": "psplash",
            "softwareType": "apt",
            "url": "",
            "version": "0.1+git0+ecc1913756-r0"
        },
        {
            "name": "python3-compression",
            "softwareType": "apt",
            "url": "",
            "version": "3.12.9-r0"
        },
        {
            "name": "python3-core",
            "softwareType": "apt",
            "url": "",
            "version": "3.12.9-r0"
        },
        {
            "name": "rauc",
            "softwareType": "apt",
            "url": "",
            "version": "1.13-r0"
        },
        {
            "name": "rauc-conf",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "rauc-grow-data-part",
            "softwareType": "apt",
            "url": "",
            "version": "1.13-r0"
        },
        {
            "name": "rauc-mark-good",
            "softwareType": "apt",
            "url": "",
            "version": "1.13-r0"
        },
        {
            "name": "rauc-service",
            "softwareType": "apt",
            "url": "",
            "version": "1.13-r0"
        },
        {
            "name": "rpcbind",
            "softwareType": "apt",
            "url": "",
            "version": "1.2.6-r0"
        },
        {
            "name": "runc-opencontainers",
            "softwareType": "apt",
            "url": "",
            "version": "1.1.14+git0+2c9f5602f0-r0"
        },
        {
            "name": "shadow",
            "softwareType": "apt",
            "url": "",
            "version": "4.14.2-r0"
        },
        {
            "name": "shadow-base",
            "softwareType": "apt",
            "url": "",
            "version": "4.14.2-r0"
        },
        {
            "name": "shadow-securetty",
            "softwareType": "apt",
            "url": "",
            "version": "4.6-r0"
        },
        {
            "name": "shared-mime-info",
            "softwareType": "apt",
            "url": "",
            "version": "2.4-r0"
        },
        {
            "name": "shared-mime-info-locale-en-gb",
            "softwareType": "apt",
            "url": "",
            "version": "2.4-r0"
        },
        {
            "name": "squashfs-tools",
            "softwareType": "apt",
            "url": "",
            "version": "4.6.1-r0"
        },
        {
            "name": "sudo",
            "softwareType": "apt",
            "url": "",
            "version": "1.9.15p5-r0"
        },
        {
            "name": "sudo-lib",
            "softwareType": "apt",
            "url": "",
            "version": "1.9.15p5-r0"
        },
        {
            "name": "sudo-sudo",
            "softwareType": "apt",
            "url": "",
            "version": "1.9.15p5-r0"
        },
        {
            "name": "systemd",
            "softwareType": "apt",
            "url": "",
            "version": "1:255.17-r0"
        },
        {
            "name": "systemd-compat-units",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "systemd-conf",
            "softwareType": "apt",
            "url": "",
            "version": "1:1.0-r0"
        },
        {
            "name": "systemd-extra-utils",
            "softwareType": "apt",
            "url": "",
            "version": "1:255.17-r0"
        },
        {
            "name": "systemd-serialgetty",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "systemd-udev-rules",
            "softwareType": "apt",
            "url": "",
            "version": "1:255.17-r0"
        },
        {
            "name": "systemd-vconsole-setup",
            "softwareType": "apt",
            "url": "",
            "version": "1:255.17-r0"
        },
        {
            "name": "tedge",
            "softwareType": "apt",
            "url": "",
            "version": "1.4.2-r0"
        },
        {
            "name": "tedge-agent",
            "softwareType": "apt",
            "url": "",
            "version": "1.4.2-r0"
        },
        {
            "name": "tedge-bootstrap",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "tedge-command-plugin",
            "softwareType": "apt",
            "url": "",
            "version": "1.0.0~rc2+git0+23af5d1f68-r0"
        },
        {
            "name": "tedge-container-plugin",
            "softwareType": "apt",
            "url": "",
            "version": "2.0.0~rc24-r0"
        },
        {
            "name": "tedge-firmware-rauc",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "tedge-inventory",
            "softwareType": "apt",
            "url": "",
            "version": "0.1.0+git0+96078aa3e3-r0"
        },
        {
            "name": "tedge-mapper-aws",
            "softwareType": "apt",
            "url": "",
            "version": "1.4.2-r0"
        },
        {
            "name": "tedge-mapper-az",
            "softwareType": "apt",
            "url": "",
            "version": "1.4.2-r0"
        },
        {
            "name": "tedge-mapper-c8y",
            "softwareType": "apt",
            "url": "",
            "version": "1.4.2-r0"
        },
        {
            "name": "tedge-mapper-collectd",
            "softwareType": "apt",
            "url": "",
            "version": "1.4.2-r0"
        },
        {
            "name": "tedge-nodered-plugin",
            "softwareType": "apt",
            "url": "",
            "version": "1.0.0-r0"
        },
        {
            "name": "tedge-sethostname",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "tini",
            "softwareType": "apt",
            "url": "",
            "version": "0.19.0-r0"
        },
        {
            "name": "u-boot-env",
            "softwareType": "apt",
            "url": "",
            "version": "1:2024.01-r0"
        },
        {
            "name": "udev",
            "softwareType": "apt",
            "url": "",
            "version": "1:255.17-r0"
        },
        {
            "name": "udev-hwdb",
            "softwareType": "apt",
            "url": "",
            "version": "1:255.17-r0"
        },
        {
            "name": "udev-rules-rpi",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "update-alternatives-opkg",
            "softwareType": "apt",
            "url": "",
            "version": "0.6.3-r0"
        },
        {
            "name": "update-rc.d",
            "softwareType": "apt",
            "url": "",
            "version": "0.8+git0+b8f9501050-r0"
        },
        {
            "name": "usbutils",
            "softwareType": "apt",
            "url": "",
            "version": "017-r0"
        },
        {
            "name": "volatile-binds",
            "softwareType": "apt",
            "url": "",
            "version": "1.0-r0"
        },
        {
            "name": "wireless-regdb-static",
            "softwareType": "apt",
            "url": "",
            "version": "2024.10.07-r0"
        },
        {
            "name": "wpa-supplicant",
            "softwareType": "apt",
            "url": "",
            "version": "2.10-r0"
        },
        {
            "name": "wpa-supplicant-cli",
            "softwareType": "apt",
            "url": "",
            "version": "2.10-r0"
        },
        {
            "name": "wpa-supplicant-passphrase",
            "softwareType": "apt",
            "url": "",
            "version": "2.10-r0"
        },
        {
            "name": "wpa-supplicant-plugins",
            "softwareType": "apt",
            "url": "",
            "version": "2.10-r0"
        },
        {
            "name": "xxhash",
            "softwareType": "apt",
            "url": "",
            "version": "0.8.2-r0"
        },
        {
            "name": "portainer-ce",
            "softwareType": "container-group",
            "url": "",
            "version": "2.21.4-1"
        }
    ],
    "c8y_SupportedConfigurations": [
        "tedge-configuration-plugin",
        "tedge-log-plugin",
        "tedge.toml"
    ],
    "c8y_SupportedLogs": [
        "software-management"
    ],
    "c8y_SupportedOperations": [
        "c8y_Command",
        "c8y_DeviceProfile",
        "c8y_DownloadConfigFile",
        "c8y_Firmware",
        "c8y_LogfileRequest",
        "c8y_RemoteAccessConnect",
        "c8y_Restart",
        "c8y_SoftwareUpdate",
        "c8y_UploadConfigFile"
    ],
    "com_cumulocity_model_Agent": {},
    "type": "thin-edge.io",
    "subtype": "dummy",
} + device