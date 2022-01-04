package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)



func ProjectFetch(db *gorm.DB) ([]interface{}, []string, error) {
	var tables []interface{}
	var sqls []string

	//加载SQL
	sqlBytes, err := ioutil.ReadFile("scripts/region.sql")
	if err != nil {
		logger.Logger().Error("load sql file error", err)
		return nil, nil, err
	}
	sql := string(sqlBytes)
	if strutil.IsNotBlank(sql) {
		sqls = append(sqls, strings.Split(sql, ";")...)
	}

	return tables, sqls, nil
}
