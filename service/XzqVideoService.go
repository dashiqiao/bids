package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"platform_report/config"
	"platform_report/utils"
)

type XzqVideoService struct{}

type VideoData struct {
	Flu string `json:"flu"` // 清晰度,(flu,流畅)
	Hd  string `json:"hd"`  // 清晰度,(hd ,高清)
	Sd  string `json:"sd"`  // 清晰度,(sd ,标清)
}

func (this *XzqVideoService) MultiPlay(fileId string, videoType int) (VideoData, error) {
	resp := VideoData{}

	conf := config.Conf{}
	config := conf.GetConf()

	url := config.DxzXzqMultiPlayUrl + fmt.Sprintf("?fileId=%v", fileId)
	fmt.Println(url)
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

	data, _ := json.Marshal(dxzRespData.Result.(map[string]interface{}))
	json.Unmarshal(data, &resp)

	return resp, nil
}
