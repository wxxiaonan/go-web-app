CREATE TABLE `user` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `user_id` bigint(20) NOT NULL,
                        `username` varchar(64) CHARACTER SET utf8mb4 NOT NULL,
                        `password` varchar(64) CHARACTER SET utf8mb4 NOT NULL,
                        `email` varchar(64) CHARACTER SET utf8mb4 DEFAULT NULL,
                        `gender` varchar(64) COLLATE utf8mb4_bin NOT NULL DEFAULT '0',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        PRIMARY KEY (`id`) USING BTREE,
                        UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE,
                        UNIQUE KEY `idx_username` (`username`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;