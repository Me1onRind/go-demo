CREATE TABLE `user_tab` (
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT(20) UNSIGNED NOT NULL,
    `name` VARCHAR(16) NOT NULL DEFAULT '',
    `email` VARCHAR(128) NOT NULL DEFAULT '',
    `create_time` INT(10) UNSIGNED NOT NULL DEFAULT 0,
    `update_time` INT(10) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY(`id`),
    UNIQUE KEY uniq_user_id(`user_id`),
    UNIQUE KEY uniq_email(`email`)
) ENGINE=InnoDB CHARSET=utf8mb4;
