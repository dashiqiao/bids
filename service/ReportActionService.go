package service

import (
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	. "platform_report/config"
	"platform_report/dao"
	"strings"
	"time"
)

type ActionParameter struct {
	Fields []*dao.ReportField
	Values map[string]interface{}
}

type Return struct {
	Code int
	Msg  string
}

type ReportAction struct {
	Server string
	Host   string
}

func NewReportAction() *ReportAction {
	cf := Conf{}
	config := cf.GetConf()
	return &ReportAction{Host: config.ReportRpc, Server: "ActionInvoke"}
}

func (c *ReportAction) Do(uuid, method string, values map[string]interface{}) *Return {
	//d := client.NewPeer2PeerDiscovery("tcp@"+c.Host, "")
	//xclient := client.NewXClient(c.Server, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	//
	//defer xclient.Close()
	invoke := NewFieldService()
	ret := &Return{}

	args := ActionParameter{
		Values: values,
		Fields: invoke.GetFields(NewReportService().GetReportIdByUuid(uuid)),
	}
	//_ = xclient.Call(context.Background(), method, args, ret)
	ai := new(ActionInvoke)
	switch method {
	case "Insert":
		_ = ai.Insert(args, ret)
	case "Update":
		_ = ai.Update(args, ret)
	case "Delete":
		_ = ai.Delete(args, ret)
	}
	return ret
}

type ActionInvoke struct{}

func (c *ActionInvoke) Insert(paras ActionParameter, reply *Return) error {
	maps := make(map[string]interface{})
	uniqueKey := make([]string, 0)
	uniqueValue := make([]interface{}, 0)
	for _, field := range paras.Fields {
		if field.IsPrimary != 1 && field.IsUnique == 1 { //不是主键还唯一
			if value, ok := paras.Values[field.Field]; ok {
				if strings.ToLower(value.(string)) != "null" {
					uniqueKey = append(uniqueKey, field.Field)
					uniqueValue = append(uniqueValue, value)
				}
			}
		}

		if field.IsPrimary != 1 {
			if value, ok := paras.Values[field.Field]; ok {
				maps[field.Field] = value
				if strings.ToLower(value.(string)) == "null" {
					maps[field.Field] = sql.NullString{}
				}
			} else {
				if strings.ToLower(field.DefaultValue) == "now" {
					maps[field.Field] = time.Now()
				} else if strings.ToLower(field.DefaultValue) == "uuid" {
					maps[field.Field] = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
				} else if strings.ToLower(field.DefaultValue) == "null" {
					maps[field.Field] = sql.NullString{}
				} else {
					maps[field.Field] = field.DefaultValue
				}
			}
		}
	}

	if len(uniqueKey) > 0 {
		if dao.Count(paras.Fields[0].Types, uniqueKey, uniqueValue) > 0 {
			reply.Code = 502
			reply.Msg = "唯一项重复,请重新输入！"
			return nil
		}
	}

	sql, args := dao.MapToInsertSql(paras.Fields[0].Types, maps)
	err := dao.Exec(sql, args)
	if err != nil {
		reply.Code = 501
		reply.Msg = err.Error()
	}
	return nil
}

func (c *ActionInvoke) Update(paras ActionParameter, reply *Return) error {
	primaryKeys := GetPrimaryKey(paras.Fields)
	primaryKeyValue := make([]interface{}, 0)
	uniqueKey := make([]string, 0)
	uniqueValue := make([]interface{}, 0)
	for _, field := range paras.Fields {
		if field.IsPrimary != 1 && field.IsUnique == 1 { //不是主键还唯一
			if value, ok := paras.Values[field.Field]; ok {
				if strings.ToLower(value.(string)) != "null" {
					uniqueKey = append(uniqueKey, field.Field)
					uniqueValue = append(uniqueValue, value)
				}
			}
		}
	}

	for _, primaryKey := range primaryKeys {
		if value, ok := paras.Values[primaryKey]; ok {
			primaryKeyValue = append(primaryKeyValue, value)
		}
	}
	if len(primaryKeys) == 0 {
		reply.Code = 500
		reply.Msg = "主键参数未找到"
		return nil
	}

	if len(uniqueKey) > 0 {
		if dao.Count(paras.Fields[0].Types, uniqueKey, uniqueValue) > 1 {
			reply.Code = 502
			reply.Msg = "唯一项重复,请重新输入！"
			return nil
		}
	}

	maps := make(map[string]interface{})
	for _, field := range paras.Fields {
		if field.IsPrimary == 1 || field.IsShow == 0 {
			continue
		}
		if value, ok := paras.Values[field.Field]; ok {
			maps[field.Field] = value
			if strings.ToLower(value.(string)) == "null" {
				maps[field.Field] = sql.NullString{}
			}
		} else {
			if strings.ToLower(field.DefaultValue) == "now" {
				maps[field.Field] = time.Now()
			} else if strings.ToLower(field.DefaultValue) == "uuid" {
				maps[field.Field] = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
			} else {
				maps[field.Field] = field.DefaultValue
			}
		}

	}
	sql, args := dao.MapToUpdateSql(paras.Fields[0].Types, maps, primaryKeys, primaryKeyValue)
	err := dao.Exec(sql, args)
	if err != nil {
		reply.Code = 501
		reply.Msg = err.Error()
	}
	return nil
}

func (c *ActionInvoke) Delete(paras ActionParameter, reply *Return) error {
	primaryKeys := GetPrimaryKey(paras.Fields)
	primaryKeyValue := make([] interface{}, 0)
	for _, primaryKey := range primaryKeys {
		if value, ok := paras.Values[primaryKey]; ok {
			primaryKeyValue = append(primaryKeyValue, value)
		}
	}
	if len(primaryKeys) == 0 {
		reply.Code = 500
		reply.Msg = "主键参数未找到"
		return nil
	}
	args := make([]interface{}, 0)
	//sql := fmt.Sprintf("DELETE FROM %v WHERE %v=?", paras.Fields[0].Types, primaryKeys[0])
	sql := fmt.Sprintf("UPDATE %v SET deleted_at = NOW()  WHERE %v=?", paras.Fields[0].Types, primaryKeys[0])
	args = append(args, primaryKeyValue[0])
	if len(primaryKeys) > 1 {
		for i := 1; i < len(primaryKeys); i++ {
			sql += fmt.Sprintf(" AND %v = ?", primaryKeys[i])
			args = append(args, primaryKeyValue[i])
		}
	}
	err := dao.Exec(sql, args)
	if err != nil {
		reply.Code = 501
		reply.Msg = err.Error()
	}
	return nil
}

func GetPrimaryKey(fields []*dao.ReportField) []string {
	//primaryKey := ""
	items := make([]string, 0)
	for _, field := range fields {
		if field.IsPrimary == 1 {
			items = append(items, field.Field)
		}
	}
	return items
}
