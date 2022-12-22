CREATE TABLE IF NOT EXISTS `users` 
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