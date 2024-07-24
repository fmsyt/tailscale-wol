#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cat << EOF > /etc/systemd/system/tailscale-wol.service
[Unit]
Description=Start wol service on tailscaled
After=tailscaled.service

[Service]
ExecStart=${SCRIPT_DIR}/wol
KillMode=process
RemainAfterExit=yes
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal
Environment="HOME=/root"
Environment="XDG_CONFIG_HOME=/root/.config"

[Install]
WantedBy=multi-user.target
EOF
