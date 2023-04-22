package util

import (
	"os"
	"path/filepath"
)

//
var RootDir string

// InferRootDir 推断出项目的根目录
func InferRootDir() {
	cwd, err := os.Getwd() // 获取当前工作目录
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		// 确保项目目录下存在 template 目录
		if exists(d + "/util") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)

	return err == nil || os.IsExist(err)
}
