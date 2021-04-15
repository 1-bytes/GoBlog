package view

import (
	"GoBlog/pkg/logger"
	"GoBlog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// Render 渲染视图
func Render(w io.Writer, data interface{}, tplFiles ...string) error {
	// 1. 设置模板相对路径
	viewDir := "resources/views/"

	// 2. 语法糖，将 articles.show 更正为 articles/show
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	// 3. 所有布局模板文件 Slice
	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// 4. 在 Slice 里新增我们的目标文件
	allFiles := append(layoutFiles, tplFiles...)

	// 5. 解析模板文件
	tmp1, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	// 6. 渲染模板
	return tmp1.ExecuteTemplate(w, "app", data)
}
