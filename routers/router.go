package routers

import (
	"platform_report/middleware"
	"platform_report/routers/api/v1"
	"platform_report/routers/api/xzq"
	"platform_report/routers/api/yyzx"

	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.JWT())

	r.POST("/api/v1/report/callback", v1.CallBack)

	apiv1 := r.Group("/api/v1/report/")
	apiv1.Use(FilterHandle(), middleware.SetUp())
	{
		//支持的筛选条件
		apiv1.GET("support", v1.Support)
		//表头定义
		apiv1.GET("header", v1.Header)
		//数据列表
		apiv1.POST("index", v1.Index)
		//数据导出
		apiv1.POST("export", v1.Export)
		//看板
		apiv1.GET("dashboard", v1.DashBoard)
		//清理缓存
		apiv1.GET("cc", v1.CC)
		//字段
		apiv1.GET("fields", v1.GetFields)
		//行为
		apiv1.POST("actions", v1.Actions)

		apiv1.GET("sign", v1.UpSign)
	}

	apiv2 := r.Group("/api/v1/tool/")
	apiv2.Use(FilterHandle(), middleware.SetUp())
	{
		apiv2.POST("tem", v1.AddTem)
		apiv2.GET("tem/list", v1.GetTemList)
		apiv2.GET("tem", v1.GetTem)
		apiv2.PUT("tem", v1.ModifyTem)
		apiv2.DELETE("tem", v1.DelTem)

		// excel 导入
		apiv2.POST("excelImport", v1.ExcelImport)
		// excel模板下载
		apiv2.GET("export/template", v1.ExportTemplate)
		// excel 导入模板
		apiv2.POST("excel/import/template", v1.ExcelTemplateImport)

		// 字典
		apiv2.GET("tec", v1.GetExcelCode)
		apiv2.PUT("tec", v1.ModifyExcelCode)
		apiv2.DELETE("tec", v1.DelExcelCode)
		apiv2.POST("tec", v1.AddExcelCode)
		apiv2.GET("tec/list", v1.GetExcelCodeList)
	}

	noFilter := r.Group("/api/v1/nofilter/")
	{
		noFilter.GET("tem/button", v1.GetTemButton)
		noFilter.GET("excel/tablelist", v1.ExcelTableList)
		noFilter.GET("tec/reload", v1.ReloadExcelCodeList)
	}

	adminv1 := r.Group("/admin/v1/")
	adminv1.Use(FilterHandle(), middleware.SetUp())
	{
		adminv1.GET("report/tablelist", v1.TableList)
		adminv1.GET("report/index", v1.ReportDefines)
		adminv1.POST("report/save", v1.AddReport)
		adminv1.POST("report/modify", v1.ModifyReport)
		adminv1.POST("report/valid", v1.SqlValid)
		adminv1.GET("report/details", v1.Details)
		adminv1.POST("report/conditions", v1.Conditions)
		adminv1.POST("report/delete", v1.Delete)
		adminv1.POST("report/rename", v1.Rename)
		adminv1.POST("report/summary", v1.SqlSummary)
		adminv1.POST("report/enable", v1.Enable)
		adminv1.POST("report/actions", v1.ReportActions)
		adminv1.POST("report/auto", v1.Auto)
		adminv1.GET("report/form", v1.AllForm)

		adminv1.POST("graph/save", v1.AddGraph)
		adminv1.POST("graph/modify", v1.ModifyGraph)
		adminv1.POST("graph/delete", v1.RemoveGraph)
		adminv1.GET("graph/index", v1.GetGraph)
		adminv1.GET("graph/detail", v1.GetGraphByID)
		adminv1.GET("graph/all", v1.GetAllGraph)
		adminv1.POST("graph/link", v1.LinkGraph)

		adminv1.GET("board/index", v1.BoardDefines)
		adminv1.POST("board/save", v1.AddBoard)
		adminv1.POST("board/delete", v1.RemoveBoard)
		adminv1.POST("board/rename", v1.RenameBoard)
		adminv1.POST("board/layout", v1.BoardLayout)

		adminv1.POST("board/target/save", v1.SaveBoardTarget)
		adminv1.POST("board/graph/save", v1.SaveBoardGraph)
		adminv1.GET("board/ls", v1.GetBoard)

		adminv1.GET("target/defines", v1.TargetDefines)
		adminv1.POST("target/save", v1.SaveTargetDefines)
		adminv1.POST("target/modify", v1.ModifyTargetDefines)
		adminv1.POST("target/delete", v1.DeleteTargetDefines)

		adminv1.GET("layout/index", v1.LayoutDefines)
		adminv1.POST("layout/modify", v1.ModifyLayout)
		adminv1.POST("layout/save", v1.SaveLayout)
		adminv1.POST("layout/delete", v1.RemoveLayout)

		adminv1.GET("analysis/defines", v1.CHDefines)
		adminv1.GET("analysis/index", v1.CHList)
		adminv1.POST("analysis/create", v1.CHAdd)
		adminv1.POST("analysis/update", v1.CHModify)
		adminv1.POST("analysis/delete", v1.CHRemove)
		adminv1.POST("analysis/percent", v1.CHPercent)
		adminv1.GET("analysis/ls", v1.CHUserList)
		adminv1.GET("analysis/actions", v1.CHUserActionsList)

		adminv1.GET("actions/index", v1.ActionDefines)
		adminv1.POST("actions/delete", v1.DeleteActions)
		adminv1.GET("actions/field", v1.ActionFields)
		adminv1.POST("actions/save", v1.SaveAction)
		adminv1.POST("actions/enable", v1.ActionEnable)
		adminv1.POST("actions/bind", v1.ActionBind)
		adminv1.POST("actions/unbind", v1.ActionUnBind)

		adminv1.GET("source/index", v1.SourceDefines)
		adminv1.POST("source/create", v1.AddSource)
		adminv1.POST("source/update", v1.UpdSource)
		adminv1.POST("source/delete", v1.RemoveSource)
	}

	xzqApi := r.Group("/api/v1/xzq")
	{
		// 获取课程列表获取课程列表
		xzqApi.POST("/course/list", xzq.CourseList)
		// 获取课节列表
		xzqApi.GET("/lesson/list", xzq.LessonList)
		// 获取文章列表
		xzqApi.POST("/essay/list", xzq.EssayList)
		// 获取展会列表
		xzqApi.POST("/exhibition/list", xzq.ExhibitionList)
		// 获取主题列表
		xzqApi.GET("/theme/list", xzq.ThemeList)
		// 一级类目
		xzqApi.GET("/first/label/list", xzq.FirstLabelList)
		// 二级类目
		xzqApi.GET("/second/label/list", xzq.SecondLabelList)
		// 获取视频地址
		xzqApi.GET("/file/video/url", xzq.VideoUrl)

		// 用户详情
		xzqApi.GET("/user/info", xzq.UserInfo)
		// 修改用户状态
		xzqApi.PUT("/user/status", xzq.ModifyUserStatus)
		// 机构详情
		xzqApi.GET("/institution/info", xzq.InstitutionInfo)
	}

	//运营中心
	yyzxApi := r.Group("/api/v1/yyzx/")
	{
		// 添加话题广告
		yyzxApi.POST("/topicad/create", yyzx.CreateTopicAd)
		// 查看话题广告
		yyzxApi.POST("/topicad/view", yyzx.ViewTopicAd)
	}

	//内容中心
	nrzxApi := r.Group("/api/v1/nrzx/")
	{
		// 审批话题
		nrzxApi.POST("/examine", yyzx.Examine)
		// 上下架
		nrzxApi.POST("/updown", yyzx.UpDown)
		// 举报审批
		nrzxApi.POST("/jb/examine", yyzx.JbExamine)
		// 举报上下架
		nrzxApi.POST("/jb/updown", yyzx.JbUpDown)
	}

	return r
}
