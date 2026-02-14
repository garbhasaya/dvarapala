-- Create "users" table
CREATE TABLE `users` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `firstname` text NOT NULL, `lastname` text NOT NULL, `email` text NOT NULL, `password` text NOT NULL, `status` integer NOT NULL DEFAULT (1), `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX `users_email_key` ON `users` (`email`);
