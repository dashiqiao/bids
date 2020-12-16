package lib

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"log"
	. "platform_report/config"
)

func InitClickHouse() *sql.DB {
	cf := Conf{}
	config := cf.GetConf()
	dbAddress := "tcp://" + config.CkHost + ":" + config.CkPort + "?"
	if config.CkUsername != "" {
		dbAddress += "username=" + config.CkUsername + "&password=" + config.CkPassword
		dbAddress += "&database=" + config.CkDatabase
	} else {
		dbAddress += "debug=true&database=" + config.CkDatabase
	}
	connect, err := sql.Open("clickhouse", dbAddress)
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
	return connect
}
