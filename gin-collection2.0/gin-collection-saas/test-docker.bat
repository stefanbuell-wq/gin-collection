@echo off
echo ================================
echo Docker Installation Test
echo ================================
echo.

echo [1/5] Checking Docker installation...
docker --version
if %errorlevel% neq 0 (
    echo ERROR: Docker not found or not running
    echo Please ensure Docker Desktop is installed and running
    pause
    exit /b 1
)

echo [2/5] Checking Docker Compose...
docker compose version
if %errorlevel% neq 0 (
    echo ERROR: Docker Compose not found
    pause
    exit /b 1
)

echo [3/5] Checking Docker daemon...
docker ps >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Docker daemon not running
    echo Please start Docker Desktop
    pause
    exit /b 1
)

echo [4/5] Testing Docker with hello-world...
docker run --rm hello-world
if %errorlevel% neq 0 (
    echo ERROR: Docker run test failed
    pause
    exit /b 1
)

echo [5/5] Docker is working correctly!
echo.
echo ================================
echo Starting Gin Collection Deployment
echo ================================
echo.

cd /d "%~dp0"

echo Creating .env file...
if not exist .env (
    copy .env.example .env
    echo .env file created
) else (
    echo .env file already exists
)

echo.
echo Validating docker-compose.yml...
docker compose config --quiet
if %errorlevel% neq 0 (
    echo ERROR: docker-compose.yml validation failed
    pause
    exit /b 1
)

echo.
echo Starting services (this may take a few minutes)...
docker compose up -d

echo.
echo Waiting for services to be ready (30 seconds)...
timeout /t 30 /nobreak >nul

echo.
echo ================================
echo Service Status
echo ================================
docker compose ps

echo.
echo ================================
echo Testing API Health
echo ================================
echo.

timeout /t 5 /nobreak >nul

curl -s http://localhost:8080/health
if %errorlevel% equ 0 (
    echo.
    echo SUCCESS: API is responding!
) else (
    echo.
    echo API not yet ready, check logs with: docker compose logs api
)

echo.
echo ================================
echo Deployment Complete
echo ================================
echo.
echo Access points:
echo   - Frontend: http://localhost:3000
echo   - API: http://localhost:8080
echo   - API Health: http://localhost:8080/health
echo.
echo Useful commands:
echo   - View logs: docker compose logs -f
echo   - Stop services: docker compose down
echo   - Check status: docker compose ps
echo.

pause
