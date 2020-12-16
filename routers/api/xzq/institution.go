package xzq

import (
	"github.com/gin-gonic/gin"
	"platform_report/service"
	"strconv"
)

func InstitutionInfo(c *gin.Context) {
	userId := c.Query("userId")
	institutionSvc := new(service.XzqInstitutionService)
	institutionSvc.UserId, _ = strconv.ParseInt(userId, 10, 64)

	list, err := institutionSvc.InstitutionInfo()
	if err != nil {
		Response(c, 400, "", err.Error())
		return
	}

	Response(c, 200, list, "")
	return
}
