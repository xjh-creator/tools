package ptime

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"testing"
)

func TestGetTodayZeroUnix(t *testing.T) {
	GetTodayZeroUnix()
}

func TestGetTomorrowZeroUnix(t *testing.T) {
	GetTomorrowZeroUnix()
}

func TestGetNowUnix(t *testing.T) {
	GetNowUnix()
}

func TestGetUnix(t *testing.T) {
	result := GetUnix(0, -3, -5)
	fmt.Println("三个月前的时间戳:", result)
	fmt.Println(gtime.New(1671177737624))
}
