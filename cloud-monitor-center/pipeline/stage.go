package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

type ProjectInitializerFetch struct {
}

func (p *ProjectInitializerFetch) Fetch(db *gorm.DB) ([]interface{}, []string, error) {
	var tables []interface{}
	var sqls []string

	//加载SQL
	sqlBytes, err := ioutil.ReadFile("script/center.sql")
	if err != nil {
		logger.Logger().Error("load sql file error", err)
		return nil, nil, err
	}
	sql := string(sqlBytes)
	if tools.IsNotBlank(sql) {
		sqls = append(sqls, strings.Split(sql, ";")...)
	}

	return tables, sqls, nil
}
