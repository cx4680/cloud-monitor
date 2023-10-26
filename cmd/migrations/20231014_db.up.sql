UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})/1024/1024/1024' WHERE metric_name = 'ecs_filesystem_free_bytes';
UPDATE t_monitor_item SET metrics_linux = '(sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE}))/1024/1024/1024' WHERE metric_name = 'ecs_disk_used';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})/1024/1024/1024' WHERE metric_name = 'ecs_filesystem_size_bytes';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})/1024/1024/1024' WHERE metric_name = 'bms_filesystem_free_bytes';
UPDATE t_monitor_item SET metrics_linux = '(sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE}))/1024/1024/1024' WHERE metric_name = 'bms_disk_used';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})/1024/1024/1024' WHERE metric_name = 'bms_filesystem_size_bytes';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})/1024/1024/1024' WHERE metric_name = 'ebms_filesystem_free_bytes';
UPDATE t_monitor_item SET metrics_linux = '(sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE}))/1024/1024/1024' WHERE metric_name = 'ebms_disk_used';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})/1024/1024/1024' WHERE metric_name = 'ebms_filesystem_size_bytes';
