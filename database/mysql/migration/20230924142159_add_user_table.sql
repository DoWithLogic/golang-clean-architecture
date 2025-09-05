-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `contact_type` ENUM('EMAIL', 'PHONE') NOT NULL,
    `contact_value` VARCHAR(320) NOT NULL,
    `birth_date` DATE NULL,
    `language` ENUM('EN', 'ID') NOT NULL DEFAULT 'ID',
    `password` VARCHAR(255) NOT NULL,
    `status` ENUM('PENDING', 'ACTIVE', 'REJECT', 'CLOSED') NOT NULL DEFAULT 'PENDING',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_contact` (`contact_type`, `contact_value`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `users`;
-- +goose StatementEnd
