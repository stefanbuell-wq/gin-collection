# GinVault API Integration

## Übersicht

Die GinVault API ermöglicht externen Anwendungen den programmatischen Zugriff auf die Gin-Sammlung. Der API-Zugang ist für **Pro** und **Enterprise** Kunden verfügbar.

---

## Tier-Vergleich: API-Zugang

| Feature | Free | Basic | Pro | Enterprise |
|---------|------|-------|-----|------------|
| **API-Zugang** | ❌ | ❌ | ✅ | ✅ |
| **API Rate Limit** | - | - | 5.000 Req/Std | 10.000 Req/Std |
| **API Key Management** | - | - | ✅ | ✅ |
| **Webhook Support** | - | - | ❌ | ✅ |

---

## Tier-Limits Übersicht

| Feature | Free | Basic | Pro | Enterprise |
|---------|------|-------|-----|------------|
| **Max. Gins** | 25 | 100 | 500 | Unbegrenzt |
| **Fotos pro Gin** | 3 | 10 | 25 | Unbegrenzt |
| **Preis/Monat** | €0 | €4,99 | €9,99 | €29,99 |
| **Preis/Jahr** | €0 | €49,99 | €99,99 | €299,99 |

---

## Authentifizierung

### API Key Format

```
sk_<32-stelliger-alphanumerischer-key>
```

Beispiel: `sk_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`

### Header-Format

```http
Authorization: Bearer sk_your_api_key_here
```

### Beispiel-Request

```bash
curl -X GET "https://api.ginvault.de/api/v1/gins" \
  -H "Authorization: Bearer sk_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6" \
  -H "Content-Type: application/json"
```

---

## API Key Management

### API Key generieren

**Voraussetzungen:**
- Pro oder Enterprise Subscription
- Benutzerrolle: `owner` oder `admin`

**Endpoint:**
```http
POST /api/v1/users/{user_id}/api-key
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "api_key": "sk_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
    "created_at": "2026-01-16T12:00:00Z"
  }
}
```

> ⚠️ **Wichtig:** Der API Key wird nur einmal angezeigt! Speichern Sie ihn sicher ab.

### API Key widerrufen

```http
DELETE /api/v1/users/{user_id}/api-key
Authorization: Bearer <jwt_token>
```

---

## Verfügbare Endpoints

### Gins

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| `GET` | `/api/v1/gins` | Alle Gins abrufen |
| `GET` | `/api/v1/gins/{id}` | Einzelnen Gin abrufen |
| `POST` | `/api/v1/gins` | Neuen Gin erstellen |
| `PUT` | `/api/v1/gins/{id}` | Gin aktualisieren |
| `DELETE` | `/api/v1/gins/{id}` | Gin löschen |

### Botanicals (Pro & Enterprise)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| `GET` | `/api/v1/botanicals` | Alle Botanicals abrufen |
| `GET` | `/api/v1/gins/{gin_id}/botanicals` | Botanicals eines Gins |
| `POST` | `/api/v1/gins/{gin_id}/botanicals` | Botanical hinzufügen |
| `DELETE` | `/api/v1/gins/{gin_id}/botanicals/{id}` | Botanical entfernen |

### Cocktails (Pro & Enterprise)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| `GET` | `/api/v1/cocktails` | Alle Cocktails abrufen |
| `GET` | `/api/v1/cocktails/{id}` | Einzelnen Cocktail abrufen |
| `POST` | `/api/v1/cocktails` | Neuen Cocktail erstellen |
| `PUT` | `/api/v1/cocktails/{id}` | Cocktail aktualisieren |
| `DELETE` | `/api/v1/cocktails/{id}` | Cocktail löschen |

### Fotos

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| `GET` | `/api/v1/gins/{gin_id}/photos` | Fotos eines Gins |
| `POST` | `/api/v1/gins/{gin_id}/photos` | Foto hochladen |
| `DELETE` | `/api/v1/gins/{gin_id}/photos/{id}` | Foto löschen |

### Export/Import (Pro & Enterprise)

| Methode | Endpoint | Beschreibung |
|---------|----------|--------------|
| `GET` | `/api/v1/export` | Gesamte Sammlung exportieren (JSON) |
| `POST` | `/api/v1/import` | Sammlung importieren |

---

## Request/Response Formate

### Gin erstellen

**Request:**
```http
POST /api/v1/gins
Authorization: Bearer sk_your_api_key
Content-Type: application/json

{
  "name": "Monkey 47",
  "brand": "Black Forest Distillers",
  "type": "Dry Gin",
  "origin_country": "Deutschland",
  "abv": 47.0,
  "price": 34.99,
  "bottle_size": 500,
  "description": "Schwarzwald Dry Gin mit 47 Botanicals",
  "taste_notes": "Komplex, würzig, Wacholder, Zitrus, Preiselbeere",
  "rating": 5,
  "is_favorite": true
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 123,
    "name": "Monkey 47",
    "brand": "Black Forest Distillers",
    "type": "Dry Gin",
    "origin_country": "Deutschland",
    "abv": 47.0,
    "price": 34.99,
    "bottle_size": 500,
    "description": "Schwarzwald Dry Gin mit 47 Botanicals",
    "taste_notes": "Komplex, würzig, Wacholder, Zitrus, Preiselbeere",
    "rating": 5,
    "is_favorite": true,
    "fill_level": 100,
    "created_at": "2026-01-16T12:00:00Z",
    "updated_at": "2026-01-16T12:00:00Z"
  }
}
```

### Gins auflisten

**Request:**
```http
GET /api/v1/gins?page=1&limit=20&sort=name&order=asc
Authorization: Bearer sk_your_api_key
```

**Response:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 123,
        "name": "Monkey 47",
        "brand": "Black Forest Distillers",
        "type": "Dry Gin",
        "abv": 47.0,
        "rating": 5,
        "photo_url": "https://..."
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

---

## Rate Limiting

### Limits nach Tier

| Tier | Requests pro Stunde |
|------|---------------------|
| Pro | 5.000 |
| Enterprise | 10.000 |

Rate Limit Header in jeder Response:

```http
X-RateLimit-Limit: 10000
X-RateLimit-Remaining: 9950
X-RateLimit-Reset: 1705410000
```

### Rate Limit überschritten

**Response (429 Too Many Requests):**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please wait before making more requests.",
    "retry_after": 3600
  }
}
```

---

## Fehler-Responses

### Authentifizierungsfehler (401)

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or missing API key"
  }
}
```

### Zugriff verweigert (403)

```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "API access requires Pro or Enterprise subscription"
  }
}
```

### Ressource nicht gefunden (404)

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Gin with ID 999 not found"
  }
}
```

### Validierungsfehler (400)

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request data",
    "details": {
      "name": "Name is required",
      "abv": "ABV must be between 0 and 100"
    }
  }
}
```

---

## Webhooks (nur Enterprise)

Nur Enterprise-Kunden können Webhooks konfigurieren, um über Änderungen benachrichtigt zu werden.

### Verfügbare Events

| Event | Beschreibung |
|-------|--------------|
| `gin.created` | Neuer Gin wurde erstellt |
| `gin.updated` | Gin wurde aktualisiert |
| `gin.deleted` | Gin wurde gelöscht |
| `photo.uploaded` | Neues Foto hochgeladen |
| `collection.exported` | Export durchgeführt |

### Webhook Payload

```json
{
  "event": "gin.created",
  "timestamp": "2026-01-16T12:00:00Z",
  "data": {
    "id": 123,
    "name": "Monkey 47",
    "brand": "Black Forest Distillers"
  },
  "tenant_id": "abc123"
}
```

### Webhook Signatur

Webhooks werden mit HMAC-SHA256 signiert:

```http
X-Webhook-Signature: sha256=abc123...
```

Verifizierung (Beispiel Node.js):
```javascript
const crypto = require('crypto');

function verifyWebhook(payload, signature, secret) {
  const expected = 'sha256=' + crypto
    .createHmac('sha256', secret)
    .update(payload)
    .digest('hex');
  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(expected)
  );
}
```

---

## SDK & Code-Beispiele

### cURL

```bash
# Alle Gins abrufen
curl -X GET "https://api.ginvault.de/api/v1/gins" \
  -H "Authorization: Bearer sk_your_api_key"

# Neuen Gin erstellen
curl -X POST "https://api.ginvault.de/api/v1/gins" \
  -H "Authorization: Bearer sk_your_api_key" \
  -H "Content-Type: application/json" \
  -d '{"name": "Hendricks", "brand": "William Grant", "abv": 41.4}'
```

### JavaScript/TypeScript

```typescript
const GINVAULT_API_KEY = 'sk_your_api_key';
const BASE_URL = 'https://api.ginvault.de/api/v1';

async function getGins() {
  const response = await fetch(`${BASE_URL}/gins`, {
    headers: {
      'Authorization': `Bearer ${GINVAULT_API_KEY}`,
      'Content-Type': 'application/json'
    }
  });
  return response.json();
}

async function createGin(gin: GinData) {
  const response = await fetch(`${BASE_URL}/gins`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${GINVAULT_API_KEY}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(gin)
  });
  return response.json();
}
```

### Python

```python
import requests

API_KEY = 'sk_your_api_key'
BASE_URL = 'https://api.ginvault.de/api/v1'

headers = {
    'Authorization': f'Bearer {API_KEY}',
    'Content-Type': 'application/json'
}

# Alle Gins abrufen
response = requests.get(f'{BASE_URL}/gins', headers=headers)
gins = response.json()

# Neuen Gin erstellen
new_gin = {
    'name': 'Tanqueray No. Ten',
    'brand': 'Tanqueray',
    'abv': 47.3,
    'price': 29.99
}
response = requests.post(f'{BASE_URL}/gins', json=new_gin, headers=headers)
```

### PHP

```php
<?php
$apiKey = 'sk_your_api_key';
$baseUrl = 'https://api.ginvault.de/api/v1';

// Alle Gins abrufen
$ch = curl_init("$baseUrl/gins");
curl_setopt($ch, CURLOPT_HTTPHEADER, [
    "Authorization: Bearer $apiKey",
    "Content-Type: application/json"
]);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$response = curl_exec($ch);
$gins = json_decode($response, true);
curl_close($ch);
```

---

## Best Practices

### Sicherheit

1. **API Key geheim halten** - Niemals in Client-Code oder Git-Repositories speichern
2. **HTTPS verwenden** - Alle Requests nur über verschlüsselte Verbindungen
3. **Key rotieren** - Regelmäßig neue API Keys generieren
4. **Minimale Berechtigungen** - Nur notwendige Endpoints nutzen

### Performance

1. **Pagination nutzen** - Große Listen in Seiten abrufen
2. **Caching implementieren** - Responses lokal zwischenspeichern
3. **Rate Limits beachten** - Requests über Zeit verteilen
4. **Batch-Operationen** - Mehrere Änderungen zusammenfassen

### Fehlerbehandlung

1. **Retries mit Backoff** - Bei 5xx Fehlern mit exponentieller Verzögerung wiederholen
2. **Rate Limit Header auswerten** - Vor Erreichen des Limits pausieren
3. **Timeouts setzen** - Requests nach 30 Sekunden abbrechen

---

## Support & Kontakt

Bei Fragen zur API-Integration:

- **E-Mail:** api-support@ginvault.de
- **Dokumentation:** https://docs.ginvault.de/api
- **Status-Seite:** https://status.ginvault.de

---

## Changelog

| Version | Datum | Änderungen |
|---------|-------|------------|
| 1.0.0 | 2026-01-16 | Initiale API-Version |

