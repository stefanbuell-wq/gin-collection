# Installation bei df.eu

## Schritt 1: FTP-Verbindung herstellen

1. Verwende einen FTP-Client (z.B. FileZilla)
2. Verbinde dich mit deinen df.eu FTP-Zugangsdaten:
   - Host: `ftp.deine-domain.de` (oder die von df.eu bereitgestellte Adresse)
   - Benutzername: Dein FTP-Username
   - Passwort: Dein FTP-Passwort
   - Port: 21

## Schritt 2: Dateien hochladen

Lade alle Dateien und Ordner in dein Webroot-Verzeichnis hoch:
- Bei df.eu ist dies meist `/html/` oder `/public_html/`
- Du kannst auch einen Unterordner wie `/html/gin-collection/` verwenden

### Wichtige Dateien:
```
/html/
├── index.html
├── manifest.json
├── service-worker.js
├── .htaccess
├── api/
│   ├── index.php
│   └── Database.php
├── assets/
│   ├── css/
│   ├── js/
│   └── images/
├── db/
│   └── schema.sql
└── uploads/
```

## Schritt 3: Verzeichnis-Berechtigungen

**Wichtig:** Diese Verzeichnisse müssen beschreibbar sein:

### Via FTP:
1. Rechtsklick auf Ordner `db`
2. "Dateiberechtigungen" → Setze auf `755` oder `775`
3. Wiederhole für Ordner `uploads`

### Via SSH (falls verfügbar):
```bash
chmod 755 db
chmod 755 uploads
```

## Schritt 4: PHP-Einstellungen prüfen (optional)

df.eu bietet meist ein PHP-Konfigurationspanel. Prüfe:
- PHP Version: mindestens 7.4 (empfohlen: 8.0+)
- SQLite-Unterstützung: sollte standardmäßig aktiviert sein
- `upload_max_filesize`: mindestens 10M (für Fotos)
- `post_max_size`: mindestens 10M

## Schritt 5: SSL/HTTPS aktivieren

Für die volle PWA-Funktionalität benötigst du HTTPS:

1. Login bei df.eu
2. Navigiere zu "SSL/TLS"
3. Aktiviere kostenloses Let's Encrypt SSL-Zertifikat
4. Warte 5-10 Minuten bis das Zertifikat aktiv ist

## Schritt 6: App testen

1. Öffne deinen Browser
2. Gehe zu: `https://deine-domain.de/gin-collection/`
3. Die App sollte sofort funktionieren
4. Die Datenbank wird automatisch beim ersten Zugriff erstellt

### Erste Schritte nach Installation:
1. Klicke auf "Hinzufügen"
2. Trage deinen ersten Gin ein
3. Teste den Barcode-Scanner (erfordert HTTPS und Kamera-Berechtigung)

## Schritt 7: Icons erstellen (optional)

Für ein professionelles App-Erlebnis erstelle zwei Icons:

### Icon 192x192px:
```
/assets/images/icon-192.png
```

### Icon 512x512px:
```
/assets/images/icon-512.png
```

**Tipp:** Du kannst online Icon-Generatoren verwenden oder ein einfaches Gin-Emoji als Platzhalter nutzen.

## Troubleshooting bei df.eu

### Problem: "Database error"
**Lösung:** 
- Prüfe Schreibrechte auf `/db/` Verzeichnis (755 oder 775)
- Stelle sicher, dass SQLite in PHP aktiviert ist

### Problem: "Permission denied" beim Upload
**Lösung:**
- Prüfe Schreibrechte auf `/uploads/` Verzeichnis (755 oder 775)

### Problem: Scanner funktioniert nicht
**Lösung:**
- Stelle sicher, dass HTTPS aktiviert ist
- Erlaube Kamera-Zugriff im Browser
- Teste mit gutem Licht

### Problem: .htaccess wird nicht geladen
**Lösung:**
- Stelle sicher, dass die Datei wirklich `.htaccess` heißt (mit Punkt am Anfang!)
- Bei manchen FTP-Clients sind "versteckte Dateien" standardmäßig ausgeblendet
- Aktiviere "Versteckte Dateien anzeigen" in deinem FTP-Client

### Problem: "Internal Server Error"
**Lösung:**
- Prüfe ob mod_rewrite aktiviert ist (sollte bei df.eu standard sein)
- Prüfe PHP Error Logs im df.eu Control Panel
- Kommentiere temporär die RewriteRules in `.htaccess` aus

## df.eu spezifische Features nutzen

### Cronjobs (optional)
Du kannst Cronjobs für automatische Backups einrichten:
```bash
# Täglich um 3 Uhr morgens
0 3 * * * cd /pfad/zu/html/gin-collection && php backup.php
```

### Datenbank-Backup
**Empfehlung:** Sichere regelmäßig die Datei `/db/gin_collection.db`
- Download via FTP
- Oder erstelle ein Backup-Script

### PHP Error Logs
Bei Problemen schaue in:
- df.eu Control Panel → "Logs" → "PHP Error Log"

## Performance-Optimierung

### Browser-Caching
Die `.htaccess` Datei aktiviert bereits Caching für statische Assets.

### Bild-Optimierung
Komprimiere hochgeladene Fotos:
- Max. 1920px Breite empfohlen
- JPEG mit 85% Qualität

## Sicherheits-Tipps

1. **Regelmäßige Backups**
   - Sichere `/db/gin_collection.db` wöchentlich
   - Sichere `/uploads/` Ordner monatlich

2. **Passwort-Schutz** (optional)
   Für zusätzliche Sicherheit kannst du einen .htpasswd Schutz einrichten:
   ```apache
   # In .htaccess hinzufügen:
   AuthType Basic
   AuthName "Gin Collection"
   AuthUserFile /pfad/zu/.htpasswd
   Require valid-user
   ```

3. **Upload-Limit**
   Die App hat bereits Sicherheits-Checks, aber zusätzlich in `.htaccess`:
   ```apache
   php_value upload_max_filesize 10M
   php_value post_max_size 10M
   ```

## Support

Bei technischen Problemen mit df.eu:
- df.eu Support kontaktieren
- Community-Forum nutzen

Bei Fragen zur App:
- Siehe README.md
- Prüfe Browser-Konsole (F12)

---

**Viel Erfolg mit deiner Gin Collection App! 🍸**
