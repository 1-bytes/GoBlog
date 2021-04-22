package middlewares

import (
	"GoBlog/pkg/session"
	"net/http"
)

func StartSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 启动会话
		session.StartSession(w, r)

		// 继续处理接下来的请求
		next.ServeHTTP(w, r)
	})
}
