<?php
// Configure secure session settings
ini_set('session.cookie_httponly', 1);
ini_set('session.cookie_samesite', 'Strict');
ini_set('session.use_strict_mode', 1);

// Start session early for authentication
session_start();

header('Content-Type: application/json');
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    exit(0);
}

require_once 'Database.php';
require_once 'Auth.php';

class GinAPI
{
    private $db;

    public function __construct()
    {
        $this->db = Database::getInstance()->getConnection();
    }

    public function handleRequest()
    {
        $method = $_SERVER['REQUEST_METHOD'];
        $path = isset($_GET['action']) ? $_GET['action'] : '';

        try {
            switch ($path) {
                // Authentication endpoints (no auth required)
                case 'register':
                    if ($method === 'POST') {
                        $this->register();
                    }
                    break;

                case 'login':
                    if ($method === 'POST') {
                        $this->login();
                    }
                    break;

                case 'logout':
                    if ($method === 'POST') {
                        $this->logout();
                    }
                    break;

                case 'me':
                    if ($method === 'GET') {
                        $this->getCurrentUser();
                    }
                    break;

                // Protected endpoints (auth required)
                case 'list':
                    if ($method === 'GET') {
                        $this->listGins();
                    }
                    break;

                case 'get':
                    if ($method === 'GET' && isset($_GET['id'])) {
                        $this->getGin($_GET['id']);
                    }
                    break;

                case 'add':
                    if ($method === 'POST') {
                        $this->addGin();
                    }
                    break;

                case 'update':
                    if ($method === 'PUT' || $method === 'POST') {
                        $this->updateGin();
                    }
                    break;

                case 'delete':
                    if ($method === 'DELETE' || $method === 'POST') {
                        $this->deleteGin();
                    }
                    break;

                case 'stats':
                    if ($method === 'GET') {
                        $this->getStats();
                    }
                    break;

                case 'search':
                    if ($method === 'GET') {
                        $this->searchGins();
                    }
                    break;

                case 'barcode':
                    if ($method === 'GET' && isset($_GET['code'])) {
                        $this->lookupBarcode($_GET['code']);
                    }
                    break;

                case 'upload':
                    if ($method === 'POST') {
                        $this->uploadPhoto();
                    }
                    break;

                // Admin endpoints (admin auth required)
                case 'admin-users':
                    if ($method === 'GET') {
                        $this->adminListUsers();
                    }
                    break;

                case 'admin-user-stats':
                    if ($method === 'GET' && isset($_GET['user_id'])) {
                        $this->adminGetUserStats($_GET['user_id']);
                    }
                    break;

                case 'admin-update-user':
                    if ($method === 'PUT' || $method === 'POST') {
                        $this->adminUpdateUser();
                    }
                    break;

                case 'admin-delete-user':
                    if ($method === 'DELETE' || $method === 'POST') {
                        $this->adminDeleteUser();
                    }
                    break;

                case 'admin-dashboard-stats':
                    if ($method === 'GET') {
                        $this->adminGetDashboardStats();
                    }
                    break;

                default:
                    $this->sendResponse(['error' => 'Invalid endpoint'], 404);
            }
        } catch (Exception $e) {
            $this->sendResponse(['error' => $e->getMessage()], 500);
        }
    }

    // Authentication methods
    private function register()
    {
        $data = json_decode(file_get_contents('php://input'), true);

        try {
            $auth = Auth::getInstance();
            $user = $auth->register(
                $data['username'] ?? null,
                $data['email'] ?? null,
                $data['password'] ?? null,
                $data['full_name'] ?? null
            );

            $this->sendResponse(['success' => true, 'user' => $user]);
        } catch (Exception $e) {
            $this->sendResponse(['error' => $e->getMessage()], 400);
        }
    }

    private function login()
    {
        $data = json_decode(file_get_contents('php://input'), true);

        try {
            $auth = Auth::getInstance();
            $user = $auth->login(
                $data['username'] ?? $data['email'] ?? null,
                $data['password'] ?? null
            );

            $this->sendResponse(['success' => true, 'user' => $user]);
        } catch (Exception $e) {
            $this->sendResponse(['error' => $e->getMessage()], 401);
        }
    }

    private function logout()
    {
        try {
            $auth = Auth::getInstance();
            $auth->logout();

            $this->sendResponse(['success' => true]);
        } catch (Exception $e) {
            $this->sendResponse(['error' => $e->getMessage()], 500);
        }
    }

    private function getCurrentUser()
    {
        try {
            $auth = Auth::getInstance();
            $user = $auth->getCurrentUser();

            if ($user) {
                $this->sendResponse(['user' => $user]);
            } else {
                $this->sendResponse(['error' => 'Not authenticated'], 401);
            }
        } catch (Exception $e) {
            $this->sendResponse(['error' => $e->getMessage()], 500);
        }
    }

    // Gin management methods (protected)
    private function listGins()
    {
        $auth = Auth::getInstance();
        $userId = $auth->requireAuth();

        $filter = isset($_GET['filter']) ? $_GET['filter'] : 'all';
        $sort = isset($_GET['sort']) ? $_GET['sort'] : 'name';

        $where = 'WHERE user_id = ?';
        $params = [$userId];

        if ($filter === 'available') {
            $where .= ' AND is_finished = 0';
        } elseif ($filter === 'finished') {
            $where .= ' AND is_finished = 1';
        }

        $orderBy = match ($sort) {
            'rating' => 'rating DESC, name',
            'price' => 'price DESC',
            'date' => 'purchase_date DESC',
            'country' => 'country, name',
            default => 'name'
        };

        $stmt = $this->db->prepare("
            SELECT id, name, brand, country, abv, price, rating, photo_url, is_finished, barcode
            FROM gins 
            $where
            ORDER BY $orderBy
        ");
        $stmt->execute($params);

        $gins = $stmt->fetchAll();
        $this->sendResponse(['gins' => $gins]);
    }

    private function getGin($id)
    {
        $auth = Auth::getInstance();
        $userId = $auth->requireAuth();

        $stmt = $this->db->prepare("SELECT * FROM gins WHERE id = ? AND user_id = ?");
        $stmt->execute([$id, $userId]);
        $gin = $stmt->fetch();

        if ($gin) {
            // Get botanicals if any
            $stmt = $this->db->prepare("
                SELECT b.name 
                FROM botanicals b
                JOIN gin_botanicals gb ON b.id = gb.botanical_id
                WHERE gb.gin_id = ?
            ");
            $stmt->execute([$id]);
            $gin['botanicals'] = $stmt->fetchAll(PDO::FETCH_COLUMN);

            $this->sendResponse(['gin' => $gin]);
        } else {
            $this->sendResponse(['error' => 'Gin not found'], 404);
        }
    }

    private function addGin()
    {
        $auth = Auth::getInstance();
        $userId = $auth->requireAuth();

        $data = json_decode(file_get_contents('php://input'), true);

        $stmt = $this->db->prepare("
            INSERT INTO gins (user_id, name, brand, country, region, abv, bottle_size, price, 
                            purchase_date, barcode, rating, tasting_notes, description, photo_url)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ");

        $stmt->execute([
            $userId,
            $data['name'] ?? null,
            $data['brand'] ?? null,
            $data['country'] ?? null,
            $data['region'] ?? null,
            $data['abv'] ?? null,
            $data['bottle_size'] ?? 700,
            $data['price'] ?? null,
            $data['purchase_date'] ?? date('Y-m-d'),
            $data['barcode'] ?? null,
            $data['rating'] ?? null,
            $data['tasting_notes'] ?? null,
            $data['description'] ?? null,
            $data['photo_url'] ?? null
        ]);

        $id = $this->db->lastInsertId();
        $this->sendResponse(['success' => true, 'id' => $id]);
    }

    private function updateGin()
    {
        $auth = Auth::getInstance();
        $userId = $auth->requireAuth();

        $data = json_decode(file_get_contents('php://input'), true);
        $id = $data['id'] ?? null;

        if (!$id) {
            $this->sendResponse(['error' => 'ID required'], 400);
            return;
        }

        // Verify ownership
        $stmt = $this->db->prepare("SELECT id FROM gins WHERE id = ? AND user_id = ?");
        $stmt->execute([$id, $userId]);
        if (!$stmt->fetch()) {
            $this->sendResponse(['error' => 'Gin not found or access denied'], 403);
            return;
        }

        $stmt = $this->db->prepare("
            UPDATE gins SET 
                name = ?, brand = ?, country = ?, region = ?, abv = ?, 
                bottle_size = ?, price = ?, purchase_date = ?, barcode = ?,
                rating = ?, tasting_notes = ?, description = ?, photo_url = ?, is_finished = ?
            WHERE id = ? AND user_id = ?
        ");

        $stmt->execute([
            $data['name'] ?? null,
            $data['brand'] ?? null,
            $data['country'] ?? null,
            $data['region'] ?? null,
            $data['abv'] ?? null,
            $data['bottle_size'] ?? 700,
            $data['price'] ?? null,
            $data['purchase_date'] ?? null,
            $data['barcode'] ?? null,
            $data['rating'] ?? null,
            $data['tasting_notes'] ?? null,
            $data['description'] ?? null,
            $data['photo_url'] ?? null,
            $data['is_finished'] ?? 0,
            $id,
            $userId
        ]);

        $this->sendResponse(['success' => true]);
    }

    private function deleteGin()
    {
        $auth = Auth::getInstance();
        $userId = $auth->requireAuth();

        $data = json_decode(file_get_contents('php://input'), true);
        $id = $data['id'] ?? $_GET['id'] ?? null;

        if (!$id) {
            $this->sendResponse(['error' => 'ID required'], 400);
            return;
        }

        $stmt = $this->db->prepare("DELETE FROM gins WHERE id = ? AND user_id = ?");
        $stmt->execute([$id, $userId]);

        if ($stmt->rowCount() === 0) {
            $this->sendResponse(['error' => 'Gin not found or access denied'], 403);
            return;
        }

        $this->sendResponse(['success' => true]);
    }

    private function getStats()
    {
        $auth = Auth::getInstance();
        $userId = $auth->requireAuth();

        $stats = [];

        // Total count
        $stmt = $this->db->prepare("SELECT COUNT(*) as total FROM gins WHERE user_id = ?");
        $stmt->execute([$userId]);
        $stats['total'] = $stmt->fetch()['total'];

        // Available vs finished
        $stmt = $this->db->prepare("SELECT COUNT(*) as available FROM gins WHERE user_id = ? AND is_finished = 0");
        $stmt->execute([$userId]);
        $stats['available'] = $stmt->fetch()['available'];
        $stats['finished'] = $stats['total'] - $stats['available'];

        // Average rating
        $stmt = $this->db->prepare("SELECT AVG(rating) as avg_rating FROM gins WHERE user_id = ? AND rating IS NOT NULL");
        $stmt->execute([$userId]);
        $stats['avg_rating'] = round($stmt->fetch()['avg_rating'], 1);

        // Total value
        $stmt = $this->db->prepare("SELECT SUM(price) as total_value FROM gins WHERE user_id = ? AND price IS NOT NULL");
        $stmt->execute([$userId]);
        $stats['total_value'] = round($stmt->fetch()['total_value'], 2);

        // Countries
        $stmt = $this->db->prepare("
            SELECT country, COUNT(*) as count 
            FROM gins 
            WHERE user_id = ? AND country IS NOT NULL 
            GROUP BY country 
            ORDER BY count DESC
        ");
        $stmt->execute([$userId]);
        $stats['countries'] = $stmt->fetchAll();

        // Top rated
        $stmt = $this->db->prepare("
            SELECT name, brand, rating, photo_url
            FROM gins 
            WHERE user_id = ? AND rating IS NOT NULL 
            ORDER BY rating DESC, name
            LIMIT 5
        ");
        $stmt->execute([$userId]);
        $stats['top_rated'] = $stmt->fetchAll();

        $this->sendResponse(['stats' => $stats]);
    }

    private function searchGins()
    {
        $auth = Auth::getInstance();
        $userId = $auth->requireAuth();

        $query = $_GET['q'] ?? '';

        if (empty($query)) {
            $this->listGins();
            return;
        }

        $stmt = $this->db->prepare("
            SELECT id, name, brand, country, abv, price, rating, photo_url, is_finished
            FROM gins 
            WHERE user_id = ? AND (name LIKE ? OR brand LIKE ? OR country LIKE ? OR tasting_notes LIKE ?)
            ORDER BY name
        ");

        $searchTerm = "%$query%";
        $stmt->execute([$userId, $searchTerm, $searchTerm, $searchTerm, $searchTerm]);
        $gins = $stmt->fetchAll();

        $this->sendResponse(['gins' => $gins]);
    }

    private function lookupBarcode($barcode)
    {
        // First check if we already have this gin
        $stmt = $this->db->prepare("SELECT * FROM gins WHERE barcode = ?");
        $stmt->execute([$barcode]);
        $existing = $stmt->fetch();

        if ($existing) {
            $this->sendResponse(['exists' => true, 'gin' => $existing]);
            return;
        }

        // Try to lookup via Open Food Facts API
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

    private function uploadPhoto()
    {
        if (!isset($_FILES['photo'])) {
            $this->sendResponse(['error' => 'No photo uploaded'], 400);
            return;
        }

        $file = $_FILES['photo'];
        $uploadDir = __DIR__ . '/../uploads/';

        if (!file_exists($uploadDir)) {
            mkdir($uploadDir, 0755, true);
        }

        // Generate unique filename
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

    // Admin methods (admin auth required)
    private function adminListUsers()
    {
        $auth = Auth::getInstance();
        $auth->requireAdmin();

        $stmt = $this->db->prepare("
            SELECT u.id, u.username, u.email, u.full_name, u.is_admin, u.created_at,
                   COUNT(g.id) as gin_count
            FROM users u
            LEFT JOIN gins g ON u.id = g.user_id
            GROUP BY u.id
            ORDER BY u.created_at DESC
        ");
        $stmt->execute();

        $users = $stmt->fetchAll();
        $this->sendResponse(['users' => $users]);
    }

    private function adminGetUserStats($userId)
    {
        $auth = Auth::getInstance();
        $auth->requireAdmin();

        // Get user info
        $stmt = $this->db->prepare("SELECT id, username, email, full_name, created_at FROM users WHERE id = ?");
        $stmt->execute([$userId]);
        $user = $stmt->fetch();

        if (!$user) {
            $this->sendResponse(['error' => 'User not found'], 404);
            return;
        }

        // Get gin count
        $stmt = $this->db->prepare("SELECT COUNT(*) as total FROM gins WHERE user_id = ?");
        $stmt->execute([$userId]);
        $user['gin_count'] = $stmt->fetch()['total'];

        // Get latest gins
        $stmt = $this->db->prepare("
            SELECT id, name, brand, created_at 
            FROM gins 
            WHERE user_id = ? 
            ORDER BY created_at DESC 
            LIMIT 5
        ");
        $stmt->execute([$userId]);
        $user['latest_gins'] = $stmt->fetchAll();

        $this->sendResponse(['user' => $user]);
    }

    private function adminUpdateUser()
    {
        $auth = Auth::getInstance();
        $adminId = $auth->requireAdmin();

        $data = json_decode(file_get_contents('php://input'), true);
        $userId = $data['user_id'] ?? null;

        if (!$userId) {
            $this->sendResponse(['error' => 'User ID required'], 400);
            return;
        }

        // Prevent admin from modifying their own admin status
        if ($userId == $adminId && isset($data['is_admin'])) {
            $this->sendResponse(['error' => 'Cannot modify your own admin status'], 403);
            return;
        }

        $updates = [];
        $params = [];

        if (isset($data['email'])) {
            $updates[] = "email = ?";
            $params[] = $data['email'];
        }

        if (isset($data['full_name'])) {
            $updates[] = "full_name = ?";
            $params[] = $data['full_name'];
        }

        if (isset($data['is_admin']) && $userId != $adminId) {
            $updates[] = "is_admin = ?";
            $params[] = $data['is_admin'] ? 1 : 0;
        }

        if (empty($updates)) {
            $this->sendResponse(['error' => 'No updates provided'], 400);
            return;
        }

        $params[] = $userId;
        $sql = "UPDATE users SET " . implode(', ', $updates) . " WHERE id = ?";

        $stmt = $this->db->prepare($sql);
        $stmt->execute($params);

        $this->sendResponse(['success' => true]);
    }

    private function adminDeleteUser()
    {
        $auth = Auth::getInstance();
        $adminId = $auth->requireAdmin();

        $data = json_decode(file_get_contents('php://input'), true);
        $userId = $data['user_id'] ?? $_GET['user_id'] ?? null;

        if (!$userId) {
            $this->sendResponse(['error' => 'User ID required'], 400);
            return;
        }

        // Prevent admin from deleting themselves
        if ($userId == $adminId) {
            $this->sendResponse(['error' => 'Cannot delete your own account'], 403);
            return;
        }

        // Delete user (CASCADE will delete all gins)
        $stmt = $this->db->prepare("DELETE FROM users WHERE id = ?");
        $stmt->execute([$userId]);

        if ($stmt->rowCount() === 0) {
            $this->sendResponse(['error' => 'User not found'], 404);
            return;
        }

        $this->sendResponse(['success' => true]);
    }

    private function adminGetDashboardStats()
    {
        $auth = Auth::getInstance();
        $auth->requireAdmin();

        $stats = [];

        // Total users
        $stmt = $this->db->query("SELECT COUNT(*) as total FROM users");
        $stats['total_users'] = $stmt->fetch()['total'];

        // Total gins (across all users)
        $stmt = $this->db->query("SELECT COUNT(*) as total FROM gins");
        $stats['total_gins'] = $stmt->fetch()['total'];

        // Recent registrations (last 5)
        $stmt = $this->db->query("
            SELECT id, username, email, full_name, created_at 
            FROM users 
            ORDER BY created_at DESC 
            LIMIT 5
        ");
        $stats['recent_users'] = $stmt->fetchAll();

        // Most active users (by gin count)
        $stmt = $this->db->query("
            SELECT u.id, u.username, u.email, COUNT(g.id) as gin_count
            FROM users u
            LEFT JOIN gins g ON u.id = g.user_id
            GROUP BY u.id
            ORDER BY gin_count DESC
            LIMIT 5
        ");
        $stats['most_active_users'] = $stmt->fetchAll();

        $this->sendResponse(['stats' => $stats]);
    }

    private function sendResponse($data, $status = 200)
    {
        http_response_code($status);
        echo json_encode($data);
        exit;
    }
}

// Initialize and handle request
$api = new GinAPI();
$api->handleRequest();
