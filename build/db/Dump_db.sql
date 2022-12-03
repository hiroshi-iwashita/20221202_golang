USE `test_db`;

DROP TABLE IF EXISTS `sample1`;

CREATE TABLE `sample1` (
  `id` int(7) NOT NULL,
  `title` varchar(256) NOT NULL,
  PRIMARY KEY (`id`)
) ;

LOCK TABLES `sample1` WRITE;

INSERT INTO `sample1` VALUES (1,'title1');

UNLOCK TABLES;