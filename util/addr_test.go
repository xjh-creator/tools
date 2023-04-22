package util

import (
	"fmt"
	"testing"
)

func TestGetFreePort(t *testing.T) {
	port, _ := GetFreePort()
	fmt.Println(port)
}
