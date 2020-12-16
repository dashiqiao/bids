package dao

import (
	"platform_report/config"
	"platform_report/lib"
	"time"
)

type ReportDefines struct {
	DataSource      int         `json:"data_source"`
	Note            string      `json:"note"`
	Tip             string      `json:"tip"`
	ReportId        int         `json:"report_id" xorm:"autoincr"`
	Uuid            string      `json:"uuid"`
	ReportName      string      `json:"report_name"`
	Sql             string      `json:"sql"`
	SqlSummary      string      `json:"sql_summary"`
	SqlCount        string      `json:"sql_count"`
	ReportType      int         `json:"report_type"`
	ReportAction    int         `json:"report_action"`
	ReportAuto      int         `json:"report_auto"`
	CreateUserId    int         `json:"create_user_id"`
	CreateUserName  interface{} `xorm:"-" json:"create_user_name"`
	Remark          string      `json:"remark"`
	InTime          time.Time   `json:"in_time"`
	ReportAttribute int         `json:"report_attribute"`
	hc              int         `json:"hc" xorm:"-"`
}

type ReportConditions struct {
	CondId   int    `json:"cond_id"`
	ParentId int    `json:"parent_id"`
	CondName string `json:"name"`
	CondType string `json:"type"`
	CondVal  string `json:"-"`
	Holder   string `json:"placeholder"`
	Conduct  int    `json:"-"`
	IsChar   int    `json:"-"`
}

type ReportSupport struct {
	ReportId int
	CondId   int
}

type ReportLink struct {
	ReportId int `json:"-"`
	LinkId   int
}

type ReportHeader struct {
	ReportId int    `json:"-"`
	EngName  string `json:"eng_name"`
	ChsName  string `json:"chs_name"`
	Daily    int    `json:"daily"`
	ShowType int    `json:"show_type"`
	Link     bool   `json:"link"`
}

type ReportGraph struct {
	DataSource int       `json:"data_source"`
	Make       int       `json:"make"`
	Tip        string    `json:"tip"`
	Auto       int       `json:"auto"`
	ReportId   int       `json:"-"`
	GraphId    string    `json:"graph_id"`
	GraphName  string    `json:"graph_name"`
	Source     int       `json:"source"`
	Sql        string    `json:"sql"`
	GraphType  int       `json:"graph_type"`
	X          string    `json:"x"`
	YChsName   string    `json:"y_chs_name"`
	YEngName   string    `json:"y_eng_name"`
	InTime     time.Time `json:"in_time"`
}

type ReportGraphRelation struct {
	ReportId int    `json:"-"`
	GraphId  string `json:"graph_id"`
}

type ReportGraphSupport struct {
	GraphId string `json:"graph_id"`
	CondId  int    `json:"cond_id"`
}

type PlatformReport struct {
}

func NewPlatformReport() *PlatformReport {
	return new(PlatformReport)
}

func (c *PlatformReport) AddReportDefines(reportName, uuid string, createUserId int) error {
	report := new(ReportDefines)
	report.Uuid = uuid
	report.ReportName = reportName
	report.CreateUserId = createUserId
	report.ReportType = 1
	report.ReportAttribute = 1
	report.ReportAuto = 1
	report.DataSource = 0
	report.Remark = ""
	report.InTime = time.Now()
	_, err := lib.InitXormMySql().Insert(report)
	return err
}

func (c *PlatformReport) ModifyReportName(reportName, uuid string) error {
	_, err := lib.InitXormMySql().SQL("update report_defines set report_name=? where uuid=?", reportName, uuid).Execute()
	return err
}

func (c *PlatformReport) ModifyReportType(uuid string, ptype int) error {
	_, err := lib.InitXormMySql().SQL("update report_defines set report_type=? where uuid=?", ptype, uuid).Execute()
	return err
}

func (c *PlatformReport) ModifyReportAction(uuid string, action int) error {
	_, err := lib.InitXormMySql().SQL("update report_defines set report_action=? where uuid=?", action, uuid).Execute()
	return err
}

func (c *PlatformReport) ModifyReportAuto(uuid string, auto int) error {
	_, err := lib.InitXormMySql().SQL("update report_defines set report_auto=? where uuid=?", auto, uuid).Execute()
	return err
}

func (c *PlatformReport) ModifyReportSqlSummray(sqlsummary, uuid, tip string) error {
	_, err := lib.InitXormMySql().SQL("update report_defines set sql_summary=?,tip=? where uuid=?", sqlsummary, tip, uuid).Execute()
	return err
}

func (c *PlatformReport) ModifyReport(reportID int, sql, sqlcount, remark, note string, headers []ReportHeader, supports []ReportSupport, links []ReportLink, datasource, attribute int) error {
	session := lib.InitXormMySql().NewSession()
	defer session.Clone()
	err := session.Begin()

	_, err = session.SQL("update report_defines set `sql` = ?,sql_count=?,remark=?,note=?,data_source = ?,report_attribute= ? where report_id = ?", sql, sqlcount, remark, note, datasource, attribute, reportID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.SQL("delete from report_header where report_id = ?", reportID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.SQL("delete from report_support where report_id = ?", reportID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.SQL("delete from report_link where report_id = ?", reportID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}

	_, err = session.InsertMulti(headers)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.InsertMulti(supports)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.InsertMulti(links)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()
	return nil
}

func (c *PlatformReport) GetReportDefines() (datas []*ReportDefines, err error) {
	err = lib.InitXormMySql().Table("report_defines").OrderBy("report_id").Find(&datas)
	return
}

func (c *PlatformReport) GetReportConditions() (datas []*ReportConditions, err error) {
	err = lib.InitXormMySql().Table("report_conditions").Find(&datas)
	return
}

func (c *PlatformReport) GetReportSupport() (datas []*ReportSupport, err error) {
	err = lib.InitXormMySql().Table("report_support").Find(&datas)
	return
}

func (c *PlatformReport) GetReportHeader() (datas []*ReportHeader, err error) {
	err = lib.InitXormMySql().Table("report_header").Find(&datas)
	return
}

func (c *PlatformReport) GetTableList() ([]map[string]interface{}, error) {
	cf := config.Conf{}
	results, err := lib.InitXormMySql().SQL(`SELECT
							table_name
						  FROM information_schema.tables
						  WHERE table_schema = ? ORDER BY table_name `, cf.GetConf().DbDatabase).Query().List()
	return results, err
}

func (c *PlatformReport) GetTableListByType() ([]map[string]interface{}, error) {
	cf := config.Conf{}
	results, err := lib.InitXormMySql().SQL(`SELECT
							table_name
						  FROM information_schema.tables
						  WHERE table_schema = ? AND TABLE_TYPE = ?  ORDER BY table_name `, cf.GetConf().DbDatabase, "BASE TABLE").Query().List()
	return results, err
}

func (c *PlatformReport) GetTableColumn(tableName interface{}) ([]map[string]interface{}, error) {
	cf := config.Conf{}
	results, err := lib.InitXormMySql().SQL(`
									SELECT
									  column_name
									FROM information_schema.columns
									WHERE table_schema = ?
									AND table_name = ?;`, cf.GetConf().DbDatabase, tableName).Query().List()
	return results, err
}

//删除
func (c *PlatformReport) RemoveReportByID(reportID string) error {
	_, err := lib.InitXormMySql().Exec("delete from report_defines where uuid=?", reportID)
	return err
}

func (c *PlatformReport) AddGraph(graphID, graphName string, mk, auto int, sql string, gtype int, x, y1, y2, tip string, supports []ReportGraphSupport, datasource int) error {
	session := lib.InitXormMySql().NewSession()
	defer session.Clone()
	err := session.Begin()
	report := new(ReportGraph)
	report.ReportId = 0
	report.GraphId = graphID
	report.GraphName = graphName
	report.Source = 3
	report.Sql = sql
	report.GraphType = gtype
	report.X = x
	report.YEngName = y1
	report.YChsName = y2
	report.Tip = tip
	report.Auto = auto
	report.Make = mk
	report.DataSource = datasource
	report.InTime = time.Now()
	_, err = session.Insert(report)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.SQL("delete from report_graph_support where graph_id = ?", graphID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.InsertMulti(supports)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()

	return nil
}

func (c *PlatformReport) ModifyGraph(graphID, graphName, sql string, gtype, auto, mk int, x, y1, y2, tip string, supports []ReportGraphSupport, datasource int) error {
	report := new(ReportGraph)
	report.GraphName = graphName
	report.Sql = sql
	report.GraphType = gtype
	report.X = x
	report.YEngName = y1
	report.YChsName = y2
	report.Tip = tip
	report.Auto = auto
	report.Make = mk
	report.DataSource = datasource
	session := lib.InitXormMySql().NewSession()
	defer session.Clone()
	err := session.Begin()
	_, err = session.Where("graph_id = ?", graphID).Cols("graph_name", "sql", "graph_type", "x", "y_eng_name", "y_chs_name", "tip", "auto", "make", "data_source").Update(report)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.SQL("delete from report_graph_support where graph_id = ?", graphID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.InsertMulti(supports)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()

	return nil
}

func (c *PlatformReport) RemoveGraph(graphID string) error {
	_, err := lib.InitXormMySql().Exec("delete from report_graph where graph_id=?", graphID)
	_, _ = lib.InitXormMySql().Exec("delete from report_graph_support where graph_id = ?", graphID)
	_, _ = lib.InitXormMySql().Exec("delete FROM report_graph_relation WHERE graph_id = ?", graphID)
	return err
}

func (c *PlatformReport) GetGraph(reportId int) (graphs []*ReportGraph, err error) {
	err = lib.InitXormMySql().SQL(`  SELECT rg.* FROM report_graph_relation rgr INNER JOIN report_graph rg
    ON rgr.graph_id = rg.graph_id WHERE rgr.report_id= ?`, reportId).Find(&graphs)
	return
}

func (c *PlatformReport) GetAllGraph() (graphs []*ReportGraph, err error) {
	err = lib.InitXormMySql().Table("report_graph").Find(&graphs)
	return
}

func (c *PlatformReport) GetAllGraphByCondition(reportId int) (graphs []*ReportGraph, err error) {
	err = lib.InitXormMySql().SQL(`SELECT
											  *
											FROM report_graph rg
											WHERE rg.graph_id IN (SELECT
												rgs.graph_id
											  FROM report_graph_support rgs
											  WHERE rgs.cond_id IN (SELECT
												  rs.cond_id
												FROM report_support rs
												WHERE rs.report_id = ?))`, reportId).Find(&graphs)
	return
}

func (c *PlatformReport) GetGraphConditions() (datas []*ReportGraphSupport, err error) {
	err = lib.InitXormMySql().Find(&datas)
	return
}

func (c *PlatformReport) AddGraphRelation(reportId int, graphId string) error {
	item := new(ReportGraphRelation)
	item.ReportId = reportId
	item.GraphId = graphId
	_, err := lib.InitXormMySql().Insert(item)
	return err
}

func (c *PlatformReport) RemoveGraphRelation(reportId int, graphId string) error {
	_, err := lib.InitXormMySql().Exec("delete from report_graph_relation where report_id = ? and graph_id=?", reportId, graphId)
	return err
}

func (c *PlatformReport) GetAllForm() (items []*ReportDefines) {
	_ = lib.InitXormMySql().Where("report_attribute = 2").Find(&items)
	return
}

func (c *PlatformReport) GetReportLinkById(id int) (items []*ReportLink) {
	_ = lib.InitXormMySql().Where("report_id = ?", id).Find(&items)
	return
}
