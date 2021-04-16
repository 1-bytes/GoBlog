package controllers

import (
	"GoBlog/app/models/article"
	"GoBlog/pkg/logger"
	"GoBlog/pkg/route"
	"GoBlog/pkg/view"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"unicode/utf8"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// --- 4.读取数据成功，显示文章 ---
		view.Render(w, article, "articles.show")
	}
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, _ *http.Request) {
	// 1. 获取结果集
	articles, err := article.GetAll()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		// 2. 加载模板
		view.Render(w, articles, "articles.index")
	}
}

// validateArticleFormData 文章表单验证
func validateArticleFormData(title string, body string) map[string]string {
	errs := make(map[string]string)

	// 验证标题
	if title == "" {
		errs["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString("title") > 40 {
		errs["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errs["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errs["body"] = "内容长度不能大于或等于10个字节"
	}
	return errs
}

// Store 创建文章页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.PostFormValue("body")

	errs := validateArticleFormData(title, body)

	// 检查是否有错误
	if len(errs) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create() //nolint:errcheck

		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+_article.GetStringID())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Title":  title,
			"Body":   body,
			"Errors": errs,
		}, "articles.create", "articles._form_field")
	}
}

// Create 创建文章页面
func (*ArticlesController) Create(w http.ResponseWriter, _ *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Edit 更新文章页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2.读取对应的文章数据
	_article, err := article.Get(id)

	// 3.如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 读取成功，显示表单
		view.Render(w, view.D{
			"Title":   _article.Title,
			"Body":    _article.Body,
			"Article": _article,
		}, "articles.edit", "articles._form_field")
	}
}

// Update 更新文章
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	// 1.获取URL参数
	id := route.GetRouteVariable("id", r)

	// 2.读取对应的文章数据
	_article, err := article.Get(id)

	// 3.如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4.未出现错误，进行表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errs := validateArticleFormData(title, body)

		if len(errs) == 0 {
			// 表单验证通过，更新数据
			_article.Title = title
			_article.Body = body
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
				"Title":   title,
				"Body":    body,
				"Article": _article,
				"Errors":  errs,
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)
	// 2. 读取对应的文章数据
	_article, err := article.Get(id)
	// 3. 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
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
