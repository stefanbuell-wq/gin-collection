#!/bin/bash

# Health check script for monitoring
# Exit 0 if healthy, 1 if unhealthy

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

API_URL="${API_URL:-http://localhost:8080}"
TIMEOUT=5

# Check API health
echo -e "${GREEN}Checking API health...${NC}"
if curl -sf --max-time $TIMEOUT "${API_URL}/health" > /dev/null; then
    echo -e "${GREEN}✓ API is healthy${NC}"
else
    echo -e "${RED}✗ API health check failed${NC}"
    exit 1
fi

# Check API ready
echo -e "${GREEN}Checking API readiness...${NC}"
if curl -sf --max-time $TIMEOUT "${API_URL}/ready" > /dev/null; then
    echo -e "${GREEN}✓ API is ready${NC}"
else
    echo -e "${RED}✗ API readiness check failed${NC}"
    exit 1
fi

# Check database connectivity
echo -e "${GREEN}Checking database...${NC}"
if docker exec gin-collection-mysql mysqladmin ping -h localhost --silent 2>/dev/null; then
    echo -e "${GREEN}✓ Database is healthy${NC}"
else
    echo -e "${RED}✗ Database check failed${NC}"
    exit 1
fi

# Check Redis
echo -e "${GREEN}Checking Redis...${NC}"
if docker exec gin-collection-redis redis-cli ping 2>/dev/null | grep -q PONG; then
    echo -e "${GREEN}✓ Redis is healthy${NC}"
else
    echo -e "${RED}✗ Redis check failed${NC}"
    exit 1
fi

# Check disk space
echo -e "${GREEN}Checking disk space...${NC}"
DISK_USAGE=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
if [ "$DISK_USAGE" -gt 90 ]; then
    echo -e "${RED}✗ Disk usage is at ${DISK_USAGE}%${NC}"
    exit 1
elif [ "$DISK_USAGE" -gt 80 ]; then
    echo -e "${YELLOW}⚠️  Disk usage is at ${DISK_USAGE}%${NC}"
else
    echo -e "${GREEN}✓ Disk usage is at ${DISK_USAGE}%${NC}"
fi

# Check memory usage
echo -e "${GREEN}Checking memory usage...${NC}"
MEM_USAGE=$(free | grep Mem | awk '{print int($3/$2 * 100)}')
if [ "$MEM_USAGE" -gt 90 ]; then
    echo -e "${RED}✗ Memory usage is at ${MEM_USAGE}%${NC}"
    exit 1
elif [ "$MEM_USAGE" -gt 80 ]; then
    echo -e "${YELLOW}⚠️  Memory usage is at ${MEM_USAGE}%${NC}"
else
    echo -e "${GREEN}✓ Memory usage is at ${MEM_USAGE}%${NC}"
fi

# Check Docker containers
echo -e "${GREEN}Checking Docker containers...${NC}"
UNHEALTHY=$(docker ps --filter "health=unhealthy" -q | wc -l)
if [ "$UNHEALTHY" -gt 0 ]; then
    echo -e "${RED}✗ $UNHEALTHY unhealthy container(s)${NC}"
    docker ps --filter "health=unhealthy" --format "table {{.Names}}\t{{.Status}}"
    exit 1
else
    echo -e "${GREEN}✓ All containers are healthy${NC}"
fi

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}All health checks passed!${NC}"
echo -e "${GREEN}================================${NC}"

exit 0
