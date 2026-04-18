# PI-DASHBOARD v2.0

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.11-yellow.svg)](https://www.python.org/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue.svg)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A professional, high-performance TUI ecosystem for Raspberry Pi 5. PI-DASHBOARD transforms your terminal into a 1080p system sentinel, fully integrated with Elgato Stream Deck hardware via a low-latency WebSocket backbone.

---

## 🚀 Overview

PI-DASHBOARD is a dual-component architecture designed for maximum performance and minimal system overhead. It consists of a **High-Res Go TUI** optimized for 1080p HDMI output and a **Python-based Controller** that interfaces with Stream Deck HID devices.

### Key Features
- 🖥️ **1080p TUI:** Built with Go & Bubble Tea, optimized for 1920x1080 TTY1 output.
- ⌨️ **HID Integration:** Direct Stream Deck control with custom-rendered key icons.
- ⚡ **Real-time Synchronization:** Full-duplex communication via WebSockets.
- 🛡️ **Secure by Design:** Localhost-only bindings and hardened Docker containers.
- 🔄 **Atomic Operations:** 5-tap safety confirmation for critical actions like server restarts.
- 🌙 **Smart Power Management:** Auto-dimming and TTY sleep timeout after 60s of inactivity.

---

## 🏗️ System Architecture

```text
[ Hardware ]          [ Controller (Python) ]          [ Dashboard (Go) ]
+-------------+        +----------------------+        +-------------------+
| Stream Deck | <----> | HID Manager          | <----> | WebSocket Hub     |
| (USB HID)   |        | UI Renderer          |        | Sytem Stats (OS)  |
+-------------+        | API Client           |        | Bubble Tea UI     |
                       +----------------------+        +-------------------+
                                  ^                             |
                                  |                             v
                                  +----------------------- [ TTY1 / HDMI ]
```

---

## 🛠️ Project Structure

- `dashboard/`: Go-based TUI application.
- `controller/`: Python-based HID and state manager.
- `docker-compose.yaml`: Unified deployment orchestration.

---

## 🚦 Quick Start

### Prerequisites
- Raspberry Pi 5 (Optimized for Raspberry Pi OS 64-bit)
- Docker & Docker Compose
- Elgato Stream Deck connected via USB

### Deployment
1. Clone the repository:
   ```bash
   git clone https://github.com/ennes-ns/pi-dashboard.git
   cd pi-dashboard
   ```

2. Launch the ecosystem:
   ```bash
   docker compose up -d --build
   ```

3. Enable the Kiosk Mode (TTY1 Injection):
   ```bash
   sudo systemctl enable --now pi-dashboard-kiosk.service
   ```

---

## 🔧 Technical Stack

- **Core:** Go 1.24, Python 3.11
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Monitoring:** [gopsutil](https://github.com/shirou/gopsutil)
- **Communication:** WebSockets (Gorilla/WS, websocket-client)
- **Containerization:** Docker (Multi-stage builds)
