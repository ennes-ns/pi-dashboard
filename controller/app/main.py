import os
import time
import threading
import requests
import websocket
import json
from StreamDeck.DeviceManager import DeviceManager
from StreamDeck.ImageHelpers import PILHelper
from PIL import Image, ImageDraw, ImageFont
from flask import Flask, request, jsonify

from plugins.pi_monitor import get_cpu_temp, get_cpu_percent
from plugins.vps_monitor import get_vps_status
from plugins.macros import handle_macro

# --- Constants & Configuration ---
DASHBOARD_WS_URL = os.getenv("DASHBOARD_WS_URL", "ws://localhost:8080/ws")
DASHBOARD_HTTP_URL = "http://localhost:8080/switch"
BRIGHTNESS_STEP = 10
SLEEP_TIMEOUT = 60

class StreamDeckManager:
    def __init__(self):
        self.deck = None
        self.exit_event = threading.Event()
        self.lock = threading.Lock()
        self.state = {
            "brightness": 50,
            "last_interaction": time.time(),
            "is_sleeping": False,
            "current_view": "home",
            "notifications": [],
            "active_alert": None
        }

    def connect(self):
        decks = DeviceManager().enumerate()
        if not decks: return False
        self.deck = decks[0]
        self.deck.open()
        self.deck.reset()
        self.deck.set_brightness(self.state["brightness"])
        self.deck.set_key_callback(self.key_callback)
        return True

    def ws_listener(self):
        def on_message(ws, message):
            print(f"WS Message: {message}")
            if "ALERT" in message:
                self.add_notification(message, level=3)

        while not self.exit_event.is_set():
            try:
                ws = websocket.WebSocketApp(DASHBOARD_WS_URL, on_message=on_message)
                ws.run_forever()
            except:
                time.sleep(5)

    def render_key(self, label, bg="black", fg="white"):
        image = PILHelper.create_image(self.deck)
        draw = ImageDraw.Draw(image)
        try: font = ImageFont.truetype("arial.ttf", 14)
        except: font = ImageFont.load_default()
        draw.rectangle([0, 0, image.width, image.height], fill=bg, outline="white")
        bbox = draw.textbbox((0, 0), label, font=font)
        x, y = (image.width-(bbox[2]-bbox[0]))//2, (image.height-(bbox[3]-bbox[1]))//2
        draw.text((x, y), label, font=font, fill=fg)
        return PILHelper.to_native_format(self.deck, image)

    def update_display_loop(self):
        while not self.exit_event.is_set():
            with self.lock:
                if not self.state["active_alert"] and not self.state["is_sleeping"]:
                    if time.time() - self.state["last_interaction"] > SLEEP_TIMEOUT:
                        self.state["is_sleeping"] = True
                        self.deck.set_brightness(0)
                if not self.state["is_sleeping"]:
                    self.deck.set_brightness(self.state["brightness"])
                    if self.state["active_alert"]:
                        img = self.render_key(f"ALERT:\n{self.state['active_alert']['msg'][:10]}", bg="darkred")
                        for i in range(self.deck.key_count()): self.deck.set_key_image(i, img)
                    elif self.state["current_view"] == "inbox":
                        self.deck.set_key_image(0, self.render_key("BACK", bg="gray"))
                    else:
                        self.render_home()
            self.exit_event.wait(0.2)

    def render_home(self):
        mappings = {
            0: ("DASH\nHOME", "darkblue"), 1: ("DASH\nNET", "darkblue"), 2: ("DASH\nDOC", "darkblue"),
            4: (f"INBOX\n[{len(self.state['notifications'])}]", "royalblue"),
            10: ("BRIGHT\n-", "teal"), 11: (f"BRIGHT\n{self.state['brightness']}%", "teal"),
            14: ("RESTART\nSRV", "orange")
        }
        for key in range(self.deck.key_count()):
            if key in mappings:
                l, c = mappings[key]
                self.deck.set_key_image(key, self.render_key(l, bg=c))
            else:
                self.deck.set_key_image(key, self.render_key("", bg="black"))

    def key_callback(self, deck, key, pressed):
        if not pressed: return
        with self.lock:
            self.state["last_interaction"] = time.time()
            if self.state["is_sleeping"]:
                self.state["is_sleeping"] = False
                return
            if self.state["current_view"] == "home":
                if key == 0: self.call_dashboard("Home")
                elif key == 1: self.call_dashboard("Network")
                elif key == 2: self.call_dashboard("Docker")
                elif key == 4: self.state["current_view"] = "inbox"
                elif key == 10: self.state["brightness"] = max(0, self.state["brightness"] - BRIGHTNESS_STEP)
                elif key == 11: self.state["brightness"] = min(100, self.state["brightness"] + BRIGHTNESS_STEP)
                elif key == 14: handle_macro("RESTART_SRV")
            elif self.state["current_view"] == "inbox" and key == 0:
                self.state["current_view"] = "home"

    def call_dashboard(self, view):
        try: requests.get(f"{DASHBOARD_HTTP_URL}?view={view}", timeout=0.5)
        except: pass

    def add_notification(self, msg, level):
        with self.lock:
            self.state["notifications"].append({"msg": msg, "level": level})
            self.state["is_sleeping"] = False

manager = StreamDeckManager()

def start_api():
    app = Flask(__name__)
    @app.route('/api/notify', methods=['POST'])
    def notify():
        data = request.json or {}
        manager.add_notification(data.get("message", "Notif"), data.get("level", 1))
        return jsonify({"status": "ok"})
    import logging
    logging.getLogger('werkzeug').setLevel(logging.ERROR)
    app.run(host='0.0.0.0', port=5000)

if __name__ == "__main__":
    while not manager.connect(): time.sleep(5)
    threading.Thread(target=start_api, daemon=True).start()
    threading.Thread(target=manager.ws_listener, daemon=True).start()
    threading.Thread(target=manager.update_display_loop, daemon=True).start()
    try:
        while True: time.sleep(1)
    except KeyboardInterrupt: manager.exit_event.set()
    manager.deck.close()
