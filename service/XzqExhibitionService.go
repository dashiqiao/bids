package service

import (
	"encoding/json"
	"errors"

	"platform_report/config"
	"platform_report/utils"
)

// 展会

type XzqExhibitionService struct {
	Page             int    `json:"page"`
	Limit            int    `json:"limit"`
	ExhibitionNumber string `json:"exhibitionNumber"` // 展会编号
	ExhibitionName   string `json:"exhibitionName"`   // 展会名称
	ReviewStatus     int    `json:"reviewStatus"`     // 审核状态
	Status           int    `json:"status"`           // 展会状态
}

type DxzExhibitionResData struct {
	PageNo    int `json:"pageNo"`
	PageSize  int `json:"pageSize"`
	Condition struct {
		ExhibitionNumber string `json:"exhibitionNumber,omitempty"` // 展会编号
		ExhibitionName   string `json:"exhibitionName,omitempty"`   // 展会名称
		ReviewStatus     int    `json:"reviewStatus,omitempty"`     // 审核状态
		Status           int    `json:"status,omitempty"`           // 展会状态
	} `json:"condition"`
}

type DxzExhibitionData struct {
	EndDate              string `json:"endDate"`              // 展会结束时间
	ExhibitionName       string `json:"exhibitionName"`       // 展会名称
	ExhibitionNumber     string `json:"exhibitionNumber"`     // 展会编号
	Id                   int64  `json:"id"`                   // 展会id
	ReviewStatus         int    `json:"reviewStatus"`         // 展会审核状态 0:草稿 1:待审核 2:审核未通过 3:审核通 过
	StartDate            string `json:"startDate"`            // 展会开始时间
	Status               int    `json:"status"`               // 展会区间状态 0:未开始 1:进行中 2:已结束
	OutlineBackgroundUrl string `json:"outlineBackgroundUrl"` // 展会封面
}

func (this *XzqExhibitionService) ExhibitionList() ([]DxzExhibitionData, error) {
	resp := make([]DxzExhibitionData, 0)

	conf := config.Conf{}
	config := conf.GetConf()

	data := new(DxzExhibitionResData)
	data.PageNo = this.Page
	data.PageSize = this.Limit
	data.Condition.ExhibitionNumber = this.ExhibitionNumber
	data.Condition.ExhibitionName = this.ExhibitionName
	data.Condition.ReviewStatus = this.ReviewStatus
	data.Condition.Status = this.Status

	dataByte, err := json.Marshal(data)
	if err != nil {
		return resp, err
	}

	_, bodyByte, err := utils.DoRequest("POST", config.DxzXzqExhibitionUrl, dataByte, "")
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
		push := DxzExhibitionData{}
		json.Unmarshal(data, &push)
		resp = append(resp, push)
	}

	return resp, nil
}
