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
INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type) VALUES ('7', '弹性裸金属', '1', 'ebms', null, null, '/productmonitoring/ebms', '0 0 0/1 * * ?', 'http://bms-manage-bms-union.product-bms-union:8081', '/compute/ebms/ops/v1/tenants/{tenantId}/servers', 'bms', '12', '云产品监控');
