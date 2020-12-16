package service

import (
	"encoding/json"
	"github.com/panjf2000/ants/v2"
	"platform_report/dao"
	"time"
)

type LogService struct {
}

func NewLogService() *LogService {
	return new(LogService)
}

func (c *LogService) AddReportSubmitLog(uuid, method, uri, data string, oprateType, ip, userId int, userName string) {
	_ = ants.Submit(func() {
		log := new(dao.ReportLog)
		log.Uuid = uuid
		log.Method = method
		log.Uri = uri
		log.PostData = data
		log.OprateType = oprateType
		log.Ip = ip
		log.UserId = userId
		log.UserName = userName
		log.InTime = time.Now()
		log.AddReportSubmitLog()
	})
}

func (c *LogService) AddReportActionLog(reportId string, types, method string, before, after map[string]interface{}, userId int, userName string) {
	log := new(dao.ReportActionLog)
	log.ReportId = reportId
	log.Types = types
	log.Method = method
	beforeJson, _ := json.Marshal(before)
	log.Before = string(beforeJson)
	afterJson, _ := json.Marshal(after)
	log.After = string(afterJson)
	log.UserId = userId
	log.UserName = userName
	log.InTime = time.Now()
	log.AddReportSubmitLog()
}
