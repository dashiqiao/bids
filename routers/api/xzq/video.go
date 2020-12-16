package xzq

import (
	"github.com/gin-gonic/gin"
	"platform_report/service"
)

func VideoUrl(c *gin.Context) {
	fileId := c.Query("fileId")

	list, err := new(service.XzqVideoService).MultiPlay(fileId, 100)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}
