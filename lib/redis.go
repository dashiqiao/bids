package lib

import (
	"fmt"
	"log"

	. "platform_report/config"

	"github.com/go-redis/redis"
)

var rcl *redis.Client

func InitRedis() {

	cf := Conf{}
	config := cf.GetConf()

	addr := fmt.Sprintf("%v:%v", config.RedisHost, config.RedisPort)
	rcl = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.RedisPassword, // no password set
		DB:       config.RedisDB,       // use default DB
		PoolSize: 1000,
	})

	if _, err := rcl.Ping().Result(); err != nil {
		log.Println("connection errors", err.Error())
	} else {
		log.Println("redis connected .....")
	}

}

func CloseRedis() {
	rcl.Close()
}
func GetRedis() *redis.Client {
	return rcl
}
