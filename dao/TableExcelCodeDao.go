package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"platform_report/lib"
	"strconv"
	"time"
)

type TableExcelCode struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`    // 字段名
	Type    string `json:"type"`    // select array
	Value   string `json:"value"`   // 值
	Tips    string `json:"tips"`    // 提示
	Conduct int    `json:"conduct"` // 1 2
}

func (this *TableExcelCode) Add() (int64, error) {
	id, err := lib.InitXormMySql().Insert(this)
	return id, err
}

func (this *TableExcelCode) GetInfoById(id int) error {
	_, err := lib.InitXormMySql().Table("table_excel_code").Where("id = ?", id).Get(this)
	return err
}

func (this *TableExcelCode) GetListByName(name, Codetype string, limit, page int) ([]TableExcelCode, error) {
	list := make([]TableExcelCode, 0)
	query := lib.InitXormMySql().Table("table_excel_code")
	if name != "" {
		query.Where("name = ? ", name)
	}
	if Codetype != "" {
		query.Where("type = ? ", Codetype)
	}

	err := query.Limit(limit, page).Find(&list)

	return list, err
}

func (this *TableExcelCode) GetList(limit, page int) ([]TableExcelCode, error) {
	list := make([]TableExcelCode, 0)
	query := lib.InitXormMySql().Table("table_excel_code")
	if this.Name != "" {
		query = query.Where("name = ? ", this.Name)
	}
	if this.Type != "" {
		query = query.Where("type = ? ", this.Type)
	}

	if limit != 0 && page != 0 {
		query = query.Limit(limit, page)
	}
	err := query.Find(&list)

	return list, err
}

func (this *TableExcelCode) GetListCount(name, codeType string) (int64, error) {
	query := lib.InitXormMySql().Table("table_excel_code")
	if name != "" {
		query = query.Where("name = ? ", name)
	}
	if codeType != "" {
		query = query.Where("type = ? ", codeType)
	}

	num, err := query.Count()

	return num, err
}

func (this *TableExcelCode) Modify() error {
	_, err := lib.InitXormMySql().Table("table_excel_code").Where("id = ? ", this.Id).Update(this)
	return err
}

func (this *TableExcelCode) Delete(id int) error {
	_, err := lib.InitXormMySql().Table("table_excel_code").Where("id = ? ", id).Delete(this)
	return err
}

func (this *TableExcelCode) GetDataBySql(sql string) ([]map[string]string, error) {
	list, err := lib.InitXormMySql().QueryString(sql)
	return list, err
}

const EXCEL_CODE_DATA_KEY = "bids_excel_code"

func (this *TableExcelCode) ReloadExcelCodeData() error {
	list, err := this.GetList(0, 0)
	if err != nil {
		log.Println(err.Error())
		return errors.New("重新加载失败！")
	}

	rcl := lib.GetRedis()
	for _, v := range list {
		key := fmt.Sprintf("%s_%v", EXCEL_CODE_DATA_KEY, v.Id)
		rcl.Del(key).Result()

		data := v.Value
		if v.Type == "select" {
			list, err := this.GetDataBySql(v.Value)
			if err != nil {
				log.Println(err.Error())
				return errors.New("重新加载失败！")
			}

			dataByte, _ := json.Marshal(list)
			data = string(dataByte)
		}

		rcl.Set(key, data, 1*time.Hour)
	}
	return nil
}

type CodeValue struct {
	K string `json:"k"`
	V string `json:"v"`
}

func (this *TableExcelCode) GetValueById(idStr string) ([]CodeValue, error) {
	rcl := lib.GetRedis()
	key := fmt.Sprintf("%s_%v", EXCEL_CODE_DATA_KEY, idStr)
	list := make([]CodeValue, 0)

	dataStr, err := rcl.Get(key).Result()
	if err != nil || dataStr == "" {
		id, _ := strconv.Atoi(idStr)
		info := TableExcelCode{}
		_, err = lib.InitXormMySql().Table("table_excel_code").Where("id = ?", id).Get(&info)
		if err != nil {
			return list, err
		}

		switch info.Type {
		case "array":
			dataStr = info.Value
		case "select":
			listMap, err := this.GetDataBySql(info.Value)
			if err != nil {
				return list, err
			}
			if len(listMap) > 0 {
				dataByte, _ := json.Marshal(listMap)
				dataStr = string(dataByte)
			} else {
				dataStr = "[]"
			}
		default:
			return list, errors.New("参数不合法！")
		}

		rcl.Set(key, dataStr, 1*time.Hour)
	}
	jerr := json.Unmarshal([]byte(dataStr), &list)
	log.Println(list)
	return list, jerr
}
