package yyzx

import (
	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/strftime"
	"github.com/spf13/cast"
	"io/ioutil"
	"platform_report/dao"
	v1 "platform_report/routers/api/v1"
	"platform_report/service"
)

func CreateTopicAd(c *gin.Context) {
	svc := service.CreateTopicAdReq{}
	err := c.ShouldBind(&svc)
	if err != nil {
		v1.ErrorStd(c, err.Error())
		return
	}
	err = service.CreateTopicAd(svc, v1.UserInfo(c).ID)
	if err != nil {
		v1.ErrorStd(c, err.Error())
		return
	}
	v1.SuccessStds(c)
}

func ViewTopicAd(c *gin.Context) {
	id := cast.ToInt64(c.Query("id"))
	ad, area, data := service.GetTopicAdByID(id)
	compontent, _ := dao.GetTopic(int64(service.TopicAd), id)
	if compontent != nil {
		ad.CntShare = compontent.CntShare
		ad.CntAnswer = compontent.CntAnswer
		ad.CntView = compontent.CntView
		ad.CntCollection = compontent.CntCollection
	}
	v1.SuccessStd(c, map[string]interface{}{
		"ad":   ad,
		"area": area,
		"data": data,
	})
}

func Examine(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	appId := fsion.Get("app_id").ValueInt()
	id := fsion.Get("id").ValueInt()
	state := fsion.Get("state").ValueInt() //1.通过 2.拒绝
	msg := fsion.Get("msg").ValueStr()
	amount := 50
	switch appId {
	case service.TopicAppid:
		func() {
			object := dao.GetTopicObject(id)
			if object.Id == 0 {
				v1.ErrorStd(c, "操作失败：记录404")
				return
			}
			if object.ExamineStatus != 0 {
				v1.SuccessStds(c)
				return
			}
			if state == 1 {
				err := dao.ExamineTopic(id, dao.Allow, "")
				if err != nil {
					v1.ErrorStd(c, err.Error())
					return
				}
				_ = dao.AssetIncrement(service.TopicAppid, id, object.UserId, dao.GXZ, amount, "", "审核通过,增加贡献值!")
				clockTime, _ := strftime.Format("%F", object.CreatedAt)
				_ = dao.ClockIn(object.UserId, clockTime)
				v1.SuccessStds(c)
			} else {
				err := dao.ExamineTopic(id, dao.NotAllow, msg)
				if err != nil {
					v1.ErrorStd(c, err.Error())
					return
				}
				v1.SuccessStds(c)
			}
		}()
	case service.Dynamic:
		func() {
			object := dao.GetDynamicObject(id)
			if object.Id == 0 {
				v1.ErrorStd(c, "操作失败：记录404")
				return
			}
			if object.ExamineStatus != 0 {
				v1.SuccessStds(c)
				return
			}
			if state == 1 {
				err := dao.ExamineDynamic(id, dao.Allow, "")
				if err != nil {
					v1.ErrorStd(c, err.Error())
					return
				}
				_ = dao.AssetIncrement(service.Dynamic, id, object.UserId, dao.GXZ, amount, "", "审核通过,增加贡献值!")
				clockTime, _ := strftime.Format("%F", object.CreatedAt)
				_ = dao.ClockIn(object.UserId, clockTime)
				v1.SuccessStds(c)
			} else {
				err := dao.ExamineDynamic(id, dao.NotAllow, msg)
				if err != nil {
					v1.ErrorStd(c, err.Error())
					return
				}
				v1.SuccessStds(c)
			}
		}()
	case service.StrategyAppid:
		func() {
			object := dao.GetStrategyObject(id)
			if object.Id == 0 {
				v1.ErrorStd(c, "操作失败：记录404")
				return
			}
			if object.ExamineStatus != 0 {
				v1.SuccessStds(c)
				return
			}
			if state == 1 {
				err := dao.ExamineStrategy(id, dao.Allow, "")
				if err != nil {
					v1.ErrorStd(c, err.Error())
					return
				}
				_ = dao.AssetIncrement(service.StrategyAppid, id, object.UserId, dao.GXZ, amount, "", "审核通过,增加贡献值!")
				clockTime, _ := strftime.Format("%F", object.CreatedAt)
				_ = dao.ClockIn(object.UserId, clockTime)
				v1.SuccessStds(c)
			} else {
				err := dao.ExamineStrategy(id, dao.NotAllow, msg)
				if err != nil {
					v1.ErrorStd(c, err.Error())
					return
				}
				v1.SuccessStds(c)
			}
		}()
	}
}

func UpDown(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	appId := fsion.Get("app_id").ValueInt()
	id := fsion.Get("id").ValueInt()
	state := fsion.Get("state").ValueInt() //0.上架 1.下架
	err := UpDownAgent(appId, id, state)
	if err != nil {
		v1.ErrorStd(c, err.Error())
		return
	}
	v1.SuccessStds(c)
}

func JbUpDown(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	appId := fsion.Get("app_id").ValueInt()
	subId := fsion.Get("sub_id").ValueInt()
	comment := fsion.Get("comment").ValueStr()
	state := fsion.Get("state").ValueInt() //0.上架 1.下架
	err := dao.AddJbAuditing(appId, subId, comment, state, v1.UserInfo(c).ID, v1.UserInfo(c).Username)
	if err != nil {
		v1.ErrorStd(c, err.Error())
		return
	}
	_ = UpDownAgent(appId, subId, state)

	v1.SuccessStds(c)
}

func JbExamine(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	fsion := gofasion.NewFasion(string(body))
	id := fsion.Get("id").ValueInt()
	appId := fsion.Get("app_id").ValueInt()
	subId := fsion.Get("sub_id").ValueInt()
	state := fsion.Get("state").ValueInt() //1.通过 2.拒绝
	category := fsion.Get("category").ValueInt()
	if id > 0 {
		_ = dao.ExamineJbById(id, state)
	} else {
		_ = dao.ExamineJbByCategory(appId, subId, state, category)
	}
	v1.SuccessStds(c)
}

func UpDownAgent(appId, id, state int) error {
	var ex dao.ExamineState
	var err error
	if state == 1 {
		ex = dao.Up
	} else {
		ex = dao.Down
	}
	switch appId {
	case service.TopicAppid:
		err = dao.ShelfTopic(id, ex)
	case service.Dynamic:
		err = dao.ShelfDynamic(id, ex)
	case service.StrategyAppid:
		err = dao.ShelfStrategy(id, ex)
	case service.TopicAnswerAppid:
		err = dao.ShelfAnswer(id, ex)
	}
	return err
}
