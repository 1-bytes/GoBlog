package middlewares

import (
	"GoBlog/pkg/auth"
	"GoBlog/pkg/flash"
	"net/http"
)

// Auth 登录用户才可以访问
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.Check() {
			flash.Warning("登录用户才能访问此页面")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next(w, r)
	}
}
