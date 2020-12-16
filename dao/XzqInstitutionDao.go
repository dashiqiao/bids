package dao

import (
	// "platform_report/config"
	"platform_report/lib"
	// "time"
)

type XzqInstitution struct {
	Id                int64  `gorm:"column:id;primary_key" json:"id" db:"id"`
	UserId            int64  `json:"userId"`
	Avatar            string `json:"avatar"`            // 用户头像
	BuiltAt           string `json:"builtAt"`           // 建校时间
	BuiltAtYear       string `json:"builtAtYear"`       // 建校年
	Channel           int    `json:"channel"`           // 所属频道
	CountryCode       string `json:"countryCode"`       // 国家代码
	CountryName       string `json:"countryName"`       // 国家名称
	CountryNameEn     string `json:"countryNameEn"`     // 国家名称(英文)
	CreateAt          string `json:"createAt"`          // 创建时间
	InstitudeId       int64  `json:"institudeId"`       // 所属机构ID，为子账号时值为机构&院校 ID，为机构&院校主账号时此字段为-1
	InstitudeType     int    `json:"institudeType"`     // 类型(1:机构2:院校)
	InstitudeTypeName string `json:"institudeTypeName"` // 类型名称
	IsCooperative     int    `json:"isCooperative"`     // 是否合作机构
	IsDxz             int    `json:"isDxz"`             // 是否大学长
	IsJjl             int    `json:"isJjl"`             // 是否金吉列
	IsSubAccount      int    `json:"isSubAccount"`      // 是否子账号(0:否1:是)
	Nick              string `json:"nick"`              // 用户昵称
	RealName          string `json:"realName"`          // 真实姓名
	Title             string `json:"title"`             // 机构&院校名称
	TitleEn           string `json:"titleEn"`           // 机构&院校名称(英文)
}

func (this *XzqInstitution) TableName() string {
	return "dz_institution"
}
func (this *XzqInstitution) GetInfo(id int64) error {
	_, err := lib.InitXormMySql().Table(this.TableName()).Where("id = ?", id).Get(this)
	return err
}

func (this *XzqInstitution) GetInfoByUserId(userId int64) error {
	_, err := lib.InitXormMySql().Table(this.TableName()).Where("user_id = ?", userId).Get(this)
	return err
}
