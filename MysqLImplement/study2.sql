DROP   DATABASE IF EXISTS `VCDRent`;

CREATE DATABASE IF NOT EXISTS `VCDRent`;

USE `VCDRent`;

-- vcd 入库
CREATE TABLE `vcd_in_records` (
  `id`          INT UNSIGNED  NOT NULL  AUTO_INCREMENT,
  `vname`       VARCHAR(30)   NOT NULL  COMMENT "vcd名",
  `author`      VARCHAR(20)   NOT NULL  COMMENT "vcd作者",
  `price`       Decimal(10,2) NOT NULL  COMMENT "vcd价格",
  `vcd_type`    VARCHAR(20)   NOT NULL  COMMENT "vcd类型",
  `time`        DATETIME      NOT NULL DEFAULT now() COMMENT "入库时间",
  `amount`      INT           NOT NULL  COMMENT "一次性入库的数量",
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

select * from `vcd_in_records`;

-- vcd 信息
CREATE TABLE `vcds` (
  `id`          INT UNSIGNED  NOT NULL  AUTO_INCREMENT,
  `name`        VARCHAR(30)   NOT NULL   COMMENT "vcd名",
  `author`      VARCHAR(20)   NOT NULL   COMMENT "vcd作者",
  `price`       Decimal(10,2) NOT NULL   COMMENT "vcd价格",
  `vcd_type`    VARCHAR(20)   NOT NULL   COMMENT "vcd类型",
  `vcdnum`      INT           NOT NULL DEFAULT 0  COMMENT "vcd数量",
  PRIMARY KEY (`id`),
  INDEX (name(30))  ,
  INDEX (author(20))  
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

select * from `vcds`;

-- 会员表
CREATE TABLE `customers` (
  `id`         INT UNSIGNED  NOT NULL  AUTO_INCREMENT,
  `phone`      VARCHAR(11)  NOT NULL UNIQUE COMMENT "顾客账号就是注册手机号或者第三方帐号",
  `password`   VARCHAR(20)  NOT NULL  COMMENT "密码",
  `nickname`   VARCHAR(25)  NOT NULL  COMMENT "昵称",
   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

select * from `customers`;


-- vcd 归还表
CREATE TABLE `back_records` (
  `id`          INT UNSIGNED  NOT NULL  AUTO_INCREMENT,
  `name`        VARCHAR(30)   NOT NULL  COMMENT "vcd名",
  `author`      VARCHAR(20)   NOT NULL  COMMENT "vcd作者",
  `phone`       VARCHAR(11)   NOT NULL  COMMENT "顾客账号就是注册手机号",
  `back_time`   DATETIME      NOT NULL DEFAULT now() COMMENT "入库时间",
  `back_num`      INT         NOT NULL  COMMENT "一次性归还的数量",
   PRIMARY KEY (`id`),
   FOREIGN KEY (`phone`) REFERENCES `customers` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

select * from `back_records`;

-- vcd 借阅表
CREATE TABLE `borrow_records` (
  `id`          INT UNSIGNED  NOT NULL  AUTO_INCREMENT,
  `name`        VARCHAR(30)  NOT NULL   COMMENT "vcd名",
  `author`      VARCHAR(20)  NOT NULL   COMMENT "vcd作者",
  `phone`       VARCHAR(11)  NOT NULL   COMMENT "顾客账号就是注册手机号",
  `borrow_time` DATETIME     NOT NULL DEFAULT now() COMMENT "借出时间",
  `borrow_num`  INT          NOT NULL  COMMENT "借出数量",
  PRIMARY KEY (`id`),
  FOREIGN KEY (`name` )  REFERENCES `vcds` (`name`),
  FOREIGN KEY (`author`) REFERENCES `vcds` (`author`),
  FOREIGN KEY (`phone`)  REFERENCES `customers` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;

select * from `borrow_records`;

-- vcd 售卖表
CREATE TABLE `sell_records` (
  `id`          INT UNSIGNED NOT NULL  AUTO_INCREMENT,
  `name`        VARCHAR(30)  NOT NULL               COMMENT "vcd名",
  `author`      VARCHAR(20)  NOT NULL               COMMENT "vcd作者",
  `phone`       VARCHAR(11)  NOT NULL               COMMENT "顾客账号就是注册手机号",
  `sell_time`   DATETIME     NOT NULL DEFAULT now() COMMENT "售卖时间",
  `sell_num`    INT          NOT NULL               COMMENT "售卖数量",
  PRIMARY KEY (`id`),
  FOREIGN KEY (`name` )  REFERENCES `vcds` (`name`),
  FOREIGN KEY (`author`) REFERENCES `vcds` (`author`),
  FOREIGN KEY (`phone`)  REFERENCES `customers` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;


-- vcd 某段时间内的销售统计表
CREATE TABLE `sell_records_statistics` (
  `id`          INT UNSIGNED NOT NULL  AUTO_INCREMENT,
  `name`        VARCHAR(30)  NOT NULL               COMMENT "vcd名",
  `author`      VARCHAR(20)  NOT NULL               COMMENT "vcd作者",
  `num`         INT          NOT NULL               COMMENT "售卖总数量",
  PRIMARY KEY (`id`),
  FOREIGN KEY (`name` )  REFERENCES `vcds` (`name`),
  FOREIGN KEY (`author`) REFERENCES `vcds` (`author`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;


-- vcd 某段时间内的借阅统计表
CREATE TABLE `borrow_records_statistics` (
  `id`          INT UNSIGNED NOT NULL  AUTO_INCREMENT,
  `name`        VARCHAR(30)  NOT NULL               COMMENT "vcd名",
  `author`      VARCHAR(20)  NOT NULL               COMMENT "vcd作者",
  `num`         INT          NOT NULL               COMMENT "借阅总数量",
  PRIMARY KEY (`id`),
  FOREIGN KEY (`name` )  REFERENCES `vcds` (`name`),
  FOREIGN KEY (`author`) REFERENCES `vcds` (`author`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4;



-- 初始化 vcd 入库触发器
DROP TRIGGER  IF EXISTS  VCDRent.after_vcd_in_records_insert;
DELIMITER $$
CREATE TRIGGER after_vcd_in_records_insert
    AFTER INSERT ON vcd_in_records
    FOR EACH ROW 
BEGIN
    UPDATE  `vcds`
    SET vcdnum = vcdnum + new.amount
    WHERE  name = new.vname  and author= new.author;
END$$
DELIMITER ;


-- 初始化的vcds 信息导入
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('Addicted to Her Sadness', 'Benjamin Kheng', 100.7, '音频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('Better with you', 'Benjamin Kheng', 200.7, '音频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('遗失的心跳', '糯米Nomi', 300.7, '音频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('I Bet', 'Ciara', 100.7, '音频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('Why You', 'Sway Bleu', 87.7, '音频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('Call On Me', 'Nelly', 56.2, '音频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('立冬', '音阙诗听 / 赵方婧', 144.8, '音频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('错位月光', 'Bo Peep / CDY', 100.7, '音频', 1);

INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('西线无战事','吕克·费特', 200.7, '视频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('野蛮人','索茜·贝肯', 200.7, '视频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('地狱尖兵','阿列克谢·克拉夫琴科', 300.7, '视频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('钢琴家','杰西卡·查斯坦', 200.7, '视频', 1);
INSERT INTO `vcds` (name ,author, price, vcd_type, vcdnum) VALUES ('教父','埃迪·雷德', 400.7, '视频', 1);

-- 初始化的vcd_in_records 信息导入
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('Addicted to Her Sadness', 'Benjamin Kheng', 100.7, '音频', 4);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('Better with you', 'Benjamin Kheng', 200.7, '音频', 5);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('遗失的心跳', '糯米Nomi', 300.7, '音频', 9);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('I Bet', 'Ciara', 100.7, '音频', 6);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('Why You', 'Sway Bleu', 87.7, '音频', 4);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('Call On Me', 'Nelly', 56.2, '音频', 3);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('立冬', '音阙诗听 / 赵方婧', 144.8, '音频', 10);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('错位月光', 'Bo Peep / CDY', 100.7, '音频', 4);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('错位月光', 'Bo Peep / CDY', 100.7, '音频', 4);
INSERT INTO `vcd_in_records` (vname ,author, price, vcd_type, amount) VALUES ('错位月光', 'Bo Peep / CDY', 100.7, '音频', 4);


-- 验证触发器
select * from `vcds`;

select * from `vcd_in_records`;


-- 初始化的会员信息导入
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769211', '123456', 'Benjamin Kheng');
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769212', '123456', 'Benjamin Kheng');
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769213', '123456', '糯米Nomi');
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769214', '123456', 'Ciara');
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769215', '123456', 'Sway Bleu');
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769216', '123456', 'Nelly');
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769217', '123456', '赵方婧');
INSERT INTO `customers` (phone, password, nickname) VALUES ('15102769218', '123456', 'Bo Peep');


select * from `customers`;


-- 初始化 借阅触发器
DROP TRIGGER  IF EXISTS  VCDRent.after_borrow_records_insert;
DELIMITER $$
CREATE TRIGGER after_borrow_records_insert
    AFTER INSERT ON borrow_records
    FOR EACH ROW 
BEGIN
    UPDATE  `vcds`
    SET vcdnum = vcdnum-new.borrow_num
    WHERE vcdnum-new.borrow_num >= 0 and name = new.name and author = new.author;
END$$
DELIMITER ;


-- 初始化的borrow_records 信息导入
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769211', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769212', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769213', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('I Bet', 'Ciara', '15102769211', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('I Bet', 'Ciara', '15102769212', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('I Bet', 'Ciara', '15102769213', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('I Bet', 'Ciara', '15102769214', 1);

INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('Why You', 'Sway Bleu', '15102769211', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('Why You', 'Sway Bleu', '15102769212', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('Why You', 'Sway Bleu', '15102769213', 1);


INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769211', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769212', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769213', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769214', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769215', 1);
INSERT INTO `borrow_records` (name ,author, phone, borrow_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769216', 1);

select * from `borrow_records`;
-- 验证触发器
select * from `vcds`;



-- 初始化 归还触发器
DROP TRIGGER  IF EXISTS  VCDRent.after_back_records_insert;
DELIMITER $$
CREATE TRIGGER after_borrow_records_insert
    AFTER INSERT ON borrow_records
    FOR EACH ROW 
BEGIN
    UPDATE  `vcds`
    SET vcdnum = vcdnum + new.back_num
    WHERE name = new.name and author = new.author;
END$$
DELIMITER ;

-- 初始化的back_records 信息导入
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769211', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769212', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769213', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('I Bet', 'Ciara', '15102769211', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('I Bet', 'Ciara', '15102769212', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('I Bet', 'Ciara', '15102769213', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('I Bet', 'Ciara', '15102769214', 1);

INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('Why You', 'Sway Bleu', '15102769211', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('Why You', 'Sway Bleu', '15102769212', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('Why You', 'Sway Bleu', '15102769213', 1);


INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769211', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769212', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769213', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769214', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769215', 1);
INSERT INTO `back_records` (name ,author, phone, back_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769216', 1);

select * from `back_records`;
-- 验证触发器
select * from `vcds`;


-- 初始化 售卖触发器
DROP TRIGGER  IF EXISTS  VCDRent.after_sell_records_insert;
DELIMITER $$
CREATE TRIGGER after_sell_records_insert
    AFTER INSERT ON sell_records
    FOR EACH ROW 
BEGIN
    UPDATE  `vcds`
    SET vcdnum = vcdnum - new.sell_num
    WHERE  vcdnum - new.sell_num >=0 and  name = new.name and author = new.author;
END$$
DELIMITER ;


-- 初始化的 sell_records 信息导入
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769211', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769212', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('Better with you', 'Benjamin Kheng', '15102769213', 1);

INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('I Bet', 'Ciara', '15102769211', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('I Bet', 'Ciara', '15102769212', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('I Bet', 'Ciara', '15102769213', 1);

INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('Why You', 'Sway Bleu', '15102769211', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('Why You', 'Sway Bleu', '15102769212', 1);


INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769211', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769212', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769213', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769214', 1);
INSERT INTO `sell_records` (name ,author, phone, sell_num) VALUES ('立冬', '音阙诗听 / 赵方婧',  '15102769215', 1);

select * from `sell_records`;
-- 验证触发器
select * from `vcds`;



-- 验证不同种类vcd数量视图
DROP VIEW IF EXISTS `vcd_kind_view`;
CREATE VIEW vcd_kind_view
AS SELECT vcd_type ,COUNT(*) as num FROM `vcds` GROUP BY  vcd_type ;
select * from vcd_kind_view;






-- 调用存储过程


-- 设置分割符为//
delimiter //  
create procedure VCDRent.get_sell_information(transdate text)
begin  -- 开始程序体
declare startdate, enddate datetime;  -- 定义变量
set startdate = date_format(transdate, '%Y-%m-%d');  -- 给起始时间赋值
set enddate = date_add(transdate, interval 1 day);  -- 截止时间赋值为1天以后
-- 重新插入数据
insert into VCDRent.sell_records_statistics(name,author, num) 
select  name, author ,COUNT(*)
from sell_records
where sell_time > startdate and sell_time < enddate  
GROUP by name, author
;
end
//
delimiter ;

CALL VCDRent.get_sell_information('2022-11-13');
select * from sell_records_statistics;


-- 设置分割符为//
delimiter //  
create procedure VCDRent.get_borrow_information(transdate text)
begin  -- 开始程序体
declare startdate, enddate datetime;  -- 定义变量
set startdate = date_format(transdate, '%Y-%m-%d');  -- 给起始时间赋值
set enddate = date_add(transdate, interval 1 day);  -- 截止时间赋值为1天以后
-- 重新插入数据
insert into VCDRent.borrow_records_statistics(name,author, num) 
select  name, author ,COUNT(*)
from borrow_records
where borrow_time > startdate and borrow_time < enddate  
GROUP by name, author
;
end
//
delimiter ;

CALL VCDRent.get_borrow_information('2022-11-13');
select * from borrow_records_statistics;

















