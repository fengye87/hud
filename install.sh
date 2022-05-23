#!/bin/bash

set -euo pipefail

function download {
  if which curl >/dev/null; then
    curl -Lo $2 $1
  else
    wget -O $2 $1
  fi
}

download https://github.com/fengye87/hud/releases/download/v1.0.0/hud /usr/local/bin/hud
chmod +x /usr/local/bin/hud
download https://raw.githubusercontent.com/fengye87/hud/main/hud.service /etc/systemd/system/hud.service

systemctl disable --now getty@tty1.service
systemctl enable --now hud.service
