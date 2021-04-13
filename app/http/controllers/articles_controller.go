package controllers

import (
	"GoBlog/app/models/article"
	"GoBlog/pkg/logger"
	"GoBlog/pkg/route"
	"GoBlog/pkg/types"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

// Show 文章详情
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2.读取对应的文章数据
	article, err := article.Get(id)

	// 3.如果出现错误
	if err != nil {
		// 判断是没找到数据 还是查询报错了
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4.读取数据成功，显示表单
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Int64ToString": types.Int64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, article)
	}
}
