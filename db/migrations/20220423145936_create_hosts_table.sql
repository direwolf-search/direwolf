-- migrate:up
CREATE TABLE IF NOT EXISTS `hosts`(
	`id` BIGINT unsigned NOT NULL AUTO_INCREMENT COMMENT 'Host identifier',
	`url` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Host URL',
	`h1` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT 'Host h1 HTML tag',
	`domain` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'Host domain',
	`content_type` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'Host domain',
	`title` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'Host title HTML tag',
	`meta` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'Host metatags',
	`md5hash` CHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'MD5 hash of full host HTML code',
	`text` LONGTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'Text of host with all HTML tags cutted',
	`status` TINYINT(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Host online status',
	`http_status` INT NOT NULL DEFAULT 0 COMMENT 'Host HTTP response code',
	`link_num` INT NULL DEFAULT NULL COMMENT 'Number of outgoing links',
	`created_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Create date',
	`updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update date',
	`visited_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last visit date',
	`visits_num` BIGINT NULL DEFAULT NULL COMMENT 'Number of outgoing links',
	KEY `hosts_url_index` (`url`) USING BTREE,
    KEY `hosts_hash_index` (`md5hash`) USING BTREE,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- migrate:down

DROP TABLE IF EXISTS `hosts`