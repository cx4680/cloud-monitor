package sys_db

import (
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

type DBInitializer struct {
	DB      *gorm.DB
	Fetches []InitializerFetch
}



func (i *DBInitializer) Initnitialization() error {
	for _, f := range i.Fetches {
		t, s, err := f(i.DB)
		if err != nil {
			return err
		}
		if e := i.DB.AutoMigrate(t...); e != nil {
			return e
		}
		for _, sql := range s {
			ns := strings.Replace(sql, "\n", "", -1)
			ns = strings.Replace(ns, "\r", "", -1)
			ns = strings.Replace(ns, "\t", "", -1)

			if strutil.IsNotEmpty(ns) {
				if e := i.DB.Exec(sql).Error; e != nil {
					return e
				}
			}
		}
	}
	return nil
}

type InitializerFetch func(db *gorm.DB) ([]interface{}, []string, error)

func CommonFetch(db *gorm.DB) ([]interface{}, []string, error) {
	var tables []interface{}
	var sqls []string
	//先将不需要保留的表删除
	if err := db.Migrator().DropTable(&commonModels.MonitorItem{}, &commonModels.MonitorProduct{}, &commonModels.ConfigItem{}); err != nil {
		return nil, nil, err
	}

	tables = append(tables, &commonModels.MonitorItem{}, &commonModels.MonitorProduct{}, &commonModels.ConfigItem{})

	tables = append(tables, &commonModels.AlertContact{})
	tables = append(tables, &commonModels.AlertContactGroup{})
	tables = append(tables, &commonModels.AlertContactGroupRel{})
	tables = append(tables, &commonModels.AlertContactInformation{})
	tables = append(tables, &commonModels.AlarmRule{})
	tables = append(tables, &commonModels.AlarmNotice{})
	tables = append(tables, &commonModels.AlarmInstance{})
	tables = append(tables, &commonModels.AlertRecord{})
	tables = append(tables, &commonModels.NotificationRecord{})
	tables = append(tables, &commonModels.ResourceGroup{}, &commonModels.ResourceResourceGroupRel{}, &commonModels.AlarmRuleGroupRel{}, &commonModels.AlarmRuleResourceRel{}, &commonModels.AlarmHandler{})

	//加载SQL
	sqlBytes, err := ioutil.ReadFile("scripts/common.sql")
	if err != nil {
		logger.Logger().Errorf("load sql file error:%v", err)
		return nil, nil, err
	}
	sql := string(sqlBytes)
	if strutil.IsNotBlank(sql) {
		sqls = append(sqls, strings.Split(sql, ";")...)
	}
	if !db.Migrator().HasTable(&commonModels.AlarmHandler{}) {
		//升级sql
		updateSqlBytes, err := ioutil.ReadFile("scripts/upgrade.sql")
		if err != nil {
			logger.Logger().Errorf("load update sql file error:%v", err)
			return nil, nil, err
		}
		updateSql := string(updateSqlBytes)
		if strutil.IsNotBlank(updateSql) {
			sqls = append(sqls, strings.Split(updateSql, ";")...)
		}
	}
	return tables, sqls, nil
}
