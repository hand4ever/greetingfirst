-- MySQL users table DDL
-- Run this script manually before starting the application:
--   mysql -u root -p demo < migrations/001_user.mysql.sql

CREATE TABLE IF NOT EXISTS `users` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `phone` VARCHAR(20) NOT NULL COMMENT 'cellphone number',
    `realname` VARCHAR(100) DEFAULT NULL,
    `username` VARCHAR(20) DEFAULT NULL,
    `age` INT DEFAULT NULL,
    `password_hash` VARCHAR(200) DEFAULT NULL,
    `created_at` DATETIME DEFAULT NULL,
    `updated_at` DATETIME DEFAULT NULL,
    `deleted_at` DATETIME DEFAULT NULL,
    `phone_active` VARCHAR(20) GENERATED ALWAYS AS (IF(deleted_at IS NULL, phone, NULL)) STORED,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_phone_active` (`phone_active`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
