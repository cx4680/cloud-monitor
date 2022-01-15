
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for monitor_product
-- ----------------------------
DROP TABLE IF EXISTS `monitor_product`;
CREATE TABLE `monitor_product` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) DEFAULT NULL,
  `status` tinyint(3) unsigned DEFAULT NULL,
  `description` varchar(256) DEFAULT NULL,
  `create_user` varchar(256) DEFAULT NULL,
  `create_time` varchar(256) DEFAULT NULL,
  `route` varchar(256) DEFAULT NULL,
  `cron` varchar(256) DEFAULT NULL,
  `host` varchar(500) DEFAULT NULL,
  `page_url` varchar(500) DEFAULT NULL,
  `abbreviation` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for monitor_item
-- ----------------------------
DROP TABLE IF EXISTS `monitor_item`;
CREATE TABLE `monitor_item` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` varchar(256) DEFAULT NULL,
  `name` varchar(256) DEFAULT NULL,
  `metric_name` varchar(256) DEFAULT NULL,
  `labels` varchar(256) DEFAULT NULL,
  `metrics_linux` varchar(256) DEFAULT NULL,
  `metrics_windows` varchar(256) DEFAULT NULL,
  `statistics` varchar(256) DEFAULT NULL,
  `unit` varchar(256) DEFAULT NULL,
  `frequency` varchar(256) DEFAULT NULL,
  `type` tinyint(3) unsigned DEFAULT NULL,
  `is_display` tinyint(3) unsigned DEFAULT NULL,
  `status` tinyint(3) unsigned DEFAULT NULL,
  `description` varchar(256) DEFAULT NULL,
  `create_user` varchar(256) DEFAULT NULL,
  `create_time` datetime DEFAULT NULL,
  `show_expression` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;


-- ----------------------------
-- Table structure for config_item
-- ----------------------------
DROP TABLE IF EXISTS `config_item`;
CREATE TABLE `config_item`
(
    `id`      varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NOT NULL,
    `pid`     varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL DEFAULT '-1' COMMENT '上级Id',
    `name`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置名称',
    `code`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置编码',
    `data`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置值',
    `sort_id` int(11)                                                       NULL DEFAULT 0 COMMENT '排序',
    `remark`  varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for notification_record
-- ----------------------------
CREATE TABLE IF NOT EXISTS `notification_record`
(
    `id`                varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NOT NULL,
    `sender_id`         varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL DEFAULT NULL COMMENT '发送人',
    `source_id`         varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL DEFAULT NULL COMMENT '通知源',
    `source_type`       int(11)                                                       NULL DEFAULT NULL COMMENT '源类型',
    `target_address`    varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '通知对象地址',
    `notification_type` int(255)                                                      NULL DEFAULT NULL COMMENT '通知类型，1：短信，2：邮箱',
    `result`            tinyint(4)                                                    NULL DEFAULT NULL COMMENT '通知结果 0:失败1:成功',
    `create_time`       datetime                                                      NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT = '通知记录'
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for alert_contact
-- ----------------------------
CREATE TABLE IF NOT EXISTS `alert_contact`
(
    `id`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编号',
    `tenant_id`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租户ID',
    `name`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '姓名',
    `status`      int(1) UNSIGNED ZEROFILL                                      NULL DEFAULT NULL COMMENT '状态 0:停用 1:正常',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '描述',
    `create_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人',
    `create_time` datetime                                                      NULL DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime                                                      NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for alert_contact_group
-- ----------------------------
CREATE TABLE IF NOT EXISTS `alert_contact_group`
(
    `id`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编号',
    `tenant_id`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '租户ID',
    `name`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '组名',
    `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '描述',
    `create_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人',
    `create_time` datetime                                                      NULL DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime                                                      NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for alert_contact_group_rel
-- ----------------------------
CREATE TABLE IF NOT EXISTS `alert_contact_group_rel`
(
    `id`          bigint(255)                                                   NOT NULL COMMENT '编号',
    `tenant_id`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '租户ID',
    `contact_id`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '联系人编号',
    `group_id`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '群组编号',
    `create_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人',
    `create_time` datetime                                                      NULL DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime                                                      NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for alert_contact_information
-- ----------------------------
CREATE TABLE IF NOT EXISTS `alert_contact_information`
(
    `id`          bigint(255)                                                   NOT NULL AUTO_INCREMENT COMMENT '编号',
    `tenant_id`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '租户ID',
    `contact_id`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '联系人编号',
    `no`          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '号码',
    `type`        int(1)                                                        NULL DEFAULT NULL COMMENT '1:电话 2:邮箱 3:蓝信',
    `is_certify`  int(1)                                                        NULL DEFAULT NULL COMMENT '是否认证 0:未认证 1:已认证',
    `active_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '激活码',
    `create_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人',
    `create_time` datetime                                                      NULL DEFAULT NULL COMMENT '创建时间',
    `update_time` datetime                                                      NULL DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for t_alarm_instance
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_alarm_instance`
(
    `alarm_rule_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL DEFAULT NULL,
    `instance_id`   varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL DEFAULT NULL,
    `create_time`   datetime                                                      NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `region_code`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    `zone_code`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    `ip`            varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    `region_name`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    `zone_name`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    `instance_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    `tenant_id`     varchar(64) CHARACTER SET utf8 COLLATE utf8_bin               NULL DEFAULT NULL COMMENT '租户id'
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for t_alarm_notice
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_alarm_notice`
(
    `alarm_rule_id`     varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '规则id',
    `contract_group_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '联系组id',
    `create_time`       datetime                                                     NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for t_alarm_rule
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_alarm_rule`
(
    `id`                varchar(32) CHARACTER SET utf8 COLLATE utf8_bin  NOT NULL,
    `monitor_type`      varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '监控类型',
    `product_type`      varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '所属产品',
    `dimensions`        int(11)                                          NULL DEFAULT NULL COMMENT '维度（1 全部资源 2 实例 ）',
    `name`              varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '规则名称',
    `metric_name`       varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL,
    `trigger_condition` json                                             NULL COMMENT '触发条件  {\r\n\"metricName\":\"cpu.util\",\r\n\"period\":60,\r\n  \"comparisonOperator\":\">=\",\r\n  \"times\":1,\r\n   \"statistics\":\"avg\",\r\n   \"threshold\":\"11\"\r\n}\r\n',
    `silences_time`     varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '冷却周期',
    `effective_start`   varchar(19) CHARACTER SET utf8 COLLATE utf8_bin  NULL DEFAULT NULL COMMENT '监控时间段-开始时间',
    `effective_end`     varchar(19) CHARACTER SET utf8 COLLATE utf8_bin  NULL DEFAULT NULL COMMENT '监控时间段-结束时间',
    `level`             int(11)                                          NULL DEFAULT NULL COMMENT '报警级别  \r\n紧急1 \r\n重要2 \r\n次要3\r\n提醒 4\r\n',
    `notify_channel`    int(11)                                          NULL DEFAULT NULL COMMENT '通知方式 1 all  2 email  3 sms ',
    `enabled`           int(11)                                          NULL DEFAULT 1 COMMENT '启用（1）禁用（0）',
    `tenant_id`         varchar(64) CHARACTER SET utf8 COLLATE utf8_bin  NULL DEFAULT NULL COMMENT '租户id',
    `create_time`       datetime                                         NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `create_user`       varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '创建人',
    `deleted`           int(11)                                          NULL DEFAULT 0 COMMENT '未删除0  删除1 ',
    `user_name`         varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NULL DEFAULT NULL COMMENT '租户名称',
    `update_time`       datetime                                         NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8
  COLLATE = utf8_bin
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for t_alert_record
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_alert_record`
(
    `id`            varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NOT NULL,
    `status`        varchar(10) CHARACTER SET utf8 COLLATE utf8_general_ci        NOT NULL DEFAULT '' COMMENT '告警状态 firing resolved',
    `tenant_id`     varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL     DEFAULT NULL,
    `rule_id`       varchar(32) CHARACTER SET utf8 COLLATE utf8_bin               NOT NULL DEFAULT '' COMMENT '规则id',
    `rule_name`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `monitor_type`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL     DEFAULT NULL COMMENT '监控类型',
    `source_type`   varchar(10) CHARACTER SET utf8 COLLATE utf8_bin               NOT NULL COMMENT '资源类型 ',
    `source_id`     varchar(32) CHARACTER SET utf8 COLLATE utf8_bin               NOT NULL COMMENT '资源id',
    `summary`       text CHARACTER SET utf8 COLLATE utf8_general_ci               NULL,
    `current_value` varchar(50) CHARACTER SET utf8 COLLATE utf8_bin               NOT NULL COMMENT '当前值',
    `start_time`    datetime                                                      NOT NULL COMMENT '告警开始时间',
    `end_time`      datetime                                                      NULL     DEFAULT NULL COMMENT '告警结束时间',
    `target_value`  varchar(50) CHARACTER SET utf8 COLLATE utf8_bin               NOT NULL COMMENT '规则定义阈值',
    `expression`    varchar(500) CHARACTER SET utf8 COLLATE utf8_bin              NOT NULL COMMENT '计算公式',
    `duration`      varchar(100) CHARACTER SET utf8 COLLATE utf8_bin              NOT NULL COMMENT '持续时间',
    `level`         int(11)                                                       NOT NULL COMMENT '告警级别 紧急 重要 次要 提醒',
    `notice_status` varchar(10) CHARACTER SET utf8 COLLATE utf8_bin               NOT NULL COMMENT '消息发送状态, error 发送失败 success 成功',
    `alarm_key`     varchar(100) CHARACTER SET utf8 COLLATE utf8_bin              NOT NULL COMMENT '告警项',
    `contact_info`  text CHARACTER SET utf8 COLLATE utf8_bin                      NULL COMMENT '联系方式',
    `region`        varchar(50) CHARACTER SET utf8 COLLATE utf8_bin               NULL     DEFAULT NULL,
    `create_time`   datetime                                                      NOT NULL,
    `update_time`   datetime                                                      NULL     DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    KEY `rule_name_idx` (`rule_name`) USING BTREE,
    KEY `resource_Id_idx` (`source_id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT = '告警记录'
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for t_user_prometheus_id
-- ----------------------------
CREATE TABLE IF NOT EXISTS `t_user_prometheus_id`
(
    `tenant_id`          varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NULL DEFAULT NULL,
    `prometheus_rule_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NOT NULL,
    `create_time`        datetime                                                      NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_user`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
    `deleted`            int(1)                                                        NULL DEFAULT 0 COMMENT '1（已删除） 0未删除',
    UNIQUE INDEX `idx` (`tenant_id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;



-- ----------------------------
-- Records of monitor_product
-- ----------------------------
INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('1', '云服务器ECS', '1', 'ecs', null, null, '/productmonitoring/ecs', '0 0 0/1 * * ?', 'http://product-ecs-ecs-manage.product-ecs', '/noauth/ecs/PageList', 'ecs');
INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('2', '弹性公网IP', '1', 'eip', null, null, '/productmonitoring/eip', '0 0 0/1 * * ?', 'http://product-eip-eip-manage.product-eip', '/eip/inner/eipInfoList', 'eip');
INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('3', '负载均衡SLB', '1', 'slb', null, null, '/productmonitoring/slb', '0 0 0/1 * * ?', 'http://product-slb-controller-slb-manage.product-slb', '/slb/inner/list', 'slb');

-- ----------------------------
-- Records of monitor_item
-- ----------------------------
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('1', '1', 'CPU使用率', 'ecs_cpu_usage', 'instance', '100 - (100 * (sum by(instance) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('2', '1', 'CPU1分钟平均负载', 'ecs_load1', 'instance', 'ecs_load1{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"windows\"}}");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('3', '1', 'CPU5分钟平均负载', 'ecs_load5', 'instance', 'ecs_load5{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"windows\"}}");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('4', '1', 'CPU15分钟平均负载', 'ecs_load15', 'instance', 'ecs_load15{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"windows\"}}");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('5', '1', '内存使用量', 'ecs_memory_used', 'instance', 'ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes', null, null, 'Byte', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('6', '1', '内存使用率', 'ecs_memory_usage', 'instance', '100 * ((ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes) / ecs_memory_MemTotal_bytes)', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('7', '1', '磁盘使用率', 'ecs_disk_usage', 'instance,device', '100 * ((ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes) / ecs_filesystem_size_bytes)', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('8', '1', '磁盘读速率', 'ecs_disk_read_rate', 'instance,device', 'irate(ecs_disk_read_bytes_total{$INSTANCE}[3m])', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('9', '1', '磁盘写速率', 'ecs_disk_write_rate', 'instance,device', 'irate(ecs_disk_written_bytes_total{$INSTANCE}[3m])', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('10', '1', '磁盘读IOPS', 'ecs_disk_read_iops', 'instance,device', 'irate(ecs_disk_reads_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('11', '1', '磁盘写IOPS', 'ecs_disk_write_iops', 'instance,device', 'irate(ecs_disk_writes_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('12', '1', '流入带宽', 'ecs_network_receive_rate', 'instance,device', 'irate(ecs_network_receive_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, 'Mbps', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('13', '1', '流出带宽', 'ecs_network_transmit_rate', 'instance,device', 'irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, 'Mbps', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('14', '1', '包接收速率', 'ecs_network_receive_packets_rate', 'instance,device', 'irate(ecs_network_receive_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('15', '1', '包发送速率', 'ecs_network_transmit_packets_rate', 'instance,device', 'irate(ecs_network_transmit_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('16', '1', '(基础)CPU的平均使用率', 'ecs_cpu_base_usage', 'instance', '100 * avg(irate(ecs_base_vcpu_seconds{$INSTANCE}[6m]))', null, null, '%', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('17', '1', '(基础)磁盘读速率', 'ecs_disk_base_read_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type="read",$INSTANCE}[6m])', null, null, 'Byte/s', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('18', '1', '(基础)磁盘写速率', 'ecs_disk_base_write_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type="write",$INSTANCE}[6m])', null, null, 'Byte/s', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('19', '1', '(基础)网卡下行带宽', 'ecs_network_base_receive_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type="rx",$INSTANCE}[6m])', null, null, 'Byte/s', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('20', '1', '(基础)网卡上行带宽', 'ecs_network_base_transmit_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type="tx",$INSTANCE}[6m])', null, null, 'Byte/s', null, '1', '1', '1', null, null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('21', '2', '出网带宽', 'eip_upstream_bandwidth', 'instance', 'sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,eip)', null, null, 'bps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('22', '2', '入网带宽', 'eip_downstream_bandwidth', 'instance', 'sum(eip_downstream_bits_rate{$INSTANCE}) by (instance,eip)', null, null, 'bps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('23', '2', '出网流量', 'eip_upstream', 'instance', '((sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,eip))/8)*60', null, null, 'Byte', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('24', '2', '入网流量', 'eip_downstream', 'instance', '((sum(eip_downstream_bits_rate{$INSTANCE}) by (instance,eip))/8)*60', null, null, 'Byte', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('25', '2', '出网带宽使用率', 'eip_upstream_bandwidth_usage', 'instance', '(sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,eip) / avg(eip_config_upstream_bandwidth{$INSTANCE}) by (instance,eip)) * 100', null, null, '%', null, null, '1', '1', null, null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('26', '3', '出网带宽', 'slb_out_bandwidth', 'instance,slb_listener_id', 'sum by(instance) (Slb_http_bps_out_rate{$INSTANCE})', null, null, 'bps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('27', '3', '入网带宽', 'slb_in_bandwidth', 'instance,slb_listener_id', 'sum by(instance) (Slb_http_bps_in_rate{$INSTANCE})', null, null, 'bps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('30', '3', '并发连接数', 'slb_max_connection', 'instance,slb_listener_id', 'sum by(instance) (Slb_all_connection_count{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('31', '3', '活跃连接数', 'slb_active_connection', 'instance,slb_listener_id', 'sum by (instance)(Slb_all_est_connection_count{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('32', '3', '非活跃连接数', 'slb_inactive_connection', 'instance,slb_listener_id', 'sum by (instance) (Slb_all_none_est_connection_count{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('33', '3', '新建连接数', 'slb_new_connection', 'instance,slb_listener_id', 'sum by(instance) (Slb_new_connection_rate{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('34', '3', '丢弃连接数', 'slb_drop_connection', 'instance,slb_listener_id', 'sum by(instance)(Slb_drop_connection_rate{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('39', '3', '异常后端云服务器数', 'slb_unhealthyserver', 'instance', 'avg by(instance) (Slb_unhealthy_server_count{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('40', '3', '正常后端云服务器数', 'slb_healthyserver', 'instance', 'avg by(instance) (Slb_healthy_server_count{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('41', '3', '7层协议查询速率', 'slb_qps', 'instance,slb_listener_id', 'sum by(instance)(Slb_request_rate{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('42', '3', '7层协议返回客户端2xx状态码数', 'slb_statuscode2xx', 'instance', 'sum by(instance) (Slb_http_2xx_rate{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('43', '3', '7层协议返回客户端3xx状态码数', 'slb_statuscode3xx', 'instance', 'sum by(instance) (Slb_http_3xx_rate{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('44', '3', '7层协议返回客户端4xx状态码数', 'slb_statuscode4xx', 'instance', 'sum by(instance) (Slb_http_4xx_rate{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('45', '3', '7层协议返回客户端5xx状态码数', 'slb_statuscode5xx', 'instance', 'sum by(instance) (Slb_http_5xx_rate{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('52', '5', '存储库总量', 'cbr_vault_size', 'instance', 'cbr_vault_size{$INSTANCE}', null, null, 'Byte', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('53', '5', '存储库使用量', 'cbr_vault_used', 'instance', 'cbr_vault_used{$INSTANCE}', null, null, 'Byte', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('54', '5', '存储库使用率', 'cbr_vault_usage_rate', 'instance', 'cbr_vault_used{$INSTANCE} / cbr_vault_size{$INSTANCE} * 100', null, null, '%', null, null, '1', '1', null, null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('55', '6', 'SNAT连接数', 'snat_connection', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('56', '6', '入方向带宽', 'inbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_recv_bytes_total_count{$INSTANCE}[1m])*8)', null, null, 'bps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('57', '6', '出方向带宽', 'outbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_send_bytes_total_count{$INSTANCE}[1m])*8)', null, null, 'bps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('58', '6', '入方向流量', 'inbound_traffic', 'instance', 'sum by (instance)(Nat_recv_bytes_total_count{$INSTANCE})', null, null, 'Byte', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('59', '6', '出方向流量', 'outbound_traffic', 'instance', 'sum by (instance)(Nat_send_bytes_total_count{$INSTANCE})', null, null, 'Byte', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('60', '6', '入方向PPS', 'inbound_pps', 'instance', 'sum by (instance)(rate(Nat_recv_packets_total_count{$INSTANCE}[1m]))', null, null, 'pps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('61', '6', '出方向PPS', 'outbound_pps', 'instance', 'sum by (instance)(rate(Nat_send_packets_total_count{$INSTANCE}[1m]))', null, null, 'pps', null, null, '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('62', '6', 'SNAT连接数使用率', 'snat_connection_ratio', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE}) / avg by (instance)(Nat_nat_max_connection_count{$INSTANCE}) *100', null, null, '%', null, null, '1', '1', null, null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('63', '1', '(基础)磁盘读IOPS', 'ecs_disk_base_read_iops', 'instance,drive', 'sum(irate(ecs_base_storage_iops_total{type="read",$INSTANCE}[15m])) by (instance,drive)', null, null, '次', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('64', '1', '(基础)磁盘写IOPS', 'ecs_disk_base_write_iops', 'instance,drive', 'sum(irate(ecs_base_storage_iops_total{type="write",$INSTANCE}[15m])) by (instance,drive)', null, null, '次', null, '1', '1', '1', null, null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('65', '1', '磁盘剩余存储量', 'ecs_filesystem_free_bytes', 'instance,device', 'ecs_filesystem_free_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('66', '1', '磁盘已用存储量', 'ecs_disk_used', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes', null, null, 'GB', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('67', '1', '磁盘存储总量', 'ecs_filesystem_size_bytes', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('68', '7', 'CPU使用率', 'bms_cpu_usage', 'instance', '100 - (100 * (sum by(instance) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))', '', '', '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('69', '7', 'CPU1分钟平均负载', 'bms_load1', 'instance', 'ecs_load1{$INSTANCE}', '', null, '', null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"windows\"}}");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('70', '7', 'CPU5分钟平均负载', 'bms_load5', 'instance', 'ecs_load5{$INSTANCE}', '', null, '', null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"windows\"}}");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('71', '7', 'CPU15分钟平均负载', 'bms_load15', 'instance', 'ecs_load15{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"windows\"}}");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('72', '7', '内存使用量', 'bms_memory_used', 'instance', 'ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes', null, null, 'Byte', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('73', '7', '内存使用率', 'bms_memory_usage', 'instance', '100 * ((ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes) / ecs_memory_MemTotal_bytes)', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('74', '7', '磁盘使用率', 'bms_disk_usage', 'instance,device', '100 * ((ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes) / ecs_filesystem_size_bytes)', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('75', '7', '磁盘读速率', 'bms_disk_read_rate', 'instance,device', 'irate(ecs_disk_read_bytes_total{$INSTANCE}[3m])', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('76', '7', '磁盘写速率', 'bms_disk_write_rate', 'instance,device', 'irate(ecs_disk_written_bytes_total{$INSTANCE}[3m])', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('77', '7', '磁盘读IOPS', 'bms_disk_read_iops', 'instance,device', 'irate(ecs_disk_reads_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('78', '7', '磁盘写IOPS', 'bms_disk_write_iops', 'instance,device', 'irate(ecs_disk_writes_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('79', '7', '流入带宽', 'bms_network_transmit_rate', 'instance,device', 'irate(ecs_network_receive_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, 'Mbps', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('80', '7', '流出带宽', 'bms_network_receive_rate', 'instance,device', 'irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, 'Mbps', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('81', '7', '包接收速率', 'bms_network_receive_packets_rate', 'instance,device', 'irate(ecs_network_receive_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('82', '7', '包发送速率', 'bms_network_transmit_packets_rate', 'instance,device', 'irate(ecs_network_transmit_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('83', '7', '(基础)CPU的平均使用率', 'bms_cpu_base_usage', 'instance', '100 * avg(irate(ecs_base_vcpu_seconds{$INSTANCE}[6m]))', null, null, '%', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('84', '7', '(基础)磁盘读速率', 'bms_disk_base_read_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type="read",$INSTANCE}[6m])', null, null, 'Byte/s', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('85', '7', '(基础)磁盘写速率', 'bms_disk_base_write_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type="write",$INSTANCE}[6m])', null, null, 'Byte/s', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('86', '7', '(基础)网卡下行带宽', 'bms_network_base_receive_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type="rx",$INSTANCE}[6m])', null, '', 'Byte/s', null, '1', '1', '1', null, null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('87', '7', '(基础)网卡上行带宽', 'bms_network_base_transmit_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type="tx",$INSTANCE}[6m])', null, null, 'Byte/s', null, '1', '1', '1', null, null, null, null);

-- ----------------------------
-- Records of config_item
-- ----------------------------
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('55', '-1', '监控周期', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('1', '-1', '统计周期', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('9', '2', '持续1个周期', '1', '1', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('10', '2', '持续3个周期', '3', '3', 1, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('11', '2', '持续5个周期', '5', '5', 2, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('12', '3', '平均值', 'Average', 'avg', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('13', '3', '最大值', 'Maximum', 'max', 1, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('14', '3', '最小值', 'Minimum', 'min', 2, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('15', '4', '大于', 'greater', '>', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('16', '4', '大于等于', 'greaterOrEqual', '>=', 1, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('17', '4', '小于', 'less', '<', 2, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('18', '4', '小于等于', 'lessOrEqual', '<=', 3, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('19', '4', '等于', 'equal', '==', 4, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('2', '-1', '持续周期', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('20', '4', '不等于', 'notEqual', '!=', 5, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('21', '-1', '概览监控项', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('22', '21', 'CPU使用率（操作系统）', null, 'ecs_cpu_usage', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('23', '21', '内存使用率（操作系统）', null, 'ecs_memory_usage', 1, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('28', '-1', '监控周期', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('29', '28', '紧急', '1', 'MAIN', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('3', '-1', '统计方式', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('30', '28', '重要', '2', 'MARJOR', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('31', '28', '次要', '3', 'MINOR', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('32', '28', '提醒', '4', 'WARN', 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('4', '-1', '对比方式', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('5', '-1', '监控数据', null, null, 0, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('51', '5', '0-3H', '0,3', '60', 1, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('52', '5', '3H-12H', '3,12', '180', 2, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('53', '5', '12H-3D', '12,72', '900', 3, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('54', '5', '3D-10D', '72,240', '2700', 4, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('6', '1', '5分钟', '300', '5m', 1, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('7', '1', '15分钟', '900', '15m', 2, null);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('8', '1', '30分钟', '1800', '30m', 3, null);






