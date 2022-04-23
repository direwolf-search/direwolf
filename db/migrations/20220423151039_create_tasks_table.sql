-- migrate:up
CREATE TABLE IF NOT EXISTS `tasks` (
	`id` BIGINT unsigned NOT NULL AUTO_INCREMENT COMMENT 'Link identifier',
	`of` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Host URL',
	`rule` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Href attribute value',
	`schedule` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Link snippet',
	`skip_next` TINYINT(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Onion address version flag',
	`created_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Link create date',
	`updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Link update date',
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- migrate:down

DROP TABLE IF EXISTS `hosts`