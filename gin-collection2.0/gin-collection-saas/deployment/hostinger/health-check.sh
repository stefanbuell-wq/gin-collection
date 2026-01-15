#!/bin/bash
#===============================================================================
# Gin Collection - Health Check Script
# Für Cron: */5 * * * * /opt/gin-collection/.../health-check.sh
#===============================================================================

# Konfiguration
ALERT_EMAIL="${ALERT_EMAIL:-}"
SLACK_WEBHOOK="${SLACK_WEBHOOK:-}"
LOG_FILE="/var/log/gin-health.log"

# Farben (für Terminal)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Logging
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

send_alert() {
    local message="$1"
    local level="$2"

    # Email Alert
    if [ -n "$ALERT_EMAIL" ]; then
        echo "$message" | mail -s "[Gin Collection] $level Alert" "$ALERT_EMAIL"
    fi

    # Slack Alert
    if [ -n "$SLACK_WEBHOOK" ]; then
        curl -s -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"[$level] $message\"}" \
            "$SLACK_WEBHOOK" > /dev/null
    fi
}

check_service() {
    local name="$1"
    local url="$2"
    local expected_code="${3:-200}"

    response=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "$url" 2>/dev/null)

    if [ "$response" = "$expected_code" ]; then
        echo -e "${GREEN}✓${NC} $name: OK ($response)"
        return 0
    else
        echo -e "${RED}✗${NC} $name: FAILED (expected $expected_code, got $response)"
        return 1
    fi
}

check_container() {
    local name="$1"

    status=$(docker inspect -f '{{.State.Status}}' "$name" 2>/dev/null)

    if [ "$status" = "running" ]; then
        echo -e "${GREEN}✓${NC} Container $name: running"
        return 0
    else
        echo -e "${RED}✗${NC} Container $name: $status"
        return 1
    fi
}

check_disk() {
    local threshold="${1:-80}"
    local usage=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')

    if [ "$usage" -lt "$threshold" ]; then
        echo -e "${GREEN}✓${NC} Disk: ${usage}% used"
        return 0
    else
        echo -e "${RED}✗${NC} Disk: ${usage}% used (threshold: ${threshold}%)"
        return 1
    fi
}

check_memory() {
    local threshold="${1:-80}"
    local usage=$(free | grep Mem | awk '{printf("%.0f", $3/$2 * 100)}')

    if [ "$usage" -lt "$threshold" ]; then
        echo -e "${GREEN}✓${NC} Memory: ${usage}% used"
        return 0
    else
        echo -e "${YELLOW}⚠${NC} Memory: ${usage}% used (threshold: ${threshold}%)"
        return 1
    fi
}

#-------------------------------------------------------------------------------
# Main Health Check
#-------------------------------------------------------------------------------
echo "========================================"
echo "  Gin Collection Health Check"
echo "  $(date '+%Y-%m-%d %H:%M:%S')"
echo "========================================"
echo ""

ERRORS=0

# Container Status
echo "=== Containers ==="
check_container "gin-api" || ((ERRORS++))
check_container "gin-frontend" || ((ERRORS++))
check_container "gin-admin" || ((ERRORS++))
check_container "gin-mysql" || ((ERRORS++))
check_container "gin-redis" || ((ERRORS++))
echo ""

# HTTP Health Checks
echo "=== HTTP Endpoints ==="
check_service "API Health" "http://localhost:8080/health" || ((ERRORS++))
check_service "Frontend" "http://localhost:3000" || ((ERRORS++))
check_service "Admin Panel" "http://localhost:3001" || ((ERRORS++))
echo ""

# System Resources
echo "=== System Resources ==="
check_disk 80 || ((ERRORS++))
check_memory 80 || ((ERRORS++))
echo ""

# Database Connectivity
echo "=== Database ==="
if docker exec gin-mysql mysqladmin ping -h localhost -u root -p"$MYSQL_ROOT_PASSWORD" 2>/dev/null | grep -q "alive"; then
    echo -e "${GREEN}✓${NC} MySQL: responding"
else
    echo -e "${RED}✗${NC} MySQL: not responding"
    ((ERRORS++))
fi

if docker exec gin-redis redis-cli ping 2>/dev/null | grep -q "PONG"; then
    echo -e "${GREEN}✓${NC} Redis: responding"
else
    echo -e "${RED}✗${NC} Redis: not responding"
    ((ERRORS++))
fi
echo ""

#-------------------------------------------------------------------------------
# Summary
#-------------------------------------------------------------------------------
echo "========================================"
if [ $ERRORS -eq 0 ]; then
    echo -e "  Status: ${GREEN}ALL SYSTEMS OPERATIONAL${NC}"
    log "Health check passed"
else
    echo -e "  Status: ${RED}$ERRORS ISSUES DETECTED${NC}"
    log "Health check failed with $ERRORS errors"
    send_alert "Health check failed with $ERRORS errors" "CRITICAL"
fi
echo "========================================"

exit $ERRORS
