package middlewares

import "net/http"

// ForceHTML 强制标头返回 HTML 内容类型
func ForceHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置Header
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}
