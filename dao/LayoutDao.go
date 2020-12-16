package dao

import (
	"platform_report/lib"
)

type ReportLayout struct {
	LayoutId int    `json:"layout_id" xorm:"not null pk autoincr INT(11)"`
	Title    string `json:"title" xorm:"VARCHAR(255)"`
	Info     string `json:"info" xorm:"TEXT"`
}

func NewReportLayout() *ReportLayout {
	return new(ReportLayout)
}

func (c *ReportLayout) Save() (int, error) {
	_, err := lib.InitXormMySql().Insert(c)
	return c.LayoutId, err
}

func (c *ReportLayout) Modify(id int, maps map[string]interface{}) error {
	_, err := lib.InitXormMySql().Table(c).ID(id).Update(maps)
	return err
}

func (c *ReportLayout) Delete(id int) error {
	_, err := lib.InitXormMySql().ID(id).Delete(c)
	return err
}

func (c *ReportLayout) Get() (ret []*ReportLayout, err error) {
	err = lib.InitXormMySql().Find(&ret)
	return
}
