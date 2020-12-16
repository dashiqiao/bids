package service

import (
	"errors"
	"platform_report/dao"
	"time"
)

type XzqUserService struct {
	UserId         int64  `json:"UserId"`         // 用户id
	Mobile         string `json:"mobile"`         // 手机号
	Nick           string `json:"nick"`           // 用户昵称
	UserName       string `json:"userName"`       // 用户名称
	UserType       int    `json:"userType"`       // 用户类型
	UserStatus     int    `json:"userStatus"`     // 用户状态
	CreateStarTime string `json:"createStarTime"` // 用户创建起始时间
	CreateEndTime  string `json:"createEndTime"`
	Country        int    `json:"country"` // 国家
	City           int    `json:"city"`    // 城市
	School         int    `json:"school"`  // 学校
}

type UserRespData struct {
	UserId            int64     `json:"userId"`
	Avatar            string    `json:"avatar" form:"avatar" gorm:"column:avatar"`                             // 头像
	RealName          string    `json:"realName" form:"realName" gorm:"column:real_name"`                      // 姓名
	Nick              string    `json:"nick" form:"nick" gorm:"column:nick"`                                   // 昵称
	Sex               int       `json:"sex" form:"sex" gorm:"column:sex"`                                      // 性别 1男2女
	Birthday          time.Time `json:"birthday" form:"birthday" gorm:"column:birthday"`                       // 生日
	Desctoption       string    `json:"desctoption" form:"desctoption" gorm:"column:desctoption"`              // 简介
	Mobile            string    `json:"mobile" form:"mobile" gorm:"column:mobile"`                             // 手机号
	RoleType          int       `json:"roleType" form:"roleType" gorm:"column:role_type"`                      // 用户类型
	CredentialsType   int       `json:"credentialsType" form:"credentialsType" gorm:"column:credentials_type"` // 证件类型
	CredentialsNo     string    `json:"credentialsNo" form:"credentialsNo" gorm:"column:credentials_no"`       // 证件号
	IsRealNameAuth    int       `json:"isRealNameAuth" form:"isRealNameAuth" gorm:"column:is_real_name_auth"`  // 认证状态 0未认证 1已认证
	XzqUserStatus     int       `json:"xzqUserStatus" form:"xzqUserStatus" gorm:"column:xzq_user_status"`      // 学长圈 用户状态 0 启用 1禁用
	UserType          int       `json:"userType" form:"userType" gorm:"column:user_type"`                      // 用户类型 1普通用户 2讲师 3顾问 4机构 5院校
	XzqTitle          int       `json:"xzqTitle" form:"xzqTitle" gorm:"column:xzq_title"`                      // 头衔
	CreatedAt         time.Time `json:"createdAt" form:"createdAt" gorm:"column:created_at"`                   // 创建时间
	ContributionPoint float64   `json:"contributionPoint" form:"contributionPoint" gorm:"contribution_point"`  // 贡献点
	Integral          float64   `json:"intergral" form:"intergral" gorm:"intergral"`                           // 积分
	Fans              int       `json:"fans" form:"fans" gorm:"fans"`                                          // 粉丝
	Attention         int       `json:"attention" form:"attention" gorm:"attention"`                           // 关注
	ClockIn           int       `json:"clockIn" form:"clockIn" gorm:"clock_in"`                                // 打卡
	Task              int       `json:"task" form:"task" gorm:"task"`                                          // 任务完成数
	Topic             int       `json:"topic"`                                                                 // 发布话题数
	Dynamic           int       `json:"dynamic"`                                                               // 发布动态数
	Raiders           int       `json:"raiders"`                                                               // 发布攻略
	Country           string    `json:"country"`                                                               // 国家
	City              string    `json:"city"`                                                                  // 城市
	School            string    `json:"school"`                                                                // 学校
	SchoolSection     int       `json:"schoolSection" form:"schoolSection" gorm:"school_section"`              // 学段
}

func (this *XzqUserService) UserInfo() (UserRespData, error) {
	resp := UserRespData{}

	user := new(dao.XzqUser)
	err := user.GetUserInfo(this.UserId)
	if err != nil {
		return resp, err
	}

	extra := new(dao.XzqUserExtra)
	err = extra.GetInfoByUserId(this.UserId)
	if err != nil {
		return resp, err
	}

	resp.UserId = user.Id
	resp.Avatar = user.Avatar
	resp.RealName = user.RealName
	resp.Nick = user.Nick
	resp.Sex = user.Sex
	resp.Birthday = user.Birthday
	resp.Desctoption = user.Desctoption
	resp.Mobile = user.Mobile
	resp.RoleType = user.RoleType
	resp.CredentialsType = user.CredentialsType
	resp.CredentialsNo = user.CredentialsNo
	resp.IsRealNameAuth = user.IsRealNameAuth
	resp.XzqUserStatus = user.XzqUserStatus
	resp.UserType = user.UserType
	resp.XzqTitle = user.XzqTitle
	resp.CreatedAt = user.CreatedAt
	// resp.ContributionPoint=user

	resp.ContributionPoint = extra.ContributionPoint
	resp.Integral = extra.Integral
	resp.Fans = extra.Fans
	resp.Attention = extra.Attention
	resp.ClockIn = extra.ClockIn
	resp.Task = extra.Task
	resp.SchoolSection = extra.SchoolSection
	return resp, nil
}

type UserStatusData struct {
	UserId int64 `json:"user_id"`
	Status int   `json:"status"`
}

func (this *XzqUserService) ModifyUserStatus(data *UserStatusData) error {
	if data.UserId <= 0 {
		return errors.New("错误的用户!")
	}
	if data.Status != 1 && data.Status != 0 {
		return errors.New("错误的状态!")
	}
	err := new(dao.XzqUser).ModifyUserStatus(data.UserId, data.Status)
	return err
}
