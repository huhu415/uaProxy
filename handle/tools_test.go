package handle

import (
	"testing"
)

func Test_isHTTP_Fail(t *testing.T) {
	// 失败测试用例 - 这些都不是有效的 HTTP 请求
	failTests := [][]byte{
		[]byte("INVALID /path HTTP/1.1"),
		[]byte("NOT-HTTP /index.html"),
		[]byte("HELLO WORLD"),
		[]byte("SSL-13"),
		[]byte("12345"),
		[]byte(""),
		[]byte(" "),
		[]byte("GEET /index.html"),
		[]byte("HTTP1.1 /index.html"), // 缺少空格
		[]byte("\r\n"),                // 只有换行
		[]byte("get /index.html"),     // 方法名小写
	}

	// 测试失败用例
	for _, tt := range failTests {
		methodLine := tt
		if len(tt) > PEEKSIZE {
			methodLine = tt[:PEEKSIZE]
		}
		t.Run("Should Fail: "+string(methodLine), func(t *testing.T) {
			if got := isHTTP(methodLine); got != false {
				t.Errorf("raw = %v, input = %v, want false for input: %s", tt, got, string(methodLine))
			}
		})
	}
}

func Test_isHTTP_Success(t *testing.T) {
	// 成功测试用例 - 这些都是有效的 HTTP 请求开头
	successTests := [][]byte{
		[]byte("PROPPATCH /api/data HTTP/1.1"),
		[]byte("SUBSCRIBE /api/data HTTP/1.1"),
		[]byte("UNSUBSCRIBE /api/data HTTP/1.1"),
		[]byte("CHECKOUT /api/data HTTP/1.1"),
		[]byte("CHECKOUT /api/data HTTP/1.1"),
		[]byte("GET /index.html HTTP/1.1"),
		[]byte("POST /api/data HTTP/1.1"),
		[]byte("HEAD /test HTTP/1.0"),
		[]byte("PUT /update HTTP/1.1"),
		[]byte("DELETE /remove HTTP/1.1"),
		[]byte("OPTIONS /check HTTP/1.1"),
		[]byte("PATCH /modify HTTP/1.1"),
	}

	// 测试成功用例
	for _, tt := range successTests {
		methodLine := tt[:PEEKSIZE]
		t.Run("Should Success: "+string(methodLine), func(t *testing.T) {
			if got := isHTTP(methodLine); got != true {
				t.Errorf("raw = %v, input = %v, want true for input: %s", tt, got, string(methodLine))
			}
		})
	}
}
