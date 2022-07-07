ALTER TABLE `t_contact_group` ADD `state` tinyint DEFAULT '1' COMMENT '状态 1启动 2删除';

CREATE TABLE `t_sync_time`
(
    `name` VARCHAR(256) NOT NULL PRIMARY KEY,
    `update_time` DATETIME
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

INSERT INTO t_sync_time (name,update_time,) VALUES ('contact','0000-01-01 00:00:00');
INSERT INTO t_sync_time (name,update_time) VALUES ('alarmRule','0000-01-01 00:00:00');
INSERT INTO t_sync_time (name,update_time) VALUES ('alarmRecord','0000-01-01 00:00:00');
