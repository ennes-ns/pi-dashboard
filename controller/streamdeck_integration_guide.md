# Stream Deck & AutoHotkey Integration Guide

This document explains the bi-directional communication protocol between the Raspberry Pi Stream Deck Daemon and any external listener (like an AutoHotkey HTTP Server on a Windows workstation, or a VPS). Agents working on the Windows Workstation can use this instruction guide to seamlessly connect AHK scripts to the Stream Deck.

## 1. Outgoing Triggers (Stream Deck ➡️ AHK Workstation)

De Stream Deck kan HTTP/REST requests (zoals POST of GET) versturen zodra er een fysieke knop wordt ingedrukt.
Dit gebeurt via de ingebouwde Python library `requests` in `app/plugins/macros.py`.

### AHK Server Vereisten
Op het Windows werkstation moet een lokale HTTP server draaien (bijvoorbeeld via een AHK HTTP server library, of een Node.js / Python wrapper rondom AHK processen). Deze server luistert typisch op een specifieke poort gelinkt aan het Tailscale IP (bijv. `http://100.x.y.z:8080/api/macro`). Zorg hierbij dat de AHK luisterservice is gebonden aan `0.0.0.0` of je Tailscale-IP in plaats van uitsluitend `localhost`, anders kan de Pi hem niet met HTTP bereiken!

### Hoe voeg je een nieuwe knop toe op de Pi? (Instructies voor de Pi-Agent)
Binnen de `streamdeck-daemon` Docker container (op de Pi) volg je deze drie stappen voor een nieuwe integratie:

1. **Definieer de API URL dynamisch:** Voeg in de `docker-compose.yml` van de Pi het werkstation IP toe als Environment Variable:
   ```yaml
   environment:
     - AHK_ENDPOINT=http://<WORKSTATION_TAILSCALE_IP>:8080/api/macro
   ```
2. **Breid de macro-handler uit (`app/plugins/macros.py`):**
   Voeg de logica toe in Python om requests te sturen en eventuele data terug te vangen.
   ```python
   AHK_URL = os.getenv("AHK_ENDPOINT", "")

   def trigger_ahk_macro(action_id):
       payload = {"action": action_id}
       try:
           # Send data to Workstation
           r = requests.post(AHK_URL, json=payload, timeout=3)
           if r.status_code == 200:
               return r.json() # Stuur eventuele JSON data terug
       except Exception:
           pass
       return None
   ```
3. **Koppel aan een fysieke toets (`app/main.py`):**
   Zoek de functie `key_change_callback` en voeg een nieuwe conditionele 'if' actie toe aan de toegewezen knop, bijvoorbeeld knop `12`.
   ```python
   elif key == 12:
       response = trigger_ahk_macro("START_GAME_MODE")
       if response and response.get("status") == "ok":
           # Eventueel het knopje op de deck kort groen laten kleuren
           pass
   ```

## 2. Incoming Triggers (AHK Workstation ➡️ Stream Deck)

Andersom kan de Windows AHK server ook direct visuele alarms en notificaties wekken op de Stream Deck. De Stream Deck applicatie draait een luisterende Flask-server intern op poort `5000` via Tailscale IP.

### Notificaties sturen vanaf het Workstation (Instructies voor AHK agent)
Het AHK script hoeft enkel een `POST` verzoek te versturen op het moment dat een automatiserings-script afrond of faalt, of er een trigger ingaat. 

**Hoofd Endpoint:** `POST http://<PI_TAILSCALE_IP>:5000/api/notify`
**Headers vereist:** `Content-Type: application/json`
**Voorbeeld Body (`JSON`):**
```json
{
  "message": "AHK Script Finished!",
  "level": 1
}
```

*   `"level": 1`: Subtiele actie. Het scherm van het deck kleurt kort, subtiel groen (3 seconden), en verdwijnt daarna achter de schermen naar de 'Inbox'. Ideaal bij `status: ok`!
*   `"level": 2`: Attentie-actie. Het blijvend met oranje en zwart knipperende waarschuwingscherm dat de knop 'Inbox' vult met deze notificatie, en eist dat de eigenaar hem fysiek op het doosje weg tikt.
*   `"level": 3`: Critical-actie. Waarschuwt direct door agressief snel rood / wit vol over te nemen. Actie vereist alvorens andere knoppen in te drukken zijn. 

## Netwerk en Beveiliging

Omdat beide machines, Pi en het Workstation, uitsluitend via de private virtual netwerken van **Tailscale** praten:
- Maak direct en onbeveiligd gebruik van `HTTP`, HTTPS certificaten zijn lokaal onnodig omwille van de beveiligde tunnel.
- Geen inlogwachtwoorden benodigd aan deze URLs in dit prototype. 
- Mocht de Pi poort `5000` onbereikbaar zijn vanaf Windows PowerShell, zorg dan dat de lokale firewall gaten toelaat over je tailscale virtuele netwerk adapter (`sudo ufw allow 5000/tcp` mocht er een strikte Linux firewall draaien, of Windows firewall).
