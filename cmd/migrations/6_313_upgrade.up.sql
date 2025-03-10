DELETE FROM t_monitor_item WHERE biz_id IN ('149','150','151','152','153','154','155','156','157','158','159','160','161','162','163','164','165');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('149', '11', '在线Broker数', 'kafka_brokers', 'instance', 'sum by (instance) (kafka_brokers{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('150', '11', '分区的数量', 'kafka_server_replicamanager_partitioncount', 'instance', 'sum by (instance) (kafka_server_replicamanager_partitioncount{$INSTANCE})', null, null, '个', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('151', '11', '生产速率', 'kafka_server_brokertopicmetrics_bytesinpersec', 'instance', 'sum by (instance) (kafka_server_brokertopicmetrics_bytesinpersec{$INSTANCE})', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('152', '11', '消费速率', 'kafka_server_brokertopicmetrics_bytesoutpersec', 'instance', 'sum by (instance) (kafka_server_brokertopicmetrics_bytesoutpersec{$INSTANCE})', null, null, 'Byte/s', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('153', '11', '消息生产速率', 'kafka_server_brokertopicmetrics_messagesinpersec', 'instance', 'sum by (instance) (kafka_server_brokertopicmetrics_messagesinpersec{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('154', '11', '落后的消费量', 'kafka_consumergroup_lag', 'instance', 'sum by (instance) (kafka_consumergroup_lag{$INSTANCE})', null, null, '个/s', null, null, '1', '1', null, null, null, null, 'chart,rule');

UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_slave_io{$INSTANCE})' WHERE biz_id = '86';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_slave_sql{$INSTANCE})' WHERE biz_id = '87';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_slave_seconds_behind_master{$INSTANCE})' WHERE biz_id = '88';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_active_connections{$INSTANCE})' WHERE biz_id = '89';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_current_connection_percent{$INSTANCE})' WHERE biz_id = '90';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_qps{$INSTANCE})' WHERE biz_id = '91';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_tps{$INSTANCE})' WHERE biz_id = '92';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_select_ps{$INSTANCE})' WHERE biz_id = '93';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_update_ps{$INSTANCE})' WHERE biz_id = '94';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_insert_ps{$INSTANCE})' WHERE biz_id = '95';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_delete_ps{$INSTANCE})' WHERE biz_id = '96';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_cpu_usage{$INSTANCE})' WHERE biz_id = '97';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_mem_usage{$INSTANCE})' WHERE biz_id = '98';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_disk_usage{$INSTANCE})' WHERE biz_id = '99';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_select_ps{$INSTANCE})' WHERE biz_id = '100';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_update_ps{$INSTANCE})' WHERE biz_id = '101';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_insert_ps{$INSTANCE})' WHERE biz_id = '102';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_delete_ps{$INSTANCE})' WHERE biz_id = '103';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_cache_hit_rate{$INSTANCE})' WHERE biz_id = '104';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_reads_ps{$INSTANCE})' WHERE biz_id = '105';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_writes_ps{$INSTANCE})' WHERE biz_id = '106';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_buffer_pool_pages_dirty{$INSTANCE})' WHERE biz_id = '107';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_buffer_pool_bytes_dirty{$INSTANCE})' WHERE biz_id = '108';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_innodb_log_waits{$INSTANCE})' WHERE biz_id = '109';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_binlog_cache_disk_use{$INSTANCE})' WHERE biz_id = '110';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_slow_queries_per_min{$INSTANCE})' WHERE biz_id = '111';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_long_query_count{$INSTANCE})' WHERE biz_id = '112';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_long_query_alert_count{$INSTANCE})' WHERE biz_id = '113';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_exec_statememt_frequency{$INSTANCE})' WHERE biz_id = '114';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_read_frequency{$INSTANCE})' WHERE biz_id = '115';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_write_frequency{$INSTANCE})' WHERE biz_id = '116';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_top_statememt_avg_exec_time{$INSTANCE})' WHERE biz_id = '117';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_top_statememt_exec_err_rate{$INSTANCE})' WHERE biz_id = '118';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance)(mysql_current_cons_num{$INSTANCE})' WHERE biz_id = '119';

UPDATE t_monitor_item SET metric_name = 'ebms_cpu_usage' WHERE biz_id = '68';
UPDATE t_monitor_item SET metric_name = 'ebms_load1' WHERE biz_id = '69';
UPDATE t_monitor_item SET metric_name = 'ebms_load5' WHERE biz_id = '70';
UPDATE t_monitor_item SET metric_name = 'ebms_load15' WHERE biz_id = '71';
UPDATE t_monitor_item SET metric_name = 'ebms_memory_used' WHERE biz_id = '72';
UPDATE t_monitor_item SET metric_name = 'ebms_memory_usage' WHERE biz_id = '73';
UPDATE t_monitor_item SET metric_name = 'ebms_disk_usage' WHERE biz_id = '74';
UPDATE t_monitor_item SET metric_name = 'ebms_disk_read_rate' WHERE biz_id = '75';
UPDATE t_monitor_item SET metric_name = 'ebms_disk_write_rate' WHERE biz_id = '76';
UPDATE t_monitor_item SET metric_name = 'ebms_disk_read_iops' WHERE biz_id = '77';
UPDATE t_monitor_item SET metric_name = 'ebms_disk_write_iops' WHERE biz_id = '78';
UPDATE t_monitor_item SET metric_name = 'ebms_network_receive_rate' WHERE biz_id = '79';
UPDATE t_monitor_item SET metric_name = 'ebms_network_transmit_rate' WHERE biz_id = '80';
UPDATE t_monitor_item SET metric_name = 'ebms_network_receive_packets_rate' WHERE biz_id = '81';
UPDATE t_monitor_item SET metric_name = 'ebms_network_transmit_packets_rate' WHERE biz_id = '82';
UPDATE t_monitor_item SET metric_name = 'ebms_filesystem_free_bytes' WHERE biz_id = '83';
UPDATE t_monitor_item SET metric_name = 'ebms_disk_used' WHERE biz_id = '84';
UPDATE t_monitor_item SET metric_name = 'ebms_filesystem_size_bytes' WHERE biz_id = '85';

UPDATE t_monitor_item SET metric_name = 'bms_network_receive_rate' WHERE biz_id = '178';
UPDATE t_monitor_item SET metric_name = 'bms_network_transmit_rate' WHERE biz_id = '179';

INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('220', '1', '进程CPU使用率', 'ecs_processes_top5Cpus', 'instance,cmd_line', 'sum by(instance,cmd_line)(ecs_processes_top5Cpus{cmd_line!="",$INSTANCE})', NULL, NULL, '%', NULL, 3, 1, 1, NULL, NULL, NULL, '{{ne .OSTYPE \"windows\"}}', 'chart');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('221', '1', '进程内存使用率', 'ecs_processes_top5Mems', 'instance,cmd_line', 'sum by(instance,cmd_line)(ecs_processes_top5Mems{cmd_line!="",$INSTANCE})', NULL, NULL, '%', NULL, 3, 1, 1, NULL, NULL, NULL, '{{ne .OSTYPE \"windows\"}}', 'chart');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('222', '1', '运行中进程数', 'ecs_procs_running', 'instance', 'sum by(instance)(ecs_procs_running{$INSTANCE})', NULL, NULL, '个', NULL, 3, 1, 1, NULL, NULL, NULL, '{{ne .OSTYPE \"windows\"}}', 'chart');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('223', '1', '进程打开文件数', 'ecs_processes_top5Fds', 'instance,cmd_line', 'sum by(instance,cmd_line)(ecs_processes_top5Fds{cmd_line!="",$INSTANCE})', NULL, NULL, '个', NULL, 3, 1, 1, NULL, NULL, NULL, '{{ne .OSTYPE \"windows\"}}', 'chart');

UPDATE t_monitor_product SET status = '1' WHERE abbreviation IN ('cgw','bms','ebms');

INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('224', '1', '上行流量', 'tenant_network_send_rate_sum', '', 'sum(eip_upstream_bits_rate{$INSTANCE})/8/1024', NULL, NULL, 'KB/s', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, NULL);
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('225', '1', '下行流量', 'tenant_network_receive_rate_sum', '', 'sum(eip_downstream_bits_rate{$INSTANCE})/8/1024', NULL, NULL, 'KB/s', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, NULL);

INSERT INTO `t_config_item` (`biz_id`, `p_biz_id`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('24', '21', '磁盘使用率（操作系统）', NULL, 'ecs_disk_usage', 2, NULL);

UPDATE t_monitor_item SET metrics_linux = 'avg by(instance)(dm_global_status_mem_used{$INSTANCE})' WHERE biz_id = '130';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,service,route)(guard_http_apirequests{$INSTANCE})' WHERE biz_id = '218';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,service)(guard_http_apirequests{$INSTANCE})' WHERE biz_id = '219';

UPDATE t_monitor_item SET type = '1' WHERE biz_id IN ('167','168','169','170','171','172','173','174','175','176','177','178','179','180','181','182','183','184');
