import requests
import os

VPS_ENDPOINT = os.getenv("VPS_MONITOR_ENDPOINT", "http://localhost/api/stats")

def get_vps_status():
    try:
        # User defined Go application endpoint returning JSON
        r = requests.get(VPS_ENDPOINT, timeout=3)
        if r.status_code == 200:
            data = r.json()
            # e.g., {"status": "online", "users": 42}
            return data.get("status", "OK")
        return "ERR"
    except Exception as e:
        return "OFFLINE"
