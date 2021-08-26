
CREATE DATABASE kpnmdb;

USE kpnmdb;

CREATE TABLE sessions (
	`uuid`      CHAR(32) NOT NULL,
	`key`       VARCHAR(64) NOT NULL,
	`value`     VARCHAR(256) NOT NULL DEFAULT "",
	`overtime`  DATETIME NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
