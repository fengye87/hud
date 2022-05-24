# HUD: Head-up Display for Linux

[![build](https://github.com/fengye87/hud/actions/workflows/build.yml/badge.svg)](https://github.com/fengye87/hud/actions/workflows/build.yml)

HUD displays machine's IP and time without user login. It's very handy when you have several VMs and need their IPs every now and then.

HUD is mainly tested on Rocky Linux 8, but it should work on most Linux distributions. Feel free to fire an issue if it's not.

## Screenshot

![screenshot](/screenshot.png)

## Installation

```bash
bash <(curl -s https://raw.githubusercontent.com/fengye87/hud/main/install.sh)
```

Or if you perfer using `wget`:

```bash
bash <(wget -qO- https://raw.githubusercontent.com/fengye87/hud/main/install.sh)
```

## Protips

✨ Sometimes the screen may look messed up, press `Ctrl+C` and everything will fall back to normal.

✨ Try `Ctrl+Alt+F2` if you want to login (and `Ctrl+Alt+F1` to back to HUD).
