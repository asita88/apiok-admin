ALTER TABLE `ok_users`
  ADD COLUMN `role` varchar(32) NOT NULL DEFAULT 'admin' COMMENT 'admin|viewer|operator' AFTER `email`;
