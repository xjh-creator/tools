package security

import (
	"fmt"
	"math"
	"strings"

	"github.com/Luzifer/go-openssl"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/grand"

	"github.com/gogf/gf/crypto/gmd5"
)

var pwdKey = "ZyZyZy"

func EncryptPass(pass string, salt string) (pwd string) {
	pwd, _ = gmd5.Encrypt(pass + salt + pwdKey)
	pwd = strings.ToLower(pwd)
	return
}

func GetSalt() string {
	return GetSaltPro(6)
}

func Decrypt(src string) string {
	o := openssl.New()
	dec, err := o.DecryptBytes("dgqd-key", []byte(src))
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}
	return string(dec)
}

func Encrypt(src string) string {
	o := openssl.New()
	dec, err := o.EncryptBytes("dgqd-key", []byte(src))
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
	}
	return string(dec)
}

var num2char = "123456789abcdefghijklmnpqrstuvwxyz"

// 10进制数转换   n 表示进制， 16 or 36
func NumToBHex(num int) string {
	numStr := ""
	n := len(num2char)
	for num != 0 {
		yu := num % n
		numStr = string(num2char[yu]) + numStr
		num = num / n
	}
	return strings.ToUpper(numStr)
}

// 36进制数转换   n 表示进制， 16 or 36
func Hex2Num(str string) int {
	str = strings.ToLower(str)
	v := 0.0
	n := len(num2char)
	length := len(str)
	for i := 0; i < length; i++ {
		s := string(str[i])
		index := strings.Index(num2char, s)
		v += float64(index) * math.Pow(float64(n), float64(length-1-i)) // 倒序
	}
	return int(v)
}

func GetSaltPro(n int) string {
	if n < 6 {
		n = 6
	}
	if n > 32 {
		n = 32
	}
	st := int(math.Pow(10, float64(n-1)))
	return gconv.String(grand.Intn(9*st) + st)
}
