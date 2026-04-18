import requests
import os
import socket
import struct

MACRO_RESTART_URL = os.getenv("MACRO_RESTART_URL", "http://localhost/api/restart")
WOL_MAC = os.getenv("WOL_MAC_ADDRESS", "D8:5E:D3:FA:13:CC")

from wakeonlan import send_magic_packet

def send_magic_packet_wrapped(mac_address):
    try:
        # Provide broadcast address or let wakeonlan use its default
        send_magic_packet(mac_address)
        return True
    except Exception:
        return False

def trigger_webhook(url, payload=None):
    try:
        r = requests.post(url, json=payload, timeout=2)
        return r.status_code == 200
    except Exception:
        return False

def handle_macro(macro_name):
    # Map macro names to actions
    if macro_name == "RESTART_SRV":
        trigger_webhook(MACRO_RESTART_URL)
    elif macro_name == "WAKE_PC":
        return send_magic_packet_wrapped(WOL_MAC)
