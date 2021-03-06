package controllers

import (
	"GoBlog/app/models/article"
	"GoBlog/app/policies"
	"GoBlog/app/requests"
	"GoBlog/pkg/auth"
	"GoBlog/pkg/config"
	"GoBlog/pkg/route"
	"GoBlog/pkg/view"
	"fmt"
	"net/http"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
	BaseController
}

// Show 文章详情
func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2.读取对应的文章数据
	article, err := article.Get(id)

	// 3.如果出现错误
	if err != nil {
		// 判断是没找到数据 还是查询报错了
		ac.ResponseForSQLError(w, err)
	} else {
		// --- 4.读取数据成功，显示文章 ---
		view.Render(w, view.D{
			"Article":          article,
			"CanModifyArticle": policies.CanModifyArticle(article),
		}, "articles.show", "articles._article_meta")
	}
}

// Index 文章列表页
func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// 1. 获取结果集
	articles, pagerData, err := article.GetAll(r, config.GetInt("pagination.perpage"))

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 2. 加载模板
		view.Render(w, view.D{
			"Articles":  articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}

// Store 创建文章页面
func (ac *ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	currentUser := auth.User()
	_article := article.Article{
		Title:  r.FormValue("title"),
		Body:   r.PostFormValue("body"),
		UserID: currentUser.ID,
	}

	// 2. 表单验证
	errs := requests.ValidateArticleForm(_article)

	// 3. 检查是否有错误
	if len(errs) == 0 {
		// 创建文章
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+_article.GetStringID())
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errs,
		}, "articles.create", "articles._form_field")
	}
}

// Create 创建文章页面
func (ac *ArticlesController) Create(w http.ResponseWriter, _ *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Edit 更新文章页面
func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2.读取对应的文章数据
	_article, err := article.Get(id)

	// 3.如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 4. 读取成功，显示表单
			view.Render(w, view.D{
				"Article": _article,
				"Errors":  view.D{},
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Update 更新文章
func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1.获取URL参数
	id := route.GetRouteVariable("id", r)

	// 2.读取对应的文章数据
	_article, err := article.Get(id)

	// 3.如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		// 4.未出现错误
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 4.1 表单验证
			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")

			errs := requests.ValidateArticleForm(_article)

			if len(errs) == 0 {
				// 表单验证通过，更新数据
				rowsAffected, err := _article.Update()

				if err != nil {
					// 数据库错误
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "500 服务器内部错误")
					return
				}

				// 更新成功，跳转到文章详情页
				if rowsAffected > 0 {
					showURL := route.Name2URL("articles.show", "id", id)
					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					fmt.Fprint(w, "您没有做任何更改！")
				}
			} else {
				// 4.3表单验证不通过，验证路由
				view.Render(w, view.D{
					"Article": _article,
					"Errors":  errs,
				}, "articles.edit", "articles._form_field")
			}
		}
	}
}

// Delete 删除文章
func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)
	// 2. 读取对应的文章数据
	_article, err := article.Get(id)
	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			// 4. 未出现错误，执行删除操作
			rowsAffected, err := _article.Delete()
			// 4.1 发生错误
			if err != nil {
				// SQL报错
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 服务器内部错误")
			} else {
				// 4.2 未发生错误
				if rowsAffected > 0 {
					// 重定向到文章列表页
					indexURL := route.Name2URL("articles.index")
					http.Redirect(w, r, indexURL, http.StatusFound)
				} else {
					// Edge case
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprintf(w, "404 文章未找到")
				}
			}
		}
	}
}
