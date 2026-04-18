import time
import threading
import websocket
from flask import Flask, request, jsonify

from modules.manager import StreamDeckManager
from modules.ui import render_key, render_home_view, render_blank_view
from modules.actions import handle_restart_logic, switch_dashboard_view

exit_event = threading.Event()
manager = StreamDeckManager(exit_event)

# State tracking voor dirty checking
last_state = {
    "view": None,
    "notif_count": -1,
    "restart_count": -1,
    "brightness": -1,
    "dash_view": None
}

def ws_listener():
    while not exit_event.is_set():
        try:
            ws = websocket.WebSocketApp("ws://localhost:8080/ws", 
                on_message=lambda ws, msg: manager.state["notifications"].append({"msg": msg, "level": 3}))
            ws.run_forever()
        except:
            exit_event.wait(5)

def update_display_loop():
    global last_state
    while not exit_event.is_set():
        with manager.lock:
            now = time.time()
            
            # Restart Cooldown (Atomic check)
            if manager.state["restart_count"] > 0 and now - manager.state["last_restart_press"] > 10:
                manager.state["restart_count"] = 0
            
            # Sleep logic
            if now - manager.state["last_interaction"] > 60:
                if not manager.state["is_sleeping"]:
                    manager.state["is_sleeping"] = True
                    manager.deck.set_brightness(0)
            
            if not manager.state["is_sleeping"]:
                manager.deck.set_brightness(manager.state["brightness"])
                
                # Dirty check: Vergelijk huidige state met de vorige rendered state
                current_notif_count = len(manager.state["notifications"])
                needs_update = (
                    manager.state["current_view"] != last_state["view"] or
                    current_notif_count != last_state["notif_count"] or
                    manager.state["restart_count"] != last_state["restart_count"] or
                    manager.state["brightness"] != last_state["brightness"]
                )
                
                if needs_update:
                    if manager.state["current_view"] == "home":
                        render_home_view(manager, render_key)
                    elif manager.state["current_view"] == "inbox":
                        render_blank_view(manager.deck, render_key)
                        manager.deck.set_key_image(0, render_key(manager.deck, "BACK", bg="gray"))
                    
                    # Cache de huidige state
                    last_state.update({
                        "view": manager.state["current_view"],
                        "notif_count": current_notif_count,
                        "restart_count": manager.state["restart_count"],
                        "brightness": manager.state["brightness"]
                    })
            
        exit_event.wait(0.1)

def key_callback(deck, key, pressed):
    if not pressed: return
    with manager.lock:
        manager.state["last_interaction"] = time.time()
        if manager.state["is_sleeping"]:
            manager.state["is_sleeping"] = False
            return
            
        if manager.state["current_view"] == "home":
            if key == 0: switch_dashboard_view(manager, "Home")
            elif key == 1: switch_dashboard_view(manager, "Network")
            elif key == 2: switch_dashboard_view(manager, "Docker")
            elif key == 4: manager.state["current_view"] = "inbox"
            elif key == 10: manager.state["brightness"] = max(0, manager.state["brightness"] - 10)
            elif key == 11: manager.state["brightness"] = min(100, manager.state["brightness"] + 10)
            elif key == 14: handle_restart_logic(manager)
        elif manager.state["current_view"] == "inbox" and key == 0:
            manager.state["current_view"] = "home"

if __name__ == "__main__":
    while not manager.connect() and not exit_event.is_set():
        time.sleep(5)
    
    manager.set_key_callback(key_callback)
    threading.Thread(target=update_display_loop, daemon=True).start()
    threading.Thread(target=ws_listener, daemon=True).start()
    
    try:
        exit_event.wait()
    except KeyboardInterrupt:
        exit_event.set()
    
    manager.close()
