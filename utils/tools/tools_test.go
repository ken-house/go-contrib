package tools

import (
	"fmt"
	"testing"
)

func TestGenerateRandStr(t *testing.T) {
	res := GenerateRandStr(10, 3)
	fmt.Println(res)
}
