package routers

import (
	"encoding/json"
	"log"
	"platform_report/middleware"
	"platform_report/pkg/e"
	"platform_report/pkg/util"
	"platform_report/routers/api"
	"platform_report/rpc"

	"github.com/gin-gonic/gin"
)

func FilterHandle() gin.HandlerFunc {
	return func(c *gin.Context) {

		// code = VerifyUserRuleRpc(c)
		// code := VerifyUserRuleRpcV1(c)
		code := VerifyUserRuleRpcV2(c)

		if code != e.SUCCESS {
			c.JSON(200, gin.H{
				"code":  code,
				"error": e.GetMsg(code),
				"data":  "",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

type RuleRpc struct {
	UserId    int
	CompanyId int
	System    string
	Method    string
	Router    string
}

func VerifyUserRuleRpc(c *gin.Context) int {
	body := RuleRpc{}
	code := e.SUCCESS

	userInfo, isExists := c.Get("userInfo")
	if !isExists {
		return e.ERROR_AUTH_CHECK_TOKEN_FAIL
	}
	user := userInfo.(*util.Claims)
	body.UserId = user.ID
	body.CompanyId = user.CompanyId
	body.System = "xzq"
	body.Method = c.Request.Method
	body.Router = c.FullPath()

	// log.Println(body)
	resp, err := rpc.RpcCall("VerifyUserRuleRpc", body)
	if err != nil {
		code = e.ERROR_RULE_AUTH_FAIL
	} else if resp.Code != e.SUCCESS {
		code = e.ERROR_RULE_AUTH_FAIL
	} else {
		c.Set("role", resp.Data)
	}

	return code
}

// 权限验证
func VerifyUserRuleRpcV1(c *gin.Context) int {
	// if code != e.SUCCESS {
	// 	return code
	// }
	// log.Println("权限验证")
	code := e.SUCCESS

	userInfo, isExists := c.Get("userInfo")
	if !isExists {
		return e.ERROR_AUTH_CHECK_TOKEN_FAIL
	}

	user := userInfo.(*middleware.Claims)
	body := make(map[string]interface{})
	body["userId"] = user.ID
	body["companyId"] = user.CompanyId
	body["system"] = "xzq"

	roleInfo := make(map[string]interface{})

	resp, err := rpc.RpcCall("GetUserRoleInfo", body)
	if err != nil {
		log.Println(err.Error())
		return e.ERROR_RULE_AUTH_FAIL
	}
	if resp.Data == nil {
		return e.ERROR_RULE_AUTH_FAIL
	}
	dataList := resp.Data.([]interface{})

	for _, v := range dataList {
		role := v.(map[string]interface{})
		if role["Types"].(int64) == 1 {
			roleInfo = role
			break
		}

		if _, ok := role["RuleList"]; ok {
			ruleList := role["RuleList"].([]interface{})
			for _, v1 := range ruleList {
				rule := v1.(map[string]interface{})
				// log.Println(rule)
				if rule["Func"].(string) == c.Request.Method && rule["Router"].(string) == c.FullPath() && rule["Status"].(int64) == 1 {
					if roleInfo["Type"] == nil || roleInfo["Type"].(int64) < role["Type"].(int64) {
						roleInfo = role
					}
				}
			}
		}

	}

	if roleInfo["Type"] == nil {
		return e.ERROR_NO_AUTH
	}
	c.Set("role", roleInfo)
	return code
}

// 权限验证
func VerifyUserRuleRpcV2(c *gin.Context) int {
	// if code != e.SUCCESS {
	// 	return code
	// }
	// log.Println("权限验证")
	code := e.SUCCESS

	userInfo, isExists := c.Get("userInfo")
	if !isExists {
		return e.ERROR_AUTH_CHECK_TOKEN_FAIL
	}

	user := userInfo.(*middleware.Claims)
	body := make(map[string]interface{})
	body["userId"] = user.ID
	body["companyId"] = user.CompanyId
	body["system"] = "xzq"

	roleInfo := api.UserRole{}
	roleList := make([]api.UserRole, 0)

	// redis := lib.GetRedisConnet()
	// key := fmt.Sprintf("%v_%v_%v_%v", "crm", "GetUserRoleInfo", user.CompanyId, user.ID)
	// infoStr := redis.Get(key).Val()
	// if infoStr != "" {
	// 	err := json.Unmarshal([]byte(infoStr), &roleList)
	// 	if err != nil {
	// 		// return list, err
	// 	}
	// } else {

	resp, err := rpc.RpcCall("GetUserRoleInfo", body)
	if err != nil {
		log.Println(err.Error())
		return e.ERROR_RULE_AUTH_FAIL
	}
	if resp.Data == nil {
		return e.ERROR_RULE_AUTH_FAIL
	}
	dataList := resp.Data.([]interface{})

	listByte, err := json.Marshal(dataList)
	if err != nil {
		// return list, err
	}

	err = json.Unmarshal(listByte, &roleList)
	if err != nil {
		// return list, err
	}
	// redis.Set(key, string(listByte), 5*time.Minute)
	// }

	log.Println(c.Request.Method, c.FullPath())
	for _, v := range roleList {

		if v.Types == 1 {
			roleInfo = v
			break
		}

		log.Println(v.RuleList)
		if v.Rules != "" {
			for _, rule := range v.RuleList {
				if rule.Func == c.Request.Method && rule.Router == c.FullPath() && rule.Status == 1 {
					if roleInfo.Type == 0 || roleInfo.Type < v.Type {
						roleInfo = v
					}
				}
			}
		}

	}

	if roleInfo.Type == 0 {
		return e.ERROR_NO_AUTH
	}
	c.Set("role", roleInfo)
	return code
}
