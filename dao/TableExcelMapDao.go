package dao

import (
	"platform_report/lib"
)

type TableExcelMap struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	TableName       string `json:"table_name"`
	FieldMap        string `json:"field_map"`
	ReportId        string `json:"report_id"`
	IsImport        int    `json:"is_import"`
	IsTemplate      int    `json:"is_template"`
	TemplateAddress string `json:"template_address"`
}

func (this *TableExcelMap) InsertList(fields []string, data []string, tem TableExcelMap) error {
	// INSERT INTO table_name (列1, 列2,...) VALUES (值1, 值2,....)
	value := make([]interface{}, 0)

	// 头部
	head := "INSERT INTO " + tem.TableName

	// 中间
	middle := " ("
	for k, v := range fields {
		if k == 0 {
			middle += v
		} else {
			middle = middle + "," + v
		}
	}
	middle += ") "

	// 尾部
	tail := " VALUES ("
	for k, v := range data {
		if k == 0 {
			tail += "?"
		} else {
			tail = tail + ",?"
		}
		// v = strings.Replace(v, `"`, `\"`, -1)
		// v = strings.Replace(v, `'`, `\'`, -1)
		value = append(value, v)
	}
	tail += ")"

	sql := head + middle + tail

	exec := make([]interface{}, 0)
	exec = append(exec, sql)
	exec = append(exec, value...)

	_, err := lib.InitXormMySql().Exec(exec...)
	if err != nil {
		return err
	}
	// res.LastInsertId()
	return nil
}

func (this *TableExcelMap) GetInfoByNameAndReportId(name, reportId string) error {
	_, err := lib.InitXormMySql().Table("table_excel_map").Where("name = ? and report_id = ? ", name, reportId).Get(this)
	return err
}

func (this *TableExcelMap) GetInfoByName(name string) error {
	_, err := lib.InitXormMySql().Table("table_excel_map").Where("name = ? ", name).Get(this)
	return err
}

func (this *TableExcelMap) GetInfoById(id int) error {
	_, err := lib.InitXormMySql().Table("table_excel_map").Where("id = ?", id).Get(this)
	return err
}

func (this *TableExcelMap) GetInfoByReportId(reportId string) error {
	_, err := lib.InitXormMySql().Table("table_excel_map").Where("report_id = ?", reportId).Get(this)
	return err
}

func (this *TableExcelMap) CheckTableName(tableName string, id int) error {
	query := lib.InitXormMySql().Table("table_excel_map").Where("table_name = ?", tableName)

	if id != 0 {
		query.Where("id != ?", id)
	}

	_, err := query.Get(this)

	return err
}

func (this *TableExcelMap) CheckName(name string, id int) error {
	query := lib.InitXormMySql().Table("table_excel_map").Where("name = ? ", name)
	if id != 0 {
		query.Where("id != ?", id)
	}

	_, err := query.Get(this)
	return err
}

func (this *TableExcelMap) CheckReportId(reportId string, id int) error {
	_, err := lib.InitXormMySql().Table("table_excel_map").Where("report_id = ? and id != ?  ", reportId, id).Get(this)
	return err
}

func (this *TableExcelMap) Add() (int64, error) {
	id, err := lib.InitXormMySql().Insert(this)
	return id, err
}

func (this *TableExcelMap) GetList(name, tableName string, limit, page int) ([]TableExcelMap, error) {
	list := make([]TableExcelMap, 0)
	query := lib.InitXormMySql().Table("table_excel_map")
	if name != "" {
		query = query.Where("name = ? ", name)
	}
	if tableName != "" {
		query = query.Where("table_name = ? ", tableName)
	}

	err := query.Limit(limit, page).Find(&list)

	return list, err
}

func (this *TableExcelMap) GetListCount(name, tableName string) (int64, error) {
	query := lib.InitXormMySql().Table("table_excel_map")
	if name != "" {
		query = query.Where("name = ? ", name)
	}
	if tableName != "" {
		query = query.Where("table_name = ? ", tableName)
	}

	num, err := query.Count()

	return num, err
}

func (this *TableExcelMap) Modify() error {
	_, err := lib.InitXormMySql().Table("table_excel_map").Cols("name", "field_map", "table_name", "report_id", "is_import", "is_template",
		"template_address").Where("id = ? ",
		this.Id).Update(this)
	return err
}

func (this *TableExcelMap) Delete(id int) error {
	_, err := lib.InitXormMySql().Table("table_excel_map").Where("id = ? ", id).Delete(this)
	return err
}
