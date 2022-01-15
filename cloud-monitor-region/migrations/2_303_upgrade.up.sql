-- 重命名 config_item
RENAME TABLE `config_item` TO `t_config_item`;
-- 重命名 monitor_product
RENAME TABLE `monitor_product` TO `t_monitor_product`;
-- 重命名 monitor_item
RENAME TABLE `monitor_item` TO `t_monitor_item`;
-- 重命名 notification_record
RENAME TABLE `notification_record` TO `t_notification_record`;
-- 重命名 t_alert_record
RENAME TABLE `t_alert_record` TO `t_alarm_record`;

-- t_config_item 改造
ALTER TABLE `t_config_item` DROP PRIMARY KEY;
ALTER TABLE `t_config_item` CHANGE `id` `biz_id` VARCHAR (50) not null default '-1' comment '业务Id';
ALTER TABLE `t_config_item` CHANGE `pid` `p_biz_id` VARCHAR (50)  not null default '-1' comment '上级业务Id';
ALTER TABLE `t_config_item` ADD `id` BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;

-- t_monitor_product
ALTER TABLE `t_monitor_product` CHANGE `id` `id` INT;
ALTER TABLE `t_monitor_product` DROP PRIMARY KEY;
ALTER TABLE `t_monitor_product` CHANGE `id` `biz_id` VARCHAR (50) not null default '-1' comment '业务Id';
ALTER TABLE `t_monitor_product` ADD `id` BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;

-- t_monitor_item
ALTER TABLE `t_monitor_item` CHANGE `id` `id` INT;
ALTER TABLE `t_monitor_item` DROP PRIMARY KEY;
ALTER TABLE `t_monitor_item` CHANGE `id` `biz_id` VARCHAR (50) not null default '-1' comment '业务Id';
ALTER TABLE `t_monitor_item` ADD `id` BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;
ALTER TABLE `t_monitor_item` CHANGE `product_id` `product_biz_id` VARCHAR (50) not null default '-1' comment '产品业务Id';

-- t_notification_record
ALTER TABLE `t_notification_record` DROP PRIMARY KEY;
ALTER TABLE `t_notification_record` CHANGE `id` `biz_id` VARCHAR (50) not null default '-1' comment '业务Id';
ALTER TABLE `t_notification_record` ADD `id` BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;
ALTER TABLE `t_notification_record` MODIFY `sender_id` VARCHAR (50) NOT NULL default '-1' comment '发送人';
ALTER TABLE `t_notification_record` MODIFY `create_time` datetime NOT NULL comment '创建时间';
-- 索引
ALTER TABLE `t_notification_record` ADD INDEX `tenant_time_idx` (`sender_id`,`create_time`);

-- 拆分 t_alarm_record 表
CREATE TABLE `t_alarm_info` (
                                `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
                                `alarm_biz_id` varchar(50) NOT NULL COMMENT '告警业务Id',
                                `summary` text COMMENT '告警额外信息',
                                `expression` text COMMENT '告警表达式',
                                `contact_info` text COMMENT '告警联系人',
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC  COMMENT='告警详情';
ALTER TABLE `t_alarm_info` ADD INDEX `alarm_biz_id_idx` (`alarm_biz_id`);
-- 迁移数据
INSERT INTO `t_alarm_info`(`alarm_biz_id`, `summary`, `contact_info`, `expression`)
SELECT `id`, `summary`, `contact_info`, `expression`  FROM `t_alarm_record` ;

-- t_alarm_record
ALTER TABLE `t_alarm_record` DROP PRIMARY KEY;
ALTER TABLE `t_alarm_record` CHANGE `id` `biz_id` VARCHAR (50) not null default '-1' comment '业务Id';
ALTER TABLE `t_alarm_record` MODIFY `tenant_id` VARCHAR (50) NOT NULL default '-1' comment '租户Id';
ALTER TABLE `t_alarm_record` ADD `id` BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;
-- 新增字段 存储过程
CREATE PROCEDURE `add_alarm_record_as_columns`()-- 新增一个存储过程
BEGIN
IF NOT EXISTS (SELECT column_name FROM information_schema.columns WHERE table_name = 't_alarm_record' and column_name = 'request_id')
THEN
ALTER TABLE `t_alarm_record` ADD `request_id` varchar(50) comment '请求Id';
ALTER TABLE `t_alarm_record` ADD `rule_source_type` tinyint not null default 1 comment '告警规则来源类型';
ALTER TABLE `t_alarm_record` ADD `rule_source_id` varchar(100) comment '告警规则来源业务Id';
END IF;
END;
-- 运行该存储过程
call add_alarm_record_as_columns();
-- 删除该存储过程
drop PROCEDURE if exists add_alarm_record_as_columns;


-- 删除字段
ALTER TABLE `t_alarm_record` drop COLUMN `summary`;
ALTER TABLE `t_alarm_record` drop COLUMN `contact_info`;
ALTER TABLE `t_alarm_record` drop COLUMN `expression`;
ALTER TABLE `t_alarm_record` drop COLUMN `notice_status`;
-- 索引

CREATE PROCEDURE `drop_index_if_exists`(IN i_table_name VARCHAR(128), IN i_index_name VARCHAR(128))
BEGIN
	DECLARE v_db varchar(100);
SELECT DATABASE() INTO v_db;

SET @tableName = i_table_name;
    SET @indexName = i_index_name;
    SET @indexExists = 0;

SELECT
    1
INTO @indexExists FROM
    INFORMATION_SCHEMA.STATISTICS
WHERE TABLE_NAME = @tableName AND INDEX_NAME = @indexName and table_schema = v_db;

SET @query = CONCAT('DROP INDEX ', @indexName, ' ON ', @tableName,' ' );
    IF @indexExists THEN
     PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
END IF;
END;

CALL drop_index_if_exists('t_alarm_record','rule_name_idx');
CALL drop_index_if_exists('t_alarm_record','resource_Id_idx');

drop PROCEDURE if exists drop_index_if_exists;


ALTER TABLE `t_alarm_record` ADD INDEX `tenant_time_idx` (`tenant_id`,`create_time`);
ALTER TABLE `t_alarm_record` ADD INDEX `biz_id_idx` (`biz_id`);



-- 重命名 alert_contact
RENAME TABLE alert_contact TO t_contact;
-- 重命名 alert_contact_group
RENAME TABLE alert_contact_group TO t_contact_group;
-- 重命名 alert_contact_group_rel
RENAME TABLE alert_contact_group_rel TO t_contact_group_rel;
-- 重命名 alert_contact_information
RENAME TABLE alert_contact_information TO t_contact_information;

-- t_contact 改造
ALTER TABLE t_contact DROP PRIMARY KEY;
ALTER TABLE t_contact CHANGE id biz_id VARCHAR (255) NOT NULL COMMENT '联系人ID';
ALTER TABLE t_contact ADD UNIQUE ( biz_id );
ALTER TABLE t_contact ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY COMMENT '主键ID';
ALTER TABLE t_contact CHANGE status state tinyint DEFAULT '1' COMMENT '状态';

-- t_contact_group 改造
ALTER TABLE t_contact_group DROP PRIMARY KEY;
ALTER TABLE t_contact_group CHANGE id biz_id VARCHAR (255) NOT NULL COMMENT '联系组ID';
ALTER TABLE t_contact_group ADD UNIQUE ( biz_id );
ALTER TABLE t_contact_group ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY COMMENT '主键ID';

-- t_contact_group_rel 改造
ALTER TABLE t_contact_group_rel DROP COLUMN id;
ALTER TABLE t_contact_group_rel ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY COMMENT '主键ID';
ALTER TABLE t_contact_group_rel CHANGE contact_id contact_biz_id VARCHAR (255) NOT NULL COMMENT '联系人ID';
ALTER TABLE t_contact_group_rel CHANGE group_id group_biz_id VARCHAR (255) NOT NULL COMMENT '联系组ID';

-- t_contact_information 改造
ALTER TABLE t_contact_information DROP COLUMN id;
ALTER TABLE t_contact_information ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY COMMENT '主键ID';
ALTER TABLE t_contact_information CHANGE contact_id contact_biz_id VARCHAR (255) NOT NULL COMMENT '联系人ID';
ALTER TABLE t_contact_information CHANGE no address VARCHAR (255) COMMENT '联系方式';
ALTER TABLE t_contact_information CHANGE is_certify state tinyint DEFAULT '0' COMMENT '状态';




CREATE TABLE IF NOT EXISTS `t_resource_group` (
    `id` varchar(256)  NOT NULL COMMENT '主键id',
    `name` varchar(256)  DEFAULT NULL COMMENT '组名称',
    `tenant_id` varchar(256)  DEFAULT NULL COMMENT '租户id',
    `product_id` varchar(256)  DEFAULT NULL COMMENT '产品id',
    `source` tinyint DEFAULT 1 COMMENT '来源:1 用户 2 弹性伸缩',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci COMMENT = '资源组'
    ROW_FORMAT = DYNAMIC;

CREATE TABLE IF NOT EXISTS `t_resource_resource_group_rel` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT  COMMENT '主键id',
    `resource_group_id` varchar(256)  DEFAULT NULL  COMMENT '资源组id',
    `resource_id` varchar(256)  DEFAULT NULL  COMMENT '资源id',
    `tenant_id` varchar(256)  DEFAULT NULL COMMENT '租户id',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci COMMENT = '资源与组关系'
    ROW_FORMAT = DYNAMIC;

CREATE TABLE IF NOT EXISTS `t_alarm_rule_group_rel` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `alarm_rule_id` varchar(256)  DEFAULT NULL COMMENT '规则id',
    `resource_group_id` varchar(256)  DEFAULT NULL COMMENT '资源组id',
    `calc_mode` tinyint DEFAULT 1 COMMENT '聚合方式 ,1:按单实例，2按组聚合 ',
    `tenant_id` varchar(256) DEFAULT NULL COMMENT '租户id',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime DEFAULT NULL  COMMENT '更新时间',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci COMMENT = '规则与组关系'
    ROW_FORMAT = DYNAMIC;


CREATE TABLE IF NOT EXISTS `t_alarm_rule_resource_rel` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
     `alarm_rule_id` varchar(256)  DEFAULT NULL,
    `resource_id` varchar(256)  DEFAULT NULL,
    `tenant_id` varchar(256)  DEFAULT NULL,
    `create_time` datetime DEFAULT NULL,
    `update_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci COMMENT = '规则与资源关系'
    ROW_FORMAT = DYNAMIC;

CREATE TABLE IF NOT EXISTS `t_alarm_handler` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `alarm_rule_id` varchar(256)  DEFAULT NULL,
    `handle_type` tinyint DEFAULT NULL COMMENT '1邮件；2 短信；3弹性',
    `handle_params` varchar(256)  DEFAULT NULL COMMENT '回调地址',
    `tenant_id` varchar(256)  DEFAULT NULL,
    `create_time` datetime DEFAULT NULL,
    `update_time` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB
    CHARACTER SET = utf8mb4
    COLLATE = utf8mb4_unicode_ci COMMENT = '规则与处理关系'
    ROW_FORMAT = DYNAMIC;

-- 表字符集及排序修改
ALTER TABLE t_alarm_handler CONVERT TO CHARACTER SET  utf8mb4 COLLATE  utf8mb4_unicode_ci;
ALTER TABLE t_alarm_rule CONVERT TO CHARACTER SET  utf8mb4 COLLATE  utf8mb4_unicode_ci;
ALTER TABLE t_resource_group CONVERT TO CHARACTER SET  utf8mb4 COLLATE  utf8mb4_unicode_ci;
ALTER TABLE t_alarm_rule_group_rel CONVERT TO CHARACTER SET  utf8mb4 COLLATE  utf8mb4_unicode_ci;
ALTER TABLE t_alarm_rule_resource_rel CONVERT TO CHARACTER SET  utf8mb4 COLLATE  utf8mb4_unicode_ci;
ALTER TABLE t_resource_resource_group_rel CONVERT TO CHARACTER SET  utf8mb4 COLLATE  utf8mb4_unicode_ci;

-- 4实例增加produt_type
-- 新增字段 存储过程
CREATE PROCEDURE  `add_alarm_instance_product_type`()
BEGIN
    IF NOT EXISTS (SELECT column_name FROM information_schema.columns WHERE table_name = 't_alarm_instance' and column_name = 'product_type' and TABLE_SCHEMA=DATABASE())
    THEN
ALTER TABLE `t_alarm_instance` ADD COLUMN `product_type` VARCHAR (256);
ALTER TABLE `t_alarm_rule` ADD COLUMN `source` varchar(256);
ALTER TABLE `t_alarm_rule` ADD COLUMN `source_type` tinyint unsigned DEFAULT '1';
UPDATE t_alarm_instance t1    LEFT JOIN t_alarm_rule t2 ON t1.alarm_rule_id = t2.id     SET t1.product_type = t2.product_type;
-- 1 notify channel  拆解为t_alarm_handler
INSERT INTO t_alarm_handler (
    `alarm_rule_id`,
    `handle_type`,
    `handle_params`,
    `tenant_id`,
    `create_time`,
    update_time
)
SELECT
    id,
    CASE notify_channel
        WHEN 2 THEN
            1
        ELSE
            2
        END AS handle_type,
    '',
    tenant_id,
    create_time,
    create_time
FROM
    t_alarm_rule
WHERE
    notify_channel != 3;

-- 2规则与实例的关系拆解
INSERT INTO t_alarm_rule_resource_rel (
    `alarm_rule_id`,
    `resource_id`,
    `tenant_id`,
    `create_time`,
    `update_time`
) SELECT
      t1.id,
      t2.instance_id,
      t1.tenant_id,
      t2.create_time,
      t2.create_time
FROM
    t_alarm_rule t1,
    t_alarm_instance t2
WHERE
        t1.id = t2.alarm_rule_id;

END IF;
END;
-- 运行该存储过程
call add_alarm_instance_product_type();
-- 删除该存储过程
drop PROCEDURE if exists add_alarm_instance_product_type;


CREATE PROCEDURE `drop_instance_key`()-- 新增一个存储过程
BEGIN
    IF  EXISTS (select CONSTRAINT_NAME  from INFORMATION_SCHEMA.KEY_COLUMN_USAGE t   where t.TABLE_NAME ='t_alarm_instance' and t.CONSTRAINT_SCHEMA=DATABASE())
    THEN
ALTER TABLE `t_alarm_instance` DROP PRIMARY KEY;
END IF;
END;
-- 运行该存储过程
call drop_instance_key();
-- 删除该存储过程
drop PROCEDURE if exists drop_instance_key;




-- t_alarm_notice
ALTER TABLE t_alarm_notice ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;


-- t_alarm_instance
RENAME TABLE t_alarm_instance TO t_resource;
ALTER TABLE `t_resource`  ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY ;
ALTER TABLE `t_resource`  ADD COLUMN `product_biz_id`  varchar(50)  COMMENT '产品业务id';
ALTER TABLE `t_resource` CHANGE COLUMN `product_type` `product_name`  varchar(256)  COMMENT '产品名称';
ALTER TABLE `t_resource`  CHANGE COLUMN `instance_id` `instance_id`  varchar(64)  NOT NULL  COMMENT '实例id';

-- t_alarm_rule
ALTER TABLE `t_alarm_rule`    ADD COLUMN `product_biz_id`  varchar(50) ;
ALTER TABLE `t_alarm_rule`  CHANGE COLUMN `product_type` `product_name`  varchar(256)  COMMENT '产品名称';
ALTER TABLE `t_alarm_rule` DROP PRIMARY KEY;
ALTER TABLE `t_alarm_rule` CHANGE COLUMN `id` `biz_id`  varchar(256) NOT NULL ;
ALTER TABLE `t_alarm_rule` ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;
ALTER TABLE `t_alarm_rule` CHANGE COLUMN `tenant_id` `tenant_id`  varchar(64) NOT NULL ;
-- t_resource_group
ALTER TABLE `t_resource_group` DROP PRIMARY KEY;
ALTER TABLE `t_resource_group` CHANGE COLUMN `id` `biz_id`  varchar(256) NOT NULL  COMMENT '业务id';
ALTER TABLE `t_resource_group` ADD id BIGINT UNSIGNED NOT NULL auto_increment PRIMARY KEY;
ALTER TABLE `t_resource_group` CHANGE COLUMN `product_id` `product_biz_id`  varchar(256) COMMENT '产品业务id';

drop table if EXISTS t_user_prometheus_id;

-- 3去重 删除多余实例
DELETE
FROM
    t_resource
WHERE
        instance_id IN (
        SELECT
            *
        FROM
            (
                SELECT
                    instance_id
                FROM
                    t_resource
                GROUP BY
                    instance_id
                HAVING
                        count(instance_id) > 1
            ) t2
    )
  AND create_time NOT IN (
    SELECT
        *
    FROM
        (
            SELECT
                min(create_time)
            FROM
                t_resource
            GROUP BY
                instance_id
            HAVING
                    count(instance_id) > 1
        ) t
);
-- 增加索引
ALTER TABLE `t_resource` ADD UNIQUE INDEX `idx_instance_id` (`instance_id`);
ALTER TABLE `t_alarm_rule` ADD UNIQUE INDEX `idx_biz_id` (`biz_id`);
ALTER TABLE `t_resource_group` ADD UNIQUE INDEX `idx_biz_id` (`biz_id`);
ALTER TABLE `t_alarm_rule` ADD INDEX  `idx_tenant_id` (`tenant_id`);

-- 4删除字段
CREATE PROCEDURE `drop_column_if_exists`(IN i_table_name VARCHAR(128), IN i_column_name VARCHAR(128))
BEGIN
	DECLARE v_db varchar(100);
SELECT DATABASE() INTO v_db;

SET @tableName = i_table_name;
    SET @icolumnName = i_column_name;
    SET @columnExists = 0;

SELECT
    1
INTO @columnExists FROM
    information_schema.columns
WHERE TABLE_NAME = @tableName AND column_name = @icolumnName and table_schema = v_db;

SET @query = CONCAT('ALTER TABLE ', @tableName, ' drop COLUMN ', @icolumnName,' ' );
    IF @columnExists THEN
     PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
END IF;
END;
CALL drop_column_if_exists('t_alarm_rule','notify_channel');
CALL drop_column_if_exists('t_resource','alarm_rule_id');
-- 删除该存储过程
drop PROCEDURE if exists drop_column_if_exists;

