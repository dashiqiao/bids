package main

import (
	_ "go.uber.org/automaxprocs"
	"platform_report/config"
	"platform_report/lib"
	"platform_report/routers"
)

func main() {
	cf := config.Conf{}
	r := routers.InitRouter()

	lib.InitRedis()

	//_ = endless.ListenAndServe(":"+cf.GetConf().Port, r)

	_ = r.Run(":" + cf.GetConf().Port)

	//endless.ListenAndServe(":"+cf.GetConf().Port, r)

}
