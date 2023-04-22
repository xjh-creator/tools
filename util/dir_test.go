package util

import (
	"fmt"
	"testing"
)

func TestInferRootDir(t *testing.T) {
	InferRootDir()

	fmt.Println(RootDir)
}
