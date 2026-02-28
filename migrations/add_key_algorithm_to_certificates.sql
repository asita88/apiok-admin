ALTER TABLE `ok_certificates` ADD COLUMN `key_algorithm` varchar(50) NOT NULL DEFAULT '' COMMENT 'Key algorithm e.g. rsa2048, rsa4096, ecdsa_p256' AFTER `ca_provider`;
