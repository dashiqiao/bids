package lib

import (
	"bytes"
	"github.com/axgle/mahonia"
	"github.com/jinzhu/now"
	"strconv"
	"strings"
	"time"
)

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
		return now.BeginningOfMonth().AddDate(0, -1, 0), now.BeginningOfMonth().AddDate(0, 0, -1)
	case "quarter":
		return now.BeginningOfQuarter(), now.EndOfQuarter()
	case "lastQuarter":
		return now.BeginningOfQuarter().AddDate(0, -3, 0), now.BeginningOfQuarter().AddDate(0, 0, -1)
	case "year":
		return now.BeginningOfYear(), now.EndOfYear()
	case "lastYear":
		return now.BeginningOfYear().AddDate(-1, 0, 0), now.BeginningOfYear().AddDate(0, 0, -1)
	default:
		return now.BeginningOfDay(), now.EndOfDay()
	}
}

func TimeFormat(now time.Time) string {
	return now.Format("2006-01-02 15:04:05")
}

func DateFormat(now time.Time) string {
	return now.Format("2006-01-02")
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func StringIpToInt(ipstring string) int {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt int = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func IpIntToString(ipInt int) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[len-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}
