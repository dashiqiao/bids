package v1

import (
	"io/ioutil"
	"platform_report/dao"
	"platform_report/service"
	"strconv"

	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
)

/**
* @api {Get} /admin/v1/board/index 看板列表
* @apiName api.BoardDefines
* @apiGroup board 看板
* @apiHeader {String} Authorization 用户授权token
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
			"board_id": "885a3868675711eab04494c69172a40b", --看板ID
			"board_name": "885a3868675711eab04494c69172a40b", --看板名称
			"create_user_id" 6: , --
			"in_time": "2020-0366T11:18:36+08:00"
			}
		]
	}
*/
func BoardDefines(c *gin.Context) {
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	ret := service.NewBoardService().GetBoardDefines()
	count := len(ret)
	m, n := (p-1)*limit, p*limit
	var maps []*dao.ReportBoard
	if n > count {
		if m > count {

		} else {
			maps = ret[m:]
		}
	} else {
		maps = ret[m:n]
	}

	c.JSON(200, gin.H{
		"code":  200,
		"data":  maps,
		"count": count,
	})
}

/**
* @api {Post} /admin/v1/board/save 添加看板
* @apiName api.AddBoard
* @apiGroup board 看板
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} name 看板名称
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
}
*/
func AddBoard(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	name := fsion.Get("name").ValueStr()
	invoke := service.NewBoardService()
	err := invoke.AddReportBoard(name, UserInfo(c).ID)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/board/delete 删除看板
* @apiName api.RemoveBoard
* @apiGroup board 看板
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} id
* @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func RemoveBoard(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("id").ValueStr()
	_ = service.NewBoardService().RemoveBoardByID(uuid)
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/board/rename 修改看板名称
* @apiName api.RenameBoard
* @apiGroup board 看板
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} id
  * @apiParam {String} name 报表名称
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func RenameBoard(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("id").ValueStr()
	name := fsion.Get("name").ValueStr()
	err := service.NewBoardService().ModifyBoardName(uuid, name)
	if err != nil {
		ErrorStd(c, err.Error())
	}
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/board/layout 修改看板布局
* @apiName api.BoardLayout
* @apiGroup board 看板
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} id
  * @apiParam {String} name 报表名称
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }
*/
func BoardLayout(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	layout := fsion.Get("layout").ValueStr()
	layoutId, err := service.NewLayoutService().Save("", layout)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	_ = service.NewBoardService().ModifyLayout(uuid, layoutId)
	new(service.AgentService).Clean(uuid)
	SuccessStds(c)
}
