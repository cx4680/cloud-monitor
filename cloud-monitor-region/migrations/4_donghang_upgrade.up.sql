UPDATE t_monitor_product SET status = '1' WHERE abbreviation IN ('bms','kafka');

UPDATE t_monitor_item SET type = '1' WHERE biz_id IN ('68','69','70','71','72','73','74','75','76','77','78','79','80','81','82','83','84','85');

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('166', '7', '物理机能耗', 'bms_total_power', 'instance', 'ipmi_sensor{$INSTANCE,name="total_power"}', null, null, 'W', null, '2', '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('167', '7', '物理机CPU能耗', 'bms_cpu_power', 'instance', 'ipmi_sensor{$INSTANCE,name="cpu_power"}', null, null, 'W', null, '2', '1', '1', null, null, null, null, 'chart');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('168', '7', '物理机MEM能耗', 'bms_memory_power', 'instance', 'ipmi_sensor{$INSTANCE,name="mem_power"}', null, null, 'W', null, '2', '1', '1', null, null, null, null, 'chart');

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

INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('169', '9', '内存使用率', 'dm_global_status_mem_use_rate', 'instance', 'dm_global_status_mem_use_rate{$INSTANCE}', null, null, '%', null, null, '1', '1', null, null, null, null, 'chart,rule');
