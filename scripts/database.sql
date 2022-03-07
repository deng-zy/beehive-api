CREATE DATABASE `beehive` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci; /*!40100 DEFAULT CHARACTER SET utf8mb4 */
USE `beehive`;
CREATE TABLE `clients` (
  `id` bigint(20) unsigned NOT NULL COMMENT '数据表主键',
  `name` varchar(64) NOT NULL COMMENT '名称',
  `secret` char(24) NOT NULL COMMENT '秘钥',
  `enabled` tinyint(1) unsigned NOT NULL COMMENT '是否已激活(1:yes, 0:no)',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='客户端';
INSERT INTO `clients` (`id`, `name`, `secret`, `enabled`, `created_at`, `updated_at`) VALUES
(1495691462048743424, 'beehive', 'GuVacPSuOwrPmjcAdWDvvkLR', 1, '2022-02-21 17:26:48', '2022-02-21 17:26:48');
CREATE TABLE `topics` (
  `id` bigint(20) unsigned NOT NULL COMMENT '数据表主键',
  `name` varchar(64) NOT NULL COMMENT '名称',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `events` (
  `id` bigint(20) unsigned NOT NULL COMMENT '数据表主键',
  `topic` varchar(64) NOT NULL COMMENT 'topic 主题',
  `payload` varchar(2048) NOT NULL COMMENT '消息',
  `publisher` varchar(64) NOT NULL COMMENT '发布者',
  `published_at` datetime DEFAULT NULL COMMENT '发布时间 发布时间可能早于创建时间',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;