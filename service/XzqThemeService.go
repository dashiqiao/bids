package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"platform_report/config"
	"platform_report/utils"
)

// 文章

type XzqThemeService struct {
}

type DxzThemeData struct {
	Channel        int    `json:"channel"`        // 频道
	CreateUserId   int    `json:"createUserId"`   // 创建人ID
	CreateUserName string `json:"createUserName"` // 创建人姓名
	CreatedAt      int    `json:"createdAt"`      //
	Deleted        string `json:"deleted"`        //
	Id             int    `json:"id"`             //
	IsEnable       int    `json:"isEnable"`       // 是否有效 1有效 0无效
	Sort           int    `json:"sort"`           // 排序
	TagName        string `json:"tagName"`        // 主题名称
	UpdateUserId   int    `json:"updateUserId"`   // 修改人ID
	UpdateUserName int    `json:"updateUserName"` // 修改人姓名
	UpdatedAt      int    `json:"updatedAt"`      //
}

func (this *XzqThemeService) ThemeList(channel string) ([]DxzThemeData, error) {
	resp := make([]DxzThemeData, 0)

	conf := config.Conf{}
	config := conf.GetConf()

	url := config.DxzXzqThemeUrl + fmt.Sprintf("?channel=%s", channel)

	log.Println(url)
	_, bodyByte, err := utils.DoRequest("GET", url, []byte{}, "")
	if err != nil {
		return resp, err
	}

	dxzRespData := new(DxzResponseData)
	err = json.Unmarshal(bodyByte, dxzRespData)
	if err != nil {
		return resp, err
	}

	if !dxzRespData.Success {
		return resp, errors.New(dxzRespData.Message)
	}

	for _, v := range dxzRespData.Result.([]interface{}) {
		data, _ := json.Marshal(v.(map[string]interface{}))
		push := DxzThemeData{}
		json.Unmarshal(data, &push)
		resp = append(resp, push)
	}

	return resp, nil
}

type DxzLabelData struct {
	Channel        int    `json:"channel"`        //
	Code           string `json:"code"`           //
	CreateUserId   int64  `json:"createUserId"`   //
	CreateUserName string `json:"createUserName"` //
	CreatedAt      string `json:"createdAt"`      //
	Id             int    `json:"id"`             //
	Deleted        bool   `json:"deleted"`        //
	IsEnable       string `json:"isEnable"`       // 可用值:0,1
	LabelIcon      string `json:"labelIcon"`      //
	Level          int    `json:"level"`          //
	Name           string `json:"name"`           //
	NameForShort   string `json:"nameForShort"`   //
	ParentCode     string `json:"parentCode"`     //
	Sort           int    `json:"sort"`           //
	UpdateUserId   int    `json:"updateUserId"`   //
	UpdateUserName string `json:"updateUserName"` //
	UpdatedAt      string `json:"updatedAt"`      //
}

func (this *XzqThemeService) FirstLabelList(channel string) ([]DxzLabelData, error) {
	resp := make([]DxzLabelData, 0)

	conf := config.Conf{}
	config := conf.GetConf()

	url := config.DxzXzqFirstLabelListUrl + fmt.Sprintf("?channel=%s", channel)
	_, bodyByte, err := utils.DoRequest("GET", url, []byte{}, "")
	if err != nil {
		return resp, err
	}

	dxzRespData := new(DxzResponseData)
	err = json.Unmarshal(bodyByte, dxzRespData)
	if err != nil {
		return resp, err
	}

	if !dxzRespData.Success {
		return resp, errors.New(dxzRespData.Message)
	}

	for _, v := range dxzRespData.Result.([]interface{}) {
		data, _ := json.Marshal(v.(map[string]interface{}))
		push := DxzLabelData{}
		json.Unmarshal(data, &push)
		resp = append(resp, push)
	}

	return resp, nil
}

func (this *XzqThemeService) SecondLabelList(channel string) ([]DxzLabelData, error) {
	resp := make([]DxzLabelData, 0)

	conf := config.Conf{}
	config := conf.GetConf()

	url := config.DxzXzqSecondLabelListUrl + fmt.Sprintf("?pCode=%v", channel)
	_, bodyByte, err := utils.DoRequest("GET", url, []byte{}, "")
	if err != nil {
		return resp, err
	}

	dxzRespData := new(DxzResponseData)
	err = json.Unmarshal(bodyByte, dxzRespData)
	if err != nil {
		return resp, err
	}

	if !dxzRespData.Success {
		return resp, errors.New(dxzRespData.Message)
	}

	for _, v := range dxzRespData.Result.([]interface{}) {
		data, _ := json.Marshal(v.(map[string]interface{}))
		push := DxzLabelData{}
		json.Unmarshal(data, &push)
		resp = append(resp, push)
	}

	return resp, nil
}
