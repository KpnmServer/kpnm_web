
USE kpnmdb;

CREATE TABLE captchas (
	`id`       CHAR(20) NOT NULL PRIMARY KEY,
	`value`    CHAR(6) NOT NULL,
	`overtime` DATETIME NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
