#!/bin/sh
# Initialiseer TTY1
stty -F /dev/tty1 sane
echo -e "\033[9;0]\033[14;0]\033[H\033[J" > /dev/tty1

# Start Dashboard
./tmp/main < /dev/tty1 > /dev/tty1
