#!/bin/bash
#===============================================================================
# Gin Collection - Hostinger VPS Server Setup
# Führe dieses Script einmalig auf einem frischen Ubuntu 22.04 VPS aus
#===============================================================================

set -e

# Farben für Output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "\n${BLUE}================================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}================================================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Prüfe ob als root ausgeführt
if [ "$EUID" -ne 0 ]; then
    print_error "Bitte als root ausführen: sudo bash setup-server.sh"
    exit 1
fi

print_header "Gin Collection - Server Setup für Hostinger VPS"

# Variablen abfragen
read -p "Domain für Frontend (z.B. gin-collection.de): " DOMAIN
read -p "Admin-Subdomain (z.B. admin.gin-collection.de): " ADMIN_DOMAIN
read -p "Email für SSL-Zertifikat: " SSL_EMAIL
read -p "GitHub Repo URL: " GITHUB_REPO

# Defaults setzen falls leer
DOMAIN=${DOMAIN:-"gin-collection.de"}
ADMIN_DOMAIN=${ADMIN_DOMAIN:-"admin.gin-collection.de"}
GITHUB_REPO=${GITHUB_REPO:-"https://github.com/stefanbuell-wq/gin-collection.git"}

echo ""
print_warning "Konfiguration:"
echo "  Domain: $DOMAIN"
echo "  Admin: $ADMIN_DOMAIN"
echo "  Email: $SSL_EMAIL"
echo "  Repo: $GITHUB_REPO"
echo ""
read -p "Fortfahren? (j/n): " CONFIRM
if [ "$CONFIRM" != "j" ]; then
    echo "Abgebrochen."
    exit 0
fi

#-------------------------------------------------------------------------------
print_header "1/7 - System Update"
#-------------------------------------------------------------------------------
apt update && apt upgrade -y
print_success "System aktualisiert"

#-------------------------------------------------------------------------------
print_header "2/7 - Basis-Pakete installieren"
#-------------------------------------------------------------------------------
apt install -y \
    curl \
    wget \
    git \
    vim \
    htop \
    ufw \
    fail2ban \
    unzip \
    ca-certificates \
    gnupg \
    lsb-release

print_success "Basis-Pakete installiert"

#-------------------------------------------------------------------------------
print_header "3/7 - Docker installieren"
#-------------------------------------------------------------------------------

# Docker GPG Key
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

# Docker Repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Docker ohne sudo für aktuellen User
if [ -n "$SUDO_USER" ]; then
    usermod -aG docker $SUDO_USER
    print_success "Docker-Gruppe für $SUDO_USER konfiguriert"
fi

systemctl enable docker
systemctl start docker

print_success "Docker installiert und gestartet"

#-------------------------------------------------------------------------------
print_header "4/7 - Nginx installieren"
#-------------------------------------------------------------------------------
apt install -y nginx
systemctl enable nginx

print_success "Nginx installiert"

#-------------------------------------------------------------------------------
print_header "5/7 - Firewall konfigurieren"
#-------------------------------------------------------------------------------
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https
ufw --force enable

print_success "Firewall konfiguriert (SSH, HTTP, HTTPS erlaubt)"

#-------------------------------------------------------------------------------
print_header "6/7 - Fail2Ban konfigurieren"
#-------------------------------------------------------------------------------
cat > /etc/fail2ban/jail.local << 'EOF'
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3

[nginx-http-auth]
enabled = true

[nginx-limit-req]
enabled = true
EOF

systemctl enable fail2ban
systemctl restart fail2ban

print_success "Fail2Ban konfiguriert"

#-------------------------------------------------------------------------------
print_header "7/7 - Projektstruktur erstellen"
#-------------------------------------------------------------------------------

# App-Verzeichnis
APP_DIR="/opt/gin-collection"
mkdir -p $APP_DIR
cd $APP_DIR

# Repository klonen
git clone $GITHUB_REPO .
cd gin-collection2.0/gin-collection-saas

# .env Template erstellen
cat > .env << EOF
# =============================================================================
# Gin Collection - Production Environment
# =============================================================================

# Application
APP_ENV=production
APP_PORT=8080
APP_BASE_URL=https://$DOMAIN
LOG_LEVEL=info

# Database (Docker MySQL)
DB_HOST=mysql
DB_PORT=3306
DB_USER=gin_app
DB_PASSWORD=$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)
DB_NAME=gin_collection
MYSQL_ROOT_PASSWORD=$(openssl rand -base64 32 | tr -d '/+=' | head -c 32)

# Redis
REDIS_URL=redis:6379

# JWT (ÄNDERN!)
JWT_SECRET=$(openssl rand -base64 64 | tr -d '/+=' | head -c 64)
JWT_EXPIRATION=24h

# CORS
CORS_ALLOWED_ORIGINS=https://$DOMAIN,https://$ADMIN_DOMAIN

# S3 Storage (Optional - bei Bedarf ausfüllen)
S3_BUCKET=
S3_REGION=eu-central-1
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_ENDPOINT=

# PayPal (Optional - bei Bedarf ausfüllen)
PAYPAL_CLIENT_ID=
PAYPAL_CLIENT_SECRET=
PAYPAL_MODE=sandbox
PAYPAL_WEBHOOK_ID=

# Email/SMTP (Optional - bei Bedarf ausfüllen)
SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=noreply@$DOMAIN
SMTP_FROM_NAME=Gin Collection
SMTP_TLS=true

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_HOUR=100

# File Upload
MAX_UPLOAD_SIZE_MB=10

# Session
SESSION_TIMEOUT_HOURS=24
REFRESH_TOKEN_EXPIRATION_DAYS=30
EOF

chmod 600 .env
print_success "Projektstruktur erstellt"
print_warning ".env Datei mit Zufalls-Passwörtern generiert: $APP_DIR/gin-collection2.0/gin-collection-saas/.env"

#-------------------------------------------------------------------------------
# Nginx Konfiguration
#-------------------------------------------------------------------------------
cat > /etc/nginx/sites-available/gin-collection << EOF
# Gin Collection - Main Frontend
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /health {
        proxy_pass http://127.0.0.1:8080/health;
    }
}

# Gin Collection - Admin Frontend
server {
    listen 80;
    server_name $ADMIN_DOMAIN;

    location / {
        proxy_pass http://127.0.0.1:3001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }

    location /admin/api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

# Aktivieren
ln -sf /etc/nginx/sites-available/gin-collection /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default
nginx -t && systemctl reload nginx

print_success "Nginx konfiguriert für $DOMAIN und $ADMIN_DOMAIN"

#-------------------------------------------------------------------------------
# Abschluss
#-------------------------------------------------------------------------------
print_header "Setup abgeschlossen!"

echo -e "${GREEN}Nächste Schritte:${NC}"
echo ""
echo "1. DNS konfigurieren:"
echo "   $DOMAIN      -> $(curl -s ifconfig.me)"
echo "   $ADMIN_DOMAIN -> $(curl -s ifconfig.me)"
echo ""
echo "2. .env Datei anpassen:"
echo "   nano $APP_DIR/gin-collection2.0/gin-collection-saas/.env"
echo ""
echo "3. SSL-Zertifikat installieren:"
echo "   certbot --nginx -d $DOMAIN -d www.$DOMAIN -d $ADMIN_DOMAIN --email $SSL_EMAIL --agree-tos"
echo ""
echo "4. Anwendung starten:"
echo "   cd $APP_DIR/gin-collection2.0/gin-collection-saas"
echo "   ./deployment/hostinger/deploy.sh"
echo ""
echo -e "${YELLOW}Wichtig: Logge dich neu ein damit Docker-Gruppe aktiv wird!${NC}"
