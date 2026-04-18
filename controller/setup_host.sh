#!/usr/bin/env bash
# Host setup script for Elgato Stream Deck UDEV rules
# Run this on the Raspberry Pi host (sudo ./setup_host.sh)

echo "Installing Elgato Stream Deck udev rules..."

cat <<EOF > /etc/udev/rules.d/70-streamdeck.rules
# Elgato Stream Deck rules to allow non-root access
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", TAG+="uaccess"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0063", TAG+="uaccess"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006c", TAG+="uaccess"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006d", TAG+="uaccess"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0080", TAG+="uaccess"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0084", TAG+="uaccess"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0090", TAG+="uaccess"
EOF

echo "Reloading udev rules..."
udevadm control --reload-rules
udevadm trigger

echo "Stream Deck udev rules installed and reloaded."
