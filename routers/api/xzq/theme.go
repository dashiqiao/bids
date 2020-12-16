package xzq

import (
	"github.com/gin-gonic/gin"
	"platform_report/service"
)

func ThemeList(c *gin.Context) {
	channel := c.Query("channel")

	list, err := new(service.XzqThemeService).ThemeList(channel)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}

func FirstLabelList(c *gin.Context) {
	channel := c.Query("channel")

	list, err := new(service.XzqThemeService).FirstLabelList(channel)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}

func SecondLabelList(c *gin.Context) {
	pCode := c.Query("pCode")
	list, err := new(service.XzqThemeService).SecondLabelList(pCode)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}
