-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "app" table
CREATE TABLE `app` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `status` integer NOT NULL DEFAULT (1), `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- Create index "app_name_key" to table: "app"
CREATE UNIQUE INDEX `app_name_key` ON `app` (`name`);
-- 2. Ensure a default app exists (required for backfill)
INSERT OR IGNORE INTO `app`
(`id`, `name`, `status`, `created_at`, `updated_at`)
VALUES
(1, 'default', 1, datetime('now'), datetime('now'));
-- Create "new_user" table
CREATE TABLE `new_user` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `firstname` text NOT NULL, `lastname` text NOT NULL, `email` text NOT NULL, `password` text NOT NULL, `status` integer NOT NULL DEFAULT (1), `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `app_id` integer NOT NULL, CONSTRAINT `user_app_users` FOREIGN KEY (`app_id`) REFERENCES `app` (`id`) ON DELETE NO ACTION);
-- Copy rows from old table "user" to new temporary table "new_user"
INSERT INTO `new_user` (`id`, `firstname`, `lastname`, `email`, `password`, `status`, `created_at`, `updated_at`, `app_id`) SELECT `id`, `firstname`, `lastname`, `email`, `password`, `status`, `created_at`, `updated_at`, 1 AS `app_id` FROM `user`;
-- Drop "user" table after copying rows
DROP TABLE `user`;
-- Rename temporary table "new_user" to "user"
ALTER TABLE `new_user` RENAME TO `user`;
-- Create index "user_email_key" to table: "user"
CREATE UNIQUE INDEX `user_email_key` ON `user` (`email`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
