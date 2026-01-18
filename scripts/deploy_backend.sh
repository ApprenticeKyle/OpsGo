#!/bin/bash
set -e

# ==========================================
# FlowGo Backend Deployment Script
# ==========================================

PROJECT_DIR="/var/www/flowgo"
LOG_FILE="/var/log/flowgo_backend_deploy.log"

exec > >(tee -a $LOG_FILE) 2>&1
echo "-----------------------------------------------------"
echo "Backend Deployment started at $(date)"

echo "[Backend] Updating code..."
cd $PROJECT_DIR
git checkout main
git fetch --all
git reset --hard origin/main

echo "[Backend] Building binary..."
# Ensure correct Go environment
export GOTOOLCHAIN=local
export GOPROXY=https://goproxy.cn,direct
/usr/local/go/bin/go build -o flowgo-server ./cmd/server/main.go

echo "[Backend] Scheduling Service Restart (Async)..."
# Asynchronous restart to allow the Go application to save state properly prior to shutdown
nohup sh -c 'sleep 3; systemctl restart flowgo' > /dev/null 2>&1 &

echo "[Backend] Deployment Successful! (Service will restart in 3 seconds)"
echo "-----------------------------------------------------"
