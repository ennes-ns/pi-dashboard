# PI-DASHBOARD: Operational Directives

## ⚠️ Critical Constraints
- **Native Display Requirement:** Het Go Dashboard MOET native draaien (via systemd) om de 1080p TTY1 stabiel aan te sturen. Docker-mapping van TTY1 bleek onbetrouwbaar voor Bubble Tea UI.
- **Zero Input Environment:** De Pi heeft GEEN toetsenbord. Alle interactie verloopt via de Stream Deck (WebSocket/HTTP).
- **Atomic Refactors:** Houd wijzigingen klein en modulair (per bestand) om fouten in de TUI-layout te voorkomen.

## 🛠️ Dashboard Management (Native)
- **Binary:** `/home/icarus/docker/pi-dashboard/dashboard/local-dash`
- **Service:** `pi-dashboard.service`
- **Build:** `cd dashboard && go build -o local-dash .`

## 📺 Kiosk Protocol
1. Het dashboard wordt beheerd door `systemd`.
2. De Stream Deck controller draait in **Docker** (`streamdeck-daemon`).
3. Communicatie verloopt via `http://localhost:8080`.
