DELETE FROM t_monitor_item WHERE biz_id IN ('220','221','222','223','224','225');

UPDATE t_monitor_product SET status = '0' WHERE abbreviation = 'cgw';
