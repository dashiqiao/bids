package dao

import (
	"platform_report/lib"
	"time"
)

type DzClockLog struct {
	Id        int `gorm:"AUTO_INCREMENT"`
	UserId    int64
	ClockDate string
	CreatedAt time.Time
}

func ClockIn(uid int64, clockDate string) error {
	log := DzClockLog{UserId: uid, ClockDate: clockDate, CreatedAt: time.Now()}
	return lib.GetDbInstance().Create(&log).Error
}
