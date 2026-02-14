-- Rename "users" table to "user"
ALTER TABLE `users` RENAME TO `user`;
-- Rename index "users_email_key" to "user_email_key"
DROP INDEX IF EXISTS `users_email_key`;
CREATE UNIQUE INDEX `user_email_key` ON `user` (`email`);
