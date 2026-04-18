# PI-DASHBOARD

A high-performance, minimalist TUI dashboard for Raspberry Pi 5, optimized for 1080p TTY display and Stream Deck integration.

## Features
- **Minimalist TUI:** Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- **1080p Optimized:** Designed specifically for 1920x1080 resolution on TTY1.
- **Dockerized:** Small footprint, easy deployment.
- **Remote Control:** Built-in HTTP listener for view switching (perfect for Stream Deck).

## Quick Start
```bash
docker compose up -d
# Inject into TTY1 (HDMI)
sudo openvt -c 1 -s -f -- docker exec -it pi-dashboard ./main
```

## Architecture
- **Language:** Go 1.24
- **Framework:** Bubble Tea / Lipgloss
- **Deployment:** Docker (multi-stage)
- **Interface:** TTY1 (direct to HDMI)

## Stream Deck Integration
Send a simple GET request to switch views:
`http://<pi-ip>:8080/switch?view=Network`
