package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"platform_report/config"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := ""
		Authorization := c.GetHeader("Authorization")
		token := strings.Split(Authorization, " ")

		if Authorization == "" {
			msg = "jwt is null"
		} else {
			userInfo, err := ParseToken(token[1])

			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					msg = "jwt is timeout"
				default:
					msg = "jwt checked fail"
				}
				c.JSON(200, gin.H{
					"code":  403,
					"error": msg,
					"data":  "",
				})
				c.Abort()
				return
			}
			c.Set("userInfo", userInfo)
		}
		if msg != "" {
			c.JSON(http.StatusOK, gin.H{
				"code":  403,
				"error": msg,
				"data":  "",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

type Claims struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	CompanyId       int    `json:"companyId"`
	Realname        string `json:"realname"`
	OrganizeId      int    `json:"organizeId"`
	SystemAuthorize string `json:"system_authorize"`
	jwt.StandardClaims
}
var jwtSecret []byte

func init() {
	cf := config.Conf{}
	config := cf.GetConf()
	jwtSecret = []byte(config.JwtSecret)
}


func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
