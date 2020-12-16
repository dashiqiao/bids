package dao

import (
	"platform_report/lib"
)

func MapToInsertSql(tableName string, info map[string]interface{}) (string, []interface{}) {
	sql := "insert into " + tableName + " ( "

	valueSql := "values ("

	valueList := make([]interface{}, 0)

	i := 0
	for k, v := range info {
		i += 1
		if i == len(info) {
			sql = sql + k + " )"
			valueSql = valueSql + "?" + " ) "
		} else {
			sql = sql + k + " , "
			valueSql = valueSql + "?" + " , "
		}
		valueList = append(valueList, v)
	}

	return sql + valueSql, valueList
}

func MapToUpdateSql(tableName string, info map[string]interface{}, key []string, id []interface{}) (string, []interface{}) {

	// UPDATE table_name SET field1=new-value1, field2=new-value2
	// [WHERE Clause]

	sql := "update " + tableName + " set  "

	whereSql := " where "

	valueList := make([]interface{}, 0)

	i := 0
	for k, v := range info {
		i += 1
		if i == len(info) {
			sql = sql + k + " = ? "

		} else {
			sql = sql + k + " = ? , "

		}
		valueList = append(valueList, v)
	}

	whereSql = whereSql + key[0] + " = ?"
	valueList = append(valueList, id[0])
	if len(key) > 1 {
		for i := 1; i < len(key); i++ {
			whereSql += " AND " + key[i] + " = ?"
			valueList = append(valueList, id[i])
		}
	}
	//fmt.Println(sql + whereSql)
	//fmt.Println("111111111", valueList)
	return sql + whereSql, valueList
}

func Count(tableName string, key []string, value []interface{}) int {
	whereSql := ""
	whereSql = whereSql + key[0] + " = ?"
	if len(key) > 1 {
		for i := 1; i < len(key); i++ {
			whereSql += " OR " + key[i] + " = ?"
		}
	}
	var cnt int
	lib.GetDbInstance().Table(tableName).Where(whereSql, value...).Count(&cnt)
	return cnt
}

func Exec(sql string, args []interface{}) error {
	return lib.GetDbInstance().Debug().Exec(sql, args...).Error
}
