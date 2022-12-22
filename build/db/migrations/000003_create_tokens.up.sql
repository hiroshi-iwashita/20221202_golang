CREATE TABLE IF NOT EXISTS `tokens`
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