from StreamDeck.ImageHelpers import PILHelper
from PIL import Image, ImageDraw, ImageFont

def render_key(deck, label, bg="black", fg="white", bar_percent=None):
    image = PILHelper.create_image(deck)
    draw = ImageDraw.Draw(image)
    try: font = ImageFont.truetype("arial.ttf", 14)
    except: font = ImageFont.load_default()
    
    # Background and Border
    draw.rectangle([0, 0, image.width, image.height], fill=bg, outline="white")
    
    # Draw Progress Bar for Brightness
    if bar_percent is not None:
        bar_height = int((image.height * bar_percent) / 100)
        draw.rectangle([0, image.height - bar_height, image.width, image.height], fill="teal")
    
    # Text
    bbox = draw.textbbox((0, 0), label, font=font)
    x, y = (image.width-(bbox[2]-bbox[0]))//2, (image.height-(bbox[3]-bbox[1]))//2
    draw.text((x, y), label, font=font, fill=fg)
    
    return PILHelper.to_native_format(deck, image)

def render_home_view(manager, ui_render_func):
    mappings = {
        0: ("DASH\nHOME", "darkblue"),
        1: ("DASH\nNET", "darkblue"),
        2: ("DASH\nDOC", "darkblue"),
        4: (f"INBOX\n[{len(manager.state['notifications'])}]", "royalblue"),
        10: ("BRIGHT\n-", "black"),
        11: (f"BRIGHT\n{manager.state['brightness']}%", "black"),
        14: (f"TAP {manager.state['restart_count']}/5\nRESTART", "darkred" if manager.state['restart_count'] > 0 else "orange")
    }
    
    for key in range(manager.deck.key_count()):
        if key in mappings:
            label, color = mappings[key]
            # Special case for brightness rendering
            bar = manager.state['brightness'] if key == 11 else None
            manager.deck.set_key_image(key, ui_render_func(manager.deck, label, bg=color, bar_percent=bar))
        else:
            manager.deck.set_key_image(key, ui_render_func(manager.deck, "", bg="black"))

def render_blank_view(deck, ui_render_func):
    img = ui_render_func(deck, "", bg="black")
    for i in range(deck.key_count()):
        deck.set_key_image(i, img)
