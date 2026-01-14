-- Drop tables in reverse order (respecting foreign key constraints)

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS tasting_sessions;
DROP TABLE IF EXISTS gin_photos;
DROP TABLE IF EXISTS gin_cocktails;
DROP TABLE IF EXISTS gin_botanicals;
DROP TABLE IF EXISTS cocktail_ingredients;
DROP TABLE IF EXISTS cocktails;
DROP TABLE IF EXISTS botanicals;
DROP TABLE IF EXISTS gins;
DROP TABLE IF EXISTS usage_metrics;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;
