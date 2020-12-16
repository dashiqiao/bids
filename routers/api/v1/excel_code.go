package v1

import (
	"encoding/json"
	"io/ioutil"
	"platform_report/service"

	"github.com/gin-gonic/gin"
)

/**
* @api {POST} /api/v1/tool/tec 新增excel相关的字段字典
* @apiName AddExcelCode
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} name 字典名称
* @apiParam {String} type 类型（select,array）
* @apiParam {String} value 字典值
* @apiParam {String} tips 提示
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/

func AddExcelCode(c *gin.Context) {

	tecSvc := new(service.TableExcelCodeService)
	dataByte, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(dataByte, tecSvc)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	err = tecSvc.AddExcelCode()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return
}

/**
* @api {PUT} /api/v1/tool/tec 修改excel相关的字段字典
* @apiName ModifyExcelCode
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id id
* @apiParam {String} name 字典名称
* @apiParam {String} type 类型（select,array）
* @apiParam {String} value 字典值
* @apiParam {String} tips 提示
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/
func ModifyExcelCode(c *gin.Context) {
	temSvc := new(service.TableExcelCodeService)
	err := c.ShouldBind(temSvc)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}
	err = temSvc.ModifyExcelCode()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return
}

/**
* @api {DELETE} /api/v1/tool/tec 删除excel相关的字段字典
* @apiName DelExcelCode
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id id
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/
func DelExcelCode(c *gin.Context) {
	id := c.Query("id")
	tecSvc := new(service.TableExcelCodeService)
	err := tecSvc.DeleteExcelCode(id)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return
}

/**
* @api {GET} /api/v1/tool/tec 获取excel相关的字段字典
* @apiName GetExcelCode
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id id
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": {
        "id": 2,
        "name": "性别",
        "type": "array",
        "value": "[{\"k\":\"男\",\"v\":\"2\"}]",
        "tips": "这是性别1",
        "conduct": 2
    },
    "msg": "ok!"
}
*/

func GetExcelCode(c *gin.Context) {
	id := c.Query("id")
	tecSvc := new(service.TableExcelCodeService)
	tec, err := tecSvc.GetExcelCodeInfo(id)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, tec, "ok!")
	return
}

/**
* @api {GET} /api/v1/tool/tec/list 获取excel相关的字段字典列表
* @apiName GetExcelCodeList
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} limit 分页数量
* @apiParam {int} page 页码
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": {
        "count": 2,
        "list": [
            {
                "id": 2,
                "name": "性别",
                "type": "array",
                "value": "[{\"k\":\"男\",\"v\":\"2\"}]",
                "tips": "这是性别1",
                "conduct": 2
            },
            {
                "id": 3,
                "name": "性别",
                "type": "array",
                "value": "[{\"k\":\"男\",\"v\":\"1\"}]",
                "tips": "这是性别",
                "conduct": 2
            }
        ]
    },
    "msg": "ok!"
}
*/
func GetExcelCodeList(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	name := c.Query("name")
	codeType := c.Query("type")

	tecSvc := new(service.TableExcelCodeService)
	tecSvc.Name = name
	tecSvc.Type = codeType

	list, err := tecSvc.GetExcelCodeList(page, limit)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "ok!")
	return
}

/**
* @api {GET} /api/v1/tool/tec/reload 刷新excel相关的字典缓存
* @apiName ReloadExcelCodeList
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/

func ReloadExcelCodeList(c *gin.Context) {

	tecSvc := new(service.TableExcelCodeService)

	err := tecSvc.ReloadExcelCodeList()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return
}
