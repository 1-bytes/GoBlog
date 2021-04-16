package controllers

import (
	"GoBlog/pkg/view"
	"net/http"
)

// AuthController 处理静态页面
type AuthController struct{}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, _ *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (*AuthController) DoRegister(_ http.ResponseWriter, _ *http.Request) {

}
