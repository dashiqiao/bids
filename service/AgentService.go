package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"platform_report/lib"
	"time"
)

type AgentService struct {
}

func (c *AgentService) hashID(args map[string]interface{}) string {
	bytes, _ := json.Marshal(args)
	md5str := fmt.Sprintf("%x", string(bytes))
	return md5str
}

func (c *AgentService) Save(uuid string, args map[string]interface{}, val gin.H) {
	key := uuid + ":" + c.hashID(args)
	lib.GetRedis().SAdd(uuid, key)
	jsonVal, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err.Error())
	}
	lib.GetRedis().Set(key, string(jsonVal), 2*time.Minute)
}

func (c *AgentService) Get(uuid string, args map[string]interface{}) gin.H {
	val := lib.GetRedis().Get(uuid + ":" + c.hashID(args)).Val()
	var items gin.H
	_ = json.Unmarshal([]byte(val), &items)
	return items
}

func (c *AgentService) Clean(uuid string) {
	invoke := lib.GetRedis()
	keys, _ := invoke.SMembers(uuid).Result()
	invoke.Del(keys...)
	invoke.Del(uuid)
}
