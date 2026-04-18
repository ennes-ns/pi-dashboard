# Troubleshooting & Lessons Learned

## 📺 Display Refresh Issues (Kiosk Mode)

### Symptom
Dashboard on TTY1 (HDMI) does not refresh after a code update or `docker compose up`.

### Cause: Docker Project Names
Docker Compose (v2) prefixes container names with the project directory name (e.g., `pi-dashboard-dashboard-1`). If a `systemd` service or a script refers to a static name like `pi-dashboard`, it will fail to find the new container instance.

### Solution: Dynamic Container ID Discovery
In the `pi-dashboard-kiosk.service`, we now use a dynamic lookup to find the correct container ID, regardless of what prefix Docker Compose assigns:

```bash
docker ps -qf "name=pi-dashboard"
```

This ensures the `openvt` command always attaches to the active dashboard process.

## ⌨️ TTY Lock-ins
If the terminal feels "stuck," you can force a terminal switch from the host using:
```bash
sudo chvt 1  # Switch to TTY1 (Dashboard)
sudo chvt 7  # Switch back to default desktop (if active)
```

## 🌐 WebSocket Timeouts
The Stream Deck controller (Python) will automatically attempt to reconnect if the Go dashboard restarts. A 5-second backoff is implemented to prevent CPU spikes during a crash loop.
