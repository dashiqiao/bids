package middleware

import (
	"bytes"
	"github.com/Anderson-Lu/gofasion/gofasion"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"platform_report/lib"
	"platform_report/service"
	"strings"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func SetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.RequestURI, "/api/v1/tool/") {
			c.Next()
			return
		}
		//开始时间
		//startTime := time.Now().Unix()
		//处理请求

		uuid, args := "", ""
		if c.Request.Method == "GET" {
			uuid = c.Query("uuid")
		} else {
			_ = c.Request.ParseForm()

			args = c.Request.PostForm.Encode()
			//fmt.Println("args ==== ", args)
			uuid = c.Request.PostFormValue("uuid")
			//fmt.Println(uuid)
			if strings.TrimSpace(args) == "" {
				//c.Copy()
				var bodyBytes []byte
				bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				args = string(bodyBytes)
				fsion := gofasion.NewFasion(args)
				uuid = fsion.Get("uuid").ValueStr()
		}
		}
		c.Next()
		//结束时间
		//endTime := time.Now().Unix()
		//if c.Request.Method == "POST" {
		//	_ = c.Request.ParseForm()
		//}

		opType := 2
		if strings.Contains(c.Request.RequestURI, "/api/v1/report/") && !strings.Contains(c.Request.RequestURI, "/api/v1/report/actions") && !strings.Contains(c.Request.RequestURI, "/api/v1/report/cc") {
			opType = 1
		}

		if strings.Contains(c.Request.RequestURI, "/admin/v1/report/index") || strings.Contains(c.Request.RequestURI, "/admin/v1/analysis/index") || strings.Contains(c.Request.RequestURI, "/admin/v1/analysis/ls") || strings.Contains(c.Request.RequestURI, "/admin/v1/analysis/actions") {
			opType = 1
		}
		userInfo, _ := c.Get("userInfo")

		service.NewLogService().AddReportSubmitLog(uuid, c.Request.Method, c.Request.RequestURI, args, opType, lib.StringIpToInt(c.ClientIP()), userInfo.(*Claims).ID, userInfo.(*Claims).Realname)

	}
}
