DROP TABLE IF EXISTS sample;
CREATE TABLE sample (
 `id` bigint NOT NULL AUTO_INCREMENT,
 `gmt_create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `gmt_modified` timestamp NOT NULL on update CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 `r_int` int NOT NULL,
 `n_int` int NULL DEFAULT NULL,
 `r_float` float NOT NULL,
 `n_float` float NULL DEFAULT NULL,
 `r_string` varchar(10) NOT NULL,
 `n_string` varchar(10) NULL DEFAULT NULL,
 `r_time` datetime NOT NULL,
 `n_time` datetime NULL DEFAULT NULL,
 `union1` varchar(10) NOT NULL,
 `union2` varchar(10) NOT NULL,
 `union3` varchar(10) NULL DEFAULT NULL,
 PRIMARY KEY (`id`),
 UNIQUE KEY `uk_n_int` (`n_int`),
 UNIQUE KEY `uk_r_int` (`r_int`),
 UNIQUE KEY `uk_union_1_union_2` (`union1`,`union2`),
 KEY `idx_n_time` (`n_time`),
 KEY `idx_r_time` (`r_time`),
 KEY `idx_union1_union3` (`union1`,`union3`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
