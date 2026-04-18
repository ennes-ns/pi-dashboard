# PI-DASHBOARD: Operational Directives

## ⚠️ Critical Constraints
- **Native Display Requirement:** The Go Dashboard MUST run natively (via systemd) to stably drive the 1080p TTY1. Docker TTY mapping proved unreliable for the Bubble Tea UI.
- **Zero Input Environment:** The Pi has NO keyboard or mouse. All interaction is handled via the Stream Deck (WebSocket/HTTP).
- **Atomic Refactors:** Keep changes small and modular (per file) to prevent regressions in the TUI layout.

## 🛠️ Dashboard Management (Native)
- **Binary Path:** `/home/icarus/docker/pi-dashboard/dashboard/local-dash`
- **Service Name:** `pi-dashboard.service`
- **Build Command:** `cd dashboard && go build -o local-dash .`

## 📺 Kiosk Protocol
1. The dashboard is managed by `systemd`.
2. The Stream Deck controller runs in **Docker** (`streamdeck-daemon`).
3. Communication is established via `http://localhost:8080`.
