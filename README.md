# PI-DASHBOARD v2.1 (Minimalist)

[![Status](https://img.shields.io/badge/Status-Stable-green.svg)](https://github.com/ennes-ns/pi-dashboard)

A professional, high-performance monitoring ecosystem for Raspberry Pi 5. PI-DASHBOARD is optimized for extreme stability and zero system overhead, utilizing the native Linux console (TTY) for the dashboard and a containerized controller for the interface.

---

## 🏗️ Architecture: The Straight TTY Model

To achieve maximum reliability and eliminate GPU-driver artifacts, PI-DASHBOARD uses a **Native Console Architecture**:

1.  **Dashboard (The Brain):** A native `btop` instance running as a dedicated systemd service. It claims `/dev/tty1` directly, providing a high-density system overview without the overhead of X11, Wayland, or terminal emulators.
2.  **Controller (The Interface):** A containerized Python application that manages Stream Deck HID communication and state management.

---

## 🚀 Deployment

### 1. Dashboard (Native)
The dashboard is managed by the `btop-kiosk.service`. It is pre-configured to:
- Hide the terminal cursor.
- Disable screen blanking and power saving.
- Bind directly to TTY1 (HDMI output).

```bash
sudo systemctl enable --now btop-kiosk.service
```

### 2. Controller (Docker)
The controller handles the physical Stream Deck interaction.

```bash
cd controller
docker compose up -d
```

---

## 🛠️ Tech Stack
- **Dashboard:** Native `btop` (16-color TTY optimized).
- **Controller:** Python 3.11, StreamDeck-Python-SDK.
- **Display:** Physical TTY1 (HDMI-A-1).
