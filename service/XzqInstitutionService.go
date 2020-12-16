package service

import (
	"platform_report/dao"
	// "time"
)

type XzqInstitutionService struct {
	UserId         int64  `json:"userId"`         // 机构id
	Name           string `json:"name"`           // 机构名称
	UserStatus     int    `json:"userStatus"`     // 状态
	Country        int    `json:"country"`        // 国家
	City           int    `json:"city"`           // 城市
	School         int    `json:"school"`         // 学校
	CreateStarTime string `json:"createStarTime"` // 用户创建起始时间
	CreateEndTime  string `json:"createEndTime"`
}

type InstitutionRestData struct {
	UserId            int64   `json:"userId"`
	Avatar            string  `json:"avatar" form:"avatar" gorm:"column:avatar"`                            // 头像
	RealName          string  `json:"realName" form:"realName" gorm:"column:real_name"`                     // 姓名
	Nick              string  `json:"nick" form:"nick" gorm:"column:nick"`                                  // 昵称
	Mobile            string  `json:"mobile" form:"mobile" gorm:"column:mobile"`                            // 手机号
	IsRealNameAuth    int     `json:"isRealNameAuth" form:"isRealNameAuth" gorm:"column:is_real_name_auth"` // 认证状态 0未认证 1已认证
	XzqUserStatus     int     `json:"xzqUserStatus" form:"xzqUserStatus" gorm:"column:xzq_user_status"`     // 学长圈 用户状态 0 启用 1禁用
	UserType          int     `json:"userType" form:"userType" gorm:"column:user_type"`                     // 用户类型 1普通用户 2讲师 3顾问 4机构 5院校
	ContributionPoint float64 `json:"contributionPoint" form:"contributionPoint" gorm:"contribution_point"` // 贡献点
	Integral          float64 `json:"intergral" form:"intergral" gorm:"intergral"`                          // 积分
	Fans              int     `json:"fans" form:"fans" gorm:"fans"`                                         // 粉丝
	Attention         int     `json:"attention" form:"attention" gorm:"attention"`                          // 关注
	ClockIn           int     `json:"clockIn" form:"clockIn" gorm:"clock_in"`                               // 打卡
	Task              int     `json:"task" form:"task" gorm:"task"`                                         // 任务完成数
	Topic             int     `json:"topic"`                                                                // 发布话题数
	Dynamic           int     `json:"dynamic"`                                                              // 发布动态数
	Raiders           int     `json:"raiders"`                                                              // 发布攻略
	Country           string  `json:"country"`                                                              // 国家
	City              string  `json:"city"`                                                                 // 城市
	School            string  `json:"school"`                                                               // 学校
	Title             string  `json:"title"`                                                                // 名称
}

func (this *XzqInstitutionService) InstitutionInfo() (*InstitutionRestData, error) {
	resp := new(InstitutionRestData)
	user := new(dao.XzqUser)
	err := user.GetUserInfo(this.UserId)
	if err != nil {
		return resp, err
	}

	institution := new(dao.XzqInstitution)
	err = institution.GetInfoByUserId(this.UserId)
	if err != nil {
		return resp, err
	}
	extra := new(dao.XzqUserExtra)
	err = extra.GetInfoByUserId(this.UserId)
	if err != nil {
		return resp, err
	}

	resp.Avatar = user.Avatar
	resp.UserId = user.Id
	resp.RealName = user.RealName
	resp.Nick = user.Nick
	resp.Mobile = user.Mobile
	resp.IsRealNameAuth = user.IsRealNameAuth
	resp.XzqUserStatus = user.XzqUserStatus
	resp.UserType = user.UserType
	resp.ContributionPoint = extra.ContributionPoint
	resp.Integral = extra.Integral
	resp.Fans = extra.Fans
	resp.Attention = extra.Attention
	resp.ClockIn = extra.ClockIn
	resp.Task = extra.Task
	resp.Country = ""
	resp.City = ""
	resp.School = ""
	resp.Title = institution.Title

	return resp, nil
}
