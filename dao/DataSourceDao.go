package dao

import "platform_report/lib"

type ReportDataSource struct {
	DsId    int    `json:"ds_id" xorm:"pk autoincr"`
	DsName  string `json:"ds_name"`
	Catalog string `json:"catalog"`
	Schema  string `json:"schema"`
}

func (c *ReportDataSource) Get(id int) *ReportDataSource {
	item := new(ReportDataSource)
	_, _ = lib.InitXormMySql().Where("ds_id = ? ", id).Get(item)
	return item
}

func (c *ReportDataSource) GetMulti() (items []*ReportDataSource, err error) {
	err = lib.InitXormMySql().Find(&items)
	return
}

func (c *ReportDataSource) Create() (int64, error) {
	id, err := lib.InitXormMySql().Insert(c)
	return id, err
}

func (c *ReportDataSource) Modify(id int, maps map[string]interface{}) error {
	_, err := lib.InitXormMySql().Table(c).ID(id).Update(maps)
	return err
}

func (c *ReportDataSource) Delete(id int) error {
	_, err := lib.InitXormMySql().ID(id).Delete(c)
	return err
}
