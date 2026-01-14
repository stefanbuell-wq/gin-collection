#!/bin/bash

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Gin Collection - Local Deployment Test${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Check Docker
echo -e "${GREEN}[1/8] Checking Docker installation...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker is not installed${NC}"
    echo "Please install Docker Desktop from: https://www.docker.com/products/docker-desktop"
    exit 1
fi

DOCKER_VERSION=$(docker --version)
echo -e "${GREEN}✓ Docker installed: $DOCKER_VERSION${NC}"

# Check Docker Compose
echo -e "${GREEN}[2/8] Checking Docker Compose...${NC}"
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
elif docker-compose --version &> /dev/null; then
    COMPOSE_CMD="docker-compose"
else
    echo -e "${RED}✗ Docker Compose not found${NC}"
    exit 1
fi

COMPOSE_VERSION=$($COMPOSE_CMD version)
echo -e "${GREEN}✓ Docker Compose: $COMPOSE_VERSION${NC}"

# Check .env file
echo -e "${GREEN}[3/8] Checking environment configuration...${NC}"
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  .env file not found, creating from template...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ .env file created${NC}"
else
    echo -e "${GREEN}✓ .env file exists${NC}"
fi

# Validate docker-compose.yml
echo -e "${GREEN}[4/8] Validating docker-compose configuration...${NC}"
if $COMPOSE_CMD config --quiet; then
    echo -e "${GREEN}✓ docker-compose.yml is valid${NC}"
else
    echo -e "${RED}✗ docker-compose.yml has errors${NC}"
    exit 1
fi

# Stop existing containers
echo -e "${GREEN}[5/8] Stopping existing containers...${NC}"
$COMPOSE_CMD down --remove-orphans 2>/dev/null || echo -e "${YELLOW}⚠️  No existing containers to stop${NC}"

# Pull images
echo -e "${GREEN}[6/8] Pulling required images...${NC}"
$COMPOSE_CMD pull

# Start services
echo -e "${GREEN}[7/8] Starting services...${NC}"
$COMPOSE_CMD up -d

# Wait for services to be ready
echo -e "${GREEN}[8/8] Waiting for services to be healthy...${NC}"
echo -e "${YELLOW}This may take 30-60 seconds...${NC}"

MAX_WAIT=120
WAITED=0
while [ $WAITED -lt $MAX_WAIT ]; do
    MYSQL_HEALTHY=$($COMPOSE_CMD ps mysql 2>/dev/null | grep -c "healthy" || echo "0")
    REDIS_HEALTHY=$($COMPOSE_CMD ps redis 2>/dev/null | grep -c "healthy" || echo "0")
    API_RUNNING=$($COMPOSE_CMD ps api 2>/dev/null | grep -c "Up" || echo "0")

    if [ "$MYSQL_HEALTHY" -eq "1" ] && [ "$REDIS_HEALTHY" -eq "1" ] && [ "$API_RUNNING" -eq "1" ]; then
        echo -e "${GREEN}✓ All services are healthy!${NC}"
        break
    fi

    echo -n "."
    sleep 2
    WAITED=$((WAITED + 2))
done

echo ""
echo ""

# Check service status
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Service Status${NC}"
echo -e "${BLUE}================================${NC}"
$COMPOSE_CMD ps

echo ""

# Health checks
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Health Checks${NC}"
echo -e "${BLUE}================================${NC}"

# MySQL Health
echo -n "MySQL: "
if docker exec gin-collection-mysql mysqladmin ping -h localhost -u root -pdev_root_password --silent 2>/dev/null; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${RED}✗ Unhealthy${NC}"
fi

# Redis Health
echo -n "Redis: "
if docker exec gin-collection-redis redis-cli ping 2>/dev/null | grep -q PONG; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${RED}✗ Unhealthy${NC}"
fi

# API Health
echo -n "API Health: "
if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Healthy${NC}"
else
    echo -e "${YELLOW}⚠️  Not responding (may still be starting)${NC}"
fi

# API Ready
echo -n "API Ready: "
if curl -sf http://localhost:8080/ready > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Ready${NC}"
else
    echo -e "${YELLOW}⚠️  Not ready (may still be starting)${NC}"
fi

echo ""

# Show API logs (last 10 lines)
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}API Logs (last 10 lines)${NC}"
echo -e "${BLUE}================================${NC}"
$COMPOSE_CMD logs --tail=10 api || echo -e "${YELLOW}API container not yet started${NC}"

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Deployment Test Complete${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo -e "${GREEN}Services are running!${NC}"
echo ""
echo "Access points:"
echo "  - Frontend: http://localhost:3000"
echo "  - API: http://localhost:8080"
echo "  - API Docs: http://localhost:8080/swagger"
echo "  - Prometheus: http://localhost:9090"
echo "  - Grafana: http://localhost:3001"
echo ""
echo "Useful commands:"
echo "  - View logs: $COMPOSE_CMD logs -f"
echo "  - Stop services: $COMPOSE_CMD down"
echo "  - Restart services: $COMPOSE_CMD restart"
echo "  - Check status: $COMPOSE_CMD ps"
echo ""

# Test API endpoint
echo -e "${BLUE}Testing API endpoint...${NC}"
echo -n "GET /health: "
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health 2>/dev/null || echo "FAILED")
if [ "$HEALTH_RESPONSE" != "FAILED" ]; then
    echo -e "${GREEN}✓ Success${NC}"
    echo "Response: $HEALTH_RESPONSE"
else
    echo -e "${YELLOW}⚠️  API not yet responding${NC}"
    echo "Try again in a few seconds: curl http://localhost:8080/health"
fi

echo ""
echo -e "${GREEN}Next steps:${NC}"
echo "1. Run migrations: ./scripts/migrate.sh up"
echo "2. Create first tenant: curl -X POST http://localhost:8080/api/v1/auth/register \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"tenant_name\":\"My Gin Collection\",\"subdomain\":\"mygins\",\"email\":\"admin@example.com\",\"password\":\"SecurePass123!\"}'"
echo "3. Open frontend: http://localhost:3000"
echo ""
