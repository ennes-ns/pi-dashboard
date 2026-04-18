import time
import requests
from plugins.macros import handle_macro

def handle_restart_logic(manager):
    now = time.time()
    # Reset count if last press was more than 10 seconds ago
    if now - manager.state["last_restart_press"] > 10:
        manager.state["restart_count"] = 0
    
    manager.state["restart_count"] += 1
    manager.state["last_restart_press"] = now
    
    if manager.state["restart_count"] >= 5:
        manager.state["restart_count"] = 0
        handle_macro("RESTART_SRV")
        return "RESTARTING..."
    return f"TAP {manager.state['restart_count']}/5"

def switch_dashboard_view(view):
    try:
        requests.get(f"http://localhost:8080/switch?view={view}", timeout=0.5)
    except:
        pass
