ALTER TABLE t_monitor_product ADD COLUMN monitor_type varchar(50);

UPDATE t_monitor_product SET monitor_type = '云产品监控';

UPDATE t_monitor_product SET status = '1' WHERE abbreviation IN ('bms','kafka','dm','postgresql');

UPDATE t_monitor_item SET type = '1' WHERE biz_id IN ('68','69','70','71','72','73','74','75','76','77','78','79','80','81','82','83','84','85');

UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_tps{$INSTANCE}[1m])' WHERE biz_id = '120';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_qps{$INSTANCE}[1m])' WHERE biz_id = '121';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_ips{$INSTANCE}[1m])' WHERE biz_id = '122';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_dps{$INSTANCE}[1m])' WHERE biz_id = '123';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_ups{$INSTANCE}[1m])' WHERE biz_id = '124';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_ddlps{$INSTANCE}[1m])' WHERE biz_id = '125';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_nioips{$INSTANCE}[1m])' WHERE biz_id = '126';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_nio_ops{$INSTANCE}[1m])' WHERE biz_id = '127';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_fio_ips{$INSTANCE}[1m])' WHERE biz_id = '128';
UPDATE t_monitor_item SET metrics_linux = 'rate(dm_global_status_fio_ops{$INSTANCE}[1m])' WHERE biz_id = '129';
UPDATE t_monitor_item SET unit = 'ms' WHERE biz_id = '136';
UPDATE t_monitor_item SET unit = 'ms' WHERE biz_id = '137';
UPDATE t_monitor_item SET metric_name = 'kafka_consumergroup_lag', metrics_linux = 'kafka_consumergroup_lag{$INSTANCE}' WHERE biz_id = '165';

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('166', '9', '内存使用率', 'dm_global_status_mem_use_rate', 'instance', 'dm_global_status_mem_use_rate{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');

DELETE FROM t_monitor_product WHERE abbreviation = 'bms';
INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type) VALUES ('4', '传统裸金属', '1', 'bms', null, null, '/productmonitoring/bms', '0 0 0/1 * * ?', 'http://bms-manage-bms-union.product-bms-union:8082', '/compute/bms/ops/v1/tenants/{tenantId}/servers', 'bms', '11', '云产品监控');
INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type) VALUES ('7', '弹性裸金属', '1', 'ebms', null, null, '/productmonitoring/ebms', '0 0 0/1 * * ?', 'http://bms-manage-bms-union.product-bms-union:8081', '/compute/ebms/ops/v1/tenants/{tenantId}/servers', 'ebms', '12', '云产品监控');

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('167', '4', 'CPU使用率', 'bms_cpu_usage', 'instance', '100 - (100 * (sum by(instance) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))', null, null, '%', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('168', '4', 'CPU1分钟平均负载', 'bms_load1', 'instance', 'ecs_load1{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('169', '4', 'CPU5分钟平均负载', 'bms_load5', 'instance', 'ecs_load5{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('170', '4', 'CPU15分钟平均负载', 'bms_load15', 'instance', 'ecs_load15{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('171', '4', '内存使用量', 'bms_memory_used', 'instance', 'ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes', null, null, 'Byte', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('172', '4', '内存使用率', 'bms_memory_usage', 'instance', '100 * ((ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes) / ecs_memory_MemTotal_bytes)', null, null, '%', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('173', '4', '磁盘使用率', 'bms_disk_usage', 'instance,device', '100 * ((ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes) / ecs_filesystem_size_bytes)', null, null, '%', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('174', '4', '磁盘读速率', 'bms_disk_read_rate', 'instance,device', 'irate(ecs_disk_read_bytes_total{$INSTANCE}[3m])', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('175', '4', '磁盘写速率', 'bms_disk_write_rate', 'instance,device', 'irate(ecs_disk_written_bytes_total{$INSTANCE}[3m])', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('176', '4', '磁盘读IOPS', 'bms_disk_read_iops', 'instance,device', 'irate(ecs_disk_reads_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('177', '4', '磁盘写IOPS', 'bms_disk_write_iops', 'instance,device', 'irate(ecs_disk_writes_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('178', '4', '流入带宽', 'bms_network_transmit_rate', 'instance,device', 'irate(ecs_network_receive_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, 'Mbps', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('179', '4', '流出带宽', 'bms_network_receive_rate', 'instance,device', 'irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, 'Mbps', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('180', '4', '包接收速率', 'bms_network_receive_packets_rate', 'instance,device', 'irate(ecs_network_receive_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('181', '4', '包发送速率', 'bms_network_transmit_packets_rate', 'instance,device', 'irate(ecs_network_transmit_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('182', '4', '磁盘剩余存储量', 'bms_filesystem_free_bytes', 'instance,device', 'ecs_filesystem_free_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('183', '4', '磁盘已用存储量', 'bms_disk_used', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes', null, null, 'GB', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('184', '4', '磁盘存储总量', 'bms_filesystem_size_bytes', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null, null, 'chart,rule');
