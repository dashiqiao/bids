package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"platform_report/dao"
	"strings"
	"time"
)

type FieldService struct {
}

func NewFieldService() *FieldService {
	return new(FieldService)
}

func (c *FieldService) DeleteActions(uuid string) {
	defines := c.GetActionsById(uuid)
	_ = dao.NewReportFieldInvoke().Delete(defines.ActionId)
	kvcache.Flush()
}

func (c *FieldService) GetActions() (ret []*dao.ReportFieldDefines, err error) {
	key := "actions"
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportFieldDefines), nil
	}
	items, err := dao.NewReportFieldInvoke().GetActions()
	if err != nil {
		return nil, err
	}
	kvcache.Set(key, items, cache.DefaultExpiration)
	return items, nil
}

func (c *FieldService) GetActionsById(uuid string) *dao.ReportFieldDefines {
	items, _ := c.GetActions()
	for _, item := range items {
		if item.Uuid == uuid {
			return item
		}
	}
	return nil
}

func (c *FieldService) SaveFields(uu, name string, fields []dao.ReportField) error {
	invoke := dao.NewReportFieldInvoke()
	actionId := 0
	if uu == "" {
		uu = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
		actionId = invoke.CreateActions(uu, name)
		if actionId == 0 {
			return errors.New("添加失败")
		}
	} else {
		item := c.GetActionsById(uu)
		if item == nil || item.ActionId == 0 {
			return errors.New("动作获取失败")
		}
		actionId = item.ActionId
		_ = invoke.UpdActions(uu, name)
	}

	for i := 0; i < len(fields); i++ {
		fields[i].ActionId = actionId
		fields[i].OrderId = i
		fields[i].CreateTime = time.Now()
	}
	err := dao.NewReportFieldInvoke().Save(actionId, fields)
	if err != nil {
		return err
	}
	NewReportService().CC()
	return nil
}

func (c *FieldService) GetFields(reportId int) []*dao.ReportField {
	key := fmt.Sprintf("fields-%v", reportId)
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportField)
	}
	datas, err := dao.NewReportFieldInvoke().GetFieldsByReportId(reportId)
	if err != nil {
		return nil
	}

	for i := 0; i < len(datas); i++ {
		if datas[i].Setting == 0 {
			continue
		}
		condition := NewReportService().GetReportConditionsById(datas[i].Setting)
		datas[i].ParentId = condition.ParentId
		datas[i].Value = condition.CondVal
		if condition.Conduct == 2 {
			mapJ, _ := dao.NewInfiniteStones().War(condition.CondVal)
			mjson, _ := json.Marshal(mapJ)
			datas[i].Value = string(mjson)
		}
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *FieldService) GetPrimaryKey(reportId int) []string {
	fields := c.GetFields(reportId)
	items := make([]string, 0)
	for _, field := range fields {
		if field.IsPrimary == 1 {
			items = append(items, field.Field)
		}
	}
	return items
}

func (c *FieldService) GetValues(reportId int, arg map[string]interface{}) map[string]interface{} {
	fields := c.GetFields(reportId)
	valueList := make([]interface{}, 0)
	sql := "SELECT * FROM " + fields[0].Types + " WHERE "
	i := 1
	for key, val := range arg {
		if i == 1 {
			sql += key + " = ? "
		} else {
			sql += " AND " + key + " = ? "
		}
		valueList = append(valueList, val)
		i++
	}
	maps, _ := dao.NewInfiniteStones().WarWithParametes(sql, valueList)
	return maps
}
