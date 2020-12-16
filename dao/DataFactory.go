package dao

import (
	"platform_report/lib"
)

type IDataSource interface {
	Do(conn, sql string, args []interface{}) ([]map[string]interface{}, error)
	Count(conn, sql string, args []interface{}) interface{}
}

type MySql struct {
}

func (c *MySql) Do(conn, sql string, args []interface{}) ([]map[string]interface{}, error) {
	return lib.InitXormMySql().SQL(sql, args...).Query().List()
}

func (c *MySql) Count(conn, sql string, args []interface{}) interface{} {
	retMap := lib.InitXormMySql().SQL(sql, args...).Query()
	count, _ := retMap.Count()
	if count > 0 {
		return retMap.Result[0]["num"]
	}
	return 0
}

type Presto struct {
}

func (c *Presto) Do(conn, sql string, args []interface{}) ([]map[string]interface{}, error) {
	return lib.ExePrestoSqlQuery(conn, sql, args...)
}

func (c *Presto) Count(conn, sql string, args []interface{}) interface{} {
	retMap, _ := c.Do(conn, sql, args)
	if len(retMap) > 0 {
		return retMap[0]["num"]
	}
	return 0
}

type DataAdapter struct {
	Conn       string
	DataSource IDataSource
}

func (c *DataAdapter) Adapter(source int) {
	if source == 0 {
		c.DataSource = new(MySql)
	} else {
		c.Conn = lib.GetPrestoUrl()
		c.DataSource = new(Presto)
	}
}

func (c *DataAdapter) Do(sql string, args []interface{}) ([]map[string]interface{}, error) {
	return c.DataSource.Do(c.Conn, sql, args)
}

func (c *DataAdapter) Count(sql string, args []interface{}) interface{} {
	return c.DataSource.Count(c.Conn, sql, args)
}
