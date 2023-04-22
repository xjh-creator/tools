package sm4

import (
	"fmt"
	"testing"
)

func TestSM4EN(t *testing.T) {
	fmt.Println(SM4EN("王木木"))
	fmt.Println(SM4EN("441900198903074480"))
	fmt.Println(SM4EN("东莞市莞太路"))
	fmt.Println(SM4EN("13434975089"))
}
