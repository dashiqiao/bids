package service

import (
	"encoding/json"
	"errors"
	"log"

	"platform_report/config"
	"platform_report/utils"
)

// 文章

type XzqEssayService struct {
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	Title       string `json:"title"`       // 标题
	Id          int    `json:"id"`          // 文章id
	RealName    string `json:"realName"`    // 用户姓名
	Nick        string `json:"nick"`        // 用户昵称
	Channel     int    `json:"channel"`     // 频道
	FirstLabel  string `json:"firstLabel"`  // 一级标签
	SecondLabel string `json:"secondLabel"` // 二级标签
	ContentTag  int    `json:"contentTag"`  // 内容主题
	State       int    `json:"state"`       // 状态
	RoleUuid    int    `json:"roleUuid"`    // 用户角色
	StartTime   string `json:"startTime"`   // 发布开始时间
	EndTime     string `json:"endTime"`     // 发布结束时间
}

type DxzEssayResData struct {
	PageNo    int `json:"pageNo"`
	PageSize  int `json:"pageSize"`
	Condition struct {
		Title       string `json:"title,omitempty"`       // 标题
		Id          int    `json:"id,omitempty"`          // 文章id
		RealName    string `json:"realName,omitempty"`    // 用户姓名
		Nick        string `json:"nick,omitempty"`        // 用户昵称
		Channel     int    `json:"channel,omitempty"`     // 频道
		FirstLabel  string `json:"firstLabel,omitempty"`  // 一级标签
		SecondLabel string `json:"secondLabel,omitempty"` // 二级标签
		ContentTag  int    `json:"contentTag,omitempty"`  // 内容主题
		State       int    `json:"state,omitempty"`       // 状态
		RoleUuid    int    `json:"roleUuid,omitempty"`    // 用户角色
		StartTime   string `json:"startTime,omitempty"`   // 发布开始时间
		EndTime     string `json:"endTime,omitempty"`     // 发布结束时间
	} `json:"condition"`
}

type DxzEssayData struct {
	Channel        int      `json:"channel"`        // 频道类型
	ChannelToDesc  string   `json:"channelToDesc"`  // 频道类型名称
	ContentTag     []string `json:"contentTag"`     // 内容主题
	FirstLabel     int64    `json:"firstLabel"`     // 一级标签id
	FirstLabelCode string   `json:"firstLabelCode"` // 一级标签code
	FirstLabelVal  string   `json:"firstLabelVal"`  // 一级标签val
	Id             int64    `json:"id"`             // 文章ID
	Nick           string   `json:"nick"`           // 昵称
	RealName       string   `json:"realName"`       // 真实姓名
	ReleaseTime    string   `json:"releaseTime"`    // 发布时间
	RoleUuid       string   `json:"roleUuid"`       // 用户角色
	SecondLabel    []string `json:"secondLabel"`    // 文章分类标签
	State          int      `json:"state"`          // 状态
	Title          string   `json:"title"`          // 文章标题
	ImgUrl         string   `json:"imgUrl"`         // 文章封面
}

func (this *XzqEssayService) EssayList() ([]DxzEssayData, error) {
	resp := make([]DxzEssayData, 0)

	conf := config.Conf{}
	config := conf.GetConf()

	data := new(DxzEssayResData)
	data.PageNo = this.Page
	data.PageSize = this.Limit
	data.Condition.Title = this.Title
	data.Condition.Id = this.Id
	data.Condition.RealName = this.RealName
	data.Condition.Nick = this.Nick
	data.Condition.Channel = this.Channel
	data.Condition.FirstLabel = this.FirstLabel
	data.Condition.SecondLabel = this.SecondLabel
	data.Condition.ContentTag = this.ContentTag
	data.Condition.State = this.State
	data.Condition.RoleUuid = this.RoleUuid
	data.Condition.StartTime = this.StartTime
	data.Condition.EndTime = this.EndTime

	dataByte, err := json.Marshal(data)
	if err != nil {
		return resp, err
	}

	_, bodyByte, err := utils.DoRequest("POST", config.DxzXzqEssayUrl, dataByte, "")
	if err != nil {
		return resp, err
	}

	dxzRespData := new(DxzResponseData)
	err = json.Unmarshal(bodyByte, dxzRespData)
	if err != nil {
		return resp, err
	}

	log.Println(string(bodyByte))

	if !dxzRespData.Success {
		return resp, errors.New(dxzRespData.Message)
	}

	for _, v := range dxzRespData.Result.([]interface{}) {
		data, _ := json.Marshal(v.(map[string]interface{}))
		push := DxzEssayData{}
		json.Unmarshal(data, &push)
		resp = append(resp, push)
	}

	return resp, nil
}
