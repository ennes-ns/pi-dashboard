# PI-DASHBOARD v2.0

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.11-yellow.svg)](https://www.python.org/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue.svg)](https://www.docker.com/)

A high-performance, minimalist TUI ecosystem for Raspberry Pi 5. PI-DASHBOARD transforms your terminal into a 1080p system sentinel, fully integrated with Elgato Stream Deck hardware.

---

## 🏗️ Architecture: The Hybrid Model

For maximum stability on 1080p TTY displays, PI-DASHBOARD uses a **Hybrid Architecture**:

1.  **Dashboard (The Brain):** A native Go application running as a systemd service. This ensures jitter-free, pixel-perfect rendering direct to HDMI via `/dev/tty1`.
2.  **Controller (The Interface):** A containerized Python application that manages Stream Deck HID communication and WebSocket alerting.

---

## 🚦 Quick Start

### 1. Build and Run Dashboard (Native)
```bash
cd dashboard
go build -o local-dash .
sudo systemctl enable --now ./pi-dashboard.service
```

### 2. Run Controller (Docker)
```bash
docker compose up -d controller
```

---

## 🛠️ Tech Stack
- **Dashboard:** Go 1.24, Bubble Tea, Lipgloss, gopsutil.
- **Controller:** Python 3.11, StreamDeck-Python-SDK, Flask.
- **Backbone:** HTTP/WebSocket on `localhost:8080`.
