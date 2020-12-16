package lib

import (
	"github.com/jinzhu/gorm"
	"sync"
	. "platform_report/config"
)

var dblock *sync.Mutex = &sync.Mutex{}
var chatDbClient *gorm.DB

func GetDbInstance() *gorm.DB {
	//double check
	if chatDbClient == nil {
		dblock.Lock()
		defer dblock.Unlock()
		if chatDbClient == nil {
			cf := Conf{}
			config := cf.GetConf()
			var err error
			chatDbClient, err = gorm.Open("mysql", config.DbUsername+":"+config.DbPassword+"@("+config.DbHost+":"+config.DbPort+")/"+config.DbDatabase+"?charset=utf8&parseTime=True&loc=Local")
			if err != nil {

			}
			chatDbClient.SingularTable(true)
			chatDbClient.DB().SetMaxIdleConns(config.DbMaxIdleConns)
			chatDbClient.DB().SetMaxOpenConns(config.DbMaxOpenConns)

		}
	}

	return chatDbClient
}