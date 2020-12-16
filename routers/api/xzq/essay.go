package xzq

import (
	"github.com/gin-gonic/gin"
	"platform_report/service"
)

func EssayList(c *gin.Context) {
	svc := service.XzqEssayService{}
	err := c.ShouldBind(&svc)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	list, err := svc.EssayList()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}
