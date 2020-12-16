package dao

import (
	"fmt"
	"platform_report/lib"
	"time"
)

type ReportAnalysis struct {
	AnalysisId   string    `json:"uuid"`
	Name         string    `json:"name"`
	Expr         string    `json:"expr"`
	Conditions   string    `json:"conditions"`
	Members      int       `json:"members"`
	CreateUserId int       `json:"-"`
	InTime       time.Time `json:"in_time"`
}

type AnalysisConditions struct {
	Op         string            `json:"op"`
	Values     string            `json:"values"`
	Flag       bool              `json:"flag"`
	Alias      string            `json:"alias"`
	DimFilters []AnalysisFilters `json:"dimFilters"`
	Id         string            `json:"id"`
	Range      []string          `json:"range"`
}

type AnalysisFilters struct {
	Dim string `json:"dim"`
	Op  string `json:"op"`
	//Values []string `json:"values"`
	Values string `json:"values"`
}

func (c *ReportAnalysis) AddAnalysis(id, name, expr, conditions string, members, createUserId int) error {
	report := new(ReportAnalysis)
	report.AnalysisId = id
	report.Name = name
	report.Conditions = conditions
	report.Expr = expr
	report.CreateUserId = createUserId
	report.Members = members
	report.InTime = time.Now()
	_, err := lib.InitXormMySql().Insert(report)

	return err
}

func (c *ReportAnalysis) ModifyAnalysis(id, name, expr, conditions string, members int) error {
	report := new(ReportAnalysis)
	report.AnalysisId = id
	report.Name = name
	report.Conditions = conditions
	report.Expr = expr
	report.Members = members
	_, err := lib.InitXormMySql().Where("analysis_id = ?", id).Cols("name", "conditions", "expr", "members").Update(report)
	return err
}

func (c *ReportAnalysis) ModifyAnalysisMembers(id string, members int) error {
	report := new(ReportAnalysis)

	report.Members = members
	_, err := lib.InitXormMySql().Where("analysis_id = ?", id).Cols("members").Update(report)

	return err
}

func (c *ReportAnalysis) RemoveAnalysis(id string) error {
	_, err := lib.InitXormMySql().Exec("delete from report_analysis where analysis_id=?", id)
	return err
}

func (c *ReportAnalysis) GetAnalysis(p, limit int) (items []*ReportAnalysis, total int64) {
	_ = lib.InitXormMySql().OrderBy("in_time").Limit(limit, (p-1)*limit).Find(&items)
	report := new(ReportAnalysis)
	total, _ = lib.InitXormMySql().Count(report)
	return
}

func (c *ReportAnalysis) GetAnalysisByID(id string) *ReportAnalysis {
	items := new(ReportAnalysis)
	_, _ = lib.InitXormMySql().Where("analysis_id = ?", id).Get(items)
	return items
}

func (c *ReportAnalysis) GetCountByName(name string) int64 {
	report := new(ReportAnalysis)
	total, _ := lib.InitXormMySql().Where("name = ?", name).Count(report)
	return total
}

func (c *ReportAnalysis) GetAnalysisActions(userId int, dt1, dt2 string, source, p, limit int) ([]map[string]interface{}, uint64) {
	//sql := fmt.Sprintf(`select *
	//							from user_behavior
	//							where user_id = %v
	//							  and created_at between '%v'
	//								and '%v'
	//							limit %v,%v`, userId, dt1, dt2, (p-1)*limit, limit)
	//sqlC := fmt.Sprintf(`select count(*) as members
	//							from user_behavior
	//							where user_id = %v
	//							  and created_at between '%v'
	//								and '%v'
	//							limit %v,%v`, userId, dt1, dt2, (p-1)*limit, limit)
	sql := " where 1=1 "
	if source == 1 {
		sql += " and operation_type = 1 "
	} else if source == 2 {
		sql += " and operation_type <> 1 "
	}
	sql += fmt.Sprintf(` and user_id = %v `, userId)
	if dt1 != "" {
		sql += fmt.Sprintf(`  and created_at between '%v' and '%v'`, dt1, dt2)
	}
	//sql += fmt.Sprintf(` limit %v,%v `, (p-1)*limit, limit)

	items, _ := NewInfiniteStones().CkQuery("select * from user_behavior_new " + sql + fmt.Sprintf(` limit %v,%v `, (p-1)*limit, limit))
	itemsC, _ := NewInfiniteStones().CkQuery("select count(*) as members from user_behavior_new " + sql)
	var total uint64
	if len(itemsC) > 0 {
		total = itemsC[0]["members"].(uint64)
	}
	return items, total
}
