CREATE TABLE if not exists `resource` (
  `id` char(64)  NOT NULL COMMENT '资源的实例Id',
  `vendor` tinyint(1) NOT NULL,
  `region` varchar(64)  NOT NULL,
  `create_at` bigint NOT NULL,
  `expire_at` bigint NOT NULL,
  `type` varchar(120)  NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255)  NOT NULL,
  `status` varchar(255) NOT NULL,
  `update_at` bigint NOT NULL,
  `sync_at` bigint NOT NULL,
  `accout` varchar(255)  NOT NULL,
  `public_ip` varchar(64) NOT NULL,
  `private_ip` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `name` (`name`) USING BTREE,
  KEY `status` (`status`),
  KEY `private_ip` (`public_ip`) USING BTREE,
  KEY `public_ip` (`public_ip`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE if not exists `host` (
  `resource_id` varchar(64) NOT NULL,
  `cpu` tinyint NOT NULL,
  `memory` int NOT NULL,
  `gpu_amount` tinyint DEFAULT NULL,
  `gpu_spec` varchar(255) DEFAULT NULL,
  `os_type` varchar(255) DEFAULT NULL,
  `os_name` varchar(255) DEFAULT NULL,
  `serial_number` varchar(120) DEFAULT NULL,
  PRIMARY KEY (`resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;;