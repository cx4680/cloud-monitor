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
INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type) VALUES ('14', 'API网关CGW', '0', 'cgw', null, null, '/productmonitoring/cgw', '0 0 0/1 * * ?', 'http://cgw-cgw-manage-admin.product-cgw.svc.cluster.local:8080', '/gateway/instance/page', 'cgw', '15', '云产品监控');

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
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('204', '13', 'configServer的CPU使用率', 'mongo_config_cpu_ratio', 'instance', 'mongo_config_cpu_ratio{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('205', '14', 'QPS', 'guard_nginx_http_current_reqs', 'instance,route', 'sum by(instance,service,route)(rate(guard_nginx_http_current_reqs{$INSTANCE}[3m]))', null, null, '次/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('206', '14', 'P90接口响应延时', 'guard_http_latency_bucket_api_p90', 'instance,route', 'histogram_quantile(0.90, sum by(instance,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('207', '14', 'P95接口响应延时', 'guard_http_latency_bucket_api_p95', 'instance,route', 'histogram_quantile(0.95, sum by(instance,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('208', '14', 'P99接口响应延时', 'guard_http_latency_bucket_api_p99', 'instance,route', 'histogram_quantile(0.99, sum by(instance,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('209', '14', 'P90服务响应延时', 'guard_http_latency_bucket_service_p90', 'instance,service', 'histogram_quantile(0.90, sum by(instance,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('210', '14', 'P95服务响应延时', 'guard_http_latency_bucket_service_p95', 'instance,service', 'histogram_quantile(0.95, sum by(instance,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('211', '14', 'P99服务响应延时', 'guard_http_latency_bucket_service_p99', 'instance,service', 'histogram_quantile(0.99, sum by(instance,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))', null, null, 'ms', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('212', '14', '接口成功率', 'guard_nginx_url_request_succ_api', 'instance,route', 'sum by(instance,service,route)(rate(guard_nginx_url_request_succ{code="200",$INSTANCE}[3m]))/sum by(instance,service,route)(rate(guard_nginx_url_request_succ{code="total",$INSTANCE}[3m]))*100', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('213', '14', '服务成功率', 'guard_nginx_url_request_succ_service', 'instance,service', 'sum by(instance,service)(rate(guard_nginx_url_request_succ{code="200",$INSTANCE}[3m]))/sum by(instance,service)(rate(guard_nginx_url_request_succ{code="total",$INSTANCE}[3m]))*100', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('214', '14', '接口入口带宽监控', 'guard_bandwidth_api_ingress', 'instance,route', 'sum by(instance,service,route)(rate(guard_bandwidth{type="ingress",$INSTANCE}[3m]))', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('215', '14', '接口出口带宽监控', 'guard_bandwidth_api_egress', 'instance,route', 'sum by(instance,service,route)(rate(guard_bandwidth{type="egress",$INSTANCE}[3m]))', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('216', '14', '服务入口带宽监控', 'guard_bandwidth_service_ingress', 'instance,service', 'sum by(instance,service)(irate(guard_bandwidth{type="ingress",$INSTANCE}[3m])) ', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('217', '14', '服务出口带宽监控', 'guard_bandwidth_service_egress', 'instance,service', 'sum by(instance,service)(rate(guard_bandwidth{type="egress",$INSTANCE}[3m]))', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('218', '14', '接口调用量', 'guard_http_apirequests', 'instance,service,route', 'sum by(instance,service,route)(guard_http_apirequests{$INSTANCE})', null, null, '次', null, null, '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('219', '14', '服务调用量', 'guard_http_servicerequests', 'instance,service', 'sum by(instance,service)(guard_http_apirequests{$INSTANCE})', null, null, '次', null, null, '1', '1', null, null, null, null, 'chart');





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


UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load1{$INSTANCE})' WHERE biz_id = '2';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load5{$INSTANCE})' WHERE biz_id = '3';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load15{$INSTANCE})' WHERE biz_id = '4';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id = '5';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id = '6';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id = '7';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '8';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '9';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '10';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '11';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '12';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '13';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '14';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '15';
UPDATE t_monitor_item SET metrics_linux = '100 * avg by(instance)(irate(ecs_base_vcpu_seconds{$INSTANCE}[6m]))' WHERE biz_id = '16';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,drive)(irate(ecs_base_storage_traffic_bytes_total{type="read",$INSTANCE}[6m]))' WHERE biz_id = '17';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,drive)(irate(ecs_base_storage_traffic_bytes_total{type="write",$INSTANCE}[6m]))' WHERE biz_id = '18';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,interface)(irate(ecs_base_network_traffic_bytes_total{type="rx",$INSTANCE}[6m]))' WHERE biz_id = '19';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,interface)(irate(ecs_base_network_traffic_bytes_total{type="tx",$INSTANCE}[6m]))' WHERE biz_id = '20';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(cbr_vault_size{$INSTANCE})' WHERE biz_id = '52';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(cbr_vault_used{$INSTANCE})' WHERE biz_id = '53';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(cbr_vault_used{$INSTANCE}) / sum by(instance)(cbr_vault_size{$INSTANCE}) * 100' WHERE biz_id = '54';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '65';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '66';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id = '67';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load1{$INSTANCE})' WHERE biz_id = '69';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load5{$INSTANCE})' WHERE biz_id = '70';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load15{$INSTANCE})' WHERE biz_id = '71';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id = '72';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id = '73';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id = '74';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '75';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '76';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '77';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '78';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '79';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '80';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '81';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '82';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '83';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '84';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id = '85';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_slave_io{$INSTANCE})' WHERE biz_id = '86';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_slave_sql{$INSTANCE})' WHERE biz_id = '87';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_slave_seconds_behind_master{$INSTANCE})' WHERE biz_id = '88';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_active_connections{$INSTANCE})' WHERE biz_id = '89';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_current_connection_percent{$INSTANCE})' WHERE biz_id = '90';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_qps{$INSTANCE})' WHERE biz_id = '91';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_tps{$INSTANCE})' WHERE biz_id = '92';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_select_ps{$INSTANCE})' WHERE biz_id = '93';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_update_ps{$INSTANCE})' WHERE biz_id = '94';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_insert_ps{$INSTANCE})' WHERE biz_id = '95';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_delete_ps{$INSTANCE})' WHERE biz_id = '96';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_cpu_usage{$INSTANCE})' WHERE biz_id = '97';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_mem_usage{$INSTANCE})' WHERE biz_id = '98';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_disk_usage{$INSTANCE})' WHERE biz_id = '99';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_select_ps{$INSTANCE})' WHERE biz_id = '100';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_update_ps{$INSTANCE})' WHERE biz_id = '101';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_insert_ps{$INSTANCE})' WHERE biz_id = '102';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_delete_ps{$INSTANCE})' WHERE biz_id = '103';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_cache_hit_rate{$INSTANCE})' WHERE biz_id = '104';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_reads_ps{$INSTANCE})' WHERE biz_id = '105';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_writes_ps{$INSTANCE})' WHERE biz_id = '106';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_buffer_pool_pages_dirty{$INSTANCE})' WHERE biz_id = '107';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_buffer_pool_bytes_dirty{$INSTANCE})' WHERE biz_id = '108';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_innodb_log_waits{$INSTANCE})' WHERE biz_id = '109';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_binlog_cache_disk_use{$INSTANCE})' WHERE biz_id = '110';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_slow_queries_per_min{$INSTANCE})' WHERE biz_id = '111';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_long_query_count{$INSTANCE})' WHERE biz_id = '112';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_long_query_alert_count{$INSTANCE})' WHERE biz_id = '113';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_exec_statememt_frequency{$INSTANCE})' WHERE biz_id = '114';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_exec_statememt_frequency{$INSTANCE})' WHERE biz_id = '115';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_write_frequency{$INSTANCE})' WHERE biz_id = '116';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_top_statememt_avg_exec_time{$INSTANCE})' WHERE biz_id = '117';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_top_statememt_exec_err_rate{$INSTANCE})' WHERE biz_id = '118';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mysql_current_cons_num{$INSTANCE})' WHERE biz_id = '119';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_tps{$INSTANCE}[1m]))' WHERE biz_id = '120';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_qps{$INSTANCE}[1m]))' WHERE biz_id = '121';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_ips{$INSTANCE}[1m]))' WHERE biz_id = '122';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_dps{$INSTANCE}[1m]))' WHERE biz_id = '123';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_ups{$INSTANCE}[1m]))' WHERE biz_id = '124';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_ddlps{$INSTANCE}[1m]))' WHERE biz_id = '125';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_nioips{$INSTANCE}[1m]))' WHERE biz_id = '126';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_nio_ops{$INSTANCE}[1m]))' WHERE biz_id = '127';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_fio_ips{$INSTANCE}[1m]))' WHERE biz_id = '128';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(rate(dm_global_status_fio_ops{$INSTANCE}[1m]))' WHERE biz_id = '129';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_mem_used{$INSTANCE})' WHERE biz_id = '130';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_cpu_use_rate{$INSTANCE})' WHERE biz_id = '131';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_sessions{$INSTANCE})' WHERE biz_id = '132';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_active_sessions{$INSTANCE})' WHERE biz_id = '133';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_task_waiting{$INSTANCE})' WHERE biz_id = '134';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_task_ready{$INSTANCE})' WHERE biz_id = '135';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_task_total_wait_time{$INSTANCE})' WHERE biz_id = '136';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_avg_wait_time{$INSTANCE})' WHERE biz_id = '137';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(dm_global_status_threads{$INSTANCE})' WHERE biz_id = '138';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_cpu_usage{$INSTANCE})' WHERE biz_id = '139';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_mem_usage{$INSTANCE})' WHERE biz_id = '140';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_disk_usage{$INSTANCE})' WHERE biz_id = '141';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_qps{$INSTANCE})' WHERE biz_id = '142';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_rqps{$INSTANCE})' WHERE biz_id = '143';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_wqps{$INSTANCE})' WHERE biz_id = '144';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_tps{$INSTANCE})' WHERE biz_id = '145';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_mean_exec_time{$INSTANCE})' WHERE biz_id = '146';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_open_ct_num{$INSTANCE})' WHERE biz_id = '147';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(pg_active_ct_num{$INSTANCE})' WHERE biz_id = '148';

UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_controller_kafkacontroller_activecontrollercount{$INSTANCE})' WHERE biz_id = '149';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_brokers{$INSTANCE})' WHERE biz_id = '150';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_controller_kafkacontroller_globalpartitioncount{$INSTANCE})' WHERE biz_id = '151';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_controller_kafkacontroller_globaltopiccount{$INSTANCE})' WHERE biz_id = '152';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_controller_kafkacontroller_offlinepartitionscount{$INSTANCE})' WHERE biz_id = '153';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_controller_kafkacontroller_preferredreplicaimbalancecount{$INSTANCE})' WHERE biz_id = '154';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_controller_kafkacontroller_topicstodeletecount{$INSTANCE})' WHERE biz_id = '155';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_controller_controllerstats_uncleanleaderelectionspersec{$INSTANCE})' WHERE biz_id = '156';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_replicamanager_leadercount{$INSTANCE})' WHERE biz_id = '157';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_replicamanager_partitioncount{$INSTANCE})' WHERE biz_id = '158';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_replicamanager_underminisrpartitioncount{$INSTANCE})' WHERE biz_id = '159';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_replicamanager_underreplicatedpartitions{$INSTANCE})' WHERE biz_id = '160';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_replicamanager_reassigningpartitions{$INSTANCE})' WHERE biz_id = '161';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_brokertopicmetrics_bytesinpersec{$INSTANCE})' WHERE biz_id = '162';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_brokertopicmetrics_bytesoutpersec{$INSTANCE})' WHERE biz_id = '163';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_server_brokertopicmetrics_messagesinpersec{$INSTANCE})' WHERE biz_id = '164';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(kafka_consumergroup_lag{$INSTANCE})' WHERE biz_id = '165';

UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(dm_global_status_mem_use_rate{$INSTANCE})' WHERE biz_id = '166';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load1{$INSTANCE})' WHERE biz_id = '168';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load5{$INSTANCE})' WHERE biz_id = '169';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_load15{$INSTANCE})' WHERE biz_id = '170';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id = '171';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id = '172';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id = '173';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '174';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '175';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '176';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '177';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '178';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '179';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '180';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '181';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '182';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '183';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id = '184';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_cpu_usage{$INSTANCE})' WHERE biz_id = '185';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_mem_usage{$INSTANCE})' WHERE biz_id = '186';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_connected_clients{$INSTANCE})' WHERE biz_id = '187';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_tps{$INSTANCE})' WHERE biz_id = '187';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,pod)(mongo_mongos_current_connections{$INSTANCE})' WHERE biz_id = '195';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,pod)(mongo_shard_current_connections{$INSTANCE})' WHERE biz_id = '196';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,pod)(mongo_config_current_connections{$INSTANCE})' WHERE biz_id = '197';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mongo_total_current_connections{$INSTANCE})' WHERE biz_id = '198';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,pod)(mongo_mongos_memory_ratio{$INSTANCE})' WHERE biz_id = '199';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mongo_config_memory_ratio{$INSTANCE})' WHERE biz_id = '200';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,pod)(mongo_shard_memory_ratio{$INSTANCE})' WHERE biz_id = '201';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,pod)(mongo_mongos_cpu_ratio{$INSTANCE})' WHERE biz_id = '202';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,pod)(mongo_shard_cpu_ratio{$INSTANCE})' WHERE biz_id = '203';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(mongo_config_cpu_ratio{$INSTANCE})' WHERE biz_id = '204';