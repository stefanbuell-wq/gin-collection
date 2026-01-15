# Changelog

Alle wichtigen Änderungen an der Gin Collection SaaS Platform werden hier dokumentiert.

---

## [1.1.0] - 2026-01-15

### Added - Super-Admin Platform

**Zeitstempel:** 2026-01-15 08:55 UTC

#### Backend (Go API)

- **Database Migration** (`internal/infrastructure/database/migrations/002_platform_admin.up.sql`)
  - Neue `platform_admins` Tabelle für Platform-Administratoren
  - Separates Authentifizierungssystem für Admins
  - Default Admin-Account: `admin@gin-collection.local` / `admin123`

- **Domain Model** (`internal/domain/models/platform_admin.go`)
  - `PlatformAdmin` Struct mit ID, Email, PasswordHash, Name, IsActive, LastLoginAt
  - `PlatformStats` Struct für Dashboard-Statistiken
  - `SystemHealth` Struct für Health-Checks

- **Repository** (`internal/repository/mysql/platform_admin_repository.go`)
  - `GetByEmail()` - Admin per Email abrufen
  - `GetByID()` - Admin per ID abrufen
  - `UpdateLastLogin()` - Letzten Login aktualisieren
  - `GetPlatformStats()` - Platform-Statistiken abrufen
  - `GetAllTenants()` - Alle Tenants mit User/Gin-Count abrufen
  - `GetAllUsers()` - Alle User aller Tenants abrufen
  - `UpdateTenantStatus()` - Tenant suspendieren/aktivieren
  - `UpdateTenantTier()` - Subscription-Tier ändern

- **Admin Service** (`internal/usecase/admin/service.go`)
  - `Login()` - Admin-Login mit separatem JWT
  - `ValidateAdminToken()` - JWT-Validierung für Admin-Tokens
  - `GetPlatformStats()` - Dashboard-Statistiken
  - `GetAllTenants()` - Tenant-Liste
  - `SuspendTenant()` / `ActivateTenant()` - Tenant-Status ändern
  - `UpdateTenantTier()` - Tier ändern

- **Middleware** (`internal/delivery/http/middleware/platform_admin.go`)
  - `RequirePlatformAdmin()` - Prüft `is_platform_admin` JWT-Claim
  - Separate JWT-Claims für Platform-Admins

- **API Handler** (`internal/delivery/http/handler/admin/handler.go`)
  - `POST /admin/api/v1/auth/login` - Admin-Login
  - `GET /admin/api/v1/auth/me` - Aktueller Admin
  - `GET /admin/api/v1/stats` - Platform-Statistiken
  - `GET /admin/api/v1/tenants` - Tenant-Liste (paginiert)
  - `POST /admin/api/v1/tenants/:id/suspend` - Tenant suspendieren
  - `POST /admin/api/v1/tenants/:id/activate` - Tenant aktivieren
  - `PUT /admin/api/v1/tenants/:id/tier` - Tier ändern
  - `GET /admin/api/v1/users` - User-Liste (paginiert)
  - `GET /admin/api/v1/health` - System-Health

- **Router** (`internal/delivery/http/router/admin_router.go`)
  - Alle Admin-Routes unter `/admin/api/v1/`
  - Middleware-geschützte Endpoints

- **Main Entry** (`cmd/api/main.go`)
  - Admin-Router initialisiert und gemountet

#### Frontend (React/TypeScript)

- **Neues Projekt** (`admin-frontend/`)
  - React 18 + TypeScript
  - Vite als Build-Tool
  - Tailwind CSS für Styling
  - React Router v6 für Navigation

- **Komponenten**
  - `Login.tsx` - Admin-Login mit JWT-Token-Speicherung
  - `Dashboard.tsx` - Statistiken-Dashboard mit Cards
  - `Tenants.tsx` - Tenant-Verwaltung (Liste, Tier ändern, Suspend/Activate)
  - `Users.tsx` - User-Übersicht aller Tenants
  - `Layout.tsx` - Sidebar-Navigation

- **API Client** (`admin-frontend/src/api.ts`)
  - Axios-basierter Client mit JWT-Interceptor
  - Alle Admin-API-Endpoints

#### Docker/Deployment

- **Dockerfile** (`Dockerfile.admin-frontend`)
  - Multi-stage Build (Node.js → Nginx)
  - Non-root User für Security
  - Health-Check konfiguriert

- **Nginx Config** (`docker/nginx-admin.conf`)
  - Separate Konfiguration für Admin-Panel
  - API-Proxy zu Backend
  - Static Asset Caching

- **Docker Compose** (`docker-compose.yml`)
  - Neuer Service `admin-frontend`
  - Port 3001 für Admin-Panel
  - CORS-Origins erweitert

### Changed

- **docker-compose.yml**
  - CORS_ALLOWED_ORIGINS um `http://localhost:3001` und `http://admin.localhost:3001` erweitert

### Security

- Separate JWT-Claims für Platform-Admins (`is_platform_admin: true`)
- Admin-Routen nur über `/admin/api/...` erreichbar
- Strikte Middleware-Prüfung auf allen Admin-Endpoints
- bcrypt Password-Hashing

---

## [1.0.0] - 2026-01-14

### Initial Release - Gin Collection SaaS Platform

**Alle 10 Phasen abgeschlossen:**

- Phase 1-2: Foundation (Clean Architecture, Docker Setup)
- Phase 3-4: Authentication & Core API
- Phase 5: Subscription System (PayPal)
- Phase 6: Advanced Features (Botanicals, Cocktails, Photos, AI)
- Phase 7: Enterprise Features (Multi-User, API Keys, Audit)
- Phase 8: Frontend (React 18, PWA)
- Phase 9: Testing & QA
- Phase 10: Deployment

---

## Versionsschema

Dieses Projekt verwendet [Semantic Versioning](https://semver.org/):

- **MAJOR**: Inkompatible API-Änderungen
- **MINOR**: Neue Features (abwärtskompatibel)
- **PATCH**: Bugfixes (abwärtskompatibel)
