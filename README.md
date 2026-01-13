# 🍸 Gin Collection - Progressive Web App

Eine moderne Progressive Web App zur Verwaltung deiner Gin-Sammlung mit Barcode-Scanner, Statistiken und Offline-Funktionalität.

## Features

### ✨ Hauptfunktionen
- 📸 **Barcode-Scanner**: Scanne Barcodes mit der Kamera und hole automatisch Produktinfos
- 📦 **Sammlung verwalten**: Übersichtliche Darstellung aller Gins mit Filtern und Sortierung
- ⭐ **Bewertungssystem**: Bewerte deine Gins mit 1-5 Sternen
- 📊 **Statistiken**: Detaillierte Übersichten über deine Sammlung
- 📷 **Foto-Upload**: Speichere Fotos deiner Flaschen
- 🔍 **Suche**: Durchsuche deine Sammlung nach Name, Marke, Land oder Notizen
- 📝 **Verkostungsnotizen**: Halte deine Tasting-Erlebnisse fest
- 📱 **PWA**: Installierbar als App auf dem Smartphone
- 🔄 **Offline-fähig**: Funktioniert auch ohne Internetverbindung

### 💾 Datenverwaltung
- Name, Marke, Land, Region
- Alkoholgehalt (ABV)
- Flaschengröße
- Preis & Kaufdatum
- Barcode
- Bewertung (1-5 Sterne)
- Verkostungsnotizen
- Beschreibung
- Foto
- Status (verfügbar/ausgetrunken)

## Installation

### Voraussetzungen
- Webserver mit PHP 7.4+ (df.eu unterstützt dies)
- SQLite-Unterstützung (standardmäßig in PHP enthalten)
- Mod_rewrite aktiviert (für saubere URLs)

### Schritt-für-Schritt Installation bei df.eu

1. **Upload der Dateien**
   - Lade alle Dateien per FTP auf deinen df.eu Webspace hoch
   - Platziere sie im Root-Verzeichnis oder in einem Unterordner (z.B. `/gin-collection/`)

2. **Verzeichnis-Berechtigungen**
   - Stelle sicher, dass folgende Verzeichnisse beschreibbar sind (chmod 755 oder 775):
     ```
     /db/
     /uploads/
     ```

3. **Datenbank initialisieren**
   - Die Datenbank wird automatisch beim ersten Aufruf erstellt
   - Die Datei wird in `/db/gin_collection.db` angelegt

4. **Icons erstellen** (optional)
   - Erstelle App-Icons für die PWA:
     - 192x192px: `/assets/images/icon-192.png`
     - 512x512px: `/assets/images/icon-512.png`
   - Du kannst einfache Platzhalter-Icons verwenden oder eigene gestalten

5. **HTTPS aktivieren** (empfohlen)
   - Für volle PWA-Funktionalität sollte HTTPS aktiviert sein
   - df.eu bietet kostenlose SSL-Zertifikate über Let's Encrypt

6. **Testen**
   - Rufe deine URL auf: `https://deine-domain.de/gin-collection/`
   - Die App sollte sofort funktionieren

## Projektstruktur

```
gin-collection/
├── index.html              # Hauptseite
├── manifest.json           # PWA Manifest
├── service-worker.js       # Service Worker für Offline-Funktionalität
├── .htaccess              # Apache Konfiguration
├── api/
│   ├── index.php          # API Endpoints
│   └── Database.php       # Datenbank-Klasse
├── assets/
│   ├── css/
│   │   └── style.css      # Stylesheet
│   ├── js/
│   │   ├── app.js         # Haupt-JavaScript
│   │   └── scanner.js     # Barcode-Scanner
│   └── images/
│       ├── icon-192.png   # PWA Icon (klein)
│       └── icon-512.png   # PWA Icon (groß)
├── db/
│   ├── schema.sql         # Datenbank-Schema
│   └── gin_collection.db  # SQLite Datenbank (wird automatisch erstellt)
└── uploads/               # Hochgeladene Fotos
```

## API Endpoints

Die API ist über `/api/index.php` erreichbar:

- `GET /api/?action=list` - Liste aller Gins
  - Parameter: `filter` (all|available|finished), `sort` (name|rating|price|country|date)
- `GET /api/?action=get&id=X` - Einzelnen Gin abrufen
- `POST /api/?action=add` - Neuen Gin hinzufügen
- `POST /api/?action=update` - Gin aktualisieren
- `POST /api/?action=delete` - Gin löschen
- `GET /api/?action=stats` - Statistiken abrufen
- `GET /api/?action=search&q=X` - Suche
- `GET /api/?action=barcode&code=X` - Barcode-Lookup
- `POST /api/?action=upload` - Foto hochladen

## Verwendung

### Gin hinzufügen
1. Klicke auf "Hinzufügen" in der Navigation
2. Optional: Klicke auf "Barcode scannen" um Produktinfos zu laden
3. Fülle die Formularfelder aus
4. Optional: Füge ein Foto hinzu
5. Klicke auf "Speichern"

### Barcode-Scanner
- Der Scanner verwendet die Kamera deines Geräts
- Halte den Barcode in den Kamerarahmen
- Die App sucht automatisch nach Produktinformationen
- Falls der Gin bereits existiert, wirst du gefragt ob du ihn ansehen möchtest

### Als App installieren
**Android:**
1. Öffne die Website in Chrome
2. Tippe auf das Menü (⋮) → "Zum Startbildschirm hinzufügen"

**iOS:**
1. Öffne die Website in Safari
2. Tippe auf Teilen → "Zum Home-Bildschirm"

**Desktop (Chrome/Edge):**
1. Öffne die Website
2. Klicke auf das ⊕ Icon in der Adressleiste
3. Oder: Menü → "App installieren"

## Technologie-Stack

- **Frontend**: HTML5, CSS3, Vanilla JavaScript
- **Backend**: PHP 8+
- **Datenbank**: SQLite
- **Barcode-Scanner**: Quagga2
- **PWA**: Service Worker, Web App Manifest
- **APIs**: Open Food Facts (Produktdaten)

## Barcode-Scanner Unterstützung

Die App unterstützt folgende Barcode-Formate:
- EAN-13 (Standard europäische Barcodes)
- EAN-8
- UPC-A
- UPC-E
- Code 128
- Code 39

## Datenschutz & Sicherheit

- Alle Daten werden lokal auf deinem Server gespeichert
- Keine Weitergabe an Dritte
- Die `.htaccess` Datei schützt sensible Dateien
- Datenbank-Verzeichnis ist nicht öffentlich zugänglich
- Uploaded Fotos sollten optional mit zusätzlichem Passwortschutz versehen werden

## Erweiterungsmöglichkeiten

### Geplante Features (optional)
- [ ] Import/Export (CSV, JSON)
- [ ] Backup-Funktion
- [ ] Mehrere Benutzer mit Login
- [ ] Botanicals-Datenbank
- [ ] Cocktail-Rezepte
- [ ] Sharing-Funktion
- [ ] Dark Mode
- [ ] Multi-Language Support

### Anpassungen
- **Styling**: Passe `/assets/css/style.css` an deine Wünsche an
- **API-Erweiterung**: Füge neue Endpoints in `/api/index.php` hinzu
- **Datenbank-Schema**: Erweitere `/db/schema.sql` nach Bedarf

## Troubleshooting

### Datenbank-Fehler
- Prüfe Schreibrechte auf `/db/` Verzeichnis
- Stelle sicher, dass SQLite in PHP aktiviert ist (`php -m | grep sqlite`)

### Scanner funktioniert nicht
- HTTPS ist erforderlich für Kamera-Zugriff
- Erteile Kamera-Berechtigung im Browser
- Teste mit verschiedenen Lichtverhältnissen

### Fotos werden nicht hochgeladen
- Prüfe Schreibrechte auf `/uploads/` Verzeichnis
- Prüfe PHP `upload_max_filesize` und `post_max_size` Einstellungen

### App lädt nicht
- Prüfe Browser-Konsole auf JavaScript-Fehler
- Stelle sicher, dass alle Dateien korrekt hochgeladen wurden
- Prüfe `.htaccess` Konfiguration

## Support & Feedback

Bei Fragen oder Problemen:
1. Prüfe die Troubleshooting-Sektion
2. Schaue in die Browser-Entwicklerkonsole
3. Prüfe die PHP Error Logs

## Lizenz

Dieses Projekt ist für den persönlichen Gebrauch erstellt.

## Credits

- Barcode-Scanner: [Quagga2](https://github.com/ericblade/quagga2)
- Produktdaten: [Open Food Facts](https://world.openfoodfacts.org/)

---

**Viel Spaß beim Verwalten deiner Gin-Sammlung! 🍸**
