#!/bin/bash
#===============================================================================
# Gin Collection - Deployment Script
# Startet oder aktualisiert die Anwendung
#===============================================================================

set -e

# Farben
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_header() {
    echo -e "\n${BLUE}================================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}================================================================${NC}\n"
}

print_success() { echo -e "${GREEN}✓ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠ $1${NC}"; }
print_error() { echo -e "${RED}✗ $1${NC}"; }

# Zum Projektverzeichnis wechseln
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_DIR"

print_header "Gin Collection - Deployment"
echo "Projekt: $PROJECT_DIR"
echo ""

# Prüfe ob .env existiert
if [ ! -f ".env" ]; then
    print_error ".env Datei nicht gefunden!"
    echo "Erstelle .env aus Template: cp .env.example .env"
    exit 1
fi

# Optionen
ACTION=${1:-"start"}

case $ACTION in
    start|up)
        print_header "Starte Anwendung..."

        # Images bauen und starten
        docker compose build --no-cache
        docker compose up -d

        print_success "Container gestartet"
        echo ""
        docker compose ps
        ;;

    update|pull)
        print_header "Aktualisiere Anwendung..."

        # Neueste Änderungen holen
        git fetch origin
        git pull origin master

        # Container neu bauen und starten (Rolling Update)
        docker compose build
        docker compose up -d --force-recreate

        # Alte Images aufräumen
        docker image prune -f

        print_success "Update abgeschlossen"
        echo ""
        docker compose ps
        ;;

    stop|down)
        print_header "Stoppe Anwendung..."
        docker compose down
        print_success "Container gestoppt"
        ;;

    restart)
        print_header "Neustart..."
        docker compose restart
        print_success "Container neu gestartet"
        docker compose ps
        ;;

    logs)
        print_header "Logs (Strg+C zum Beenden)"
        docker compose logs -f --tail=100
        ;;

    logs-api)
        docker compose logs -f api --tail=100
        ;;

    logs-frontend)
        docker compose logs -f frontend --tail=100
        ;;

    status)
        print_header "Status"
        docker compose ps
        echo ""
        echo "Ressourcen:"
        docker stats --no-stream
        ;;

    shell-api)
        docker compose exec api sh
        ;;

    shell-db)
        docker compose exec mysql mysql -u gin_app -p gin_collection
        ;;

    shell-redis)
        docker compose exec redis redis-cli
        ;;

    backup)
        print_header "Erstelle Datenbank-Backup..."
        BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
        docker compose exec -T mysql mysqldump -u root -p"$MYSQL_ROOT_PASSWORD" gin_collection > "$BACKUP_FILE"
        gzip "$BACKUP_FILE"
        print_success "Backup erstellt: ${BACKUP_FILE}.gz"
        ;;

    clean)
        print_header "Aufräumen..."
        docker compose down -v --remove-orphans
        docker system prune -af
        print_success "System aufgeräumt"
        ;;

    health)
        print_header "Health Check"
        echo "API:"
        curl -s http://localhost:8080/health | jq . 2>/dev/null || curl -s http://localhost:8080/health
        echo ""
        echo "Frontend:"
        curl -sI http://localhost:3000 | head -1
        echo ""
        echo "Admin:"
        curl -sI http://localhost:3001 | head -1
        ;;

    *)
        echo "Gin Collection - Deployment Script"
        echo ""
        echo "Verwendung: $0 [BEFEHL]"
        echo ""
        echo "Befehle:"
        echo "  start, up      Startet alle Container"
        echo "  update, pull   Aktualisiert von Git und startet neu"
        echo "  stop, down     Stoppt alle Container"
        echo "  restart        Startet Container neu"
        echo "  logs           Zeigt alle Logs"
        echo "  logs-api       Zeigt nur API Logs"
        echo "  logs-frontend  Zeigt nur Frontend Logs"
        echo "  status         Zeigt Status und Ressourcen"
        echo "  shell-api      Shell im API Container"
        echo "  shell-db       MySQL Shell"
        echo "  shell-redis    Redis CLI"
        echo "  backup         Erstellt DB Backup"
        echo "  health         Health Check aller Services"
        echo "  clean          Entfernt alle Container und Images"
        echo ""
        ;;
esac
