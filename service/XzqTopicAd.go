package service

import (
	"platform_report/dao"
	"time"
)

type CreateTopicAdReq struct {
	Title   string        `json:"title"`
	Content string        `json:"content"`
	Data    []TopicAdData `json:"data"`
	City    []TopicAdCity `json:"city"`
}

type TopicAdData struct {
	Id          int    `json:"id"`
	Type        int    `json:"type"` //1.课程 2.文章
	ChannelName string `json:"channel_name"`
	FirstLabel  string `json:"first_label"`
	SecondLabel string `json:"second_label"`
	Uri         string `json:"uri"`
}

type TopicAdCity struct {
	Country int `json:"country"`
	City    int `json:"city"`
}

func CreateTopicAd(ad CreateTopicAdReq, userId int) error {
	req := dao.DzTopic{Title: ad.Title, Content: ad.Content, UserId: userId,
		CreatedAt: time.Now(), UpdatedAt: time.Now()}
	area := make([]dao.DzTopicAdArea, 0, len(ad.City))
	data := make([]dao.DzTopicAdData, 0, len(ad.Data))
	for _, item := range ad.City {
		area = append(area, dao.DzTopicAdArea{
			CountryId: item.Country,
			CityId:    item.City,
		})
	}
	for _, item := range ad.Data {
		data = append(data, dao.DzTopicAdData{
			DxzId:       item.Id,
			ChannelName: item.ChannelName,
			FirstLabel:  item.FirstLabel,
			SecondLabel: item.SecondLabel,
			Cover:       item.Uri,
		})
	}
	subId, err := dao.CreateTopicAd(req, area, data)
	if err != nil {
		return err
	}
	_, _ = dao.RegisterTopic(int64(TopicAd), subId, int64(userId))
	return nil
}

func GetTopicAdByID(id int64) (dao.DzTopic, []dao.DzTopicAdArea, []dao.DzTopicAdData) {
	return dao.GetTopicAdByID(id)
}
