package pipeline

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
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
		log.Println("load sql file error", err)
		return nil, nil, err
	}
	sql := string(sqlBytes)
	if tools.IsNotBlank(sql) {
		sqls = append(sqls, strings.Split(sql, ";")...)
	}

	return tables, sqls, nil
}
