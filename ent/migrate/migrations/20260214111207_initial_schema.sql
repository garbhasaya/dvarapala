-- Create "user" table
CREATE TABLE `user` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `firstname` text NOT NULL, `lastname` text NOT NULL, `email` text NOT NULL, `password` text NOT NULL, `status` integer NOT NULL DEFAULT (1), `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- Create index "user_email_key" to table: "user"
CREATE UNIQUE INDEX `user_email_key` ON `user` (`email`);
