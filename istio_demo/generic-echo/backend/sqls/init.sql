CREATE DATABASE arpg CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE `arpg`.`player` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(128) NULL,
  `account` VARCHAR(32) NOT NULL,
  `password` VARCHAR(32) NOT NULL ,
  PRIMARY KEY (`id`,`account`)) 