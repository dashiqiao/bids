package xzq

import (
	// "encoding/json"
	// "io/ioutil"
	// "platform_report/dao"
	"platform_report/service"
	"strconv"

	// "strconv"

	// "github.com/Anderson-Lu/gofasion/gofasion"
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

func CourseList(c *gin.Context) {
	svc := service.XzqCourseService{}
	err := c.ShouldBind(&svc)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	list, err := svc.CourseList()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}

func LessonList(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	list, err := new(service.XzqCourseService).LessonList(id)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}
