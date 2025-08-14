-- Rollback migration users table
USE user_service_db;
DROP TABLE IF EXISTS users;