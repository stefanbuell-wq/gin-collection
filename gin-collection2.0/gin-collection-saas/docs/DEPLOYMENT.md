# Deployment Guide

This guide covers deploying the Gin Collection SaaS application to production.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start (Local Development)](#quick-start-local-development)
- [Production Deployment](#production-deployment)
- [Database Setup](#database-setup)
- [Monitoring](#monitoring)
- [Backup & Restore](#backup--restore)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **Operating System:** Linux (Ubuntu 22.04 LTS recommended)
- **CPU:** 2+ cores
- **RAM:** 4GB minimum (8GB recommended)
- **Disk:** 20GB minimum (SSD recommended)
- **Docker:** 24.0+
- **Docker Compose:** 2.20+

### Required Services

- **MySQL:** 8.0+ (Managed database recommended for production)
- **Redis:** 7.0+ (Managed cache recommended for production)
- **S3-Compatible Storage:** AWS S3, MinIO, or compatible service
- **PayPal Business Account:** For subscription payments

### Domain & SSL

- Domain name with DNS control
- SSL certificate (Let's Encrypt recommended)

## Quick Start (Local Development)

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/gin-collection-saas.git
cd gin-collection-saas
```

### 2. Setup Environment

```bash
cp .env.example .env
# Edit .env with your local configuration
```

### 3. Start Services

```bash
docker-compose up -d
```

### 4. Run Migrations

```bash
./scripts/migrate.sh up
```

### 5. Access Application

- **Frontend:** http://localhost:3000
- **API:** http://localhost:8080
- **API Docs:** http://localhost:8080/swagger
- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3001

### 6. Create First Tenant

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

## Production Deployment

### Option 1: Automated Setup (Ubuntu/Debian)

```bash
# Run as root or with sudo
sudo ./scripts/setup_prod.sh
```

This script will:
- ✅ Install Docker and Docker Compose
- ✅ Create directory structure
- ✅ Setup environment variables
- ✅ Configure SSL certificates (Let's Encrypt)
- ✅ Setup firewall rules
- ✅ Create systemd service
- ✅ Setup log rotation

### Option 2: Manual Setup

#### Step 1: Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### Step 2: Application Setup

```bash
# Create directory
sudo mkdir -p /opt/gin-collection
cd /opt/gin-collection

# Copy files
sudo cp docker-compose.prod.yml /opt/gin-collection/
sudo cp -r monitoring /opt/gin-collection/

# Setup environment
sudo cp .env.example /opt/gin-collection/.env
sudo nano /opt/gin-collection/.env  # Edit with production values
```

#### Step 3: SSL Setup

```bash
# Install certbot
sudo apt install certbot

# Get certificate
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# Link certificates
sudo mkdir -p /opt/gin-collection/ssl
sudo ln -s /etc/letsencrypt/live/yourdomain.com/fullchain.pem /opt/gin-collection/ssl/cert.pem
sudo ln -s /etc/letsencrypt/live/yourdomain.com/privkey.pem /opt/gin-collection/ssl/key.pem
```

#### Step 4: Start Services

```bash
cd /opt/gin-collection
docker-compose -f docker-compose.prod.yml up -d
```

#### Step 5: Run Migrations

```bash
./scripts/migrate.sh up
```

#### Step 6: Verify Deployment

```bash
./scripts/health_check.sh
```

## Database Setup

### Local MySQL (Development)

MySQL runs in Docker container automatically with `docker-compose up`.

### Managed MySQL (Production Recommended)

#### AWS RDS

1. Create MySQL 8.0 instance
2. Configure security group (allow port 3306 from API servers)
3. Create database: `CREATE DATABASE gin_collection;`
4. Update `.env`:
   ```
   DB_HOST=your-rds-endpoint.rds.amazonaws.com
   DB_PORT=3306
   DB_USER=admin
   DB_PASSWORD=your-secure-password
   DB_NAME=gin_collection
   ```

#### DigitalOcean Managed Database

1. Create MySQL 8.0 cluster
2. Add trusted sources (API server IPs)
3. Create database via control panel
4. Download CA certificate
5. Update `.env` with connection details

### Running Migrations

```bash
# Apply all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down 1

# Check current version
./scripts/migrate.sh version

# Force specific version (use carefully!)
./scripts/migrate.sh force 5
```

## Monitoring

### Prometheus + Grafana

Metrics and monitoring are included in `docker-compose.prod.yml`.

**Access Grafana:**
- URL: http://your-server:3001
- Default user: `admin`
- Password: Set in `.env` (`GRAFANA_PASSWORD`)

**Pre-configured Dashboards:**
- API Performance
- Database Metrics
- Redis Metrics
- System Resources

### Key Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `http_requests_total` | Total HTTP requests | N/A |
| `http_request_duration_seconds` | Request latency | P95 > 500ms |
| `http_requests_errors` | Error rate | > 5% |
| `db_connections_in_use` | DB connections | > 90% of max |
| `tenant_gin_count` | Gins per tenant | Approaching limit |
| `s3_upload_errors_total` | S3 upload failures | > 5 in 10min |

### Alerts

Alerts are configured in `monitoring/alerts/api_alerts.yml`.

**Critical Alerts:**
- API Down (> 1 minute)
- High Error Rate (> 5%)
- Database Connection Pool Exhaustion

**Warning Alerts:**
- High Response Time (P95 > 500ms)
- High Memory Usage (> 512MB)
- Multiple PayPal Webhook Failures

### Setting up Alertmanager (Optional)

```yaml
# monitoring/alertmanager.yml
global:
  slack_api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'

route:
  receiver: 'slack-notifications'
  group_by: ['alertname']

receivers:
  - name: 'slack-notifications'
    slack_configs:
      - channel: '#alerts'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ end }}'
```

## Backup & Restore

### Automated Backups

Setup cron job for daily backups:

```bash
# Edit crontab
sudo crontab -e

# Add daily backup at 2 AM
0 2 * * * /opt/gin-collection/scripts/backup.sh >> /var/log/gin-collection/backup.log 2>&1
```

### Manual Backup

```bash
sudo /opt/gin-collection/scripts/backup.sh
```

Backups are stored in `/opt/gin-collection/backups/` and retained for 30 days.

**Backup includes:**
- MySQL database (all databases)
- Redis data
- Application configuration

### Restore from Backup

```bash
# List available backups
ls -lh /opt/gin-collection/backups/

# Restore specific backup
sudo ./scripts/restore.sh 20240114_020000
```

**⚠️ Warning:** Restore will overwrite existing data!

### S3 Backup (Recommended for Production)

Set `S3_BACKUP_BUCKET` in `.env` to enable automatic S3 uploads:

```bash
S3_BACKUP_BUCKET=gin-collection-backups
```

Backups will be uploaded to:
- `s3://gin-collection-backups/backups/mysql/`
- `s3://gin-collection-backups/backups/redis/`
- `s3://gin-collection-backups/backups/app/`

## Environment Variables

### Required Variables

```bash
# Database
DB_HOST=mysql                           # Database host
DB_PORT=3306                            # Database port
DB_USER=gin_app                         # Database user
DB_PASSWORD=CHANGE_ME                   # Database password
DB_NAME=gin_collection                  # Database name

# JWT
JWT_SECRET=CHANGE_ME_256_BIT_SECRET     # 256-bit random secret
JWT_EXPIRATION=24h                      # Token expiration

# S3
S3_BUCKET=gin-collection-photos         # S3 bucket name
S3_REGION=eu-central-1                  # AWS region
S3_ACCESS_KEY=YOUR_ACCESS_KEY           # AWS access key
S3_SECRET_KEY=YOUR_SECRET_KEY           # AWS secret key

# PayPal
PAYPAL_CLIENT_ID=YOUR_CLIENT_ID         # PayPal client ID
PAYPAL_CLIENT_SECRET=YOUR_SECRET        # PayPal secret
PAYPAL_MODE=live                        # live or sandbox
PAYPAL_WEBHOOK_ID=YOUR_WEBHOOK_ID       # Webhook ID
```

### Optional Variables

```bash
# Application
APP_ENV=production                      # Environment
LOG_LEVEL=info                          # Log level
APP_BASE_URL=https://yourdomain.com     # Base URL

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_HOUR=100        # Default rate limit

# CORS
CORS_ALLOWED_ORIGINS=https://yourdomain.com  # Allowed origins

# Monitoring
GRAFANA_PASSWORD=CHANGE_ME              # Grafana admin password
```

## Systemd Service

### Service Management

```bash
# Start service
sudo systemctl start gin-collection

# Stop service
sudo systemctl stop gin-collection

# Restart service
sudo systemctl restart gin-collection

# Check status
sudo systemctl status gin-collection

# Enable auto-start on boot
sudo systemctl enable gin-collection

# View logs
sudo journalctl -u gin-collection -f
```

### Service File

Located at `/etc/systemd/system/gin-collection.service`:

```ini
[Unit]
Description=Gin Collection SaaS
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/gin-collection
ExecStart=/usr/local/bin/docker-compose -f docker-compose.prod.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.prod.yml down

[Install]
WantedBy=multi-user.target
```

## SSL/TLS Configuration

### Let's Encrypt (Recommended)

```bash
# Initial setup
sudo certbot certonly --standalone -d yourdomain.com

# Auto-renewal (already setup by certbot)
sudo certbot renew --dry-run
```

### Manual Certificate

```bash
# Copy certificate files
sudo cp fullchain.pem /opt/gin-collection/ssl/cert.pem
sudo cp privkey.pem /opt/gin-collection/ssl/key.pem
```

### Nginx Configuration for SSL

If using external nginx (not dockerized):

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /opt/gin-collection/ssl/cert.pem;
    ssl_certificate_key /opt/gin-collection/ssl/key.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## PayPal Integration

### 1. Create PayPal App

1. Login to [PayPal Developer](https://developer.paypal.com/)
2. Go to My Apps & Credentials
3. Create new app
4. Note Client ID and Secret

### 2. Create Subscription Plans

```bash
# Use PayPal API or Dashboard to create plans
# Plan IDs should match in your code:
# - PLAN_BASIC_MONTHLY
# - PLAN_BASIC_YEARLY
# - PLAN_PRO_MONTHLY
# - PLAN_PRO_YEARLY
```

### 3. Setup Webhooks

1. In PayPal Dashboard, go to Webhooks
2. Create webhook with URL: `https://yourdomain.com/api/v1/webhooks/paypal`
3. Subscribe to events:
   - `BILLING.SUBSCRIPTION.ACTIVATED`
   - `BILLING.SUBSCRIPTION.CANCELLED`
   - `BILLING.SUBSCRIPTION.SUSPENDED`
   - `PAYMENT.SALE.COMPLETED`
4. Note Webhook ID and add to `.env`

### 4. Test Webhooks

```bash
# Use PayPal's webhook simulator in the dashboard
# Or test with curl:
curl -X POST https://yourdomain.com/api/v1/webhooks/paypal \
  -H "Content-Type: application/json" \
  -d @test_webhook.json
```

## Scaling

### Horizontal Scaling (Multiple API Instances)

Update `docker-compose.prod.yml`:

```yaml
api:
  deploy:
    replicas: 3  # Run 3 instances
```

Add load balancer (nginx, HAProxy, etc.).

### Database Scaling

- **Read Replicas:** Configure MySQL read replicas for read-heavy workloads
- **Connection Pooling:** Already configured in Go code
- **Caching:** Redis caching enabled for expensive queries

### Storage Scaling

- **S3:** Auto-scales, no configuration needed
- **CDN:** Add CloudFlare or CloudFront for photo delivery

## Security Checklist

- [ ] Change all default passwords in `.env`
- [ ] Generate strong JWT secret (256-bit)
- [ ] Enable firewall (UFW) and allow only necessary ports
- [ ] Setup SSL/TLS certificates
- [ ] Configure CORS to only allow your domain
- [ ] Review database user permissions
- [ ] Enable audit logging
- [ ] Setup automated backups
- [ ] Configure fail2ban for SSH
- [ ] Enable MySQL SSL connection
- [ ] Setup monitoring alerts
- [ ] Review PayPal webhook signature verification
- [ ] Enable rate limiting
- [ ] Regular security updates (`apt update && apt upgrade`)

## Troubleshooting

### API Won't Start

```bash
# Check logs
docker logs gin-collection-api

# Common issues:
# - Database connection failed: Check DB_HOST, DB_PASSWORD
# - Port already in use: Change API_PORT in .env
# - Missing environment variables: Review .env file
```

### Database Connection Issues

```bash
# Test database connection
docker exec gin-collection-mysql mysql -u root -p -e "SELECT 1"

# Check if database exists
docker exec gin-collection-mysql mysql -u root -p -e "SHOW DATABASES"

# Verify environment variables
docker exec gin-collection-api env | grep DB_
```

### High Memory Usage

```bash
# Check container stats
docker stats

# Restart API with memory limit
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

### Slow API Response

```bash
# Check database queries
# Look for N+1 queries or missing indexes

# Enable query logging in MySQL
docker exec gin-collection-mysql mysql -u root -p -e "SET GLOBAL general_log = 'ON'"

# Check Prometheus for slow endpoints
# Look at http_request_duration_seconds metric
```

### PayPal Webhooks Not Working

```bash
# Check webhook handler logs
docker logs gin-collection-api | grep webhook

# Verify webhook signature
# Check PAYPAL_WEBHOOK_ID is correct

# Test webhook manually
curl -X POST http://localhost:8080/api/v1/webhooks/paypal \
  -H "Content-Type: application/json" \
  -d '{"event_type": "BILLING.SUBSCRIPTION.ACTIVATED"}'
```

### S3 Upload Failures

```bash
# Check S3 credentials
docker exec gin-collection-api env | grep S3_

# Test S3 access
aws s3 ls s3://gin-collection-photos/ \
  --region eu-central-1

# Check upload errors in logs
docker logs gin-collection-api | grep "S3 upload"
```

## Performance Optimization

### Database

- Ensure all tables have proper indexes
- Enable query caching in MySQL
- Use read replicas for reporting queries
- Regular `ANALYZE TABLE` and `OPTIMIZE TABLE`

### API

- Enable Go's built-in profiler: `import _ "net/http/pprof"`
- Use connection pooling (already configured)
- Enable Redis caching for expensive queries
- Use CDN for static assets

### Frontend

- Enable gzip compression in nginx
- Use browser caching for assets
- Lazy load images
- Code splitting with React.lazy()

## Maintenance

### Regular Tasks

**Daily:**
- Monitor error logs
- Check health endpoints
- Review alert notifications

**Weekly:**
- Review disk space usage
- Check backup completion
- Review database slow query log
- Update dependencies (security patches)

**Monthly:**
- Review access logs
- Update SSL certificates (if not auto-renewed)
- Database optimization (ANALYZE/OPTIMIZE)
- Review and archive old audit logs

**Quarterly:**
- Security audit
- Performance review
- Capacity planning
- Disaster recovery drill

## Support

For issues or questions:

1. Check this documentation
2. Review logs: `docker logs gin-collection-api`
3. Run health check: `./scripts/health_check.sh`
4. Check GitHub Issues: https://github.com/yourusername/gin-collection-saas/issues

## License

See LICENSE file for details.
