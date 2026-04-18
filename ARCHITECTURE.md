# Technical Architecture: Native TTY Model

## Design Philosophy
PI-DASHBOARD follows the principle of **Maximum Native Performance**. By stripping away graphical layers (Wayland/X11), we ensure the system resources are dedicated to actual workloads while maintaining a robust monitoring interface.

## Component Breakdown

### 🖥️ Native Dashboard (btop)
- **Role:** High-Density System Monitoring.
- **Process:** Managed by `btop-kiosk.service` running as `root`.
- **TTY Target:** Direct hardware binding to `/dev/tty1`.
- **Configuration:** No-skin (Default theme) to ensure compatibility with 16-color Linux Console limitations.

### ⌨️ Containerized Controller (Python)
- **Role:** HID Management & State Persistence.
- **Container:** `streamdeck-daemon` (Managed via `./controller/docker-compose.yml`).
- **Safety:** Implements Atomic 5-tap confirmation for critical server actions.

## Deployment & Security
- **Hardened Kiosk:** The default Linux login (`getty@tty1`) is masked to prevent login prompt leakage.
- **Power Management:** Screen blanking is disabled at the terminal level via `setterm`.
- **Zero Input:** The TTY is configured as output-only, optimized for headless operation.
