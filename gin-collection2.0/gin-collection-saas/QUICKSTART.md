# Quick Start Guide

## Prerequisites

### 1. Install Docker Desktop

**Windows:**
1. Download Docker Desktop: https://www.docker.com/products/docker-desktop
2. Run installer
3. Restart computer
4. **Start Docker Desktop** (wichtig!)
   - Suchen Sie "Docker Desktop" im Windows-Startmenü
   - Warten Sie, bis das Icon grün ist (30-60 Sekunden)
5. Verify installation:
   ```bash
   docker --version
   docker compose version
   ```

**⚠️ WICHTIG:** Docker Desktop muss laufen, bevor Sie Services starten können!

**macOS:**
```bash
brew install --cask docker
# OR download from https://www.docker.com/products/docker-desktop
```

**Linux:**
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# Logout and login again
```

### 2. Clone Repository

```bash
git clone https://github.com/yourusername/gin-collection-saas.git
cd gin-collection-saas
```

## Local Deployment (5 Minutes)

### Option 1A: Windows Automated Script (Empfohlen)

**Für Windows-Benutzer:**

1. Stellen Sie sicher, dass Docker Desktop läuft (Icon ist grün)
2. Öffnen Sie den Projektordner im Explorer
3. Doppelklicken Sie auf: **`test-docker.bat`**

Das Script macht automatisch:
- ✅ Docker Installation Check
- ✅ Docker Daemon Check
- ✅ .env Datei erstellen
- ✅ Services starten
- ✅ Health Checks durchführen
- ✅ API testen
- ✅ Status anzeigen

### Option 1B: Linux/Mac Automated Script

```bash
# Make script executable
chmod +x test-deployment.sh

# Run deployment test
./test-deployment.sh
```

This script will:
- ✅ Check Docker installation
- ✅ Validate configuration
- ✅ Start all services
- ✅ Wait for health checks
- ✅ Show service status
- ✅ Test API endpoints

### Option 2: Manual Steps

#### Step 1: Create Environment File

```bash
cp .env.example .env
```

**Important:** For production, edit `.env` and change:
- `JWT_SECRET` to a secure 256-bit random string
- `DB_PASSWORD` to a strong password
- PayPal and S3 credentials

#### Step 2: Start Services

```bash
# Using docker compose (new)
docker compose up -d

# OR using docker-compose (old)
docker-compose up -d
```

#### Step 3: Wait for Services

Services take 30-60 seconds to be ready.

Check status:
```bash
docker compose ps
```

Expected output:
```
NAME                        STATUS              PORTS
gin-collection-api          Up (healthy)        8080/tcp
gin-collection-frontend     Up (healthy)        3000/tcp
gin-collection-mysql        Up (healthy)        3306/tcp
gin-collection-redis        Up (healthy)        6379/tcp
```

#### Step 4: Run Database Migrations

```bash
# Make script executable
chmod +x scripts/migrate.sh

# Run migrations
./scripts/migrate.sh up
```

#### Step 5: Create First Tenant

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_name": "My Gin Collection",
    "subdomain": "mygins",
    "email": "admin@example.com",
    "password": "SecurePassword123!"
  }'
```

Expected response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "role": "owner"
  },
  "tenant": {
    "id": 1,
    "name": "My Gin Collection",
    "subdomain": "mygins",
    "tier": "free"
  }
}
```

## Access the Application

### Frontend (User)
- **URL:** http://localhost:3000
- **Login:** Use the email/password from registration

### Admin Panel (Platform Admin)
- **URL:** http://localhost:3001
- **Email:** `admin@gin-collection.local`
- **Passwort:** `admin123`
- **Funktionen:** Tenant-Management, User-Übersicht, Statistiken

### API
- **URL:** http://localhost:8080
- **Health Check:** http://localhost:8080/health
- **API Docs:** http://localhost:8080/swagger (if enabled)

### Monitoring
- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3002 (wenn konfiguriert)
  - Default login: `admin` / `admin`

## Testing the Application

### 1. Create a Gin

```bash
# Save your token from registration
TOKEN="your_jwt_token_here"

curl -X POST http://localhost:8080/api/v1/gins \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hendricks Gin",
    "brand": "Hendricks",
    "country": "Scotland",
    "gin_type": "London Dry",
    "abv": 41.4,
    "bottle_size_ml": 700,
    "price": 29.99,
    "currency": "EUR",
    "availability_status": "available"
  }'
```

### 2. List Gins

```bash
curl -X GET http://localhost:8080/api/v1/gins \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Get Statistics

```bash
curl -X GET http://localhost:8080/api/v1/gins/stats \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Search Gins

```bash
curl -X GET "http://localhost:8080/api/v1/gins/search?q=hendricks" \
  -H "Authorization: Bearer $TOKEN"
```

## Useful Commands

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
docker compose logs -f mysql
docker compose logs -f redis
```

### Stop Services

```bash
docker compose down
```

### Restart Services

```bash
docker compose restart
```

### Stop and Remove Everything (including volumes)

```bash
docker compose down -v
```

### Rebuild Images

```bash
docker compose build --no-cache
docker compose up -d
```

### Access Container Shell

```bash
# API container
docker exec -it gin-collection-api sh

# MySQL container
docker exec -it gin-collection-mysql mysql -u root -pdev_root_password gin_collection

# Redis container
docker exec -it gin-collection-redis redis-cli
```

## Troubleshooting

### Services Won't Start

```bash
# Check logs
docker compose logs

# Check if ports are in use
# Windows
netstat -ano | findstr "8080"
netstat -ano | findstr "3306"

# Linux/Mac
lsof -i :8080
lsof -i :3306
```

### Database Connection Error

```bash
# Check MySQL is running
docker exec gin-collection-mysql mysqladmin ping -h localhost -u root -pdev_root_password

# Check database exists
docker exec gin-collection-mysql mysql -u root -pdev_root_password -e "SHOW DATABASES;"

# Recreate database
docker compose down
docker volume rm gin-collection-saas_mysql_data
docker compose up -d
./scripts/migrate.sh up
```

### API Won't Start

```bash
# Check API logs
docker compose logs api

# Common issues:
# 1. Database not ready - wait 30 seconds and check again
# 2. Environment variables missing - check .env file
# 3. Port 8080 in use - change API_PORT in .env
```

### Frontend Won't Load

```bash
# Check frontend logs
docker compose logs frontend

# Check nginx configuration
docker exec gin-collection-frontend cat /etc/nginx/nginx.conf

# Rebuild frontend
docker compose build frontend --no-cache
docker compose up -d frontend
```

### Permission Denied (Linux)

```bash
# If you get permission errors with Docker
sudo usermod -aG docker $USER
# Then logout and login again

# If you get permission errors with scripts
chmod +x scripts/*.sh
chmod +x test-deployment.sh
```

## Development Workflow

### 1. Backend Development

```bash
# Option A: Run in Docker (recommended)
docker compose up -d
docker compose logs -f api

# Option B: Run locally
go run cmd/api/main.go
```

### 2. Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev

# Access at http://localhost:5173
```

### 3. Running Tests

```bash
# Backend tests
go test ./...
go test -cover ./...

# Integration tests (requires Docker)
go test ./tests/integration/...

# Frontend tests
cd frontend
npm test
npm run test:coverage
```

## Next Steps

1. ✅ Explore the frontend at http://localhost:3000
2. ✅ Read the [API Documentation](docs/API.md)
3. ✅ Review [Testing Guide](tests/README.md)
4. ✅ Check [Deployment Guide](docs/DEPLOYMENT.md) for production
5. ✅ Set up PayPal sandbox for testing subscriptions
6. ✅ Configure S3 bucket for photo uploads

## Support

If you encounter issues:

1. Check the logs: `docker compose logs`
2. Run health check: `./scripts/health_check.sh`
3. Review [Troubleshooting Guide](docs/DEPLOYMENT.md#troubleshooting)
4. Create an issue on GitHub

## Clean Up

To completely remove the deployment:

```bash
# Stop and remove containers, networks, volumes
docker compose down -v

# Remove images
docker rmi gin-collection-api gin-collection-frontend

# Remove .env file (optional)
rm .env
```

---

**Happy Coding! 🚀**
