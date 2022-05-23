# HUD: Head-up Display for Linux

[![build](https://github.com/fengye87/hud/actions/workflows/build.yml/badge.svg)](https://github.com/fengye87/hud/actions/workflows/build.yml)

HUD displays machine's IP and time without login process. It's very handy when you clone several VMs and want to quickly gather their IPs.

HUD is mainly tested on Rocky Linux 8, but it should work on most Linux distributions. Feel free to fire an issue if it's not.

## Screenshot

![screenshot](/screenshot.png)

## Installation

```bash
bash <(curl -s https://raw.githubusercontent.com/fengye87/hud/main/install.sh)
```

Or if you perfer using _wget_:

```bash
bash <(wget -qO- https://raw.githubusercontent.com/fengye87/hud/main/install.sh)
```

## Known Issue

Sometimes the screen's text may seem messed up, press Ctrl+C and everything will come back normal
