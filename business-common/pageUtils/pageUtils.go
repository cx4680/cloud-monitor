package pageUtils

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"unsafe"
)

func Paginate(pageSize int, current int, sql string, sqlParam []interface{}, retPoint unsafe.Pointer) *vo.PageVO {
	pageSize = correctionPageSize(pageSize)

	db := global.DB
	var total int
	db.Raw("select count(1) from ( "+sql+") t ", sqlParam...).Scan(&total)

	pages := total / pageSize
	if total%pageSize != 0 {
		pages += 1
	}
	if current < 1 || pages == 0 {
		current = 1
	} else if current > pages {
		current = pages
	}
	if total > 0 {
		offset := ((current) - 1) * (pageSize)
		sqlParam = append(sqlParam, pageSize)
		sqlParam = append(sqlParam, offset)
		db.Raw(sql+"  limit ? offset ?", sqlParam...).Find(retPoint)
	}

	return &vo.PageVO{
		Records: retPoint,
		Current: current,
		Size:    pageSize,
		Total:   total,
		Pages:   pages,
	}
}

func correctionPageSize(pageSize int) int {
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 5000 {
		pageSize = 5000
	}
	return pageSize
}
