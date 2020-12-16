package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"platform_report/dao"
	"platform_report/service"
	"strconv"
	"strings"

	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
)

/**
* @api {Get} /api/v1/report/fields 字段列表
* @apiName api.GetFields
* @apiGroup fields 字段
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} uuid 报表ID
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*   {
    "code": 200,
    "data": [
        {
            "action_id": 27,
            "field_id": 6,
            "types": "sem_channel",
            "field": "id",  --数据库字段名称
            "name": "主键", -- 显示字段名称
            "form_type": "text", -- 字段类型 text / select / datetime
            "default_value": "", -- 默认值
            "max_length": 0, --最大长度
            "is_primary": 1,
            "is_unique": 1,
            "is_null": 1, -- 是否必填 1.是
            "is_show": 0, -- 是否显示 1.是 0.不显示
            "input_tips": "", -- 输入提示
			"parent_id": 0, -- 父级条件ID
            "setting": 0 , -- 条件ID
			"value":"", --- 下拉菜单的值
            "order_id": 0, -- 排序ID
            "create_time": "2020-04-28T09:49:42+08:00"
        }
    ],
    "values": {
        "channel_name": "55m+5bqm",
        "created_at": "MjAxOC0wOS0yNSAxNjozODozMg==",
        "deleted_at": null,
        "id": 1,
        "status": 0,
        "updated_at": null
    }
}
*/
func GetFields(c *gin.Context) {
	uuid := c.Query("uuid")
	invoke := service.NewFieldService()
	reportId := service.NewReportService().GetReportIdByUuid(uuid)
	primaryKeys := invoke.GetPrimaryKey(reportId)
	resMap := make(map[string]interface{})
	for _, primaryKey := range primaryKeys {
		resMap[primaryKey] = c.Query(primaryKey)
	}
	fields := invoke.GetFields(reportId)
	values := invoke.GetValues(reportId, resMap)
	c.JSON(http.StatusOK, gin.H{
		"code":   200,
		"flelds": fields,
		"values": values,
	})
}

/**
* @api {Post} /api/v1/report/actions 动作
* @apiName api.Actions
* @apiGroup fields 字段
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} uuid 报表ID
* @apiParam {String} action Insert/Update/Delete
* @apiParam {String} .. 其他参数
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*   {
    "code": 200,
    "msg": ""
}
*/
func Actions(c *gin.Context) {
	uuid := c.PostForm("uuid")
	action := c.PostForm("action")
	_, boolVal := c.GetPostFormMap("status")
	if boolVal {
	}
	m := c.Request.PostForm
	params := make(map[string]interface{})
	for k, v := range m {
		params[k] = v[0]
	}
	invoke := service.NewFieldService()
	reportId := service.NewReportService().GetReportIdByUuid(uuid)

	//fields := invoke.GetFields(reportId)
	primaryKeys := invoke.GetPrimaryKey(reportId)
	resMap := make(map[string]interface{})
	for _, primaryKey := range primaryKeys {
		resMap[primaryKey] = c.PostForm(primaryKey)
	}
	//before := invoke.GetValues(reportId, resMap)
	ret := service.NewReportAction().Do(uuid, action, params)
	if ret.Code != 0 {
		ErrorStd(c, ret.Msg)
		return
	}
	service.NewReportService().CC()
	//_ = ants.Submit(func() {
	//	//操作日志
	//	service.NewLogService().AddReportActionLog(uuid, fields[0].Types, action, before, params, UserInfo(c).ID, UserInfo(c).Realname)
	//})
	new(service.AgentService).Clean(uuid)
	SuccessStds(c)
}

//--------------------------------------

/**
* @api {Get} /admin/v1/actions/index 动作列表
* @apiName api.ActionDefines
* @apiGroup actions
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
	}
*/
func ActionDefines(c *gin.Context) {
	title := c.Query("title")
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	items, _ := service.NewFieldService().GetActions()
	if title != "" {
		newRet := make([]*dao.ReportFieldDefines, 0)
		for _, val := range items {
			if strings.Contains(val.Name, title) {
				newRet = append(newRet, val)
			}
		}
		items = newRet
	}
	count := len(items)
	m, n := (p-1)*limit, p*limit
	var maps []*dao.ReportFieldDefines
	if n > count {
		if m > count {
		} else {
			maps = items[m:]
		}
	} else {
		maps = items[m:n]
	}

	c.JSON(200, gin.H{
		"code":  200,
		"data":  maps,
		"count": count,
	})
}

/**
* @api {Get} /admin/v1/actions/field 动作详情
* @apiName api.ActionFields
* @apiGroup actions
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} action_id 动作ID
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "msg": "ok"
	"data":
	}
*/
func ActionFields(c *gin.Context) {
	actionId, _ := strconv.Atoi(c.DefaultQuery("action_id", "0"))
	items, _ := dao.NewReportFieldInvoke().GetFieldsByActionId(actionId)

	c.JSON(200, gin.H{
		"code":      200,
		"data":      items,
		"condition": service.NewReportService().GetReportSelectConditions(),
	})
}

/**
  * @api {Post} /admin/v1/actions/save 添加或修改动作
  * @apiName api.Action
  * @apiGroup actions
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 动作ID
  * @apiParam {String} title 动作标题
  * @apiParam {String} fields json==>格式如下 ：
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
		  "code": 200,
		  "msg": "ok"
  		}
        -----------json-------------
        [{
            "types": "sem_channel",
            "field": "id",  --数据库字段名称
            "name": "...", -- 显示字段名称
            "form_type": "text", -- 字段类型 text / select / datetime
            "default_value": "", -- 默认值
            "max_length": 0, --最大长度
            "is_primary": 1, --是否主键
            "is_unique": 1,  --是否唯一
            "is_null": 1, -- 是否必填 1.是
            "is_show": 0, -- 是否显示 1.是 0.不显示
            "input_tips": "", -- 输入提示
            "setting": 0, -- 条件ID
            "order_id": 0, -- 排序ID
  		},{...}}
*/
func SaveAction(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	title := fsion.Get("title").ValueStr()
	fieldsJson := fsion.Get("fields").ValueStr()
	var fields []dao.ReportField
	err := json.Unmarshal([]byte(fieldsJson), &fields)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	invoke := service.NewFieldService()
	err = invoke.SaveFields(uuid, title, fields)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/actions/enable 修改报表动作开关
  * @apiName api.ActionEnable
  * @apiGroup actions
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 报表ID
  * @apiParam {Int} enable  1|2|4|8 ==> 对应 增、删、改、查
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
		  "code": 200,
		  "msg": "ok"
  		}
*/
func ActionEnable(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	enable := fsion.Get("enable").ValueInt()
	err := service.NewReportService().ModifyReportAction(uuid, enable)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/actions/bind 报表关联
  * @apiName api.ActionBind
  * @apiGroup actions
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 报表ID
  * @apiParam {Int} action_id  动作id
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
		  "code": 200,
		  "msg": "ok"
  		}
*/
func ActionBind(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	action := fsion.Get("action_id").ValueInt()
	reportId := service.NewReportService().GetReportIdByUuid(uuid)
	err := dao.NewReportFieldInvoke().Bind(reportId, action)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/actions/unbind 报表取消关联
  * @apiName api.ActionUnBind
  * @apiGroup actions
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 报表ID
  * @apiParam {Int} action_id  动作id
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
		  "code": 200,
		  "msg": "ok"
  		}
*/
func ActionUnBind(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	action := fsion.Get("action_id").ValueInt()
	reportId := service.NewReportService().GetReportIdByUuid(uuid)
	err := dao.NewReportFieldInvoke().UnBind(reportId, action)
	if err != nil {
		ErrorStd(c, err.Error())
		return
	}
	SuccessStds(c)
}

/**
  * @api {Post} /admin/v1/report/actions 报表关联的动作
  * @apiName api.ReportActions
  * @apiGroup actions
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid 报表ID
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
		  "code": 200,
		  "msg": "ok"
  		}
*/
func ReportActions(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	reportId := service.NewReportService().GetReportIdByUuid(uuid)
	//fmt.Println("reportId===", reportId)
	item, _ := dao.NewReportFieldInvoke().GetActionsByReportId(reportId)
	//item := service.NewFieldService().GetActionsById(uuid)
	SuccessStd(c, item)
}

/**
  * @api {Post} /admin/v1/actions/delete 删除动作
  * @apiName api.DeleteActions
  * @apiGroup actions
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
		  "code": 200,
		  "msg": "ok"
  		}
*/
func DeleteActions(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	uuid := fsion.Get("uuid").ValueStr()
	service.NewFieldService().DeleteActions(uuid)
	SuccessStds(c)
}

func Test(c *gin.Context) {

}
