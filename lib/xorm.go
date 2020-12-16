package lib

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	. "platform_report/config"
)

func InitXormMySql() *xorm.Engine {
	cf := Conf{}
	config := cf.GetConf()
	dbAddress := config.DbUsername + ":" + config.DbPassword + "@tcp(" + config.DbHost + ":" + config.DbPort + ")/" + config.
		DbDatabase + "?charset=utf8"

	engine, _ := xorm.NewEngine("mysql", dbAddress)
	engine.SetMaxOpenConns(config.DbMaxOpenConns)
	engine.SetMaxIdleConns(config.DbMaxIdleConns)
	//_ = engine.RegisterSqlTemplate(xorm.Pongo2("./sql", ".stpl"))
	_ = engine.Ping()

	engine.ShowSQL(true)
	return engine
}
