-- 清除303版本测试可能存在的脏数据
DELETE FROM `t_monitor_product` WHERE biz_id in ('5','6','7');

INSERT INTO `t_monitor_product` (`biz_id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('5', '云备份CBR', '1', 'cbr', null, null, '/productmonitoring/cbr', '0 0 0/1 * * ?', 'http://product-backup-backup-manage.product-backup', '/noauth/backup/vault/pageList', 'cbr');
INSERT INTO `t_monitor_product` (`biz_id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('6', 'NAT网关', '1', 'nat', null, null, '/productmonitoring/nat', '0 0 0/1 * * ?', 'http://product-nat-controller-nat-manage.product-nat-gw', '/nat-gw/inner/nat/page', 'nat');
INSERT INTO `t_monitor_product` (`biz_id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('7', '裸金属服务器', '1', 'bms', null, null, '/productmonitoring/bms', '0 0 0/1 * * ?', 'http://bms-manage.product-bms:8080', '/compute/bms/ops/v1', 'bms');
INSERT INTO `t_monitor_product` (`biz_id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('8', '云数据库MySQL', '1', 'mysql', null, null, '/productmonitoring/mysql', '0 0 0/1 * * ?', 'http://product-mysql-rds-mysql-manage.product-mysql.svc.cluster.local:8888', '/v1/mysql/instance', 'mysql');
INSERT INTO `t_monitor_product` (`biz_id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('9', '云数据库达梦RDBDM', '1', 'dm', null, null, '/productmonitoring/dm', '0 0 0/1 * * ?', 'http://product-dm-rds-dm-manage.product-dm.svc.cluster.local:8888', '/v1/dm/instance', 'dm');
INSERT INTO `t_monitor_product` (`biz_id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('10', '云数据库PostgreSQL', '1', 'postgresql', null, null, '/productmonitoring/postgresql', '0 0 0/1 * * ?', 'http://product-postgresql-rdb-pg-manage.product-pg.svc.cluster.local:8888', '/v1/postgresql/instance', 'postgresql');
INSERT INTO `t_monitor_product` (`biz_id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`, `cron`, `host`, `page_url`, `abbreviation`) VALUES ('11', '消息队列Kafka', '1', 'kafka', null, null, '/productmonitoring/kafka', '0 0 0/1 * * ?', 'http://cmq-kafka.product-cmq-kafka:8080', '/kafka/v1/cluster', 'kafka');

DELETE FROM `t_monitor_item` WHERE biz_id in ('83','84','85','86','87');

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('83', '7', '磁盘剩余存储量', 'bms_filesystem_free_bytes', 'instance,device', 'ecs_filesystem_free_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('84', '7', '磁盘已用存储量', 'bms_disk_used', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes', null, null, 'GB', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('85', '7', '磁盘存储总量', 'bms_filesystem_size_bytes', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('86', '8', '主从复制IO线程状态', 'mysql_slave_io', 'instance', 'mysql_slave_io{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"Basic\"}}");
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('87', '8', '主从复制SQL线程状态', 'mysql_slave_sql', 'instance', 'mysql_slave_sql{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"Basic\"}}");
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('88', '8', '复制延迟', 'mysql_slave_seconds_behind_master', 'instance', 'mysql_slave_seconds_behind_master{$INSTANCE}', null, null, 's', null, '2', '1', '1', null, null, null, "{{ne .OSTYPE \"Basic\"}}");

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('89', '8', '活跃连接数', 'mysql_active_connections', 'instance', 'mysql_active_connections{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('90', '8', '当前连接占比', 'mysql_current_connection_percent', 'instance', 'mysql_current_connection_percent{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('91', '8', 'QPS', 'mysql_qps', 'instance', 'mysql_qps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('92', '8', 'TPS', 'mysql_tps', 'instance', 'mysql_tps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('93', '8', '每秒查询数量', 'mysql_select_ps', 'instance', 'mysql_select_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('94', '8', '每秒更新数量', 'mysql_update_ps', 'instance', 'mysql_update_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('95', '8', '每秒插入数量', 'mysql_insert_ps', 'instance', 'mysql_insert_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('96', '8', '每秒删除数量', 'mysql_delete_ps', 'instance', 'mysql_delete_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('97', '8', 'CPU使用率', 'mysql_cpu_usage', 'instance', 'mysql_cpu_usage{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('98', '8', '内存使用率', 'mysql_mem_usage', 'instance', 'mysql_mem_usage{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('99', '8', '磁盘使用率', 'mysql_disk_usage', 'instance', 'mysql_disk_usage{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('100', '8', 'InnoDB每秒查询行数', 'mysql_innodb_select_ps', 'instance', 'mysql_innodb_select_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('101', '8', 'InnoDB每秒更新行数', 'mysql_innodb_update_ps', 'instance', 'mysql_innodb_update_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('102', '8', 'InnoDB每秒插入行数', 'mysql_innodb_insert_ps', 'instance', 'mysql_innodb_insert_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('103', '8', 'InnoDB每秒删除行数', 'mysql_innodb_delete_ps', 'instance', 'mysql_innodb_delete_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('104', '8', 'InnoDB缓存命中率', 'mysql_innodb_cache_hit_rate', 'instance', 'mysql_innodb_cache_hit_rate{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('105', '8', 'InnoDB每秒读次数', 'mysql_innodb_reads_ps', 'instance', 'mysql_innodb_reads_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('106', '8', 'InnoDB每秒写次数', 'mysql_innodb_writes_ps', 'instance', 'mysql_innodb_writes_ps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('107', '8', 'InnoDB脏页数量', 'mysql_innodb_buffer_pool_pages_dirty', 'instance', 'mysql_innodb_buffer_pool_pages_dirty{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('108', '8', 'InnoDB脏页大小', 'mysql_innodb_buffer_pool_bytes_dirty', 'instance', 'mysql_innodb_buffer_pool_bytes_dirty{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('109', '8', 'InnoDB日志写等待', 'mysql_innodb_log_waits', 'instance', 'mysql_innodb_log_waits{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('110', '8', '大事务', 'mysql_binlog_cache_disk_use', 'instance', 'mysql_binlog_cache_disk_use{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('111', '8', '每分钟慢查询数量', 'mysql_slow_queries_per_min', 'instance', 'mysql_slow_queries_per_min{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('112', '8', '长时间执行SQL(执行时间超过600秒)', 'mysql_long_query_count', 'instance', 'mysql_long_query_count{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('113', '8', '长时间执行SQL报警(执行时间超过1800秒)', 'mysql_long_query_alert_count', 'instance', 'mysql_long_query_alert_count{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('114', '8', '执行语句频率', 'mysql_exec_statememt_frequency', 'instance', 'mysql_exec_statememt_frequency{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('115', '8', '读频率', 'mysql_read_frequency', 'instance', 'mysql_exec_statememt_frequency{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('116', '8', '写频率', 'mysql_write_frequency', 'instance', 'mysql_write_frequency{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('117', '8', 'TOP语句平均执行时间', 'mysql_top_statememt_avg_exec_time', 'instance', 'mysql_top_statememt_avg_exec_time{$INSTANCE}', null, null, 'us', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('118', '8', 'TOP语句执行错误率', 'mysql_top_statememt_exec_err_rate', 'instance', 'mysql_top_statememt_exec_err_rate{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('119', '8', '当前打开连接数', 'mysql_current_cons_num', 'instance', 'mysql_current_cons_num{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('120', '9', '每秒事务数', 'dm_global_status_tps', 'instance', 'dm_global_status_tps{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('121', '9', '每秒执行select SQL语句数', 'dm_global_status_qps', 'instance', 'dm_global_status_qps{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('122', '9', '每秒执行insert SQL语句数', 'dm_global_status_ips', 'instance', 'dm_global_status_ips{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('123', '9', '每秒执行delete SQL语句数', 'dm_global_status_dps', 'instance', 'dm_global_status_dps{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('124', '9', '每秒执行update SQL语句数', 'dm_global_status_ups', 'instance', 'dm_global_status_ups{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('125', '9', '每秒执行DDL SQL语句数', 'dm_global_status_ddlps', 'instance', 'dm_global_status_ddlps{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('126', '9', '每秒从客户端接收字节数', 'dm_global_status_nioips', 'instance', 'dm_global_status_nioips{$INSTANCE}', null, null, 'Byte', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('127', '9', '每秒往客户端发送字节数', 'dm_global_status_nio_ops', 'instance', 'dm_global_status_nio_ops{$INSTANCE}', null, null, 'Byte', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('128', '9', '每秒读取字节数', 'dm_global_status_fio_ips', 'instance', 'dm_global_status_fio_ips{$INSTANCE}', null, null, 'Byte', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('129', '9', '每秒写入字节数', 'dm_global_status_fio_ops', 'instance', 'dm_global_status_fio_ips{$INSTANCE}', null, null, 'Byte', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('130', '9', '内存占用字节数', 'dm_global_status_mem_used', 'instance', 'dm_global_status_mem_used{$INSTANCE}', null, null, 'Byte', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('131', '9', 'CPU使用率', 'dm_global_status_cpu_use_rate', 'instance', 'dm_global_status_cpu_use_rate{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('132', '9', '总会话数', 'dm_global_status_sessions', 'instance', 'dm_global_status_sessions{$INSTANCE}', null, null, null, null, '2', null, '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('133', '9', '活动会话数', 'dm_global_status_active_sessions', 'instance', 'dm_global_status_active_sessions{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('134', '9', '等待处理任务数', 'dm_global_status_task_waiting', 'instance', 'dm_global_status_task_waiting{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('135', '9', '已处理任务数', 'dm_global_status_task_ready', 'instance', 'dm_global_status_task_ready{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('136', '9', '已处理任务的总等待时间', 'dm_global_status_task_total_wait_time', 'instance', 'dm_global_status_task_total_wait_time{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('137', '9', '已处理任务的平均等待时间', 'dm_global_status_avg_wait_time', 'instance', 'dm_global_status_avg_wait_time{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('138', '9', '线程数', 'dm_global_status_threads', 'instance', 'dm_global_status_threads{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('139', '10', 'CPU使用率', 'pg_cpu_usage', 'instance', 'pg_cpu_usage{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('140', '10', '内存使用率', 'pg_mem_usage', 'instance', 'pg_mem_usage{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('141', '10', '磁盘使用率', 'pg_disk_usage', 'instance', 'pg_disk_usage{$INSTANCE}', null, null, '%', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('142', '10', 'QPS', 'pg_qps', 'instance', 'pg_qps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('143', '10', '读QPS', 'pg_rqps', 'instance', 'pg_rqps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('144', '10', '写QPS', 'pg_wqps', 'instance', 'pg_wqps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('145', '10', 'TPS', 'pg_tps', 'instance', 'pg_tps{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('146', '10', '最长平均执行时间', 'mean_exec_time', 'instance', 'mean_exec_time{$INSTANCE}', null, null, 'ms', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('147', '10', '当前打开连接数', 'pg_open_ct_num', 'instance', 'pg_open_ct_num{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('148', '10', '当前活跃连接数', 'pg_active_ct_num', 'instance', 'pg_active_ct_num{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null);

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('149', '11', 'Broker控制器', 'kafka_brokers', 'instance', 'kafka_brokers{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('150', '11', '在线Broker数', 'kafka_controller_kafkacontroller_activecontrollercount', 'instance', 'kafka_controller_kafkacontroller_activecontrollercount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('151', '11', '全局分区数', 'kafka_controller_kafkacontroller_globalpartitioncount', 'instance', 'kafka_controller_kafkacontroller_globalpartitioncount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('152', '11', '全局主题数', 'kafka_controller_kafkacontroller_globaltopiccount', 'instance', 'kafka_controller_kafkacontroller_globaltopiccount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('153', '11', '离线主题数', 'kafka_controller_kafkacontroller_offlinepartitionscount', 'instance', 'kafka_controller_kafkacontroller_offlinepartitionscount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('154', '11', '未平衡的副本数', 'kafka_controller_kafkacontroller_preferredreplicaimbalancecount', 'instance', 'kafka_controller_kafkacontroller_offlinepartitionscount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('155', '11', '删除的主题数', 'kafka_controller_kafkacontroller_topicstodeletecount', 'instance', 'kafka_controller_kafkacontroller_topicstodeletecount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('156', '11', '未同步的选举', 'kafka_controller_controllerstats_uncleanleaderelectionspersec', 'instance', 'kafka_controller_controllerstats_uncleanleaderelectionspersec{$INSTANCE}', null, null, '次/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('157', '11', 'Leader数量', 'kafka_server_replicamanager_leadercount', 'instance', 'kafka_server_replicamanager_leadercount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('158', '11', '分区的数量', 'kafka_server_replicamanager_partitioncount', 'instance', 'kafka_server_replicamanager_partitioncount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('159', '11', '最小未同步分区的数量', 'kafka_server_replicamanager_underminisrpartitioncount', 'instance', 'kafka_server_replicamanager_underminisrpartitioncount{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('160', '11', '没有副本的分区数', 'kafka_server_replicamanager_underreplicatedpartitions', 'instance', 'kafka_server_replicamanager_underreplicatedpartitions{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('161', '11', '再分配的分区数', 'kafka_server_replicamanager_reassigningpartitions', 'instance', 'kafka_server_replicamanager_reassigningpartitions{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('162', '11', '生产速率', 'kafka_server_brokertopicmetrics_bytesinpersec', 'instance', 'kafka_server_brokertopicmetrics_bytesinpersec{$INSTANCE}', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('163', '11', '消费速率', 'kafka_server_brokertopicmetrics_bytesoutpersec', 'instance', 'kafka_server_brokertopicmetrics_bytesoutpersec{$INSTANCE}', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('164', '11', '消息生产速率', 'kafka_server_brokertopicmetrics_messagesinpersec', 'instance', 'kafka_server_brokertopicmetrics_messagesinpersec{$INSTANCE}', null, null, '个/s', null, '2', '1', '1', null, null, null, null);
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('165', '11', '落后的消费量', 'kafka_server_fetcherlagmetrics_consumerlag', 'instance', 'kafka_server_fetcherlagmetrics_consumerlag{$INSTANCE}', null, null, '个', null, '2', '1', '1', null, null, null, null);

ALTER TABLE `t_monitor_item` ADD `display` VARCHAR (256) DEFAULT 'chart,rule,scaling' COMMENT '展示位置';

ALTER TABLE `t_monitor_product` ADD `sort` INT COMMENT '排序';

UPDATE t_monitor_product SET sort = 1 WHERE abbreviation = 'ecs';
UPDATE t_monitor_product SET sort = 2 WHERE abbreviation = 'eip';
UPDATE t_monitor_product SET sort = 3 WHERE abbreviation = 'slb';
UPDATE t_monitor_product SET sort = 4 WHERE abbreviation = 'cbr';
UPDATE t_monitor_product SET sort = 5 WHERE abbreviation = 'nat';
UPDATE t_monitor_product SET sort = 6 WHERE abbreviation = 'bms';
UPDATE t_monitor_product SET sort = 7 WHERE abbreviation = 'mysql';
UPDATE t_monitor_product SET sort = 8 WHERE abbreviation = 'dm';
UPDATE t_monitor_product SET sort = 9 WHERE abbreviation = 'postgresql';
UPDATE t_monitor_product SET sort = 10 WHERE abbreviation = 'kafka';
