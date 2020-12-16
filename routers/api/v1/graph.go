package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"platform_report/dao"
	"platform_report/lib"
	"platform_report/service"
	"strconv"
	"strings"

	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
)

/**
  * @api {Post} /admin/v1/graph/save 添加图表
  * @apiName api.AddGraph
  * @apiGroup graph
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} name  图表名称
  * @apiParam {String} sql
  * @apiParam {Int} auto  1.智能模式 0.自定义模式
  * @apiParam {Int} make 0.使用字段名 1.使用字段值
  * @apiParam {String} sql sql
  * @apiParam {Int} type  1.折线图 2.饼图 3.柱状 4.漏斗
  * @apiParam {String} x  x轴 列名,如 report_date
  * @apiParam {String} y_eng_name  y轴  英文列名 ["a1","a3","a2"]
  * @apiParam {String} y_chs_name  y轴  中文列表 ["a1","a3","a2"]
  * @apiParam {String} supports  [1,2,3,4]  --条件ID
  * @apiParam {Int} preview  1.预览，不入库
  * @apiParam {String} [tip] 提示
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }

*/
func AddGraph(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	name := fsion.Get("name").ValueStr()
	sql := fsion.Get("sql").ValueStr()
	gtype := fsion.Get("type").ValueInt()
	x := fsion.Get("x").ValueStr()
	yEngName := fsion.Get("y_eng_name").ValueStr()
	yChsName := fsion.Get("y_chs_name").ValueStr()
	preview := fsion.Get("preview").ValueInt()
	tip := fsion.Get("tip").ValueStr()
	auto := fsion.Get("auto").ValueInt()
	dataSource := fsion.Get("data_source").ValueInt()
	supportsJson := fsion.Get("supports").ValueStr()
	var supports []int
	_ = json.Unmarshal([]byte(supportsJson), &supports)
	mk := fsion.Get("make").ValueInt()
	if preview == 1 {
		_, boolVal := c.GetPostFormMap("status")
		if boolVal {
		}
		m := c.Request.PostForm
		params := make(map[string]interface{})
		for k, v := range m {
			params[k] = v[0]
		}

		t1, t2 := lib.GetTimeByType("yesterday")
		params["time1"], params["time2"] = lib.TimeFormat(t1), lib.TimeFormat(t2)

		combine := service.NewCombineService(params, nil)
		graph := combine.GetReportGraphPreviewIndex(name, sql, mk, gtype, auto, x, yEngName, yChsName, dataSource)

		c.JSON(http.StatusOK, gin.H{
			"code":  http.StatusOK,
			"graph": graph,
		})
	} else {
		err := service.NewGraphService().AddGraph(name, sql, gtype, mk, auto, x, yEngName, yChsName, tip, supports, dataSource)
		if err != nil {
			ErrorStd(c, err.Error())
			return
		}
		SuccessStds(c)
	}
}

/**
  * @api {Post} /admin/v1/graph/modify 修改图表
  * @apiName api.ModifyGraph
  * @apiGroup graph
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} graph  图表ID
  * @apiParam {String} name  图表名称
  * @apiParam {String} sql
  * @apiParam {Int} auto  1.智能模式 0.自定义模式
  * @apiParam {Int} make 0.使用字段名 1.使用字段值
  * @apiParam {String} sql
  * @apiParam {Int} type  1.折线图 2.饼图 3.柱状 4.漏斗
  * @apiParam {String} x  x轴 列名,如 report_date
  * @apiParam {String} y_eng_name  y轴  英文列名 ["a1","a3","a2"]
  * @apiParam {String} y_chs_name  y轴  中文列表 ["a1","a3","a2"]
  * @apiParam {String} supports  [1,2,3,4]  --条件ID
  * @apiParam {String} [tip] 提示
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func ModifyGraph(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	graph := fsion.Get("graph").ValueStr()
	name := fsion.Get("name").ValueStr()
	sql := fsion.Get("sql").ValueStr()
	gtype := fsion.Get("type").ValueInt()
	auto := fsion.Get("auto").ValueInt()
	x := fsion.Get("x").ValueStr()
	yEngName := fsion.Get("y_eng_name").ValueStr()
	yChsName := fsion.Get("y_chs_name").ValueStr()
	tip := fsion.Get("tip").ValueStr()
	supportsJson := fsion.Get("supports").ValueStr()
	mk := fsion.Get("make").ValueInt()
	dataSource := fsion.Get("data_source").ValueInt()
	var supports []int
	_ = json.Unmarshal([]byte(supportsJson), &supports)

	err := service.NewGraphService().ModifyGraph(graph, name, sql, gtype, mk, auto, x, yEngName, yChsName, tip, supports, dataSource)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/graph/delete 删除图表
  * @apiName api.RemoveGraph
  * @apiGroup graph
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} graph  图表ID
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func RemoveGraph(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	graph := fsion.Get("graph").ValueStr()
	err := service.NewGraphService().RemoveGraph(graph)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Get} /admin/v1/graph/index 图表列表
  * @apiName api.GetGraph
  * @apiGroup graph
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid  报表ID
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok",
  		"data": [
        {
            "graph_id": "69162dbc1e6e4495a3de7156e89d1df3", --图表ID
            "graph_name": "test01", --图表名称
            "source": 1, --数据源 1.明细 2.汇总 3.自定义SQL
            "sql": "",
            "graph_type": 1, --1.折线  2.饼
            "x": "report_date",
            "y_chs_name": "",
            "y_eng_name": "",
            "in_time": "2020-04-16T14:30:38+08:00"
        }
    ]
  }
*/
func GetGraph(c *gin.Context) {
	uuid := c.Query("uuid")
	graphs := service.NewGraphService().GetGraphList(uuid)
	SuccessStd(c, graphs)
}

/**
  * @api {Get} /admin/v1/graph/all 全部图表
  * @apiName api.GetAllGraph
  * @apiGroup graph
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} [uuid]  报表ID   查看与此报表匹配的图表
  * @apiParam {String} [title] 模糊搜索标题
  * @apiParam {Int} [type] 1.折线图 2.饼图 3.柱状 4.漏斗
  * @apiParam {Int} p  当前页
  * @apiParam {Int} limit 分页大小
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok",
      "count":100,
  		"data": [
        {
            "graph_id": "69162dbc1e6e4495a3de7156e89d1df3", --图表ID
            "graph_name": "test01", --图表名称
            "source": 1, --数据源 1.明细 2.汇总 3.自定义SQL
            "sql": "",
            "graph_type": 1, --1.折线  2.饼
            "x": "report_date",
            "y_chs_name": "",
            "y_eng_name": "",
            "in_time": "2020-04-16T14:30:38+08:00"
        }
    ]
  }
*/
func GetAllGraph(c *gin.Context) {
	title := c.Query("title")
	uuid := c.Query("uuid")
	stype, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	var ret []*dao.ReportGraph
	if uuid != "" {
		ret = service.NewGraphService().GetAllGraphListByCondition(uuid)
	} else {
		ret = service.NewGraphService().GetAllGraphList()
	}

	if title != "" {
		newRet := make([]*dao.ReportGraph, 0)
		for _, val := range ret {
			if strings.Contains(val.GraphName, title) {
				newRet = append(newRet, val)
			}
		}
		ret = newRet
	}
	if stype > 0 {
		newRet := make([]*dao.ReportGraph, 0)
		for _, val := range ret {
			if val.GraphType == stype {
				newRet = append(newRet, val)
			}
		}
		ret = newRet
	}

	count := len(ret)
	m, n := (p-1)*limit, p*limit
	var maps []*dao.ReportGraph
	if n > count {
		if m > count {
		} else {
			maps = ret[m:]
		}
	} else {
		maps = ret[m:n]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  maps,
		"count": count,
	})
}

/**
  * @api {Get} /admin/v1/graph/detail 图表详情
  * @apiName api.GetGraphByID
  * @apiGroup graph
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} graph 图表ID
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func GetGraphByID(c *gin.Context) {
	//uuid := c.Query("u/uid")
	graph := c.Query("graph")
	graphs := service.NewGraphService().GetGraphByGraphId(graph)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    graphs,
		"support": service.NewGraphService().GetGraphConditions(graph),
	})
}

/**
  * @api {Post} /admin/v1/graph/link 关联图表
  * @apiName api.LinkGraph
  * @apiGroup graph
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 报表ID
  * @apiParam {String} graph 图表ID
  * @apiParam {Int} state 1.关联  0.移除关联
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func LinkGraph(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	graph := fsion.Get("graph").ValueStr()
	state := fsion.Get("state").ValueInt()
	if state == 0 {
		err := service.NewGraphService().UnLinkGraph(uuid, graph)
		if err != nil {
			ErrorStd(c, err.Error())
			return
		}
	} else {
		err := service.NewGraphService().LinkGraph(uuid, graph)
		if err != nil {
			ErrorStd(c, err.Error())
			return
		}
	}
	new(service.AgentService).Clean(uuid)
	SuccessStds(c)
}
