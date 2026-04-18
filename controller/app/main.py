import time
import threading
import websocket
from flask import Flask, request, jsonify

from modules.manager import StreamDeckManager
from modules.ui import render_key, render_home_view, render_blank_view
from modules.actions import handle_restart_logic, switch_dashboard_view

exit_event = threading.Event()
manager = StreamDeckManager(exit_event)

def ws_listener():
    while not exit_event.is_set():
        try:
            ws = websocket.WebSocketApp("ws://localhost:8080/ws", 
                on_message=lambda ws, msg: manager.state["notifications"].append({"msg": msg, "level": 3}))
            ws.run_forever()
        except:
            exit_event.wait(5)

def update_display_loop():
    while not exit_event.is_set():
        with manager.lock:
            # Sleep logic (60s timeout)
            if time.time() - manager.state["last_interaction"] > 60:
                manager.state["is_sleeping"] = True
                manager.deck.set_brightness(0)
            
            if not manager.state["is_sleeping"]:
                manager.deck.set_brightness(manager.state["brightness"])
                if manager.state["current_view"] == "home":
                    render_home_view(manager, render_key)
                elif manager.state["current_view"] == "inbox":
                    render_blank_view(manager.deck, render_key)
                    manager.deck.set_key_image(0, render_key(manager.deck, "BACK", bg="gray"))
            
        exit_event.wait(0.2)

def key_callback(deck, key, pressed):
    if not pressed: return
    with manager.lock:
        manager.state["last_interaction"] = time.time()
        if manager.state["is_sleeping"]:
            manager.state["is_sleeping"] = False
            return
            
        if manager.state["current_view"] == "home":
            if key == 0: switch_dashboard_view("Home")
            elif key == 1: switch_dashboard_view("Network")
            elif key == 2: switch_dashboard_view("Docker")
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
