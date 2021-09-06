
USE kpnmdb;

CREATE TABLE users (
	`id`             CHAR(32) PRIMARY KEY,
	`username`       VARCHAR(32) UNIQUE NOT NULL,
	`email`          VARCHAR(64) UNIQUE NOT NULL,
	`password`       CHAR(64) NOT NULL,
	`frozen`         INTEGER(1) NOT NULL,
	`description`    VARCHAR(255) NOT NULL DEFAULT ""
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
