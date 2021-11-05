-- -----------------------
-- Database: `direwolf_db`
-- -----------------------

-- --------------------------------------------------------

SET FOREIGN_KEY_CHECKS=0;
SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

-- ----------------
-- Create database
-- ----------------

CREATE DATABASE IF NOT EXISTS `direwolf_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `direwolf_db`;

-- --------------------
-- Create table `hosts`
-- --------------------

DROP TABLE IF EXISTS `hosts`;
CREATE TABLE IF NOT EXISTS `hosts`(
	`id` BIGINT unsigned NOT NULL AUTO_INCREMENT COMMENT 'Host identifier',
	`url` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Host URL',
	`h1` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Host h1 HTML tag',
	`title` TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Host title HTML tag',
	`hash` CHAR(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'MD5 hash of full host HTML code',
	`text` LONGTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'Text of host with all HTML tags cutted',
	`status` TINYINT(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Host online status',
	`http_status` VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'Host HTTP response code',
	`created_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Create date',
	`updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update date',
	KEY `hosts_url_index` (`url`) USING BTREE,
    KEY `hosts_hash_index` (`hash`) USING BTREE,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------
-- Create table `links`
-- --------------------

CREATE TABLE `links` (
	`id` BIGINT unsigned NOT NULL AUTO_INCREMENT COMMENT 'Link identifier',
	`from_id` BIGINT unsigned NULL COMMENT 'Identifier from hosts table',
	`body` VARCHAR(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Href attribute value',
	`snippet` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Link snippet',
	`is_v3` TINYINT(1) unsigned NOT NULL DEFAULT 0 COMMENT 'Onion address version flag',
	`created_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Link create date',
	`updated_at` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Link update date',
	KEY `links_from_id_index` (`from_id`) USING BTREE,
KEY `links_body_snippet_index` (`body`,`snippet`) USING BTREE,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Foreign key constraints
--

--
-- Foreign key constraints for `links` table
--
ALTER TABLE `links`
  ADD CONSTRAINT `links_from_id_foreign` FOREIGN KEY (`from_id`) REFERENCES `hosts` (`id`) ON DELETE SET NULL;

SET FOREIGN_KEY_CHECKS=1;
COMMIT;