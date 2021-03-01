# create databases
CREATE DATABASE IF NOT EXISTS `test` CHARSET utf8mb4;

# select database to use for this test
USE test ;

CREATE TABLE IF NOT EXISTS `todo`
(
    `id`            int(10) unsigned    NOT NULL AUTO_INCREMENT,
    `title`         varchar(255)        NOT NULL,
    `content`       mediumtext,
    `stat`          tinyint(1)          NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

LOCK TABLES `todo` WRITE;
/*!40000 ALTER TABLE `todo` DISABLE KEYS */;

UNLOCK TABLES;

