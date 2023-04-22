package security

import (
	"fmt"
	"testing"
)

func TestNumToBHex(t *testing.T) {
	t.Log(NumToBHex(100000))
	t.Log(Hex2Num("3JI7"))
}

func TestEncrypt(t *testing.T) {
	name := "admin"
	password := "123456"
	fmt.Println(Encrypt(name))
	fmt.Println(Encrypt(password))
	fmt.Println(Decrypt(name))
	fmt.Println(Decrypt(password))
}
