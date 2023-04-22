package ptime

import (
	"fmt"
	"time"
)

// GetNowUnix 获取当前点时间戳
func GetNowUnix() {
	fmt.Println(time.Now().UnixMilli())
}

// GetTodayZeroUnix 获取当天 0 点时间戳
func GetTodayZeroUnix() int64 {
	t := time.Now()
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	timeSamp := addTime.UnixMilli()
	return timeSamp
}

// GetTomorrowZeroUnix 获取明天 0 点时间戳
func GetTomorrowZeroUnix() int64 {
	nowTimeStr := time.Now().Format("2006-01-02")
	//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
	t2, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
	tomTimeStamp := t2.AddDate(0, 0, 1).UnixMilli()
	return tomTimeStamp
}

// GetYesterdayZeroUnix 获取昨天 0 点时间戳
func GetYesterdayZeroUnix() int64 {
	ts := time.Now().AddDate(0, 0, -1)
	timeYesterday := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
	timeStampYesterday := timeYesterday.UnixMilli()
	return timeStampYesterday
}

// GetUnix 获取某个时间的时间戳
func GetUnix(year, month, day int) int64 {
	ts := time.Now().AddDate(year, month, day)
	fmt.Printf("oneDayLater: start.AddDate(0, -3, -5) = %v\n", ts)
	timeYesterday := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
	timeStampYesterday := timeYesterday.UnixMilli()
	return timeStampYesterday
}
