#!/bin/bash

set -e  # Falla si algo sale mal

cd /root/neploy
git pull origin main
go build -o neploy ./cmd/app
sudo mv neploy /opt/neploy/
sudo chmod +x /opt/neploy/neploy
sudo systemctl restart neploy.service
