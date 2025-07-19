# ************************************************************
# Sequel Ace SQL dump
# 版本号： 20080
#
# https://sequel-ace.com/
# https://github.com/Sequel-Ace/Sequel-Ace
#
# 主机: 127.0.0.1 (MySQL 8.4.5)
# 数据库: roselle
# 生成时间: 2025-07-19 03:43:49 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# 转储表 ad_whitelist
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ad_whitelist`;

CREATE TABLE `ad_whitelist` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `open_id` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `is_valid` tinyint NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



# 转储表 captcha
# ------------------------------------------------------------

DROP TABLE IF EXISTS `captcha`;

CREATE TABLE `captcha` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `open_id` varchar(256) COLLATE utf8mb4_general_ci NOT NULL,
  `str` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `md5` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
  `is_valid` tinyint NOT NULL DEFAULT '1',
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



# 转储表 gsc
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gsc`;

CREATE TABLE `gsc` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `work_title` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `work_author` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `work_dynasty` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `translation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `intro` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `baidu_wiki` varchar(256) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `audio_id` int NOT NULL DEFAULT '0',
  `foreword` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `annotation_` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `appreciation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `master_comment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
  `layout` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT 'indent',
  PRIMARY KEY (`ID`),
  KEY `IDX_WORK_AUDIO_ID` (`audio_id`),
  KEY `IDX_LAYOUT` (`layout`),
  FULLTEXT KEY `IDX_FULL_TEXT` (`work_author`,`work_title`,`work_dynasty`,`content`) /*!50100 WITH PARSER `ngram` */ ,
  FULLTEXT KEY `idx_work_title` (`work_title`) /*!50100 WITH PARSER `ngram` */ ,
  FULLTEXT KEY `idx_content` (`content`,`foreword`) /*!50100 WITH PARSER `ngram` */ ,
  FULLTEXT KEY `idx_work_author` (`work_author`) /*!50100 WITH PARSER `ngram` */ ,
  FULLTEXT KEY `idx_work_dynasty` (`work_dynasty`) /*!50100 WITH PARSER `ngram` */ 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



# 转储表 user_feedback
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_feedback`;

CREATE TABLE `user_feedback` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `open_id` varchar(256) COLLATE utf8mb4_general_ci NOT NULL,
  `gsc_id` int NOT NULL,
  `feedback_type` int NOT NULL,
  `remark` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



# 转储表 user_like_gsc
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_like_gsc`;

CREATE TABLE `user_like_gsc` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `open_id` varchar(128) NOT NULL,
  `gsc_id` int NOT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_open_id` (`open_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
