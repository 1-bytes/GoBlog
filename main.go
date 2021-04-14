package main

import (
	"GoBlog/bootstrap"
	"net/http"
	"strings"
)

// forceHTMLMiddleware 中间件,用于设置返回的Header中的ContentType.
func forceHTMLMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置Header
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		// 继续处理请求
		h.ServeHTTP(w, r)
	})
}

// removeTrailingSlash 中间件,路由清理末尾的斜杠.
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

// main ...
func main() {
	//database.Initialize()
	//db := database.DB

	bootstrap.SetupDB()
	router := bootstrap.SetupRoute()

	// 中间件：强制内容为 HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
