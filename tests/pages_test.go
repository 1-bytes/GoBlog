package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHomePage(t *testing.T) {
	baseUrl := "http://localhost:3000"

	// 1. 请求 -> 模拟用户访问浏览器
	resp, err := http.Get(baseUrl)

	// 2. 检测 -> 是否无错误且返回状态 200
	assert.NoError(t, err, "有错误发生，err不为空")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "应返回状态码 200")
}
