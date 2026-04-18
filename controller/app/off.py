from StreamDeck.DeviceManager import DeviceManager
import time

if __name__ == "__main__":
    streamdecks = DeviceManager().enumerate()
    if len(streamdecks) > 0:
        deck = streamdecks[0]
        deck.open()
        deck.reset()
        deck.set_brightness(0)
        # Clear all keys to be sure it's black
        for key in range(deck.key_count()):
            deck.set_key_image(key, None)
        deck.close()
        print("Stream Deck brightness set to 0 and keys cleared.")
    else:
        print("No Stream Deck found.")
