package xzq

import (
	"github.com/gin-gonic/gin"
	"platform_report/service"
)

func ExhibitionList(c *gin.Context) {
	svc := service.XzqExhibitionService{}
	err := c.ShouldBind(&svc)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	list, err := svc.ExhibitionList()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}
