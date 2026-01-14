#!/bin/bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
MIGRATIONS_DIR="./internal/infrastructure/database/migrations"
DIRECTION="${1:-up}"

echo -e "${GREEN}Running database migrations...${NC}"

# Load environment variables
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
elif [ -f "/opt/gin-collection/.env" ]; then
    export $(cat /opt/gin-collection/.env | grep -v '^#' | xargs)
else
    echo -e "${RED}Error: .env file not found${NC}"
    exit 1
fi

# Construct database URL
DB_URL="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?multiStatements=true&parseTime=true"

# Check if golang-migrate is installed
if ! command -v migrate &> /dev/null; then
    echo -e "${YELLOW}golang-migrate not found. Installing...${NC}"

    # Detect OS
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    if [ "$ARCH" = "x86_64" ]; then
        ARCH="amd64"
    elif [ "$ARCH" = "aarch64" ]; then
        ARCH="arm64"
    fi

    VERSION="v4.17.0"
    URL="https://github.com/golang-migrate/migrate/releases/download/${VERSION}/migrate.${OS}-${ARCH}.tar.gz"

    echo "Downloading from: $URL"
    curl -L "$URL" | tar xvz
    sudo mv migrate /usr/local/bin/
    sudo chmod +x /usr/local/bin/migrate

    echo -e "${GREEN}✓ golang-migrate installed${NC}"
fi

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}Migrations directory not found: $MIGRATIONS_DIR${NC}"
    exit 1
fi

# Run migrations
case "$DIRECTION" in
    up)
        echo -e "${GREEN}Running UP migrations...${NC}"
        migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" up
        echo -e "${GREEN}✓ Migrations completed${NC}"
        ;;
    down)
        echo -e "${YELLOW}Running DOWN migrations...${NC}"
        read -p "Are you sure you want to rollback? (yes/no) " -r
        if [ "$REPLY" = "yes" ]; then
            # Default to rolling back 1 migration
            STEPS="${2:-1}"
            migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" down "$STEPS"
            echo -e "${GREEN}✓ Rollback completed${NC}"
        else
            echo -e "${YELLOW}Rollback cancelled${NC}"
        fi
        ;;
    version)
        echo -e "${GREEN}Current migration version:${NC}"
        migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" version
        ;;
    force)
        if [ -z "$2" ]; then
            echo -e "${RED}Usage: $0 force <version>${NC}"
            exit 1
        fi
        echo -e "${YELLOW}Forcing version to $2...${NC}"
        migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" force "$2"
        echo -e "${GREEN}✓ Version forced${NC}"
        ;;
    *)
        echo -e "${RED}Unknown direction: $DIRECTION${NC}"
        echo "Usage: $0 {up|down|version|force}"
        exit 1
        ;;
esac

# Show current version
echo ""
echo -e "${GREEN}Current database version:${NC}"
migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" version || echo "No migrations applied yet"
