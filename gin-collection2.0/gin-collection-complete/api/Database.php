<?php
// Database configuration and initialization

class Database {
    private static $instance = null;
    private $db;
    private $dbPath;

    private function __construct() {
        $this->dbPath = __DIR__ . '/../db/gin_collection.db';
        
        // Create database directory if it doesn't exist
        $dbDir = dirname($this->dbPath);
        if (!file_exists($dbDir)) {
            mkdir($dbDir, 0755, true);
        }

        try {
            $this->db = new PDO('sqlite:' . $this->dbPath);
            $this->db->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            $this->db->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_ASSOC);
            
            // Enable foreign keys
            $this->db->exec('PRAGMA foreign_keys = ON;');
            
            // Initialize database schema
            $this->initSchema();
        } catch (PDOException $e) {
            error_log("Database connection failed: " . $e->getMessage());
            throw $e;
        }
    }

    public static function getInstance() {
        if (self::$instance === null) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    public function getConnection() {
        return $this->db;
    }

    private function initSchema() {
        $schemaFile = __DIR__ . '/../db/schema.sql';
        
        if (file_exists($schemaFile)) {
            $schema = file_get_contents($schemaFile);
            $this->db->exec($schema);
        }
    }

    public function query($sql, $params = []) {
        try {
            $stmt = $this->db->prepare($sql);
            $stmt->execute($params);
            return $stmt;
        } catch (PDOException $e) {
            error_log("Query failed: " . $e->getMessage());
            throw $e;
        }
    }

    public function lastInsertId() {
        return $this->db->lastInsertId();
    }
}
