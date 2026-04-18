import os

def get_cpu_temp():
    try:
        # Assuming run with host /sys mounted if necessary, or privileged mode works
        with open('/sys/class/thermal/thermal_zone0/temp', 'r') as f:
            temp_c = int(f.read().strip()) / 1000.0
        return f"{temp_c:.1f}°C"
    except Exception as e:
        return "N/A"

def get_cpu_usage():
    # Simple reading from /proc/stat
    try:
        with open('/proc/stat', 'r') as f:
            lines = f.readlines()
        for line in lines:
            if line.startswith('cpu '):
                parts = line.split()
                # user, nice, system, idle, iowait, irq, softirq
                idle = float(parts[4])
                total = sum(float(p) for p in parts[1:8])
                return idle, total
    except Exception:
        pass
    return 0, 0

last_idle, last_total = 0, 0

def get_cpu_percent():
    global last_idle, last_total
    idle, total = get_cpu_usage()
    if total == 0:
        return "N/A"
    idle_delta = idle - last_idle
    total_delta = total - last_total
    last_idle, last_total = idle, total
    
    if total_delta == 0:
        return "0%"
        
    usage = 100.0 * (1.0 - idle_delta / total_delta)
    return f"{usage:.1f}%"
