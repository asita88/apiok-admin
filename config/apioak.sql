/*
Navicat MySQL Data Transfer

Source Server         : 127.0.0.1
Source Server Version : 80044
Source Host           : 127.0.0.1:3306
Source Database       : apiok

Target Server Type    : MYSQL
Target Server Version : 80044
File Encoding         : 65001

Date: 2025-12-08 17:10:51
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for apiok_data
-- ----------------------------
DROP TABLE IF EXISTS `apiok_data`;
CREATE TABLE `apiok_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型: services/routers/plugins/upstreams/certificates/upstream_nodes',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源名称',
  `data` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'JSON格式的配置数据',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_type_name` (`type`,`name`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据存储表';

-- ----------------------------
-- Records of apiok_data
-- ----------------------------
INSERT INTO `apiok_data` VALUES ('33', 'upstreams', 'up-m00a4aXgZul43zU', '{\"name\":\"up-m00a4aXgZul43zU\",\"algorithm\":\"ROUNDROBIN\",\"connect_timeout\":30000,\"write_timeout\":3000,\"read_timeout\":3000,\"enabled\":true,\"nodes\":[{\"name\":\"un-rmxmDAvLlJt7dpe\"},{\"name\":\"un-tPdCAkBvo9xZOwz\"},{\"name\":\"un-PiXhYWPg81UqbNG\"}]}', '2025-12-05 18:50:39', '2025-12-08 15:57:47');
INSERT INTO `apiok_data` VALUES ('34', 'upstream_nodes', 'un-rmxmDAvLlJt7dpe', '{\"name\":\"un-rmxmDAvLlJt7dpe\",\"address\":\"127.0.0.1\",\"port\":8080,\"weight\":100,\"health\":\"HEALTH\",\"check\":{\"enabled\":false,\"tcp\":true,\"method\":\"\",\"host\":\"\",\"uri\":\"/\",\"interval\":1,\"timeout\":1},\"tags\":{\"123\":\"456\",\"789\":\"123\"}}', '2025-12-05 18:50:39', '2025-12-08 15:57:47');
INSERT INTO `apiok_data` VALUES ('38', 'routers', 'rt-cMQH880AhzqxQnO', '{\"name\":\"rt-cMQH880AhzqxQnO\",\"methods\":[\"ALL\"],\"paths\":[\"/apple\"],\"enabled\":true,\"headers\":{},\"service\":{\"name\":\"sv-DQ0CASH03pMddn3\"},\"upstream\":{\"name\":\"up-BsYcvp9VClgXIX2\"},\"plugins\":[],\"client_max_body_size\":100}', '2025-12-05 18:58:54', '2025-12-08 17:03:38');
INSERT INTO `apiok_data` VALUES ('40', 'services', 'sv-GdO0rkKOo5yJ9RM', '{\"name\":\"sv-GdO0rkKOo5yJ9RM\",\"protocols\":[\"http\"],\"hosts\":[\"aa.apple.com\"],\"ports\":[80],\"plugins\":[{\"name\":\"pc-dTETOfGzqFfQYje\"}],\"enabled\":true}', '2025-12-05 19:22:39', '2025-12-08 17:03:30');
INSERT INTO `apiok_data` VALUES ('41', 'upstreams', 'up-BsYcvp9VClgXIX2', '{\"name\":\"up-BsYcvp9VClgXIX2\",\"algorithm\":\"ROUNDROBIN\",\"connect_timeout\":30000,\"write_timeout\":3000,\"read_timeout\":3000,\"enabled\":true,\"nodes\":[{\"name\":\"un-RTjg5YiOqd2jFVH\"}]}', '2025-12-05 19:49:33', '2025-12-08 17:03:16');
INSERT INTO `apiok_data` VALUES ('42', 'upstream_nodes', 'un-RTjg5YiOqd2jFVH', '{\"name\":\"un-RTjg5YiOqd2jFVH\",\"address\":\"127.0.0.2\",\"port\":8080,\"weight\":100,\"health\":\"HEALTH\",\"check\":{\"enabled\":false,\"tcp\":true,\"method\":\"\",\"host\":\"\",\"uri\":\"/\",\"interval\":1,\"timeout\":1}}', '2025-12-05 19:49:33', '2025-12-08 17:03:16');
INSERT INTO `apiok_data` VALUES ('43', 'upstream_nodes', 'un-tPdCAkBvo9xZOwz', '{\"name\":\"un-tPdCAkBvo9xZOwz\",\"address\":\"127.0.0.2\",\"port\":8080,\"weight\":100,\"health\":\"HEALTH\",\"check\":{\"enabled\":false,\"tcp\":true,\"method\":\"\",\"host\":\"\",\"uri\":\"/\",\"interval\":1,\"timeout\":1}}', '2025-12-08 15:57:47', '2025-12-08 15:57:47');
INSERT INTO `apiok_data` VALUES ('44', 'upstream_nodes', 'un-PiXhYWPg81UqbNG', '{\"name\":\"un-PiXhYWPg81UqbNG\",\"address\":\"127.0.0.3\",\"port\":8080,\"weight\":100,\"health\":\"HEALTH\",\"check\":{\"enabled\":false,\"tcp\":true,\"method\":\"\",\"host\":\"\",\"uri\":\"/\",\"interval\":1,\"timeout\":1}}', '2025-12-08 15:57:47', '2025-12-08 15:57:47');
INSERT INTO `apiok_data` VALUES ('45', 'plugins', 'pc-dTETOfGzqFfQYje', '{\"config\":{\"time_window\":100,\"count\":1000},\"key\":\"limit-count\",\"name\":\"pc-dTETOfGzqFfQYje\"}', '2025-12-08 17:03:30', '2025-12-08 17:03:30');
INSERT INTO `apiok_data` VALUES ('46', 'services', 'sv-DQ0CASH03pMddn3', '{\"name\":\"sv-DQ0CASH03pMddn3\",\"protocols\":[\"http\"],\"hosts\":[\"bb.apple.com\"],\"ports\":[80],\"plugins\":[],\"enabled\":true}', '2025-12-08 17:03:32', '2025-12-08 17:03:32');

-- ----------------------------
-- Table structure for apiok_sync_hash
-- ----------------------------
DROP TABLE IF EXISTS `apiok_sync_hash`;
CREATE TABLE `apiok_sync_hash` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `hash_key` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '哈希键名',
  `hash_value` text COLLATE utf8mb4_unicode_ci COMMENT '哈希值（JSON格式）',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_hash_key` (`hash_key`)
) ENGINE=InnoDB AUTO_INCREMENT=2499017 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='同步哈希表';

-- ----------------------------
-- Records of apiok_sync_hash
-- ----------------------------
INSERT INTO `apiok_sync_hash` VALUES ('2499016', 'sync/update', '{\"old\":\"bd0a57b3a6ce39f912f9babfa19404cc\",\"new\":\"bd0a57b3a6ce39f912f9babfa19404cc\"}', '2025-12-08 15:51:04');

-- ----------------------------
-- Table structure for oak_certificates
-- ----------------------------
DROP TABLE IF EXISTS `oak_certificates`;
CREATE TABLE `oak_certificates` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Certificate id',
  `sni` varchar(150) NOT NULL DEFAULT '' COMMENT 'SNI',
  `certificate` text NOT NULL COMMENT 'Certificate content',
  `private_key` text NOT NULL COMMENT 'Private key content',
  `enable` tinyint unsigned NOT NULL DEFAULT '2' COMMENT 'Certificate enable  1:on  2:off',
  `expired_at` timestamp NULL DEFAULT NULL COMMENT 'Expiration time',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  KEY `IDX_SNI` (`sni`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Certificates';

-- ----------------------------
-- Records of oak_certificates
-- ----------------------------

-- ----------------------------
-- Table structure for oak_plugins
-- ----------------------------
DROP TABLE IF EXISTS `oak_plugins`;
CREATE TABLE `oak_plugins` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Plugin id',
  `plugin_key` varchar(20) NOT NULL DEFAULT '' COMMENT 'Plugin key',
  `icon` varchar(50) NOT NULL DEFAULT '' COMMENT 'Plugin icon',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT 'Plugin type',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT 'Plugin description',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_KEY` (`plugin_key`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Plugins';

-- ----------------------------
-- Records of oak_plugins
-- ----------------------------
INSERT INTO `oak_plugins` VALUES ('1', 'pl-dIhZpgqcCHQzNgT', 'cors', 'icon-cors', '3', '配置服务端CORS（Cross-Origin Resource Sharing，跨域资源共享）的响应头信息', '2025-12-04 15:20:19', '2025-12-04 15:20:19');
INSERT INTO `oak_plugins` VALUES ('2', 'pl-5xO9hzfcHJtpcQT', 'mock', 'icon-mock', '99', '配置模拟API数据，且请求不会转发到上游', '2025-12-04 15:20:19', '2025-12-04 15:20:19');
INSERT INTO `oak_plugins` VALUES ('3', 'pl-xZjvnLQfq2i5GTS', 'key-auth', 'icon-key-auth', '1', '配置身份验证密钥（key密钥字符串）', '2025-12-04 15:20:19', '2025-12-04 15:20:19');
INSERT INTO `oak_plugins` VALUES ('4', 'pl-0FnmajmiO7C8PtX', 'jwt-auth', 'icon-jwt-auth', '1', '配置用于JWT身份验证的密钥', '2025-12-04 15:20:19', '2025-12-04 15:20:19');
INSERT INTO `oak_plugins` VALUES ('5', 'pl-m5BzSXbCQfGzoQi', 'limit-req', 'icon-limit-req', '2', '使用漏桶算法限制客户端对服务的请求速率', '2025-12-04 15:20:19', '2025-12-04 15:20:19');
INSERT INTO `oak_plugins` VALUES ('6', 'pl-rLYsoeNVfPUMUAA', 'limit-conn', 'icon-limit-conn', '2', '限制客户端对服务的并发请求数', '2025-12-04 15:20:19', '2025-12-04 15:20:19');
INSERT INTO `oak_plugins` VALUES ('7', 'pl-XZxaqOgRZsBKpoE', 'limit-count', 'icon-limit-count', '2', '限制客户端在指定的时间范围内对服务的总请求数', '2025-12-04 15:20:19', '2025-12-04 15:20:19');

-- ----------------------------
-- Table structure for oak_plugin_configs
-- ----------------------------
DROP TABLE IF EXISTS `oak_plugin_configs`;
CREATE TABLE `oak_plugin_configs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Plugin config id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Plugin config name',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT 'Plugin relation type 1:service  2:router',
  `target_id` char(20) NOT NULL DEFAULT '' COMMENT 'Target id',
  `plugin_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Plugin res id',
  `plugin_key` varchar(20) NOT NULL DEFAULT '' COMMENT 'Plugin key',
  `config` text NOT NULL COMMENT 'Plugin configuration',
  `enable` tinyint unsigned NOT NULL DEFAULT '2' COMMENT 'Plugin config enable  1:on  2:off',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Plugin Configs';

-- ----------------------------
-- Records of oak_plugin_configs
-- ----------------------------
INSERT INTO `oak_plugin_configs` VALUES ('1', 'pc-cIvvY2IdZ1U9eQD', 'plugin-limit-req', '2', 'rt-Vany67MXLtgTHwG', 'pl-m5BzSXbCQfGzoQi', 'limit-req', '{\"rate\":1,\"burst\":500}', '1', '2025-12-04 15:56:32', '2025-12-04 15:56:32');

-- ----------------------------
-- Table structure for oak_routers
-- ----------------------------
DROP TABLE IF EXISTS `oak_routers`;
CREATE TABLE `oak_routers` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Router id',
  `service_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service id',
  `upstream_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Upstream id',
  `router_name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Router name',
  `request_methods` varchar(150) NOT NULL DEFAULT '' COMMENT 'Request method',
  `router_path` varchar(200) NOT NULL DEFAULT '' COMMENT 'Routing path',
  `enable` tinyint unsigned NOT NULL DEFAULT '2' COMMENT 'Router enable  1:on  2:off',
  `release` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'Service release status 1:unpublished  2:to be published  3:published',
  `client_max_body_size` bigint unsigned DEFAULT NULL COMMENT 'Maximum request body size in bytes',
  `chunked_transfer_encoding` tinyint unsigned DEFAULT NULL COMMENT 'Chunked transfer encoding 1:enable 2:disable',
  `proxy_buffering` tinyint unsigned DEFAULT NULL COMMENT 'Proxy buffering 1:enable 2:disable',
  `proxy_cache` text COMMENT 'Proxy cache configuration (JSON)',
  `proxy_set_header` text COMMENT 'Proxy set header configuration (JSON)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Routers';

-- ----------------------------
-- Records of oak_routers
-- ----------------------------
INSERT INTO `oak_routers` VALUES ('2', 'rt-cMQH880AhzqxQnO', 'sv-DQ0CASH03pMddn3', 'up-BsYcvp9VClgXIX2', 'sss', 'ALL', '/apple', '1', '3', '100', null, null, null, null, '2025-12-05 18:58:50', '2025-12-08 17:03:37');

-- ----------------------------
-- Table structure for oak_services
-- ----------------------------
DROP TABLE IF EXISTS `oak_services`;
CREATE TABLE `oak_services` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Service name',
  `protocol` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'Protocol  1:HTTP  2:HTTPS  3:HTTP&HTTPS',
  `enable` tinyint unsigned NOT NULL DEFAULT '2' COMMENT 'Service enable  1:on  2:off',
  `release` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'Service release status 1:unpublished  2:to be published  3:published',
  `client_max_body_size` bigint unsigned DEFAULT NULL COMMENT 'Maximum request body size in bytes',
  `chunked_transfer_encoding` tinyint unsigned DEFAULT NULL COMMENT 'Chunked transfer encoding 1:enable 2:disable',
  `proxy_buffering` tinyint unsigned DEFAULT NULL COMMENT 'Proxy buffering 1:enable 2:disable',
  `proxy_cache` text COMMENT 'Proxy cache configuration (JSON)',
  `proxy_set_header` text COMMENT 'Proxy set header configuration (JSON)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  KEY `IDX_NAME` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Services';

-- ----------------------------
-- Records of oak_services
-- ----------------------------
INSERT INTO `oak_services` VALUES ('3', 'sv-GdO0rkKOo5yJ9RM', 'aa', '1', '1', '2', null, null, null, null, null, '2025-12-05 19:22:32', '2025-12-08 17:07:08');
INSERT INTO `oak_services` VALUES ('4', 'sv-DQ0CASH03pMddn3', 'bb', '1', '1', '3', null, null, null, null, null, '2025-12-08 16:13:03', '2025-12-08 17:03:32');

-- ----------------------------
-- Table structure for oak_service_domains
-- ----------------------------
DROP TABLE IF EXISTS `oak_service_domains`;
CREATE TABLE `oak_service_domains` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Domain id',
  `service_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service id',
  `domain` varchar(50) NOT NULL DEFAULT '' COMMENT 'Domain name',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_SERVICE_ID_DOMAIN` (`service_res_id`,`domain`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Service domains';

-- ----------------------------
-- Records of oak_service_domains
-- ----------------------------
INSERT INTO `oak_service_domains` VALUES ('5', 'sd-aAREmDKmlalbqcU', 'sv-DQ0CASH03pMddn3', 'bb.apple.com', '2025-12-08 16:13:03', '2025-12-08 16:13:03');
INSERT INTO `oak_service_domains` VALUES ('6', 'sd-6lC97dnRG8juBeC', 'sv-GdO0rkKOo5yJ9RM', 'aa.apple.com', '2025-12-08 16:13:24', '2025-12-08 16:13:24');

-- ----------------------------
-- Table structure for oak_upstreams
-- ----------------------------
DROP TABLE IF EXISTS `oak_upstreams`;
CREATE TABLE `oak_upstreams` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Upstream id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'Upstream name',
  `algorithm` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Load balancing algorithm  1:round robin  2:chash',
  `connect_timeout` int unsigned NOT NULL DEFAULT '1' COMMENT 'Connect timeout',
  `write_timeout` int unsigned NOT NULL DEFAULT '1' COMMENT 'Write timeout',
  `read_timeout` int unsigned NOT NULL DEFAULT '1' COMMENT 'Read timeout',
  `enable` tinyint unsigned NOT NULL DEFAULT '2' COMMENT 'Enable  1:on  2:off',
  `release` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'Release status 1:unpublished  2:to be published  3:published',
  `check_enabled` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Health check enabled  0:false  1:true',
  `check_tcp` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'Health check type  0:http  1:tcp',
  `check_method` varchar(10) NOT NULL DEFAULT '' COMMENT 'HTTP method for health check',
  `check_host` varchar(150) NOT NULL DEFAULT '' COMMENT 'HTTP Host header for health check',
  `check_uri` varchar(150) NOT NULL DEFAULT '/' COMMENT 'HTTP URI for health check',
  `check_interval` int unsigned NOT NULL DEFAULT '1' COMMENT 'Health check interval in seconds',
  `check_timeout` int unsigned NOT NULL DEFAULT '1' COMMENT 'Health check timeout in seconds',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Upstreams';

-- ----------------------------
-- Records of oak_upstreams
-- ----------------------------
INSERT INTO `oak_upstreams` VALUES ('3', 'up-BsYcvp9VClgXIX2', 'test1', '1', '30000', '3000', '3000', '1', '3', '0', '1', '', '', '/', '1', '1', '2025-12-05 18:12:04', '2025-12-08 17:03:16');
INSERT INTO `oak_upstreams` VALUES ('4', 'up-m00a4aXgZul43zU', 'test2', '1', '30000', '3000', '3000', '1', '3', '0', '1', '', '', '/', '1', '1', '2025-12-05 18:31:39', '2025-12-08 15:57:47');

-- ----------------------------
-- Table structure for oak_upstream_nodes
-- ----------------------------
DROP TABLE IF EXISTS `oak_upstream_nodes`;
CREATE TABLE `oak_upstream_nodes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Service node id',
  `upstream_res_id` char(20) NOT NULL DEFAULT '' COMMENT 'Upstream id',
  `node_ip` varchar(60) NOT NULL DEFAULT '' COMMENT 'Node IP',
  `ip_type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'IP Type  1:IPV4  2:IPV6',
  `node_port` smallint unsigned NOT NULL DEFAULT '0' COMMENT 'Node port',
  `node_weight` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Node weight',
  `health` tinyint unsigned NOT NULL DEFAULT '1' COMMENT 'Health type  1:HEALTH  2:UNHEALTH',
  `health_check` tinyint unsigned NOT NULL DEFAULT '2' COMMENT 'Health check  1:on  2:off',
  `tags` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_UPSTREAM_ID_NODE_IP_PORT` (`upstream_res_id`,`node_ip`,`node_port`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Upstream nodes';

-- ----------------------------
-- Records of oak_upstream_nodes
-- ----------------------------
INSERT INTO `oak_upstream_nodes` VALUES ('4', 'un-RTjg5YiOqd2jFVH', 'up-BsYcvp9VClgXIX2', '127.0.0.2', '1', '8080', '100', '1', '2', null, '2025-12-05 18:19:07', '2025-12-08 15:58:03');
INSERT INTO `oak_upstream_nodes` VALUES ('8', 'un-rmxmDAvLlJt7dpe', 'up-m00a4aXgZul43zU', '127.0.0.1', '1', '8080', '100', '1', '2', '{\"123\":\"456\",\"789\":\"123\"}', '2025-12-05 18:50:36', '2025-12-08 15:57:30');
INSERT INTO `oak_upstream_nodes` VALUES ('9', 'un-tPdCAkBvo9xZOwz', 'up-m00a4aXgZul43zU', '127.0.0.2', '1', '8080', '100', '1', '2', '', '2025-12-08 15:56:53', '2025-12-08 15:57:30');
INSERT INTO `oak_upstream_nodes` VALUES ('10', 'un-PiXhYWPg81UqbNG', 'up-m00a4aXgZul43zU', '127.0.0.3', '1', '8080', '100', '1', '2', '', '2025-12-08 15:57:30', '2025-12-08 15:57:30');

-- ----------------------------
-- Table structure for oak_users
-- ----------------------------
DROP TABLE IF EXISTS `oak_users`;
CREATE TABLE `oak_users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'User iD',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT 'User name',
  `password` char(32) NOT NULL DEFAULT '' COMMENT 'Password',
  `email` varchar(80) NOT NULL DEFAULT '' COMMENT 'Email',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_EMAIL` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Users';

-- ----------------------------
-- Records of oak_users
-- ----------------------------
INSERT INTO `oak_users` VALUES ('1', 'us-Ew2h6VglDSz5Jgi', 'apple', '550e1bafe077ff0b0b67f4e32f29d751', 'apple@apple.com', '2025-12-04 15:25:01', '2025-12-04 15:25:01');

-- ----------------------------
-- Table structure for oak_user_tokens
-- ----------------------------
DROP TABLE IF EXISTS `oak_user_tokens`;
CREATE TABLE `oak_user_tokens` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_id` char(20) NOT NULL DEFAULT '' COMMENT 'User tokenID',
  `token` text NOT NULL COMMENT 'Token',
  `user_email` varchar(80) NOT NULL DEFAULT '' COMMENT 'Email',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Update time',
  `expired_at` timestamp NULL DEFAULT NULL COMMENT 'Expired time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_ID` (`res_id`),
  UNIQUE KEY `UNIQ_USER_EMAIL` (`user_email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='User token';

-- ----------------------------
-- Records of oak_user_tokens
-- ----------------------------
INSERT INTO `oak_user_tokens` VALUES ('2', 'ut-fRxOwcPbZkePkrB', 'eyJlbmNyeXB0aW9uIjoiYXBwbGVAYXBwbGUuY29tIiwidGltZXN0YW1wIjoiNzk3NmY5YTIzZWZmYjA5MDAzYzZiMzdkMGYwYzljZmYiLCJzZWNyZXQiOiJGUF9WbFdWM292TXY1SHNnUkJFRXIzbzgwWEttOXFlbkVjc1dQWGJPQW8wPSIsImlzc3VlciI6InphbmVoeSJ9', 'apple@apple.com', '2025-12-08 15:56:21', '2025-12-08 17:08:01', '2025-12-08 19:08:01');
