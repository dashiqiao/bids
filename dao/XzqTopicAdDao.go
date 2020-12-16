package dao

import (
	"platform_report/lib"
	"time"
)

type DzTopic struct {
	Id            int64 `gorm:"AUTO_INCREMENT"`
	Title         string
	Content       string
	ExamineStatus int `gorm:"default:'1'"`
	TopicType     int `gorm:"default:'2'"`
	UserId        int
	CntView       int `gorm:"-"`
	CntShare      int `gorm:"-"`
	CntCollection int `gorm:"-"`
	CntAnswer     int `gorm:"-"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DzTopicAdArea struct {
	Id        int `gorm:"AUTO_INCREMENT"`
	TopicId   int64
	CountryId int
	CityId    int
}

type DzTopicAdData struct {
	Id          int `gorm:"AUTO_INCREMENT"`
	TopicId     int64
	DzType      int
	DxzId       int
	ChannelName string
	Cover       string
	FirstLabel  string
	SecondLabel string
}

func CreateTopicAd(ad DzTopic, areas []DzTopicAdArea, datas []DzTopicAdData) (int64, error) {
	tx := lib.GetDbInstance().Begin()
	err := tx.Create(&ad).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, area := range areas {
		area.TopicId = ad.Id
		err := tx.Create(&area).Error
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	for _, data := range datas {
		data.TopicId = ad.Id
		data.DzType = 2
		err := tx.Create(&datas).Error
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return ad.Id, nil
}

func GetTopicAdByID(id int64) (DzTopic, []DzTopicAdArea, []DzTopicAdData) {
	var (
		item DzTopic
		area []DzTopicAdArea
		data []DzTopicAdData
	)
	lib.GetDbInstance().Model(DzTopic{}).Where("id = ?", id).Scan(&item)
	if item.Id == 0 {
		return item, area, data
	}
	lib.GetDbInstance().Where("topic_id = ?", id).Find(&area)
	lib.GetDbInstance().Where("topic_id = ?", id).Find(&data)
	return item, area, data
}
