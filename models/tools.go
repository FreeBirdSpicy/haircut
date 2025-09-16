package models

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

const Timestamp = "2006-01-02"
const Timestamp12 = "2006-01-02 03:04:05"
const Timestamp24 = "2006-01-02 15:04:05"

var CurMonth string

type MonthMenu struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	On    string `json:"on"`
}

func init() {
	CurMonth = time.Now().Format("2006-01")
}

// 时间戳转换成日期
func TimestampToTime(timestamp int) string {
	return time.Unix(int64(timestamp), 0).Format(Timestamp24)
}

// 日期转换成时间戳
func TimeToTimestamp(timeStr string) int {
	t, _ := time.Parse(Timestamp24, timeStr)
	return int(t.Unix())
}

// 获取当前时间戳
func GetTimestamp() int {
	return int(time.Now().Unix())
}

// 获取当前日期
func GetDate() string {
	return time.Now().Format(Timestamp)
}

// 获取当前时间
func GetTime() string {
	return time.Now().Format(Timestamp24)
}

// 任意类型变量转换成字符串
func ToString(value interface{}) string {
	var result string
	if value == nil {
		return result
	}

	switch v := value.(type) {
	case string:
		result = v
	case int, int8, int16, int32, int64:
		result = strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		result = strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
	case float32:
		result = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		result = strconv.FormatFloat(v, 'f', -1, 64)
	case []byte:
		result = string(v)
	default:
		jsonValue, _ := json.Marshal(v)
		result = string(jsonValue)
	}
	return result
}

// 获取最近半年的月份
func GetLastHalfYear() []MonthMenu {
	monthMenuList := []MonthMenu{}
	for i := 0; i < 6; i++ {
		name := time.Now().AddDate(0, -i, 0).Format("1月")
		value := time.Now().AddDate(0, -i, 0).Format("2006-01")
		item := MonthMenu{
			Name:  name,
			Value: value,
		}
		monthMenuList = append(monthMenuList, item)
	}
	return monthMenuList
}

// 获取指定月份的第一天和最后一天日期
func GetMonthFirstAndLast(month string) (string, string) {
	t, _ := time.Parse("2006-01", month)
	start := t.Format(Timestamp)
	end := t.AddDate(0, 1, -1).Format(Timestamp)
	return start, end
}
