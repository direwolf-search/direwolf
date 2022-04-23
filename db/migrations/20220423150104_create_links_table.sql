-- migrate:up

CREATE TABLE IF NOT EXISTS `links`  (
	`id` BIGINT unsigned NOT NULL AUTO_INCREMENT COMMENT 'Link identifier',
	`from_id` BIGINT unsigned NULL COMMENT 'Identifier from hosts table',
	`from` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Host URL',
	`body` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Href attribute value',
	`snippet` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Link snippet',
	`is_v3` TINYINT(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Onion address version flag',
	`created_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Link create date',
	`updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Link update date',
	KEY `links_from_id_index` (`from_id`) USING BTREE,
    KEY `links_body_snippet_index` (`body`,`snippet`) USING BTREE,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- migrate:down

DROP TABLE IF EXISTS `links`

