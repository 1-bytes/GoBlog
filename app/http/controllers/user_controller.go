package controllers

import (
	"GoBlog/app/models/article"
	"GoBlog/app/models/user"
	"GoBlog/pkg/logger"
	"GoBlog/pkg/route"
	"GoBlog/pkg/view"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
}

func (*UserController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)
	// 2. 读取对应的文章数据
	_user, err := user.Get(id)
	// 3. 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 用户未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示文章列表
		articles, err := article.GetByUserID(_user.GetStringID())
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			view.Render(w, view.D{
				"Articles": articles,
			}, "articles.index", "articles._article_meta")
		}
	}
}
