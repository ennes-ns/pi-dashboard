import time
import threading
from StreamDeck.DeviceManager import DeviceManager

class StreamDeckManager:
    def __init__(self, exit_event):
        self.deck = None
        self.exit_event = exit_event
        self.lock = threading.Lock()
        self.state = {
            "brightness": 50,
            "last_interaction": time.time(),
            "is_sleeping": False,
            "current_view": "home",
            "notifications": [],
            "restart_count": 0,
            "last_restart_press": 0
        }

    def connect(self):
        decks = DeviceManager().enumerate()
        if not decks: return False
        self.deck = decks[0]
        self.deck.open()
        self.deck.reset()
        self.deck.set_brightness(self.state["brightness"])
        return True

    def set_key_callback(self, callback):
        if self.deck:
            self.deck.set_key_callback(callback)

    def close(self):
        if self.deck:
            self.deck.close()
