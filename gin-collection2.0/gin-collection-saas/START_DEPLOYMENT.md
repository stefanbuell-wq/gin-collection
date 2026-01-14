# 🚀 Deployment Starten - Schnellanleitung

## ✅ Docker ist installiert!

**Ihre Installation:**
- Docker Version: 29.1.3
- Docker Compose: v5.0.0-desktop.1

## 📋 Nächste Schritte

### Schritt 1: Docker Desktop starten

**WICHTIG:** Docker Desktop muss laufen, damit die Container gestartet werden können.

1. **Docker Desktop öffnen:**
   - Suchen Sie nach "Docker Desktop" im Windows-Startmenü
   - Oder doppelklicken Sie auf das Docker-Icon auf dem Desktop
   - Warten Sie, bis Docker Desktop vollständig gestartet ist (ca. 30-60 Sekunden)

2. **Status prüfen:**
   - Das Docker-Icon in der Taskleiste sollte grün sein
   - Es sollte "Docker Desktop is running" anzeigen

### Schritt 2: Deployment testen (Option 1 - Einfach)

**Mit dem Batch-Script (Windows):**

1. Öffnen Sie den Ordner im Explorer:
   ```
   E:\Web-Projekte\Gin-App\gin-collection2.0\gin-collection-saas\
   ```

2. Doppelklicken Sie auf:
   ```
   test-docker.bat
   ```

3. Das Script wird automatisch:
   - ✅ Docker Installation prüfen
   - ✅ .env Datei erstellen
   - ✅ Services starten (MySQL, Redis, API, Frontend)
   - ✅ Health Checks durchführen
   - ✅ Status anzeigen

### Schritt 3: Deployment testen (Option 2 - Manuell)

**Mit Docker Compose:**

1. Öffnen Sie PowerShell oder CMD

2. Navigieren Sie zum Projektordner:
   ```powershell
   cd E:\Web-Projekte\Gin-App\gin-collection2.0\gin-collection-saas
   ```

3. Starten Sie die Services:
   ```powershell
   docker compose up -d
   ```

4. Warten Sie 30-60 Sekunden bis alle Services gestartet sind

5. Prüfen Sie den Status:
   ```powershell
   docker compose ps
   ```

### Schritt 4: Services prüfen

**Erwartetes Ergebnis:**

```
NAME                        STATUS              PORTS
gin-collection-api          Up (healthy)        0.0.0.0:8080->8080/tcp
gin-collection-frontend     Up (healthy)        0.0.0.0:3000->8080/tcp
gin-collection-mysql        Up (healthy)        0.0.0.0:3306->3306/tcp
gin-collection-redis        Up (healthy)        0.0.0.0:6379->6379/tcp
```

### Schritt 5: Anwendung testen

**Im Browser öffnen:**

1. **Frontend:**
   - URL: http://localhost:3000
   - Sie sollten die Login-Seite sehen

2. **API Health Check:**
   - URL: http://localhost:8080/health
   - Antwort sollte sein: `{"status":"ok"}`

3. **API Ready Check:**
   - URL: http://localhost:8080/ready
   - Prüft ob DB + Redis verfügbar sind

## 🔍 Troubleshooting

### Problem: "Docker daemon not running"

**Lösung:**
1. Docker Desktop starten
2. Warten bis das Icon grün ist
3. Erneut versuchen

### Problem: "Port already in use"

**Lösung:**
Prüfen Sie, ob die Ports bereits belegt sind:

```powershell
# Port 8080 (API)
netstat -ano | findstr "8080"

# Port 3000 (Frontend)
netstat -ano | findstr "3000"

# Port 3306 (MySQL)
netstat -ano | findstr "3306"
```

Wenn Ports belegt sind:
1. Stoppen Sie die andere Anwendung
2. Oder ändern Sie die Ports in `.env`

### Problem: Services starten nicht

**Lösung:**
Logs anschauen:

```powershell
# Alle Logs
docker compose logs

# Nur API Logs
docker compose logs api

# Nur MySQL Logs
docker compose logs mysql

# Live Logs (folgen)
docker compose logs -f
```

### Problem: "Cannot connect to database"

**Lösung:**
1. Warten Sie 30-60 Sekunden (MySQL braucht Zeit zum Starten)
2. Prüfen Sie MySQL Status:
   ```powershell
   docker compose ps mysql
   ```
3. Wenn "unhealthy" - Logs prüfen:
   ```powershell
   docker compose logs mysql
   ```

## 📝 Nützliche Befehle

### Services verwalten

```powershell
# Status anzeigen
docker compose ps

# Services starten
docker compose up -d

# Services stoppen
docker compose down

# Services neu starten
docker compose restart

# Logs anzeigen
docker compose logs -f

# Nur einen Service neu starten
docker compose restart api
```

### Container Shell öffnen

```powershell
# API Container
docker exec -it gin-collection-api sh

# MySQL Container
docker exec -it gin-collection-mysql mysql -u root -pdev_root_password gin_collection

# Redis Container
docker exec -it gin-collection-redis redis-cli
```

### Aufräumen

```powershell
# Alles stoppen
docker compose down

# Stoppen + Volumes löschen (= Datenbank zurücksetzen)
docker compose down -v

# Stoppen + Images löschen
docker compose down --rmi all
```

## 🎯 Nach erfolgreichem Start

### 1. Database Migrations ausführen

```powershell
# Linux/Mac
./scripts/migrate.sh up

# Windows (Git Bash)
bash scripts/migrate.sh up

# Oder manuell im MySQL Container
docker exec -it gin-collection-mysql mysql -u root -pdev_root_password gin_collection < internal/infrastructure/database/migrations/001_initial_schema.up.sql
```

### 2. Ersten Tenant erstellen

**Mit curl (PowerShell):**

```powershell
curl -X POST http://localhost:8080/api/v1/auth/register `
  -H "Content-Type: application/json" `
  -d '{
    "tenant_name": "My Gin Collection",
    "subdomain": "mygins",
    "email": "admin@example.com",
    "password": "SecurePassword123!"
  }'
```

**Oder mit Postman/Insomnia:**
- Method: POST
- URL: http://localhost:8080/api/v1/auth/register
- Body (JSON):
  ```json
  {
    "tenant_name": "My Gin Collection",
    "subdomain": "mygins",
    "email": "admin@example.com",
    "password": "SecurePassword123!"
  }
  ```

### 3. Im Frontend anmelden

1. Öffnen Sie: http://localhost:3000
2. Klicken Sie auf "Login"
3. Geben Sie die Credentials ein:
   - Email: `admin@example.com`
   - Password: `SecurePassword123!`

### 4. Ersten Gin anlegen

Nach dem Login können Sie über das Frontend Ihren ersten Gin hinzufügen!

## 🎓 Weiterführende Dokumentation

- **QUICKSTART.md** - Detaillierte Schnellanleitung
- **docs/DEPLOYMENT.md** - Umfassender Deployment Guide
- **docs/PROGRESS_REPORT.md** - Projekt-Übersicht
- **tests/README.md** - Testing Guide

## ✅ Checkliste

- [ ] Docker Desktop gestartet
- [ ] Docker Icon ist grün
- [ ] `test-docker.bat` ausgeführt
- [ ] Alle 4 Services laufen
- [ ] Health Checks sind grün
- [ ] http://localhost:8080/health zeigt "ok"
- [ ] http://localhost:3000 zeigt Login-Seite
- [ ] Migrations ausgeführt
- [ ] Erster Tenant erstellt
- [ ] Im Frontend eingeloggt
- [ ] Erster Gin angelegt

---

## 🆘 Support

Bei Problemen:
1. Logs prüfen: `docker compose logs`
2. Status prüfen: `docker compose ps`
3. Health Check Script: `bash scripts/health_check.sh`
4. Dokumentation: `docs/DEPLOYMENT.md`

**Viel Erfolg! 🚀**
