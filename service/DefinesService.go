package service

import (
	"fmt"
	"platform_report/dao"
	"platform_report/lib"
	"platform_report/pkg/util"
	"strconv"
)

type DefinesService struct {
	prefix string
}

func NewDefinesService() *DefinesService {
	return &DefinesService{prefix: "report"}
}

func (c *DefinesService) GetKey(uuid string) string {
	return c.prefix + ":" + uuid
}

func (c *DefinesService) RemoveCache(uuid string) {
	_, err := lib.GetRedis().Del(c.GetKey(uuid)).Result()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *DefinesService) GetOne(uuid string) *dao.ReportDefines {
	resMap := lib.GetRedis().HGetAll(c.GetKey(uuid)).Val()
	if len(resMap) > 0 {
		return c.Fill(resMap)
	}
	items, _ := new(dao.DefinesDao).GetMulti([]string{uuid})
	if len(items) > 0 {
		lib.GetRedis().HMSet(c.GetKey(uuid), util.StructToMap(items[0]))
		return items[0]
	}
	return nil
}

func (c *DefinesService) Get(name string, p, limit int) ([]*dao.ReportDefines, int64) {
	invoke := new(dao.DefinesDao)
	ids, count, _ := invoke.Get(name, p, limit)
	retMap := make([]*dao.ReportDefines, 0, limit)
	notExist := make([]string, 0)
	for _, id := range ids {
		resMap := lib.GetRedis().HGetAll(c.GetKey(id)).Val()
		if len(resMap) > 0 {
			retMap = append(retMap, c.Fill(resMap))
		} else {
			notExist = append(notExist, id)
		}
	}
	if len(notExist) > 0 {
		notExistRes, _ := invoke.GetMulti(notExist)
		if notExistRes != nil {
			retMap = append(retMap, notExistRes...)
			for _, res := range notExistRes {
				lib.GetRedis().HMSet(c.GetKey(res.Uuid), util.StructToMap(res))
			}
		}
	}
	return retMap, count
}

func (c *DefinesService) Fill(resMap map[string]string) *dao.ReportDefines {
	if len(resMap) > 0 {
		item := new(dao.ReportDefines)
		item.ReportId, _ = strconv.Atoi(resMap["report_id"])
		item.Uuid = resMap["uuid"]
		item.ReportName = resMap["report_name"]
		item.Sql = resMap["sql"]
		item.SqlCount = resMap["sql_count"]
		item.SqlSummary = resMap["sql_summary"]
		item.ReportType, _ = strconv.Atoi(resMap["report_type"])
		item.ReportAction, _ = strconv.Atoi(resMap["report_action"])
		item.ReportAuto, _ = strconv.Atoi(resMap["report_auto"])
		item.Tip = resMap["tip"]
		item.Note = resMap["note"]
		item.Remark = resMap["remark"]
		item.DataSource, _ = strconv.Atoi(resMap["data_source"])
		item.ReportAttribute, _ = strconv.Atoi(resMap["report_attribute"])
		return item
	}
	return nil
}
