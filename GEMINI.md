# PI-DASHBOARD: Operational Directives

## ⚠️ Critical Constraints
- **Zero Input Environment:** De Pi heeft GEEN toetsenbord of muis. Het dashboard mag NOOIT wachten op input. Alle interactie verloopt via de Stream Deck (WebSocket/HTTP).
- **TTY1 Primary:** De HDMI-output is gekoppeld aan TTY1. Gebruik `openvt -c 1` voor injectie.
- **Atomic Refactors:** Houd wijzigingen klein en modulair (per bestand) om fouten in de TUI-layout te voorkomen.

## 🛠️ Dashboard Management
- **Binary Path:** `/app/tmp/main` (gegenereerd door Air).
- **Restart Logic:** De 5-tap restart op de Stream Deck is de enige manier om de controller/dashboard veilig te herstarten zonder SSH.

## 📺 Kiosk Protocol
1. Stop `getty@tty1.service`.
2. Disable terminal blanking (`setterm -blank 0`).
3. Injecteer Docker exec via `openvt`.
