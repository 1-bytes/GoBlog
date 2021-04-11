package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestAllPages(t *testing.T) {
	baseUrl := "http://localhost:3000"

	// 1. 声明加初始化测试数据
	var tests = []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", http.StatusOK},
		{"GET", "/about", http.StatusOK},
		{"GET", "/notfound", http.StatusNotFound},
		{"GET", "/articles", http.StatusOK},
		{"GET", "/articles/create", http.StatusOK},
		{"GET", "/articles/3", http.StatusOK},
		{"GET", "/articles/3/edit", http.StatusOK},
		{"POST", "/articles/3", http.StatusOK},
		{"POST", "/articles", http.StatusOK},
		{"POST", "/articles/1/delete", http.StatusNotFound},
	}

	// 2. 遍历所有测试
	for _, test := range tests {
		t.Logf("当前请求 URL: %v \n", test.url)

		var (
			resp *http.Response
			err  error
		)

		// 2.1 请求以获取响应
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseUrl+test.url, data)
		default:
			resp, err = http.Get(baseUrl + test.url)
		}
		// 2.2 断言
		assert.NoError(t, err, "请求 "+test.url+" 时报错")
		assert.Equal(t, test.expected, resp.StatusCode, test.url+" 应返回状态码 "+strconv.Itoa(resp.StatusCode))
	}
}
