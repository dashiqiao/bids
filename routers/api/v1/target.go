package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"platform_report/service"
	"strconv"

	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
)

/**
* @api {Get} /admin/v1/target/defines 指标列表
* @apiName api.TargetDefines
* @apiGroup target 指标
* @apiHeader {String} Authorization 用户授权token
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"data":
		[
			{
			"target_id:":"xxx" --指标ID
			"icon": "XXX", --图标
			"eng_name": "XXXX", --英文名字
			"chs_name": "XXXX": --中文名字
			"report_id":"", --链接 报表ID
			"report_name":"" ,--报表名字
			"formula":"{X}#addition#{Y}#subtraction#{@@day}" 算式 "
			}
		]
	}
*/
func TargetDefines(c *gin.Context) {
	title := c.Query("title")
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	limit = 100
	items, _ := service.NewTargetService().GetTargetDefines(title, p, limit)
	invoke := service.NewReportService()
	for i := 0; i < len(items); i++ {
		if items[i].ReportId == "" {
			continue
		}
		defines, _ := invoke.GetReportDefinesByUuid(items[i].ReportId)
		if defines != nil {
			items[i].ReportName = defines.ReportName
		}
	}
	c.JSON(200, gin.H{
		"code":  200,
		"data":  items,
		"count": 10,
	})
	//SuccessStd(c, items)
}

/**
* @api {Post} /admin/v1/target/save 添加指标
* @apiName api.SaveTargetDefines
* @apiGroup target 指标
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} eng_name 英文名称
* @apiParam {String} chs_name 中文名称
* @apiParam {String} icon 图标
* @apiParam {String} [link] 链接
* @apiParam {String} [formula] 算式
* @apiParam {String} [tip] 提示
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"msg":""
     }
*/
func SaveTargetDefines(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	engName := fsion.Get("eng_name").ValueStr()
	chsName := fsion.Get("chs_name").ValueStr()
	icon := fsion.Get("icon").ValueStr()
	link := fsion.Get("link").ValueStr()
	formula := fsion.Get("formula").ValueStr()
	tip := fsion.Get("tip").ValueStr()
	err := service.NewTargetService().SaveTargetDefines(icon, engName, chsName, link, formula, tip)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/target/modify 修改指标
* @apiName api.ModifyTargetDefines
* @apiGroup target 指标
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} id 指标ID
* @apiParam {String} eng_name 英文名称
* @apiParam {String} chs_name 中文名称
* @apiParam {String} icon 图标
* @apiParam {String} [link] 链接
* @apiParam {String} [formula] 算式
* @apiParam {String} [tip] 提示
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"msg":""
     }

*/
func ModifyTargetDefines(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("id").ValueStr()
	engName := fsion.Get("eng_name").ValueStr()
	chsName := fsion.Get("chs_name").ValueStr()
	icon := fsion.Get("icon").ValueStr()
	link := fsion.Get("link").ValueStr()
	formula := fsion.Get("formula").ValueStr()
	tip := fsion.Get("tip").ValueStr()
	err := service.NewTargetService().ModifyTarget(id, icon, engName, chsName, link, formula, tip)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/target/delete 删除指标
* @apiName api.DeleteTargetDefines
* @apiGroup target 指标
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} id 指标ID
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"msg":""
     }

*/
func DeleteTargetDefines(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("id").ValueStr()
	err := service.NewTargetService().RemoveTargetByID(id)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
* @api {Get}/admin/v1/board/ls 读取看板
* @apiName api.GetBoardTarget
* @apiGroup target 指标
* @apiHeader {String} Authorization 用户授权token
* @apiHeader {String} board_id 看板ID
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"title":"xx",--指标卡片名称
	"target":[
			{
				"target_id":"", --指标ID
				"target_name":"", -- 指标名称
			}
		]
     }
	"graph":[]

*/
func GetBoard(c *gin.Context) {
	boardId := c.Query("board_id")
	fmt.Println(boardId)
	invoke := service.NewTargetService()
	title, layout, number, targets, graph := invoke.GetBoardTarget(boardId)
	target := make(map[string]interface{}, 0)
	//target["title"] = title
	target["data"] = targets
	target["number"] = number
	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"title":  title,
		"layout": layout,
		"target": target,
		"graph":  graph,
	})
}

/**
* @api {Post} /admin/v1/board/target/save （第一步）看板关联指标
* @apiName api.SaveBoardTarget
* @apiGroup target 指标
* @apiHeader {String} Authorization 用户授权token
* @apiHeader {String} board_id 看板ID
* @apiHeader {String} card_name 卡片名称
* @apiHeader {String} target_id 一组指标ID，格式：["111","2222"]
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"msg":""
     }

*/
func SaveBoardTarget(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	cardName := fsion.Get("card_name").ValueStr()
	boardId := fsion.Get("board_id").ValueStr()
	targetJson := fsion.Get("target_id").ValueStr()
	err := service.NewTargetService().SaveBoardTarget(cardName, boardId, targetJson)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	new(service.AgentService).Clean(boardId)
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/board/graph/save （第二步）看板关联图表/报表
* @apiName api.SaveBoardGraph
* @apiGroup target 指标
* @apiHeader {String} Authorization 用户授权token
* @apiHeader {String} json 格式如下
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"msg":""
     }
--------------请求参数----------------
 	{
    "board_id": "xxx", --看板ID
    "layout": 1, --布局
	"cards":[
			{
				"card_name":"xx",--卡片名称
				"graph_id":"xx",--图表ID
				"report_id":"xx",--报表ID
				"no":1--位置编号
			}
		]
     }
*/
func SaveBoardGraph(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	boardID := fsion.Get("board_id").ValueStr()
	layout := fsion.Get("layout").ValueInt()
	cardsJson := fsion.Get("cards").ValueStr()
	var cards []service.Cards
	err := json.Unmarshal([]byte(cardsJson), &cards)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	err = service.NewTargetService().SaveBoardGraph(boardID, cards)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	_ = service.NewBoardService().ModifyLayout(boardID, layout)
	new(service.AgentService).Clean(boardID)
	SuccessStds(c)
}
