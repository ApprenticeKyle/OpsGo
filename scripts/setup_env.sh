#!/bin/bash
set -e

# Setup Env Script for FlowGo (Universal: Ubuntu/Debian/CentOS/AliLinux)

# Detect Package Manager
if command -v apt-get &> /dev/null; then
    PKG_MANAGER="apt-get"
    echo "Detected apt-get (Debian/Ubuntu)"
elif command -v yum &> /dev/null; then
    PKG_MANAGER="yum"
    echo "Detected yum (CentOS/RHEL/AliLinux)"
else
    echo "Error: Neither apt-get nor yum found. Manual installation required."
    exit 1
fi

echo "Updating system..."
if [ "$PKG_MANAGER" = "apt-get" ]; then
    apt-get update && apt-get upgrade -y
    apt-get install -y git curl wget nginx build-essential
elif [ "$PKG_MANAGER" = "yum" ]; then
    yum update -y
    # Install EPEL for Nginx if not present (often needed on CentOS)
    yum install -y epel-release || true
    yum install -y git curl wget nginx gcc make
fi

echo "Installing Database..."
if [ "$PKG_MANAGER" = "apt-get" ]; then
    apt-get install -y mysql-server
    systemctl start mysql
    systemctl enable mysql
elif [ "$PKG_MANAGER" = "yum" ]; then
    # Try installing MySQL, fallback to MariaDB if exact package not found
    # AliLinux often has mysql-server or dnf-mysql
    if yum list installed mysql-server &> /dev/null; then
        echo "MySQL already installed."
    else
        echo "Attempting to install MySQL/MariaDB..."
        # Quickest way: Mariadb is usually default in standard repos
        yum install -y mariadb-server || yum install -y mysql-server
    fi
    systemctl start mariadb || systemctl start mysqld
    systemctl enable mariadb || systemctl enable mysqld
fi

echo "Installing Go..."
# Remove old
rm -rf /usr/local/go
# Download
# Download
wget https://golang.google.cn/dl/go1.24.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
# Cleanup
rm go1.24.0.linux-amd64.tar.gz

# Configure PATH
if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
fi

# Configure Go Proxy
/usr/local/go/bin/go env -w GOPROXY=https://goproxy.cn,direct

echo "Installing Node.js..."
if [ "$PKG_MANAGER" = "apt-get" ]; then
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt-get install -y nodejs
elif [ "$PKG_MANAGER" = "yum" ]; then
    curl -fsSL https://rpm.nodesource.com/setup_20.x | bash -
    yum install -y nodejs
fi

echo "Configuring Nginx..."
mkdir -p /var/www/html /var/www/flowgo

# Determine Nginx User
if [ "$PKG_MANAGER" = "apt-get" ]; then
   NGINX_USER="www-data"
else
   NGINX_USER="nginx"
fi
# Ensure permissions
chown -R $NGINX_USER:$NGINX_USER /var/www/html

cat > /etc/nginx/nginx.conf <<EOF
user $NGINX_USER;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    
    server {
        listen 80;
        server_name _;
        root /var/www/html;
        index index.html;

        location / {
            try_files \$uri \$uri/ /index.html;
        }

        location /api/ {
            proxy_pass http://127.0.0.1:8080/api/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
        }
    }
}
EOF

systemctl restart nginx

echo "Creating Systemd Service for FlowGo..."
# Check for mysql service name again to correct 'After=' directive
DB_SERVICE="mysql.service"
if systemctl list-units --full -all | grep -Fq "mariadb.service"; then
    DB_SERVICE="mariadb.service"
elif systemctl list-units --full -all | grep -Fq "mysqld.service"; then
    DB_SERVICE="mysqld.service"
fi

cat > /etc/systemd/system/flowgo.service <<EOF
[Unit]
Description=FlowGo API Server
After=network.target $DB_SERVICE

[Service]
User=root
WorkingDirectory=/var/www/flowgo
ExecStart=/var/www/flowgo/flowgo-server
Restart=always
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload

echo "Setup Complete! Please re-login or run 'source ~/.bashrc' to update PATH."
echo "NOTE: If you installed MariaDB/MySQL, run 'mysql_secure_installation' manually if needed."
