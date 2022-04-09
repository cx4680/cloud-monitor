ALTER TABLE t_monitor_product DROP COLUMN monitor_type;

UPDATE t_monitor_product SET status = '0' WHERE abbreviation IN ('bms','kafka','dm','postgresql');

UPDATE t_monitor_item SET type = null WHERE biz_id IN ('68','69','70','71','72','73','74','75','76','77','78','79','80','81','82','83','84','85');

DELETE FROM t_monitor_item WHERE biz_id IN ('166','167','168','169','170','171','172','173','174','175','176','177','178','179','180','181','182','183','184');

DELETE FROM t_monitor_product WHERE abbreviation IN ('bms','ebms');
INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation) VALUES ('7', '裸金属服务器', '0', 'bms', null, null, '/productmonitoring/bms', '0 0 0/1 * * ?', 'http://bms-manage.product-bms:8080', '/compute/bms/ops/v1', 'bms');
