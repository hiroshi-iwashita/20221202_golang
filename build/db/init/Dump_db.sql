CREATE DATABASE IF NOT EXISTS `test_db`;
USE `test_db`;

-- DROP SCHEMA IF EXISTS `test_db`;
-- CREATE SCHEMA `test_db`;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` 
    (
        `id` int(11) NOT NULL AUTO_INCREMENT,
        `user_id` VARCHAR(36) NOT NULL,
        `first_name` VARCHAR(191) NULL,
        `last_name` VARCHAR(191) NULL,
        `email` VARCHAR(191) NOT NULL,
        `password` VARCHAR(191) NOT NULL,
        `email_verified_at` DATETIME(3) NULL,
        `created_at` DATETIME(3) NOT NULL,
        `updated_at` DATETIME(3) NOT NULL,
        `deleted_at` DATETIME(3) NULL,
        PRIMARY KEY (`id`, `user_id`),
        CONSTRAINT `UK_users` 
            UNIQUE (`email`),
        INDEX `IDX_users_user_id` (`user_id`)
    ) 
    DEFAULT CHARACTER SET `utf8mb4`
    COLLATE `utf8mb4_unicode_ci`
;

DROP TABLE IF EXISTS `tokens`;
CREATE TABLE `tokens`
    (
        `id` int(11) NOT NULL AUTO_INCREMENT,
        `user_id` VARCHAR(36) NOT NULL,
        `token` VARCHAR(191) NOT NULL,
        `token_hash` BLOB NOT NULL,
        `created_at` DATETIME(3) NOT NULL,
        `updated_at` DATETIME(3) NOT NULL,
        `expire_at` DATETIME(3) NOT NULL,
        PRIMARY KEY (`id`),
        CONSTRAINT `UK_tokens` 
            UNIQUE (`user_id`, `token`),
        CONSTRAINT `FK_tokens_user_id` 
            FOREIGN KEY (`user_id`) 
            REFERENCES `users` (`user_id`)
            ON UPDATE CASCADE
            ON DELETE CASCADE
    )
    DEFAULT CHARACTER SET `utf8mb4`
    COLLATE `utf8mb4_unicode_ci`
;