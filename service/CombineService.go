package service

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/orderedmap"
	"github.com/xormplus/core"
	"gopkg.in/flosch/pongo2.v3"
	"platform_report/dao"
	"platform_report/pkg/util"
	"regexp"
	"strings"
)

type CombineService struct {
	params   map[string]interface{}
	defines  *dao.ReportDefines
	cacheKey string
}

type GraphMap struct {
	ID     string                 `json:"id"`
	Figure int                    `json:"figure"`
	Title  string                 `json:"title"`
	X      []interface{}          `json:"x"`
	Y      *orderedmap.OrderedMap `json:"y"`
	Number int                    `json:"number"`
	Tip    string                 `json:"tip"`
	Source int                    `json:"source"`
	Uuid   string                 `json:"uuid"`
}

type YAxle struct {
}

func NewCombineService(p map[string]interface{}, d *dao.ReportDefines) *CombineService {
	mjson, _ := json.Marshal(p)
	key := string(mjson)
	return &CombineService{params: p, defines: d, cacheKey: key}
}

func (c *CombineService) GetReportDataIndex(px string, p, limit int) ([]map[string]interface{}, interface{}) {
	reportInvoke := NewReportService()
	da := new(dao.DataAdapter)
	da.Adapter(c.defines.DataSource)
	if c.defines.ReportAuto == Smart {
		where := reportInvoke.GetWhere(c.defines.Uuid, c.defines.DataSource, c.params, true)
		countsql := "select count(0) as num from (" + c.defines.Sql + ") as tab " + where
		count := da.Count(countsql, nil)
		sql := "select * from (" + c.defines.Sql + ") as tab "
		if px != "" {
			px = " order by " + px
		}
		if c.defines.DataSource > 0 {
			sql = sql + where + fmt.Sprintf(" and rn between %d and %d ", p, limit) + px
		} else {
			sql = sql + where + px + fmt.Sprintf(" limit %d,%d", p, limit)
		}

		fmt.Println(sql)
		maps, _ := da.Do(sql, nil)
		return maps, count
	} else {
		count := da.Count(FormatSql(c.defines.SqlCount, c.params))
		maps, _ := da.Do(FormatSql(c.defines.Sql, c.params))
		return maps, count
	}
}

func (c *CombineService) GetReportSummaryIndex() []map[string]interface{} {
	if c.defines.ReportType&REPORTSummary == 0 {
		return make([]map[string]interface{}, 0)
	}
	da := new(dao.DataAdapter)
	da.Adapter(c.defines.DataSource)
	reportInvoke := NewReportService()
	sql := ""
	if c.defines.ReportAuto == Smart {

		where := reportInvoke.GetWhere(c.defines.Uuid, c.defines.DataSource, c.params, false)
		if strings.Contains(strings.ToLower(c.defines.SqlSummary), "union all") {
			sqls := strings.Split(strings.ToLower(c.defines.SqlSummary), "union all")
			for i := 0; i < len(sqls); i++ {
				if strings.Contains(strings.ToLower(c.defines.SqlSummary), "where") {
					sqls[i] = sqls[i] + where
				} else {
					sqls[i] = sqls[i] + " where 1=1 " + where
				}
				sql = sql + sqls[i]
				if i != len(sqls)-1 {
					sql = sql + " union all "
				}
			}
		}

		//else if strings.Contains(strings.ToLower(c.defines.SqlSummary), "where") {
		//
		//	sql = c.defines.SqlSummary + where
		//} else {
		//	sql = c.defines.SqlSummary + " where 1=1 " + where
		//}
		sql = c.getGraphWhere(c.defines.SqlSummary, where)
		maps, _ := da.Do(sql, nil)
		return maps

	} else {
		maps, _ := da.Do(FormatSql(c.defines.SqlSummary, c.params))
		return maps
	}
}

func (c *CombineService) GetReportGraphIndex() []*GraphMap {
	ret := make([]*GraphMap, 0)
	if c.defines.ReportType&REPORTChart == 0 {
		return ret
	}
	graphs := NewGraphService().GetGraphList(c.defines.Uuid)
	if graphs == nil {
		return ret
	}
	reportInvoke := NewReportService()
	where := reportInvoke.GetWhere(c.defines.Uuid, c.defines.DataSource, c.params, false) //AND ....
	da := new(dao.DataAdapter)

	for _, graph := range graphs {
		da.Adapter(graph.DataSource)
		sql := graph.Sql
		var datasource []map[string]interface{}
		var err error
		if c.defines.ReportAuto == Smart {
			sql = c.getGraphWhere(sql, where)
			datasource, err = da.Do(sql, nil)
		} else {
			datasource, err = da.Do(FormatSql(sql, c.params))
		}
		if err != nil {
			fmt.Println(err.Error())
			return ret
		}
		gm := new(GraphMap)

		gm.Figure = graph.GraphType
		gm.Title = graph.GraphName
		gm.ID = graph.GraphId
		gm.Tip = graph.Tip
		x := make([]interface{}, 0)
		y := orderedmap.New()

		var engName, chsName []string
		_ = json.Unmarshal([]byte(graph.YEngName), &engName)
		_ = json.Unmarshal([]byte(graph.YChsName), &chsName)
		for _, column := range datasource {
			if X(graph.GraphType) {
				if value, ok := column[graph.X]; ok {
					x = append(x, value)
				}
			}
			if graph.Make == 0 {
				for i := 0; i < len(engName); i++ {
					if value, ok := column[engName[i]]; ok {
						c.appendToMap(chsName[i], value, y)
					}
				}
			} else {
				if value, ok := column[chsName[0]]; ok {
					c.appendToMap(util.ConversionString(value), column[engName[0]], y)
				}
			}

		}
		gm.X = x
		gm.Y = y

		ret = append(ret, gm)
	}
	return ret
}

func (c *CombineService) GetReportGraphPreviewIndex(title, sql string, mk, gtype, auto int, xName, yEngName, yChsName string, datasource int) []*GraphMap {
	ret := make([]*GraphMap, 0)
	//reportInvoke := NewReportService()
	var data []map[string]interface{}
	var err error
	da := new(dao.DataAdapter)
	da.Adapter(datasource)
	if auto == Smart {
		data, err = da.Do(sql+" limit 10 ", nil)
	} else {
		data, err = da.Do(FormatSql(sql, c.params))
	}

	if err != nil {
		fmt.Println(err.Error())
		return ret
	}
	gm := new(GraphMap)
	gm.Figure = gtype
	gm.Title = title
	x := make([]interface{}, 0)
	y := orderedmap.New()
	var engName, chsName []string
	_ = json.Unmarshal([]byte(yEngName), &engName)
	_ = json.Unmarshal([]byte(yChsName), &chsName)
	for _, column := range data {
		if X(gtype) {
			if value, ok := column[xName]; ok {
				x = append(x, value)
			}
		}
		if mk == 0 {
			for i := 0; i < len(engName); i++ {
				if value, ok := column[engName[i]]; ok {
					c.appendToMap(chsName[i], value, y)
				}
			}
		} else {
			if value, ok := column[chsName[0]]; ok {
				c.appendToMap(util.ConversionString(value), column[engName[0]], y)
			}
		}
	}
	gm.X = x
	gm.Y = y

	ret = append(ret, gm)
	return ret
}

func (c *CombineService) GetReportGraphDash(graphs []*dao.ReportBoardCard, vtype, time1, time2 string) []*GraphMap {
	ret := make([]*GraphMap, 0)
	reportInvoke := NewReportService()
	invoke := NewGraphService()
	da := new(dao.DataAdapter)
	for _, g := range graphs {
		gm := new(GraphMap)
		gm.Source = g.Source
		gm.Number = g.Number
		gm.Title = g.CardName
		if g.Source == 3 {
			gm.Uuid = g.ReportId
			ret = append(ret, gm)
			continue
		}
		graph := invoke.GetGraphByGraphId(g.GraphId)
		if graph == nil {
			continue
		}
		da.Adapter(graph.DataSource)

		where := reportInvoke.GetWhere2(g.GraphId, graph.DataSource, vtype, time1, time2, false) //AND ....
		sql := graph.Sql
		var datasource []map[string]interface{}
		var err error
		if graph.Auto == Smart {
			sql = c.getGraphWhere(sql, where)
			//fmt.Println("============", sql)
			datasource, err = da.Do(sql, nil)
		} else {
			datasource, err = da.Do(FormatSql(sql, c.params))
		}

		//fmt.Println("================================22222")
		if err != nil {
			fmt.Println(g.GraphId, sql, err.Error())
		}

		//fmt.Println("sssssssssss", datasource)
		gm.Figure = graph.GraphType
		gm.ID = graph.GraphId
		x := make([]interface{}, 0)
		y := orderedmap.New()
		var engName, chsName []string
		_ = json.Unmarshal([]byte(graph.YEngName), &engName)
		_ = json.Unmarshal([]byte(graph.YChsName), &chsName)
		for _, column := range datasource {
			if X(graph.GraphType) {
				if value, ok := column[graph.X]; ok {
					x = append(x, value)
				}
			}
			if graph.Make == 0 {
				for i := 0; i < len(engName); i++ {
					if value, ok := column[engName[i]]; ok {
						c.appendToMap(chsName[i], value, y)
					}
				}
			} else {
				if value, ok := column[chsName[0]]; ok {
					c.appendToMap(util.ConversionString(value), column[engName[0]], y)
				}
			}
		}
		gm.X = x
		gm.Y = y
		ret = append(ret, gm)
	}
	return ret
}

func (c *CombineService) autoCompletion(uuid string, datas []map[string]interface{}) []map[string]interface{} {
	if len(datas) == 0 {
		return datas
	}
	reportInvoke := NewReportService()
	headers, _ := reportInvoke.GetReportHeaderByUuid(uuid)

	for _, header := range headers {
		if _, ok := datas[0][header.EngName]; !ok {
			datas[0][header.EngName] = "-"
		}
	}
	return datas
}

func (c *CombineService) appendToMap(key string, value interface{}, maps *orderedmap.OrderedMap) *orderedmap.OrderedMap {
	if arr, ok := maps.Get(key); ok {
		newArr := arr.([]interface{})
		newArr = append(newArr, value)
		maps.Set(key, newArr)
	} else {
		maps.Set(key, []interface{}{value})
	}
	//if arr, ok := maps[key]; ok {
	//	arr = append(arr, value)
	//	maps[key] = arr
	//} else {
	//	data := make([]interface{}, 0)
	//	data = append(data, value)
	//	maps[key] = data
	//}
	return maps
}

func (c *CombineService) getGraphWhere(sql, where string) string {
	sql = strings.ToLower(sql)
	ok1, _ := regexp.MatchString(`(where)([\s\S]+)(order\s+by)`, sql) //同时包含 where order by
	ok2, _ := regexp.MatchString(`(where)([\s\S]+)(group\s+by)`, sql) //同时包含 where group by
	ok3, _ := regexp.MatchString(`(order\s+by)`, sql)                 //包含 order by
	ok4, _ := regexp.MatchString(`(group\s+by)`, sql)                 //只包含 group by
	ok5, _ := regexp.MatchString(`(where)`, sql)                      //只包含where
	ok6, _ := regexp.MatchString(`(limit)`, sql)                      //包含limit
	//fmt.Println(ok1, ok2, ok3, ok4, ok5)
	if ok1 {
		whereRegex := regexp.MustCompile(`(where)([\s\S]+)(order\s+by)([\s\S]+)`)
		whereSql := whereRegex.FindString(sql)
		ok, _ := regexp.MatchString(`(group\s+by)`, whereSql) //剩余包含group by
		if ok {
			groupRegex := regexp.MustCompile(`(group\s+by)([\s\S]+)`)
			groupSql := groupRegex.FindString(sql)
			//fmt.Println(groupSql)
			funny := strings.ReplaceAll(whereSql, groupSql, "") //原条件
			//fmt.Println(funny)
			sql = strings.ReplaceAll(sql, funny, funny+where)
		} else {
			groupRegex := regexp.MustCompile(`(order\s+by)([\s\S]+)`)
			groupSql := groupRegex.FindString(sql)
			//fmt.Println(groupSql)
			funny := strings.ReplaceAll(whereSql, groupSql, "")
			//fmt.Println(funny)
			sql = strings.ReplaceAll(sql, funny, funny+where)
			//sql = strings.ReplaceAll(sql, funny, where)
			//fmt.Println(funny1)
		}
	} else if ok2 {
		whereRegex := regexp.MustCompile(`(where)([\s\S]+)(group\s+by)([\s\S]+)`)
		whereSql := whereRegex.FindString(sql)
		ok, _ := regexp.MatchString(`(group\s+by)`, whereSql) //剩余包含group by
		if ok {
			groupRegex := regexp.MustCompile(`(group\s+by)([\s\S]+)`)
			groupSql := groupRegex.FindString(sql)
			//fmt.Println(groupSql)
			funny := strings.ReplaceAll(whereSql, groupSql, "")
			sql = strings.ReplaceAll(sql, funny, funny+where)
		}
	} else if ok4 {
		whereRegex := regexp.MustCompile(`(group\s+by)([\s\S]+)`)
		whereSql := whereRegex.FindString(sql)
		funny := strings.ReplaceAll(sql, whereSql, "")
		sql = funny + " where 1=1 " + where + whereSql
	} else if ok3 {
		whereRegex := regexp.MustCompile(`(order\s+by)([\s\S]+)`)
		whereSql := whereRegex.FindString(sql)
		funny := strings.ReplaceAll(sql, whereSql, "")
		sql = funny + " where 1=1 " + where + whereSql
		//select + where +order
	} else if ok5 {
		sql = sql + where
	} else if ok6 {
		whereRegex := regexp.MustCompile(`(limit)([\s\S]+)`)
		whereSql := whereRegex.FindString(sql)
		funny := strings.ReplaceAll(sql, whereSql, "")
		sql = funny + " where 1=1 " + where + whereSql
	} else {
		sql = sql + " where 1=1 " + where
	}
	//fmt.Println(sql)
	//return sql + " limit 30"
	return sql
}

func FormatSql(sql string, args map[string]interface{}) (string, []interface{}) {
	tpl, _ := pongo2.FromString(sql)
	out, err := tpl.Execute(pongo2.Context(args))
	if err != nil {
		fmt.Println(err)
	}
	sqlPar, sqlArgPar, err := core.MapToSlice(out, &args)
	if err != nil {
		fmt.Println(err)
	}
	return sqlPar, sqlArgPar
}

func X(gtype int) bool {
	return gtype == LINE || gtype == BAR || gtype == HorizontalBar || gtype == XLine || gtype == XBar
}
