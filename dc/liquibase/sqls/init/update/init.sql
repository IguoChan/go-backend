DROP TABLE IF EXISTS `go_backend`.`student`;
CREATE TABLE `go_backend`.`student` (
    `id` bigint(4) unsigned NOT NULL AUTO_INCREMENT,
    `stu_id` varchar(64) NOT NULL,
    `first_name` varchar(64) NOT NULL,
    `last_name` varchar(64) NOT NULL,
    `email` varchar(128) NOT NULL,
    `phone_number` varchar(32) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_stu_id` (`stu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;