package controllers

import (
	"GoBlog/pkg/view"
	"fmt"
	"net/http"
)

// PagesController 处理静态页面
type PagesController struct{}

// AboutFormData 关于页面表单数据
type AboutFormData struct {
	Title, Body string
}

// About 关于
func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	view.Render(w, AboutFormData{
		Title: "关于我们",
		Body:  "此博客是用以记录编程笔记，如您有反馈或建议，请联系 QQ:123456",
	}, "pages.about")
}

// Home 首页
func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello，这里是blog</h1>")
}

// NotFound 404页面
func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}
