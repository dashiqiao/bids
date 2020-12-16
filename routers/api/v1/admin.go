package v1

import (
	"encoding/json"
	"io/ioutil"
	"platform_report/dao"
	"platform_report/service"
	"strconv"

	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
)

/**
* @api {Get} /admin/v1/report/tablelist 获取数据表结构
* @apiName api.TableList
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
*
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
		"code": 200,
		"msg": "ok"
		"data":
			{
				"columns": [
					{
						"column_name": "id"  --列名
					}],
				"table_name": "dxz_course_languages_statistics_info" --表名
			}
      }
*/
func TableList(c *gin.Context) {
	invoke := service.NewReportService()
	ret, err := invoke.GetReportTables()
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStd(c, ret)
}

/**
* @api {Get} /admin/v1/report/index 报表列表
* @apiName api.ReportDefines
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} [title] 模糊搜索标题
* @apiParam {Int} p  当前页
* @apiParam {Int} limit 分页大小
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	"data":
		[
			{
			"report_id": 1, --报表ID
			"uuid": "885a3868675711eab04494c69172a40b", --报表ID
			"report_name": "测试1", --报表名称
			"sql": "SELECT * FROM test ",
			"in_time": "2020-03-16T11:18:36+08:00"
			}
		]
	}
*/
func ReportDefines(c *gin.Context) {
	title := c.Query("title")
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	items, count := service.NewDefinesService().Get(title, p, limit)
	invoke := service.NewUserService()
	for i := 0; i < len(items); i++ {
		items[i].CreateUserName = invoke.GetUserName(items[i].CreateUserId)
	}

	c.JSON(200, gin.H{
		"code":  200,
		"data":  items,
		"count": count,
	})
}

/**
* @api {Post} /admin/v1/report/save 添加报表
* @apiName api.AddReport
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} name 报表名称
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
}
*/
func AddReport(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	name := fsion.Get("name").ValueStr()
	invoke := service.NewReportService()
	err := invoke.AddReportTables(name, UserInfo(c).ID)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/report/modify 修改报表
* @apiName api.ModifyReport
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} json 格式在下面
* @apiSuccessExample Success-TXResponse:
* HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	}
	--json
	{
		"report_id":1 --报表ID
		"note":"", --备注
		"remark":"", --备注
		"sql":"",  --明细sql
        "sql_count":"",  --集合sql
		"headers":[
						{ "eng_name":"xxx","chs_name":"xxx"}
				  ], --表头
		"supports":[1,2,3,4]  --条件ID
	}
*/
func ModifyReport(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	reportID := fsion.Get("report_id").ValueInt()
	sql := fsion.Get("sql").ValueStr()
	sqlCount := fsion.Get("sql_count").ValueStr()
	remark := fsion.Get("remark").ValueStr()
	note := fsion.Get("note").ValueStr()
	dataSource := fsion.Get("data_source").ValueInt()
	attribute := fsion.Get("attribute").ValueInt()
	headersJson := fsion.Get("headers").ValueStr()
	supportsJson := fsion.Get("supports").ValueStr()
	linkJson := fsion.Get("link").ValueStr()
	var links []int
	if linkJson != "" {
		err := json.Unmarshal([]byte(linkJson), &links)
		if err != nil {
			ErrorStd(c, err.Error())
			return
		}
	}

	var headers []service.Header
	err := json.Unmarshal([]byte(headersJson), &headers)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	var supports []int
	err = json.Unmarshal([]byte(supportsJson), &supports)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	invoke := service.NewReportService()
	err = invoke.ModifyReportTables(reportID, sql, sqlCount, remark, note, headers, supports, links, dataSource, attribute)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}

	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/report/valid 验证SQL合法性
* @apiName api.SqlValid
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} sql
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
}
*/
func SqlValid(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	auto := fsion.Get("auto").ValueDefaultInt32(1)
	sql := fsion.Get("sql").ValueStr()
	dataSource := fsion.Get("data_source").ValueInt()
	//fmt.Println(sql)
	err := service.Valid(sql)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	//defines, _ := service.NewReportService().GetReportDefinesByUuid(uuid)
	var ret []map[string]interface{}

	da := new(dao.DataAdapter)
	da.Adapter(dataSource)
	if auto != service.Smart {
		args := make(map[string]interface{})
		args["px"] = ""
		args["time1"] = ""
		args["p"] = 0
		args["limit"] = 10
		//
		//tpl, _ := pongo2.FromString(sql)
		//sql, _ = tpl.Execute(pongo2.Context(args))
		//ret, err = dao.NewInfiniteStones().PtrWarWithParametes(service.FormatSql(sql, args))
		ret, err = da.Do(service.FormatSql(sql, args))
	} else {
		ret, err = da.Do(sql+"  LIMIT 10", nil)
	}

	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStd(c, ret)
}

/**
* @api {Get} /admin/v1/report/details 获取报表详情(包括检索条件、表头等)
* @apiName api.Details
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} uuid
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
* {
    "code": 200,
    "data": {
        "defines": {  ---主体结构
            "report_id": 4,
            "uuid": "b7c2424b699011eab04494c69172a40b",
            "report_name": "报表20200319",
			"report_type": 7, ---1.明细  2.汇总  4，图表  计算方法按位与或, 例 (7&2) >0 ,则汇总项开启，(7&4) >0 ,则有图表
            "sql": "select channel_name, \r\n       category_name,  \r\n\t\t\t sum(consumption) consumption, \r\n\t\t\t sum(cash) cash, \r\n\t\t\t sum(exhibition) exhibition, \r\n\t\t\t sum(click) click, \r\n\t\t\t sum(pv) pv, \r\n\t\t\t sum(uv) uv, \r\n\t\t\t sum(register) register,  \r\n\t\t\t sum(cash)/sum(register) AS cr\r\n  from  dxz_gio_studying_abroad_track\r\n  #WHERE\r\n  group by  channel_name,   \r\n           category_name",
            "in_time": "2020-03-19T00:00:00+08:00"
        },
        "headers": [  ---表头
            {
                "eng_name": "channel_name",
                "chs_name": "渠道名称"
            },
            {
                "eng_name": "category_name",
                "chs_name": "分类名称"
            },
            {
                "eng_name": "consumption",
                "chs_name": "消费"
            },
            {
                "eng_name": "cash",
                "chs_name": "现金"
            },
            {
                "eng_name": "exhibition",
                "chs_name": "展示"
            },
            {
                "eng_name": "click",
                "chs_name": "点击"
            },
            {
                "eng_name": "pv",
                "chs_name": "pv量"
            },
            {
                "eng_name": "uv",
                "chs_name": "uv量"
            },
            {
                "eng_name": "register",
                "chs_name": "注册量"
            },
            {
                "eng_name": "cr",
                "chs_name": "注册成本"
            }
        ],
        "support": [   ---条件
            {
                "cond_id": 2,
                "name": "report_date",
                "type": "select",
                "val": "[{\"k\":\"昨天\",\"v\":\"yesterday\"},{\"k\":\"今天\",\"v\":\"day\"},{\"k\":\"本周\",\"v\":\"week\"},{\"k\":\"上周\",\"v\":\"lastWeek\"},{\"k\":\"本月\",\"v\":\"month\"},{\"k\":\"上月\",\"v\":\"lastMonth\"},{\"k\":\"本季度\",\"v\":\"quarter\"},{\"k\":\"上季度\",\"v\":\"lastQuarter\"},{\"k\":\"本年\",\"v\":\"year\"},{\"k\":\"去年\",\"v\":\"lastYear\"},{\"k\":\"自定义\",\"v\":\"custom\"}]",
                "placeholder": "时间范围"
            }
        ]
    }
}
*/
func Details(c *gin.Context) {
	reportID := c.Query("uuid")
	ret := make(map[string]interface{})
	invoke := service.NewReportService()
	ret["defines"], _ = invoke.GetReportDefinesByUuid(reportID)
	ret["headers"], _ = invoke.GetReportHeaderByUuid(reportID)
	ret["support"], _ = invoke.GetReportConditionsByUuid(reportID)
	ret["link"] = new(dao.PlatformReport).GetReportLinkById(invoke.GetReportIdByUuid(reportID))
	SuccessStd(c, ret)
}

/**
* @api {POST} /admin/v1/report/conditions 获取某个表所支持的检索条件
* @apiName api.Conditions
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} sql sql
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	"data":{}
}
*/
func Conditions(c *gin.Context) {
	//body, _ := ioutil.ReadAll(c.Request.Body)
	//fsion := gofasion.NewFasion(string(body))
	//sql := fsion.Get("sql").ValueStr()

	//err := service.Valid(sql)
	//if err != nil {
	//
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"data": nil,
	//		"msg":  err.Error(),
	//	})
	//	return
	//}

	//columns, err := dao.NewInfiniteStones().Valid(sql + "  LIMIT 1")
	//
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": 200,
	//		"data": nil,
	//		"msg":  err.Error(),
	//	})
	//	return
	//}

	invoke := service.NewReportService()
	conditions := invoke.GetReportConditions()
	//ret := make([]*dao.ReportConditions, 0, 5)
	//reMark := make(map[string]bool)
	//fmt.Println(columns)
	//for _, item := range columns {
	//	for cl, _ := range item {
	//		for _, cond := range conditions {
	//			if cond.CondName == cl {
	//				if _, ok := reMark[cl]; !ok {
	//					ret = append(ret, cond)
	//					reMark[cl] = true
	//				}
	//			}
	//		}
	//	}
	//	continue
	//}
	SuccessStd(c, conditions)
}

/**
* @api {Post} /admin/v1/report/remove 删除报表
* @apiName api.Delete
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} uuid
* @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func Delete(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	_ = dao.NewPlatformReport().RemoveReportByID(uuid)
	service.NewReportService().CC()
	service.NewDefinesService().RemoveCache(uuid)
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/report/rename 修改报表名称
* @apiName api.Rename
* @apiGroup admin 管理
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid
  * @apiParam {String} name 报表名称
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func Rename(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	name := fsion.Get("name").ValueStr()
	err := service.NewReportService().ModifyReportName(name, uuid)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/report/summary 修改报表汇总
  * @apiName api.SqlSummary
  * @apiGroup admin 管理
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid
  * @apiParam {String} sql
  * @apiParam {String} [tip]
  * @apiParam {Int} enable  0.禁用 1.启用
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func SqlSummary(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	sql := fsion.Get("sql").ValueStr()
	tip := fsion.Get("tip").ValueStr()
	enable := fsion.Get("enable").ValueInt()
	err := service.NewReportService().ModifyReportSqlSummary(uuid, sql, tip, enable)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/report/enable 修改报表可见性
  * @apiName api.Enable
  * @apiGroup admin 管理
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid
  * @apiParam {Int} enable  1|2|4 ??
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func Enable(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	enable := fsion.Get("enable").ValueInt()
	err := service.NewReportService().ModifyReportType(uuid, enable)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/report/auto  是否启用sql智能模式
  * @apiName api.Auto
  * @apiGroup admin 管理
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid
  * @apiParam {Int} auto  1.是 0.否
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func Auto(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	auto := fsion.Get("auto").ValueInt()
	err := service.NewReportService().ModifyReportAuto(uuid, auto)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

func AllForm(c *gin.Context) {
	//service.NewReportService().get

	items, _ := service.NewDefinesService().Get("", 1, 1000)
	SuccessStd(c, items)
}
