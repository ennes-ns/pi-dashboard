# PI-DASHBOARD: Master Directives

## Tonal & Structural Mandates
- **No Cringe:** Vermijd alle "AI-marketing" termen (Sentinel, Nexus, etc.). Wees direct en technisch.
- **Kiosk Target:** Alle output is geoptimaliseerd voor 1920x1080 TTY1 (HDMI-A-1).
- **Go Focused:** Gebruik Go met `Bubble Tea` voor de TUI.

## Deployment Workflow
1. **Build:** Docker multi-stage build voor een minimale binary.
2. **Inject:** Gebruik `openvt -c 1 -s -f -- docker exec -it ...` om direct op HDMI te tonen.
3. **Persist:** Systemd of Docker restart policies zorgen voor weergave na reboot.

## Stream Deck Integration
- Het dashboard luistert op poort 8080 naar JSON-events van de Stream Deck container.
- Knop 1 -> Dashboard View A
- Knop 2 -> Dashboard View B (Btop/Logs/Stats)
