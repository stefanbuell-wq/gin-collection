# GinVault - Backlog & Open Points

> Letzte Aktualisierung: 2026-01-16

---

## In Arbeit

_Aktuell keine offenen Aufgaben_

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

### Infrastruktur
- [ ] PayPal Plan IDs in Produktion konfigurieren
- [ ] S3 Storage für Produktion einrichten
- [ ] SMTP für E-Mail-Versand konfigurieren

---

## Offen - Mittlere Priorität

### Features
- [ ] Barcode-Scanner optimieren (bessere Kamera-Unterstützung)
- [ ] Gin-Import aus CSV/Excel
- [ ] Cocktail-Rezept-Verwaltung erweitern
- [ ] Botanicals-Datenbank mit Vorschlägen
- [ ] Dark/Light Mode Toggle

### Admin Panel
- [ ] Platform Admin Dashboard erweitern
- [ ] Tenant-Statistiken verbessern
- [ ] Audit-Log für Admin-Aktionen

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

---

## Erledigt

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

## Notizen

_Platz für allgemeine Notizen und Ideen_

