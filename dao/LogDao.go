package dao

import (
	"platform_report/lib"
	"time"
)

type ReportLog struct {
	Uuid       string
	Method     string
	Uri        string
	PostData   string
	OprateType int
	UserId     int
	UserName   string
	Ip         int
	InTime     time.Time
}

func (c *ReportLog) AddReportSubmitLog() {
	_, _ = lib.InitXormMySql().Insert(c)
}

type ReportActionLog struct {
	ReportId string
	Types    string
	Method   string
	Before   string
	After    string
	UserId   int
	UserName string
	InTime   time.Time
}

func (c *ReportActionLog) AddReportSubmitLog() {
	_, _ = lib.InitXormMySql().Insert(c)
}
