package v1

import (
	"encoding/json"
	"fmt"
	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"platform_report/dao"
	"platform_report/service"
	"strconv"
)

/**
  * @api {Get} /admin/v1/analysis/defines 用户分组条件定义
  * @apiName api.CHDefines
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }

*/
func CHDefines(c *gin.Context) {
	type Quota struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
		Id    string      `json:"id"`
	}

	type Condition struct {
		Name      string      `json:"name"`
		Id        string      `json:"id"`
		ValueType string      `json:"value_type"`
		Value     interface{} `json:"value"`
		Range     []string    `json:"range"`
	}
	type Dict struct {
		K string      `json:"k"`
		V interface{} `json:"v"`
	}
	ranges := make([]Dict, 0, 9)
	ranges = append(ranges, Dict{K: "昨天", V: "yesterday"})
	ranges = append(ranges, Dict{K: "本周", V: "week"})
	ranges = append(ranges, Dict{K: "上周", V: "lastWeek"})
	ranges = append(ranges, Dict{K: "本月", V: "month"})
	ranges = append(ranges, Dict{K: "上月", V: "lastMonth"})
	ranges = append(ranges, Dict{K: "本季度", V: "quarter"})
	ranges = append(ranges, Dict{K: "上季度", V: "lastQuarter"})
	ranges = append(ranges, Dict{K: "本年", V: "year"})
	ranges = append(ranges, Dict{K: "去年", V: "lastYear"})

	dicts := make([]Dict, 0, 13)
	dicts = append(dicts, Dict{K: "文章", V: "1"})
	dicts = append(dicts, Dict{K: "课程&商品", V: "2"})
	dicts = append(dicts, Dict{K: "免费公开课", V: "3"})
	dicts = append(dicts, Dict{K: "问题", V: "4"})
	dicts = append(dicts, Dict{K: "相册", V: "5"})
	dicts = append(dicts, Dict{K: "APP小课程	", V: "6"})
	dicts = append(dicts, Dict{K: "文章评论", V: "7"})
	dicts = append(dicts, Dict{K: "话题", V: "8"})
	dicts = append(dicts, Dict{K: "院校主页", V: "9"})
	dicts = append(dicts, Dict{K: "个人主页", V: "10"})
	dicts = append(dicts, Dict{K: "问题回答", V: "11"})
	dicts = append(dicts, Dict{K: "课程&商品评论", V: "12"})
	dicts = append(dicts, Dict{K: "APP小课程评论", V: "13"})

	quota := make([]Quota, 0, 2)
	quota = append(quota, Quota{Name: "用户做过", Value: true, Id: "ud"})
	quota = append(quota, Quota{Name: "用户没做过", Value: false, Id: "und"})
	conditions := make([]Condition, 0, 9)
	conditions = append(conditions, Condition{Name: "活跃天数", Id: "hyts", ValueType: "天", Value: "", Range: []string{"ud"}})
	conditions = append(conditions, Condition{Name: "浏览", Id: "view", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "收藏", Id: "fav", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "分享", Id: "share", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "点赞", Id: "good", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "评论", Id: "comment", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "回答", Id: "reply", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "打赏", Id: "award", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "站内转发", Id: "distribute", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	conditions = append(conditions, Condition{Name: "订单", Id: "order", ValueType: "次", Value: "", Range: []string{"ud", "und"}})
	filter := make([]Condition, 0, 4)
	channel := make([]Dict, 0, 4)
	channel = append(channel, Dict{K: "语言培训", V: "1"})
	channel = append(channel, Dict{K: "出国留学", V: "2"})
	channel = append(channel, Dict{K: "海外移民", V: "4"})
	channel = append(channel, Dict{K: "院校直通车", V: "8"})

	yypx := make([]Dict, 0)
	cglx := make([]Dict, 0)
	cglx = append(cglx, Dict{K: "美国", V: "4087"})
	cglx = append(cglx, Dict{K: "英国", V: "4095"})
	cglx = append(cglx, Dict{K: "澳大利亚", V: "4103"})
	cglx = append(cglx, Dict{K: "加拿大", V: "4111"})
	cglx = append(cglx, Dict{K: "新西兰", V: "4119"})
	cglx = append(cglx, Dict{K: "日本", V: "4127"})
	cglx = append(cglx, Dict{K: "韩国", V: "4134"})
	cglx = append(cglx, Dict{K: "中国香港", V: "4141"})
	cglx = append(cglx, Dict{K: "俄罗斯", V: "4148"})
	cglx = append(cglx, Dict{K: "乌克兰", V: "4156"})
	cglx = append(cglx, Dict{K: "爱尔兰", V: "4163"})
	cglx = append(cglx, Dict{K: "意大利", V: "4169"})
	cglx = append(cglx, Dict{K: "西班牙", V: "4175"})
	cglx = append(cglx, Dict{K: "白俄罗斯", V: "4181"})
	cglx = append(cglx, Dict{K: "新加坡", V: "4188"})
	cglx = append(cglx, Dict{K: "马来西亚", V: "4195"})
	cglx = append(cglx, Dict{K: "丹麦", V: "4202"})
	cglx = append(cglx, Dict{K: "德国", V: "4208"})
	cglx = append(cglx, Dict{K: "法国", V: "4215"})
	cglx = append(cglx, Dict{K: "挪威", V: "4221"})
	cglx = append(cglx, Dict{K: "瑞典", V: "4227"})
	cglx = append(cglx, Dict{K: "芬兰", V: "4233"})
	cglx = append(cglx, Dict{K: "荷兰", V: "4239"})
	cglx = append(cglx, Dict{K: "瑞士", V: "4245"})
	cglx = append(cglx, Dict{K: "中国澳门", V: "4252"})
	cglx = append(cglx, Dict{K: "奥地利", V: "4258"})

	yypx = append(yypx, Dict{K: "托福", V: "4265"})
	yypx = append(yypx, Dict{K: "雅思", V: "4274"})
	yypx = append(yypx, Dict{K: "GRE", V: "4283"})
	yypx = append(yypx, Dict{K: "GMAT", V: "4289"})
	yypx = append(yypx, Dict{K: "PTE", V: "4296"})
	yypx = append(yypx, Dict{K: "LSAT", V: "4303"})
	yypx = append(yypx, Dict{K: "日语", V: "4309"})
	yypx = append(yypx, Dict{K: "韩语", V: "4316"})
	yypx = append(yypx, Dict{K: "法语", V: "4324"})
	yypx = append(yypx, Dict{K: "德语", V: "4331"})
	yypx = append(yypx, Dict{K: "俄语", V: "4340"})
	yypx = append(yypx, Dict{K: "意大利语", V: "4349"})
	yypx = append(yypx, Dict{K: "西班牙语", V: "4358"})
	yypx = append(yypx, Dict{K: "AP", V: "4367"})
	yypx = append(yypx, Dict{K: "ACT", V: "4375"})
	yypx = append(yypx, Dict{K: "SAT", V: "4383"})
	yypx = append(yypx, Dict{K: "SSAT", V: "4390"})
	yypx = append(yypx, Dict{K: "ALEVEL", V: "4397"})
	yypx = append(yypx, Dict{K: "AEAS", V: "4407"})
	yypx = append(yypx, Dict{K: "考研英语", V: "4412"})
	yypx = append(yypx, Dict{K: "英语四六级", V: "4417"})
	yypx = append(yypx, Dict{K: "多邻国", V: "78848"})

	hwym := make([]Dict, 0)
	yxztc := make([]Dict, 0)

	hwym = append(hwym, Dict{K: "美国", V: "4582"})
	hwym = append(hwym, Dict{K: "加拿大", V: "4592"})
	hwym = append(hwym, Dict{K: "澳大利亚", V: "4610"})
	hwym = append(hwym, Dict{K: "爱尔兰", V: "4621"})
	hwym = append(hwym, Dict{K: "希腊", V: "4625"})
	hwym = append(hwym, Dict{K: "英国", V: "4628"})
	hwym = append(hwym, Dict{K: "新西兰", V: "4634"})
	hwym = append(hwym, Dict{K: "西班牙", V: "4642"})
	hwym = append(hwym, Dict{K: "葡萄牙", V: "4645"})
	hwym = append(hwym, Dict{K: "马耳他", V: "4649"})
	hwym = append(hwym, Dict{K: "安提瓜", V: "4654"})
	hwym = append(hwym, Dict{K: "土耳其", V: "4658"})
	hwym = append(hwym, Dict{K: "格林纳达", V: "4661"})
	hwym = append(hwym, Dict{K: "塞浦路斯", V: "4665"})
	hwym = append(hwym, Dict{K: "圣基茨", V: "4669"})
	hwym = append(hwym, Dict{K: "多米尼克", V: "4673"})
	hwym = append(hwym, Dict{K: "中国香港", V: "4677"})
	hwym = append(hwym, Dict{K: "其他", V: "4783"})

	yxztc = append(yxztc, Dict{K: "美国", V: "4423"})
	yxztc = append(yxztc, Dict{K: "英国", V: "4431"})
	yxztc = append(yxztc, Dict{K: "澳大利亚", V: "4439"})
	yxztc = append(yxztc, Dict{K: "加拿大", V: "4447"})
	yxztc = append(yxztc, Dict{K: "新西兰", V: "4455"})
	yxztc = append(yxztc, Dict{K: "日本", V: "4463"})
	yxztc = append(yxztc, Dict{K: "韩国", V: "4470"})
	yxztc = append(yxztc, Dict{K: "中国香港", V: "4477"})

	yxztc = append(yxztc, Dict{K: "俄罗斯", V: "4484"})
	yxztc = append(yxztc, Dict{K: "乌克兰", V: "4492"})
	yxztc = append(yxztc, Dict{K: "爱尔兰", V: "4499"})
	yxztc = append(yxztc, Dict{K: "意大利", V: "4505"})
	yxztc = append(yxztc, Dict{K: "西班牙", V: "4511"})
	yxztc = append(yxztc, Dict{K: "白俄罗斯", V: "4517"})
	yxztc = append(yxztc, Dict{K: "新加坡", V: "4524"})
	yxztc = append(yxztc, Dict{K: "马来西亚", V: "4531"})
	yxztc = append(yxztc, Dict{K: "丹麦", V: "4538"})
	yxztc = append(yxztc, Dict{K: "德国", V: "4544"})
	yxztc = append(yxztc, Dict{K: "法国", V: "4551"})
	yxztc = append(yxztc, Dict{K: "瑞典", V: "4557"})
	yxztc = append(yxztc, Dict{K: "芬兰", V: "4563"})
	yxztc = append(yxztc, Dict{K: "荷兰", V: "4569"})
	yxztc = append(yxztc, Dict{K: "瑞士", V: "4575"})

	filter = append(filter, Condition{Name: "资源类型", Id: "type", Value: dicts, ValueType: "type"})
	filter = append(filter, Condition{Name: "资源标题", Id: "resource_title", Value: "", ValueType: "resource_title"})
	filter = append(filter, Condition{Name: "资源赛道", Id: "channel_id", Value: channel, ValueType: "channel_id"})
	filter = append(filter, Condition{Name: "语言培训", Id: "first_label_id", Value: yypx, ValueType: "yypx"})
	filter = append(filter, Condition{Name: "出国留学", Id: "first_label_id", Value: cglx, ValueType: "cglx"})
	filter = append(filter, Condition{Name: "海外移民", Id: "first_label_id", Value: hwym, ValueType: "hwym"})
	filter = append(filter, Condition{Name: "院校直通车", Id: "first_label_id", Value: yxztc, ValueType: "yxztc"})

	c.JSON(http.StatusOK, gin.H{
		"code":       200,
		"quota":      quota,
		"conditions": conditions,
		"filter":     filter,
		"range":      ranges,
	})
}

/**
  * @api {Post} /admin/v1/analysis/create 创建用户分组
  * @apiName api.CHAdd
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} name 分组名称
  * @apiParam {String} expr 表达式
  * @apiParam {String} conditions 条件
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }

*/
func CHAdd(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	name := fsion.Get("name").ValueStr()
	expr := fsion.Get("expr").ValueStr()
	conditions := fsion.Get("conditions").ValueStr()

	var analysis []dao.AnalysisConditions
	err := json.Unmarshal([]byte(conditions), &analysis)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	nas := service.NewAnalysisService()
	sql := nas.FormatSql(expr, analysis)
	members := nas.GetTotalMembers(service.FillCountHeader(sql))
	err = service.NewAnalysisService().Add(name, expr, conditions, int(members), UserInfo(c).ID)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
	})
}

/**
  * @api {Post} /admin/v1/analysis/update 修改用户分组
  * @apiName api.CHModify
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid
  * @apiParam {String} name 分组名称
  * @apiParam {String} expr 表达式
  * @apiParam {String} conditions 条件
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }

*/
func CHModify(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("uuid").ValueStr()
	name := fsion.Get("name").ValueStr()
	expr := fsion.Get("expr").ValueStr()
	conditions := fsion.Get("conditions").ValueStr()
	var analysis []dao.AnalysisConditions
	err := json.Unmarshal([]byte(conditions), &analysis)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	nas := service.NewAnalysisService()
	sql := nas.FormatSql(expr, analysis)
	members := nas.GetTotalMembers(service.FillCountHeader(sql))
	err = nas.Modify(id, name, expr, conditions, int(members))
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
	})
}

/**
  * @api {Post} /admin/v1/analysis/delete 删除用户分组
  * @apiName api.CHRemove
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} uuid
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "msg": "ok"
  }

*/
func CHRemove(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("uuid").ValueStr()
	fmt.Println("==========", id)
	err := service.NewAnalysisService().Remove(id)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
	})
}

/**
  * @api {Get} /admin/v1/analysis/index 用户分组列表
  * @apiName api.CHList
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {int} p  当前页
  * @apiParam {int} limit 分页大小
  * @apiParam {string} [uuid] 分组ID
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *  {
    "code": 200,
    "data": [
        {
            "uuid": "xxxx",
            "name": "测试别删",  -- 分组名称
            "expr": "A and B or C",  -- 表达式
            "conditions": "[{\"op\":\"=\",\"values\":[1,10],\"flag\":true,\"alias\":\"A\",\"dimFilters\":[{\"dim\":\"resource_title\",\"op\":\"not like\",\"values\":[\"www.growingio.com\"]}],\"id\":\"view\",\"range\":\"2020-01-01,2020-07-10\"},{\"flag\":false,\"alias\":\"B\",\"id\":\"view\",\"range\":\"2020-01-01,2020-07-10\"},{\"flag\":true,\"alias\":\"C\",\"id\":\"hyts\",\"range\":\"2020-01-01,2020-07-10\",\"op\":\"=\",\"values\":[1,10]}]",
            "members": 0 --人数
        }
    ],
    "total": 1
}

*/
func CHList(c *gin.Context) {
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	uuid := c.Query("uuid")
	if uuid != "" {
		item := service.NewAnalysisService().GetAnalysisByID(uuid)
		var analysis []dao.AnalysisConditions
		_ = json.Unmarshal([]byte(item.Conditions), &analysis)

		nas := service.NewAnalysisService()
		sql := nas.FormatSql(item.Expr, analysis)
		members := nas.GetTotalMembers(service.FillCountHeader(sql))
		_ = nas.ModifyAnalysisMembers(uuid, int(members))
		item.Members = int(members)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": item,
		})
	} else {
		items, total := service.NewAnalysisService().GetAnalysis(p, limit)
		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"data":  items,
			"total": total,
		})
	}
}

/**
  * @api {Post} /admin/v1/analysis/percent 占比
  * @apiName api.CHPercent
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {String} expr 表达式
  * @apiParam {String} conditions 条件
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
      "code": 200,
      "members": 100,
	  "percent": 98.10,
  }

*/
func CHPercent(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	expr := fsion.Get("expr").ValueStr()
	conditions := fsion.Get("conditions").ValueStr()
	var analysis []dao.AnalysisConditions
	err := json.Unmarshal([]byte(conditions), &analysis)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	nas := service.NewAnalysisService()
	sql := nas.FormatSql(expr, analysis)

	total := nas.GetTotalMembers("select count(distinct  user_id) as members from user_behavior_new")
	members := nas.GetTotalMembers(service.FillCountHeader(sql))

	fmt.Println(sql)
	percent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(members)/float64(total)), 64)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"members": members,
		"percent": percent * 100,
	})
}

/**
  * @api {Get} /admin/v1/analysis/ls 用户列表
  * @apiName api.CHUserList
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {int} p  当前页
  * @apiParam {int} limit 分页大小
  * @apiParam {string} uuid 分组ID
  * @apiSuccessExample Success-TXResponse:
  *     HTTP/1.1 200 OK
  *     {
			"code": 200,
			"data": [
				{
					"times": 6,  -- 操作次数
					"user_id": 10001, --用户ID
					"visit": "2020-07-09T00:00:00+08:00"  --最后操作时间
				}
			],
			"total": 1 --总行数
  }

*/
func CHUserList(c *gin.Context) {
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	uuid := c.Query("uuid")
	nas := service.NewAnalysisService()
	analysis := nas.GetAnalysisByID(uuid)
	if analysis == nil {
		return
	}
	var conditions []dao.AnalysisConditions
	//fmt.Println(analysis.Conditions)
	_ = json.Unmarshal([]byte(analysis.Conditions), &conditions)
	sql := nas.FormatSql(analysis.Expr, conditions)

	items, _ := dao.NewInfiniteStones().CkQuery(service.FillListHeader(sql, p, limit))
	for i := 0; i < len(items); i++ {
		items[i]["times"] = nas.GetVisitTimes(items[i]["user_id"])
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  items,
		"total": nas.GetTotalMembers(service.FillCountHeader(sql)),
	})

}

/**
  * @api {Get} /admin/v1/analysis/actions 用户行为列表
  * @apiName api.CHUserActionsList
  * @apiGroup analysis
  * @apiHeader {String} Authorization 用户授权token
  * @apiParam {int} p  当前页
  * @apiParam {int} limit  分页大小
  * @apiParam {string} user_id 用户ID
  * @apiParam {string} dt1 时间范围
  * @apiParam {string} dt2 时间范围
  * @apiParam {string} type 0.全部  1.访问类型  2. 其他所有类型
  * @apiSuccessExample Success-TXResponse:
  * HTTP/1.1 200 OK
  * {
    "code": 200,
    "data": [
        {
            "created_at": "2020-07-09T00:00:00+08:00", --时间
            "id": 1,
            "operation_type": 1,
            "resource_author_avatar": "", --访问类型 1 浏览 2 收藏 3分享 4 点赞 5 评论 6 回答 7 打赏 8 站内转发
            "resource_author_id": 10002,
            "resource_author_name": "",
            "resource_author_type": 0,
            "resource_id": 100,  --资源ID
            "resource_introduction": "",
            "resource_label": "",
            "resource_title": "", --资源标题
            "type": 2, --资源类型 1 文章 2 课程&商品 3 免费公开课 4 问题 5 相册 6 APP小课程 7 文章评论 8 话题 9 院校主页 10 个人主页 11 问题回答 12 课程&商品评论 13 APP小课程评论
            "user_id": 10001
        }
    ],
    "total": 6
}
*/
func CHUserActionsList(c *gin.Context) {
	p, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userId, _ := strconv.Atoi(c.DefaultQuery("user_id", "0"))
	dt1 := c.Query("dt1")
	dt2 := c.Query("dt2")
	source, _ := strconv.Atoi(c.DefaultQuery("type", "0"))

	nas := service.NewAnalysisService()
	//dt1, dt2 := service.FormatRange(r)
	items, total := nas.GetAnalysisActions(userId, dt1, dt2, source, p, limit)

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  items,
		"total": total,
	})

}
