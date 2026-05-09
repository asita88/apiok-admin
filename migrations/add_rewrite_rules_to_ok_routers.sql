ALTER TABLE `ok_routers` ADD COLUMN `rewrite_rules` text COMMENT 'APISIX-style rewrite (JSON)' AFTER `proxy_set_header`;
