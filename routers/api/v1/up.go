package v1

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	. "platform_report/config"
	"strings"
	"time"
)

func UpSign(c *gin.Context) {
	cf := Conf{}
	config := cf.GetConf()
	fmt.Println(config.UpyunSecret)
	c.JSON(200, gin.H{
		"code":   200,
		"bucket": "xuezhangquan",
		"uri":    "http://v0.api.upyun.com/xuezhangquan",
		//"bucket":        "dxz-marketing",
		//"uri":           "http://v0.api.upyun.com/dxz-marketing",
		"authorization": sign(config.UpyunOprator, md5Str(config.UpyunSecret), "POST", "/", makeRFC1123Date(time.Now()), "", ""),
	})
}

func md5Str(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
func makeRFC1123Date(d time.Time) string {
	utc := d.UTC().Format(time.RFC1123)
	return strings.Replace(utc, "UTC", "GMT", -1)
}
func base64ToStr(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func sign(key, secret, method, uri, date, policy, md5 string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	elems := []string{}
	for _, v := range []string{method, uri, date, policy, md5} {
		if v != "" {
			elems = append(elems, v)
		}
	}
	value := strings.Join(elems, "&")
	//fmt.Println(value)
	mac.Write([]byte(value))
	signStr := base64ToStr(mac.Sum(nil))
	return "UPYUN " + key + ":" + signStr
}
