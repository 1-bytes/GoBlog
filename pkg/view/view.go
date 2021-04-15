package view

import (
	"GoBlog/pkg/logger"
	"GoBlog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

func Render(w io.Writer, name string, data interface{}) error {
	// 1. 设置模板相对路径
	viewDir := "resources/views/"

	// 2. 语法糖，将 articles.show 更正为 articles/show
	name = strings.Replace(name, ".", "/", -1)

	// 3. 所有布局模板文件 Slice
	files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
	logger.LogError(err)

	// 4. 在 Slice 里新增我们的目标文件
	newFiles := append(files, viewDir+"/articles/show.gohtml")

	// 5. 解析模板文件
	tmp1, err := template.New("show.gohtml").Funcs(template.FuncMap{
		"RouteName2URL": route.Name2URL,
	}).ParseFiles(newFiles...)
	logger.LogError(err)

	// 6. 渲染模板
	return tmp1.ExecuteTemplate(w, "app", data)
}
