package dao

import (
	// "platform_report/config"
	"platform_report/lib"
	"time"
)

type XzqUser struct {
	Id                  int64     `gorm:"column:id;primary_key" json:"id" db:"id"`
	UserId              int64     `json:"userId" form:"userId" gorm:"column:user_id"`
	Avatar              string    `json:"avatar" form:"avatar" gorm:"column:avatar"`                                         // 头像
	RealName            string    `json:"realName" form:"realName" gorm:"column:real_name"`                                  // 姓名
	Nick                string    `json:"nick" form:"nick" gorm:"column:nick"`                                               // 昵称
	Sex                 int       `json:"sex" form:"sex" gorm:"column:sex"`                                                  // 性别 1男2女
	Birthday            time.Time `json:"birthday" form:"birthday" gorm:"column:birthday"`                                   // 生日
	Desctoption         string    `json:"desctoption" form:"desctoption" gorm:"column:desctoption"`                          // 简介
	Mobile              string    `json:"mobile" form:"mobile" gorm:"column:mobile"`                                         // 手机号
	RoleType            int       `json:"roleType" form:"roleType" gorm:"column:role_type"`                                  // 用户类型
	CredentialsType     int       `json:"credentialsType" form:"credentialsType" gorm:"column:credentials_type"`             // 证件类型
	CredentialsNo       string    `json:"credentialsNo" form:"credentialsNo" gorm:"column:credentials_no"`                   // 证件号
	IsRealNameAuth      int       `json:"isRealNameAuth" form:"isRealNameAuth" gorm:"column:is_real_name_auth"`              // 认证状态 0未认证 1已认证
	Country             string    `json:"country" form:"country" gorm:"column:country"`                                      // 国家
	IsJjl               int       `json:"isJjl" form:"isJjl" gorm:"column:is_jjl"`                                           // 是否金吉列
	IsColleges          int       `json:"isColleges" form:"isColleges" gorm:"column:is_colleges"`                            // 是否是合作院校
	InstitutionId       int64     `json:"institutionId" form:"institutionId" gorm:"column:institution_id"`                   // 机构id
	RegistrationChannel int       `json:"registrationChannel" form:"registrationChannel" gorm:"column:registration_channel"` // 注册渠道
	XzqOpenid           string    `json:"xzqOpenid" form:"xzqOpenid" gorm:"column:xzq_openid"`                               // 学长圈openid
	XzqTitle            int       `json:"xzqTitle" form:"xzqTitle" gorm:"column:xzq_title"`                                  // 头衔
	XzqIntroduction     string    `json:"xzqIntroduction" form:"xzqIntroduction" gorm:"column:xzq_introduction"`             // 学长圈简介
	CreatedAt           time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at"`                               // 创建时间

	RoleUuidVal   string `json:"roleUuidVal" form:"roleUuidVal" gorm:"column:role_uuid_val"`       // 角色名称
	RoleUuid      string `json:"roleUuid" form:"roleUuid" gorm:"column:role_uuid"`                 // 角色uuid
	SexVal        string `json:"sexVal" form:"sexVal" gorm:"column:sex_val"`                       // 性别内容
	IsTeacher     int    `json:"isTeacher" form:"isTeacher" gorm:"column:is_teacher"`              // 是否是讲师
	InstitudeId   int64  `json:"institudeId" form:"institudeId" gorm:"column:institude_id"`        // 所属院校id
	Channel       int    `json:"channel" form:"channel" gorm:"column:channel"`                     // 频道
	UserType      int    `json:"userType" form:"userType" gorm:"column:user_type"`                 // 用户类型 1普通用户 2机构
	XzqUserStatus int    `json:"xzqUserStatus" form:"xzqUserStatus" gorm:"column:xzq_user_status"` // 学长圈 用户状态 0 启用 1 禁用
}

func (this *XzqUser) TableName() string {
	return "dz_user"
}
func (this *XzqUser) GetUserInfo(userId int64) error {
	_, err := lib.InitXormMySql().Table(this.TableName()).Where("id = ?", userId).Get(this)
	return err
}

func (this *XzqUser) ModifyUserStatus(userId int64, status int) error {
	_, err := lib.InitXormMySql().Table(this.TableName()).SQL("update "+this.TableName()+" set xzq_user_status = ? where id = ?", status,
		userId).Execute()
	return err
}
