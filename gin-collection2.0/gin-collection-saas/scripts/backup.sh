#!/bin/bash

set -e

# Configuration
BACKUP_DIR="/opt/gin-collection/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Starting backup process...${NC}"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Load environment variables
if [ -f "/opt/gin-collection/.env" ]; then
    export $(cat /opt/gin-collection/.env | grep -v '^#' | xargs)
else
    echo -e "${RED}Error: .env file not found${NC}"
    exit 1
fi

# Backup MySQL database
echo -e "${GREEN}[1/3] Backing up MySQL database...${NC}"
MYSQL_BACKUP_FILE="$BACKUP_DIR/mysql_$TIMESTAMP.sql.gz"

# If using Docker MySQL
if docker ps | grep -q gin-collection-mysql; then
    docker exec gin-collection-mysql mysqldump \
        -u root \
        -p"${DB_ROOT_PASSWORD}" \
        --all-databases \
        --single-transaction \
        --routines \
        --triggers \
        --events \
        | gzip > "$MYSQL_BACKUP_FILE"
else
    # If using external MySQL
    mysqldump \
        -h "${DB_HOST}" \
        -u "${DB_USER}" \
        -p"${DB_PASSWORD}" \
        "${DB_NAME}" \
        --single-transaction \
        --routines \
        --triggers \
        | gzip > "$MYSQL_BACKUP_FILE"
fi

echo -e "${GREEN}✓ MySQL backup saved to $MYSQL_BACKUP_FILE${NC}"

# Backup Redis data
echo -e "${GREEN}[2/3] Backing up Redis data...${NC}"
REDIS_BACKUP_FILE="$BACKUP_DIR/redis_$TIMESTAMP.rdb"

if docker ps | grep -q gin-collection-redis; then
    docker exec gin-collection-redis redis-cli SAVE
    docker cp gin-collection-redis:/data/dump.rdb "$REDIS_BACKUP_FILE"
    echo -e "${GREEN}✓ Redis backup saved to $REDIS_BACKUP_FILE${NC}"
else
    echo -e "${YELLOW}⚠️  Redis container not found, skipping Redis backup${NC}"
fi

# Backup application files and configuration
echo -e "${GREEN}[3/3] Backing up application files...${NC}"
APP_BACKUP_FILE="$BACKUP_DIR/app_config_$TIMESTAMP.tar.gz"

tar -czf "$APP_BACKUP_FILE" \
    -C /opt/gin-collection \
    .env \
    docker-compose.prod.yml \
    monitoring \
    2>/dev/null || echo -e "${YELLOW}⚠️  Some files may not exist${NC}"

echo -e "${GREEN}✓ Application backup saved to $APP_BACKUP_FILE${NC}"

# Calculate backup size
BACKUP_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
echo -e "${GREEN}Total backup size: $BACKUP_SIZE${NC}"

# Clean up old backups
echo -e "${GREEN}Cleaning up backups older than $RETENTION_DAYS days...${NC}"
find "$BACKUP_DIR" -name "*.gz" -type f -mtime +$RETENTION_DAYS -delete
find "$BACKUP_DIR" -name "*.rdb" -type f -mtime +$RETENTION_DAYS -delete
find "$BACKUP_DIR" -name "*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete

echo -e "${GREEN}✓ Cleanup complete${NC}"

# Optional: Upload to S3
if [ -n "$S3_BACKUP_BUCKET" ]; then
    echo -e "${GREEN}Uploading backups to S3...${NC}"

    if command -v aws &> /dev/null; then
        aws s3 cp "$MYSQL_BACKUP_FILE" "s3://${S3_BACKUP_BUCKET}/backups/mysql/"
        aws s3 cp "$REDIS_BACKUP_FILE" "s3://${S3_BACKUP_BUCKET}/backups/redis/" || true
        aws s3 cp "$APP_BACKUP_FILE" "s3://${S3_BACKUP_BUCKET}/backups/app/"
        echo -e "${GREEN}✓ Backups uploaded to S3${NC}"
    else
        echo -e "${YELLOW}⚠️  AWS CLI not installed, skipping S3 upload${NC}"
    fi
fi

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Backup completed successfully!${NC}"
echo -e "${GREEN}================================${NC}"
echo ""
echo "Backup files:"
echo "  - MySQL: $MYSQL_BACKUP_FILE"
echo "  - Redis: $REDIS_BACKUP_FILE"
echo "  - App Config: $APP_BACKUP_FILE"
echo ""
