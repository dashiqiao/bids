package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"platform_report/service"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	return
}

/**
* @api {POST} /api/v1/tool/excelImport excel导入
* @apiName ExcelImport
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} report_id 报表id
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/
func ExcelImport(c *gin.Context) {

	// 上传文件
	name, fileName, err := service.ExcelUpload(c, "upload/")
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	reportId := c.PostForm("report_id")

	fileNameList := strings.Split(fileName, ".")
	temSvc := new(service.TableExcelMapService)
	tem, err := temSvc.GetInfoByNameAndReportId(fileNameList[0], reportId)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}
	if tem.Id == 0 {
		Response(c, 400, "", "查找不到对应的excel相关配置!")
		return
	}

	errList, err := service.ExcelImport(tem, name)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	if len(errList) != 0 {
		Response(c, 401, errList, "导入存在失败数据！")
		return
	}

	Response(c, 200, "", "ok!")
	return
}

/**
* @api {POST} /api/v1/tool/tem 新增表与excel映射
* @apiName AddTem
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {String} name 文件名
* @apiParam {String} table_name 表名
* @apiParam {String} field_map 映射关系
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/

func AddTem(c *gin.Context) {
	temSvc := new(service.TableExcelMapService)
	dataByte, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(dataByte, temSvc)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}
	if temSvc.Name == "" || temSvc.TableName == "" || temSvc.FieldMap == "" {
		Response(c, 400, "", "缺少参数!")
		return
	}

	err = temSvc.AddTem()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return
}

/**
* @api {GET} /api/v1/tool/tem/list 获取映射列表
* @apiName GetTemList
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} page 页码
* @apiParam {int} limit 数量
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": {
        "count": 1,
        "list": [
            {
                "id": 1,
                "name": "数据",
                "table_name": "table_excel_map",
                "field_map": "{\"中文名\":\"name\",\"表名\":\"table_name\",\"字段映射\":\"field_map\",\"状态\":\"status\"}",
                "status": 1
            }
        ]
    },
    "msg": "ok!"
}
*/
func GetTemList(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")
	name := c.Query("name")
	tableName := c.Query("table_name")

	temSvc := new(service.TableExcelMapService)
	temSvc.Name = name
	temSvc.TableName = tableName

	list, err := temSvc.GetList(page, limit)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "ok!")
	return
}

/**
* @api {PUT} /api/v1/tool/tem 修改表与excel映射
* @apiName ModifyTem
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id int
* @apiParam {String} name 文件名
* @apiParam {int} status 状态
* @apiParam {String} field_map 映射关系
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/
func ModifyTem(c *gin.Context) {
	temSvc := new(service.TableExcelMapService)
	err := c.ShouldBind(temSvc)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}
	err = temSvc.Modify()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return

}

/**
* @api {DELETE} /api/v1/tool/tem 删除表与excel映射
* @apiName DelTem
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id int
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": "",
    "msg": "ok!"
}
*/
func DelTem(c *gin.Context) {
	id := c.Query("id")
	temSvc := new(service.TableExcelMapService)
	err := temSvc.DeleteTem(id)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return
}

/**
* @api {GET} /api/v1/tool/tem?id=xxx 获取单表的配置
* @apiName GetTem
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id int
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": {
                "id": 1,
                "name": "数据",
                "table_name": "table_excel_map",
                "field_map": "{\"中文名\":\"name\",\"表名\":\"table_name\",\"字段映射\":\"field_map\",\"状态\":\"status\"}",
                "status": 1
            },
    "msg": "ok!"
}
*/
func GetTem(c *gin.Context) {
	id := c.Query("id")
	temSvc := new(service.TableExcelMapService)
	tem, err := temSvc.GetInfo(id)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, tem, "ok!")
	return
}

/**
* @api {GET} /api/v1/tool/export/template 导入模板
* @apiName ExportTemplate
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id int
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     返回excel流
 */
func ExportTemplate(c *gin.Context) {
	id := c.Query("id")
	temSvc := new(service.TableExcelMapService)
	tem, err := temSvc.GetInfo(id)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}
	if tem.TemplateAddress != "" {

		pwd, _ := os.Getwd()
		fileContent, err := ioutil.ReadFile(pwd + "/" + tem.TemplateAddress)
		if err != nil {
			Response(c, 400, "", err.Error())
			return
		}
		// 设置返回头并返回数据
		c.Header("Content-type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment;filename="+tem.Name+".xlsx")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Data(http.StatusOK, "application/octet-stream", fileContent)
		return
	}

	data := make([]map[string]string, 0)
	err = json.Unmarshal([]byte(tem.FieldMap), &data)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	c.Header("Content-type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment;filename="+tem.Name+".xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	excel := excelize.NewFile()
	sheetName := "Sheet1"
	sheet := excel.NewSheet(sheetName)
	// style, _ := excel.NewStyle(`{"font":{"bold":true},"fill":{"type":"pattern","color":["#FFFF00"] ,"pattern":1}}`)
	excel.MergeCell(sheetName, "A1", Div(len(data))+"1")
	// log.Println(len(data))
	// log.Println(Div(len(data)))

	// 	style, err := f.NewStyle(`{"alignment":{"horizontal":"center","ident":1,"justify_last_line":true,"reading_order":0,"relative_indent":1,"shrink_to_fit":true,"text_rotation":45,"vertical":"","wrap_text":true}}`)
	// if err != nil {
	//     fmt.Println(err)
	// }
	// err = f.SetCellStyle("Sheet1", "H9", "H9", style)
	style, _ := excel.NewStyle(`{"alignment":{"wrap_text":true}}`)
	excel.SetCellStyle(sheetName, "A1", "A1", style)
	excel.SetRowHeight(sheetName, 1, 150)
	// SetColWidth
	// err := f.SetColWidth("Sheet1", "A", "H", 20)

	content := "说明:\n1.已配置字典的字段,请自行匹配设置。\n2.字段位置可以任意修改,但字段名不可修改,不可新增其它字段名。\n3.如果删除字段,请修改第一行的单元格合并数量,保证合并的数量与列的数量一致。\n4.如果日期格式为datetime," +
		"请把格式设置为字符串。\n5.该行不可删除。"
	excel.SetCellValue(sheetName, "A1", content)

	clumnNum := 0
	for _, v := range data {
		for k, _ := range v {
			if k != "code_id" {
				clumnNum++
				sheetPosition := Div(clumnNum) + "2"
				excel.SetCellValue(sheetName, sheetPosition, k)
			}
		}
	}

	// 合并单元格
	// err := f.MergeCell("Sheet1", "D3", "E9")
	// func (f *File) MergeCell(sheet, hcell, vcell string) error
	excel.SetActiveSheet(sheet)

	if err := excel.Write(c.Writer); err != nil {
		Response(c, 400, "", err.Error())
		return
	}

}

/**
* @api {GET} /api/v1/tool/tem?report_id=xxx 获取单表的配置
* @apiName GetTem
* @apiGroup tool 工具
* @apiHeader {String} Authorization 用户授权token
* @apiParam {int} id int
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
    "code": 200,
    "data": {
                "id": 1,
                "name": "数据",
                "table_name": "table_excel_map",
                "field_map": "{\"中文名\":\"name\",\"表名\":\"table_name\",\"字段映射\":\"field_map\",\"状态\":\"status\"}",
                "status": 1
            },
    "msg": "ok!"
}
*/
func GetTemButton(c *gin.Context) {
	reportId := c.Query("report_id")
	temSvc := new(service.TableExcelMapService)
	tem, err := temSvc.GetInfoByReportId(reportId)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, tem, "ok!")
	return
}

/**
* @api {Get} /admin/v1/tool/excel/tablelist 获取数据表结构
* @apiName api.TableList
* @apiGroup admin 管理
* @apiHeader {String} Authorization 用户授权token
*
* @apiSuccessExample Success-TXResponse:
*     HTTP/1.1 200 OK
*     {
		"code": 200,
		"msg": "ok"
		"data":
			{
				"columns": [
					{
						"column_name": "id"  --列名
					}],
				"table_name": "dxz_course_languages_statistics_info" --表名
			}
      }
*/
func ExcelTableList(c *gin.Context) {
	invoke := new(service.TableExcelMapService)
	ret, err := invoke.GetReportTables()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": ret,
	})
}

// 上传模板
func ExcelTemplateImport(c *gin.Context) {

	// 上传文件
	address, _, err := service.ExcelUpload(c, "upload/excelTemplate/")
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	idStr := c.PostForm("id")
	id, _ := strconv.Atoi(idStr)
	err = new(service.TableExcelMapService).ExcelTemplateImport(id, address)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "ok!")
	return
}
