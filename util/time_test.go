package util

import (
	"fmt"
	"testing"
)

func TestStringToTime(t *testing.T) {
	//timeStr := "0000-00-00 00:00:00"
	timeStr := "2022-01-02 15:04:05"
	res, _ := StringToTime(timeStr)
	fmt.Println(res)
}
