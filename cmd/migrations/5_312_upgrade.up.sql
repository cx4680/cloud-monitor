create table t_alarm_item
(
    id                bigint(20) unsigned not null primary key AUTO_INCREMENT,
    rule_biz_id       varchar(100),
    metric_code       varchar(100),
    trigger_condition json comment '条件表达式',
    `level`           smallint,
    silences_time     varchar(50) comment '告警间隔'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

alter table t_alarm_item comment '告警规则项';

-- 转移数据
INSERT INTO t_alarm_item(rule_biz_id, metric_code, trigger_condition, `level`, silences_time)
SELECT id,
       metric_name,
       trigger_condition,
       `level`,
       silences_time
FROM t_alarm_rule;

-- 修改字段
alter table t_alarm_rule change column metric_name metric_code varchar (100);

ALTER TABLE t_alarm_rule DROP trigger_condition;

alter table t_alarm_rule
    add `type` smallint comment '1:单指标, 2:多指标',
    add template_biz_id varchar(100) comment '模板Id',
    add combination smallint comment '逻辑组合：1:或, 2:且',
    add `period` smallint,
    add times smallint;


-- 告警规则模板
CREATE TABLE `t_alarm_rule_template`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `biz_id`          varchar(256) NOT NULL,
    `monitor_type`    varchar(255) DEFAULT NULL COMMENT '监控类型',
    `product_biz_id`  varchar(100) DEFAULT NULL COMMENT '产品名称',
    `name`            varchar(255) DEFAULT NULL COMMENT '规则名称',
    `metric_code`     varchar(100) DEFAULT NULL,
    `silences_time`   varchar(255) DEFAULT NULL COMMENT '冷却周期',
    `effective_start` varchar(19)  DEFAULT NULL COMMENT '监控时间段-开始时间',
    `effective_end`   varchar(19)  DEFAULT NULL COMMENT '监控时间段-结束时间',
    `level`           int(11) DEFAULT NULL COMMENT '报警级别  \r\n紧急1 \r\n重要2 \r\n次要3\r\n提醒 4\r\n',
    `create_time`     datetime     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `type`            smallint(6) DEFAULT NULL COMMENT '1:单指标, 2:多指标',
    `combination`     smallint(6) DEFAULT NULL COMMENT '逻辑组合：1:或, 2:且',
    `period`          smallint(6) DEFAULT NULL,
    `times`           smallint(6) DEFAULT NULL,
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;
alter table t_alarm_rule_template comment '告警规则模板';

CREATE TABLE `t_alarm_item_template`
(
    `id`                   bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `rule_template_biz_id` varchar(100) DEFAULT NULL,
    `metric_code`          varchar(100) DEFAULT NULL,
    `trigger_condition`    json         DEFAULT NULL COMMENT '条件表达式',
    `level`                smallint(6) DEFAULT NULL,
    `silences_time`        varchar(50)  DEFAULT NULL COMMENT '告警间隔',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;
alter table t_alarm_item_template comment '告警规则模板项';

create table t_tenant_alarm_template_rel
(
    id              bigint not null auto_increment,
    tenant_id       varchar(100),
    template_biz_id varchar(100),
    create_time     datetime,
    primary key (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;
alter table t_tenant_alarm_template_rel comment '租户告警模板关系表';






INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type) VALUES ('12', '缓存数据库Redis', '1', 'redis', null, null, '/productmonitoring/redis', '0 0 0/1 * * ?', 'http://product-redis-ndb-redis-manage.product-redis.svc.cluster.local:8888', '/v1/redis/instance', 'redis', '13', '云产品监控');
INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type) VALUES ('13', '非关系型数据库NoSQL', '1', 'mongo', null, null, '/productmonitoring/mongo', '0 0 0/1 * * ?', 'http://dbaas-manage.product-dbaas.svc.cluster.local', '/v2/mongo/instance', 'mongo', '14', '云产品监控');
INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type) VALUES ('14', 'API网关CGW', '1', 'cgw', null, null, '/productmonitoring/cgw', '0 0 0/1 * * ?', 'http://cgw-cgw-manage-admin.product-cgw.svc.cluster.local:8080', '/gateway/instance/page', 'cgw', '15', '云产品监控');

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('185', '12', 'CPU使用率', 'redis_cpu_usage', 'instance', 'redis_cpu_usage{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('186', '12', '内存使用率', 'redis_mem_usage', 'instance', 'redis_mem_usage{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('187', '12', '客户端连接数', 'redis_connected_clients', 'instance', 'redis_connected_clients{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('188', '12', 'TPS', 'redis_tps', 'instance', 'redis_tps{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');
-- INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('189', '12', '每秒操作数', 'redis_instantaneous_ops', 'instance', 'redis_instantaneous_ops{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');
-- INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('190', '12', '阻塞客户端数量', 'redis_blocked_clients', 'instance', 'redis_blocked_clients{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');
-- INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('191', '12', '缓存命中率', 'redis_hit_rate', 'instance', 'redis_hit_rate{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');
-- INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('192', '12', '内存碎片率', 'redis_mem_fragmentation_ratio', 'instance', 'redis_mem_fragmentation_ratio{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');
-- INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('193', '12', '内存峰值', 'redis_used_memory', 'instance', 'redis_used_memory{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');
-- INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('194', '12', '占用内存大小', 'redis_used_memory_human', 'instance', 'redis_used_memory_human{$INSTANCE}', null, null, null, null, null, '1', '1', null, null, null, null, 'chart');

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('195', '13', 'mongos客户端当前连接数', 'mongo_mongos_current_connections', 'instance,pod', 'mongo_mongos_current_connections{$INSTANCE}', null, null, '个', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('196', '13', 'shard客户端当前连接数', 'mongo_shard_current_connections', 'instance,pod', 'mongo_shard_current_connections{$INSTANCE}', null, null, '个', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('197', '13', 'configServer客户端当前连接数', 'mongo_config_current_connections', 'instance,pod', 'mongo_config_current_connections{$INSTANCE}', null, null, '个', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('198', '13', 'mongo各角色总当前连接数', 'mongo_total_current_connections', 'instance', 'mongo_total_current_connections{$INSTANCE}', null, null, '个', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('199', '13', '每个mongos的内存使用率', 'mongo_mongos_memory_ratio', 'instance,pod', 'mongo_mongos_memory_ratio{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('200', '13', 'configServer的内存使用率', 'mongo_config_memory_ratio', 'instance', 'mongo_config_memory_ratio{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('201', '13', '每个分片的内存使用率', 'mongo_shard_memory_ratio', 'instance,pod', 'mongo_shard_memory_ratio{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('202', '13', '每个mongos的CPU使用率', 'mongo_mongos_cpu_ratio', 'instance,pod', 'mongo_mongos_cpu_ratio{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('203', '13', '每个shard的CPU使用率', 'mongo_shard_cpu_ratio', 'instance,pod', 'mongo_shard_cpu_ratio{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('204', '13', 'configServer的内存使用率', 'mongo_config_cpu_ratio', 'instance', 'mongo_config_cpu_ratio{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('205', '14', 'QPS', 'guard_nginx_http_current_reqs', 'instance,route', 'sum by(instance,service,route)(rate(guard_nginx_http_current_reqs{$INSTANCE}[3m]))', null, null, '次/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('206', '14', 'P90接口响应延时', 'guard_http_latency_bucket_api_p90', 'instance,route', 'histogram_quantile(0.90, sum by(instance,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('207', '14', 'P95接口响应延时', 'guard_http_latency_bucket_api_p95', 'instance,route', 'histogram_quantile(0.95, sum by(instance,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('208', '14', 'P99接口响应延时', 'guard_http_latency_bucket_api_p99', 'instance,route', 'histogram_quantile(0.99, sum by(instance,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('209', '14', 'P90服务响应延时', 'guard_http_latency_bucket_service_p90', 'instance,service', 'histogram_quantile(0.90, sum by(instance,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('210', '14', 'P95服务响应延时', 'guard_http_latency_bucket_service_p95', 'instance,service', 'histogram_quantile(0.95, sum by(instance,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('211', '14', 'P99服务响应延时', 'guard_http_latency_bucket_service_p99', 'instance,service', 'histogram_quantile(0.99, sum by(instance,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('212', '14', '接口成功率', 'guard_nginx_url_request_succ_api', 'instance,route', 'sum by(instance,service,route)(rate(guard_nginx_url_request_succ{code="200",$INSTANCE}[3m]))/sum by(instance,service,route)(rate(guard_nginx_url_request_succ{code="total",$INSTANCE}[3m]))', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('213', '14', '服务成功率', 'guard_nginx_url_request_succ_service', 'instance,service', 'sum by(instance,service)(rate(guard_nginx_url_request_succ{code="200",$INSTANCE}[3m]))/sum by(instance,service)(rate(guard_nginx_url_request_succ{code="total",$INSTANCE}[3m]))', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('214', '14', '接口入口带宽监控', 'guard_bandwidth_api_ingress', 'instance,route', 'sum by(instance,service,route)(rate(guard_bandwidth{type="ingress",$INSTANCE}[3m]))', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('215', '14', '接口出口带宽监控', 'guard_bandwidth_api_egress', 'instance,route', 'sum by(instance,service,route)(rate(guard_bandwidth{type="egress",$INSTANCE}[3m]))', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('216', '14', '服务入口带宽监控', 'guard_bandwidth_service_ingress', 'instance,service', 'sum by(instance,service)(irate(guard_bandwidth{type="ingress",$INSTANCE}[3m])) ', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('217', '14', '服务出口带宽监控', 'guard_bandwidth_service_egress', 'instance,service', 'sum by(instance,service)(rate(guard_bandwidth{type="egress",$INSTANCE}[3m]))', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('218', '14', '接口调用量', 'guard_http_apirequests', 'instance,service,route', 'guard_http_apirequests{$INSTANCE}', null, null, '次', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('219', '14', '服务调用量', 'guard_http_servicerequests', 'instance,service', 'guard_http_apirequests{$INSTANCE}', null, null, '次', null, null, '1', '1', null, null, null, null, 'chart');





INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('1', '云产品监控', '1', 'sd_ecs_cpu_u', 'ecs_cpu_usage', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);
INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('2', '云产品监控', '1', 'sd_ecs_mem_u', 'ecs_memory_usage', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);
INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('3', '云产品监控', '1', 'sd_ecs_disk_u', 'ecs_disk_usage', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);
INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('4', '云产品监控', '8', 'sd_rdb_con_u', 'mysql_current_connection_percent', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);
INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('5', '云产品监控', '8', 'sd_rdb_disk_u', 'mysql_disk_usage', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);
INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('6', '云产品监控', '8', 'sd_rdb_read_d', 'mysql_slave_seconds_behind_master', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);
INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('7', '云产品监控', '3', 'sd_slb_lost_c', 'slb_drop_connection', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);
INSERT INTO `t_alarm_rule_template` (`biz_id`, `monitor_type`, `product_biz_id`, `name`, `metric_code`, `silences_time`, `effective_start`, `effective_end`, `level`, `create_time`, `type`, `combination`, `period`, `times`) VALUES ('8', '云产品监控', '12', 'sd_redis_cpu_u', 'redis_cpu_usage', '3小时', '00:00', '23:59', NULL, CURRENT_TIMESTAMP, 1, NULL, NULL, NULL);

INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('1', 'ecs_cpu_usage', '{"unit": "%", "times": 1, "labels": "instance", "period": 300, "threshold": 95, "metricCode": "ecs_cpu_usage", "metricName": "CPU使用率", "statistics": "Maximum"
, "comparisonOperator": "greaterOrEqual"}', 2, '3小时');
INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('2', 'ecs_memory_usage', '{"unit": "%", "times": 1, "labels": "instance", "period": 300, "threshold": 95, "metricCode": "ecs_memory_usage", "metricName": "内存使用率", "statistics": "Maximum"
, "comparisonOperator": "greaterOrEqual"}', 2, '3小时');
INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('3', 'ecs_disk_usage', '{"unit": "%", "times": 1, "labels": "instance", "period": 300, "threshold": 95, "metricCode": "ecs_disk_usage", "metricName": "磁盘使用率", "statistics": "Maximum"
, "comparisonOperator": "greaterOrEqual"}', 2, '3小时');
INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('4', 'mysql_current_connection_percent', '{"unit": "%", "times": 1, "labels": "instance", "period": 300, "threshold": 80, "metricCode": "mysql_current_connection_percent", "metricName": "当前连接占比", "statistics": "Maximum"
, "comparisonOperator": "greaterOrEqual"}', 2, '3小时');
INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('5', 'mysql_disk_usage', '{"unit": "%", "times": 1, "labels": "instance", "period": 300, "threshold": 80, "metricCode": "mysql_disk_usage", "metricName": "磁盘使用率", "statistics": "Maximum"
, "comparisonOperator": "greaterOrEqual"}', 2, '3小时');
INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('6', 'mysql_slave_seconds_behind_master', '{"unit": "s", "times": 1, "labels": "instance", "period": 300, "threshold": 30, "metricCode": "mysql_slave_seconds_behind_master", "metricName": "复制延迟", "statistics": "Maximum"
, "comparisonOperator": "greaterOrEqual"}', 2, '3小时');
INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('7', 'slb_drop_connection', '{"unit": "个/s", "times": 1, "labels": "instance", "period": 300, "threshold": 0, "metricCode": "slb_drop_connection", "metricName": "丢弃连接数", "statistics": "Maximum"
, "comparisonOperator": "greater"}', 3, '3小时');
INSERT INTO `t_alarm_item_template` (`rule_template_biz_id`, `metric_code`, `trigger_condition`, `level`, `silences_time`) VALUES ('8', 'redis_cpu_usage', '{"unit": "%", "times": 1, "labels": "instance", "period": 300, "threshold": 80, "metricCode": "redis_cpu_usage", "metricName": "CPU使用率", "statistics": "Maximum"
, "comparisonOperator": "greaterOrEqual"}', 3, '3小时');



ALTER TABLE t_alarm_record MODIFY COLUMN source_type varchar(100) COMMENT '资源类型';
ALTER TABLE t_alarm_record MODIFY COLUMN rule_id varchar(50) COMMENT '规则id';
