package service

import (
	"bytes"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"platform_report/dao"
	"platform_report/lib"
	"platform_report/pkg/util"
	"strings"
	"time"
)

type AnalysisService struct {
}

func NewAnalysisService() *AnalysisService {
	return new(AnalysisService)
}

func (c *AnalysisService) Add(name, expr, conditions string, members, createUserId int) error {
	ID := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	if c.GetCountByName(name) > 0 {
		return errors.New("分群名称已存在,请重新输入")
	}
	err := new(dao.ReportAnalysis).AddAnalysis(ID, name, expr, conditions, members, createUserId)
	return err
}

func (c *AnalysisService) Modify(id, name, expr, conditions string, members int) error {
	if c.GetCountByName(name) > 1 {
		return errors.New("分群名称已存在,请重新输入")
	}

	err := new(dao.ReportAnalysis).ModifyAnalysis(id, name, expr, conditions, members)
	return err
}

func (c *AnalysisService) ModifyAnalysisMembers(id string, members int) error {
	err := new(dao.ReportAnalysis).ModifyAnalysisMembers(id, members)
	return err
}

func (c *AnalysisService) Remove(id string) error {
	err := new(dao.ReportAnalysis).RemoveAnalysis(id)
	return err
}

func (c *AnalysisService) GetAnalysis(p, limit int) ([]*dao.ReportAnalysis, int64) {
	return new(dao.ReportAnalysis).GetAnalysis(p, limit)
}

func (c *AnalysisService) GetAnalysisByID(id string) *dao.ReportAnalysis {
	return new(dao.ReportAnalysis).GetAnalysisByID(id)
}

func (c *AnalysisService) GetCountByName(name string) int64 {
	return new(dao.ReportAnalysis).GetCountByName(name)
}

func (c *AnalysisService) GetAnalysisActions(userId int, dt1, dt2 string, source, p, limit int) ([]map[string]interface{}, uint64) {
	return new(dao.ReportAnalysis).GetAnalysisActions(userId, dt1, dt2, source, p, limit)
}

func (c *AnalysisService) FormatSql(expr string, conditions []dao.AnalysisConditions) string {
	var buffer bytes.Buffer
	for _, item := range conditions {
		val := FormatConditionValue(item.Id)
		//dt1, dt2 := FormatRange(item.Range)
		if item.Flag {
			buffer.WriteString(fmt.Sprintf(" user_id in (select user_id from user_behavior_new where "))
			//if val != 0 {
			buffer.WriteString(fmt.Sprintf(" operation_type = %v", val))
			buffer.WriteString(fmt.Sprintf(" and created_at between '%s' and '%s'", item.Range[0], item.Range[1]))
			for _, filter := range item.DimFilters {
				buffer.WriteString(BuilderFilter(filter))
			}
			buffer.WriteString(" group by user_id,operation_type having count(*) ")
			//if len(item.Values) == 1 {
			buffer.WriteString(item.Op)
			//buffer.WriteString(strconv.Itoa(val))
			buffer.WriteString(item.Values)
			//} else {
			//	buffer.WriteString(fmt.Sprintf(" between  %v and %v ", item.Values[0], item.Values[1]))
			//}
			//} else {
			//	buffer.WriteString(fmt.Sprintf(" created_at between '%s' and '%s'", item.Range[0], item.Range[1]))
			//	for _, filter := range item.DimFilters {
			//		buffer.WriteString(BuilderFilter(filter))
			//	}
			//	buffer.WriteString(" group by user_id,created_at having count(*) ")
			//if len(item.Values) == 1 {
			//buffer.WriteString(item.Op)
			//buffer.WriteString(strconv.Itoa(val))
			//} else {
			//	buffer.WriteString(fmt.Sprintf(" between  %v and %v ", item.Values[0], item.Values[1]))
			//}
			//}
			buffer.WriteString(" ) ")
		} else {
			buffer.WriteString(fmt.Sprintf(" user_id not in (select user_id from user_behavior_new where operation_type = %v and created_at between '%s' and '%s' group by user_id ) ", val, item.Range[0], item.Range[1]))
		}
		expr = strings.ReplaceAll(expr, item.Alias, buffer.String())
		buffer.Reset()
	}
	return " from user_behavior_new where " + expr
}

func BuilderFilter(filter dao.AnalysisFilters) string {
	var buffer bytes.Buffer
	buffer.WriteString(" and ")
	buffer.WriteString(filter.Dim)
	buffer.WriteString("  ")
	buffer.WriteString(filter.Op)
	buffer.WriteString("  ")

	if filter.Op == "in" || filter.Op == "not in" {
		if filter.Dim == "resource_title" {
			buffer.WriteString("('" + filter.Values + "')")
		} else {
			buffer.WriteString("(" + filter.Values + ")")
		}
	} else if filter.Op == "like" || filter.Op == "not like" {
		buffer.WriteString("'%" + filter.Values + "%'")
	} else {
		if filter.Dim == "resource_title" {
			buffer.WriteString("'" + filter.Values + "'")
		} else {
			buffer.WriteString(filter.Values)
		}

	}
	return buffer.String()
}

func FillCountHeader(sql string) string {
	return "select count(distinct user_id) as members " + sql
}

func FillListHeader(sql string, p, limit int) string {
	//and created_at between addDays(now(), -30) and now() and operation_type= 1
	return " select user_id,0 as times,max(created_at) as visit " + sql + fmt.Sprintf("   group by  user_id limit %v,%v", (p-1)*limit, limit)
}

func (c *AnalysisService) GetTotalMembers(sql string) uint64 {
	maps, _ := dao.NewInfiniteStones().CkQuery(sql)
	if len(maps) > 0 {
		return maps[0]["members"].(uint64)
	}
	return 0
}

func (c *AnalysisService) GetVisitTimes(userId interface{}) uint64 {
	//and operation_type = 1
	now := time.Now()
	m, _ := time.ParseDuration("-720h")
	m1 := now.Add(m)
	maps, _ := dao.NewInfiniteStones().CkQuery(`select count(0) as times
        from user_behavior_new
        where user_id= ` + util.ConversionString(userId) + `  and created_at between '` + lib.DateFormat(m1) + `'  and '` + lib.DateFormat(now) + `'
          `)
	if len(maps) > 0 {
		return maps[0]["times"].(uint64)
	}
	return 0
}

func FormatRange(r string) (string, string) {
	times := strings.Split(r, ",")
	if len(times) >= 2 {
		return times[0], times[1]
	}
	return "", ""
}

func FormatArray(r []string) string {
	str := ""
	for _, item := range r {
		str += "'" + item + "' ,"
	}
	return str[0 : len(str)-1]
}

func FormatConditionValue(op string) int {
	switch op {
	case "view":
		return 1
	case "fav":
		return 2
	case "share":
		return 3
	case "good":
		return 4
	case "comment":
		return 5
	case "reply":
		return 6
	case "award":
		return 7
	case "distribute":
		return 8
	case "order":
		return 0
	default:
		return 0
	}
}
