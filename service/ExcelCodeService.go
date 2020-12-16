package service

import (
	"encoding/json"
	"errors"
	"log"
	"platform_report/dao"
	"strconv"
	"strings"
)

type TableExcelCodeService struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`  // 字段名
	Type  string `json:"type"`  // select array
	Value string `json:"value"` // 值
	Tips  string `json:"tips"`  // 提示
	// Conduct int    `json:"conduct"` // 1 2
}

// 新增
func (this *TableExcelCodeService) AddExcelCode() error {
	if this.Name == "" || this.Value == "" || this.Type == "" {
		return errors.New("缺少参数!")
	}

	tec := new(dao.TableExcelCode)
	tec.Name = this.Name
	tec.Type = this.Type
	tec.Value = this.Value
	tec.Tips = this.Tips
	tec.Conduct = GetConduct(this.Type)
	if tec.Conduct == 0 {
		return errors.New("类型错误！")
	}

	ok := CheckTypeAndValue(this.Type, this.Value)
	if !ok {
		return errors.New("类型与值不匹配！")
	}

	if tec.Type == "select" {
		list, err := new(dao.TableExcelCode).GetDataBySql(tec.Value)
		if err != nil {
			log.Println(err.Error())
			return errors.New("重新加载失败！")
		}

		dataByte, _ := json.Marshal(list)

		checkList := make([]dao.CodeValue, 0)
		err = json.Unmarshal(dataByte, &checkList)
		if err != nil {
			return errors.New("数据类型错误！")
		}

	}

	_, err := tec.Add()
	return err
}

func (this *TableExcelCodeService) ModifyExcelCode() error {

	if this.Name == "" || this.Id == 0 || this.Value == "" || this.Type == "" {
		return errors.New("缺少参数!")
	}

	tec := new(dao.TableExcelCode)
	err := tec.GetInfoById(this.Id)
	if err != nil {
		return err
	}
	if tec.Name != this.Name {
		tec.Name = this.Name
	}
	if tec.Tips != this.Tips {
		tec.Tips = this.Tips
	}

	tec.Conduct = GetConduct(this.Type)
	if tec.Conduct == 0 {
		return errors.New("类型错误！")
	}

	if tec.Type != this.Type {
		tec.Type = this.Type
	}
	if tec.Value != this.Value {
		ok := CheckTypeAndValue(this.Type, this.Value)
		if !ok {
			return errors.New("类型与值不匹配！")
		}
		tec.Value = this.Value
	}

	err = tec.Modify()

	return nil
}

func (this *TableExcelCodeService) DeleteExcelCode(idStr string) error {
	if idStr == "" {
		return errors.New("id不能为空!")
	}
	id, _ := strconv.Atoi(idStr)
	tec := new(dao.TableExcelCode)
	err := tec.GetInfoById(id)
	if err != nil {
		return err
	}
	if tec.Name == "" {
		return errors.New("无数据!")
	}

	err = tec.Delete(id)
	return err
}

func (this *TableExcelCodeService) GetExcelCodeInfo(idStr string) (dao.TableExcelCode, error) {
	tec := new(dao.TableExcelCode)

	if idStr == "" {
		return *tec, errors.New("id不能为空!")
	}
	id, _ := strconv.Atoi(idStr)
	err := tec.GetInfoById(id)

	return *tec, err
}

func (this *TableExcelCodeService) GetExcelCodeList(pageStr, limitStr string) (map[string]interface{}, error) {

	resp := make(map[string]interface{}, 0)

	page := 1
	limit := 15

	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			return resp, errors.New("page不是数字!")
		}
		if pageInt > 0 {
			page = pageInt
		}
	}

	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			return resp, errors.New("limit不是数字!")
		}
		if limitInt > 0 {
			limit = limitInt
		}
	}
	offset := (page - 1) * limit

	list, err := new(dao.TableExcelCode).GetListByName(this.Name, this.Type, limit, offset)
	if err != nil {
		return resp, err
	}
	count, err := new(dao.TableExcelCode).GetListCount(this.Name, this.Type)
	if err != nil {
		return resp, err
	}

	resp["list"] = list
	resp["count"] = count
	return resp, nil
}

func (this *TableExcelCodeService) ReloadExcelCodeList() error {
	err := new(dao.TableExcelCode).ReloadExcelCodeData()
	return err
}

func CheckTypeAndValue(codeType, value string) bool {
	switch codeType {
	case "select":
		if !strings.Contains(value, "select") && !strings.Contains(value, "SELECT") {
			return false
		}
		if !strings.Contains(value, "from") && !strings.Contains(value, "FROM") {
			return false
		}
		return true
	case "array":
		// [{"k":"显示","v":0},{"k":"隐藏","v":1}]
		list := make([]dao.CodeValue, 0)
		err := json.Unmarshal([]byte(value), &list)
		if err != nil {
			return false
		}
		if len(list) == 0 {
			return false
		}
		for _, v := range list {
			if v.K == "" || v.V == "" {
				return false
			}
		}
		return true
	}

	return false
}
func GetConduct(codeType string) int {
	switch codeType {
	case "select":
		return 1
	case "array":
		return 2
	}

	return 0
}
