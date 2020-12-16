package lib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/prestodb/presto-go-client/presto"
	. "platform_report/config"
	"reflect"
)

func GetPrestoUrl() string {
	//catalog, schema string
	cf := Conf{}
	config := cf.GetConf()
	prestoUrl := "http://" + config.PrestoUsername
	if config.PrestoPassword != "" {
		prestoUrl += ":" + config.PrestoPassword
	}
	//prestoUrl += "@" + config.PrestoHost + ":" + config.PrestoPort + "?catalog=" + catalog + "&schema=" + schema
	prestoUrl += "@" + config.PrestoHost + ":" + config.PrestoPort
	return prestoUrl
}

//func InitPresto() *sql.DB {
//	db, err := sql.Open("presto", GetPrestoUrl())
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil
//	}
//	return db
//}

type jsonNullInt64 struct {
	sql.NullInt64
}

func (v jsonNullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)
}

type jsonNullFloat64 struct {
	sql.NullFloat64
}

func (v jsonNullFloat64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Float64)
}

type jsonNullTime struct {
	mysql.NullTime
}

func (v jsonNullTime) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Time)
}

var jsonNullInt64Type = reflect.TypeOf(jsonNullInt64{})
var jsonNullFloat64Type = reflect.TypeOf(jsonNullFloat64{})
var jsonNullTimeType = reflect.TypeOf(jsonNullTime{})
var nullInt64Type = reflect.TypeOf(sql.NullInt64{})
var nullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
var nullTimeType = reflect.TypeOf(mysql.NullTime{})

func ExePrestoSqlQuery(conn, sqlExe string, args ...interface{}) ([]map[string]interface{}, error) {
	db, err := sql.Open("presto", conn)
	if err != nil {
		fmt.Println("presto1=====1=====>", err.Error())
		return nil, fmt.Errorf("can't connect to presto error: %v", err)
	}
	rows, err := db.Query(sqlExe, args...)
	if err != nil {
		fmt.Println("presto1=====2=====>", sqlExe, err.Error())
		return nil, fmt.Errorf("can't connect to presto error: %v", err)
	}
	defer rows.Close()
	if rows == nil {
		return nil, nil
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("column error: %v", err)
	}

	ct, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("column type error: %v", err)
	}

	types := make([]reflect.Type, len(ct))
	for i, tp := range ct {
		st := tp.ScanType()
		if st == nil {
			return nil, fmt.Errorf("scantype is null for column: %v", err)
		}
		switch st {
		case nullInt64Type:
			types[i] = jsonNullInt64Type
		case nullFloat64Type:
			types[i] = jsonNullFloat64Type
		case nullTimeType:
			types[i] = jsonNullTimeType
		default:
			types[i] = st
		}
	}
	values := make([]interface{}, len(ct))
	var slice []map[string]interface{}
	for rows.Next() {
		for i := range values {
			values[i] = reflect.New(types[i]).Interface()
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan values: %v", err)
		}
		data := make(map[string]interface{})
		for i, v := range values {
			//fmt.Println(reflect.TypeOf(v).Kind())
			data[columns[i]] = v
		}
		slice = append(slice, data)
	}
	var mapResult []map[string]interface{}
	jsons, _ := json.Marshal(slice)
	err = json.Unmarshal(jsons, &mapResult)
	return mapResult, err
}
