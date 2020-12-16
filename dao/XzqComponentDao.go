package dao

import (
	"fmt"
	"github.com/spf13/cast"
	"platform_report/lib"
	"time"
)

type DzComponentTopic struct {
	Id            int64 `gorm:"AUTO_INCREMENT"`
	AppId         int64
	SubId         int64
	TagOpen       int
	UserId        int64
	CntPraise     int
	CntReply      int
	CntShare      int
	CntAnswer     int
	CntView       int
	CntCollection int
	CreatedAt     time.Time
}

func GetComponentTopicKey(appId, subId int64) string {
	return fmt.Sprintf("component:topic:%d:%d", appId, subId)
}

func RegisterTopic(appId, subId, uid int64) (DzComponentTopic, error) {
	var dct DzComponentTopic
	total := cast.ToInt(lib.GetRedis().HLen(GetComponentTopicKey(appId, subId)).String())
	if total > 0 {
		return dct, nil
	}
	lib.GetDbInstance().Model(DzComponentTopic{}).Where("app_id = ? and sub_id = ?", appId, subId).Scan(&dct)
	if dct.AppId > 0 {
		_ = SaveComponentTopicToCache(appId, subId, dct)
		return dct, nil
	}
	dct = DzComponentTopic{AppId: appId, SubId: subId, TagOpen: 256,
		UserId: uid, CreatedAt: time.Now()}
	e := lib.GetDbInstance().Create(&dct).Error
	if e == nil {
		_ = SaveComponentTopicToCache(appId, subId, dct)
	}
	return dct, nil
}

func SaveComponentTopicToCache(appId, subId int64, dct DzComponentTopic) error {
	key := GetComponentTopicKey(appId, subId)
	_ = lib.GetRedis().HSet(key, "id", dct.Id)
	_ = lib.GetRedis().HSet(key, "app_id", dct.AppId)
	_ = lib.GetRedis().HSet(key, "sub_id", dct.SubId)
	_ = lib.GetRedis().HSet(key, "tag_open", dct.TagOpen)
	_ = lib.GetRedis().HSet(key, "user_id", dct.UserId)
	_ = lib.GetRedis().HSet(key, "cnt_praise", dct.CntPraise)
	_ = lib.GetRedis().HSet(key, "cnt_reply", dct.CntReply)
	_ = lib.GetRedis().HSet(key, "cnt_share", dct.CntShare)
	_ = lib.GetRedis().HSet(key, "cnt_answer", dct.CntAnswer)
	_ = lib.GetRedis().HSet(key, "cnt_view", dct.CntView)
	_ = lib.GetRedis().HSet(key, "cnt_collection", dct.CntCollection)
	return nil
}

func GetTopic(appId, subId int64) (*DzComponentTopic, error) {
	items, err := GetMultiTopic(appId, subId)
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items[0], nil
	}
	return nil, nil
}

func GetMultiTopic(appId int64, subId ...int64) ([]*DzComponentTopic, error) {
	items := make([]*DzComponentTopic, 0, len(subId))
	for _, val := range subId {
		mp, err := lib.GetRedis().HGetAll(GetComponentTopicKey(appId, val)).Result()
		if err != nil {
			continue
		}
		item := new(DzComponentTopic)
		if len(mp) == 0 {
			topic, err := RegisterTopic(appId, val, 0)
			if err != nil {
				continue
			}
			item.Id = topic.Id
			item.AppId = topic.AppId
			item.SubId = topic.SubId
			item.TagOpen = topic.TagOpen
			item.UserId = topic.UserId
			item.CntPraise = topic.CntPraise
			item.CntReply = topic.CntReply
			item.CntShare = topic.CntShare
			item.CntAnswer = topic.CntAnswer
			item.CntView = topic.CntView
			item.CntCollection = topic.CntCollection
			item.CreatedAt = topic.CreatedAt
		} else {
			for k, v := range mp {
				switch k {
				case "id":
					item.Id = cast.ToInt64(v)
				case "app_id":
					item.AppId = cast.ToInt64(v)
				case "sub_id":
					item.SubId = cast.ToInt64(v)
				case "tag_open":
					item.TagOpen = cast.ToInt(v)
				case "user_id":
					item.UserId = cast.ToInt64(v)
				case "cnt_praise":
					item.CntPraise = cast.ToInt(v)
				case "cnt_reply":
					item.CntReply = cast.ToInt(v)
				case "cnt_share":
					item.CntShare = cast.ToInt(v)
				case "cnt_answer":
					item.CntAnswer = cast.ToInt(v)
				case "cnt_view":
					item.CntView = cast.ToInt(v)
				case "cnt_collection":
					item.CntCollection = cast.ToInt(v)
				}
			}
		}
		items = append(items, item)
	}
	return items, nil
}

func ComponentTopicNumDo(appId, subId, num int64, field string) error {
	sql := ""
	if num > 0 {
		sql = fmt.Sprintf(" %v= %v + %v", field, field, num)
	} else {
		sql = fmt.Sprintf(" %v= %v - %v", field, field, num*-1)
	}
	err := lib.GetDbInstance().Exec("update dz_component_topic set "+sql+" where app_id = ? and sub_id = ?", appId, subId).Error
	if err != nil {
		fmt.Println("11111", err.Error())
		return err
	}
	lib.GetRedis().HIncrBy(GetComponentTopicKey(appId, subId), field, num)
	return nil
}

func IncrementComponentTopicPraise(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, 1, "cnt_praise")
}

func DecrementComponentTopicPraise(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, -1, "cnt_praise")
}

func IncrementComponentTopicReply(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, 1, "cnt_reply")
}

func DecrementComponentTopicReply(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, -1, "cnt_reply")
}

func IncrementComponentTopicShare(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, 1, "cnt_share")
}

func IncrementComponentTopicView(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, 1, "cnt_view")
}

func IncrementComponentTopicAnswer(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, 1, "cnt_answer")
}

func DecrementComponentTopicAnswer(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, -1, "cnt_answer")
}

func IncrementComponentTopicCollection(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, 1, "cnt_collection")
}

func DecrementComponentTopicCollection(appId, subId int64) error {
	return ComponentTopicNumDo(appId, subId, -1, "cnt_collection")
}
