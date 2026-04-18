# PI-DASHBOARD v2.0 (WIP)

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.11-yellow.svg)](https://www.python.org/)
[![Status](https://img.shields.io/badge/Status-Work%20in%20Progress-orange.svg)](https://github.com/ennes-ns/pi-dashboard)

**⚠️ Status: Work In Progress** - De architectuur is onlangs gemigreerd naar een hybride model (Native + Docker). View-switching en real-time alerting zijn in actieve ontwikkeling.

---

## 🏗️ Architecture: The Hybrid Model

For maximum stability on 1080p TTY displays, PI-DASHBOARD uses a **Hybrid Architecture**:

1.  **Dashboard (The Brain):** A native Go application running as a systemd service. This ensures jitter-free, pixel-perfect rendering direct to HDMI via `/dev/tty1`.
2.  **Controller (The Interface):** A containerized Python application that manages Stream Deck HID communication and WebSocket alerting.

---

## 🛠️ Current Focus
- [ ] Stabilizing view-switching synchronization.
- [ ] Optimizing TTY rendering for high-resolution displays.
- [ ] Finalizing Stream Deck action modules.

---

## 🚦 Quick Start

### 1. Build and Run Dashboard (Native)
```bash
cd dashboard
go build -o local-dash .
sudo systemctl enable --now pi-dashboard.service
```

### 2. Run Controller (Docker)
```bash
docker compose up -d controller
```
