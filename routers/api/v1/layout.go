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
* @api {Get} /admin/v1/layout/index 布局列表
* @apiName api.LayoutDefines
* @apiGroup layout 布局
* @apiHeader {String} Authorization 用户授权token
* @apiParam {Int} p  当前页
* @apiParam {Int} limit 分页大小
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok",
	"data":
		[
			{
			"layout_id": 1, --ID
			"title": "这是一个新的布局" --标题
			"info": "xxxxx" --布局json
			}
		],
	"count":10 --总行数
	}
*/
func LayoutDefines(c *gin.Context) {
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	invoke := service.NewLayoutService()
	ret, err := invoke.Get()
	count := len(ret)
	m, n := (p-1)*limit, p*limit
	var maps []*dao.ReportLayout
	if n > count {
		if m > count {
		} else {
			maps = ret[m:]
		}
	} else {
		maps = ret[m:n]
	}
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	c.JSON(200, gin.H{
		"code":  200,
		"data":  maps,
		"count": count,
	})
}

/**
* @api {Post} /admin/v1/layout/modify 修改(标题/布局)
* @apiName api.ModifyLayout
* @apiGroup layout 布局
* @apiHeader {String} Authorization 用户授权token
* @apiParam {Int} id  布局ID
* @apiParam {String} [title]  标题
* @apiParam {String} [info]   布局信息
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	}
*/
func ModifyLayout(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("id").ValueInt()
	title := fsion.Get("title").ValueStr()
	info := fsion.Get("info").ValueStr()
	invoke := service.NewLayoutService()
	if title != "" {
		_ = invoke.ModifyTitle(id, title)
	}
	if info != "" {
		_ = invoke.ModifyInfo(id, info)
	}
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/layout/save 添加布局
* @apiName api.SaveLayout
* @apiGroup layout 布局
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} title  标题
* @apiParam {String} info   布局信息
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	}
*/
func SaveLayout(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	title := fsion.Get("title").ValueStr()
	info := fsion.Get("info").ValueStr()
	invoke := service.NewLayoutService()
	_, err := invoke.Save(title, info)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
* @api {Post} /admin/v1/layout/delete 删除布局
* @apiName api.RemoveLayout
* @apiGroup layout 布局
* @apiHeader {String} Authorization 用户授权token
* @apiParam {Int} id  布局ID
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	}
*/
func RemoveLayout(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("id").ValueInt()
	invoke := service.NewLayoutService()
	_ = invoke.Delete(id)
	SuccessStds(c)
}
