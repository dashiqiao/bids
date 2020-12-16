package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"platform_report/middleware"
)

func UserInfo(c *gin.Context) *middleware.Claims {
	userInfo, _ := c.Get("userInfo")
	return userInfo.(*middleware.Claims)
}

func CallBack(c *gin.Context) {
	//service.NewLogService().AddReportSubmitLog("", "POST", "/api/v1/report/callback", "", 0, lib.StringIpToInt(c.ClientIP()), UserInfo(c).ID, UserInfo(c).Realname)
	//c.JSON(http.StatusOK, nil)
}

func ErrorStd(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  msg,
	})
}

func SuccessStd(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}

func SuccessStds(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}
