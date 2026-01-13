# 🍸 Installation auf atlas-bergedorf.de/GinVault

Anleitung für die **Erstinstallation** der Gin Collection PWA auf `atlas-bergedorf.de/GinVault`.

> [!IMPORTANT]
> Diese Anleitung ist für eine **Neuinstallation**. Es wird **keine Migration** durchgeführt, sondern eine komplett neue Datenbank erstellt.

---

## 📋 Voraussetzungen

- ✅ FTP-Zugang zu atlas-bergedorf.de
- ✅ SSH-Zugang (empfohlen, aber optional)
- ✅ PHP 7.4+ mit SQLite-Unterstützung
- ✅ HTTPS (für PWA-Funktionen)

---

## 🚀 Installation (4 Schritte)

### Schritt 1: Dateien hochladen

#### Via FTP (z.B. FileZilla):

1. Verbinde dich mit `ftp.atlas-bergedorf.de`
2. Navigiere zum Verzeichnis für die Domain
3. Erstelle den Ordner `GinVault` (falls noch nicht vorhanden)
4. Lade **alle Dateien** aus dem Projekt in `/GinVault/` hoch:

```
/GinVault/
├── index.html
├── login.html
├── manifest.json
├── service-worker.js
├── .htaccess
├── api/
│   ├── index.php
│   ├── Database.php
│   └── Auth.php
├── assets/
│   ├── css/
│   ├── js/
│   └── images/
├── db/
│   ├── schema.sql
│   ├── install.php
│   └── migrate.php
└── uploads/
```

> [!CAUTION]
> Stelle sicher, dass die `.htaccess` Datei hochgeladen wird! Manche FTP-Clients blenden versteckte Dateien aus.

---

### Schritt 2: Verzeichnis-Berechtigungen setzen

#### Via SSH (empfohlen):

```bash
cd /pfad/zu/GinVault
chmod 755 db
chmod 755 uploads
```

#### Via FTP:

1. Rechtsklick auf Ordner `db` → "Dateiberechtigungen" → `755`
2. Rechtsklick auf Ordner `uploads` → "Dateiberechtigungen" → `755`

---

### Schritt 3: Datenbank erstellen

#### Via SSH:

```bash
cd /pfad/zu/GinVault
php db/install.php
```

**Erwartete Ausgabe:**

```
🍸 Gin Collection - Fresh Installation
========================================

Checking prerequisites...

PHP Version: 8.1.0 ✓
SQLite3 Extension: ✓
PDO SQLite: ✓
Database directory writable: ✓
Uploads directory writable: ✓

✅ All prerequisites met!

Creating new database...
Loading schema from schema.sql...
Creating tables and indexes...
Creating default admin user...
✓ Admin user created (ID: 1)

✅ Installation completed successfully!

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Default admin credentials:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Username: admin
  Email:    admin@gin-collection.local
  Password: Admin123!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚠️  IMPORTANT: Change the admin password after first login!
```

#### Kein SSH-Zugang?

Falls du keinen SSH-Zugang hast, erstelle eine temporäre Datei `install_web.php` im Hauptverzeichnis:

```php
<?php
// Temporäre Web-Installation
// WICHTIG: Diese Datei nach Installation LÖSCHEN!

require_once __DIR__ . '/db/install.php';

// Setze Content-Type auf Text
header('Content-Type: text/plain; charset=utf-8');

// Führe Installation aus
$install = new FreshInstall();
$install->checkPrerequisites();
$install->run();
```

Dann:
1. Rufe auf: `https://atlas-bergedorf.de/GinVault/install_web.php`
2. Notiere die Admin-Zugangsdaten
3. **LÖSCHE** `install_web.php` sofort nach der Installation!

---

### Schritt 4: Testen und Admin-Passwort ändern

1. Öffne: `https://atlas-bergedorf.de/GinVault/login.html`
2. Login mit:
   - **Username:** `admin`
   - **Password:** `Admin123!`
3. ✅ Du solltest zur Hauptseite weitergeleitet werden
4. **WICHTIG:** Ändere sofort das Admin-Passwort!

---

## 🔒 Sicherheit

### Admin-Passwort ändern

> [!CAUTION]
> Das Default-Passwort `Admin123!` ist **NICHT sicher** für den Produktivbetrieb!

**So änderst du das Passwort:**

1. Einloggen als Admin
2. Öffne Browser DevTools (F12) → Console
3. Führe aus:
```javascript
fetch('api/?action=change-password', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({
        old_password: 'Admin123!',
        new_password: 'DeinSicheresPasswort123!@#'
    })
}).then(r => r.json()).then(console.log);
```

> [!NOTE]
> Eine Passwort-Änderungs-Funktion im UI kann später hinzugefügt werden.

### HTTPS prüfen

Stelle sicher, dass HTTPS aktiv ist:
```bash
curl -I https://atlas-bergedorf.de/GinVault/
```

Falls nicht, aktiviere SSL im Hosting-Panel.

---

## 👥 Weitere Benutzer erstellen

1. Öffne `https://atlas-bergedorf.de/GinVault/login.html`
2. Klicke auf Tab **"Registrieren"**
3. Fülle das Formular aus:
   - Username (eindeutig)
   - E-Mail (eindeutig)
   - Vollständiger Name (optional)
   - Passwort (min. 8 Zeichen)
4. Klicke **"Registrieren"**

Jeder Benutzer hat seine **eigene private Gin-Sammlung**!

---

## 🧪 Funktionstest

### Test 1: Login

```
✓ Öffne login.html
✓ Login: admin / Admin123!
✓ Weiterleitung zu index.html
✓ Benutzername "admin" im Header sichtbar
```

### Test 2: Gin hinzufügen

```
✓ Klicke auf "Hinzufügen" Button
✓ Fülle Formular aus (Name ist Pflichtfeld)
✓ Speichern
✓ Gin erscheint in der Liste
```

### Test 3: Datenisolation

```
✓ Als Admin einloggen
✓ Gin "Hendrick's" hinzufügen
✓ Ausloggen (🚪 Button)
✓ Neuen User registrieren
✓ Als neuer User einloggen
✓ Gin von Admin NICHT sichtbar
```

---

## 🔧 Troubleshooting

### Problem: "Database error" beim ersten Aufruf

**Lösung:**
```bash
# Prüfe Berechtigungen
ls -la db/
chmod 755 db/
```

### Problem: "Permission denied" beim Schreiben

**Lösung:**
```bash
# Setze Berechtigungen für Webserver-User
chown -R www-data:www-data db/ uploads/
chmod 755 db/ uploads/
```

### Problem: Login funktioniert nicht

**Prüfe:**
1. PHP Sessions aktiviert?
2. Cookies im Browser aktiviert?
3. HTTPS aktiv bei `session.cookie_secure = 1`?

**Debug:**
```bash
# Prüfe PHP Error Log
tail -f /var/log/php-errors.log
```

### Problem: "401 Unauthorized" bei allen Requests

**Lösung:**
- Browser-Cookies aktivieren
- Cookie-Blocker deaktivieren
- Session-Verzeichnis prüfen:
```bash
php -i | grep session.save_path
```

### Problem: Datenbank wurde bereits erstellt

Falls du versehentlich die Installation mehrmals ausführst:

```bash
# Datenbank löschen und neu erstellen
rm db/gin_collection.db
php db/install.php
```

---

## 📁 Wichtige Dateien

### Datenbank
- `db/gin_collection.db` - SQLite Datenbank (wird bei Installation erstellt)

### Konfiguration
- `.htaccess` - Apache Konfiguration (URL Rewriting)
- `manifest.json` - PWA Manifest

### API
- `api/index.php` - Haupt-API Endpoint
- `api/Auth.php` - Authentifizierungs-Klasse
- `api/Database.php` - Datenbank-Wrapper

---

## 🔄 Backup

### Regelmäßige Backups erstellen

**Datenbank sichern:**
```bash
# Via SSH
cp db/gin_collection.db db/backup_$(date +%Y%m%d).db

# Via FTP
# Lade db/gin_collection.db herunter
```

**Uploads sichern:**
```bash
# Via SSH
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz uploads/

# Via FTP
# Lade kompletten uploads/ Ordner herunter
```

### Automatisches Backup (Cronjob)

```bash
# Täglich um 3 Uhr morgens
0 3 * * * cd /pfad/zu/GinVault && cp db/gin_collection.db db/backup_$(date +\%Y\%m\%d).db
```

---

## 📱 PWA Installation

Nach erfolgreicher Installation können Benutzer die App auf ihrem Smartphone installieren:

1. Öffne `https://atlas-bergedorf.de/GinVault/` im Browser
2. Browser zeigt "Zum Startbildschirm hinzufügen" an
3. App verhält sich wie eine native App

**Voraussetzungen:**
- ✅ HTTPS aktiv
- ✅ `manifest.json` vorhanden
- ✅ Service Worker registriert
- ✅ Icons vorhanden (192x192 und 512x512 px)

---

## ✅ Installations-Checkliste

- [ ] Dateien hochgeladen nach `/GinVault/`
- [ ] `.htaccess` Datei vorhanden
- [ ] Berechtigungen gesetzt (`db/` und `uploads/` → 755)
- [ ] Installation ausgeführt (`php db/install.php`)
- [ ] Admin-Login getestet
- [ ] Admin-Passwort geändert
- [ ] HTTPS aktiv
- [ ] Backup-Strategie definiert
- [ ] Weitere Benutzer können sich registrieren

---

## 🎉 Fertig!

Deine Gin Collection PWA ist jetzt unter **https://atlas-bergedorf.de/GinVault/** verfügbar!

### Nächste Schritte:

1. ✅ Lade Freunde ein, sich zu registrieren
2. ✅ Beginne mit dem Erfassen deiner Gin-Sammlung
3. ✅ Nutze den Barcode-Scanner für schnelles Hinzufügen
4. ✅ Erstelle Tasting Notes und Bewertungen

**Viel Spaß mit deiner digitalen Gin-Sammlung! 🍸**

---

## 📞 Support

Bei Problemen:
1. Browser DevTools → Console auf Fehler prüfen
2. PHP Error Log prüfen
3. Datenbank-Integrität prüfen: `sqlite3 db/gin_collection.db ".schema"`

Weitere Dokumentation:
- [README.md](README.md) - Allgemeine Übersicht
- [FEATURES.md](FEATURES.md) - Feature-Liste
- [SETUP_MULTI_TENANCY.md](SETUP_MULTI_TENANCY.md) - Multi-Tenancy Details
