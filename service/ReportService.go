package service

import (
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"platform_report/dao"
	"platform_report/lib"
	"regexp"
	"strings"
	"time"
)

type ReportService struct {
}

type Conditions struct {
	CondId   int    `json:"cond_id"`
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Val      string `json:"val"`
	Holder   string `json:"placeholder"`
	Conduct  int    `json:"-"`
	IsChar   int    `json:"-"`
}

const (
	REPORTDetail  = 1
	REPORTSummary = 2
	REPORTChart   = 4

	REPORTSource_1 = 1
	REPORTSource_2 = 2
	REPORTSource_3 = 3

	LINE          = 1
	PIE           = 2
	BAR           = 3
	FUNNEL        = 4
	HorizontalBar = 8
	XLine         = 9
	XBar          = 13

	Insert = 1
	Delete = 2
	Update = 4
	Select = 8

	Smart = 1
)

func NewReportService() *ReportService {
	return new(ReportService)
}

var kvcache = cache.New(10*time.Minute, 30*time.Minute)

func (c *ReportService) SetReportDataCache(key string, data *[]map[string]interface{}, count interface{}) {
	kvcache.Set(key, data, cache.DefaultExpiration)
	kvcache.Set(key+"-count", count, cache.DefaultExpiration)
}

func (c *ReportService) SetReportSummaryCache(key string, data []map[string]interface{}) {
	kvcache.Set("summary-"+key, data, cache.DefaultExpiration)
}

func (c *ReportService) GetReportSummaryCache(key string) ([]map[string]interface{}, error) {
	if x, found := kvcache.Get("summary-" + key); found {
		return x.([]map[string]interface{}), nil
	}
	return nil, errors.New("not fount")
}

func (c *ReportService) SetReportGraphCache(key string, data []map[string]interface{}) {
	kvcache.Set("graph-"+key, data, cache.DefaultExpiration)
}

func (c *ReportService) GetReportGraphCache(key string) ([]map[string]interface{}, error) {
	if x, found := kvcache.Get("graph-" + key); found {
		return x.([]map[string]interface{}), nil
	}
	return nil, errors.New("not fount")
}

func (c *ReportService) GetReportDataCache(key string) (*[]map[string]interface{}, interface{}, error) {
	var count interface{}
	if x, found := kvcache.Get(key + "-count"); found {
		count = x.(interface{})
	} else {
		return nil, 0, errors.New("....")
	}
	if x, found := kvcache.Get(key); found {
		return x.(*[]map[string]interface{}), count, nil
	} else {
		return nil, 0, errors.New("....")
	}
	//return nil, 0, errors.New("....")
}

func (c *ReportService) GetReportConditions() []*dao.ReportConditions {
	key := "report_conditions"
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportConditions)
	}
	datas, err := dao.NewPlatformReport().GetReportConditions()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *ReportService) GetReportConditionsById(id int) *dao.ReportConditions {
	conditions := c.GetReportConditions()
	for _, condition := range conditions {
		if condition.CondId == id {
			return condition
		}
	}
	return nil
}

func (c *ReportService) GetReportSelectConditions() []*dao.ReportConditions {
	conditions := c.GetReportConditions()
	ret := make([]*dao.ReportConditions, 0)
	if conditions != nil {
		for _, condition := range conditions {
			if condition.CondType == "select" {
				ret = append(ret, condition)
			}
		}
	}
	return ret
}

func (c *ReportService) getReportHeader() []*dao.ReportHeader {
	key := "report_header"
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportHeader)
	}
	datas, err := dao.NewPlatformReport().GetReportHeader()
	if err != nil {
		return nil
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *ReportService) getReportSupport() []*dao.ReportSupport {
	key := "report_support"
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportSupport)
	}
	datas, err := dao.NewPlatformReport().GetReportSupport()
	if err != nil {
		return nil
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *ReportService) CC() {
	kvcache.Flush()
}

func (c *ReportService) GetReportConditionsByUuid(uuid string) ([]Conditions, error) {
	defines := NewDefinesService().GetOne(uuid)
	ret := make([]Conditions, 0)
	if defines == nil {
		return ret, nil
	}
	supports := c.getReportSupport()
	conditions := c.GetReportConditions()

	for _, val := range supports {
		if val.ReportId == defines.ReportId {
			for _, cond := range conditions {
				if val.CondId == cond.CondId {
					ret = append(ret, Conditions{ParentId: cond.ParentId, CondId: cond.CondId, Name: cond.CondName, Type: cond.CondType, Val: cond.CondVal, Holder: cond.Holder, IsChar: cond.IsChar, Conduct: cond.Conduct})
				}
			}
		}
	}
	return ret, nil
}

func (c *ReportService) GetReportHeaderByUuid(uuid string) ([]*dao.ReportHeader, error) {
	defines := NewDefinesService().GetOne(uuid)
	if defines == nil {
		return nil, nil
	}
	header := c.getReportHeader()
	ret := make([]*dao.ReportHeader, 0)
	for _, val := range header {
		if val.ReportId == defines.ReportId {
			ret = append(ret, val)
		}
	}
	return ret, nil
}

func (c *ReportService) GetReportIdByUuid(uuid string) int {
	item, _ := c.GetReportDefinesByUuid(uuid)
	if item == nil {
		return 0
	}
	return item.ReportId
}

func (c *ReportService) GetReportDefinesByUuid(uuid string) (*dao.ReportDefines, error) {
	defines := NewDefinesService().GetOne(uuid)
	if defines == nil {
		return nil, errors.New("报表数据为空")
	}
	return defines, nil
}

func (c *ReportService) GetReportById(id int) *dao.ReportDefines {
	return new(dao.DefinesDao).GetById(id)
}

func (c *ReportService) GetReportTables() (*[]map[string]interface{}, error) {
	key := "tables"
	if x, found := kvcache.Get(key); found {
		return x.(*[]map[string]interface{}), nil
	}
	reportInvoke := dao.NewPlatformReport()
	tables, err := reportInvoke.GetTableList()
	if err != nil {
		return nil, err
	}
	ret := make([]map[string]interface{}, 0, 10)
	for _, val := range tables {
		//fmt.Println(val["table_name"])
		columns, err := reportInvoke.GetTableColumn(val["table_name"])
		if err != nil {
			return nil, err
		}
		item := make(map[string]interface{})
		item["table_name"] = val["table_name"]
		item["columns"] = columns
		ret = append(ret, item)
	}
	kvcache.Set(key, &ret, cache.DefaultExpiration)
	return &ret, nil
}

func (c *ReportService) AddReportTables(name string, createUserId int) error {
	if name == "" {
		return errors.New("请输入报表名称")
	}
	if new(dao.DefinesDao).Exist(name) > 0 {
		return errors.New("报表已存在,请重新输入")
	}
	uuid := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	invoke := dao.NewPlatformReport()
	err := invoke.AddReportDefines(name, uuid, createUserId)
	c.CC()
	return err
}

func (c *ReportService) ModifyReportName(name, uuid string) error {
	if name == "" {
		return errors.New("请输入报表名称")
	}
	//if new(dao.DefinesDao).Exist(name) > 1 {
	//	return errors.New("报表已存在,请重新输入")
	//}
	invoke := dao.NewPlatformReport()
	err := invoke.ModifyReportName(name, uuid)
	c.CC()
	new(AgentService).Clean(uuid)
	NewDefinesService().RemoveCache(uuid)
	return err
}

func (c *ReportService) ModifyReportType(uuid string, ptype int) error {
	invoke := dao.NewPlatformReport()
	err := invoke.ModifyReportType(uuid, ptype)
	c.CC()
	new(AgentService).Clean(uuid)
	NewDefinesService().RemoveCache(uuid)
	return err
}

func (c *ReportService) ModifyReportAction(uuid string, action int) error {
	invoke := dao.NewPlatformReport()
	err := invoke.ModifyReportAction(uuid, action)
	c.CC()
	new(AgentService).Clean(uuid)
	NewDefinesService().RemoveCache(uuid)
	return err
}

func (c *ReportService) ModifyReportAuto(uuid string, auto int) error {
	invoke := dao.NewPlatformReport()
	err := invoke.ModifyReportAuto(uuid, auto)
	c.CC()
	new(AgentService).Clean(uuid)
	NewDefinesService().RemoveCache(uuid)
	return err
}

func (c *ReportService) ModifyReportSqlSummary(uuid, sql, tip string, enable int) error {
	if sql == "" {
		return errors.New("请输入sql")
	}
	err := Valid(sql)
	if err != nil {
		return err
	}
	invoke := dao.NewPlatformReport()
	err = invoke.ModifyReportSqlSummray(sql, uuid, tip)
	if err == nil {
		defines, _ := c.GetReportDefinesByUuid(uuid)
		if defines != nil {
			if enable == 1 && (defines.ReportType&REPORTSummary) == 0 {
				ptype := defines.ReportType ^ REPORTSummary
				_ = c.ModifyReportType(uuid, ptype)
			}

			if enable == 0 && (defines.ReportType&REPORTSummary) > 0 {
				ptype := defines.ReportType ^ REPORTSummary
				_ = c.ModifyReportType(uuid, ptype)
			}
			//if (defines.ReportType & REPORTSummary) == 0 {
			//	ptype := defines.ReportType ^ REPORTSummary
			//	c.ModifyReportType(uuid, ptype)
			//}
		}
	}
	c.CC()
	new(AgentService).Clean(uuid)
	NewDefinesService().RemoveCache(uuid)
	return err
}

type Header struct {
	EngName  string `json:"eng_name"`
	ChsName  string `json:"chs_name"`
	Daily    int    `json:"daily"`
	ShowType int    `json:"show_type"`
	Link     bool   `json:"link"`
}

func (c *ReportService) ModifyReportTables(reportID int, sql, sqlcount, remark, tip string, headers []Header, supports, links []int, datasource, attribute int) error {
	if sql == "" {
		return errors.New("请输入Sql")
	}
	if reportID == 0 {
		return errors.New("报表ID错误")
	}
	err := Valid(sql)
	if err != nil {
		return err
	}
	rheaders := make([]dao.ReportHeader, 0, len(headers))
	rsupports := make([]dao.ReportSupport, 0, len(supports))
	rlinks := make([]dao.ReportLink, 0, len(links))
	for _, val := range headers {
		rheaders = append(rheaders, dao.ReportHeader{ReportId: reportID, EngName: val.EngName,
			ChsName: val.ChsName, Daily: val.Daily, ShowType: val.ShowType, Link: val.Link})
	}
	for _, val := range supports {
		rsupports = append(rsupports, dao.ReportSupport{ReportId: reportID, CondId: val})
	}
	for _, val := range links {
		rlinks = append(rlinks, dao.ReportLink{ReportId: reportID, LinkId: val})
	}
	invoke := dao.NewPlatformReport()

	err = invoke.ModifyReport(reportID, sql, sqlcount, remark, tip, rheaders, rsupports, rlinks, datasource, attribute)
	c.CC()
	uuid := c.GetReportById(reportID).Uuid
	new(AgentService).Clean(uuid)
	NewDefinesService().RemoveCache(uuid)
	return err
}

func Valid(sql string) error {
	sql = strings.ToLower(sql)
	ok1, _ := regexp.MatchString(`(delete\s)`, sql)
	ok2, _ := regexp.MatchString(`(update\s)`, sql)
	ok3, _ := regexp.MatchString(`(drop\s)`, sql)
	//fmt.Println(ok1, ok2, ok3)
	if ok1 || ok2 || ok3 {
		return errors.New("sql中包含危险关键字[drop,delete,update],只允许select")
	}
	return nil
}

func (c *ReportService) GetWhere(uuid string, dc int, params map[string]interface{}, needWhere bool) string {
	conditions, _ := c.GetReportConditionsByUuid(uuid)
	vData := ""
	//for _, val := range conditions {
	//	if val.Type == "date" {
	//		vData = val.Name
	//		if _, ok := params[val.Name]; !ok {
	//			params[val.Name] = "yesterday"
	//		}
	//	}
	//}

	//允许只传时间访问
	if _, ok := params["time1"]; ok {
		if vData != "" {
			params[vData] = "custom"
		}
	}

	where := ""
	if needWhere {
		where = " WHERE 1=1 "
	}
	for key, p := range params {
		for _, val := range conditions {
			if key == val.Name {
				if val.Type == "date" { //时间控件
					var time1, time2 string
					if p.(string) == "custom" {
						time1 = params["time1"].(string)
						time2 = params["time2"].(string) + " 23:59:59 "
					} else {
						t1, t2 := lib.GetTimeByType(p.(string))
						time1 = lib.TimeFormat(t1)
						time2 = lib.TimeFormat(t2) + " 23:59:59"
					}
					if dc == 0 {
						where += " AND " + key + " between " + "'" + time1 + "' AND " + "'" + time2 + "' "
					} else {
						where += " AND " + key + " between timestamp  " + "'" + time1 + "' AND " + " timestamp  '" + time2 + "' "
					}

				} else {
					if val.IsChar == 1 {
						where += " AND " + key + " like '%" + p.(string) + "%' "
					} else {
						where += " AND " + key + "=" + p.(string) + " "
					}
				}
			}
		}
	}
	return where
}

func (c *ReportService) GetWhere2(uuid string, dc int, vtype, dt1, dt2 string, needWhere bool) string {
	conditions := NewGraphService().GetGraphConditions(uuid)
	params := make(map[string]interface{})
	for _, val := range conditions {
		if val.Type == "date" {
			params[val.Name] = vtype
			params["time1"] = dt1
			params["time2"] = dt2
		}
	}
	where := ""
	if needWhere {
		where = " WHERE 1=1 "
	}
	for key, p := range params {
		for _, val := range conditions {
			if key == val.Name {
				if val.Type == "date" { //时间控件
					var time1, time2 string
					if p.(string) == "custom" {
						time1 = params["time1"].(string)
						time2 = params["time2"].(string) + " 23:59:59"
					} else {
						t1, t2 := lib.GetTimeByType(p.(string))
						time1 = lib.TimeFormat(t1)
						time2 = lib.TimeFormat(t2) + " 23:59:59 "
					}
					//where += " AND " + key + " between " + "'" + time1 + "' AND " + "'" + time2 + "'"
					if dc == 0 {
						where += " AND " + key + " between " + "'" + time1 + "' AND " + "'" + time2 + "' "
					} else {
						where += " AND " + key + " between timestamp  " + "'" + time1 + "' AND " + " timestamp  '" + time2 + "' "
					}

				} else {
					if val.IsChar == 1 {
						where += " AND " + key + " like '%" + p.(string) + "%' "
					} else {
						where += " AND " + key + "=" + p.(string) + " "
					}
				}
			}
		}
	}
	return where
}

func (c *ReportService) GetGraphWhere(uuid string, params map[string]interface{}, needWhere bool) string {
	conditions := NewGraphService().GetGraphConditions(uuid)
	//for _, val := range conditions {
	//	if val.Type == "date" {
	//		if _, ok := params[val.Name]; !ok {
	//			params[val.Name] = "yesterday"
	//		}
	//	}
	//}
	where := ""
	if needWhere {
		where = " WHERE 1=1 "
	}
	for key, p := range params {
		for _, val := range conditions {
			if key == val.Name {
				if val.Type == "date" { //时间控件
					var time1, time2 string
					if p.(string) == "custom" {
						time1 = params["time1"].(string)
						time2 = params["time2"].(string) + " 23:59:59 "
					} else {
						t1, t2 := lib.GetTimeByType(p.(string))
						time1 = lib.TimeFormat(t1)
						time2 = lib.TimeFormat(t2)
					}
					where += " AND " + key + " between " + "'" + time1 + "' AND " + "'" + time2 + "'"
				} else {
					if val.IsChar == 1 {
						where += " AND " + key + " like '%" + p.(string) + "%' "
					} else {
						where += " AND " + key + "=" + p.(string) + " "
					}
				}
			}
		}
	}
	return where
}
