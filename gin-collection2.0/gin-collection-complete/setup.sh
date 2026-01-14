#!/bin/bash

# Gin Collection - Setup Script
# Dieses Script bereitet die Installation vor

echo "🍸 Gin Collection - Installation"
echo "=================================="
echo ""

# Check PHP version
echo "Prüfe PHP Version..."
PHP_VERSION=$(php -v | head -n 1 | cut -d " " -f 2 | cut -d "." -f 1,2)
echo "PHP Version: $PHP_VERSION"

if (( $(echo "$PHP_VERSION < 7.4" | bc -l) )); then
    echo "⚠️  Warnung: PHP 7.4+ wird empfohlen!"
fi

# Check SQLite
echo ""
echo "Prüfe SQLite..."
if php -m | grep -q sqlite3; then
    echo "✓ SQLite3 ist verfügbar"
else
    echo "✗ SQLite3 ist nicht verfügbar - Installation wird fehlschlagen!"
fi

# Create necessary directories
echo ""
echo "Erstelle Verzeichnisse..."
mkdir -p db uploads

# Set permissions
echo "Setze Berechtigungen..."
chmod 755 db uploads
chmod 644 .htaccess

echo ""
echo "✓ Setup abgeschlossen!"
echo ""
echo "Nächste Schritte:"
echo "1. Lade alle Dateien auf deinen Webserver hoch"
echo "2. Stelle sicher, dass /db/ und /uploads/ beschreibbar sind"
echo "3. Rufe die Webseite auf um die App zu starten"
echo "4. Erstelle Icons für die PWA (192x192 und 512x512 px)"
echo ""
echo "Dokumentation: siehe README.md"
echo ""
