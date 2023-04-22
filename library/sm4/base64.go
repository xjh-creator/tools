package sm4

import "encoding/base64"

// Base64EncodeByStr Base64加密
func Base64EncodeByStr(str string) string { return Base64Encode([]byte(str)) }
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Base64DecodeByByte Base64解密
func Base64DecodeByByte(str []byte) ([]byte, error) { return Base64Decode(string(str)) }
func Base64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}
