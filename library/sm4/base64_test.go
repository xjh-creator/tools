package sm4

import (
	"fmt"
	"testing"
)

func TestBase64EncodeByStr(t *testing.T) {
	a := Base64Encode([]byte("小明"))
	fmt.Println("---", a)
	result, _ := Base64Decode(a)
	fmt.Println(string(result))
}

func TestBase64DecodeByByte(t *testing.T) {
	result, _ := Base64Decode("6I6e5b6u5Yqd5a+8IOaWh+aYjuWHuuihjA==")
	fmt.Println("+++", string(result))
}
