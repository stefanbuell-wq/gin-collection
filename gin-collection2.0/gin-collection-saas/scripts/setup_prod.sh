#!/bin/bash

set -e

echo "================================"
echo "Gin Collection Production Setup"
echo "================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root or with sudo${NC}"
    exit 1
fi

# Check system requirements
echo -e "${GREEN}[1/10] Checking system requirements...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}docker-compose not found. Installing...${NC}"
    curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

echo -e "${GREEN}✓ Docker and docker-compose installed${NC}"

# Create directory structure
echo -e "${GREEN}[2/10] Creating directory structure...${NC}"
mkdir -p /opt/gin-collection
mkdir -p /opt/gin-collection/backups
mkdir -p /opt/gin-collection/logs
mkdir -p /var/log/gin-collection

echo -e "${GREEN}✓ Directories created${NC}"

# Copy application files
echo -e "${GREEN}[3/10] Copying application files...${NC}"
if [ ! -f "docker-compose.prod.yml" ]; then
    echo -e "${RED}docker-compose.prod.yml not found. Please run this script from the project root.${NC}"
    exit 1
fi

cp docker-compose.prod.yml /opt/gin-collection/
cp -r monitoring /opt/gin-collection/

echo -e "${GREEN}✓ Files copied${NC}"

# Setup environment file
echo -e "${GREEN}[4/10] Setting up environment variables...${NC}"
if [ ! -f "/opt/gin-collection/.env" ]; then
    echo -e "${YELLOW}No .env file found. Creating from template...${NC}"
    cp .env.example /opt/gin-collection/.env
    echo -e "${YELLOW}⚠️  Please edit /opt/gin-collection/.env with production values${NC}"
    read -p "Press enter to continue after editing .env file..."
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

# Generate JWT secret if not set
if grep -q "change_me_to_256_bit_random_secret" /opt/gin-collection/.env; then
    echo -e "${YELLOW}Generating JWT secret...${NC}"
    JWT_SECRET=$(openssl rand -base64 64 | tr -d '\n')
    sed -i "s|JWT_SECRET=.*|JWT_SECRET=${JWT_SECRET}|" /opt/gin-collection/.env
    echo -e "${GREEN}✓ JWT secret generated${NC}"
fi

# Setup SSL certificates
echo -e "${GREEN}[5/10] Setting up SSL certificates...${NC}"
read -p "Do you want to setup Let's Encrypt SSL? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if ! command -v certbot &> /dev/null; then
        echo -e "${YELLOW}Installing certbot...${NC}"
        apt-get update
        apt-get install -y certbot
    fi

    read -p "Enter domain name: " DOMAIN
    certbot certonly --standalone -d $DOMAIN

    mkdir -p /opt/gin-collection/ssl
    ln -sf /etc/letsencrypt/live/$DOMAIN/fullchain.pem /opt/gin-collection/ssl/cert.pem
    ln -sf /etc/letsencrypt/live/$DOMAIN/privkey.pem /opt/gin-collection/ssl/key.pem

    echo -e "${GREEN}✓ SSL certificates setup${NC}"
else
    echo -e "${YELLOW}⚠️  Skipping SSL setup. You can set it up later.${NC}"
fi

# Setup firewall
echo -e "${GREEN}[6/10] Configuring firewall...${NC}"
if command -v ufw &> /dev/null; then
    ufw allow 22/tcp
    ufw allow 80/tcp
    ufw allow 443/tcp
    ufw --force enable
    echo -e "${GREEN}✓ Firewall configured${NC}"
else
    echo -e "${YELLOW}⚠️  UFW not found. Please configure firewall manually.${NC}"
fi

# Pull Docker images
echo -e "${GREEN}[7/10] Pulling Docker images...${NC}"
cd /opt/gin-collection
docker-compose -f docker-compose.prod.yml pull

echo -e "${GREEN}✓ Docker images pulled${NC}"

# Run database migrations
echo -e "${GREEN}[8/10] Running database migrations...${NC}"
read -p "Do you want to run database migrations? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Assume migrations are run via the API container on startup
    echo -e "${YELLOW}Migrations will run automatically on first startup${NC}"
else
    echo -e "${YELLOW}⚠️  Skipping migrations${NC}"
fi

# Setup systemd service
echo -e "${GREEN}[9/10] Setting up systemd service...${NC}"
cat > /etc/systemd/system/gin-collection.service <<EOF
[Unit]
Description=Gin Collection SaaS
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/gin-collection
ExecStart=/usr/local/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.prod.yml down
ExecReload=/usr/local/bin/docker-compose -f docker-compose.prod.yml restart

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable gin-collection

echo -e "${GREEN}✓ Systemd service created${NC}"

# Setup log rotation
echo -e "${GREEN}[10/10] Setting up log rotation...${NC}"
cat > /etc/logrotate.d/gin-collection <<EOF
/var/log/gin-collection/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 root root
    sharedscripts
}
EOF

echo -e "${GREEN}✓ Log rotation configured${NC}"

# Start services
echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Setup completed successfully!${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review /opt/gin-collection/.env and update with production values"
echo "2. Start the application: systemctl start gin-collection"
echo "3. Check status: systemctl status gin-collection"
echo "4. View logs: docker-compose -f /opt/gin-collection/docker-compose.prod.yml logs -f"
echo "5. Access Grafana at http://your-server:3001 (default password in .env)"
echo ""
echo -e "${YELLOW}Important:${NC}"
echo "- Setup your DNS records to point to this server"
echo "- Configure your PayPal webhook URL"
echo "- Setup S3 bucket and update credentials"
echo "- Review security settings and firewall rules"
echo ""

read -p "Do you want to start the application now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    systemctl start gin-collection
    echo -e "${GREEN}✓ Application started${NC}"
    echo ""
    echo "Waiting for services to be ready..."
    sleep 30
    docker-compose -f /opt/gin-collection/docker-compose.prod.yml ps
else
    echo -e "${YELLOW}You can start the application later with: systemctl start gin-collection${NC}"
fi
