package dao

import (
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"os"
	"platform_report/lib"
	"reflect"
	"strconv"
	"strings"
)

type InfiniteStones struct {
}

func NewInfiniteStones() *InfiniteStones {
	return new(InfiniteStones)
}

func (c *InfiniteStones) Valid(sql string) ([]map[string]interface{}, error) {
	ret, err := lib.InitXormMySql().SQL(sql).Query().List()
	return ret, err
}

func (c *InfiniteStones) War(sql string) (*[]map[string]interface{}, error) {
	ret, err := lib.InitXormMySql().SQL(sql).Query().List()
	return &ret, err
}

func (c *InfiniteStones) WarWithParametes(sql string, args []interface{}) (map[string]interface{}, error) {
	ret, err := lib.InitXormMySql().SQL(sql, args...).Query().List()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret[0], err
	//return nil, err
}

func (c *InfiniteStones) PtrWarWithParametes(sql string, args []interface{}) (*[]map[string]interface{}, error) {
	ret, err := lib.InitXormMySql().SQL(sql, args...).Query().List()
	return &ret, err
}

func (c *InfiniteStones) CopyWarWithParametes(sql string, args []interface{}) ([]map[string]interface{}, error) {
	ret, err := lib.InitXormMySql().SQL(sql, args...).Query().List()
	return ret, err
}

func (c *InfiniteStones) Count(sql string) interface{} {
	retMap := lib.InitXormMySql().SQL(sql).Query()
	count, _ := retMap.Count()
	if count > 0 {
		return retMap.Result[0]["num"]
	}
	return 0
}

func (c *InfiniteStones) CountWithParametes(sql string, args []interface{}) interface{} {
	retMap := lib.InitXormMySql().SQL(sql, args...).Query()
	count, _ := retMap.Count()
	if count > 0 {
		return retMap.Result[0]["num"]
	}
	return 0
}

func (c *InfiniteStones) SqlTemplate(sql string, args map[string]interface{}) []map[string]interface{} {
	templateId := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	_ = os.Remove("./sql/" + templateId + ".stpl")
	file, _ := os.OpenFile("./sql/"+templateId+".stpl", os.O_RDWR|os.O_CREATE, 0766);
	_, _ = file.WriteString(sql)
	_ = file.Close()
	results, err := lib.InitXormMySql().SqlTemplateClient(templateId+".stpl", &args).Query().List()
	fmt.Println(results)
	if err != nil {
		fmt.Println(err.Error())
	}
	return results

}

func MapTryConvert(data []map[string]string) []map[string]interface{} {
	dataInfo := make([]map[string]interface{}, 0)
	for _, v := range data {
		push := make(map[string]interface{})
		for key, value := range v {
			if i, err := strconv.ParseInt(value, 10, 64); err == nil {
				push[key] = i
			} else if i, err := strconv.ParseFloat(value, 64); err == nil {
				push[key] = i
			} else {
				push[key] = value
			}
		}
		dataInfo = append(dataInfo, push)
	}
	return dataInfo
}

func (c *InfiniteStones) CkQuery(sql string) ([]map[string]interface{}, error) {
	rows, err := lib.InitClickHouse().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return Rows2Interfaces(rows)
}

func Rows2Interfaces(rows *sql.Rows) (resultsSlice []map[string]interface{}, err error) {
	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		result, err := row2mapInterface(rows, fields)
		if err != nil {
			return nil, err
		}
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func row2mapInterface(rows *sql.Rows, fields []string) (resultsMap map[string]interface{}, err error) {
	resultsMap = make(map[string]interface{}, len(fields))
	scanResultContainers := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		var scanResultContainer interface{}
		scanResultContainers[i] = &scanResultContainer
	}
	if err := rows.Scan(scanResultContainers...); err != nil {
		return nil, err
	}

	for ii, key := range fields {
		// log.Println(ii, key)
		resultsMap[key] = reflect.Indirect(reflect.ValueOf(scanResultContainers[ii])).Interface()
	}
	return
}
