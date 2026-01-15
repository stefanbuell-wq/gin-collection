# Changelog - Gin Collection SaaS

Alle dokumentierten Änderungen am Projekt.

---

## [2026-01-15] - Frontend Integration & Enterprise Features

### Zusammenfassung
Vollständige Integration des Frontends mit den Backend-APIs, Implementierung des Email-Systems für Benutzereinladungen, Erstellung eines Database Migration CLI Tools, und Befüllung der Datenbank mit Testdaten für alle Tier-Levels.

---

### Phase 2: User Management Frontend Integration

#### Geänderte Dateien

**`frontend/src/pages/Users.tsx`**
- Komplette Überarbeitung mit Anbindung an echte APIs
- Loading States für alle Operationen
- Error Handling mit Toast-Benachrichtigungen
- API-Calls implementiert:
  - `GET /api/v1/users` - Benutzerliste laden
  - `POST /api/v1/users/invite` - Benutzer einladen
  - `PUT /api/v1/users/:id` - Benutzer aktualisieren (Rolle ändern)
  - `DELETE /api/v1/users/:id` - Benutzer löschen
- Einladungsformular mit Email, Vorname, Nachname, Rolle
- Bestätigungsdialog für Löschaktionen
- Rollenänderung über Dropdown-Menü

**`frontend/src/api/client.ts`**
- Fix für Tenant-Subdomain-Header bei localhost-Entwicklung
- Liest `X-Tenant-Subdomain` aus dem persistierten Auth-Storage
```typescript
function getTenantSubdomain(): string | null {
  try {
    const authStorage = localStorage.getItem('auth-storage');
    if (authStorage) {
      const parsed = JSON.parse(authStorage);
      return parsed?.state?.tenant?.subdomain || null;
    }
  } catch { }
  return null;
}
```

---

### Phase 3: Subscription/Payment Frontend Integration

#### Neue Dateien

**`frontend/src/pages/SubscriptionSuccess.tsx`**
- PayPal Success Callback-Seite
- Aktiviert Subscription nach erfolgreicher PayPal-Zahlung
- Zeigt Bestätigungsmeldung mit Tier-Upgrade-Information
- Automatische Weiterleitung zur Subscription-Seite

**`frontend/src/pages/SubscriptionCancel.tsx`**
- PayPal Cancel Callback-Seite
- Informiert Benutzer über abgebrochene Zahlung
- Option zur Rückkehr zur Subscription-Übersicht

#### Geänderte Dateien

**`frontend/src/pages/Subscription.tsx`**
- Komplette Überarbeitung mit echten API-Calls
- Anzeige des aktuellen Subscription-Status
- Upgrade-Flow mit PayPal-Weiterleitung
- Cancel-Subscription mit Bestätigungsmodal
- API-Calls implementiert:
  - `GET /api/v1/subscriptions/current` - Aktuelles Abo laden
  - `GET /api/v1/subscriptions/plans` - Verfügbare Pläne laden
  - `POST /api/v1/subscriptions/create` - PayPal-Checkout starten
  - `POST /api/v1/subscriptions/cancel` - Abo kündigen
- Visuelle Darstellung der Tier-Features
- Billing-Cycle Anzeige (Monthly/Yearly)

**`frontend/src/index.tsx`** (Routes)
- Neue Routen hinzugefügt:
  - `/subscription/success` - PayPal Success Handler
  - `/subscription/cancel` - PayPal Cancel Handler

---

### Phase 4: Enterprise Features - Email System

#### Neue Dateien

**`internal/infrastructure/external/email.go`**
- Vollständiger SMTP Email-Client mit TLS-Support
- HTML Email-Templates:
  - `user_invitation` - Benutzereinladung
  - `password_reset` - Passwort zurücksetzen
  - `welcome` - Willkommens-Email
  - `subscription_confirmation` - Abo-Bestätigung
- Template-basiertes Email-Rendering
- Konfigurierbar via Environment Variables

```go
type EmailClient struct {
    config    *EmailConfig
    templates map[string]*template.Template
}

type EmailConfig struct {
    Host       string
    Port       int
    Username   string
    Password   string
    FromEmail  string
    FromName   string
    TLS        bool
    SkipVerify bool
}

// Verfügbare Methoden:
func (c *EmailClient) Send(data *EmailData) error
func (c *EmailClient) SendUserInvitation(data *UserInvitationData) error
func (c *EmailClient) SendPasswordReset(data *PasswordResetData) error
func (c *EmailClient) SendWelcome(data *WelcomeData) error
func (c *EmailClient) SendSubscriptionConfirmation(data *SubscriptionConfirmationData) error
```

#### Geänderte Dateien

**`pkg/config/config.go`**
- Neue `SMTPConfig` Struktur hinzugefügt:
```go
type SMTPConfig struct {
    Host       string
    Port       int
    Username   string
    Password   string
    FromEmail  string
    FromName   string
    TLS        bool
    SkipVerify bool
}
```
- SMTP-Konfiguration wird aus Environment Variables geladen

**`internal/usecase/user/service.go`**
- EmailClient und BaseURL zum Service hinzugefügt
- `InviteUser` sendet jetzt echte Einladungs-Emails
- Verbesserte `generateTempPassword` mit crypto/rand für sichere Passwörter
```go
type Service struct {
    userRepo      repositories.UserRepository
    tenantRepo    repositories.TenantRepository
    auditLogRepo  repositories.AuditLogRepository
    emailClient   *external.EmailClient  // NEU
    baseURL       string                  // NEU
}
```

**`cmd/api/main.go`**
- Email-Client Initialisierung hinzugefügt
- User-Service erhält emailClient und baseURL Parameter
```go
emailClient := external.NewEmailClient(&external.EmailConfig{
    Host:       cfg.SMTP.Host,
    Port:       cfg.SMTP.Port,
    Username:   cfg.SMTP.Username,
    Password:   cfg.SMTP.Password,
    FromEmail:  cfg.SMTP.FromEmail,
    FromName:   cfg.SMTP.FromName,
    TLS:        cfg.SMTP.TLS,
    SkipVerify: cfg.SMTP.SkipVerify,
})

userService := userUsecase.NewService(
    userRepo,
    tenantRepo,
    auditLogRepo,
    emailClient,  // NEU
    cfg.App.BaseURL,  // NEU
)
```

**`docker-compose.yml`**
- SMTP Environment Variables hinzugefügt:
```yaml
SMTP_HOST: ${SMTP_HOST:-}
SMTP_PORT: ${SMTP_PORT:-587}
SMTP_USERNAME: ${SMTP_USERNAME:-}
SMTP_PASSWORD: ${SMTP_PASSWORD:-}
SMTP_FROM_EMAIL: ${SMTP_FROM_EMAIL:-noreply@gin-collection.local}
SMTP_FROM_NAME: ${SMTP_FROM_NAME:-Gin Collection}
SMTP_TLS: ${SMTP_TLS:-true}
SMTP_SKIP_VERIFY: ${SMTP_SKIP_VERIFY:-false}
```

---

### Phase 4: Database Migration CLI Tool

#### Neue Dateien

**`cmd/migrate/main.go`**
- Vollständiges Database Migration CLI Tool
- Unterstützte Commands:
  - `up` - Alle pending Migrations ausführen
  - `down` - Migrations rückgängig machen
  - `status` - Status aller Migrations anzeigen
  - `create` - Neue Migration erstellen
- Features:
  - Automatische `schema_migrations` Tabelle
  - Step-Limit für up/down Operationen
  - SQL-Statement-Splitting (respektiert String-Literale)
  - Versions-Tracking

```bash
# Verwendung:
go run cmd/migrate/main.go -command=up
go run cmd/migrate/main.go -command=down -steps=1
go run cmd/migrate/main.go -command=status
go run cmd/migrate/main.go -command=create -name="add_feature_x"

# Flags:
-host       Database host (default: localhost, env: DB_HOST)
-port       Database port (default: 3306, env: DB_PORT)
-user       Database user (default: gin_app, env: DB_USER)
-password   Database password (env: DB_PASSWORD)
-database   Database name (default: gin_collection, env: DB_NAME)
-path       Migrations path (default: ./internal/infrastructure/database/migrations)
-command    Command: up, down, status, create
-steps      Number of migrations (0 = all)
-name       Migration name (for create)
```

**`internal/infrastructure/database/migrations/003_add_invite_tokens.up.sql`**
```sql
CREATE TABLE IF NOT EXISTS invite_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_invite_tokens_token (token),
    INDEX idx_invite_tokens_email (email),
    INDEX idx_invite_tokens_expires (expires_at)
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_password_reset_token (token),
    INDEX idx_password_reset_expires (expires_at)
);
```

**`internal/infrastructure/database/migrations/003_add_invite_tokens.down.sql`**
```sql
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS invite_tokens;
```

---

### Testdaten

#### Erstellte Tenants

| ID | Name | Tier | Subdomain | Gins |
|----|------|------|-----------|------|
| 1 | Stefans Ginsammlung | enterprise | sbs | 2 |
| 2 | Free Demo | free | test123 | 5 |
| 3 | Basic Sammlung | basic | basic-demo | 10 |
| 4 | Pro Sammlung | pro | pro-demo | 25 |
| 5 | Enterprise Bar | enterprise | enterprise-demo | 30 |

#### Erstellte Benutzer

Alle Benutzer haben das Passwort: `Test123456`

| Email | Rolle | Tenant |
|-------|-------|--------|
| stefan.buell@gmail.com | owner | Stefans Ginsammlung |
| test@test.com | owner | Free Demo |
| basic@demo.local | owner | Basic Sammlung |
| pro@demo.local | owner | Pro Sammlung |
| enterprise@demo.local | owner | Enterprise Bar |
| admin@enterprise.local | admin | Enterprise Bar |
| member@enterprise.local | member | Enterprise Bar |
| viewer@enterprise.local | viewer | Enterprise Bar |

#### Testdaten pro Tier

**Free Tier (5 Gins):**
- Gordon's, Tanqueray, Bombay Sapphire, Beefeater, Hendrick's

**Basic Tier (10 Gins):**
- Gin Mare, Tanqueray No. Ten, Roku, The Botanist, Monkey 47
- Sipsmith, Plymouth, Malfy Con Limone, Elephant Gin, Star of Bombay

**Pro Tier (25 Gins):**
- Premium Auswahl inkl. Nikka Coffey, Hernö, No. 3, Martin Miller's
- Japanische, Schottische, Englische und internationale Gins

**Enterprise Tier (30 Gins):**
- Komplette Bar-Ausstattung mit Klassikern und Premium-Marken
- Ki No Bi, Ferdinand Saar, Monkey 47, Archie Rose, etc.

---

### Bug Fixes

1. **TypeScript Error: 'currentSubscription' declared but never read**
   - Fix: Variable wird jetzt im Billing-Cycle Display verwendet

2. **Migration CLI 'create' Command Database Error**
   - Fix: Create-Command wird vor Database-Verbindung separat behandelt

3. **API "Tenant not found" Errors**
   - Fix: X-Tenant-Subdomain Header wird aus persistiertem Auth-Storage gelesen

---

### Konfiguration

#### Neue Environment Variables

```bash
# SMTP Configuration (Email)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=your_password
SMTP_FROM_EMAIL=noreply@gin-collection.local
SMTP_FROM_NAME=Gin Collection
SMTP_TLS=true
SMTP_SKIP_VERIFY=false
```

---

### Login-Informationen für Tests

```bash
# Free Demo Account
Email: test@test.com
Passwort: Test123456
URL: http://localhost:5173

# Basic Account
Email: basic@demo.local
Passwort: Test123456

# Pro Account
Email: pro@demo.local
Passwort: Test123456

# Enterprise Accounts (verschiedene Rollen)
Email: enterprise@demo.local (owner)
Email: admin@enterprise.local (admin)
Email: member@enterprise.local (member)
Email: viewer@enterprise.local (viewer)
Passwort für alle: Test123456
```

---

### Nächste Schritte

1. **SMTP Server konfigurieren** - Für echte Email-Versendung
2. **PayPal Sandbox testen** - Subscription-Flow verifizieren
3. **Admin Panel erweitern** - Weitere Super-Admin Funktionen
4. **E2E Tests** - Frontend-Backend Integration testen

---

## Datei-Übersicht

### Neue Dateien (7)
1. `internal/infrastructure/external/email.go` - Email-Client
2. `cmd/migrate/main.go` - Migration CLI Tool
3. `internal/infrastructure/database/migrations/003_add_invite_tokens.up.sql`
4. `internal/infrastructure/database/migrations/003_add_invite_tokens.down.sql`
5. `frontend/src/pages/SubscriptionSuccess.tsx`
6. `frontend/src/pages/SubscriptionCancel.tsx`
7. `docs/CHANGELOG.md` - Diese Dokumentation

### Geänderte Dateien (7)
1. `frontend/src/pages/Users.tsx` - API Integration
2. `frontend/src/pages/Subscription.tsx` - API Integration
3. `frontend/src/api/client.ts` - Tenant Header Fix
4. `frontend/src/index.tsx` - Neue Routes
5. `pkg/config/config.go` - SMTP Config
6. `internal/usecase/user/service.go` - Email Integration
7. `cmd/api/main.go` - Email Client Init
8. `docker-compose.yml` - SMTP Env Vars

**Total: 14 Dateien erstellt/geändert**
