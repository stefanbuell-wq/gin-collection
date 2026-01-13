# 🔐 Multi-Tenancy Setup Guide

Schnellanleitung zur Einrichtung der Mandantenfähigkeit für die Gin Collection PWA.

---

## 📋 Voraussetzungen

- ✅ PHP 7.4+ mit SQLite-Unterstützung
- ✅ Webserver (Apache/Nginx) mit mod_rewrite
- ✅ HTTPS empfohlen (für sichere Sessions)

---

## 🚀 Installation (3 Schritte)

### Schritt 1: Datenbank Migration

Führe das Migrations-Script aus, um die Benutzer-Tabelle zu erstellen:

```bash
cd /pfad/zur/gin-collection
php db/migrate.php
```

**Erwartete Ausgabe:**
```
Starting database migration...
Creating users table...
Creating default admin user...
Adding user_id column to gins table...
✓ Migration completed successfully!

Default admin credentials:
  Username: admin
  Email: admin@gin-collection.local
  Password: Admin123!

⚠️  IMPORTANT: Please change the admin password after first login!
```

### Schritt 2: Verzeichnis-Berechtigungen

Stelle sicher, dass der Webserver Schreibrechte hat:

```bash
chmod 755 db/
chmod 644 db/gin_collection.db
chmod 755 uploads/
```

### Schritt 3: Testen

1. Öffne `https://deine-domain.de/gin-collection/login.html`
2. Login mit Admin-Credentials:
   - **Username:** `admin`
   - **Password:** `Admin123!`
3. ✅ Du solltest zur Hauptseite weitergeleitet werden

---

## ⚙️ Wichtige Konfiguration

### 1. Admin-Passwort ändern

> [!CAUTION]
> **SOFORT nach dem ersten Login das Passwort ändern!**

Das Default-Passwort ist **nicht sicher** für den Produktivbetrieb.

### 2. HTTPS aktivieren

Für sichere Sessions ist HTTPS **dringend empfohlen**:

- Bei df.eu: Let's Encrypt SSL-Zertifikat aktivieren
- Oder: Eigenes SSL-Zertifikat installieren

### 3. Session-Konfiguration (Optional)

In `api/index.php` sind bereits sichere Session-Einstellungen konfiguriert:

```php
ini_set('session.cookie_httponly', 1);
ini_set('session.cookie_samesite', 'Strict');
ini_set('session.use_strict_mode', 1);
```

Für HTTPS zusätzlich aktivieren:
```php
ini_set('session.cookie_secure', 1);  // Nur über HTTPS
```

---

## 👥 Benutzer-Verwaltung

### Neuen Benutzer erstellen

1. Öffne `login.html`
2. Klicke auf Tab "Registrieren"
3. Fülle Formular aus:
   - Username (erforderlich, eindeutig)
   - E-Mail (erforderlich, eindeutig)
   - Vollständiger Name (optional)
   - Passwort (min. 8 Zeichen)
4. Klicke "Registrieren"

### Bestehende Gins

Alle Gins, die **vor der Migration** existierten, sind automatisch dem **Admin-Benutzer** zugeordnet.

---

## 🔒 Sicherheits-Features

### Implementiert

- ✅ **Bcrypt Password Hashing** (cost factor 12)
- ✅ **Session-basierte Authentifizierung**
- ✅ **SQL Injection Protection** (Prepared Statements)
- ✅ **Datenisolation** (jeder User sieht nur seine Gins)
- ✅ **Ownership Verification** (vor Update/Delete)
- ✅ **Secure Session Cookies** (HTTPOnly, SameSite)

### Empfohlene Zusatz-Maßnahmen

- 🔐 **HTTPS** für Production
- 🚫 **Rate Limiting** für Login-Versuche
- 📧 **Email-Verifikation** bei Registrierung
- 🔑 **Password Reset** Funktion
- 🔐 **2FA** (optional)

---

## 🧪 Funktionstest

### Test 1: Login

```
1. Öffne login.html
2. Login: admin / Admin123!
3. ✅ Weiterleitung zu index.html
4. ✅ Benutzername im Header sichtbar
```

### Test 2: Datenisolation

```
1. Als User A einloggen
2. Gin "Hendrick's" hinzufügen
3. Ausloggen (🚪 Button)
4. Als User B einloggen
5. ✅ Gin von User A NICHT sichtbar
6. Gin "Bombay Sapphire" hinzufügen
7. ✅ Nur eigener Gin sichtbar
```

### Test 3: API-Sicherheit

```
1. Ausloggen
2. Browser DevTools → Console
3. fetch('api/?action=list').then(r => r.json()).then(console.log)
4. ✅ Fehler: 401 Unauthorized
```

---

## 📁 Geänderte Dateien

### Backend
- ✅ `db/schema.sql` - Erweitert um `users` Tabelle
- ✅ `db/migrate.php` - **NEU** - Migrations-Script
- ✅ `api/Auth.php` - **NEU** - Authentifizierungs-Klasse
- ✅ `api/index.php` - Erweitert um Auth-Endpunkte

### Frontend
- ✅ `login.html` - **NEU** - Login/Registrierungs-Seite
- ✅ `assets/js/auth.js` - **NEU** - Auth-Modul
- ✅ `index.html` - Erweitert um User-Info & Logout
- ✅ `assets/css/style.css` - Erweitert um User-Info Styles

---

## 🔧 Troubleshooting

### Problem: "Migration already completed"

**Lösung:** Die Datenbank wurde bereits migriert. Kein Handlungsbedarf.

### Problem: "Permission denied" beim Schreiben

**Lösung:** 
```bash
chmod 755 db/
chmod 644 db/gin_collection.db
```

### Problem: Login funktioniert nicht

**Prüfe:**
1. PHP Sessions aktiviert? (`session_start()` funktioniert?)
2. Cookies aktiviert im Browser?
3. HTTPS bei `session.cookie_secure = 1`?

### Problem: "401 Unauthorized" bei allen Requests

**Lösung:** Session-Cookie wird nicht gesetzt. Prüfe:
- Browser-Cookies aktiviert
- Kein Cookie-Blocker aktiv
- Session-Verzeichnis beschreibbar

---

## 🚀 Nächste Schritte (Optional)

### Phase 2: OAuth2 SSO Integration

Für Google/Facebook Login:

1. **Google OAuth2**
   - Google Cloud Console → Credentials
   - OAuth 2.0 Client ID erstellen
   - Redirect URI: `https://deine-domain.de/gin-collection/api/?action=oauth-callback`

2. **Facebook Login**
   - Facebook Developers → App erstellen
   - Facebook Login aktivieren
   - App ID und Secret notieren

3. **Implementation**
   - OAuth2 Library installieren
   - Callback-Handler in `api/index.php`
   - SSO-Buttons in `login.html`

---

## 📞 Support

Bei Problemen:
1. Browser DevTools → Console auf Fehler prüfen
2. PHP Error Log prüfen
3. Datenbank-Integrität prüfen: `sqlite3 db/gin_collection.db ".schema"`

---

## ✅ Checkliste

- [ ] Migration ausgeführt (`php db/migrate.php`)
- [ ] Admin-Login getestet
- [ ] Admin-Passwort geändert
- [ ] HTTPS aktiviert (empfohlen)
- [ ] Berechtigungen gesetzt
- [ ] Datenisolation getestet
- [ ] Backup der Datenbank erstellt

**Status:** ✅ Multi-Tenancy ist einsatzbereit!
