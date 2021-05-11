package view

import (
	"GoBlog/pkg/auth"
	"GoBlog/pkg/flash"
	"GoBlog/pkg/logger"
	"GoBlog/pkg/route"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// D 是 map[string]interface{} 的简写
type D map[string]interface{}

// Render 渲染通用视图
func Render(w io.Writer, data D, tplFiles ...string) {
	RenderTemplate(w, "app", data, tplFiles...)
}

// RenderSimple 渲染简单视图
func RenderSimple(w io.Writer, data D, tplFiles ...string) {
	RenderTemplate(w, "simple", data, tplFiles...)
}

// RenderTemplate 渲染视图
func RenderTemplate(w io.Writer, name string, data D, tplFiles ...string) {
	// 1. 通用模板数据
	data["isLogined"] = auth.Check()
	data["loginUser"] = auth.User()
	data["flash"] = flash.All()

	// 2. 生成模板文件
	allFiles := getTemplateFiles(tplFiles...)

	// 3. 解析所有模板文件
	tmp1, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	// 4. 渲染模板
	tmp1.ExecuteTemplate(w, name, data)
}

// RenderTemplate 获取模板文件的内容
func getTemplateFiles(tplFiles ...string) []string {
	// 1. 设置模板相对路径
	viewDir := "resources/views/"

	// 2. 遍历传参文件列表 Slice，设置正确的路径，支持 dir.filename 语法糖
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	// 3. 所有布局模板文件 Slice
	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// 4. 合并所有文件
	return append(layoutFiles, tplFiles...)
}

// RenderErrorTemplate 渲染错误视图
func RenderSimpleMessage(w http.ResponseWriter, errMsg interface{}, httpStatus int) {
	w.WriteHeader(httpStatus)
	RenderSimple(w, D{
		"Title": "错误",
		"Body":  errMsg,
	}, "pages.message")
}
