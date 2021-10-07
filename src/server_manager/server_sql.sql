
USE kpnmdb;

CREATE TABLE `servers` (
	`id`          VARCHAR(32) PRIMARY KEY,
	`name`        VARCHAR(64) NOT NULL,
	`version`     VARCHAR(32) NOT NULL,
	`description` VARCHAR(255) NOT NULL,
	`addrstr`     VARCHAR(111) NOT NULL,
	`group`       CHAR(32) NOT NULL,
	FOREIGN KEY servers_group(`group`) REFERENCES `groups`(`id`) ON UPDATE CASCADE ON DELETE CASCADE,
	`status`      INTEGER(1) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
