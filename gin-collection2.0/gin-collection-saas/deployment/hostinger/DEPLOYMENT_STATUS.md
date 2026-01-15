# GinVault - Deployment Status

**Datum:** 15. Januar 2026
**Server:** Hostinger VPS (srv1273559)
**Domain:** ginvault.cloud

---

## Aktueller Status

| Komponente | Status | Details |
|------------|--------|---------|
| VPS Server | ✅ Läuft | Ubuntu 22.04, Docker installiert |
| Docker | ✅ Läuft | Version 29.1.4 |
| MySQL | ✅ Healthy | Container: gin-collection-mysql |
| Redis | ✅ Healthy | Container: gin-collection-redis |
| API (Go) | ✅ Healthy | Container: gin-collection-api, Port 8080 |
| Frontend | ⚠️ Problem | Container läuft, aber JS lädt nicht |
| Admin | ⚠️ Problem | Container läuft, aber JS lädt nicht |
| Nginx | ✅ Läuft | Reverse Proxy konfiguriert |
| SSL | ✅ Aktiv | Let's Encrypt Zertifikat |

---

## Offene Probleme

### Frontend/Admin JavaScript lädt nicht
- **Symptom:** Seite zeigt nur "GinVault" Titel, keine App-Inhalte
- **Vermutung:** Nginx-Config im Container oder Asset-Pfade
- **Zu prüfen:**
  ```bash
  docker exec gin-collection-frontend cat /usr/share/nginx/html/index.html
  docker exec gin-collection-frontend cat /etc/nginx/nginx.conf
  curl -sI https://ginvault.cloud/assets/index-*.js
  ```

---

## Server-Zugangsdaten

```bash
# SSH Zugang
ssh root@[VPS_IP]

# Projektverzeichnis
cd ~/gin-collection/gin-collection2.0/gin-collection-saas

# Docker Status
docker compose ps

# Logs anzeigen
docker compose logs -f
docker compose logs api --tail 50
docker compose logs frontend --tail 50
```

---

## Installierte Komponenten

### Docker Container
```
gin-collection-mysql      mysql:8.0                    Port 3306
gin-collection-redis      redis:7-alpine               Port 6379
gin-collection-api        gin-collection-saas-api      Port 8080
gin-collection-frontend   gin-collection-saas-frontend Port 3000
gin-collection-admin      gin-collection-saas-admin    Port 3001
```

### Nginx Sites
- `/etc/nginx/sites-available/gin-collection`
- Proxy für ginvault.cloud → localhost:3000
- Proxy für admin.ginvault.cloud → localhost:3001
- API Proxy für /api/ → localhost:8080

### SSL Zertifikat
- Let's Encrypt via Certbot
- Domains: ginvault.cloud, www.ginvault.cloud, admin.ginvault.cloud
- Auto-Renewal aktiv

---

## Umgebungsvariablen (.env)

```bash
# Pfad: ~/gin-collection/gin-collection2.0/gin-collection-saas/.env

APP_ENV=production
APP_PORT=8080
APP_BASE_URL=https://ginvault.cloud

DB_HOST=mysql
DB_PORT=3306
DB_USER=gin_app
DB_PASSWORD=[generiert]
DB_NAME=gin_collection

REDIS_URL=redis:6379

JWT_SECRET=[generiert]
JWT_EXPIRATION=24h

CORS_ALLOWED_ORIGINS=https://ginvault.cloud,https://www.ginvault.cloud,https://admin.ginvault.cloud
```

---

## Nützliche Befehle

```bash
# Container neustarten
docker compose restart

# Alle Container stoppen
docker compose down

# Container mit Rebuild starten
docker compose up -d --build

# Logs live verfolgen
docker compose logs -f

# In Container Shell
docker exec -it gin-collection-api sh
docker exec -it gin-collection-frontend sh

# Nginx testen
sudo nginx -t
sudo systemctl reload nginx

# SSL erneuern
sudo certbot renew

# Datenbank Backup
docker exec gin-collection-mysql mysqldump -u root -p gin_collection > backup.sql
```

---

## Deployment-Schritte (für Neuinstallation)

1. **Server Setup:**
   ```bash
   cd ~/gin-collection/gin-collection2.0/gin-collection-saas
   sudo bash deployment/hostinger/setup-server.sh
   ```

2. **Umgebungsvariablen:**
   ```bash
   nano .env
   # Werte eintragen
   ```

3. **Container bauen und starten:**
   ```bash
   docker compose build --no-cache
   docker compose up -d
   ```

4. **Nginx konfigurieren:**
   ```bash
   sudo cp deployment/hostinger/nginx-ginvault.conf /etc/nginx/sites-available/gin-collection
   sudo ln -sf /etc/nginx/sites-available/gin-collection /etc/nginx/sites-enabled/
   sudo rm -f /etc/nginx/sites-enabled/default
   sudo nginx -t && sudo systemctl reload nginx
   ```

5. **SSL aktivieren:**
   ```bash
   sudo certbot --nginx -d ginvault.cloud -d www.ginvault.cloud -d admin.ginvault.cloud
   ```

---

## Git Repository

- **URL:** https://github.com/stefanbuell-wq/gin-collection
- **Branch:** master
- **Letzter Commit:** db10ea6 (feat: Add nginx config for ginvault.cloud)

### Aktualisieren:
```bash
cd ~/gin-collection/gin-collection2.0/gin-collection-saas
git pull
docker compose build --no-cache
docker compose up -d
```

---

## Nächste Schritte

1. [ ] Frontend JavaScript-Loading Problem beheben
2. [ ] Admin Panel testen
3. [ ] API Endpoints verifizieren
4. [ ] Datenbank-Migration prüfen
5. [ ] Admin-User anlegen
6. [ ] Test-Tenant erstellen
7. [ ] Monitoring einrichten (optional)
8. [ ] Backup-Cron einrichten

---

## Kontakt

- **Entwickler:** Stefan Buell
- **Email:** stefan.buell@gmail.com
