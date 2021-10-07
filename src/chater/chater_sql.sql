
USE kpnmdb;

CREATE TABLE `groups` (
	`id`           CHAR(32) PRIMARY KEY
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- ALTER TABLE `groups` `id`           CHAR(32) PRIMARY KEY;
ALTER TABLE `groups` ADD `name`         VARCHAR(32) UNIQUE NOT NULL;
ALTER TABLE `groups` ADD `description`  VARCHAR(255) NOT NULL;
ALTER TABLE `groups` ADD `owner`        CHAR(32) NOT NULL;
ALTER TABLE `groups` ADD FOREIGN KEY groups_owner_id(`owner`) REFERENCES users(`id`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `groups` ADD `type`         INTEGER(1) NOT NULL;

CREATE TABLE `members` (
	`group_id`     CHAR(32) NOT NULL,
	FOREIGN KEY members_group_id(`group_id`) REFERENCES `groups`(`id`) ON UPDATE CASCADE ON DELETE CASCADE,
	`user_id`      CHAR(32) NOT NULL,
	PRIMARY KEY (`group_id`, `user_id`),
	`type`         INTEGER(1) NOT NULL,
	`last_read`    DATETIME NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- CREATE TABLE `messages_group_46ec1bcf-e187-44d9-919b-6781a5c612bc` (
-- 	`id` CHAR(32) PRIMARY KEY,
-- 	`date` DATETIME NOT NULL,
-- 	`owner` CHAR(32) NOT NULL,
-- 	`type` INTEGER(1) UNSIGNED NOT NULL,
-- 	`data` TEXT NOT NULL,
-- 	`sdata` VARCHAR(255) NOT NULL,
-- 	`nextid` CHAR(32) NOT NULL,
-- 	FOREIGN KEY messages_next_id(`nextid`) REFERENCES `messages_group_46ec1bcf-e187-44d9-919b-6781a5c612bc`(`id`) ON UPDATE CASCADE ON DELETE CASCADE,
-- 	`isend` BOOLEAN NOT NULL
-- )ENGINE=InnoDB DEFAULT CHARSET=utf8;


