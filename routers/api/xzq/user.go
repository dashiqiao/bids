package xzq

import (
	"github.com/gin-gonic/gin"
	"platform_report/service"
	"strconv"
)

func UserInfo(c *gin.Context) {
	userId := c.Query("userId")
	userSvc := new(service.XzqUserService)
	userSvc.UserId, _ = strconv.ParseInt(userId, 10, 64)

	list, err := userSvc.UserInfo()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}

func ModifyUserStatus(c *gin.Context) {
	userStatus := new(service.UserStatusData)
	err := c.ShouldBind(userStatus)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	err = new(service.XzqUserService).ModifyUserStatus(userStatus)
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, "", "")
	return
}
