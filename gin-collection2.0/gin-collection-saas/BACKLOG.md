# GinVault - Backlog & Open Points

> Letzte Aktualisierung: 2026-01-19
> Quellen: Security Audit, Phase 10 Summary, Super-Admin Plan, SaaS Migration Plan

---

## In Arbeit

_Aktuell keine offenen Aufgaben_

---

## 🚀 Launch-Vorbereitung

> Aus Phase 10 Summary - Production Deployment Checklist

### Sofort (diese Woche)
- [ ] PayPal Sandbox einrichten und testen
- [x] S3 Bucket für Production konfigurieren - ✅ Cloudflare R2 (10GB Free Tier)
- [ ] Backup/Restore-Skripte testen (`scripts/backup.sh`, `scripts/restore.sh`)

### Kurzfristig (nächste 2-4 Wochen)
- [ ] Production Server Setup
- [ ] Domain und SSL-Zertifikate konfigurieren
- [ ] Deploy to Staging Environment
- [ ] User Acceptance Testing (UAT)
- [ ] Performance Testing
- [ ] Security Audit durchführen
- [ ] Load Testing (Ziel: 1000+ req/s)

### Launch
- [ ] Production Deployment
- [ ] Monitoring-Dashboards verifizieren (Prometheus/Grafana)
- [ ] Soft Launch (Beta-User einladen)
- [ ] Marketing-Ankündigung
- [ ] Full Public Launch 🚀

### Post-Launch
- [ ] Metriken und Alerts monitoren
- [ ] User-Feedback sammeln und auswerten
- [ ] Performance-Optimierungen basierend auf echten Daten
- [ ] Regelmäßige Security-Updates

---

## 🛡️ Super-Admin Platform

> Aus Super-Admin Plan (smooth-singing-peacock.md)

### Admin Dashboard - Statistiken
- [ ] Tenant-Übersicht: Anzahl aktiv/suspendiert/gekündigt
- [ ] User-Gesamtzahl über alle Tenants
- [ ] Gin-Gesamtzahl über alle Tenants
- [ ] Revenue-Anzeige pro Tier (Free/Basic/Pro/Enterprise)
- [ ] Neue Registrierungen Chart (letzte 7/30 Tage)
- [ ] Storage-Nutzung Übersicht

### Tenant-Management (Admin)
- [ ] Tenant-Liste mit Filter und Suche
- [ ] Tenant-Details: User, Gins, Subscription-Status
- [ ] Tenant suspendieren / reaktivieren
- [ ] Subscription-Tier manuell ändern (Override)
- [ ] Tenant soft-delete mit Bestätigung

### User-Management (Admin)
- [ ] Alle User aller Tenants anzeigen
- [ ] User-Suche nach Email und Name
- [ ] User-Detailansicht (Tenant, Rolle, Aktivität)
- [ ] User deaktivieren/aktivieren

### System Health Dashboard
- [ ] API-Status Live-Anzeige
- [ ] Database Connection Status
- [ ] Redis Connection Status
- [ ] S3/Storage Status
- [ ] Background Jobs Status (falls vorhanden)

---

## 🏢 Enterprise Features

> Aus SaaS Migration Plan (vivid-cooking-phoenix.md)

### Multi-Tenancy Enterprise
- [ ] Separate Datenbank pro Enterprise-Tenant (DB Provisioning)
- [ ] Multi-User Support: Owner, Admin, Member, Viewer Rollen
- [ ] Custom Branding: Logo, Farben, Domain
- [ ] SLA Monitoring & Health Checks pro Tenant
- [ ] Audit Logging für Enterprise-Tenants

### Enterprise API
- [ ] API-Key Authentifizierung für Enterprise
- [ ] Webhook-System für Enterprise (Events: gin.created, gin.updated, etc.)
- [ ] Rate Limits: 10000 req/h für Enterprise

### GDPR & Compliance
- [ ] Data Export API (alle Nutzerdaten als JSON)
- [ ] Account Deletion API (vollständige Löschung)
- [ ] Data Retention Policies konfigurierbar

---

## 🧪 Testing & QA

> Aus Phase 10 Summary und SaaS Migration Plan

### Load Testing
- [ ] 1000 concurrent requests/sec Test
- [ ] 10,000 Tenants Simulation
- [ ] Database Connection Pool unter Last testen
- [ ] Redis Rate Limiting Performance verifizieren

### Integration Tests
- [ ] Tenant-Isolation Tests (Cross-Tenant-Leak Prevention)
- [ ] Tier-Enforcement Tests (Feature Gates)
- [ ] PayPal Subscription Flow E2E Test

### Security Tests
- [ ] Penetration Testing durchführen
- [ ] SQL Injection Versuche testen
- [ ] XSS Prevention verifizieren
- [ ] CSRF Token Validation testen

---

## 🔴 Sicherheit - KRITISCH

> Security Audit durchgeführt am 2026-01-18

### Sofort-Maßnahmen (24-48 Stunden)

#### 1. JWT Secret austauschen
- [x] Echtes 256-bit Secret generieren - ✅ Erledigt am 2026-01-18
- [ ] In sicherem Secret Manager speichern (nicht in .env) - Optional für später
- [ ] Validierung im Code: Reject schwache Secrets in Production - Optional

#### 2. Secrets aus Git entfernen
- [x] `.env` aus Git-History entfernen - ✅ War nie committed
- [x] `.env` zu `.gitignore` hinzufügen - ✅ Bereits vorhanden
- [ ] Pre-commit Hook für Secret-Scanning einrichten (git-secrets)
- [ ] Alle Passwörter/API-Keys rotieren (empfohlen bei Production)

#### 3. CSRF-Schutz implementieren - ✅ Erledigt am 2026-01-18
- [x] CSRF-Token Middleware für POST/PUT/DELETE Requests
- [x] Token-Generierung mit `crypto/rand`
- [x] Double-Submit Cookie Pattern + Redis Server-Side Storage
- [x] Secure Cookie in Production
- [x] Frontend: CSRF Token bei App-Start und Login laden
- [x] Frontend: Token in X-CSRF-Token Header bei POST/PUT/DELETE/PATCH senden
- [x] Automatischer Token-Refresh bei CSRF-Fehler

#### 4. Rate Limiting aktivieren - ✅ Erledigt am 2026-01-19
- [x] Login: 5 Versuche pro 15 Min pro IP
- [x] Registrierung: 3 pro Stunde pro IP
- [x] Password Reset: 3 pro Stunde pro E-Mail
- [x] API-Endpoints: 100 Requests pro Stunde pro Tenant (war bereits implementiert)
- [x] Fallback In-Memory Rate Limiting wenn Redis unavailable

#### 5. Tokens aus localStorage entfernen - ✅ Erledigt am 2026-01-19
- [x] JWT in HttpOnly Cookie speichern
- [x] `Secure` Flag setzen (nur HTTPS)
- [x] `SameSite=Strict` setzen
- [x] Refresh Token ebenfalls in HttpOnly Cookie

---

## 🟠 Sicherheit - HOCH

### 6. HTTPS erzwingen
- [ ] nginx: HTTP → HTTPS Redirect (301)
- [ ] HSTS Header: `Strict-Transport-Security: max-age=31536000; includeSubDomains`
- [ ] Minimum TLS 1.3 konfigurieren

### 7. Security Headers hinzufügen
```nginx
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'wasm-unsafe-eval';" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "DENY" always;
add_header X-Permitted-Cross-Domain-Policies "none" always;
add_header Permissions-Policy "geolocation=(), microphone=(), camera=(self)" always;
```

### 8. Token-Blacklist für Logout - ✅ Erledigt am 2026-01-19
- [x] Redis-basierte Token-Blacklist implementieren
- [x] Tokens invalidieren bei: Logout, Password Change, Password Reset
- [x] Blacklist-TTL = JWT-TTL (24h)

### 9. File-Upload Sicherheit - ✅ Erledigt am 2026-01-19
- [x] Magic-Byte Validierung (nicht nur Extension)
- [x] Max Upload Size Middleware (50MB für Bilder)
- [x] Dateinamen sanitizen (UUID + Extension) - war bereits implementiert
- [ ] Virus-Scan Integration (ClamAV) - optional, später

### 10. Passwort-Policy verschärfen
- [ ] Minimum 12 Zeichen
- [ ] Groß-/Kleinbuchstaben + Zahlen + Sonderzeichen erforderlich
- [ ] Common-Password-Check (haveibeenpwned API)
- [ ] Password History (letzte 5 nicht wiederverwendbar)

---

## 🟡 Sicherheit - MITTEL

### Authentifizierung & Autorisierung
- [ ] MFA für Admin-Accounts implementieren
- [ ] IP-Whitelist Option für Admin-Panel
- [ ] API-Key Rotation Mechanismus (jährlich)
- [ ] API-Key Expiration implementieren

### Logging & Monitoring
- [ ] Structured Logging (JSON Format)
- [ ] Sensitive Daten in Logs maskieren
- [ ] Security-Events separat loggen
- [ ] Alerting bei verdächtigen Aktivitäten

### Dependency Security
- [ ] `npm audit` in CI/CD Pipeline
- [ ] `go list -m -json all | nancy` für Go Dependencies
- [ ] Dependabot für GitHub aktivieren
- [ ] Monatliche Security-Updates

### CORS einschränken
- [ ] Localhost-Origins aus Production entfernen
- [ ] Nur spezifische Production-Domains whitelisten
- [ ] Wildcard `*` Support entfernen

### Sonstiges
- [ ] `/.well-known/security.txt` erstellen
- [ ] Vulnerability Disclosure Policy dokumentieren
- [ ] SMTP SKIP_VERIFY in Production verbieten
- [ ] Docker Images pinnen (z.B. `alpine:3.19` statt `alpine:latest`)

---

## Offen - Hohe Priorität

### GIN Tasting Anleitung (Premium Feature)
**Beschreibung:** Hochwertige, professionelle GIN Tasting Anleitung im GinVault-Design (ginvault.cloud Style)

| Tier | Features |
|------|----------|
| **Basic** | PDF-Download der Tasting-Anleitung |
| **Pro** | PDF + Digitale Unterstützung in der GinVault App (interaktive Anleitung, Notizen, Bewertungen während Tasting) |
| **Enterprise** | Alles aus Pro + Komplette digitale Plattform zur Umsetzung in der Gastronomie (Event-Management, Gäste-Einladungen, Live-Voting, Ergebnis-Präsentation, Branding) |

**Akzeptanzkriterien:**
- [ ] PDF-Design im GinVault Premium-Style (Dark Theme, Gold Akzente)
- [ ] Tasting-Ablauf mit Schritten (Optik, Geruch, Geschmack, Abgang)
- [ ] Bewertungsbogen / Scoring-System
- [ ] **Basic:** Download-Button auf Subscription-Seite
- [ ] **Pro:** In-App Tasting-Modus mit Timer, Schritt-für-Schritt Anleitung
- [ ] **Pro:** Tasting-Notizen speichern und mit Gin verknüpfen
- [ ] **Enterprise:** Tasting-Events erstellen und verwalten
- [ ] **Enterprise:** Gäste per Link/QR-Code einladen (ohne Account)
- [ ] **Enterprise:** Live-Dashboard mit Ergebnissen
- [ ] **Enterprise:** White-Label / Custom Branding für Events

---

### Backend
- [ ] Webhook-System für Enterprise implementieren

### Frontend
- [ ] API Key Management UI für Pro/Enterprise User
- [ ] Webhook-Konfiguration UI für Enterprise

### FAQ & Dokumentation

#### FAQ für die App (In-App Help)
- [ ] "Was kann ich mit GinVault machen?" - Feature-Übersicht
- [ ] "Wie füge ich einen Gin hinzu?" - Schritt-für-Schritt
- [ ] "Wie funktioniert der Barcode-Scanner?"
- [ ] "Wie erstelle ich Tasting-Notizen?"
- [ ] "Was bedeuten die verschiedenen Abos?"
- [ ] Erreichbar über Help-Icon in der App

#### FAQ für die Webseite (Landing Page)
- [ ] Feature-Erklärung: Was ist GinVault?
- [ ] Subscription-Vergleichstabelle:
  - Free: X Gins, 1 Foto, Basic Features
  - Basic: X Gins, X Fotos, [Features]
  - Pro: Unbegrenzt, API-Zugang, AI, Export
  - Enterprise: Multi-User, Custom Branding, Webhooks
- [ ] Preisübersicht (monatlich/jährlich)
- [ ] "Für wen ist GinVault?" - Zielgruppen
- [ ] "Kann ich upgraden/downgraden?"
- [ ] "Wie kündige ich mein Abo?"
- [ ] Kontakt & Support-Infos

### Infrastruktur

#### PayPal Integration (Geschätzter Aufwand: ~2 Stunden)
- [ ] PayPal Developer Account erstellen (https://developer.paypal.com)
- [ ] Sandbox App erstellen (Dashboard → Apps & Credentials → Create App)
  - Client ID notieren
  - Client Secret notieren
- [ ] Billing Plans in PayPal anlegen:
  - [ ] Basic Monthly (4,99€/Monat)
  - [ ] Basic Yearly (49,99€/Jahr)
  - [ ] Pro Monthly (9,99€/Monat)
  - [ ] Pro Yearly (99,99€/Jahr)
  - [ ] Enterprise Monthly (29,99€/Monat)
  - [ ] Enterprise Yearly (299,99€/Jahr)
- [ ] Webhook einrichten (URL: /api/v1/webhooks/paypal)
  - Events: BILLING.SUBSCRIPTION.ACTIVATED, CANCELLED, SUSPENDED, PAYMENT.SALE.COMPLETED
  - Webhook ID notieren
- [ ] Environment Variables auf Server setzen:
  - PAYPAL_CLIENT_ID
  - PAYPAL_CLIENT_SECRET
  - PAYPAL_MODE=sandbox (später: live)
  - PAYPAL_WEBHOOK_ID
- [ ] Plan IDs im Code hinterlegen (internal/domain/models/subscription.go)
- [ ] Sandbox-Tests durchführen
- [ ] Live schalten (PAYPAL_MODE=live, neue Live-Credentials)

#### Weitere Infrastruktur
- [x] S3 Storage für Produktion einrichten - ✅ Cloudflare R2 konfiguriert am 2026-01-18
- [ ] S3 Presigned URLs implementieren (1h Expiry für private Bilder) - *erst bei 9GB Storage relevant*
- [x] CDN für Static Assets - ✅ Cloudflare R2 Public URL
- [ ] MySQL Transparent Data Encryption aktivieren
- [x] SMTP für E-Mail-Versand konfiguriert (Hostinger, info@ginvault.cloud)

#### Monitoring & Observability
- [ ] Prometheus Metriken Dashboard erstellen
- [ ] Grafana Visualisierungen konfigurieren
- [ ] Alert-Rules aktivieren:
  - API Down (>1 Minute)
  - High Error Rate (>5%)
  - Database Connection Pool Exhaustion
  - Disk Space Low (<10%)
- [ ] Structured Logging (JSON) in Production aktivieren

---

### Gin-Lexikon nach Ländern
**Beschreibung:** Umfassende Gin-Enzyklopädie, gegliedert nach Herkunftsländern

| Format | Inhalt |
|--------|--------|
| **GinVault-Modul** | In-App Lexikon mit Suchfunktion, Filter nach Land/Region, Verlinkung zur eigenen Sammlung |
| **Buch/PDF** | Premium-Publikation als Nachschlagewerk, evtl. Print-on-Demand oder E-Book |

**Länder-Kapitel:**
- [ ] Deutsche Gins (Schwarzwald, Bayern, Berlin, etc.)
- [ ] Englische Gins (London Dry, Plymouth, etc.)
- [ ] Schottische Gins
- [ ] Spanische Gins
- [ ] Niederländische Gins (Genever-Tradition)
- [ ] Amerikanische Gins
- [ ] Japanische Gins
- [ ] Weitere Länder (Australien, Südafrika, etc.)

**Inhalte pro Land:**
- [ ] Geschichte & Tradition der Gin-Herstellung
- [ ] Typische Botanicals der Region
- [ ] Top-Brennereien mit Portraits
- [ ] Empfohlene Gins (Klassiker + Geheimtipps)
- [ ] Regionale Tonic-Pairings

**Akzeptanzkriterien:**
- [ ] Mindestens 5 Länder zum Launch
- [ ] Pro Land: 10-20 Gin-Einträge mit Details
- [ ] Integration mit Sammlung ("Habe ich" / "Möchte ich")
- [ ] Suchbar und filterbar
- [ ] Optional: Buchversion als Premium-Download (Pro/Enterprise)

---

## Offen - Mittlere Priorität

### Features
- [ ] Barcode-Scanner optimieren (bessere Kamera-Unterstützung)
- [ ] Gin-Import aus CSV/Excel
- [ ] Cocktail-Rezept-Verwaltung erweitern
- [ ] Botanicals-Datenbank mit Vorschlägen
- [ ] Dark/Light Mode Toggle

### Admin Panel
- [ ] Platform Admin Dashboard erweitern (siehe Super-Admin Platform Sektion)
- [ ] Tenant-Statistiken verbessern
- [ ] Audit-Log für Admin-Aktionen
- [ ] Server-Deployment UI (Git Pull, Docker Restart)
- [ ] Backup-Management UI (manueller Backup-Trigger)
- [ ] System-Logs Viewer im Admin Panel

### Performance
- [ ] Redis Caching für häufige Abfragen
- [ ] Bild-Optimierung (WebP, Thumbnails)
- [ ] Lazy Loading für Gin-Listen

---

## Offen - Niedrige Priorität

### Nice-to-have
- [ ] PWA Push-Benachrichtigungen
- [ ] Gin-Sharing (öffentliche Links)
- [ ] Sammlung-Statistiken exportieren (PDF)
- [ ] Multi-Language Support (EN, FR)
- [ ] Gin-Vergleichs-Feature
- [ ] Wunschliste für Gins
- [ ] Tasting Themen Basic nur Pdf, Pro Pdf + Digitale unterstützung, Enterprise wie Pro + Komplette plattform für gastronomie

### Technische Schulden
- [ ] Unit Tests erweitern (Coverage > 80%)
- [ ] E2E Tests mit Playwright
- [ ] API Documentation (Swagger/OpenAPI)
- [ ] Error Tracking (Sentry Integration)

### CI/CD Pipeline
- [ ] GitHub Actions Workflows aktivieren (`.github/workflows/ci.yml`)
- [ ] Automatic Deployment to Staging bei Push auf `develop`
- [ ] Automatic Deployment to Production bei Push auf `main`
- [ ] Security Scanning in Pipeline (Trivy, Gosec)
- [ ] Slack/Discord Notifications bei Deployment

---

## Erledigt

### Projekt-Meilensteine

**Phase 1-10 der SaaS-Migration abgeschlossen** ✅
> Dokumentiert in `docs/PHASE_10_SUMMARY.md`

- [x] Phase 1: Foundation (Go Module, Docker, Config)
- [x] Phase 2: Domain & Repositories
- [x] Phase 3: Auth & Multi-Tenancy
- [x] Phase 4: Core Gin Features (CRUD, Search, Export)
- [x] Phase 5: Subscriptions & Feature-Gating
- [x] Phase 6: Advanced Features (Botanicals, Cocktails, AI)
- [x] Phase 7: Enterprise Features (Multi-User, API Keys)
- [x] Phase 8: Frontend (React SPA, PWA)
- [x] Phase 9: Migration & Testing
- [x] Phase 10: Deployment (Docker, CI/CD, Monitoring)

### 2026-01-19
- [x] **Token-Blacklist für JWT-Invalidierung implementiert** (Backend)
  - Redis-basierte Token-Blacklist für sofortige Token-Invalidierung
  - JTI (JWT ID) zu Token-Claims hinzugefügt für eindeutige Identifikation
  - Token wird bei Logout zur Blacklist hinzugefügt
  - Alle User-Tokens werden bei Password Change/Reset invalidiert
  - Blacklist-Check in Auth-Middleware (RequireAuth, OptionalAuth)
  - Graceful Degradation wenn Redis nicht verfügbar
  - Redis-Keys: `blacklist:{jti}` (TTL: verbleibende Token-Lebenszeit), `user_revoked:{userID}` (TTL: 24h)
  - Dateien: tokenblacklist.go (neu), jwt.go, auth.go, auth_handler.go, service.go, main.go
- [x] **File-Upload Sicherheit implementiert** (Backend)
  - Magic-Byte-Validierung statt Extension-basierter Prüfung
  - Prüft tatsächlichen Dateiinhalt (JPEG, PNG, GIF, WebP)
  - Verhindert Upload von getarnten Dateien (exe→jpg)
  - Upload-Size-Limit Middleware (50MB für Bilder)
  - Schützt vor Memory-Exhaustion-Angriffen
  - Dateien: filevalidation.go (neu), upload.go (neu), service.go, router.go
- [x] **HttpOnly JWT Cookies implementiert** (Backend + Frontend)
  - JWT Access Token in HttpOnly Cookie (statt localStorage)
  - Refresh Token in HttpOnly Cookie (Path: /api/v1/auth)
  - Secure Flag in Production, SameSite=Strict
  - Backend: Cookie-Helper (pkg/utils/cookies.go)
  - Backend: Auth-Handler setzt Cookies bei Login/Register, löscht bei Logout
  - Backend: Auth-Middleware liest Token aus Cookie ODER Header (Backward-Compat für API-Keys)
  - Frontend: localStorage Token-Speicherung entfernt
  - Frontend: Refresh-Flow angepasst (Cookie wird automatisch gesendet)
  - Dateien: cookies.go, config.go, auth_handler.go, auth.go, main.go, client.ts, authStore.ts
- [x] **Rate Limiting erweitert** (Backend)
  - Login: 5 Versuche pro 15 Min pro IP (reduziert von 10)
  - Registration: 3 pro Stunde pro IP (neu)
  - Password Reset: 3 pro Stunde pro Email (neu)
  - Token Validation: 3 pro Stunde pro IP (neu)
  - In-Memory Fallback wenn Redis unavailable (neu)
  - Dateien: rate_limit.go, router.go

### 2026-01-18
- [x] **CSRF-Schutz implementiert** (Backend + Frontend)
  - Backend: csrf.go Middleware mit Double-Submit Cookie + Redis Validation
  - Token-Generierung mit crypto/rand (32 Bytes, 24h Expiry)
  - Secure Cookie Flag in Production
  - Endpoint: GET /api/v1/csrf-token
  - Frontend: Token bei App-Start und Login laden
  - Axios Interceptor sendet X-CSRF-Token Header bei POST/PUT/DELETE/PATCH
  - Automatischer Token-Refresh bei CSRF-Fehler (403 + CSRF_* Code)
  - Dateien: csrf.go, router.go, main.go, client.ts, authStore.ts, App.tsx
- [x] **Barcode-Scanner Button im Dashboard aktiviert**
  - Pulsierender Scanner-Button unten rechts war nur visuell
  - Jetzt mit BarcodeScanner-Komponente verbunden
  - Scannt Barcode → API-Lookup → Navigation zu GinCreate mit vorausgefüllten Daten
  - Mobile Touch-Support verbessert (z-index, pointer-events)
  - Dateien: Dashboard.tsx, Dashboard.css, GinCreate.tsx
- [x] **Mobile Double-Click Bug Fix** (GinCreate.tsx)
  - Problem: Auf Mobile musste man 2x auf "Speichern" klicken
  - Ursache: Mobile Touch-Events werden bei Form-Submit anders behandelt
  - Lösung: Button von `type="submit"` auf `type="button"` geändert mit expliziten `onClick` und `onTouchEnd` Handlern
- [x] **Upgrade Modal als Overlay** (GinCreate.tsx, GinCreate.css)
  - Modal erscheint jetzt als fixed Overlay über der Seite
  - Backdrop mit Blur-Effekt
  - Zentrierte Darstellung auf allen Geräten
- [x] **Debug-Modus für Mobile-Testing**
  - URL-Parameter `?debug=1` aktiviert sichtbares Debug-Panel
  - Zeigt letzte 20 Log-Einträge mit Timestamps
  - Hilfreich für Debugging auf Mobile-Geräten ohne DevTools
- [x] **Gin-Limit Enforcement mit Upgrade-Benachrichtigung**
  - Backend gibt `upgrade_required: true` zurück bei 403
  - Enthält: `limit`, `current_count`, `current_tier`
  - ginStore.ts extrahiert Upgrade-Info aus Fehler-Response
  - Modal zeigt aktuellen Tier und Limit an
- [x] **Password Reset Feature** (Backend + Frontend)
  - Forgot Password Seite (ForgotPassword.tsx)
  - Reset Password Seite (ResetPassword.tsx)
  - Backend-Endpoints: `/forgot-password`, `/reset-password`, `/validate-reset-token`
  - Token-basiertes Reset mit E-Mail-Versand
- [x] **GinVault E-Mail Templates**
  - Branding auf GinVault umgestellt
  - Dark Theme Design

### 2026-01-16
- [x] Verkostungsnotizen/Tasting Sessions pro Gin implementiert
  - Backend: Repository, Service, Handler
  - API-Endpoints: GET/POST/PUT/DELETE /gins/:id/tastings
  - Frontend: TastingSessions-Komponente mit GinVault-Design
  - Mehrere Verkostungen pro Gin mit Datum, Bewertung, Notizen
- [x] AI-Integration mit Ollama (lokal, kostenlos)
- [x] GinVault Dark Theme durchgängig implementiert
- [x] Layout.tsx auf Vault-Theme umgestellt
- [x] Tier-Werte Backend/Frontend synchronisiert
- [x] API-Dokumentation erstellt (docs/API-INTEGRATION.md)
- [x] Pro-Tier: API-Zugang hinzugefügt
- [x] API Key Middleware für Pro-Tier aktiviert
- [x] Rate Limiting implementiert (Redis-basiert, Tier-abhängig)

---

## Neue Anforderung hinzufügen

```markdown
### [Titel der Anforderung]
**Priorität:** Hoch / Mittel / Niedrig
**Beschreibung:**
[Beschreibung der Anforderung]

**Akzeptanzkriterien:**
- [ ] Kriterium 1
- [ ] Kriterium 2
```

---

## Vision / Roadmap - Langfristige Features

### 1. Automatische Etikettenerkennung (Label Recognition)
**Priorität:** Hoch
**Beschreibung:** Ein Foto der Flasche macht GinVault zum intelligenten Erkennungssystem:
- Marke automatisch erkennen
- Botanicals identifizieren
- Alkoholgehalt auslesen
- Herkunft bestimmen
- Preisrange schätzen

**USP:** Massiver UX-Boost und starker Differentiator gegenüber allen existierenden Gin-Apps.

---

### 2. KI-gestützte Aromenanalyse ("Aroma-Coach")
**Priorität:** Hoch
**Beschreibung:** GinVault als intelligenter Geschmacks-Berater:
- Nutzer geben Lieblingsgins ein
- KI erkennt Muster (z.B. "floral + citrus")
- GinVault schlägt neue Gins vor, die exakt ins Profil passen

**Vision:** Wie Spotify-Discover - nur für Gin.

---

### 3. Händler- und Brennerei-Dashboards (B2B-Modul)
**Priorität:** Mittel
**Beschreibung:** Ein B2B-Modul für Hersteller mit Analytics:
- Welche Gins werden am häufigsten gesammelt
- Welche Aromen im Trend sind
- Welche Zielgruppen welche Gins bevorzugen

**Potenzial:** Macht GinVault für die Industrie extrem wertvoll.

---

### 4. Limited Editions exklusiv für GinVault
**Priorität:** Mittel
**Beschreibung:** Kooperationen mit Brennereien:
- "GinVault Edition No. 1"
- Exklusive Batch-Releases
- Nur für Premium-Mitglieder

**Potenzial:** Schafft Begehrlichkeit und wiederkehrende Umsätze.

---

### 5. Gamification & Achievements
**Priorität:** Mittel
**Beschreibung:** Sammler lieben Status. Beispiele:
- "10 Tastings abgeschlossen"
- "Botanical-Master: 50 Aromen erkannt"
- "Rare Bottle Collector"
- Badges und Level-System

**Potenzial:** Erhöht Retention und Community-Dynamik massiv.

---

### 6. Social Features mit echtem Mehrwert
**Priorität:** Mittel
**Beschreibung:** Nicht nur Likes, sondern:
- Tasting-Vergleiche mit Freunden
- Gemeinsame Tasting-Sessions
- "Flavor Match Score" zwischen Nutzern
- Challenges ("Taste 5 Mediterranean Gins this month")

**Vision:** Das macht GinVault lebendig und community-driven.

---

### 7. Integration mit Bars & Events
**Priorität:** Mittel
**Beschreibung:** GinVault als digitaler Begleiter für reale Erlebnisse:
- Bars integrieren ihre Gin-Karte
- Nutzer scannen Gin im Restaurant
- Tasting wird automatisch gespeichert
- Events über GinVault buchen

**Vision:** Verbindet Online und Offline nahtlos.

---

### 8. Erweiterte Abo-Modelle mit echtem Mehrwert
**Priorität:** Hoch
**Beschreibung:** Premium-Features die sich "lohnen":
- Unbegrenzte Sammlung
- Exklusive Gins
- Deep-Analytics
- KI-Empfehlungen
- Early Access zu Limited Editions
- Rabattcodes bei Partnern

---

### 9. GinVault als Geschenkprodukt
**Priorität:** Niedrig
**Beschreibung:** GinVault in Geschenkboxen integrieren:
- 3 Gins + Premium-Abo für 3 Monate
- QR-Code führt direkt zum Tasting-Erlebnis

**Potenzial:** Perfektes Weihnachts- und Geburtstagsprodukt.

---

### 10. API für Shops & Hersteller
**Priorität:** Mittel
**Beschreibung:** Shops können:
- GinVault-Bewertungen anzeigen
- Aromenprofile integrieren
- "Passt zu deinem Geschmack"-Empfehlungen nutzen

**Vision:** GinVault wird zur Infrastruktur des Gin-Markts.

---

### 11. Community-Ranking & Awards
**Priorität:** Niedrig
**Beschreibung:** Jährlicher "GinVault Community Award":
- Beste Gins
- Beste Newcomer
- Beste Brennerei

**Potenzial:** Schafft Reichweite und Presse-Aufmerksamkeit.

---

### 12. GinVault für Firmen (Corporate Tasting)
**Priorität:** Mittel
**Beschreibung:** Firmen lieben Tasting-Events:
- Digitale Tasting-Boxen
- Moderierte Sessions
- Firmenprofile
- Mitarbeiter-Challenges

**Potenzial:** Extrem lukrativer B2B-Markt.

---

### 13. Erweiterte Cocktail-Features
**Priorität:** Mittel
**Beschreibung:** Viele Gin-Fans sind auch Cocktail-Fans:
- Cocktail-Rezepte basierend auf eigener Sammlung
- "Was kann ich mit meinen Gins mixen?"
- KI-Mixing-Assistent

**Potenzial:** Erweitert die Zielgruppe enorm.

---

### 14. Sammler-Wertentwicklung & Rare-Bottle-Tracking
**Priorität:** Niedrig
**Beschreibung:** Für Premium-Sammler:
- Preisentwicklung tracken
- Seltenheitsindex
- Marktwert der eigenen Sammlung

**Vision:** GinVault als "Gin-Portfolio-Tracker".

---

### 15. Tasting-Box-Ökosystem
**Priorität:** Mittel
**Beschreibung:** Boxen nicht nur verkaufen, sondern:
- Boxen mit Partnern co-branden
- Boxen als Abo anbieten
- Boxen als Onboarding-Tool für neue Nutzer

**Potenzial:** Verstärkt das Flywheel und schafft wiederkehrende Umsätze.

---

## Notizen

_Platz für allgemeine Notizen und Ideen_

