ALTER TABLE `ok_certificates` ADD COLUMN `issuer` varchar(255) NOT NULL DEFAULT '' COMMENT 'Certificate issuer' AFTER `key_algorithm`;
