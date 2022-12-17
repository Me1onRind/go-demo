CREATE TABLE `id_creator_tab` (
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `id_type` INT(10) UNSIGNED NOT NULL,
    `offset` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
    `step` INT(10) UNSIGNED NOT NULL,
    `create_time` INT(10) UNSIGNED NOT NULL DEFAULT 0,
    `update_time` INT(10) UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY(`id`),
    UNIQUE KEY uniq_user_id(`id_type`)
) ENGINE=InnoDB CHARSET=utf8mb4;
