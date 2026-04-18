# Stream Deck Daemon

## What it does
This service provides a headless, custom-built daemon for managing an Elgato Stream Deck Original connected directly to the Raspberry Pi. It monitors local hardware statistics (CPU, Temp) and remote VPS status, displaying them dynamically on the Stream Deck's LCD keys. It also offers macro buttons to trigger webhooks.

## How it works
The service runs in a privileged Docker container mapped to the host's `/dev/bus/usb` to access the Stream Deck via `libusb` and `hidapi`. 
A Python script (`app/main.py`) continuously polls endpoints and hardware files (like `/proc/stat` and `/sys/class/thermal/`), and renders text graphics onto the buttons using the `Pillow` library.

## How to use it
1. Make sure the host UDEV rules are installed by running `sudo ./setup_host.sh` on the Pi once.
2. Run the container: `docker compose up -d` (or via your master podman/docker compose file).
3. **Customization:** To add new endpoints, edit the scripts in `app/plugins/`. Because `./app` is volume-mapped, you do not need to rebuild the image; simply run `docker compose restart`.
