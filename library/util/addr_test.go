package util

import (
	"fmt"
	"testing"
)

func TestUrlResolve(t *testing.T) {
	str := UrlResolve("https://www.myoumuamua.com/fzappdowntest61", "../fzmanage/index.html")
	fmt.Println(str)
}
