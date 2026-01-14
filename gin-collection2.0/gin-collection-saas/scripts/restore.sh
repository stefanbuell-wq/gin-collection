#!/bin/bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
BACKUP_DIR="/opt/gin-collection/backups"

if [ $# -eq 0 ]; then
    echo -e "${RED}Usage: $0 <backup_timestamp>${NC}"
    echo ""
    echo "Available backups:"
    ls -lh "$BACKUP_DIR" | grep -E "mysql_|redis_|app_" | awk '{print $9}'
    exit 1
fi

TIMESTAMP=$1

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root or with sudo${NC}"
    exit 1
fi

# Load environment variables
if [ -f "/opt/gin-collection/.env" ]; then
    export $(cat /opt/gin-collection/.env | grep -v '^#' | xargs)
else
    echo -e "${RED}Error: .env file not found${NC}"
    exit 1
fi

echo -e "${YELLOW}⚠️  WARNING: This will restore data from backup${NC}"
echo -e "${YELLOW}This operation will OVERWRITE existing data!${NC}"
echo ""
read -p "Are you sure you want to continue? (yes/no) " -r
if [ "$REPLY" != "yes" ]; then
    echo -e "${YELLOW}Restore cancelled${NC}"
    exit 0
fi

# Stop services
echo -e "${GREEN}[1/4] Stopping services...${NC}"
cd /opt/gin-collection
docker-compose -f docker-compose.prod.yml stop

echo -e "${GREEN}✓ Services stopped${NC}"

# Restore MySQL
MYSQL_BACKUP_FILE="$BACKUP_DIR/mysql_$TIMESTAMP.sql.gz"
if [ -f "$MYSQL_BACKUP_FILE" ]; then
    echo -e "${GREEN}[2/4] Restoring MySQL database...${NC}"

    if docker ps -a | grep -q gin-collection-mysql; then
        docker start gin-collection-mysql
        sleep 10

        gunzip < "$MYSQL_BACKUP_FILE" | docker exec -i gin-collection-mysql \
            mysql -u root -p"${DB_ROOT_PASSWORD}"

        echo -e "${GREEN}✓ MySQL database restored${NC}"
    else
        echo -e "${YELLOW}⚠️  MySQL container not found${NC}"
    fi
else
    echo -e "${RED}MySQL backup file not found: $MYSQL_BACKUP_FILE${NC}"
    exit 1
fi

# Restore Redis
REDIS_BACKUP_FILE="$BACKUP_DIR/redis_$TIMESTAMP.rdb"
if [ -f "$REDIS_BACKUP_FILE" ]; then
    echo -e "${GREEN}[3/4] Restoring Redis data...${NC}"

    if docker ps -a | grep -q gin-collection-redis; then
        docker start gin-collection-redis
        sleep 5
        docker stop gin-collection-redis

        docker cp "$REDIS_BACKUP_FILE" gin-collection-redis:/data/dump.rdb

        echo -e "${GREEN}✓ Redis data restored${NC}"
    else
        echo -e "${YELLOW}⚠️  Redis container not found, skipping${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Redis backup file not found, skipping${NC}"
fi

# Restore application configuration
APP_BACKUP_FILE="$BACKUP_DIR/app_config_$TIMESTAMP.tar.gz"
if [ -f "$APP_BACKUP_FILE" ]; then
    echo -e "${GREEN}[4/4] Restoring application configuration...${NC}"

    read -p "Do you want to restore application config? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        tar -xzf "$APP_BACKUP_FILE" -C /opt/gin-collection/
        echo -e "${GREEN}✓ Application configuration restored${NC}"
    else
        echo -e "${YELLOW}⚠️  Skipping app config restore${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  App config backup not found, skipping${NC}"
fi

# Restart services
echo -e "${GREEN}Restarting services...${NC}"
docker-compose -f docker-compose.prod.yml up -d

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Restore completed successfully!${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo "Waiting for services to be ready..."
sleep 30

docker-compose -f docker-compose.prod.yml ps
