# Docker Setup Guide

## 🐳 Docker Installation & Setup

Dieses Dokument erklärt, wie Sie Docker für die Gin Collection SaaS Plattform einrichten.

## Windows

### 1. Docker Desktop installieren

1. **Download:**
   - Besuchen Sie: https://www.docker.com/products/docker-desktop
   - Klicken Sie auf "Download for Windows"
   - Speichern Sie die Datei

2. **Installation:**
   - Doppelklicken Sie auf die heruntergeladene .exe Datei
   - Folgen Sie dem Installationsassistenten
   - Akzeptieren Sie die Standardeinstellungen
   - Klicken Sie auf "Install"

3. **Neustart:**
   - Nach der Installation wird ein Neustart empfohlen
   - Starten Sie Ihren Computer neu

4. **Erste Schritte:**
   - Nach dem Neustart öffnet sich Docker Desktop automatisch
   - Akzeptieren Sie die Nutzungsbedingungen
   - Optional: Erstellen Sie ein Docker Hub Konto (nicht erforderlich)
   - Warten Sie, bis Docker vollständig gestartet ist

### 2. Docker Desktop verwenden

**Docker ist bereit, wenn:**
- ✅ Das Docker-Icon in der Taskleiste (unten rechts) **grün** ist
- ✅ Beim Klick auf das Icon steht "Docker Desktop is running"
- ✅ Sie können `docker --version` in der Kommandozeile ausführen

**So starten Sie Docker Desktop:**
- Windows-Taste drücken
- "Docker Desktop" eingeben
- Klicken Sie auf "Docker Desktop"
- Warten Sie 30-60 Sekunden

### 3. Testen

Öffnen Sie PowerShell oder CMD:

```powershell
# Docker Version prüfen
docker --version
# Erwartete Ausgabe: Docker version 29.1.3, build ...

# Docker Compose prüfen
docker compose version
# Erwartete Ausgabe: Docker Compose version v5.0.0-desktop.1

# Test-Container ausführen
docker run hello-world
# Sollte eine Willkommensnachricht anzeigen
```

## macOS

### 1. Docker Desktop installieren

**Option A: Homebrew (Empfohlen)**

```bash
brew install --cask docker
```

**Option B: Manueller Download**

1. Besuchen Sie: https://www.docker.com/products/docker-desktop
2. Klicken Sie auf "Download for Mac"
3. Wählen Sie die richtige Version:
   - **Apple Silicon (M1/M2/M3):** ARM64
   - **Intel:** AMD64
4. Öffnen Sie die .dmg Datei
5. Ziehen Sie Docker in den Applications-Ordner

### 2. Docker Desktop starten

1. Öffnen Sie Docker aus dem Applications-Ordner
2. Akzeptieren Sie die Berechtigungsanfrage
3. Warten Sie, bis Docker gestartet ist
4. Das Docker-Icon in der Menu Bar sollte erscheinen

### 3. Testen

```bash
docker --version
docker compose version
docker run hello-world
```

## Linux (Ubuntu/Debian)

### 1. Docker installieren

**Automatisches Install-Script:**

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

**Manuell:**

```bash
# Alte Versionen entfernen
sudo apt remove docker docker-engine docker.io containerd runc

# Abhängigkeiten installieren
sudo apt update
sudo apt install ca-certificates curl gnupg lsb-release

# Docker GPG Key hinzufügen
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Repository hinzufügen
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Docker installieren
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

### 2. Benutzer zu Docker-Gruppe hinzufügen

```bash
sudo usermod -aG docker $USER
```

**Wichtig:** Loggen Sie sich aus und wieder ein, damit die Änderungen wirksam werden!

### 3. Docker beim Boot starten

```bash
sudo systemctl enable docker
sudo systemctl start docker
```

### 4. Testen

```bash
docker --version
docker compose version
docker run hello-world
```

## Häufige Probleme

### Windows: "Docker daemon not running"

**Problem:** Docker Desktop wurde nicht gestartet oder ist abgestürzt.

**Lösung:**
1. Suchen Sie "Docker Desktop" im Startmenü
2. Starten Sie die Anwendung
3. Warten Sie, bis das Icon grün wird
4. Versuchen Sie es erneut

### Windows: "WSL 2 installation is incomplete"

**Problem:** Windows Subsystem for Linux (WSL 2) ist nicht installiert.

**Lösung:**
```powershell
# In PowerShell als Administrator:
wsl --install
wsl --set-default-version 2
```

Computer neu starten und Docker Desktop erneut öffnen.

### macOS: "Cannot connect to Docker daemon"

**Problem:** Docker Desktop läuft nicht.

**Lösung:**
1. Öffnen Sie Docker Desktop aus Applications
2. Warten Sie, bis das Icon in der Menu Bar erscheint
3. Versuchen Sie es erneut

### Linux: "permission denied while trying to connect"

**Problem:** Benutzer ist nicht in der docker-Gruppe.

**Lösung:**
```bash
sudo usermod -aG docker $USER
# Logout und Login
# Oder:
newgrp docker
```

### "Port already in use"

**Problem:** Ein anderer Container oder Prozess verwendet bereits die Ports.

**Lösung:**

**Windows:**
```powershell
# Prüfen welcher Prozess Port 8080 verwendet
netstat -ano | findstr "8080"

# Prozess beenden (ersetzen Sie PID)
taskkill /PID <PID> /F
```

**Linux/macOS:**
```bash
# Prüfen welcher Prozess Port 8080 verwendet
lsof -i :8080

# Prozess beenden
kill -9 <PID>
```

**Oder:** Ändern Sie die Ports in `.env`:
```env
API_PORT=8081
FRONTEND_PORT=3002
ADMIN_PORT=3003
DB_PORT=3307
```

## Docker Compose Version

Die Gin Collection SaaS verwendet **Docker Compose V2** (integriert in Docker Desktop).

**Alte Version (deprecated):**
```bash
docker-compose up -d
```

**Neue Version (empfohlen):**
```bash
docker compose up -d
```

Beide funktionieren, aber die neue Version ist schneller und besser integriert.

## Nützliche Docker Befehle

### Container verwalten

```bash
# Alle laufenden Container anzeigen
docker ps

# Alle Container anzeigen (auch gestoppte)
docker ps -a

# Container stoppen
docker stop <container-name>

# Container starten
docker start <container-name>

# Container neu starten
docker restart <container-name>

# Container löschen
docker rm <container-name>

# Container logs anzeigen
docker logs <container-name>
docker logs -f <container-name>  # Live logs
```

### Images verwalten

```bash
# Alle Images anzeigen
docker images

# Image herunterladen
docker pull mysql:8.0

# Image löschen
docker rmi <image-name>

# Ungenutzte Images löschen
docker image prune
```

### System aufräumen

```bash
# Alle gestoppten Container löschen
docker container prune

# Alle ungenutzten Images löschen
docker image prune -a

# Alle ungenutzten Volumes löschen
docker volume prune

# Alles aufräumen (Vorsicht!)
docker system prune -a --volumes
```

### Netzwerk & Volumes

```bash
# Netzwerke anzeigen
docker network ls

# Volumes anzeigen
docker volume ls

# Volume inspizieren
docker volume inspect <volume-name>
```

## Docker Desktop Einstellungen

### Ressourcen anpassen

Docker Desktop → Settings → Resources:

**Empfohlene Einstellungen:**
- **CPUs:** 4 (minimum 2)
- **Memory:** 4 GB (minimum 2 GB)
- **Swap:** 1 GB
- **Disk:** 60 GB

**Für Gin Collection SaaS:**
- MySQL benötigt ca. 1-2 GB RAM
- Redis benötigt ca. 100-500 MB RAM
- API benötigt ca. 200-500 MB RAM
- Frontend (Nginx) benötigt ca. 50 MB RAM
- Admin-Frontend (Nginx) benötigt ca. 50 MB RAM
- **Gesamt:** Ca. 2-3 GB RAM

### File Sharing (Windows/macOS)

Docker Desktop → Settings → Resources → File Sharing:

Stellen Sie sicher, dass Ihr Projektordner freigegeben ist:
```
E:\Web-Projekte\Gin-App\gin-collection2.0
```

## Performance Optimierung

### Windows (WSL 2)

**`.wslconfig` erstellen:**

Erstellen Sie: `C:\Users\<YourUsername>\.wslconfig`

```ini
[wsl2]
memory=4GB
processors=4
swap=1GB
```

Computer neu starten.

### macOS

**VirtioFS aktivieren:**

Docker Desktop → Settings → General:
- ✅ Enable VirtioFS accelerated directory sharing

### Linux

**Storage Driver optimieren:**

In `/etc/docker/daemon.json`:

```json
{
  "storage-driver": "overlay2"
}
```

```bash
sudo systemctl restart docker
```

## Sicherheit

### Docker Socket Zugriff (Linux)

**Problem:** Root-Zugriff erforderlich.

**Besser:** Benutzer zur docker-Gruppe hinzufügen (siehe oben).

### Container-Isolation

Docker-Container laufen isoliert, aber teilen sich den Kernel mit dem Host.

**Best Practices:**
- ✅ Verwenden Sie non-root User in Containern (bereits in unseren Dockerfiles)
- ✅ Limitieren Sie Ressourcen (CPU, RAM)
- ✅ Verwenden Sie private Netzwerke
- ✅ Scannen Sie Images auf Schwachstellen

### Image Scanning

```bash
# Mit Docker Scout (integriert in Docker Desktop)
docker scout cves mysql:8.0

# Mit Trivy
docker run aquasec/trivy image mysql:8.0
```

## Weiterführende Ressourcen

- **Docker Dokumentation:** https://docs.docker.com/
- **Docker Desktop:** https://docs.docker.com/desktop/
- **Docker Compose:** https://docs.docker.com/compose/
- **Best Practices:** https://docs.docker.com/develop/dev-best-practices/

## Support

Bei Problemen mit Docker:

1. **Docker Desktop Logs:** Settings → Troubleshoot → Show logs
2. **Community Forum:** https://forums.docker.com/
3. **Stack Overflow:** Tag `docker`
4. **GitHub Issues:** https://github.com/docker/for-win/issues (Windows)
5. **GitHub Issues:** https://github.com/docker/for-mac/issues (macOS)

---

**Zurück zur Hauptdokumentation:** [DEPLOYMENT.md](DEPLOYMENT.md)
