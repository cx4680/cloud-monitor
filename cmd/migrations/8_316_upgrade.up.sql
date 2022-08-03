UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Cpus{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '220';
UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Mems{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '221';
UPDATE t_monitor_item SET metrics_linux = 'ecs_procs_running{$INSTANCE}' WHERE biz_id = '222';
UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Fds{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '223';

UPDATE t_monitor_product SET page_url = '/inner/cmq/v1/kafka/cluster/list', status = '1' WHERE abbreviation IN ('kafka');
