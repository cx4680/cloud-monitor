package pageUtils

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"gorm.io/gorm"
)

func Paginate(pageSize int, current int, sql string, sqlParam []interface{}, list interface{}, db *gorm.DB) *vo.PageVO {
	var total int
	db.Raw("select count(1) from ( "+sql+") t ", sqlParam...).Scan(&total)
	if pageSize < 0 {
		pageSize = 10
	} else if pageSize > 5000 {
		pageSize = 5000
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages += 1
	}
	if current < 1 || pages == 0 {
		current = 1
	} else if current > pages {
		current = pages
	}

	offset := ((current) - 1) * (pageSize)
	sqlParam = append(sqlParam, pageSize)
	sqlParam = append(sqlParam, offset)
	db.Raw(sql+"  limit ? offset ?", sqlParam...).Find(list)

	return &vo.PageVO{
		Records: list,
		Current: current,
		Size:    pageSize,
		Total:   total,
		Pages:   pages,
	}
}
