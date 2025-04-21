-- Create "perms" table
CREATE TABLE `perms` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `code` text NOT NULL, `name` text NOT NULL, `description` text NULL, `perm_group` integer NULL, CONSTRAINT `perms_perm_groups_group` FOREIGN KEY (`perm_group`) REFERENCES `perm_groups` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "perms_code_key" to table: "perms"
CREATE UNIQUE INDEX `perms_code_key` ON `perms` (`code`);
-- Create "perm_groups" table
CREATE TABLE `perm_groups` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `code` text NOT NULL, `name` text NOT NULL);
-- Create index "perm_groups_code_key" to table: "perm_groups"
CREATE UNIQUE INDEX `perm_groups_code_key` ON `perm_groups` (`code`);
