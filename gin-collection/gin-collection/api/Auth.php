<?php
// Authentication and User Management Class

class Auth
{
    private $db;
    private static $instance = null;

    private function __construct()
    {
        $this->db = Database::getInstance()->getConnection();

        // Start session if not already started
        if (session_status() === PHP_SESSION_NONE) {
            session_start();
        }
    }

    public static function getInstance()
    {
        if (self::$instance === null) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    /**
     * Register a new user
     */
    public function register($username, $email, $password, $fullName = null)
    {
        // Validate input
        if (empty($username) || empty($email) || empty($password)) {
            throw new Exception('Username, email, and password are required');
        }

        // Validate email format
        if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
            throw new Exception('Invalid email format');
        }

        // Validate password strength
        if (strlen($password) < 8) {
            throw new Exception('Password must be at least 8 characters long');
        }

        // Check if username already exists
        $stmt = $this->db->prepare("SELECT id FROM users WHERE username = ?");
        $stmt->execute([$username]);
        if ($stmt->fetch()) {
            throw new Exception('Username already exists');
        }

        // Check if email already exists
        $stmt = $this->db->prepare("SELECT id FROM users WHERE email = ?");
        $stmt->execute([$email]);
        if ($stmt->fetch()) {
            throw new Exception('Email already exists');
        }

        // Hash password with bcrypt (cost factor 12)
        $passwordHash = password_hash($password, PASSWORD_BCRYPT, ['cost' => 12]);

        // Insert user
        $stmt = $this->db->prepare("
            INSERT INTO users (username, email, password_hash, full_name)
            VALUES (?, ?, ?, ?)
        ");

        $stmt->execute([$username, $email, $passwordHash, $fullName]);

        $userId = $this->db->lastInsertId();

        return [
            'id' => $userId,
            'username' => $username,
            'email' => $email,
            'full_name' => $fullName
        ];
    }

    /**
     * Login user
     */
    public function login($usernameOrEmail, $password)
    {
        // Validate input
        if (empty($usernameOrEmail) || empty($password)) {
            throw new Exception('Username/email and password are required');
        }

        // Find user by username or email
        $stmt = $this->db->prepare("
            SELECT id, username, email, password_hash, full_name, is_admin 
            FROM users 
            WHERE username = ? OR email = ?
        ");
        $stmt->execute([$usernameOrEmail, $usernameOrEmail]);
        $user = $stmt->fetch();

        if (!$user) {
            throw new Exception('Invalid credentials');
        }

        // Verify password
        if (!password_verify($password, $user['password_hash'])) {
            throw new Exception('Invalid credentials');
        }

        // Create session
        $_SESSION['user_id'] = $user['id'];
        $_SESSION['username'] = $user['username'];
        $_SESSION['email'] = $user['email'];
        $_SESSION['full_name'] = $user['full_name'];
        $_SESSION['is_admin'] = $user['is_admin'];
        $_SESSION['login_time'] = time();

        // Regenerate session ID for security
        session_regenerate_id(true);

        return [
            'id' => $user['id'],
            'username' => $user['username'],
            'email' => $user['email'],
            'full_name' => $user['full_name'],
            'is_admin' => $user['is_admin']
        ];
    }

    /**
     * Logout user
     */
    public function logout()
    {
        // Clear session data
        $_SESSION = [];

        // Destroy session cookie
        if (isset($_COOKIE[session_name()])) {
            setcookie(session_name(), '', time() - 3600, '/');
        }

        // Destroy session
        session_destroy();

        return true;
    }

    /**
     * Check if user is authenticated
     */
    public function isAuthenticated()
    {
        return isset($_SESSION['user_id']) && !empty($_SESSION['user_id']);
    }

    /**
     * Get current authenticated user
     */
    public function getCurrentUser()
    {
        if (!$this->isAuthenticated()) {
            return null;
        }

        return [
            'id' => $_SESSION['user_id'],
            'username' => $_SESSION['username'],
            'email' => $_SESSION['email'],
            'full_name' => $_SESSION['full_name'],
            'is_admin' => $_SESSION['is_admin'] ?? 0
        ];
    }

    /**
     * Get current user ID
     */
    public function getCurrentUserId()
    {
        if (!$this->isAuthenticated()) {
            return null;
        }

        return $_SESSION['user_id'];
    }

    /**
     * Require authentication (middleware)
     * Throws exception if not authenticated
     */
    public function requireAuth()
    {
        if (!$this->isAuthenticated()) {
            http_response_code(401);
            throw new Exception('Authentication required');
        }

        return $this->getCurrentUserId();
    }

    /**
     * Change password for current user
     */
    public function changePassword($currentPassword, $newPassword)
    {
        $userId = $this->requireAuth();

        // Validate new password
        if (strlen($newPassword) < 8) {
            throw new Exception('New password must be at least 8 characters long');
        }

        // Get current password hash
        $stmt = $this->db->prepare("SELECT password_hash FROM users WHERE id = ?");
        $stmt->execute([$userId]);
        $user = $stmt->fetch();

        if (!$user) {
            throw new Exception('User not found');
        }

        // Verify current password
        if (!password_verify($currentPassword, $user['password_hash'])) {
            throw new Exception('Current password is incorrect');
        }

        // Hash new password
        $newPasswordHash = password_hash($newPassword, PASSWORD_BCRYPT, ['cost' => 12]);

        // Update password
        $stmt = $this->db->prepare("UPDATE users SET password_hash = ? WHERE id = ?");
        $stmt->execute([$newPasswordHash, $userId]);

        return true;
    }

    /**
     * Update user profile
     */
    public function updateProfile($fullName = null, $email = null)
    {
        $userId = $this->requireAuth();

        $updates = [];
        $params = [];

        if ($fullName !== null) {
            $updates[] = "full_name = ?";
            $params[] = $fullName;
        }

        if ($email !== null) {
            // Validate email
            if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
                throw new Exception('Invalid email format');
            }

            // Check if email already exists (for other users)
            $stmt = $this->db->prepare("SELECT id FROM users WHERE email = ? AND id != ?");
            $stmt->execute([$email, $userId]);
            if ($stmt->fetch()) {
                throw new Exception('Email already exists');
            }

            $updates[] = "email = ?";
            $params[] = $email;
            $_SESSION['email'] = $email;
        }

        if (empty($updates)) {
            throw new Exception('No updates provided');
        }

        $params[] = $userId;
        $sql = "UPDATE users SET " . implode(', ', $updates) . " WHERE id = ?";

        $stmt = $this->db->prepare($sql);
        $stmt->execute($params);

        if ($fullName !== null) {
            $_SESSION['full_name'] = $fullName;
        }

        return $this->getCurrentUser();
    }

    /**
     * Check if current user is admin
     */
    public function isAdmin()
    {
        if (!$this->isAuthenticated()) {
            return false;
        }

        return isset($_SESSION['is_admin']) && $_SESSION['is_admin'] == 1;
    }

    /**
     * Require admin privileges (middleware)
     * Throws exception if not admin
     */
    public function requireAdmin()
    {
        if (!$this->isAuthenticated()) {
            http_response_code(401);
            throw new Exception('Authentication required');
        }

        if (!$this->isAdmin()) {
            http_response_code(403);
            throw new Exception('Admin privileges required');
        }

        return $this->getCurrentUserId();
    }
}
