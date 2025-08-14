-- Rollback migration orders tables
USE order_service_db;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
