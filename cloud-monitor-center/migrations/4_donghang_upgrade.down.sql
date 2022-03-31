UPDATE t_monitor_product SET status = '0' WHERE abbreviation IN ('bms','kafka');

UPDATE t_monitor_item SET type = null WHERE biz_id IN ('68','69','70','71','72','73','74','75','76','77','78','79','80','81','82','83','84','85');

DELETE FROM t_monitor_item WHERE biz_id in ('166','167','168','169');