package dao

import (
	"platform_report/lib"
)

type DefinesDao struct {
}

func (c *DefinesDao) Get(name string, p, limit int) ([]string, int64, error) {
	var
	(
		ids   []string
		err   error
		count int64
	)

	if name == "" {
		err = lib.InitXormMySql().Table(ReportDefines{}).Cols("uuid").Asc("in_time").Limit(limit, (p-1)*limit).Find(&ids)
		count, _ = lib.InitXormMySql().Table(ReportDefines{}).Count()
	} else {
		err = lib.InitXormMySql().Table(ReportDefines{}).Cols("uuid").Where("report_name like ?", "%"+name+"%").Cols("uuid").Asc("in_time").Limit(limit, (p-1)*limit).Find(&ids)
		count, _ = lib.InitXormMySql().Table(ReportDefines{}).Where("report_name like ?", "%"+name+"%").Count()
	}
	return ids, count, err
}

func (c *DefinesDao) GetMulti(ids []string) (ret []*ReportDefines, err error) {
	err = lib.InitXormMySql().In("uuid", ids).Find(&ret)
	return
}

func (c *DefinesDao) GetById(id int) *ReportDefines {
	item := new(ReportDefines)
	_, _ = lib.InitXormMySql().Where("report_id = ? ", id).Get(item)
	return item
}

func (c *DefinesDao) Exist(name string) int64 {
	count, _ := lib.InitXormMySql().Table(ReportDefines{}).Where("report_name = ?", name).Count()
	return count
}
