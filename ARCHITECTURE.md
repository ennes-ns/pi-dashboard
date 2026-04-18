# Technical Architecture

## Design Philosophy
PI-DASHBOARD is built on the principle of **Atomic Separation**. The display logic (Dashboard) and the control logic (Controller) are decoupled, communicating via a low-latency local WebSocket.

## 1. Dashboard (The Brain)
Written in **Go 1.24**, the dashboard is responsible for:
- Gathering system metrics (CPU, MEM, LOAD, UPTIME) every 1s.
- Rendering a high-resolution TUI using **Bubble Tea** and **Lipgloss**.
- Managing a **WebSocket Hub** to push real-time alerts to connected controllers.
- Exposing a REST API for view-switching commands.

## 2. Controller (The Interface)
Written in **Python 3.11**, the controller manages:
- **USB HID Communication:** Direct interface with the Stream Deck using the `StreamDeck` library.
- **Dynamic Rendering:** Generating 72x72px icons using `Pillow` with real-time status overlays (e.g., brightness bars).
- **State Persistence:** Managing brightness levels, notification inboxes, and safety counters.
- **Atomic Actions:** Implementing multi-tap confirmation logic for high-impact system commands.

## 3. Deployment & Security
- **Hardened Containers:** Both services run in minimal Alpine/Slim images.
- **Network Isolation:** The WebSocket server binds to `127.0.0.1`, preventing external access even within the local network.
- **TTY Injection:** Uses `openvt` to attach the Docker TTY session directly to `/dev/tty1`, bypassing the need for X11 or Wayland.
