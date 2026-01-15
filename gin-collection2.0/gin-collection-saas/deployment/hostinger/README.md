# Gin Collection - Hostinger VPS Deployment

Vollständige Anleitung zum Deployment der Gin Collection SaaS Platform auf einem Hostinger VPS.

## Voraussetzungen

### Hostinger VPS Anforderungen

| Spezifikation | Minimum | Empfohlen |
|---------------|---------|-----------|
| RAM | 2 GB | 4 GB |
| CPU | 1 vCore | 2 vCores |
| Storage | 20 GB SSD | 40 GB SSD |
| OS | Ubuntu 22.04 | Ubuntu 22.04 |
| Bandbreite | 1 TB | Unlimited |

### Empfohlener Hostinger Plan

- **VPS 2** (~9€/Monat) für kleine bis mittlere Last
- **VPS 4** (~16€/Monat) für größere Nutzerzahlen

---

## Schnellstart

### 1. VPS bestellen und SSH-Zugang einrichten

```bash
# Von deinem lokalen Rechner
ssh root@DEINE_VPS_IP
```

### 2. Setup-Script ausführen

```bash
# Auf dem VPS
curl -fsSL https://raw.githubusercontent.com/stefanbuell-wq/gin-collection/master/gin-collection2.0/gin-collection-saas/deployment/hostinger/setup-server.sh -o setup.sh
chmod +x setup.sh
sudo bash setup.sh
```

Das Script fragt nach:
- Domain (z.B. `gin-collection.de`)
- Admin-Subdomain (z.B. `admin.gin-collection.de`)
- E-Mail für SSL-Zertifikat
- GitHub Repository URL

### 3. DNS konfigurieren

In der Hostinger DNS-Verwaltung oder deinem Domain-Provider:

| Typ | Name | Wert | TTL |
|-----|------|------|-----|
| A | @ | DEINE_VPS_IP | 3600 |
| A | www | DEINE_VPS_IP | 3600 |
| A | admin | DEINE_VPS_IP | 3600 |

### 4. SSL-Zertifikat installieren

```bash
sudo certbot --nginx \
  -d gin-collection.de \
  -d www.gin-collection.de \
  -d admin.gin-collection.de \
  --email deine@email.de \
  --agree-tos
```

### 5. Umgebungsvariablen anpassen

```bash
cd /opt/gin-collection/gin-collection2.0/gin-collection-saas
nano .env
```

Wichtige Einstellungen:
- `JWT_SECRET` - Bereits generiert, bei Bedarf ändern
- `DB_PASSWORD` - Bereits generiert
- `S3_*` - Für Foto-Upload konfigurieren
- `PAYPAL_*` - Für Zahlungen konfigurieren
- `SMTP_*` - Für E-Mail-Versand konfigurieren

### 6. Anwendung starten

```bash
./deployment/hostinger/deploy.sh start
```

---

## Deployment Script Befehle

```bash
# Starten
./deploy.sh start

# Aktualisieren (Git Pull + Neustart)
./deploy.sh update

# Stoppen
./deploy.sh stop

# Neustart
./deploy.sh restart

# Logs anzeigen
./deploy.sh logs
./deploy.sh logs-api
./deploy.sh logs-frontend

# Status prüfen
./deploy.sh status

# Health Check
./deploy.sh health

# Datenbank-Backup
./deploy.sh backup

# Shell-Zugang
./deploy.sh shell-api
./deploy.sh shell-db
./deploy.sh shell-redis

# Aufräumen (Vorsicht!)
./deploy.sh clean
```

---

## Systemd Service (Autostart)

Für automatischen Start nach Server-Neustart:

```bash
# Service-Datei kopieren
sudo cp /opt/gin-collection/gin-collection2.0/gin-collection-saas/deployment/hostinger/gin-collection.service /etc/systemd/system/

# Service aktivieren
sudo systemctl daemon-reload
sudo systemctl enable gin-collection

# Service steuern
sudo systemctl start gin-collection
sudo systemctl stop gin-collection
sudo systemctl restart gin-collection
sudo systemctl status gin-collection
```

---

## Architektur

```
                    ┌─────────────────┐
                    │   Cloudflare    │ (Optional)
                    │   oder direkt   │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │     Nginx       │ :80/:443
                    │  Reverse Proxy  │
                    └────────┬────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│   Frontend    │   │  Admin Panel  │   │    API        │
│   :3000       │   │   :3001       │   │   :8080       │
│   (React)     │   │   (React)     │   │   (Go)        │
└───────────────┘   └───────────────┘   └───────┬───────┘
                                                │
                              ┌─────────────────┼─────────────────┐
                              │                                   │
                              ▼                                   ▼
                    ┌─────────────────┐                 ┌─────────────────┐
                    │     MySQL       │                 │     Redis       │
                    │     :3306       │                 │     :6379       │
                    └─────────────────┘                 └─────────────────┘
```

---

## Ports

| Service | Container Port | Host Port | Extern |
|---------|----------------|-----------|--------|
| Frontend | 8080 | 3000 | Nein (via Nginx) |
| Admin | 8080 | 3001 | Nein (via Nginx) |
| API | 8080 | 8080 | Nein (via Nginx) |
| MySQL | 3306 | - | Nein |
| Redis | 6379 | - | Nein |
| Nginx | 80, 443 | 80, 443 | Ja |

---

## Backup & Restore

### Automatisches Backup (Cron)

```bash
# Crontab bearbeiten
crontab -e

# Tägliches Backup um 3:00 Uhr
0 3 * * * /opt/gin-collection/gin-collection2.0/gin-collection-saas/deployment/hostinger/deploy.sh backup >> /var/log/gin-backup.log 2>&1
```

### Manuelles Backup

```bash
./deploy.sh backup
# Erstellt: backup_20240115_030000.sql.gz
```

### Restore

```bash
# Backup entpacken
gunzip backup_20240115_030000.sql.gz

# In Container importieren
docker exec -i gin-mysql mysql -u gin_app -p gin_collection < backup_20240115_030000.sql
```

---

## Monitoring

### Ressourcen prüfen

```bash
# Container Status
docker stats

# Disk Space
df -h

# Memory
free -m

# Logs
./deploy.sh logs
```

### Health Endpoints

- Frontend: `https://gin-collection.de/health`
- Admin: `https://admin.gin-collection.de/health`
- API: `https://gin-collection.de/api/health`

---

## Troubleshooting

### Container startet nicht

```bash
# Logs prüfen
docker compose logs api
docker compose logs mysql

# Container manuell starten für Fehlerdetails
docker compose up api
```

### Datenbank-Verbindungsfehler

```bash
# MySQL Status prüfen
docker compose exec mysql mysqladmin -u root -p status

# Logs prüfen
docker compose logs mysql
```

### SSL-Probleme

```bash
# Zertifikat erneuern
sudo certbot renew

# Nginx Konfiguration testen
sudo nginx -t

# Nginx neu laden
sudo systemctl reload nginx
```

### Speicherplatz voll

```bash
# Docker aufräumen
docker system prune -af
docker volume prune -f

# Alte Logs löschen
truncate -s 0 /var/lib/docker/containers/*/*-json.log
```

---

## Updates

### Anwendung aktualisieren

```bash
cd /opt/gin-collection/gin-collection2.0/gin-collection-saas
./deployment/hostinger/deploy.sh update
```

### System aktualisieren

```bash
sudo apt update && sudo apt upgrade -y
sudo reboot
```

---

## Sicherheit

### Bereits konfiguriert

- [x] UFW Firewall (nur 22, 80, 443)
- [x] Fail2Ban (SSH Schutz)
- [x] Non-root Container
- [x] Interne Docker-Netzwerke

### Empfohlene Zusätze

1. **SSH-Key statt Passwort**
   ```bash
   # Auf lokalem Rechner
   ssh-copy-id root@DEINE_VPS_IP

   # Auf VPS: Passwort-Login deaktivieren
   sudo nano /etc/ssh/sshd_config
   # PasswordAuthentication no
   sudo systemctl restart sshd
   ```

2. **Cloudflare (Optional)**
   - DDoS-Schutz
   - CDN für statische Assets
   - Zusätzliche SSL-Schicht

3. **Regelmäßige Updates**
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

---

## Support

Bei Problemen:
1. Logs prüfen: `./deploy.sh logs`
2. GitHub Issues: https://github.com/stefanbuell-wq/gin-collection/issues
