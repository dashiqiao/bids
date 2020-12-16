package service

import (
	"encoding/json"
	"errors"
	"log"

	"platform_report/config"
	"platform_report/utils"
)

type XzqCourseService struct {
	Page           int    `json:"page"`           // 页码
	Limit          int    `json:"limit"`          // 数量
	GoodsId        string `json:"goodsId"`        // 商品Id
	ProductName    string `json:"productName"`    // 课程标题
	TeacherName    string `json:"teacherName"`    // 讲师姓名
	CategoryCode   int    `json:"categoryCode"`   // 上课模式
	SalesStatus    int    `json:"salesStatus"`    // 状态
	AuthorType     int    `json:"authorType"`     // 角色类型
	ChannelCode    int    `json:"channelCode"`    // 频道
	FirstLabel     string `json:"firstLabel"`     // 一级标签code
	SecondLabel    string `json:"secondLabel"`    // 二级标签code
	Tag            int    `json:"tag"`            // 主题ID
	UpdatedAtStart string `json:"updatedAtStart"` // 更新日期检索开始日期 yyyy-MM-dd
	UpdatedAtEnd   string `json:"updatedAtEnd"`   // 更新日期检索结束日期 yyyy-MM-dd
}

// categoryCode 为空的时候 查询的是直播课
// categoryCode = 100010005 是录播课
type DxzCourseResData struct {
	PageNo    int `json:"pageNo"`
	PageSize  int `json:"pageSize"`
	Condition struct {
		GoodsId        string `json:"goodsId,omitempty"`
		ProductName    string `json:"productName,omitempty"`
		TeacherName    string `json:"teacherName,omitempty"`
		CategoryCode   int    `json:"categoryCode,omitempty"`
		SalesStatus    int    `json:"salesStatus,omitempty"`
		AuthorType     int    `json:"authorType,omitempty"`
		ChannelCode    int    `json:"channelCode,omitempty"`
		FirstLabel     string `json:"firstLabel,omitempty"`
		SecondLabel    string `json:"secondLabel,omitempty"`
		Tag            int    `json:"tag,omitempty"`
		UpdatedAtStart string `json:"updatedAtStart,omitempty"`
		UpdatedAtEnd   string `json:"updatedAtEnd,omitempty"`
	} `json:"condition"`
}

type DxzResponseData struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	MessageArgs []string    `json:"messageArgs"`
	Result      interface{} `json:"result"`
	Success     bool        `json:"success"`
	PageInfo    interface{} `json:"pageInfo"`
}

type DxzCourseData struct {
	Id                int64    `json:"id"`                // 课程id
	AuthorName        string   `json:"authorName"`        // 角色类型名称
	AuthorType        int      `json:"authorType"`        // 角色类型
	CategoryCode      int      `json:"categoryCode"`      // 上课模式
	CategoryName      string   `json:"categoryName"`      // 上课模式名称
	ChannelCode       int      `json:"channelCode"`       // 频道号
	ChannelName       string   `json:"channelName"`       // 频道名称
	CreatedAt         string   `json:"createdAt"`         // 创建时间
	Discount          int      `json:"discount"`          // 优惠价格
	FirstLabel        string   `json:"firstLabel"`        // 一级标签
	Price             int      `json:"price"`             // 销售价格
	ProductName       string   `json:"productName"`       // 课程标题
	SalesStatus       int      `json:"salesStatus"`       // 销售状态
	SalesStatusString string   `json:"salesStatusString"` // 销售状态名称
	SecondLabelList   []string `json:"secondLabelList"`   // 二级标签数组
	TagList           []string `json:"tagList"`           // 主题数组
	TeacherName       string   `json:"teacherName"`       // 讲师姓名
	UpdatedAt         string   `json:"updatedAt"`         // 更新时间

	ClassId       int64  `json:"classId"`       // 班次id
	MainVideo     string `json:"mainVideo"`     // 视频地址
	MainImage     string `json:"mainImage"`     // 图片封面
	MainVideUrl   string `json:"mainvideUrl"`   // 视频封面地址
	MainVideoName string `json:"mainVideoName"` // 视频封面名称
	FileId        string `json:"fileId"`        // 录播课 用来交换 视频的 主键

}

func (this *XzqCourseService) CourseList() ([]DxzCourseData, error) {
	conf := config.Conf{}
	config := conf.GetConf()

	resp := make([]DxzCourseData, 0)

	data := new(DxzCourseResData)
	data.PageNo = this.Page
	data.PageSize = this.Limit
	data.Condition.GoodsId = this.GoodsId
	data.Condition.ProductName = this.ProductName
	data.Condition.TeacherName = this.TeacherName
	data.Condition.CategoryCode = this.CategoryCode
	data.Condition.SalesStatus = this.SalesStatus
	data.Condition.AuthorType = this.AuthorType
	data.Condition.ChannelCode = this.ChannelCode
	data.Condition.FirstLabel = this.FirstLabel
	data.Condition.SecondLabel = this.SecondLabel
	data.Condition.Tag = this.Tag
	data.Condition.UpdatedAtEnd = this.UpdatedAtEnd
	data.Condition.UpdatedAtStart = this.UpdatedAtStart

	dataByte, err := json.Marshal(data)
	if err != nil {
		return resp, err
	}

	log.Println(string(dataByte))
	_, bodyByte, err := utils.DoRequest("POST", config.DxzXzqCourseUrl, dataByte, "")
	if err != nil {
		return resp, err
	}

	dxzRespData := new(DxzResponseData)
	err = json.Unmarshal(bodyByte, dxzRespData)
	if err != nil {
		return resp, err
	}

	log.Println(string(bodyByte))

	if !dxzRespData.Success {
		return resp, errors.New(dxzRespData.Message)
	}

	for _, v := range dxzRespData.Result.([]interface{}) {
		data, _ := json.Marshal(v.(map[string]interface{}))
		push := DxzCourseData{}
		json.Unmarshal(data, &push)
		resp = append(resp, push)
	}

	return resp, nil
}

type LessonRequest struct {
	Id int64 `json:"id"` // 班次id
}
type LessonResponse struct {
	ClassName     string `json:"className"`     // 班次名称
	CreatedId     int64  `json:"createdId"`     // 创建人ID
	DisableCount  int    `json:"disableCount"`  // 班次未审核数量
	Discount      int    `json:"discount"`      // 优惠价格
	ExamineStatus int    `json:"examineStatus"` // 审核状态 0:待审核，1:未通过，2:通过
	FailedReasons string `json:"failedReasons"` // 审核不通过原因
	GoodsId       int64  `json:"goodsId"`       // 课程ID
	GoodsItemsVo  []struct {
		CoursewareName        string `json:"coursewareName"`        // 课件名称
		CoursewareSafetyChain string `json:"coursewareSafetyChain"` // 课件防盗链
		CoursewareSize        int    `json:"coursewareSize"`        // 课件大小
		CoursewareUrl         string `json:"coursewareUrl"`         // 课件url
		CreatedId             int64  `json:"createdId"`             // 创建人id
		FailedReasons         string `json:"failedReasons"`         // 审核不通过原因
		FileId                string `json:"fileId"`                // 录播视频file_id
		FileName              string `json:"fileName"`              // 视频名称
		GoodsId               int64  `json:"goodsId"`               // 课程id
		Id                    int64  `json:"id"`                    //
		ItemName              string `json:"itemName"`              // 课节名称
		ItemNumber            int    `json:"itemNumber"`            // 课节序号
		ItemStatus            int    `json:"itemStatus"`            // 课节审核状态 0待审核、1审核未通过、2审核通过、3草稿箱(未提交审核)、4未上 integer(int32) 传
		PlanFromDate          string `json:"planFromDate"`          // 计划上课时间
		PlanToDate            string `json:"planToDate"`            // 计划下课课时间
		UpdatedId             int64  `json:"updatedId"`             // 修改人id
		VideoDuration         int    `json:"videoDuration"`         // 视频时长(秒)
	} `json:"goodsItemsVo"` // 班次名称
	GoodsSaleQuantity  int    `json:"goodsSaleQuantity"`  // 销售数量
	GoodsStockQuantity int    `json:"goodsStockQuantity"` // 库存
	GoodsTotalQuantity int    `json:"goodsTotalQuantity"` // 限定人数
	Id                 int64  `json:"id"`                 // 班次ID
	IsReview           int    `json:"isReview"`           // 是否支持回看 0:否，1:是
	Order              int64  `json:"order"`              // 班次序号
	Period             string `json:"period"`             // 班次周期
	PlanStatus         int    `json:"planStatus"`         // 计划状态 0:上架未开课，1:教学中，2:已结课
	SaleStatus         int    `json:"saleStatus"`         // 班次状态(1:待开班次 2:教学中班次 3:已结班次)

}

// 课节列表
func (this *XzqCourseService) LessonList(id int64) (LessonResponse, error) {
	conf := config.Conf{}
	config := conf.GetConf()

	resp := LessonResponse{}

	data := new(LessonRequest)
	data.Id = id
	dataByte, err := json.Marshal(data)
	if err != nil {
		return resp, err
	}

	log.Println(string(dataByte))
	_, bodyByte, err := utils.DoRequest("POST", config.DxzXzqLessonUrl, dataByte, "")
	if err != nil {
		return resp, err
	}

	dxzRespData := new(DxzResponseData)
	err = json.Unmarshal(bodyByte, dxzRespData)
	if err != nil {
		return resp, err
	}

	log.Println(string(bodyByte))

	if !dxzRespData.Success {
		return resp, errors.New(dxzRespData.Message)
	}

	dataResp, _ := json.Marshal(dxzRespData.Result.(map[string]interface{}))
	json.Unmarshal(dataResp, &resp)

	return resp, nil
}
