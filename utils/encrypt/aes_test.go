package encrypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEncrypt(t *testing.T) {
	key := "1234567812345678" // 必须16位
	iv := "8765432187654321"  // 必须16位
	encryStr, err := AesEncrypt("您好，北京a", key, iv)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	fmt.Println(encryStr)
	assert.Equal(t, err, nil)
}

func TestAesDecrypt(t *testing.T) {
	key := "1234567812345678" // 必须16位
	iv := "8765432187654321"  // 必须16位
	contentStr, err := AesDecrypt("LsC6PCqBrcmcHPw+0m1ODwtG280fCMav9+CBtKdkv+g=", key, iv)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	fmt.Println(contentStr)
	assert.Equal(t, err, nil)
}
