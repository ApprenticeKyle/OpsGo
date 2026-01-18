#!/bin/bash
set -e

# ==========================================
# FlowGo Frontend Deployment Script
# ==========================================

FRONTEND_REPO_DIR="/var/www/flowboard_source"
FRONTEND_WEB_ROOT="/var/www/html"
LOG_FILE="/var/log/flowgo_frontend_deploy.log"

exec > >(tee -a $LOG_FILE) 2>&1
echo "-----------------------------------------------------"
echo "Frontend Deployment started at $(date)"

echo "[Frontend] Updating code..."
if [ ! -d "$FRONTEND_REPO_DIR" ]; then
    echo "Error: Frontend repo not found at $FRONTEND_REPO_DIR"
    echo "Please clone it first: git clone git@github.com:ApprenticeKyle/FlowBoard.git $FRONTEND_REPO_DIR"
    exit 1
fi

cd $FRONTEND_REPO_DIR
git checkout main
git fetch --all
git reset --hard origin/main

echo "[Frontend] Installing dependencies & Building..."
export PATH=$PATH:/usr/local/bin # Ensure node is in path
npm install --silent
npm run build

echo "[Frontend] Updating Nginx root..."
rm -rf "$FRONTEND_WEB_ROOT"/*
cp -r dist/* "$FRONTEND_WEB_ROOT"/

echo "[Frontend] Deployment Successful!"
echo "-----------------------------------------------------"
