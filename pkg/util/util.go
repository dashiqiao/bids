package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"platform_report/lib"
	"reflect"
	"strconv"
	"strings"
	"time"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/jinzhu/now"
)

func StructToMap(v interface{}) map[string]interface{} {
	bytes, _ := ffjson.Marshal(v)
	ret := make(map[string]interface{})
	_ = ffjson.Unmarshal(bytes, &ret)
	return ret
}


func TimeIntToString(ts int64) string {
	if ts == 0 {
		return ""
	}
	tm := time.Unix(ts, 0)
	return tm.Format("2006-01-02 15:04:05")
}
func InArrayStr(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func InArrayStrIndex(items []string, item string) int {
	i := 0
	for _, eachItem := range items {
		if eachItem == item {
			return i
		}
		i++
	}
	return 0
}
func InArray(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func InArrayInt64(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

/**
 * 格式化字节大小
 * @param  number $size     字节数
 * @param  string $delimiter 数字和单位分隔符
 * @return string               格式化后的带单位的大小
 * @author
 */
func FormatBytes(size int) string {
	// $units = array('B', 'KB', 'MB', 'GB', 'TB', 'PB');
	// for ($i = 0; $size >= 1024 && $i < 5; $i++) $size /= 1024;
	// return round($size, 2) . $delimiter . $units[$i];
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	i := 0
	sizeF := float64(size)
	for ; sizeF >= 1024 && i < 5; i++ {
		sizeF = sizeF / 1024
	}
	// sizeF, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", sizeF), 64)
	return fmt.Sprintf("%.2f %v", sizeF, units[i])
}

/*取得文件后缀*/
func GetExtension(fileName string) string {
	nameList := strings.Split(fileName, ".")
	return nameList[len(nameList)-1]

}

func GetFullPath(filename string) string {
	// 没写完
	url := "" + filename
	return url
}

func ConversionString(val interface{}) string {
	str := ""
	//fmt.Println(reflect.TypeOf(val).Kind())
	switch reflect.TypeOf(val).Kind() {
	case reflect.Slice:
		break
	case reflect.String:
		str = val.(string)
	case reflect.Int:
		str = fmt.Sprintf("%v", val.(int))
	case reflect.Int8:
		str = fmt.Sprintf("%v", val.(int8))
	case reflect.Int32:
		str = fmt.Sprintf("%v", val.(int32))
	case reflect.Int64:
		str = fmt.Sprintf("%v", val.(int64))
	case reflect.Uint:
		str = fmt.Sprintf("%v", val.(uint))
	case reflect.Uint8:
		str = fmt.Sprintf("%v", val.(uint8))
	case reflect.Uint16:
		str = fmt.Sprintf("%v", val.(uint16))
	case reflect.Uint32:
		str = fmt.Sprintf("%v", val.(uint32))
	case reflect.Uint64:
		str = fmt.Sprintf("%v", val.(uint64))
	case reflect.Float32:
		str = fmt.Sprintf("%v", val.(float32))
	case reflect.Float64:
		str = fmt.Sprintf("%v", val.(float64))
	case reflect.Map:
	case reflect.Struct:
	case reflect.Uintptr:
	case reflect.UnsafePointer:
	case reflect.Array:
	case reflect.Bool:
	case reflect.Chan:
	case reflect.Complex64:
	case reflect.Complex128:
	}
	return str
}

func ConversionMap(data map[string]interface{}) map[string]interface{} {
	resp := make(map[string]interface{})
	for k, v := range data {
		log.Println("字段名:", k, "  字段类型:", reflect.TypeOf(data[k]).Kind(), "  字段值:", v)
		// log.Println(reflect.TypeOf(data[k]))
		// log.Println(reflect.TypeOf(data[k]).Kind())
		// log.Println(v)
	}

	// switch reflect.TypeOf(val).Kind() {
	// case reflect.Slice:
	// 	break
	// case reflect.String:
	// 	str = val.(string)
	// case reflect.Int:
	// 	str = fmt.Sprintf("%v", val.(int))
	// case reflect.Int8:
	// 	str = fmt.Sprintf("%v", val.(int8))
	// case reflect.Int32:
	// 	str = fmt.Sprintf("%v", val.(int32))
	// case reflect.Int64:
	// 	str = fmt.Sprintf("%v", val.(int64))
	// case reflect.Uint:
	// 	str = fmt.Sprintf("%v", val.(uint))
	// case reflect.Uint8:
	// 	str = fmt.Sprintf("%v", val.(uint8))
	// case reflect.Uint16:
	// 	str = fmt.Sprintf("%v", val.(uint16))
	// case reflect.Uint32:
	// 	str = fmt.Sprintf("%v", val.(uint32))
	// case reflect.Uint64:
	// 	str = fmt.Sprintf("%v", val.(uint64))
	// case reflect.Float32:
	// 	str = fmt.Sprintf("%v", val.(float32))
	// case reflect.Float64:
	// 	str = fmt.Sprintf("%v", val.(float64))
	// case reflect.Map:
	// case reflect.Struct:
	// case reflect.Uintptr:
	// case reflect.UnsafePointer:
	// case reflect.Array:
	// case reflect.Bool:
	// case reflect.Chan:
	// case reflect.Complex64:
	// case reflect.Complex128:
	// }
	return resp
}

func Int64ToInt(num int64) int {
	if num == int64(0) {
		return 0
	}
	numStr := strconv.FormatInt(num, 10)
	numInt, _ := strconv.Atoi(numStr)
	return numInt
}

func RemoveSpace(str string) string {
	// 去除空格
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	return str
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}
func TimeParse(str string) time.Time {
	dt, err := time.Parse("2006-01-02", str) //15:04:05
	if err != nil {
		fmt.Println(err.Error())
	}
	return dt
}

func GetSubDays(t1, t2 time.Time) int {
	hours := t2.Sub(t1).Hours()
	if hours <= 0 {
		return 1
	}
	if (int(hours) % 24) == 0 {
		return int(hours)/24 + 1
	}
	return int(hours)/24 + 2
}

func GetPkDays(t1, t2 string) (string, string) {
	w1, w2 := TimeParse(t1), TimeParse(t2)
	num := GetSubDays(w1, w2)
	return lib.DateFormat(w1.AddDate(0, 0, num*-1)), lib.DateFormat(w2.AddDate(0, 0, num*-1))

}

func GetTimeByType(format string) (time.Time, time.Time) {
	now.WeekStartDay = time.Monday
	switch format {
	case "yesterday":
		return now.BeginningOfDay().AddDate(0, 0, -1), now.EndOfDay().AddDate(0, 0, -1)
	case "week":
		return now.BeginningOfWeek(), now.EndOfWeek()
	case "lastWeek":
		return now.BeginningOfWeek().AddDate(0, 0, -7), now.EndOfWeek().AddDate(0, 0, -7)
	case "month":
		return now.BeginningOfMonth(), now.EndOfMonth()
	case "lastMonth":
		return now.BeginningOfMonth().AddDate(0, -1, 0), now.EndOfMonth().AddDate(0, -1, 0)
	case "quarter":
		return now.BeginningOfQuarter(), now.EndOfQuarter()
	case "lastQuarter":
		return now.BeginningOfQuarter().AddDate(0, 3, 0), now.EndOfQuarter().AddDate(0, 3, 0)
	case "year":
		return now.BeginningOfYear(), now.EndOfYear()
	case "lastYear":
		return now.BeginningOfYear().AddDate(-1, 0, 0), now.EndOfYear().AddDate(-1, 0, 0)
	default:
		return now.BeginningOfDay(), now.EndOfDay()
	}
}

func Post(url string, data interface{}, contentType string, token string) (content string, err error) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		// log.Println(err.Error())
		// panic(err)
		return content, err
	}
	if contentType != "" {
		req.Header.Add("content-type", contentType)
	}

	if token != "" {
		req.Header.Add("Authorization", token)
	}

	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		// panic(error)
		return content, err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

func Get(url string, token string) (content string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// log.Println(err.Error())
		// panic(err)
		return content, err
	}

	if token != "" {
		req.Header.Add("Authorization", token)
	}

	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		// panic(error)
		return content, err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

type GetVerifyAuthorityResp struct {
	Code int       `json:"code"`
	Data AdminRule `json:"data"`
	Msg  string    `json:"msg"`
}
type AdminRule struct {
	Id        int    `json:"id"`
	Types     int    `json:"types"`      // 0系统设置 1工作台 2客户管理 3项目管理 4人力资源 5财务管理 6商业智能
	Title     string `json:"title"`      // 名称
	Name      string `json:"name"`       // 定义
	Func      string `json:"func"`       // 方法 配合权限管理
	Router    string `json:"router"`     // 路由 配合菜单管理
	Level     int    `json:"level"`      // 级别 1模块 2控制器 3操作  1 1级菜单 2 2级菜单 3 3级菜单
	Pid       int    `json:"pid"`        // 父id，默认0
	Status    int    `json:"status"`     // 状态，1启用，0禁用
	Dtime     int64  `json:"dtime"`      // 删除时间
	CompanyId int    `json:"company_id"` // 公司id
}
