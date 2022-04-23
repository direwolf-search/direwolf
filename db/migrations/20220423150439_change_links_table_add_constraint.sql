-- migrate:up

ALTER TABLE `links`
  ADD CONSTRAINT `links_from_id_foreign` FOREIGN KEY (`from_id`) REFERENCES `hosts` (`id`) ON DELETE SET NULL;

-- migrate:down

ALTER TABLE `links`
  DROP FOREIGN KEY `links_from_id_foreign`;