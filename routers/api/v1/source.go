package v1

import (
	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"platform_report/dao"
)

/**
* @api {Get} /admin/v1/source/index 数据源列表
* @apiName api.SourceDefines
* @apiGroup source
* @apiHeader {String} Authorization 用户授权token
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
	"data":
		[
			{
			"ds_id:":"xxx" --数据源ID
			"ds_name": "XXX" --数据源名称
			}
		]
	}
*/
func SourceDefines(c *gin.Context) {
	items, _ := new(dao.ReportDataSource).GetMulti()
	c.JSON(200, gin.H{
		"code": 200,
		"data": items,
	})
}

/**
* @api {Post} /admin/v1/source/create 添加数据源
* @apiName api.AddSource
* @apiGroup source
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} name 数据源名称
* @apiParam {String} catalog 配置文件名称
* @apiParam {String} schema 数据库名称
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
        "code": 200
	 }
*/
func AddSource(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	name := fsion.Get("name").ValueStr()
	catalog := fsion.Get("catalog").ValueStr()
	schema := fsion.Get("schema").ValueStr()
	invoke := &dao.ReportDataSource{
		DsName:  name,
		Catalog: catalog,
		Schema:  schema,
	}
	item, err := invoke.Create()
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStd(c, item)
}

/**
* @api {Post} /admin/v1/source/update 修改数据源
* @apiName api.UpdSource
* @apiGroup source
* @apiHeader {String} Authorization 用户授权token
* @apiParam {Int} id 数据源ID
* @apiParam {String} name 数据源名称
* @apiParam {String} catalog 配置文件名称
* @apiParam {String} schema 数据库名称
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
        "code": 200
	 }
*/
func UpdSource(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("id").ValueInt()
	name := fsion.Get("name").ValueStr()
	catalog := fsion.Get("catalog").ValueStr()
	schema := fsion.Get("schema").ValueStr()
	maps := make(map[string]interface{})
	maps["ds_name"] = name
	maps["catalog"] = catalog
	maps["schema"] = schema
	err := new(dao.ReportDataSource).Modify(id, maps)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}


/**
* @api {Post} /admin/v1/source/delete 删除数据源
* @apiName api.RemoveSource
* @apiGroup source
* @apiHeader {String} Authorization 用户授权token
* @apiParam {Int} id 数据源ID
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
        "code": 200
	 }
*/
func RemoveSource(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("id").ValueInt()
	err := new(dao.ReportDataSource).Delete(id)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}
