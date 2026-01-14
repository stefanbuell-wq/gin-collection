# Gin Collection SaaS - Multi-Tenant Platform

Eine moderne, Go-basierte SaaS-Plattform zur Verwaltung von Gin-Sammlungen mit Multi-Tenancy, Subscription-Tiers und PayPal-Integration.

## 🎯 Features

### Core Features
- ✅ **Multi-Tenancy**: Hybrid-Modell (Shared DB für Free/Basic/Pro, Separate DB für Enterprise)
- ✅ **4 Subscription Tiers**: Free, Basic, Pro, Enterprise
- ✅ **54 Business Features**: Vollständige Feature-Parität mit PHP-App
- ✅ **PayPal Integration**: Subscription Management
- ✅ **REST API**: Vollständige REST API mit JWT-Authentifizierung
- ✅ **Docker-Ready**: Deployment-agnostisch

### Subscription Tiers

| Feature | Free | Basic | Pro | Enterprise |
|---------|------|-------|-----|------------|
| Max Gins | 10 | 50 | ∞ | ∞ |
| Photos/Gin | 1 | 3 | 10 | 50 |
| Botanicals | ❌ | ❌ | ✅ | ✅ |
| Cocktails | ❌ | ❌ | ✅ | ✅ |
| AI Suggestions | ❌ | ❌ | ✅ | ✅ |
| Export/Import | ❌ | ❌ | ✅ | ✅ |
| Multi-User | ❌ | ❌ | ❌ | ✅ |
| API Access | ❌ | ❌ | ❌ | ✅ |
| Separate DB | ❌ | ❌ | ❌ | ✅ |

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make (optional)

### 1. Clone & Setup
```bash
cd gin-collection-saas
cp .env.example .env
# Edit .env with your configuration
```

### 2. Install Dependencies
```bash
make deps
# OR
go mod tidy
```

### 3. Start Development Environment
```bash
make dev
# OR
docker-compose -f docker/docker-compose.yml up -d
```

### 4. Run Database Migrations
```bash
# Migrations werden automatisch beim Docker-Start ausgeführt
# Oder manuell:
docker exec -i gin-mysql mysql -ugin_app -pgin_password gin_collection < internal/infrastructure/database/migrations/001_initial_schema.up.sql
```

### 5. Start API Server
```bash
make run
# OR
go run cmd/api/main.go
```

Die API ist dann verfügbar unter: `http://localhost:8080`

## 📁 Projektstruktur

```
gin-collection-saas/
├── cmd/
│   └── api/                 # API Server Entry Point
├── internal/
│   ├── domain/              # Domain Models & Business Logic
│   │   ├── models/          # Entities (Tenant, User, Gin, etc.)
│   │   ├── repositories/    # Repository Interfaces
│   │   └── errors/          # Domain Errors
│   ├── usecase/             # Business Logic
│   │   ├── gin/             # Gin CRUD, Search, Export
│   │   ├── auth/            # Authentication
│   │   ├── subscription/    # Subscription Management
│   │   └── tenant/          # Tenant Management
│   ├── delivery/            # HTTP Layer
│   │   └── http/
│   │       ├── handler/     # HTTP Handlers
│   │       ├── middleware/  # Middleware (Auth, Tenant, etc.)
│   │       └── router/      # Route Definitions
│   ├── repository/          # Data Access
│   │   └── mysql/           # MySQL Implementations
│   └── infrastructure/      # External Services
│       ├── database/        # Database Connection & Migrations
│       ├── storage/         # S3 Storage
│       └── external/        # External APIs (PayPal, etc.)
├── pkg/                     # Shared Packages
│   ├── config/              # Configuration Management
│   ├── logger/              # Logging
│   └── validator/           # Input Validation
├── docker/                  # Docker Files
└── migrations/              # Data Migrations (SQLite → MySQL)
```

## 🛠️ Development

### Available Make Commands
```bash
make help          # Show all available commands
make deps          # Install Go dependencies
make build         # Build the application
make run           # Run the application locally
make dev           # Start Docker development environment
make test          # Run tests
make test-coverage # Run tests with coverage
make lint          # Run linters
make docker-build  # Build Docker image
make docker-up     # Start all Docker services
make docker-down   # Stop all Docker services
make migrate-up    # Run database migrations
make migrate-down  # Rollback database migrations
```

### Running Tests
```bash
# Unit tests
go test ./...

# Integration tests (requires Docker)
docker-compose -f docker/docker-compose.yml up -d
go test -tags=integration ./tests/integration/...

# Coverage
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📡 API Endpoints

### Authentication
```
POST   /api/v1/auth/register     # Register new tenant + user
POST   /api/v1/auth/login        # Login
POST   /api/v1/auth/refresh      # Refresh JWT token
```

### Tenants
```
GET    /api/v1/tenants/current   # Get current tenant info
PUT    /api/v1/tenants/current   # Update tenant settings
GET    /api/v1/tenants/usage     # Get usage metrics
```

### Subscriptions
```
GET    /api/v1/subscriptions/current  # Get current subscription
GET    /api/v1/subscriptions/plans    # List available plans
POST   /api/v1/subscriptions/upgrade  # Upgrade tier
POST   /api/v1/subscriptions/cancel   # Cancel subscription
```

### Gins
```
GET    /api/v1/gins              # List gins (with filters)
POST   /api/v1/gins              # Create gin
GET    /api/v1/gins/:id          # Get gin details
PUT    /api/v1/gins/:id          # Update gin
DELETE /api/v1/gins/:id          # Delete gin
GET    /api/v1/gins/search       # Search gins
GET    /api/v1/gins/stats        # Get statistics
POST   /api/v1/gins/export       # Export data (Pro+)
```

### Botanicals (Pro+)
```
GET    /api/v1/botanicals                    # List all botanicals
GET    /api/v1/gins/:id/botanicals           # Get gin botanicals
PUT    /api/v1/gins/:id/botanicals           # Update gin botanicals
```

### Photos
```
GET    /api/v1/gins/:id/photos               # List photos
POST   /api/v1/gins/:id/photos               # Upload photo
DELETE /api/v1/gins/:id/photos/:photo_id     # Delete photo
```

## 🔐 Environment Variables

Siehe `.env.example` für alle verfügbaren Konfigurationsoptionen.

Wichtige Variablen:
- `JWT_SECRET`: **MUSS** in Production geändert werden
- `DB_PASSWORD`: Datenbank-Passwort
- `PAYPAL_CLIENT_ID` / `PAYPAL_CLIENT_SECRET`: PayPal Credentials
- `S3_BUCKET` / `AWS_ACCESS_KEY_ID`: S3 Storage Credentials

## 🐳 Docker Deployment

### Development
```bash
docker-compose -f docker/docker-compose.yml up
```

### Production
```bash
# Build
docker build -f docker/Dockerfile.api -t gin-collection-api:latest .

# Run
docker run -d \
  --name gin-api \
  -p 8080:8080 \
  --env-file .env \
  gin-collection-api:latest
```

## 📊 Database Schema

Das Schema umfasst folgende Haupttabellen:
- **tenants**: SaaS-Tenants (Kunden)
- **users**: Benutzer pro Tenant
- **subscriptions**: Subscription-Informationen
- **gins**: Gin-Sammlung (tenant-scoped)
- **botanicals**: Botanicals (shared reference data)
- **cocktails**: Cocktail-Rezepte (shared reference data)
- **gin_botanicals**, **gin_photos**, **gin_cocktails**: Relationen

Siehe `internal/infrastructure/database/migrations/001_initial_schema.up.sql` für Details.

## 🔄 Migration (SQLite → MySQL)

Für bestehende SQLite-Daten:
```bash
go run migrations/sqlite_to_mysql/migrate.go \
  --sqlite-db=/path/to/old.db \
  --mysql-dsn="user:pass@tcp(localhost:3306)/gin_collection" \
  --tenant-subdomain=mycollection
```

## 🧪 Testing Strategy

- **Unit Tests**: Repository, UseCase, Middleware
- **Integration Tests**: API Endpoints, Database Operations
- **E2E Tests**: Complete User Journeys
- **Security Tests**: Tenant Isolation, RBAC, SQL Injection

## 📈 Monitoring

### Health Checks
```bash
# Liveness
curl http://localhost:8080/health

# Readiness (checks DB + Redis)
curl http://localhost:8080/ready

# Metrics (Prometheus format)
curl http://localhost:8080/metrics
```

## 🎯 Implementation Roadmap

### ✅ Phase 1-2: Foundation (COMPLETED)
- [x] Projektstruktur & Clean Architecture
- [x] Docker Setup
- [x] MySQL Schema mit Multi-Tenancy
- [x] Config Management
- [x] Domain Models
- [x] Tenant Router

### ✅ Phase 3-4: Authentication & Core API (COMPLETED)
- [x] JWT Service
- [x] Auth Middleware
- [x] Tenant Middleware
- [x] RBAC Middleware
- [x] Gin Repository & Service
- [x] Gin HTTP Handlers
- [x] Search & Filter Logic
- [x] Statistics & Export

### ✅ Phase 5: Subscription System (COMPLETED)
- [x] PayPal Integration
- [x] Tier Enforcement Middleware
- [x] Usage Tracking
- [x] Upgrade/Downgrade Flow
- [x] Webhook Handling

### ✅ Phase 6: Advanced Features (COMPLETED)
- [x] Botanicals & Cocktails
- [x] Photo Upload (AWS S3)
- [x] AI Suggestions (Similar Gins)
- [x] Export/Import (JSON, CSV)
- [x] Barcode Scanner Integration

### ✅ Phase 7: Enterprise Features (COMPLETED)
- [x] Multi-User Support (4 Rollen)
- [x] API Key Authentication
- [x] Audit Logging
- [x] Separate DB Provisioning für Enterprise

### ✅ Phase 8: Frontend (COMPLETED)
- [x] React 18 + TypeScript SPA
- [x] PWA Service Worker
- [x] Subscription UI & Payment Flow
- [x] Zustand State Management
- [x] Responsive Design (Tailwind CSS)

### ✅ Phase 9: Testing & QA (COMPLETED)
- [x] Integration Tests (Tenant Isolation)
- [x] Tier Enforcement Tests
- [x] E2E Subscription Tests
- [x] Security Tests (OWASP Top 10)
- [x] Load Tests (k6)
- [x] Frontend Tests (Vitest)

### ✅ Phase 10: Deployment (COMPLETED)
- [x] Docker Images (Multi-stage)
- [x] Docker Compose (Dev & Prod)
- [x] GitHub Actions CI/CD
- [x] Monitoring (Prometheus + Grafana)
- [x] Backup & Restore Scripts
- [x] Production Setup Script
- [x] Deployment Documentation

## 📝 License

Proprietary - © 2026 Gin Collection SaaS

## 🤝 Support

Bei Fragen oder Problemen:
1. Prüfen Sie die Dokumentation
2. Prüfen Sie die Logs: `docker logs gin-api`
3. Erstellen Sie ein Issue im Repository

---

**Status:** ✅ **PRODUCTION READY** - All 10 Phases Complete!
