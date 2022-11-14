CREATE DATABASE IF NOT EXISTS `bookstack`;

USE `bookstack`;

CREATE TABLE `user` (
  `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sid`        VARCHAR(10)  NOT NULL UNIQUE COMMENT "学生学号",
  `password`   VARCHAR(20)  NOT NULL  COMMENT "密码",
  `username`   VARCHAR(25)  NOT NULL  COMMENT "姓名",
  `phone`      VARCHAR(11)  NOT NULL  COMMENT "联系方式",
  `book_num`   TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "已经借阅的数量",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;


CREATE TABLE `manager` (
  `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `mid`        VARCHAR(10)  NOT NULL UNIQUE COMMENT "管理员账号",
  `password`   VARCHAR(20)  NOT NULL  COMMENT "密码",
  `username`   VARCHAR(25)  NOT NULL  COMMENT "姓名",
  `phone`      VARCHAR(11)  NOT NULL  COMMENT "联系方式",
  `permissionlevel`    INT  NOT NULL  DEFAULT 0 COMMENT "权限等级-1-2-3",
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `books` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `ISBN`        VARCHAR(13)  NOT NULL UNIQUE COMMENT "图书唯一资源标志符",
  `name`        VARCHAR(20)  NOT NULL  COMMENT "图书名",
  `author`      VARCHAR(20)  NOT NULL  COMMENT "作者",
  `publishing_house`    VARCHAR(20) NOT NULL  COMMENT "出版社",
  `state`       TINYINT(1) NOT NULL DEFAULT 0 COMMENT "是否被借出0-在馆 1-借出",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `borrow_record` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sid`         VARCHAR(10)  NOT NULL UNIQUE COMMENT "学生学号",
  `ISBN`        VARCHAR(13)  NOT NULL UNIQUE COMMENT "图书唯一资源标志符",
  `borrow_time` DATETIME     NOT NULL DEFAULT now() COMMENT "借出时间",
  `should_back_time` DATETIME    NOT NULL  COMMENT "借出时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `fine_record` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sid`         VARCHAR(10)  NOT NULL UNIQUE COMMENT "学生学号",
  `ISBN`        VARCHAR(13)  NOT NULL UNIQUE COMMENT "被罚款图书唯一资源标志符",
  `borrow_time` DATETIME     NOT NULL  COMMENT "借出时间",
  `back_time`   DATETIME     NOT NULL DEFAULT now() COMMENT "归还时间",
  `fine_time`   DATETIME     NOT NULL DEFAULT now() COMMENT "罚款时间",
  `amount`      INT          NOT NULL DEFAULT 0 COMMENT "金额",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `back_record` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `sid`         VARCHAR(10)  NOT NULL UNIQUE COMMENT "学生学号",
  `ISBN`        VARCHAR(13)  NOT NULL UNIQUE COMMENT "图书唯一资源标志符",
  `borrow_time` DATETIME     NOT NULL  COMMENT "借出时间",
  `back_time`   DATETIME     NOT NULL DEFAULT now() COMMENT "归还时间",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

CREATE TABLE `book_num` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`        VARCHAR(20)  NOT NULL  COMMENT "图书名",
  `author`      VARCHAR(20)  NOT NULL  COMMENT "作者",
  `publishing_house`    VARCHAR(20) NOT NULL  COMMENT "出版社",
  `num`       TINYINT(1) NOT NULL DEFAULT 0 COMMENT "在馆数量",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;



