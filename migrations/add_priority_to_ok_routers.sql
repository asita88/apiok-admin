ALTER TABLE `ok_routers` ADD COLUMN `priority` int NOT NULL DEFAULT 0 COMMENT 'match order, larger matches first' AFTER `router_path`;
