# Technical Architecture: Hybrid Native-Container Model

## Why Hybrid?
During development, we discovered that Docker's TTY mapping for high-resolution terminal applications (Bubble Tea) introduced significant artifacts and input/output latency on Raspberry Pi 5. To achieve industrial-grade stability, we moved the display engine to a native process while keeping the control plane containerized.

## Component Breakdown

### 🖥️ Native Dashboard (Go)
- **Role:** High-Res Display & System Metrics.
- **Process:** Managed by `pi-dashboard.service`.
- **TTY Target:** Direct binding to `/dev/tty1`.
- **API:** Exposes a local-only server on port 8080 for view switching and WebSocket alerts.

### ⌨️ Containerized Controller (Python)
- **Role:** HID Management & UI Logic.
- **Container:** `streamdeck-daemon`.
- **Connectivity:** Uses `network_mode: host` to communicate with the native Dashboard.
- **Safety:** Implements Atomic 5-tap confirmation for critical server actions.

## Communication Flow
1. User presses button on **Stream Deck**.
2. **Controller (Docker)** receives HID event.
3. **Controller** sends HTTP GET to `localhost:8080/switch`.
4. **Dashboard (Native)** updates TUI state.
5. **Dashboard** pushes real-time system alerts back to **Controller** via WebSocket.
