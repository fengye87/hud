#!/bin/bash

wget -O /usr/local/bin/hud https://github.com/fengye87/hud/releases/download/v0.1.0/hud
chmod +x /usr/local/bin/hud
wget -O /etc/systemd/system/hud.service https://raw.githubusercontent.com/fengye87/hud/main/hud.service

systemctl disable --now getty@.service
systemctl enable --now hud.service
