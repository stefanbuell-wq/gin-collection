#!/bin/bash
#===============================================================================
# Gin Collection - Automated Backup Script
# Für Cron: 0 3 * * * /opt/gin-collection/.../backup.sh
#===============================================================================

set -e

# Konfiguration
BACKUP_DIR="/opt/gin-collection/backups"
RETENTION_DAYS=7
DATE=$(date +%Y%m%d_%H%M%S)
PROJECT_DIR="/opt/gin-collection/gin-collection2.0/gin-collection-saas"

# Backup-Verzeichnis erstellen
mkdir -p "$BACKUP_DIR"

# In Projektverzeichnis wechseln
cd "$PROJECT_DIR"

# .env laden für Passwörter
source .env

echo "[$(date)] Starting backup..."

#-------------------------------------------------------------------------------
# MySQL Backup
#-------------------------------------------------------------------------------
echo "[$(date)] Backing up MySQL..."
docker compose exec -T mysql mysqldump \
    -u root \
    -p"$MYSQL_ROOT_PASSWORD" \
    --single-transaction \
    --routines \
    --triggers \
    gin_collection > "$BACKUP_DIR/db_$DATE.sql"

gzip "$BACKUP_DIR/db_$DATE.sql"
echo "[$(date)] MySQL backup: db_$DATE.sql.gz"

#-------------------------------------------------------------------------------
# Redis Backup (RDB Snapshot)
#-------------------------------------------------------------------------------
echo "[$(date)] Backing up Redis..."
docker compose exec -T redis redis-cli BGSAVE
sleep 2
docker cp gin-redis:/data/dump.rdb "$BACKUP_DIR/redis_$DATE.rdb" 2>/dev/null || true
echo "[$(date)] Redis backup: redis_$DATE.rdb"

#-------------------------------------------------------------------------------
# .env Backup (verschlüsselt)
#-------------------------------------------------------------------------------
echo "[$(date)] Backing up .env..."
cp .env "$BACKUP_DIR/env_$DATE.bak"
chmod 600 "$BACKUP_DIR/env_$DATE.bak"

#-------------------------------------------------------------------------------
# Alte Backups löschen
#-------------------------------------------------------------------------------
echo "[$(date)] Cleaning old backups (older than $RETENTION_DAYS days)..."
find "$BACKUP_DIR" -type f -mtime +$RETENTION_DAYS -delete

#-------------------------------------------------------------------------------
# Backup-Größe anzeigen
#-------------------------------------------------------------------------------
echo ""
echo "=== Backup Summary ==="
ls -lh "$BACKUP_DIR"/*$DATE* 2>/dev/null || echo "No files found"
echo ""
echo "Total backup size:"
du -sh "$BACKUP_DIR"

echo ""
echo "[$(date)] Backup completed successfully!"
