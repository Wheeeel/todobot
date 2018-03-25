--
-- DbNinja v3.2.7 for MySQL
--
-- Dump date: 2018-03-25 04:18:00 (UTC)
-- Database: tasklist
--

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
/*!40101 SET NAMES utf8 */;

CREATE DATABASE `tasklist` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

USE `tasklist`;

--
-- Structure for table: active_task_instance
--
CREATE TABLE `active_task_instance` (
  `instance_uuid` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `task_id` int(11) NOT NULL,
  `instance_state` int(11) NOT NULL,
  `reminder_state` int(11) DEFAULT NULL,
  `reminder_expression` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  `notify_to_id` bigint(20) NOT NULL,
  `start_at` timestamp NULL DEFAULT NULL,
  `end_at` timestamp NULL DEFAULT NULL,
  `wander_times` int(11) NOT NULL DEFAULT '0',
  `cooldown` int(11) DEFAULT '30' COMMENT '冷却时间, 单位为秒',
  `phrase_group_uuid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '语料组 ID, 什么语料组会显示给用户作为提醒信息',
  UNIQUE KEY `uk_instance_uuid` (`instance_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


--
-- Structure for table: phrases
--
CREATE TABLE `phrases` (
  `uuid` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '语料 UUID',
  `phrase` longtext COLLATE utf8_unicode_ci COMMENT '用户自定义语料，长度不能超过 50 个中文字符',
  `create_by` int(11) DEFAULT NULL COMMENT '创建语料的用户ID',
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建语料时间',
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新记录时间',
  `show` varchar(50) COLLATE utf8_unicode_ci DEFAULT 'yes' COMMENT '本句子是否可以被展示出来（有问题句会被屏蔽）',
  `group_uuid` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '语料属于的语料组，用于和提醒信息关联到一起',
  UNIQUE KEY `uk_uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='用户自定义语料表';


--
-- Structure for table: task_done
--
CREATE TABLE `task_done` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `task_id` int(10) NOT NULL,
  `by` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=636 DEFAULT CHARSET=utf8;


--
-- Structure for table: tasks
--
CREATE TABLE `tasks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(10) NOT NULL,
  `content` longtext NOT NULL COMMENT 'task content',
  `enroll_cnt` int(10) NOT NULL DEFAULT '4',
  `chat_id` bigint(20) DEFAULT NULL,
  `create_by` int(11) NOT NULL DEFAULT '0' COMMENT '创建task的用户ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=747 DEFAULT CHARSET=utf8;


--
-- Structure for table: users
--
CREATE TABLE `users` (
  `uuid` varchar(255) CHARACTER SET utf8 NOT NULL,
  `id` int(11) NOT NULL,
  `user_name` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `disp_name` varchar(1024) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `dont_track` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT 'no',
  `moyu_phrase_uuid` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '开启一个摸鱼任务的时候的默认提醒语料组',
  `phrase_uuid` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '开启一个正式任务的时候的默认语料组'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

