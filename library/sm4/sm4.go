package sm4

import (
	"bytes"
	"crypto/cipher"
	"github.com/pkg/errors"

	"github.com/tjfoc/gmsm/sm4"
)

var (
	// 128比特密钥
	KEY = []byte("*&JOjuikNJ)(&$^o")
	// 128比特iv
	IV = []byte("7ajkshdkad&)1KKo")
)

func SM4EN(str string) string {
	if str == "" {
		return ""
	}
	//res, err := sm4.Sm4Ecb(KEY, []byte(str), true)
	res, err := SM4Encrypt(KEY, IV, []byte(str))
	if err != nil {
		//glog.Warnf("国密4加密失败:%s", err)
		return str
	}

	return Base64Encode(res)
}

func SM4DE(str string) string {
	if str == "" {
		return ""
	}
	code, err := Base64Decode(str)
	//res, err := sm4.Sm4Ecb(KEY, code, false)
	res, err := SM4Decrypt(KEY, IV, code)
	if err != nil {
		//glog.Warnf("%s 国密4解密失败:%s", str, err)
		return str
	}
	return string(res)
}

func SM4DE2(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	code, err := Base64Decode(str)
	//res, err := sm4.Sm4Ecb(KEY, code, false)
	res, err := SM4Decrypt(KEY, IV, code)
	if err != nil {
		//glog.Warnf("%s 国密4解密失败:%s", str, err)
		return "", err
	}
	return string(res), nil
}

var (
	enblockMode cipher.BlockMode
	enblockSize int
	deblockMode cipher.BlockMode
)

func SM4Encrypt(key, iv, plainText []byte) ([]byte, error) {
	//if enblockMode == nil {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//}
	origData := pkcs5Padding(plainText, blockSize)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

func SM4Decrypt(key, iv, cipherText []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)

	origData := make([]byte, len(cipherText))
	if len(cipherText)%blockMode.BlockSize() != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	if len(origData) < len(cipherText) {
		return nil, errors.New("crypto/cipher: output smaller than input")
	}
	blockMode.CryptBlocks(origData, cipherText)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

// pkcs5填充
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return nil
	}
	unpadding := int(src[length-1])
	if length < unpadding {
		return nil
	}
	return src[:(length - unpadding)]
}
