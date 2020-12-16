package service

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"platform_report/dao"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

func PathIfNotExistMkdir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, err
	}
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func ExcelUpload(c *gin.Context, preDir string) (string, string, error) {
	pathExist, err := PathIfNotExistMkdir(preDir)
	if err != nil || !pathExist {
		return "", "", err
	}
	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	rand.Seed(time.Now().UnixNano())
	tmpName := time.Now().Format("20060102150405") + "rand_" + strconv.Itoa(rand.Intn(1000)) + RandString(5)
	newFileName := tmpName + filepath.Ext(file.Filename)
	err = c.SaveUploadedFile(file, preDir+newFileName)
	if err != nil {
		return "", "", err
	}
	return preDir + newFileName, file.Filename, nil
}

func ExcelImport(tem dao.TableExcelMap, name string) ([]string, error) {
	fieldMap := make([]map[string]string, 0)
	fieldList := make([]string, 0)
	excelFieldList := make([]string, 0)
	errList := make([]string, 0)

	if !strings.Contains(tem.FieldMap, "{") || !strings.Contains(tem.FieldMap, "}") {
		return errList, errors.New("字段映射无法解析！")
	}

	if !strings.Contains(tem.FieldMap, "[") || !strings.Contains(tem.FieldMap, "]") {
		return errList, errors.New("字段映射无法解析！")
	}

	err := json.Unmarshal([]byte(tem.FieldMap), &fieldMap)
	if err != nil {
		return errList, errors.New("字段映射无法解析！")
	}

	// 打开文件
	xlFile, err := xlsx.OpenFile(name)
	if err != nil {
		log.Println(err.Error())
		return errList, errors.New("打开文件失败！")
	}

	// 遍历 行
	for index, row := range xlFile.Sheets[0].Rows {
		data := make([]string, 0)

		column := 0
		isExistFieldMap := make(map[string]int)
		// 遍历 列
		for _, cell := range row.Cells {
			column++
			if index == 1 {
				// 第二行是字段名
				fieldName := cell.String()
				field := ""

				// 获取excel字段对应的mysql字段
				for _, v := range fieldMap {
					if _, ok := v[fieldName]; ok {
						field = v[fieldName]
					}
				}

				// 没有获取到对应的字段 或者对应的字段为空
				if field == "" {
					return errList, errors.New("第" + strconv.Itoa(column) + "列的列名不在映射表中！")
				}

				// 判断是否存在重复字段
				if _, ok := isExistFieldMap[field]; !ok {
					isExistFieldMap[field] = 1
					fieldList = append(fieldList, field)
					excelFieldList = append(excelFieldList, fieldName)
				} else {
					return errList, errors.New("不允许存在重复列！")
				}

			} else if index > 1 {
				// 第二行开始才是数据
				value := cell.String()
				data = append(data, value)
			}
		}

		if len(data) < len(excelFieldList) {
			num := len(excelFieldList) - len(data)
			for i := 0; i < num; i++ {
				data = append(data, "")
			}
		} else if len(data) > len(excelFieldList) {
			data = data[0:len(excelFieldList)]
		}

		if index > 1 {
			// 数据效验
			data, err := CheckData(fieldMap, excelFieldList, data, errList)
			if err != nil {
				log.Println("report_id:" + tem.ReportId + ",第" + strconv.Itoa(index+1) + "行数据导入失败，失败原因:" + err.Error())
				errList = append(errList, "第"+strconv.Itoa(index+1)+"行数据导入失败，失败原因:"+err.Error())
				continue
			}

			// 开始插入数据 或者等着批量插入 暂时使用单条插入
			err = new(dao.TableExcelMap).InsertList(fieldList, data, tem)
			if err != nil {
				log.Println("report_id:" + tem.ReportId + ",第" + strconv.Itoa(index+1) + "行数据导入失败，失败原因:" + err.Error())
				errList = append(errList, "第"+strconv.Itoa(index+1)+"行数据导入失败，失败原因:"+err.Error())
				continue
			}
		}

	}

	return errList, nil
}

// 数据效验
func CheckData(fieldMap []map[string]string, excelFieldList []string, data []string, errList []string) ([]string,
	error) {

	for i := 0; i < len(excelFieldList); i++ {
		codeId := ""
		// 先去获取 看看这个字段是否配置了 code_id
		for _, v := range fieldMap {
			if _, ok := v[excelFieldList[i]]; ok {
				if _, ook := v["code_id"]; ook {
					codeId = v["code_id"]
				}
			}
		}
		if codeId != "" {
			// 获取 table_excel_code 的 value
			codeList, err := new(dao.TableExcelCode).GetValueById(codeId)
			if err != nil {
				// return err
				log.Println(err.Error())
				return data, errors.New("当前" + excelFieldList[i] + "值验证失败！")
			}

			isExist := false
			// 做 数据修正 从中文转变为数据从k转变为v
			// [{"k":"显示","v":0},{"k":"隐藏","v":1}] 数据样例
			for _, v := range codeList {
				if v.K == data[i] {
					isExist = true
					data[i] = v.V
				}
			}

			if !isExist {
				return data, errors.New("当前" + excelFieldList[i] + "值不合法")
			}
		}

	}

	return data, nil
}

type TableExcelMapService struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	TableName       string `json:"table_name"`
	FieldMap        string `json:"field_map"`
	ReportId        string `json:"report_id"`
	IsImport        int    `json:"is_import"`
	IsTemplate      int    `json:"is_template"`
	TemplateAddress string `json:"template_address"`
}

// 获取
func (this *TableExcelMapService) GetInfoByNameAndReportId(name, reportId string) (dao.TableExcelMap, error) {
	tem := new(dao.TableExcelMap)
	err := tem.GetInfoByReportId(reportId)
	if err != nil {
		return *tem, err
	}
	return *tem, nil
}

func (this *TableExcelMapService) GetInfoByName(name string) (dao.TableExcelMap, error) {
	tem := new(dao.TableExcelMap)
	err := tem.GetInfoByName(name)
	if err != nil {
		return *tem, err
	}
	return *tem, nil
}

// 新增
func (this *TableExcelMapService) AddTem() error {
	tem := new(dao.TableExcelMap)
	tem.Name = this.Name
	tem.TableName = this.TableName
	tem.FieldMap = this.FieldMap
	tem.IsTemplate = 1

	fieldMap := make([]map[string]string, 0)

	if !strings.Contains(tem.FieldMap, "{") || !strings.Contains(tem.FieldMap, "}") {
		return errors.New("字段映射格式错误！")
	}

	if !strings.Contains(tem.FieldMap, "[") || !strings.Contains(tem.FieldMap, "]") {
		return errors.New("字段映射格式错误！")
	}

	err := json.Unmarshal([]byte(tem.FieldMap), &fieldMap)
	if err != nil {
		return errors.New("字段映射格式错误！")
	}
	if len(fieldMap) == 0 {
		return errors.New("字段映射不允许为空列表！")
	}

	for _, v := range fieldMap {
		for key, value := range v {
			if (key == "" || value == "") && key != "code_id" {
				return errors.New("字段映射不允许存在空字符串！")
			}
		}
	}

	temData := new(dao.TableExcelMap)
	err = temData.CheckName(this.Name, 0)
	if err != nil {
		return err
	}

	if temData.Id != 0 {
		return errors.New("文件名已存在!")
	}

	err = temData.CheckTableName(this.TableName, 0)
	if err != nil {
		return err
	}

	if temData.Id != 0 {
		return errors.New("表的配置已存在!")
	}

	_, err = tem.Add()
	return err
}

func (this *TableExcelMapService) GetList(pageStr, limitStr string) (map[string]interface{}, error) {

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

	list, err := new(dao.TableExcelMap).GetList(this.Name, this.TableName, limit, offset)
	if err != nil {
		return resp, err
	}
	count, err := new(dao.TableExcelMap).GetListCount(this.Name, this.TableName)
	if err != nil {
		return resp, err
	}

	resp["list"] = list
	resp["count"] = count
	return resp, nil
}

func (this *TableExcelMapService) Modify() error {

	if this.FieldMap == "" || this.Name == "" || this.Id == 0 {
		return errors.New("缺少参数!")
	}

	tem := new(dao.TableExcelMap)
	err := tem.GetInfoById(this.Id)
	if err != nil {
		return err
	}
	if tem.Name == "" {
		return errors.New("excel名不能为空！")
	}

	if this.Name != "" && tem.Name != this.Name {
		tem.Name = this.Name
	}
	if this.FieldMap != "" && tem.FieldMap != this.FieldMap {
		tem.FieldMap = this.FieldMap
	}

	if this.TableName != "" && tem.TableName != this.TableName {
		tem.TableName = this.TableName
	}
	if this.ReportId != "" && this.ReportId != this.ReportId {
		return errors.New("该模板已有报表配置！")
	}

	tem.ReportId = this.ReportId
	tem.IsImport = this.IsImport
	tem.IsTemplate = this.IsTemplate
	tem.TemplateAddress = this.TemplateAddress
	fieldMap := make([]map[string]string, 0)

	if !strings.Contains(tem.FieldMap, "{") || !strings.Contains(tem.FieldMap, "}") {
		return errors.New("字段映射格式错误！")
	}
	if !strings.Contains(tem.FieldMap, "[") || !strings.Contains(tem.FieldMap, "]") {
		return errors.New("字段映射格式错误！")
	}

	err = json.Unmarshal([]byte(tem.FieldMap), &fieldMap)
	if err != nil {
		return errors.New("字段映射格式错误！")
	}

	if len(fieldMap) == 0 {
		return errors.New("字段映射不允许为空列表！")
	}

	for _, v := range fieldMap {
		for key, value := range v {
			if (key == "" || value == "") && key != "code_id" {
				return errors.New("字段映射不允许存在空字符串！")
			}
		}
	}

	temData := new(dao.TableExcelMap)
	err = temData.GetInfoByReportId(this.ReportId)
	if err != nil {
		return errors.New("获取报表配置错误！")
	}
	// if this.ReportId != "" {
	// 	err = temData.CheckReportId(this.ReportId, this.Id)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	if temData.Id != 0 && temData.Id != this.Id {
	// 		return errors.New("该报表已有配置!")
	// 	}
	// }

	err = temData.CheckName(this.Name, this.Id)
	if err != nil {
		return err
	}

	if temData.Id != 0 && temData.Id != this.Id {
		return errors.New("文件名已存在!")
	}

	err = tem.Modify()
	return err
}

func (this *TableExcelMapService) DeleteTem(idStr string) error {
	if idStr == "" {
		return errors.New("id不能为空!")
	}
	id, _ := strconv.Atoi(idStr)
	tem := new(dao.TableExcelMap)
	err := tem.GetInfoById(id)
	if err != nil {
		return err
	}
	if tem.Name == "" {
		return errors.New("无数据!")
	}

	err = tem.Delete(id)
	return err
}

func (this *TableExcelMapService) GetInfo(idStr string) (dao.TableExcelMap, error) {
	tem := new(dao.TableExcelMap)

	if idStr == "" {
		return *tem, errors.New("id不能为空!")
	}
	id, _ := strconv.Atoi(idStr)
	err := tem.GetInfoById(id)

	return *tem, err
}

func (this *TableExcelMapService) GetInfoByReportId(reportIdStr string) (dao.TableExcelMap, error) {
	tem := new(dao.TableExcelMap)

	if reportIdStr == "" {
		return *tem, errors.New("id不能为空!")
	}
	err := tem.GetInfoByReportId(reportIdStr)

	return *tem, err
}

func (c *TableExcelMapService) GetReportTables() (*[]map[string]interface{}, error) {
	reportInvoke := dao.NewPlatformReport()
	tables, err := reportInvoke.GetTableListByType()
	if err != nil {
		return nil, err
	}
	ret := make([]map[string]interface{}, 0, 10)
	for _, val := range tables {
		columns, err := reportInvoke.GetTableColumn(val["table_name"])
		if err != nil {
			return nil, err
		}
		item := make(map[string]interface{})
		item["table_name"] = val["table_name"]
		item["columns"] = columns
		ret = append(ret, item)
	}
	return &ret, nil
}

func (c *TableExcelMapService) ExcelTemplateImport(id int, address string) error {
	if id == 0 {
		return errors.New("参数不合法！")
	}

	fieldMap := make([]map[string]string, 0)
	tem := new(dao.TableExcelMap)
	err := tem.GetInfoById(id)
	if err != nil {
		log.Println(err.Error())
		return errors.New("获取配置信息失败！")
	}

	err = json.Unmarshal([]byte(tem.FieldMap), &fieldMap)
	if err != nil {
		return errors.New("字段映射无法解析！")
	}

	// 打开文件
	xlFile, err := xlsx.OpenFile(address)
	if err != nil {
		log.Println(err.Error())
		return errors.New("打开文件失败！")
	}

	// 遍历 行
	for index, row := range xlFile.Sheets[0].Rows {
		column := 0
		isExistFieldMap := make(map[string]int)
		// 遍历 列
		for _, cell := range row.Cells {
			column++
			if index == 1 {
				// 第二行是字段名
				fieldName := cell.String()
				field := ""
				if fieldName != "code_id" {
					for _, v := range fieldMap {
						if _, ok := v[fieldName]; ok {
							field = v[fieldName]
						}
					}
				}

				if field == "" {
					return errors.New("第" + strconv.Itoa(column) + "列的列名不在映射表中！")
				}

				// 判断是否存在重复字段
				if _, ok := isExistFieldMap[field]; !ok {
					isExistFieldMap[field] = 1
				} else {
					return errors.New("不允许存在重复列！")
				}

			}
		}
	}

	tem.TemplateAddress = address
	err = tem.Modify()

	return err
}
