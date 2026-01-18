# GinVault - Backlog & Open Points

> Letzte Aktualisierung: 2026-01-18

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

#### PayPal Integration (Geschätzter Aufwand: ~2 Stunden)
- [ ] PayPal Developer Account erstellen (https://developer.paypal.com)
- [ ] Sandbox App erstellen (Dashboard → Apps & Credentials → Create App)
  - Client ID notieren
  - Client Secret notieren
- [ ] Billing Plans in PayPal anlegen:
  - [ ] Basic Monthly (4,99€/Monat)
  - [ ] Basic Yearly (49,99€/Jahr)
  - [ ] Pro Monthly (9,99€/Monat)
  - [ ] Pro Yearly (99,99€/Jahr)
  - [ ] Enterprise Monthly (29,99€/Monat)
  - [ ] Enterprise Yearly (299,99€/Jahr)
- [ ] Webhook einrichten (URL: /api/v1/webhooks/paypal)
  - Events: BILLING.SUBSCRIPTION.ACTIVATED, CANCELLED, SUSPENDED, PAYMENT.SALE.COMPLETED
  - Webhook ID notieren
- [ ] Environment Variables auf Server setzen:
  - PAYPAL_CLIENT_ID
  - PAYPAL_CLIENT_SECRET
  - PAYPAL_MODE=sandbox (später: live)
  - PAYPAL_WEBHOOK_ID
- [ ] Plan IDs im Code hinterlegen (internal/domain/models/subscription.go)
- [ ] Sandbox-Tests durchführen
- [ ] Live schalten (PAYPAL_MODE=live, neue Live-Credentials)

#### Weitere Infrastruktur
- [ ] S3 Storage für Produktion einrichten
- [x] SMTP für E-Mail-Versand konfiguriert (Hostinger, info@ginvault.cloud)

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

### 2026-01-18
- [x] **Mobile Double-Click Bug Fix** (GinCreate.tsx)
  - Problem: Auf Mobile musste man 2x auf "Speichern" klicken
  - Ursache: Mobile Touch-Events werden bei Form-Submit anders behandelt
  - Lösung: Button von `type="submit"` auf `type="button"` geändert mit expliziten `onClick` und `onTouchEnd` Handlern
- [x] **Upgrade Modal als Overlay** (GinCreate.tsx, GinCreate.css)
  - Modal erscheint jetzt als fixed Overlay über der Seite
  - Backdrop mit Blur-Effekt
  - Zentrierte Darstellung auf allen Geräten
- [x] **Debug-Modus für Mobile-Testing**
  - URL-Parameter `?debug=1` aktiviert sichtbares Debug-Panel
  - Zeigt letzte 20 Log-Einträge mit Timestamps
  - Hilfreich für Debugging auf Mobile-Geräten ohne DevTools
- [x] **Gin-Limit Enforcement mit Upgrade-Benachrichtigung**
  - Backend gibt `upgrade_required: true` zurück bei 403
  - Enthält: `limit`, `current_count`, `current_tier`
  - ginStore.ts extrahiert Upgrade-Info aus Fehler-Response
  - Modal zeigt aktuellen Tier und Limit an
- [x] **Password Reset Feature** (Backend + Frontend)
  - Forgot Password Seite (ForgotPassword.tsx)
  - Reset Password Seite (ResetPassword.tsx)
  - Backend-Endpoints: `/forgot-password`, `/reset-password`, `/validate-reset-token`
  - Token-basiertes Reset mit E-Mail-Versand
- [x] **GinVault E-Mail Templates**
  - Branding auf GinVault umgestellt
  - Dark Theme Design

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

## Vision / Roadmap - Langfristige Features

### 1. Automatische Etikettenerkennung (Label Recognition)
**Priorität:** Hoch
**Beschreibung:** Ein Foto der Flasche macht GinVault zum intelligenten Erkennungssystem:
- Marke automatisch erkennen
- Botanicals identifizieren
- Alkoholgehalt auslesen
- Herkunft bestimmen
- Preisrange schätzen

**USP:** Massiver UX-Boost und starker Differentiator gegenüber allen existierenden Gin-Apps.

---

### 2. KI-gestützte Aromenanalyse ("Aroma-Coach")
**Priorität:** Hoch
**Beschreibung:** GinVault als intelligenter Geschmacks-Berater:
- Nutzer geben Lieblingsgins ein
- KI erkennt Muster (z.B. "floral + citrus")
- GinVault schlägt neue Gins vor, die exakt ins Profil passen

**Vision:** Wie Spotify-Discover - nur für Gin.

---

### 3. Händler- und Brennerei-Dashboards (B2B-Modul)
**Priorität:** Mittel
**Beschreibung:** Ein B2B-Modul für Hersteller mit Analytics:
- Welche Gins werden am häufigsten gesammelt
- Welche Aromen im Trend sind
- Welche Zielgruppen welche Gins bevorzugen

**Potenzial:** Macht GinVault für die Industrie extrem wertvoll.

---

### 4. Limited Editions exklusiv für GinVault
**Priorität:** Mittel
**Beschreibung:** Kooperationen mit Brennereien:
- "GinVault Edition No. 1"
- Exklusive Batch-Releases
- Nur für Premium-Mitglieder

**Potenzial:** Schafft Begehrlichkeit und wiederkehrende Umsätze.

---

### 5. Gamification & Achievements
**Priorität:** Mittel
**Beschreibung:** Sammler lieben Status. Beispiele:
- "10 Tastings abgeschlossen"
- "Botanical-Master: 50 Aromen erkannt"
- "Rare Bottle Collector"
- Badges und Level-System

**Potenzial:** Erhöht Retention und Community-Dynamik massiv.

---

### 6. Social Features mit echtem Mehrwert
**Priorität:** Mittel
**Beschreibung:** Nicht nur Likes, sondern:
- Tasting-Vergleiche mit Freunden
- Gemeinsame Tasting-Sessions
- "Flavor Match Score" zwischen Nutzern
- Challenges ("Taste 5 Mediterranean Gins this month")

**Vision:** Das macht GinVault lebendig und community-driven.

---

### 7. Integration mit Bars & Events
**Priorität:** Mittel
**Beschreibung:** GinVault als digitaler Begleiter für reale Erlebnisse:
- Bars integrieren ihre Gin-Karte
- Nutzer scannen Gin im Restaurant
- Tasting wird automatisch gespeichert
- Events über GinVault buchen

**Vision:** Verbindet Online und Offline nahtlos.

---

### 8. Erweiterte Abo-Modelle mit echtem Mehrwert
**Priorität:** Hoch
**Beschreibung:** Premium-Features die sich "lohnen":
- Unbegrenzte Sammlung
- Exklusive Gins
- Deep-Analytics
- KI-Empfehlungen
- Early Access zu Limited Editions
- Rabattcodes bei Partnern

---

### 9. GinVault als Geschenkprodukt
**Priorität:** Niedrig
**Beschreibung:** GinVault in Geschenkboxen integrieren:
- 3 Gins + Premium-Abo für 3 Monate
- QR-Code führt direkt zum Tasting-Erlebnis

**Potenzial:** Perfektes Weihnachts- und Geburtstagsprodukt.

---

### 10. API für Shops & Hersteller
**Priorität:** Mittel
**Beschreibung:** Shops können:
- GinVault-Bewertungen anzeigen
- Aromenprofile integrieren
- "Passt zu deinem Geschmack"-Empfehlungen nutzen

**Vision:** GinVault wird zur Infrastruktur des Gin-Markts.

---

### 11. Community-Ranking & Awards
**Priorität:** Niedrig
**Beschreibung:** Jährlicher "GinVault Community Award":
- Beste Gins
- Beste Newcomer
- Beste Brennerei

**Potenzial:** Schafft Reichweite und Presse-Aufmerksamkeit.

---

### 12. GinVault für Firmen (Corporate Tasting)
**Priorität:** Mittel
**Beschreibung:** Firmen lieben Tasting-Events:
- Digitale Tasting-Boxen
- Moderierte Sessions
- Firmenprofile
- Mitarbeiter-Challenges

**Potenzial:** Extrem lukrativer B2B-Markt.

---

### 13. Erweiterte Cocktail-Features
**Priorität:** Mittel
**Beschreibung:** Viele Gin-Fans sind auch Cocktail-Fans:
- Cocktail-Rezepte basierend auf eigener Sammlung
- "Was kann ich mit meinen Gins mixen?"
- KI-Mixing-Assistent

**Potenzial:** Erweitert die Zielgruppe enorm.

---

### 14. Sammler-Wertentwicklung & Rare-Bottle-Tracking
**Priorität:** Niedrig
**Beschreibung:** Für Premium-Sammler:
- Preisentwicklung tracken
- Seltenheitsindex
- Marktwert der eigenen Sammlung

**Vision:** GinVault als "Gin-Portfolio-Tracker".

---

### 15. Tasting-Box-Ökosystem
**Priorität:** Mittel
**Beschreibung:** Boxen nicht nur verkaufen, sondern:
- Boxen mit Partnern co-branden
- Boxen als Abo anbieten
- Boxen als Onboarding-Tool für neue Nutzer

**Potenzial:** Verstärkt das Flywheel und schafft wiederkehrende Umsätze.

---

## Notizen

_Platz für allgemeine Notizen und Ideen_

