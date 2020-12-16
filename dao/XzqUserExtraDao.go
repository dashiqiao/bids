package dao

import (
	// "platform_report/config"
	"platform_report/lib"
	// "time"
)

type XzqUserExtra struct {
	Id                int64   `gorm:"column:id;primary_key" json:"id" db:"id"`
	UserId            int64   `gorm:"column:user_id" json:"userId" db:"user_id"`                            // 用户id
	DxzId             int64   `gorm:"column:dxz_id" json:"dxzId" db:"dxz_id"`                               // 大学长id
	Country           int64   `gorm:"column:country" json:"country" db:"country"`                           // 国家
	City              int64   `gorm:"column:city" json:"city" db:"city"`                                    // 城市
	School            int64   `gorm:"column:school" json:"school" db:"school"`                              // 学校
	SystemLanguage    int64   `gorm:"column:system_language" json:"systemLanguage" db:"system_language"`    // 系统语言
	Title             int     `json:"title" form:"title" gorm:"column:title"`                               // 学长圈头衔
	ContributionPoint float64 `json:"contributionPoint" form:"contributionPoint" gorm:"contribution_point"` // 贡献点
	Integral          float64 `json:"intergral" form:"intergral" gorm:"intergral"`                          // 积分
	Fans              int     `json:"fans" form:"fans" gorm:"fans"`                                         // 粉丝
	Attention         int     `json:"attention" form:"attention" gorm:"attention"`                          // 关注
	ClockIn           int     `json:"clockIn" form:"clockIn" gorm:"clock_in"`                               // 打卡
	Task              int     `json:"task" form:"task" gorm:"task"`                                         // 任务完成数
	TagIds            string  `json:"tagIds" form:"tagIds" gorm:"tag_ids"`                                  // 标签列表
	SchoolSection     int     `json:"schoolSection" form:"schoolSection" gorm:"school_section"`             // 学段

}

func (this *XzqUserExtra) TableName() string {
	return "dz_user_extra"
}
func (this *XzqUserExtra) GetInfo(id int64) error {
	_, err := lib.InitXormMySql().Table(this.TableName()).Where("id = ?", id).Get(this)
	return err
}

func (this *XzqUserExtra) GetInfoByUserId(userId int64) error {
	_, err := lib.InitXormMySql().Table(this.TableName()).Where("user_id = ?", userId).Get(this)
	return err
}
