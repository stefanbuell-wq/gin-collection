<?php
header('Content-Type: application/json');
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    exit(0);
}

require_once 'Database.php';

class GinAPI {
    private $db;

    public function __construct() {
        $this->db = Database::getInstance()->getConnection();
    }

    public function handleRequest() {
        $method = $_SERVER['REQUEST_METHOD'];
        $path = isset($_GET['action']) ? $_GET['action'] : '';

        try {
            switch ($path) {
                // Bestehende Endpoints
                case 'list': if ($method === 'GET') $this->listGins(); break;
                case 'get': if ($method === 'GET' && isset($_GET['id'])) $this->getGin($_GET['id']); break;
                case 'add': if ($method === 'POST') $this->addGin(); break;
                case 'update': if ($method === 'PUT' || $method === 'POST') $this->updateGin(); break;
                case 'delete': if ($method === 'DELETE' || $method === 'POST') $this->deleteGin(); break;
                case 'stats': if ($method === 'GET') $this->getStats(); break;
                case 'search': if ($method === 'GET') $this->searchGins(); break;
                case 'barcode': if ($method === 'GET' && isset($_GET['code'])) $this->lookupBarcode($_GET['code']); break;
                case 'upload': if ($method === 'POST') $this->uploadPhoto(); break;
                
                // Neue Endpoints
                case 'botanicals': if ($method === 'GET') $this->getBotanicals(); break;
                case 'gin-botanicals':
                    if ($method === 'GET' && isset($_GET['gin_id'])) $this->getGinBotanicals($_GET['gin_id']);
                    elseif ($method === 'POST') $this->saveGinBotanicals();
                    break;
                case 'cocktails': if ($method === 'GET') $this->getCocktails(); break;
                case 'cocktail': if ($method === 'GET' && isset($_GET['id'])) $this->getCocktail($_GET['id']); break;
                case 'gin-cocktails': if ($method === 'GET' && isset($_GET['gin_id'])) $this->getGinCocktails($_GET['gin_id']); break;
                case 'photos':
                    if ($method === 'GET' && isset($_GET['gin_id'])) $this->getGinPhotos($_GET['gin_id']);
                    elseif ($method === 'POST') $this->addGinPhoto();
                    elseif ($method === 'DELETE') $this->deleteGinPhoto();
                    break;
                case 'ai-suggestions': if ($method === 'GET' && isset($_GET['gin_id'])) $this->getAISuggestions($_GET['gin_id']); break;
                case 'export': if ($method === 'GET') $this->exportData(); break;
                case 'import': if ($method === 'POST') $this->importData(); break;
                
                default: $this->sendResponse(['error' => 'Invalid endpoint'], 404);
            }
        } catch (Exception $e) {
            $this->sendResponse(['error' => $e->getMessage()], 500);
        }
    }

    private function listGins() {
        $filter = $_GET['filter'] ?? 'all';
        $sort = $_GET['sort'] ?? 'name';
        $type = $_GET['type'] ?? null;
        
        $where = [];
        if ($filter === 'available') $where[] = 'is_finished = 0';
        elseif ($filter === 'finished') $where[] = 'is_finished = 1';
        if ($type) $where[] = "gin_type = " . $this->db->quote($type);
        
        $whereClause = count($where) > 0 ? 'WHERE ' . implode(' AND ', $where) : '';
        
        $orderBy = match($sort) {
            'rating' => 'rating DESC, name',
            'price' => 'price DESC',
            'date' => 'purchase_date DESC',
            'country' => 'country, name',
            'fill' => 'fill_level DESC, name',
            default => 'name'
        };

        $stmt = $this->db->query("
            SELECT id, name, brand, country, gin_type, abv, price, rating, photo_url, 
                   is_finished, fill_level, barcode
            FROM gins $whereClause ORDER BY $orderBy
        ");

        $gins = $stmt->fetchAll();
        $this->sendResponse(['gins' => $gins]);
    }

    private function getGin($id) {
        $stmt = $this->db->prepare("SELECT * FROM gins WHERE id = ?");
        $stmt->execute([$id]);
        $gin = $stmt->fetch();

        if ($gin) {
            $stmt = $this->db->prepare("
                SELECT b.id, b.name, b.category, gb.prominence
                FROM botanicals b
                JOIN gin_botanicals gb ON b.id = gb.botanical_id
                WHERE gb.gin_id = ?
            ");
            $stmt->execute([$id]);
            $gin['botanicals'] = $stmt->fetchAll();

            $stmt = $this->db->prepare("SELECT * FROM gin_photos WHERE gin_id = ? ORDER BY is_primary DESC, created_at");
            $stmt->execute([$id]);
            $gin['photos'] = $stmt->fetchAll();

            $stmt = $this->db->prepare("
                SELECT c.* FROM cocktails c
                JOIN gin_cocktails gc ON c.id = gc.cocktail_id
                WHERE gc.gin_id = ?
            ");
            $stmt->execute([$id]);
            $gin['cocktails'] = $stmt->fetchAll();

            $this->sendResponse(['gin' => $gin]);
        } else {
            $this->sendResponse(['error' => 'Gin not found'], 404);
        }
    }

    private function addGin() {
        $data = json_decode(file_get_contents('php://input'), true);

        $stmt = $this->db->prepare("
            INSERT INTO gins (name, brand, country, region, gin_type, abv, bottle_size, 
                            fill_level, price, current_market_value, purchase_date, purchase_location,
                            barcode, rating, nose_notes, palate_notes, finish_notes, general_notes,
                            description, photo_url, recommended_tonic, recommended_garnish)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ");

        $stmt->execute([
            $data['name'] ?? null,
            $data['brand'] ?? null,
            $data['country'] ?? null,
            $data['region'] ?? null,
            $data['gin_type'] ?? null,
            $data['abv'] ?? null,
            $data['bottle_size'] ?? 700,
            $data['fill_level'] ?? 100,
            $data['price'] ?? null,
            $data['current_market_value'] ?? null,
            $data['purchase_date'] ?? date('Y-m-d'),
            $data['purchase_location'] ?? null,
            $data['barcode'] ?? null,
            $data['rating'] ?? null,
            $data['nose_notes'] ?? null,
            $data['palate_notes'] ?? null,
            $data['finish_notes'] ?? null,
            $data['general_notes'] ?? null,
            $data['description'] ?? null,
            $data['photo_url'] ?? null,
            $data['recommended_tonic'] ?? null,
            $data['recommended_garnish'] ?? null
        ]);

        $id = $this->db->lastInsertId();
        
        if (isset($data['botanicals']) && is_array($data['botanicals'])) {
            $this->saveBotanicals($id, $data['botanicals']);
        }

        $this->sendResponse(['success' => true, 'id' => $id]);
    }

    private function updateGin() {
        $data = json_decode(file_get_contents('php://input'), true);
        $id = $data['id'] ?? null;

        if (!$id) {
            $this->sendResponse(['error' => 'ID required'], 400);
            return;
        }

        $stmt = $this->db->prepare("
            UPDATE gins SET 
                name = ?, brand = ?, country = ?, region = ?, gin_type = ?, abv = ?, 
                bottle_size = ?, fill_level = ?, price = ?, current_market_value = ?,
                purchase_date = ?, purchase_location = ?, barcode = ?,
                rating = ?, nose_notes = ?, palate_notes = ?, finish_notes = ?, 
                general_notes = ?, description = ?, photo_url = ?, is_finished = ?,
                recommended_tonic = ?, recommended_garnish = ?
            WHERE id = ?
        ");

        $stmt->execute([
            $data['name'] ?? null,
            $data['brand'] ?? null,
            $data['country'] ?? null,
            $data['region'] ?? null,
            $data['gin_type'] ?? null,
            $data['abv'] ?? null,
            $data['bottle_size'] ?? 700,
            $data['fill_level'] ?? 100,
            $data['price'] ?? null,
            $data['current_market_value'] ?? null,
            $data['purchase_date'] ?? null,
            $data['purchase_location'] ?? null,
            $data['barcode'] ?? null,
            $data['rating'] ?? null,
            $data['nose_notes'] ?? null,
            $data['palate_notes'] ?? null,
            $data['finish_notes'] ?? null,
            $data['general_notes'] ?? null,
            $data['description'] ?? null,
            $data['photo_url'] ?? null,
            $data['is_finished'] ?? 0,
            $data['recommended_tonic'] ?? null,
            $data['recommended_garnish'] ?? null,
            $id
        ]);

        if (isset($data['botanicals'])) {
            $this->saveBotanicals($id, $data['botanicals']);
        }

        $this->sendResponse(['success' => true]);
    }

    private function deleteGin() {
        $data = json_decode(file_get_contents('php://input'), true);
        $id = $data['id'] ?? $_GET['id'] ?? null;

        if (!$id) {
            $this->sendResponse(['error' => 'ID required'], 400);
            return;
        }

        $stmt = $this->db->prepare("DELETE FROM gins WHERE id = ?");
        $stmt->execute([$id]);

        $this->sendResponse(['success' => true]);
    }

    private function getStats() {
        $stats = [];

        $stmt = $this->db->query("SELECT COUNT(*) as total FROM gins");
        $stats['total'] = $stmt->fetch()['total'];

        $stmt = $this->db->query("SELECT COUNT(*) as available FROM gins WHERE is_finished = 0");
        $stats['available'] = $stmt->fetch()['available'];
        $stats['finished'] = $stats['total'] - $stats['available'];

        $stmt = $this->db->query("SELECT AVG(rating) as avg_rating FROM gins WHERE rating IS NOT NULL");
        $stats['avg_rating'] = round($stmt->fetch()['avg_rating'], 1);

        $stmt = $this->db->query("SELECT SUM(price) as total_value FROM gins WHERE price IS NOT NULL AND is_finished = 0");
        $stats['total_value'] = round($stmt->fetch()['total_value'], 2);

        $stmt = $this->db->query("SELECT SUM(current_market_value) as market_value FROM gins WHERE current_market_value IS NOT NULL AND is_finished = 0");
        $stats['current_market_value'] = round($stmt->fetch()['market_value'], 2);

        $stmt = $this->db->query("SELECT gin_type, COUNT(*) as count FROM gins WHERE gin_type IS NOT NULL GROUP BY gin_type ORDER BY count DESC");
        $stats['gin_types'] = $stmt->fetchAll();

        $stmt = $this->db->query("SELECT country, COUNT(*) as count FROM gins WHERE country IS NOT NULL GROUP BY country ORDER BY count DESC");
        $stats['countries'] = $stmt->fetchAll();

        $stmt = $this->db->query("SELECT id, name, brand, rating, photo_url FROM gins WHERE rating IS NOT NULL ORDER BY rating DESC, name LIMIT 5");
        $stats['top_rated'] = $stmt->fetchAll();

        $stmt = $this->db->query("
            SELECT b.name, COUNT(*) as count
            FROM botanicals b
            JOIN gin_botanicals gb ON b.id = gb.botanical_id
            GROUP BY b.id
            ORDER BY count DESC
            LIMIT 10
        ");
        $stats['top_botanicals'] = $stmt->fetchAll();

        $stmt = $this->db->query("
            SELECT 
                CASE 
                    WHEN fill_level = 100 THEN 'Voll'
                    WHEN fill_level >= 75 THEN '75-99%'
                    WHEN fill_level >= 50 THEN '50-74%'
                    WHEN fill_level >= 25 THEN '25-49%'
                    WHEN fill_level > 0 THEN '1-24%'
                    ELSE 'Leer'
                END as level,
                COUNT(*) as count
            FROM gins
            WHERE is_finished = 0
            GROUP BY level
        ");
        $stats['fill_levels'] = $stmt->fetchAll();

        $this->sendResponse(['stats' => $stats]);
    }

    private function searchGins() {
        $query = $_GET['q'] ?? '';
        
        if (empty($query)) {
            $this->listGins();
            return;
        }

        $stmt = $this->db->prepare("
            SELECT id, name, brand, country, gin_type, abv, price, rating, photo_url, is_finished, fill_level
            FROM gins 
            WHERE name LIKE ? OR brand LIKE ? OR country LIKE ? OR gin_type LIKE ? 
                  OR nose_notes LIKE ? OR palate_notes LIKE ? OR finish_notes LIKE ? OR general_notes LIKE ?
            ORDER BY name
        ");

        $searchTerm = "%$query%";
        $stmt->execute(array_fill(0, 8, $searchTerm));
        $gins = $stmt->fetchAll();

        $this->sendResponse(['gins' => $gins]);
    }

    private function lookupBarcode($barcode) {
        $stmt = $this->db->prepare("SELECT * FROM gins WHERE barcode = ?");
        $stmt->execute([$barcode]);
        $existing = $stmt->fetch();

        if ($existing) {
            $this->sendResponse(['exists' => true, 'gin' => $existing]);
            return;
        }

        $url = "https://world.openfoodfacts.org/api/v0/product/{$barcode}.json";
        $response = @file_get_contents($url);

        if ($response) {
            $data = json_decode($response, true);
            
            if ($data['status'] === 1) {
                $product = $data['product'];
                $ginData = [
                    'exists' => false,
                    'name' => $product['product_name'] ?? '',
                    'brand' => $product['brands'] ?? '',
                    'country' => $product['countries'] ?? '',
                    'abv' => $product['alcohol_volume'] ?? null,
                    'photo_url' => $product['image_url'] ?? null,
                    'barcode' => $barcode
                ];
                $this->sendResponse($ginData);
                return;
            }
        }

        $this->sendResponse(['exists' => false, 'found' => false, 'barcode' => $barcode]);
    }

    private function uploadPhoto() {
        if (!isset($_FILES['photo'])) {
            $this->sendResponse(['error' => 'No photo uploaded'], 400);
            return;
        }

        $file = $_FILES['photo'];
        $uploadDir = __DIR__ . '/../uploads/';
        
        if (!file_exists($uploadDir)) {
            mkdir($uploadDir, 0755, true);
        }

        $extension = pathinfo($file['name'], PATHINFO_EXTENSION);
        $filename = uniqid('gin_') . '.' . $extension;
        $filepath = $uploadDir . $filename;

        if (move_uploaded_file($file['tmp_name'], $filepath)) {
            $url = '/uploads/' . $filename;
            $this->sendResponse(['success' => true, 'url' => $url]);
        } else {
            $this->sendResponse(['error' => 'Upload failed'], 500);
        }
    }

    // Neue Funktionen folgen im nächsten Teil...
    private function getBotanicals() {
        $stmt = $this->db->query("SELECT * FROM botanicals ORDER BY category, name");
        $botanicals = $stmt->fetchAll();
        
        $grouped = [];
        foreach ($botanicals as $botanical) {
            $category = $botanical['category'] ?? 'Sonstige';
            if (!isset($grouped[$category])) $grouped[$category] = [];
            $grouped[$category][] = $botanical;
        }
        
        $this->sendResponse(['botanicals' => $botanicals, 'grouped' => $grouped]);
    }

    private function getGinBotanicals($ginId) {
        $stmt = $this->db->prepare("
            SELECT b.*, gb.prominence
            FROM botanicals b
            JOIN gin_botanicals gb ON b.id = gb.botanical_id
            WHERE gb.gin_id = ?
            ORDER BY b.name
        ");
        $stmt->execute([$ginId]);
        $this->sendResponse(['botanicals' => $stmt->fetchAll()]);
    }

    private function saveGinBotanicals() {
        $data = json_decode(file_get_contents('php://input'), true);
        $ginId = $data['gin_id'];
        $botanicals = $data['botanicals'] ?? [];

        $stmt = $this->db->prepare("DELETE FROM gin_botanicals WHERE gin_id = ?");
        $stmt->execute([$ginId]);

        $this->saveBotanicals($ginId, $botanicals);
        $this->sendResponse(['success' => true]);
    }

    private function saveBotanicals($ginId, $botanicals) {
        if (empty($botanicals)) return;

        $stmt = $this->db->prepare("INSERT INTO gin_botanicals (gin_id, botanical_id, prominence) VALUES (?, ?, ?)");

        foreach ($botanicals as $botanical) {
            $botanicalId = is_array($botanical) ? $botanical['id'] : $botanical;
            $prominence = is_array($botanical) ? ($botanical['prominence'] ?? 'notable') : 'notable';
            $stmt->execute([$ginId, $botanicalId, $prominence]);
        }
    }

    private function getCocktails() {
        $stmt = $this->db->query("SELECT * FROM cocktails ORDER BY difficulty, name");
        $this->sendResponse(['cocktails' => $stmt->fetchAll()]);
    }

    private function getCocktail($id) {
        $stmt = $this->db->prepare("SELECT * FROM cocktails WHERE id = ?");
        $stmt->execute([$id]);
        $cocktail = $stmt->fetch();

        if ($cocktail) {
            $stmt = $this->db->prepare("SELECT * FROM cocktail_ingredients WHERE cocktail_id = ? ORDER BY id");
            $stmt->execute([$id]);
            $cocktail['ingredients'] = $stmt->fetchAll();
            $this->sendResponse(['cocktail' => $cocktail]);
        } else {
            $this->sendResponse(['error' => 'Cocktail not found'], 404);
        }
    }

    private function getGinCocktails($ginId) {
        $stmt = $this->db->prepare("
            SELECT c.* FROM cocktails c
            JOIN gin_cocktails gc ON c.id = gc.cocktail_id
            WHERE gc.gin_id = ?
        ");
        $stmt->execute([$ginId]);
        $this->sendResponse(['cocktails' => $stmt->fetchAll()]);
    }

    private function getGinPhotos($ginId) {
        $stmt = $this->db->prepare("SELECT * FROM gin_photos WHERE gin_id = ? ORDER BY is_primary DESC, created_at DESC");
        $stmt->execute([$ginId]);
        $this->sendResponse(['photos' => $stmt->fetchAll()]);
    }

    private function addGinPhoto() {
        $data = json_decode(file_get_contents('php://input'), true);
        
        $stmt = $this->db->prepare("
            INSERT INTO gin_photos (gin_id, photo_url, photo_type, caption, is_primary)
            VALUES (?, ?, ?, ?, ?)
        ");
        
        $stmt->execute([
            $data['gin_id'],
            $data['photo_url'],
            $data['photo_type'] ?? 'bottle',
            $data['caption'] ?? null,
            $data['is_primary'] ?? 0
        ]);
        
        $this->sendResponse(['success' => true, 'id' => $this->db->lastInsertId()]);
    }

    private function deleteGinPhoto() {
        $data = json_decode(file_get_contents('php://input'), true);
        $stmt = $this->db->prepare("DELETE FROM gin_photos WHERE id = ?");
        $stmt->execute([$data['id']]);
        $this->sendResponse(['success' => true]);
    }

    private function getAISuggestions($ginId) {
        $stmt = $this->db->prepare("SELECT * FROM gins WHERE id = ?");
        $stmt->execute([$ginId]);
        $gin = $stmt->fetch();

        if (!$gin) {
            $this->sendResponse(['error' => 'Gin not found'], 404);
            return;
        }

        $suggestions = [];

        if ($gin['country']) {
            $stmt = $this->db->prepare("
                SELECT id, name, brand, rating, photo_url
                FROM gins WHERE country = ? AND id != ? ORDER BY rating DESC LIMIT 3
            ");
            $stmt->execute([$gin['country'], $ginId]);
            $suggestions['same_country'] = $stmt->fetchAll();
        }

        if ($gin['rating']) {
            $stmt = $this->db->prepare("
                SELECT id, name, brand, rating, photo_url
                FROM gins WHERE rating = ? AND id != ? ORDER BY RANDOM() LIMIT 3
            ");
            $stmt->execute([$gin['rating'], $ginId]);
            $suggestions['similar_rating'] = $stmt->fetchAll();
        }

        $stmt = $this->db->prepare("
            SELECT DISTINCT g.id, g.name, g.brand, g.rating, g.photo_url, COUNT(*) as shared
            FROM gins g
            JOIN gin_botanicals gb1 ON g.id = gb1.gin_id
            JOIN gin_botanicals gb2 ON gb1.botanical_id = gb2.botanical_id
            WHERE gb2.gin_id = ? AND g.id != ?
            GROUP BY g.id ORDER BY shared DESC, g.rating DESC LIMIT 3
        ");
        $stmt->execute([$ginId, $ginId]);
        $suggestions['shared_botanicals'] = $stmt->fetchAll();

        $this->sendResponse(['suggestions' => $suggestions]);
    }

    private function exportData() {
        $format = $_GET['format'] ?? 'json';
        $stmt = $this->db->query("SELECT * FROM gins ORDER BY name");
        $gins = $stmt->fetchAll();

        foreach ($gins as &$gin) {
            $stmt = $this->db->prepare("
                SELECT b.name, gb.prominence
                FROM botanicals b
                JOIN gin_botanicals gb ON b.id = gb.botanical_id
                WHERE gb.gin_id = ?
            ");
            $stmt->execute([$gin['id']]);
            $gin['botanicals'] = $stmt->fetchAll();
        }

        if ($format === 'csv') {
            $this->exportCSV($gins);
        } else {
            $this->sendResponse([
                'export_date' => date('Y-m-d H:i:s'),
                'total_gins' => count($gins),
                'gins' => $gins
            ]);
        }
    }

    private function exportCSV($gins) {
        header('Content-Type: text/csv');
        header('Content-Disposition: attachment; filename="gin-collection-' . date('Y-m-d') . '.csv"');
        
        $output = fopen('php://output', 'w');
        fputcsv($output, ['Name', 'Marke', 'Land', 'Typ', 'ABV', 'Preis', 'Bewertung', 'Füllstand', 'Kaufdatum', 'Händler', 'Barcode']);
        
        foreach ($gins as $gin) {
            fputcsv($output, [
                $gin['name'], $gin['brand'], $gin['country'], $gin['gin_type'], $gin['abv'],
                $gin['price'], $gin['rating'], $gin['fill_level'], $gin['purchase_date'],
                $gin['purchase_location'], $gin['barcode']
            ]);
        }
        
        fclose($output);
        exit;
    }

    private function importData() {
        $data = json_decode(file_get_contents('php://input'), true);
        
        if (!isset($data['gins']) || !is_array($data['gins'])) {
            $this->sendResponse(['error' => 'Invalid import data'], 400);
            return;
        }

        $imported = 0;
        $errors = [];

        foreach ($data['gins'] as $ginData) {
            try {
                if (isset($ginData['barcode'])) {
                    $stmt = $this->db->prepare("SELECT id FROM gins WHERE barcode = ?");
                    $stmt->execute([$ginData['barcode']]);
                    if ($stmt->fetch()) {
                        $errors[] = "Gin mit Barcode {$ginData['barcode']} existiert bereits";
                        continue;
                    }
                }

                $stmt = $this->db->prepare("
                    INSERT INTO gins (name, brand, country, gin_type, abv, price, rating)
                    VALUES (?, ?, ?, ?, ?, ?, ?)
                ");
                $stmt->execute([
                    $ginData['name'], $ginData['brand'] ?? null, $ginData['country'] ?? null,
                    $ginData['gin_type'] ?? null, $ginData['abv'] ?? null,
                    $ginData['price'] ?? null, $ginData['rating'] ?? null
                ]);
                
                $imported++;
            } catch (Exception $e) {
                $errors[] = "Fehler bei {$ginData['name']}: " . $e->getMessage();
            }
        }

        $this->sendResponse(['success' => true, 'imported' => $imported, 'errors' => $errors]);
    }

    private function sendResponse($data, $status = 200) {
        http_response_code($status);
        echo json_encode($data);
        exit;
    }
}

$api = new GinAPI();
$api->handleRequest();
