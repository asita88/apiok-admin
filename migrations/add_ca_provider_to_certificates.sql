ALTER TABLE `ok_certificates` ADD COLUMN `ca_provider` varchar(50) NOT NULL DEFAULT '' COMMENT 'CA provider e.g. letsencrypt, manual' AFTER `sni`;
