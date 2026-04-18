# PI-DASHBOARD Architecture

## Componenten
1. **Core Dashboard (Go):** 
   - TUI framework: `github.com/charmbracelet/bubbletea`
   - Styling: `github.com/charmbracelet/lipgloss`
   - HTTP Listener: Ontvangt commando's van Stream Deck (bijv. `POST /view/switch`).
2. **Docker Container:**
   - Base image: `golang:alpine` voor build, `alpine` voor runtime.
   - Privileged of device mapping voor TTY toegang indien nodig (of via `openvt` host-side).
3. **Stream Deck Bridge:**
   - Container die de fysieke Stream Deck hardware uitleest (USB HID).
   - Stuurt HTTP requests naar de PI-DASHBOARD container op basis van knop-acties.

## Data Flow
Hardware -> Stream Deck Container -> [HTTP/JSON] -> PI-DASHBOARD (Go) -> TTY1 (HDMI)
