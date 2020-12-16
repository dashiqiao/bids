package v1

import (
	"encoding/json"
	"platform_report/dao"
	"platform_report/pkg/util"
	"platform_report/service"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

/**
* @api {Get} /api/v1/report/support 获取报表支持的搜素条件
* @apiName api.Support
* @apiGroup report 报表
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} uuid 报表ID
*
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	"data":{
		"cond_id":1 , --条件ID
  		"parent_id": 0, --父级ID，为0则没有,注意>0时，val[json]新增 p 字段，为父级数据ID
		"name": "report_date",  --请求时参数名称
		"type": "select",       --类型 下拉框/输入框
		"val": "[{\"k\":\"昨天\",\"v\":\"yesterday\"},{\"k\":\"今天\",\"v\":\"day\"},{\"k\":\"本周\",\"v\":\"week\"},{\"k\":\"上周\",\"v\":\"lastWeek\"},{\"k\":\"本月\",\"v\":\"month\"},{\"k\":\"上月\",\"v\":\"lastMonth\"},{\"k\":\"本季度\",\"v\":\"quarter\"},{\"k\":\"上季度\",\"v\":\"lastQuarter\"},{\"k\":\"本年\",\"v\":\"year\"},{\"k\":\"去年\",\"v\":\"lastYear\"},{\"k\":\"自定义\",\"v\":\"custom\"}]",
		"placeholder": "时间范围" --提示输入内容
		}
}
*/
func Support(c *gin.Context) {
	uuid := c.Query("uuid")
	reportInvoke := service.NewReportService()
	datas, err := reportInvoke.GetReportConditionsByUuid(uuid)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}

	for i := 0; i < len(datas); i++ {
		if datas[i].Conduct == 2 {
			mapJ, _ := dao.NewInfiniteStones().War(datas[i].Val)
			mjson, _ := json.Marshal(mapJ)
			datas[i].Val = string(mjson)
		}
	}
	SuccessStd(c, datas)
}

/**
* @api {Get} /api/v1/report/header 获取报头
* @apiName api.Header
* @apiGroup report 报表
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} uuid 报表ID
*
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	"data":[
			{
			"eng_name": "channel_name", --英文字段
			"chs_name": "渠道名称"       --中文字段
			}]
	}
*/
func Header(c *gin.Context) {
	uuid := c.Query("uuid")
	reportInvoke := service.NewReportService()
	datas, err := reportInvoke.GetReportHeaderByUuid(uuid)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStd(c, datas)
}

/**
* @api {Post} /api/v1/report/index 获取数据
* @apiName api.Index
* @apiGroup report 报表
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} uuid 报表ID
* @apiParam {String} xxxx 其他参数...
* @apiParam {Int} p  当前页
* @apiParam {Int} limit 分页大小
* @apiParam {String} px 排序
*
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	"remark":"",--表底
	"data":[],  --明细数据
	"summary":[], --汇总数据
	"graph":[
  				{
                 "id" :"xxxxx",
                 "title" :"xxxxx", ---图表标题
            	 "figure": 1, --1.折线  2.饼图  3.柱状 4.漏斗
                 "x": [ "2020-03-31","2020-03-30"],
		"y": {
                "今日头条": [
                    "64",
                    "76"]
				}
     ],  --图表数据
	}
*/
func Index(c *gin.Context) {
	_, boolVal := c.GetPostFormMap("status")
	if boolVal {
	}
	m := c.Request.PostForm
	params := make(map[string]interface{})
	for k, v := range m {
		params[k] = v[0]
	}
	p, _ := strconv.Atoi(c.PostForm("p"))
	limit, _ := strconv.Atoi(c.PostForm("limit"))

	px := c.PostForm("px")

	if uuid, ok := params["uuid"]; ok {
		agent := new(service.AgentService)
		h := agent.Get(uuid.(string), params)
		if len(h) > 0 {
			c.JSON(200, h)
			return
		}

		reportInvoke := service.NewReportService()
		define, _ := reportInvoke.GetReportDefinesByUuid(uuid.(string))
		if define == nil {
			ErrorStd(c, "报表定义未找到")
			return
		}
		if define.DataSource > 0 {
			limit = p * limit
			params["limit"] = limit
		}
		p = (p - 1) * limit
		params["p"] = p

		links := new(dao.PlatformReport).GetReportLinkById(define.ReportId)
		childs := make([]*dao.ReportDefines, 0)
		if links != nil {
			for _, link := range links {
				d := reportInvoke.GetReportById(link.LinkId)
				if d == nil {
					continue
				}
				childs = append(childs, d)
			}

		}
		combine := service.NewCombineService(params, define)
		xmaps, xcount := combine.GetReportDataIndex(px, p, limit) //明细数据
		xsummary := make([]map[string]interface{}, 0)
		var xgraph []*service.GraphMap
		if p == 0 {
			xsummary = combine.GetReportSummaryIndex() //汇总数据
			xgraph = combine.GetReportGraphIndex()     //图表数据
		}
		h = gin.H{
			"defines": define,
			"links":   childs,
			"data":    xmaps,
			"action":  define.ReportAction,
			"summary": xsummary,
			"graph":   xgraph,
			"count":   xcount,
			"remark":  define.Remark,
			"note":    define.Note,
			"tip":     define.Tip,
			"code":    200,
		}
		agent.Save(uuid.(string), params, h)
		c.JSON(200, h)
	}
}

/**
  * @api {Post} /api/v1/report/export 数据导出
  * @apiName api.Export
  * @apiGroup report 报表
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 报表ID
  * @apiParam {String} xxxx 其他参数...
  * @apiParam {String} px 排序
  *
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200
}
*/
func Export(c *gin.Context) {
	_, boolVal := c.GetPostFormMap("status")
	if boolVal {
	}
	m := c.Request.PostForm
	params := make(map[string]interface{})
	for k, v := range m {
		params[k] = v[0]
	}

	px := c.PostForm("px")
	if uuid, ok := params["uuid"]; ok {
		reportInvoke := service.NewReportService()
		define, err := reportInvoke.GetReportDefinesByUuid(uuid.(string))
		if err != nil {
			ErrorStd(c, err.Error())
			return
		}
		sql := ""
		var maps []map[string]interface{}
		da := new(dao.DataAdapter)
		da.Adapter(define.DataSource)
		if define.ReportAuto == service.Smart {
			where := reportInvoke.GetWhere(uuid.(string), define.DataSource, params, true)
			sql = "select * from (" + define.Sql + ") as tab "
			if px != "" {
				px = " order by " + px
			}
			sql = sql + where + px + " limit 5000"
			maps, err = da.Do(sql, nil)
		} else {
			params["p"] = 0
			params["limit"] = 5000
			maps, err = da.Do(service.FormatSql(define.Sql, params))
		}
		//" limit 10000"
		if err != nil {
			ErrorStd(c, err.Error())
			return
		}
		c.Header("Content-type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment;filename="+define.ReportName+".xlsx")
		c.Header("Content-Transfer-Encoding", "binary")
		excel := excelize.NewFile()
		headers, _ := reportInvoke.GetReportHeaderByUuid(uuid.(string))
		sheetName := "Sheet1"
		sheet := excel.NewSheet(sheetName)
		style, _ := excel.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","color":["#FFFF00"] ,"pattern":1}}`)
		for idx, header := range headers {
			sheetPosition := Div(idx+1) + "1"
			excel.SetCellValue(sheetName, sheetPosition, header.ChsName)
			excel.SetCellStyle(sheetName, sheetPosition, sheetPosition, style)
			//fmt.Println(excel.GetColWidth("Sheet1",header.ChsName))
		}

		for lineNum, m := range maps {
			clumnNum := 0
			for _, header := range headers {
				ok := false
				for key, value := range m {
					if header.EngName != key {
						continue
					}
					ok = true
					clumnNum++
					sheetPosition := Div(clumnNum) + strconv.Itoa(lineNum+2)
					excel.SetCellValue(sheetName, sheetPosition, value)
					//excel.SetColWidth("Sheet1",sheetPosition,sheetPosition,64)
					//if value == nil || value == "" || value == sql {
					//	fmt.Println(66666666)
					//	value = "nbsp;"
					//}
					//value = value.(string)

					//switch value.(type) {
					//case string:
					//	excel.SetCellValue(sheetName, sheetPosition, value.(string))
					//	break
					//case int:
					//	excel.SetCellValue(sheetName, sheetPosition, value.(int))
					//	break
					//case float64:
					//	excel.SetCellValue(sheetName, sheetPosition, value.(float64))
					//	break
					//default:
					//	excel.SetCellValue(sheetName, sheetPosition, value)
					//}
				}

				if !ok {
					clumnNum++
					sheetPosition := Div(clumnNum) + strconv.Itoa(lineNum+2)
					excel.SetCellValue(sheetName, sheetPosition, "")
				}
			}
		}
		excel.SetActiveSheet(sheet)
		//fmt.Println(define.ReportName)
		//_ = excel.Write(c.Writer)
		//filename := define.ReportName + ".xlsx"
		if err := excel.Write(c.Writer); err != nil {
			c.JSON(200, gin.H{
				"code": 501,
				"msg":  err.Error(),
			})
			return
		}
		//c.FileAttachment(filename, filename)
		//c.JSON(200, gin.H{
		//	"code": 200,
		//})
	} else {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "请提供uuid",
		})
		return
	}
}

func TestExport(c *gin.Context) {
	// 列标题
	c.Header("Content-type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment;filename=图书.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	titles := []string{
		"姓名", "年龄", "性别",
	}
	// 数据源
	data := []map[string]interface{}{
		{"name": "jack", "age": "", "gender": "男"},
		{"name": "mary", "age": 28, "gender": "女"},
	}

	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")

	for clumnNum, v := range titles {
		sheetPosition := Div(clumnNum+1) + "1"

		f.SetCellValue("Sheet1", sheetPosition, v)
	}
	for lineNum, v := range data {
		// Set value of a cell.
		clumnNum := 0
		for _, vv := range v {
			clumnNum++
			sheetPosition := Div(clumnNum) + strconv.Itoa(lineNum+2)
			//f.SetCellValue("Sheet1", sheetPosition, vv)
			switch vv.(type) {
			case string:
				f.SetCellValue("Sheet1", sheetPosition, vv.(string))
				break
			case int:
				f.SetCellValue("Sheet1", sheetPosition, vv.(int))
				break
			case float64:
				f.SetCellValue("Sheet1", sheetPosition, vv.(float64))
				break
			}
		}
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save xlsx file by the given path.

	//if err := f.SaveAs("Book2.xlsx"); err != nil {
	//	println(err.Error())
	//}
	//c.FileAttachment()
	_ = f.Write(c.Writer)
	//c.FileAttachment("Book2.xlsx", "Book2.xlsx")
}

func CC(c *gin.Context) {
	reportInvoke := service.NewReportService()
	reportInvoke.CC()
	c.JSON(200, gin.H{
		"code": 200,
	})
}

func Div(Num int) string {
	var (
		Str  string = ""
		k    int
		temp []int //保存转化后每一位数据的值，然后通过索引的方式匹配A-Z
	)
	//用来匹配的字符A-Z
	Slice := []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	if Num > 26 { //数据大于26需要进行拆分
		for {
			k = Num % 26 //从个位开始拆分，如果求余为0，说明末尾为26，也就是Z，如果是转化为26进制数，则末尾是可以为0的，这里必须为A-Z中的一个
			if k == 0 {
				temp = append(temp, 26)
				k = 26
			} else {
				temp = append(temp, k)
			}
			Num = (Num - k) / 26 //减去Num最后一位数的值，因为已经记录在temp中
			if Num <= 26 { //小于等于26直接进行匹配，不需要进行数据拆分
				temp = append(temp, Num)
				break
			}
		}
	} else {
		return Slice[Num]
	}
	for _, value := range temp {
		Str = Slice[value] + Str //因为数据切分后存储顺序是反的，所以Str要放在后面
	}
	return Str
}

/**
  * @api {Get} /api/v1/report/dashboard 看板数据读取
  * @apiName api.DashBoard
  * @apiGroup report 报表
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 看板ID
  * @apiParam {String} date 时间 ==> yesterday ...
  * @apiParam {String} [time1] 时间
  * @apiParam {String} [time2] 时间
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  	"data":[
  	}
*/
func DashBoard(c *gin.Context) {
	dashboard := c.Query("uuid")
	date := c.Query("date")
	time1 := c.Query("time1")
	time2 := c.Query("time2")
	args := make(map[string]interface{})
	args["date"] = date
	args["time1"] = time1
	args["time2"] = time2

	agent := new(service.AgentService)
	h := agent.Get(dashboard, args)
	if len(h) > 0 {
		c.JSON(200, h)
		return
	}
	title, layout, number, targets, graph := service.NewTargetService().GetBoardTarget(dashboard)
	//fmt.Println(graph)
	b, _ := json.Marshal(targets)
	dst := make([]*dao.ReportTargetDefines, len(targets))
	_ = json.Unmarshal(b, &dst)

	layoutInfo, _ := service.NewLayoutService().GetOne(layout)
	layoutJson := ""
	if layoutInfo != nil {
		layoutJson = layoutInfo.Info
	}

	boardInvoke := service.NewBoardService()
	target := make(map[string]interface{}, 0)
	target["title"] = title
	t1, t2 := util.GetPkDays(time1, time2)
	pkItems := boardInvoke.GetTargetDashBoard(dst, date, t1, t2)
	items := boardInvoke.GetTargetDashBoard(targets, date, time1, time2)
	for i := 0; i < len(items); i++ {
		for j := 0; j < len(pkItems); j++ {
			if items[i].EngName == pkItems[j].EngName {
				items[i].PK = pkItems[j].Value
			}
		}
	}

	target["data"] = items
	target["number"] = number
	graphInvoke := service.NewCombineService(nil, nil)
	graphs := graphInvoke.GetReportGraphDash(graph, date, time1, time2)
	h = gin.H{
		"layout": layoutJson,
		"target": target,
		"graph":  graphs,
		"code":   200,
	}
	agent.Save(dashboard, args, h)
	c.JSON(200, h)
}
