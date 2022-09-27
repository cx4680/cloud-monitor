UPDATE t_monitor_product SET page_url = '/slb/inner/monitor/list' WHERE abbreviation = 'slb';
UPDATE t_monitor_product SET iam_page_url = '/slb/list' WHERE abbreviation IN ('slb');
UPDATE t_monitor_product SET iam_page_url = '/nat-gw/nat/page' WHERE abbreviation IN ('nat');
UPDATE t_monitor_product SET iam_page_url = '/nat-gw/nat/page' WHERE abbreviation IN ('nat-e');
